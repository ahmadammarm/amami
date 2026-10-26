import { request } from '../client';
import type { Agenda, AgendaDocumentation } from '../../types/agenda';
import type { AgendaPayload } from '../../schemas/agenda';

export const agendaService = {
  getAgendas(): Promise<Agenda[]> {
    return request<Agenda[]>({ method: 'GET', url: '/agenda' });
  },

  createAgenda(payload: AgendaPayload): Promise<Agenda> {
    return request<Agenda>({ method: 'POST', url: '/agenda', data: payload });
  },

  updateAgenda(id: string, payload: AgendaPayload): Promise<Agenda> {
    return request<Agenda>({ method: 'PUT', url: `/agenda/${id}`, data: payload });
  },

  deleteAgenda(id: string): Promise<void> {
    return request<void>({ method: 'DELETE', url: `/agenda/${id}` });
  },

  getDocumentation(agendaId: string): Promise<AgendaDocumentation[]> {
    return request<AgendaDocumentation[]>({ method: 'GET', url: `/agenda/${agendaId}/documentation` });
  },

  uploadDocumentation(agendaId: string, files: File[]): Promise<AgendaDocumentation[]> {
    const form = new FormData();
    files.forEach((file) => form.append('files', file));
    return request<AgendaDocumentation[]>({
      method: 'POST',
      url: `/agenda/${agendaId}/documentation`,
      data: form,
      // Overrides the client's JSON default; axios fills in the multipart boundary.
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },

  deleteDocumentation(agendaId: string, docId: string): Promise<void> {
    return request<void>({ method: 'DELETE', url: `/agenda/${agendaId}/documentation/${docId}` });
  },
};
