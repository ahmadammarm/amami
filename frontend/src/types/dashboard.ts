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

export interface DashboardMetricsResponse {
  finance: FinanceMetrics;
  zakat: ZakatMetrics;
  qurban: QurbanMetrics;
  jamaah: JamaahMetrics;
  logistics: LogisticsMetrics;
  system: SystemMetrics;
}
