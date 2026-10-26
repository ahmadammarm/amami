import { z } from 'zod';

export const packageSchema = z.object({
  name: z.string().trim().min(1, 'Nama paket wajib diisi').max(100),
  type: z.enum(['KAMBING', 'DOMBA', 'SAPI', 'SAPI_KOLEKTIF']),
  price: z.number({ invalid_type_error: 'Harga wajib diisi' }).int().positive('Harga harus lebih dari 0'),
  year_hijri: z.number({ invalid_type_error: 'Tahun wajib diisi' }).int().min(1400).max(1600),
  stock_total: z.number({ invalid_type_error: 'Stok wajib diisi' }).int().positive('Stok minimal 1'),
});

export type PackagePayload = z.infer<typeof packageSchema>;

export const animalSchema = z.object({
  tag_number: z.string().trim().min(1, 'Nomor tag wajib diisi').max(50),
  type: z.enum(['KAMBING', 'DOMBA', 'SAPI']),
  weight: z.number({ invalid_type_error: 'Berat wajib diisi' }).positive('Berat harus lebih dari 0').lt(10000),
  status: z.enum(['AVAILABLE', 'ASSIGNED', 'SLAUGHTERED']),
  vendor_info: z.string().trim().max(1000),
});

export type AnimalPayload = z.infer<typeof animalSchema>;

export const bookingSchema = z.object({
  shohibul_id: z.string().uuid('Pilih shohibul'),
  package_id: z.number({ required_error: 'Pilih paket' }).int().positive('Pilih paket'),
  payment_status: z.enum(['PENDING', 'PAID']),
});

export type BookingPayload = z.infer<typeof bookingSchema>;
