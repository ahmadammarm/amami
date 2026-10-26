import { z } from 'zod';

export const jamaahSchema = z.object({
  full_name: z.string().trim().min(1, 'Nama lengkap wajib diisi').max(100, 'Nama maksimal 100 karakter'),
  phone: z.string().trim().max(20, 'Nomor telepon maksimal 20 karakter'),
  address: z.string().trim().max(1000, 'Alamat terlalu panjang'),
  is_mustahik: z.boolean(),
});

export type JamaahPayload = z.infer<typeof jamaahSchema>;

export const mustahikSurveySchema = z.object({
  total_income: z
    .number({ required_error: 'Penghasilan wajib diisi', invalid_type_error: 'Penghasilan wajib diisi' })
    .int('Penghasilan harus bilangan bulat')
    .min(0, 'Penghasilan tidak boleh negatif'),
  dependents_count: z
    .number({ required_error: 'Jumlah tanggungan wajib diisi', invalid_type_error: 'Jumlah tanggungan wajib diisi' })
    .int()
    .min(0, 'Tidak boleh negatif')
    .max(30, 'Maksimal 30 tanggungan'),
  house_status: z.enum(['OWNED', 'FAMILY', 'RENTED', 'NONE'], { required_error: 'Pilih status rumah' }),
  has_disability: z.boolean(),
  criteria_notes: z.string().trim().max(2000, 'Catatan terlalu panjang'),
});

export type MustahikSurveyPayload = z.infer<typeof mustahikSurveySchema>;
