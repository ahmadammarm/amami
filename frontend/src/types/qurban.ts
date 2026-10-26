import type { ISODateString } from './api';

export type PackageType = 'KAMBING' | 'DOMBA' | 'SAPI' | 'SAPI_KOLEKTIF';
export type AnimalType = 'KAMBING' | 'DOMBA' | 'SAPI';
export type AnimalStatus = 'AVAILABLE' | 'ASSIGNED' | 'SLAUGHTERED';
export type PaymentStatus = 'PENDING' | 'PAID';

export interface QurbanPackage {
  id: number;
  name: string;
  type: PackageType;
  price: number;
  year_hijri: number;
  stock_total: number;
  stock_remaining: number;
}

export interface QurbanBooking {
  id: string;
  shohibul_id: string;
  shohibul_name: string;
  package_id: number;
  package_name: string;
  booking_date: ISODateString;
  payment_status: PaymentStatus;
  total_amount: number;
}

export interface QurbanAnimal {
  id: string;
  tag_number: string;
  type: AnimalType;
  weight: number;
  status: AnimalStatus;
  vendor_info: string;
}
