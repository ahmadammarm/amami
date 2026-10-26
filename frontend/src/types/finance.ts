import type { ISODateString } from './api';

export type TransactionType = 'CREDIT' | 'DEBIT';

export interface Fund {
  id: number;
  name: string;
  code: string;
  current_balance: number;
}

export interface Transaction {
  id: string;
  fund_id: number;
  type: TransactionType;
  amount: number;
  category: string;
  description: string;
  reference_id: string | null;
  /** Set when this entry reverses another transaction. */
  reversal_of_id: string | null;
  /** Set when another entry has reversed this transaction. */
  reversed_by_id: string | null;
  created_by: string;
  created_at: ISODateString;
}
