import { z } from 'zod';

export const agendaSchema = z.object({
  title: z.string().trim().min(1, 'Judul wajib diisi').max(200, 'Judul maksimal 200 karakter'),
  description: z.string().trim().max(5000, 'Deskripsi terlalu panjang'),
  day: z.string().trim().min(1, 'Hari wajib diisi').max(20),
  time: z.string().trim().min(1, 'Waktu wajib diisi').max(20),
  location: z.string().trim().max(100, 'Lokasi maksimal 100 karakter'),
  status: z.enum(['SCHEDULED', 'ONGOING', 'COMPLETED', 'CANCELLED']),
});

export type AgendaPayload = z.infer<typeof agendaSchema>;

/** Mirrors the backend upload rules (MAX_UPLOAD_MB default and allowed types). */
export const DOCUMENTATION_MAX_BYTES = 10 * 1024 * 1024;
export const DOCUMENTATION_MAX_FILES = 10;
export const DOCUMENTATION_ACCEPT = 'image/jpeg,image/png,image/webp,application/pdf';
