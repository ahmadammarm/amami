package dashboard

type DashboardMetricsResponse struct {
	Finance   FinanceMetrics   `json:"finance"`
	Zakat     ZakatMetrics     `json:"zakat"`
	Qurban    QurbanMetrics    `json:"qurban"`
	Jamaah    JamaahMetrics    `json:"jamaah"`
	Logistics LogisticsMetrics `json:"logistics"`
	System    SystemMetrics    `json:"system"`
}

type FinanceMetrics struct {
	TotalBalance  int64 `json:"total_balance"`
	MonthlyInflow int64 `json:"monthly_inflow"`
}

type ZakatMetrics struct {
	Collected int64 `json:"collected"`
	Mustahik  int64 `json:"mustahik"`
}

type QurbanMetrics struct {
	Cows  int64 `json:"cows"`
	Goats int64 `json:"goats"`
}

type JamaahMetrics struct {
	Total           int64 `json:"total"`
	NewlyRegistered int64 `json:"newly_registered"`
}

type LogisticsMetrics struct {
	UpcomingAgendas int64 `json:"upcoming_agendas"`
	ActiveLoans     int64 `json:"active_loans"`
	BrokenAssets    int64 `json:"broken_assets"`
}

type SystemMetrics struct {
	Uptime         string `json:"uptime"`
	DBStatus       string `json:"db_status"`
	ActiveSessions int64  `json:"active_sessions"`
}
