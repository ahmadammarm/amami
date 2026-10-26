import type { ISODateString } from './api';

export type ZakatType = 'FITRAH_BERAS' | 'MAAL' | 'INFAQ';
export type ZakatUnit = 'IDR' | 'KG' | 'LITER';

export interface ZakatDonation {
  id: string;
  transaction_id: string | null;
  muzakki_id: string;
  muzakki_name: string;
  zakat_type: string;
  amount: number;
  unit: ZakatUnit;
  description: string;
  created_at: ISODateString;
}

export interface ZakatDistribution {
  id: string;
  transaction_id: string | null;
  mustahik_id: string;
  mustahik_name: string;
  amount: number;
  unit: ZakatUnit;
  description: string;
  distributed_at: ISODateString;
}
