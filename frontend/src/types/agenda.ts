export interface Agenda {
  id: string;
  title: string;
  description: string;
  day: string;
  time: string;
  location: string;
  status: 'SCHEDULED' | 'ONGOING' | 'COMPLETED' | 'CANCELLED';
  created_by_id: string | null;
}

export interface CreateAgendaPayload {
  title: string;
  description: string;
  day: string;
  time: string;
  location: string;
  status: 'SCHEDULED' | 'ONGOING' | 'COMPLETED' | 'CANCELLED';
}
