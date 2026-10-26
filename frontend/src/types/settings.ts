export interface MosqueProfile {
  name: string;
  address: string;
  phone: string;
  logo: string;
  legal_yayasan_id: string;
  bank_name: string;
  bank_account_name: string;
  bank_account_number: string;
  zakat_fitrah_amount: string;
  active_hijri_year: string;
}

/** Only the name is required by the backend; every other field may be left out. */
export type MosqueProfilePayload = Pick<MosqueProfile, 'name'> & Partial<Omit<MosqueProfile, 'name'>>;

export interface SMTPConfig {
  host: string;
  port: string;
  username: string;
  from: string;
}

export interface SMTPConfigPayload extends SMTPConfig {
  password: string;
}

export interface TestSMTPPayload {
  host: string;
  port: string;
  username: string;
  password: string;
}

export interface SystemHealth {
  database: boolean;
  smtp: boolean;
}

export interface AuditLog {
  id: number;
  user: string;
  action: string;
  entity: string;
  timestamp: string;
}

export interface AuditLogList {
  logs: AuditLog[];
  total_count: number;
  page: number;
  limit: number;
}
