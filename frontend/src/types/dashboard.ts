export interface FinanceMetrics {
  total_balance: number;
  monthly_inflow: number;
}

export interface ZakatMetrics {
  collected: number;
  mustahik: number;
}

export interface QurbanMetrics {
  cows: number;
  goats: number;
}

export interface JamaahMetrics {
  total: number;
  newly_registered: number;
}

export interface LogisticsMetrics {
  upcoming_agendas: number;
  active_loans: number;
  broken_assets: number;
}

export interface SystemMetrics {
  uptime: string;
  db_status: string;
  active_sessions: number;
}

// --- Analytics (Chart Data) ---

export interface MonthlyFinancialPoint {
  month: string;
  income: number;
  expense: number;
}

export interface AssetStatusBreakdown {
  good: number;
  loaned: number;
  needs_repair: number;
  broken: number;
}

export interface MonthlyJamaahPoint {
  month: string;
  new_members: number;
}

export interface ZakatBreakdownMetrics {
  total_donated: number;
  total_distributed: number;
}

export interface QurbanTypeBreakdown {
  sapi: number;
  kambing: number;
  domba: number;
}

export interface AnalyticsMetrics {
  monthly_finance: MonthlyFinancialPoint[];
  asset_breakdown: AssetStatusBreakdown;
  qurban_breakdown: QurbanTypeBreakdown;
  jamaah_growth: MonthlyJamaahPoint[];
  zakat_breakdown: ZakatBreakdownMetrics;
}

export interface DashboardMetricsResponse {
  finance: FinanceMetrics;
  zakat: ZakatMetrics;
  qurban: QurbanMetrics;
  jamaah: JamaahMetrics;
  logistics: LogisticsMetrics;
  system: SystemMetrics;
  analytics: AnalyticsMetrics;
}
