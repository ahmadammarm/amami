package seed

import (
	"log"
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"gorm.io/gorm"
)

// Seed inserts the super admin, default settings and sample data.
// Roles and permissions are synced separately by SyncRBAC on every start.
func Seed(db *gorm.DB) error {
	log.Println("Seeding database...")

	// 1. Roles & permissions
	if err := SyncRBAC(db); err != nil {
		return err
	}

	var superAdminRole domain.Role
	if err := db.Where("name = ?", domain.RoleSuperAdmin).First(&superAdminRole).Error; err != nil {
		return err
	}

	// 2. Seed Super Admin User
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}
	superAdmin := domain.User{
		Username:     "admin",
		Email:        "admin@amami.org",
		PasswordHash: hashedPassword,
		RoleID:       superAdminRole.ID,
		Status:       domain.UserActive,
	}

	if err := db.FirstOrCreate(&superAdmin, domain.User{Username: "admin"}).Error; err != nil {
		return err
	}

	// 3. Seed System Settings
	defaultSettings := []domain.SystemSetting{
		{Key: "mosque_name", Value: "Masjid Agung Al-Hikmah"},
		{Key: "mosque_address", Value: "Jl. Merdeka No. 1, Kota Amanah"},
		{Key: "mosque_phone", Value: "021-1234567"},
		{Key: "mosque_logo", Value: "/icons.svg"},
		{Key: "legal_yayasan_id", Value: "YAYASAN-AH-001"},
		{Key: "smtp_host", Value: "smtp.example.com"},
		{Key: "smtp_port", Value: "587"},
		{Key: "smtp_user", Value: "admin@amami.org"},
		{Key: "smtp_from", Value: "no-reply@amami.org"},
	}

	for _, s := range defaultSettings {
		if err := db.FirstOrCreate(&s, domain.SystemSetting{Key: s.Key}).Error; err != nil {
			return err
		}
	}

	// 4. Seed Default Funds
	defaultFunds := []domain.Fund{
		{Name: "Kas Operasional Masjid", Code: "KAS-OPR", CurrentBalance: 0},
		{Name: "Kas Pembangunan", Code: "KAS-PEM", CurrentBalance: 0},
		{Name: "Kas Anak Yatim", Code: "KAS-YAT", CurrentBalance: 0},
	}
	for i := range defaultFunds {
		if err := db.FirstOrCreate(&defaultFunds[i], domain.Fund{Code: defaultFunds[i].Code}).Error; err != nil {
			return err
		}
	}

	// 5. Seed Jamaah Data
	log.Println("Seeding Jamaah...")
	jamaahNames := []string{
		"Budi Santoso", "Ahmad Yusuf", "Siti Aminah", "Andi Pratama", "Nisa Kamilia",
		"Wahyu Hidayat", "Ratna Sari", "Irfan Hakim", "Dian Sastrowardoyo", "Muhammad Fadel",
		"Reza Rahadian", "Ayu Ting Ting", "Rafi Ahmad", "Baim Wong", "Deddy Corbuzier",
	}
	var jamaahs []domain.Jamaah
	for _, name := range jamaahNames {
		j := domain.Jamaah{
			FullName:      name,
			Phone:         "0812345678" + string(rune('0'+(len(jamaahs)%10))),
			Address:       "Jl. Kebaikan No. " + string(rune('1'+(len(jamaahs)%9))) + ", Jakarta",
			IsMustahik:    len(jamaahs)%5 == 0,
			MustahikScore: float64(len(jamaahs) % 5 * 20),
		}
		db.Create(&j)
		jamaahs = append(jamaahs, j)
	}

	// 6. Seed Transactions and Zakat
	log.Println("Seeding Transactions and Zakat...")
	categories := []string{"INFAQ_JUMAT", "INFAQ_KOTAK", "ZAKAT_FITRAH", "ZAKAT_MAAL"}
	expenseCategories := []string{"OPERASIONAL", "PEMELIHARAAN", "KEGIATAN"}

	// Seed last 6 months data
	now := time.Now()
	for m := 0; m < 6; m++ {
		monthDate := now.AddDate(0, -m, 0)

		// Incomes (4 per month)
		for i := 0; i < 4; i++ {
			amount := int64(500000 + (i * 250000))
			txDate := monthDate.AddDate(0, 0, -(i * 5))
			cat := categories[i%len(categories)]

			tx := domain.Transaction{
				FundID:    defaultFunds[0].ID,
				Type:      domain.TransactionCredit,
				Amount:    amount,
				Category:  cat,
				CreatedBy: superAdmin.ID,
				CreatedAt: txDate,
			}
			db.Create(&tx)

			// If Zakat, create Zakat Donation record
			if cat == "ZAKAT_FITRAH" || cat == "ZAKAT_MAAL" {
				zakat := domain.ZakatDonation{
					TransactionID: &tx.ID,
					MuzakkiID:     jamaahs[i%len(jamaahs)].ID,
					ZakatType:     cat,
					AmountOrQty:   float64(amount),
					Unit:          domain.ZakatUnitIDR,
					CreatedAt:     txDate,
				}
				db.Create(&zakat)
			}
			defaultFunds[0].CurrentBalance += amount
		}

		// Expenses (2 per month)
		for i := 0; i < 2; i++ {
			amount := int64(200000 + (i * 150000))
			txDate := monthDate.AddDate(0, 0, -(i * 10))

			tx := domain.Transaction{
				FundID:    defaultFunds[0].ID,
				Type:      domain.TransactionDebit,
				Amount:    amount,
				Category:  expenseCategories[i%len(expenseCategories)],
				CreatedBy: superAdmin.ID,
				CreatedAt: txDate,
			}
			db.Create(&tx)
			defaultFunds[0].CurrentBalance -= amount
		}
	}
	db.Save(&defaultFunds[0])

	// 7. Seed Qurban Packages & Bookings
	log.Println("Seeding Qurban...")
	packages := []domain.QurbanPackage{
		{Name: "Kambing Premium", Type: "KAMBING", Price: 3500000, YearHijri: 1447, StockTotal: 50, StockRemaining: 40},
		{Name: "Kambing Standar", Type: "KAMBING", Price: 2800000, YearHijri: 1447, StockTotal: 100, StockRemaining: 75},
		{Name: "Sapi Kolektif (1/7)", Type: "SAPI_KOLEKTIF", Price: 3000000, YearHijri: 1447, StockTotal: 70, StockRemaining: 60},
	}
	for i := range packages {
		db.Create(&packages[i])
	}

	for i := 0; i < 8; i++ {
		booking := domain.QurbanBooking{
			ShohibulID:    jamaahs[i].ID,
			PackageID:     packages[i%3].ID,
			PaymentStatus: "PAID",
			TotalAmount:   packages[i%3].Price,
			BookingDate:   now.AddDate(0, 0, -i),
		}
		db.Create(&booking)
	}

	// 8. Seed Inventory Assets
	log.Println("Seeding Inventory...")
	assets := []domain.Asset{
		{Name: "AC Daikin 2 PK", SKU: "AST-AC-001", PurchasePrice: 6000000, CurrentStatus: domain.AssetGood, Location: "Ruang Utama"},
		{Name: "AC Daikin 2 PK", SKU: "AST-AC-002", PurchasePrice: 6000000, CurrentStatus: domain.AssetGood, Location: "Ruang Utama"},
		{Name: "Sound System Yamaha", SKU: "AST-SND-001", PurchasePrice: 15000000, CurrentStatus: domain.AssetGood, Location: "Ruang Kontrol"},
		{Name: "Karpet Turki 50m", SKU: "AST-KRP-001", PurchasePrice: 25000000, CurrentStatus: domain.AssetNeedsRepair, Location: "Ruang Utama"},
		{Name: "Vacuum Cleaner", SKU: "AST-VAC-001", PurchasePrice: 2000000, CurrentStatus: domain.AssetGood, Location: "Gudang"},
		{Name: "Genset Honda 5000W", SKU: "AST-GEN-001", PurchasePrice: 12000000, CurrentStatus: domain.AssetGood, Location: "Halaman Samping"},
		{Name: "Kamera Sony A6400", SKU: "AST-CAM-001", PurchasePrice: 13000000, CurrentStatus: domain.AssetLoaned, Location: "Ruang Multimedia"},
		{Name: "Proyektor Epson", SKU: "AST-PRJ-001", PurchasePrice: 5500000, CurrentStatus: domain.AssetGood, Location: "Ruang Serbaguna"},
		{Name: "Papan Tulis Kaca", SKU: "AST-BRD-001", PurchasePrice: 1500000, CurrentStatus: domain.AssetGood, Location: "Ruang Serbaguna"},
		{Name: "Mimbar Jati", SKU: "AST-MBR-001", PurchasePrice: 8000000, CurrentStatus: domain.AssetGood, Location: "Ruang Utama"},
	}
	for i := range assets {
		db.Create(&assets[i])
	}

	// 9. Seed Audit Logs
	for i := 1; i <= 25; i++ {
		action := "UPDATE"
		if i%3 == 0 {
			action = "LOGIN"
		}
		if i%5 == 0 {
			action = "CREATE"
		}

		log := domain.AuditLog{
			UserID:    superAdmin.ID,
			Action:    action,
			Entity:    "System Configuration",
			CreatedAt: time.Now().Add(time.Duration(-i) * time.Hour),
		}
		db.Create(&log)
	}

	log.Println("Seeding completed successfully!")
	return nil
}
