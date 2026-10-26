import { z } from 'zod';

const digitsOnly = (label: string) =>
  z.string().trim().regex(/^\d*$/, `${label} hanya boleh berisi angka`);

export const mosqueProfileSchema = z.object({
  name: z.string().trim().min(1, 'Nama masjid wajib diisi'),
  address: z.string().trim().min(1, 'Alamat wajib diisi'),
  phone: z.string().trim().min(1, 'Nomor telepon wajib diisi'),
  legal_yayasan_id: z.string().trim().min(1, 'Nomor legal yayasan wajib diisi'),
  bank_name: z.string().trim(),
  bank_account_name: z.string().trim(),
  bank_account_number: digitsOnly('Nomor rekening'),
  zakat_fitrah_amount: digitsOnly('Besaran zakat fitrah'),
  active_hijri_year: digitsOnly('Tahun hijriah'),
});

export type MosqueProfileFormValues = z.infer<typeof mosqueProfileSchema>;
