import apiClient from '../client';
import type { 
  Agenda,
  CreateAgendaPayload
} from '../../types/agenda';
import type { ApiResponse } from '../../types/api';

export const agendaService = {
  // Agenda
  async getAgendas(): Promise<Agenda[]> {
    const response = await apiClient.get<ApiResponse<Agenda[]>>('/agenda');
    return response.data.data;
  },

  async createAgenda(payload: CreateAgendaPayload): Promise<Agenda> {
    const response = await apiClient.post<ApiResponse<Agenda>>('/agenda', payload);
    return response.data.data;
  },

  async updateAgenda(id: string, payload: CreateAgendaPayload): Promise<Agenda> {
    const response = await apiClient.put<ApiResponse<Agenda>>(`/agenda/${id}`, payload);
    return response.data.data;
  }
};
