import type { ISODateString } from './api';

export type AgendaStatus = 'SCHEDULED' | 'ONGOING' | 'COMPLETED' | 'CANCELLED';

export interface Agenda {
  id: string;
  title: string;
  description: string;
  day: string;
  time: string;
  location: string;
  status: AgendaStatus;
  created_by_id: string | null;
}

export interface AgendaDocumentation {
  id: string;
  agenda_id: string;
  /** Server-relative URL; pass through assetUrl() before rendering. */
  file_url: string;
  file_name: string;
  file_type: 'image/jpeg' | 'image/png' | 'image/webp' | 'application/pdf';
  file_size: number;
  uploaded_by: string | null;
  uploaded_at: ISODateString;
}
