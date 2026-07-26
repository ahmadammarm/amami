export interface CollectZakatPayload {
  muzakki_id: string;
  zakat_type: string;
  fund_id?: number;
  amount: number;
  unit: string;
  description: string;
}

export interface DistributeZakatPayload {
  mustahik_id: string;
  fund_id?: number;
  amount: number;
  unit: string;
  description: string;
}

export interface ZakatDonation {
  id: string;
  transaction_id?: string;
  muzakki_id: string;
  muzakki_name: string;
  zakat_type: string;
  amount: number;
  unit: string;
  description: string;
  created_at: string;
}

export interface ZakatDistribution {
  id: string;
  transaction_id?: string;
  mustahik_id: string;
  mustahik_name: string;
  amount: number;
  unit: string;
  description: string;
  distributed_at: string;
}

export interface PaginatedDonationResponse {
  data: ZakatDonation[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface PaginatedDistributionResponse {
  data: ZakatDistribution[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}
