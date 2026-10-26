import type { ISODateString } from './api';

export type AssetStatus = 'GOOD' | 'LOANED' | 'NEEDS_REPAIR' | 'REPAIR' | 'BROKEN';

export interface Asset {
  id: string;
  name: string;
  sku: string;
  purchase_date: ISODateString | null;
  purchase_price: number;
  current_status: AssetStatus;
  location: string;
}

export interface AssetLoan {
  id: string;
  asset_id: string;
  jamaah_id: string;
  loan_date: ISODateString;
  due_date: ISODateString | null;
  return_date: ISODateString | null;
  condition_notes: string;
}
