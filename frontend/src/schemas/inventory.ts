import { z } from 'zod';

export const assetSchema = z.object({
  name: z.string().trim().min(1, 'Nama aset wajib diisi').max(255),
  sku: z.string().trim().min(1, 'SKU wajib diisi').max(50, 'SKU maksimal 50 karakter'),
  /** yyyy-mm-dd from a date input, or null when unknown. */
  purchase_date: z.string().nullable(),
  purchase_price: z
    .number({ invalid_type_error: 'Harga wajib diisi' })
    .int('Harga harus bilangan bulat')
    .min(0, 'Harga tidak boleh negatif'),
  current_status: z.enum(['GOOD', 'LOANED', 'NEEDS_REPAIR', 'REPAIR', 'BROKEN']),
  location: z.string().trim().max(100, 'Lokasi maksimal 100 karakter'),
});

export type AssetFormValues = z.infer<typeof assetSchema>;

/** API payload: dates are sent as full ISO timestamps. */
export type AssetPayload = Omit<AssetFormValues, 'purchase_date'> & { purchase_date: string | null };

export function toAssetPayload(values: AssetFormValues): AssetPayload {
  return {
    ...values,
    purchase_date: values.purchase_date ? new Date(values.purchase_date).toISOString() : null,
  };
}
