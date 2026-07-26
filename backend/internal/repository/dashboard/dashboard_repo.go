package dashboard

import (
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetTotalFundBalance() (int64, error)
	GetMonthlyInflow() (int64, error)
	GetTotalZakatCollected() (int64, error)
	GetTotalMustahik() (int64, error)
	GetQurbanAnimalsCount(animalType string) (int64, error)
	GetTotalJamaah() (int64, error)
	GetNewlyRegisteredJamaah() (int64, error)
	GetUpcomingAgendasCount() (int64, error)
	GetActiveLoansCount() (int64, error)
	GetBrokenAssetsCount() (int64, error)
	PingDB() error
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetTotalFundBalance() (int64, error) {
	var total int64
	err := r.db.Model(&domain.Fund{}).Select("COALESCE(SUM(current_balance), 0)").Row().Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetMonthlyInflow() (int64, error) {
	var total int64
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	err := r.db.Model(&domain.Transaction{}).
		Where("type = ? AND created_at >= ?", "CREDIT", startOfMonth).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetTotalZakatCollected() (int64, error) {
	// For now, checking Zakat transactions based on category or we can just mock it if table is empty
	var total int64
	err := r.db.Model(&domain.Transaction{}).
		Where("category ILIKE ?", "%zakat%").
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetTotalMustahik() (int64, error) {
	var count int64
	err := r.db.Model(&domain.MustahikData{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetQurbanAnimalsCount(animalType string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.QurbanAnimal{}).Where("type = ?", animalType).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalJamaah() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Jamaah{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetNewlyRegisteredJamaah() (int64, error) {
	var count int64
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	err := r.db.Model(&domain.Jamaah{}).Where("created_at >= ?", startOfMonth).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetUpcomingAgendasCount() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Agenda{}).Where("start_time >= ?", time.Now()).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetActiveLoansCount() (int64, error) {
	var count int64
	err := r.db.Model(&domain.AssetLoan{}).Where("return_date IS NULL").Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetBrokenAssetsCount() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Asset{}).Where("current_status = ?", "REPAIR").Count(&count).Error
	return count, err
}

func (r *dashboardRepository) PingDB() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
