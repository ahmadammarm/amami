export interface Fund {
  id: number;
  name: string;
  code: string;
  current_balance: number;
}

export interface Transaction {
  id: string;
  fund_id: number;
  type: 'CREDIT' | 'DEBIT';
  amount: number;
  category: string;
  reference_id?: string;
  created_by: string;
  created_at: string;
}

export interface CreateFundPayload {
  name: string;
  code: string;
}

export interface CreateTransactionPayload {
  fund_id: number;
  type: 'CREDIT' | 'DEBIT';
  amount: number;
  category: string;
  description?: string;
  reference_id?: string;
}

export interface PaginatedTransactionResponse {
  data: Transaction[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}
