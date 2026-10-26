import { z } from 'zod';

export const transactionSchema = z.object({
  fund_id: z.number({ required_error: 'Pilih kas' }).int().positive('Pilih kas'),
  type: z.enum(['CREDIT', 'DEBIT']),
  amount: z
    .number({ required_error: 'Nominal wajib diisi', invalid_type_error: 'Nominal wajib diisi' })
    .int('Nominal harus bilangan bulat')
    .positive('Nominal harus lebih dari 0'),
  category: z.string().trim().min(1, 'Kategori wajib diisi').max(30, 'Kategori maksimal 30 karakter'),
  description: z.string().trim().max(255, 'Keterangan maksimal 255 karakter'),
});

export type CreateTransactionPayload = z.infer<typeof transactionSchema>;

export const reversalSchema = z.object({
  reason: z.string().trim().min(3, 'Alasan minimal 3 karakter').max(255, 'Alasan maksimal 255 karakter'),
});

export type ReverseTransactionPayload = z.infer<typeof reversalSchema>;

export const fundSchema = z.object({
  name: z.string().trim().min(1, 'Nama kas wajib diisi').max(100),
  code: z.string().trim().min(1, 'Kode kas wajib diisi').max(20),
});

export type CreateFundPayload = z.infer<typeof fundSchema>;
