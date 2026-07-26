import apiClient from '../client';
import type { 
  Asset, 
  CreateAssetPayload, 
  CreateAssetLoanPayload, 
  ReturnAssetLoanPayload
} from '../../types/inventory';
import type { ApiResponse } from '../../types/api';

export const inventoryService = {
  // Assets
  async getAssets(page = 1, limit = 100): Promise<{ items: Asset[], total: number }> {
    const response = await apiClient.get(`/inventory/assets?page=${page}&limit=${limit}`);
    return response.data.data;
  },

  async createAsset(payload: CreateAssetPayload): Promise<Asset> {
    const response = await apiClient.post<ApiResponse<Asset>>('/inventory/assets', payload);
    return response.data.data;
  },

  async updateAsset(id: string, payload: CreateAssetPayload): Promise<Asset> {
    const response = await apiClient.put<ApiResponse<Asset>>(`/inventory/assets/${id}`, payload);
    return response.data.data;
  },

  async deleteAsset(id: string): Promise<void> {
    await apiClient.delete(`/inventory/assets/${id}`);
  },

  // Loans
  async createLoan(assetId: string, payload: CreateAssetLoanPayload): Promise<void> {
    await apiClient.post(`/inventory/assets/${assetId}/loans`, payload);
  },

  async returnLoan(loanId: string, payload: ReturnAssetLoanPayload): Promise<void> {
    await apiClient.put(`/inventory/assets/loans/${loanId}/return`, payload);
  }
};
