import { z } from 'zod';

const amountField = z
  .number({ required_error: 'Jumlah wajib diisi', invalid_type_error: 'Jumlah wajib diisi' })
  .positive('Jumlah harus lebih dari 0');

/**
 * Cash (IDR) must be posted to a fund in whole rupiah; goods (KG/LITER) are
 * recorded without touching the ledger.
 */
function refineLedgerPosting(
  value: { fund_id: number | null; amount: number; unit: 'IDR' | 'KG' | 'LITER' },
  ctx: z.RefinementCtx,
): void {
  if (value.fund_id === null) {
    if (value.unit === 'IDR') {
      ctx.addIssue({ code: z.ZodIssueCode.custom, path: ['fund_id'], message: 'Pilih kas untuk zakat uang' });
    }
    return;
  }
  if (value.unit !== 'IDR') {
    ctx.addIssue({ code: z.ZodIssueCode.custom, path: ['unit'], message: 'Hanya satuan IDR yang dapat dicatat ke kas' });
  }
  if (!Number.isInteger(value.amount)) {
    ctx.addIssue({ code: z.ZodIssueCode.custom, path: ['amount'], message: 'Nominal ke kas harus bilangan bulat' });
  }
}

export const collectZakatSchema = z
  .object({
    muzakki_id: z.string().uuid('Pilih muzakki'),
    zakat_type: z.enum(['FITRAH_BERAS', 'MAAL', 'INFAQ']),
    fund_id: z.number().int().positive().nullable(),
    amount: amountField,
    unit: z.enum(['IDR', 'KG', 'LITER']),
    description: z.string().trim().max(1000),
  })
  .superRefine(refineLedgerPosting);

export type CollectZakatPayload = z.infer<typeof collectZakatSchema>;

export const distributeZakatSchema = z
  .object({
    mustahik_id: z.string().uuid('Pilih mustahik'),
    fund_id: z.number().int().positive().nullable(),
    amount: amountField,
    unit: z.enum(['IDR', 'KG', 'LITER']),
    description: z.string().trim().max(1000),
  })
  .superRefine(refineLedgerPosting);

export type DistributeZakatPayload = z.infer<typeof distributeZakatSchema>;
