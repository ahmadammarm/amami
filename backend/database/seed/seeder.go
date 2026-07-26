package seed

import (
	"log"
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	log.Println("Seeding database...")

	// 1. Seed Permissions
	permissions := []domain.Permission{
		{Code: "user:manage", Description: "Manage users and roles"},
		{Code: "ledger:read", Description: "Read financial reports"},
		{Code: "ledger:write", Description: "Create transactions"},
		{Code: "zakat:manage", Description: "Manage zakat collection and distribution"},
		{Code: "qurban:manage", Description: "Manage qurban packages and bookings"},
		{Code: "inventory:manage", Description: "Manage mosque assets"},
		{Code: "jamaah:manage", Description: "Manage community database and scoring"},
		{Code: "agenda:manage", Description: "Manage mosque events and documentation"},
		{Code: "settings:manage", Description: "Manage mosque profile and system settings"},
	}

	for _, p := range permissions {
		if err := db.FirstOrCreate(&p, domain.Permission{Code: p.Code}).Error; err != nil {
			return err
		}
	}

	// 2. Seed Roles
	roles := []domain.Role{
		{Name: "SUPER_ADMIN", Description: "Full system governance and management access"},
		{Name: "BENDAHARA", Description: "Financial and Zakat management access"},
		{Name: "TAKMIR", Description: "Broad operational management access (All except User/System management)"},
		{Name: "SEKRETARIS", Description: "Administrative and documentation access"},
		{Name: "JAMAAH", Description: "Community member access (Personal data only)"},
	}

	for i := range roles {
		if err := db.FirstOrCreate(&roles[i], domain.Role{Name: roles[i].Name}).Error; err != nil {
			return err
		}
	}

	// 3. Map Permissions to Roles (RBAC)
	
	// Helper to find permission ID by code
	getPermID := func(code string) uint {
		var p domain.Permission
		db.Where("code = ?", code).First(&p)
		return p.ID
	}

	// SUPER_ADMIN: Gets everything
	var allPerms []domain.Permission
	db.Find(&allPerms)
	for _, p := range allPerms {
		db.FirstOrCreate(&domain.RolePermission{RoleID: roles[0].ID, PermissionID: p.ID})
	}

	// BENDAHARA: Finance, Zakat, and reading Jamaah
	bendaharaPerms := []string{"ledger:read", "ledger:write", "zakat:manage", "jamaah:manage"}
	for _, code := range bendaharaPerms {
		db.FirstOrCreate(&domain.RolePermission{RoleID: roles[1].ID, PermissionID: getPermID(code)})
	}

	// TAKMIR: Broad access (Everything except user:manage and settings:manage)
	takmirPerms := []string{"ledger:read", "zakat:manage", "qurban:manage", "inventory:manage", "jamaah:manage", "agenda:manage"}
	for _, code := range takmirPerms {
		db.FirstOrCreate(&domain.RolePermission{RoleID: roles[2].ID, PermissionID: getPermID(code)})
	}

	// SEKRETARIS: Administration, Agenda, and reading reports
	sekretarisPerms := []string{"ledger:read", "jamaah:manage", "agenda:manage", "settings:manage"}
	for _, code := range sekretarisPerms {
		db.FirstOrCreate(&domain.RolePermission{RoleID: roles[3].ID, PermissionID: getPermID(code)})
	}

	// 4. Seed Super Admin User
	hashedPassword, _ := utils.HashPassword("admin123")
	superAdmin := domain.User{
		Username:     "admin",
		Email:        "admin@amami.org",
		PasswordHash: hashedPassword,
		RoleID:       roles[0].ID,
		Status:       "ACTIVE",
	}

	if err := db.FirstOrCreate(&superAdmin, domain.User{Username: "admin"}).Error; err != nil {
		return err
	}

	// 5. Seed System Settings
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

	// 5.5 Seed Default Funds
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

	// 6. Seed Jamaah Data
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

	// 7. Seed Transactions and Zakat
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
			txDate := monthDate.AddDate(0, 0, -(i*5))
			cat := categories[i%len(categories)]
			
			tx := domain.Transaction{
				FundID:    defaultFunds[0].ID,
				Type:      "CREDIT",
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
					Unit:          "IDR",
					CreatedAt:     txDate,
				}
				db.Create(&zakat)
			}
			defaultFunds[0].CurrentBalance += amount
		}

		// Expenses (2 per month)
		for i := 0; i < 2; i++ {
			amount := int64(200000 + (i * 150000))
			txDate := monthDate.AddDate(0, 0, -(i*10))
			
			tx := domain.Transaction{
				FundID:    defaultFunds[0].ID,
				Type:      "DEBIT",
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

	// 8. Seed Qurban Packages & Bookings
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

	// 9. Seed Inventory Assets
	log.Println("Seeding Inventory...")
	assets := []domain.Asset{
		{Name: "AC Daikin 2 PK", SKU: "AST-AC-001", PurchasePrice: 6000000, CurrentStatus: "GOOD", Location: "Ruang Utama"},
		{Name: "AC Daikin 2 PK", SKU: "AST-AC-002", PurchasePrice: 6000000, CurrentStatus: "GOOD", Location: "Ruang Utama"},
		{Name: "Sound System Yamaha", SKU: "AST-SND-001", PurchasePrice: 15000000, CurrentStatus: "GOOD", Location: "Ruang Kontrol"},
		{Name: "Karpet Turki 50m", SKU: "AST-KRP-001", PurchasePrice: 25000000, CurrentStatus: "NEEDS_REPAIR", Location: "Ruang Utama"},
		{Name: "Vacuum Cleaner", SKU: "AST-VAC-001", PurchasePrice: 2000000, CurrentStatus: "GOOD", Location: "Gudang"},
		{Name: "Genset Honda 5000W", SKU: "AST-GEN-001", PurchasePrice: 12000000, CurrentStatus: "GOOD", Location: "Halaman Samping"},
		{Name: "Kamera Sony A6400", SKU: "AST-CAM-001", PurchasePrice: 13000000, CurrentStatus: "LOANED", Location: "Ruang Multimedia"},
		{Name: "Proyektor Epson", SKU: "AST-PRJ-001", PurchasePrice: 5500000, CurrentStatus: "GOOD", Location: "Ruang Serbaguna"},
		{Name: "Papan Tulis Kaca", SKU: "AST-BRD-001", PurchasePrice: 1500000, CurrentStatus: "GOOD", Location: "Ruang Serbaguna"},
		{Name: "Mimbar Jati", SKU: "AST-MBR-001", PurchasePrice: 8000000, CurrentStatus: "GOOD", Location: "Ruang Utama"},
	}
	for i := range assets {
		db.Create(&assets[i])
	}

	// 10. Seed Audit Logs
	for i := 1; i <= 25; i++ {
		action := "UPDATE"
		if i%3 == 0 { action = "LOGIN" }
		if i%5 == 0 { action = "CREATE" }
		
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
