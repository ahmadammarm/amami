export interface QurbanPackage {
  id: number;
  name: string;
  type: string;
  price: number;
  year_hijri: number;
  stock_total: number;
  stock_remaining: number;
}

export interface CreatePackagePayload {
  name: string;
  type: string;
  price: number;
  year_hijri: number;
  stock_total: number;
}

export interface QurbanBooking {
  id: string;
  shohibul_id: string;
  package_id: number;
  booking_date: string;
  payment_status: string;
  total_amount: number;
}

export interface CreateBookingPayload {
  shohibul_id: string;
  package_id: number;
  payment_status: string;
  total_amount: number;
}

export interface QurbanAnimal {
  id: string;
  tag_number: string;
  type: string;
  weight: number;
  status: string;
  vendor_info: string;
}

export interface CreateAnimalPayload {
  tag_number: string;
  type: string;
  weight: number;
  status: string;
  vendor_info: string;
}
