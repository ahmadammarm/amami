import apiClient from '../client';
import type { 
  Asset, 
  CreateAssetPayload, 
  CreateAssetLoanPayload, 
  ReturnAssetLoanPayload,
  Agenda,
  CreateAgendaPayload
} from '../../types/logistics';
import type { ApiResponse } from '../../types/api';

export const logisticsService = {
  // Assets
  async getAssets(page = 1, limit = 100): Promise<{ items: Asset[], total: number }> {
    const response = await apiClient.get(`/logistics/assets?page=${page}&limit=${limit}`);
    return response.data.data;
  },

  async createAsset(payload: CreateAssetPayload): Promise<Asset> {
    const response = await apiClient.post<ApiResponse<Asset>>('/logistics/assets', payload);
    return response.data.data;
  },

  async updateAsset(id: string, payload: CreateAssetPayload): Promise<Asset> {
    const response = await apiClient.put<ApiResponse<Asset>>(`/logistics/assets/${id}`, payload);
    return response.data.data;
  },

  async deleteAsset(id: string): Promise<void> {
    await apiClient.delete(`/logistics/assets/${id}`);
  },

  // Loans
  async createLoan(assetId: string, payload: CreateAssetLoanPayload): Promise<void> {
    await apiClient.post(`/logistics/assets/${assetId}/loans`, payload);
  },

  async returnLoan(loanId: string, payload: ReturnAssetLoanPayload): Promise<void> {
    await apiClient.put(`/logistics/assets/loans/${loanId}/return`, payload);
  },

  // Agenda
  async getAgendas(): Promise<Agenda[]> {
    const response = await apiClient.get<ApiResponse<Agenda[]>>('/logistics/agenda');
    return response.data.data;
  },

  async createAgenda(payload: CreateAgendaPayload): Promise<Agenda> {
    const response = await apiClient.post<ApiResponse<Agenda>>('/logistics/agenda', payload);
    return response.data.data;
  },

  async updateAgenda(id: string, payload: CreateAgendaPayload): Promise<Agenda> {
    const response = await apiClient.put<ApiResponse<Agenda>>(`/logistics/agenda/${id}`, payload);
    return response.data.data;
  }
};
