package dashboard

import (
	"github.com/ahmadammarm/amami/backend/internal/dto/dashboard"
	dashRepo "github.com/ahmadammarm/amami/backend/internal/repository/dashboard"
)

type DashboardService interface {
	GetMetrics() (*dashboard.DashboardMetricsResponse, error)
}

type dashboardService struct {
	repo dashRepo.DashboardRepository
}

func NewDashboardService(repo dashRepo.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetMetrics() (*dashboard.DashboardMetricsResponse, error) {
	var metrics dashboard.DashboardMetricsResponse

	// Finance
	balance, _ := s.repo.GetTotalFundBalance()
	inflow, _ := s.repo.GetMonthlyInflow()
	metrics.Finance = dashboard.FinanceMetrics{
		TotalBalance:  balance,
		MonthlyInflow: inflow,
	}

	// Zakat
	zakatCol, _ := s.repo.GetTotalZakatCollected()
	mustahik, _ := s.repo.GetTotalMustahik()
	metrics.Zakat = dashboard.ZakatMetrics{
		Collected: zakatCol,
		Mustahik:  mustahik,
	}

	// Qurban
	cows, _ := s.repo.GetQurbanAnimalsCount("Sapi")
	goats, _ := s.repo.GetQurbanAnimalsCount("Kambing")
	metrics.Qurban = dashboard.QurbanMetrics{
		Cows:  cows,
		Goats: goats,
	}

	// Jamaah
	totalJamaah, _ := s.repo.GetTotalJamaah()
	newJamaah, _ := s.repo.GetNewlyRegisteredJamaah()
	metrics.Jamaah = dashboard.JamaahMetrics{
		Total:           totalJamaah,
		NewlyRegistered: newJamaah,
	}

	// Logistics
	agendas, _ := s.repo.GetUpcomingAgendasCount()
	loans, _ := s.repo.GetActiveLoansCount()
	broken, _ := s.repo.GetBrokenAssetsCount()
	metrics.Logistics = dashboard.LogisticsMetrics{
		UpcomingAgendas: agendas,
		ActiveLoans:     loans,
		BrokenAssets:    broken,
	}

	// System Health
	dbStatus := "Healthy"
	if err := s.repo.PingDB(); err != nil {
		dbStatus = "Down"
	}
	metrics.System = dashboard.SystemMetrics{
		Uptime:         "99.9%", // Simulated
		DBStatus:       dbStatus,
		ActiveSessions: 1, // Simulated
	}

	return &metrics, nil
}
