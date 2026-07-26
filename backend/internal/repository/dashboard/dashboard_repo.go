package dashboard

import (
	"fmt"
	"time"

	dashDTO "github.com/ahmadammarm/amami/backend/internal/dto/dashboard"
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	// Existing scalar metrics
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
	// New analytics (chart data)
	GetMonthlyFinanceTrend(months int) ([]dashDTO.MonthlyFinancialPoint, error)
	GetAssetStatusBreakdown() (dashDTO.AssetStatusBreakdown, error)
	GetQurbanTypeBreakdown() (dashDTO.QurbanTypeBreakdown, error)
	GetMonthlyJamaahGrowth(months int) ([]dashDTO.MonthlyJamaahPoint, error)
	GetZakatDonationVsDistribution() (dashDTO.ZakatBreakdownMetrics, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// --- Existing Scalar Methods ---

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
	var total int64
	err := r.db.Model(&domain.Transaction{}).
		Where("category ILIKE ?", "%zakat%").
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetTotalMustahik() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Jamaah{}).Where("is_mustahik = ?", true).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetQurbanAnimalsCount(animalType string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.QurbanAnimal{}).Where("UPPER(type) = UPPER(?)", animalType).Count(&count).Error
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
	err := r.db.Model(&domain.Agenda{}).Where("status = ?", "SCHEDULED").Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetActiveLoansCount() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Asset{}).Where("current_status = ?", "LOANED").Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetBrokenAssetsCount() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Asset{}).Where("current_status IN (?, ?, ?)", "REPAIR", "NEEDS_REPAIR", "BROKEN").Count(&count).Error
	return count, err
}

func (r *dashboardRepository) PingDB() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// --- New Analytics Methods (Optimized Single-Query) ---

// gormTable returns the actual table name GORM uses for a given model.
func (r *dashboardRepository) gormTable(model interface{}) string {
	stmt := &gorm.Statement{DB: r.db}
	stmt.Parse(model)
	return stmt.Schema.Table
}

func (r *dashboardRepository) GetMonthlyFinanceTrend(months int) ([]dashDTO.MonthlyFinancialPoint, error) {
	type row struct {
		Month   string
		Income  int64
		Expense int64
	}
	var rows []row
	intervalExpr := fmt.Sprintf("%d months", months)
	table := r.gormTable(&domain.Transaction{})
	sqlQuery := `
		SELECT
			TO_CHAR(DATE_TRUNC('month', created_at), 'Mon YYYY') AS month,
			COALESCE(SUM(CASE WHEN type = 'CREDIT' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'DEBIT'  THEN amount ELSE 0 END), 0) AS expense
		FROM ` + table + `
		WHERE created_at >= NOW() - INTERVAL '` + intervalExpr + `'
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY DATE_TRUNC('month', created_at) ASC
	`
	if err := r.db.Raw(sqlQuery).Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]dashDTO.MonthlyFinancialPoint, len(rows))
	for i, r := range rows {
		result[i] = dashDTO.MonthlyFinancialPoint{
			Month:   r.Month,
			Income:  r.Income,
			Expense: r.Expense,
		}
	}
	return result, nil
}

func (r *dashboardRepository) GetAssetStatusBreakdown() (dashDTO.AssetStatusBreakdown, error) {
	type row struct {
		Good        int64
		Loaned      int64
		NeedsRepair int64
		Broken      int64
	}
	var result row
	table := r.gormTable(&domain.Asset{})
	sqlQuery := `
		SELECT
			COUNT(*) FILTER (WHERE current_status = 'GOOD')                          AS good,
			COUNT(*) FILTER (WHERE current_status = 'LOANED')                       AS loaned,
			COUNT(*) FILTER (WHERE current_status IN ('REPAIR', 'NEEDS_REPAIR'))     AS needs_repair,
			COUNT(*) FILTER (WHERE current_status = 'BROKEN')                       AS broken
		FROM ` + table + `
	`
	if err := r.db.Raw(sqlQuery).Scan(&result).Error; err != nil {
		return dashDTO.AssetStatusBreakdown{}, err
	}
	return dashDTO.AssetStatusBreakdown{
		Good:        result.Good,
		Loaned:      result.Loaned,
		NeedsRepair: result.NeedsRepair,
		Broken:      result.Broken,
	}, nil
}

func (r *dashboardRepository) GetQurbanTypeBreakdown() (dashDTO.QurbanTypeBreakdown, error) {
	type row struct {
		Type  string
		Count int64
	}
	var rows []row
	table := r.gormTable(&domain.QurbanAnimal{})
	sqlQuery := `
		SELECT UPPER(type) AS type, COUNT(*) AS count
		FROM ` + table + `
		GROUP BY UPPER(type)
	`
	if err := r.db.Raw(sqlQuery).Scan(&rows).Error; err != nil {
		return dashDTO.QurbanTypeBreakdown{}, err
	}
	var breakdown dashDTO.QurbanTypeBreakdown
	for _, row := range rows {
		switch row.Type {
		case "SAPI":
			breakdown.Sapi = row.Count
		case "KAMBING":
			breakdown.Kambing = row.Count
		case "DOMBA":
			breakdown.Domba = row.Count
		}
	}
	return breakdown, nil
}

func (r *dashboardRepository) GetMonthlyJamaahGrowth(months int) ([]dashDTO.MonthlyJamaahPoint, error) {
	type row struct {
		Month      string
		NewMembers int64
	}
	var rows []row
	intervalExpr := fmt.Sprintf("%d months", months)
	table := r.gormTable(&domain.Jamaah{})
	sqlQuery := `
		SELECT
			TO_CHAR(DATE_TRUNC('month', created_at), 'Mon YYYY') AS month,
			COUNT(*) AS new_members
		FROM ` + table + `
		WHERE created_at >= NOW() - INTERVAL '` + intervalExpr + `'
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY DATE_TRUNC('month', created_at) ASC
	`
	if err := r.db.Raw(sqlQuery).Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]dashDTO.MonthlyJamaahPoint, len(rows))
	for i, r := range rows {
		result[i] = dashDTO.MonthlyJamaahPoint{
			Month:      r.Month,
			NewMembers: r.NewMembers,
		}
	}
	return result, nil
}

func (r *dashboardRepository) GetZakatDonationVsDistribution() (dashDTO.ZakatBreakdownMetrics, error) {
	type result struct {
		TotalDonated     int64
		TotalDistributed int64
	}
	var res result
	donationTable := r.gormTable(&domain.ZakatDonation{})
	distributionTable := r.gormTable(&domain.ZakatDistribution{})
	sqlQuery := `
		SELECT
			COALESCE((SELECT CAST(SUM(amount) AS BIGINT) FROM ` + donationTable + `), 0) AS total_donated,
			COALESCE((SELECT CAST(SUM(amount) AS BIGINT) FROM ` + distributionTable + `), 0) AS total_distributed
	`
	if err := r.db.Raw(sqlQuery).Scan(&res).Error; err != nil {
		return dashDTO.ZakatBreakdownMetrics{}, err
	}
	return dashDTO.ZakatBreakdownMetrics{
		TotalDonated:     res.TotalDonated,
		TotalDistributed: res.TotalDistributed,
	}, nil
}
