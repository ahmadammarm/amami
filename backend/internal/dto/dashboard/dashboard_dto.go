package dashboard

type DashboardMetricsResponse struct {
	Finance   FinanceMetrics   `json:"finance"`
	Zakat     ZakatMetrics     `json:"zakat"`
	Qurban    QurbanMetrics    `json:"qurban"`
	Jamaah    JamaahMetrics    `json:"jamaah"`
	Logistics LogisticsMetrics `json:"logistics"`
	System    SystemMetrics    `json:"system"`
	Analytics AnalyticsMetrics `json:"analytics"`
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

// --- Analytics (Chart Data) ---

type MonthlyFinancialPoint struct {
	Month   string `json:"month"`
	Income  int64  `json:"income"`
	Expense int64  `json:"expense"`
}

type AssetStatusBreakdown struct {
	Good        int64 `json:"good"`
	Loaned      int64 `json:"loaned"`
	NeedsRepair int64 `json:"needs_repair"`
	Broken      int64 `json:"broken"`
}

type MonthlyJamaahPoint struct {
	Month      string `json:"month"`
	NewMembers int64  `json:"new_members"`
}

type ZakatBreakdownMetrics struct {
	TotalDonated     int64 `json:"total_donated"`
	TotalDistributed int64 `json:"total_distributed"`
}

type QurbanTypeBreakdown struct {
	Sapi   int64 `json:"sapi"`
	Kambing int64 `json:"kambing"`
	Domba  int64 `json:"domba"`
}

type AnalyticsMetrics struct {
	MonthlyFinance  []MonthlyFinancialPoint `json:"monthly_finance"`
	AssetBreakdown  AssetStatusBreakdown    `json:"asset_breakdown"`
	QurbanBreakdown QurbanTypeBreakdown     `json:"qurban_breakdown"`
	JamaahGrowth    []MonthlyJamaahPoint    `json:"jamaah_growth"`
	ZakatBreakdown  ZakatBreakdownMetrics   `json:"zakat_breakdown"`
}
