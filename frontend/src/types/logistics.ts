export interface Asset {
  id: string;
  name: string;
  sku: string;
  purchase_date: string | null;
  purchase_price: number;
  current_status: 'GOOD' | 'REPAIR' | 'BROKEN';
  location: string;
}

export interface CreateAssetPayload {
  name: string;
  sku: string;
  purchase_date: string | null;
  purchase_price: number;
  current_status: 'GOOD' | 'REPAIR' | 'BROKEN';
  location: string;
}

export interface AssetLoan {
  id: string;
  asset_id: string;
  jamaah_id: string;
  loan_date: string;
  due_date: string | null;
  return_date: string | null;
  condition_notes: string;
}

export interface CreateAssetLoanPayload {
  jamaah_id: string;
  loan_date: string;
  due_date: string | null;
  condition_notes: string;
}

export interface ReturnAssetLoanPayload {
  return_date: string;
  condition_notes: string;
}

export interface Agenda {
  id: string;
  title: string;
  description: string;
  start_time: string;
  end_time: string;
  location: string;
  status: 'SCHEDULED' | 'ONGOING' | 'COMPLETED' | 'CANCELLED';
  created_by_id: string | null;
}

export interface CreateAgendaPayload {
  title: string;
  description: string;
  start_time: string;
  end_time: string;
  location: string;
  status: 'SCHEDULED' | 'ONGOING' | 'COMPLETED' | 'CANCELLED';
}
