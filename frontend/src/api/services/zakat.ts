import apiClient from '../client';
import type { 
  ZakatDonation, 
  ZakatDistribution,
  CollectZakatPayload, 
  DistributeZakatPayload, 
  PaginatedDonationResponse, 
  PaginatedDistributionResponse 
} from '../../types/zakat';

export const zakatService = {
  async getDonations(page: number = 1, limit: number = 10): Promise<PaginatedDonationResponse> {
    const response = await apiClient.get<PaginatedDonationResponse>(`/zakat/donations?page=${page}&limit=${limit}`);
    return response.data;
  },

  async getDistributions(page: number = 1, limit: number = 10): Promise<PaginatedDistributionResponse> {
    const response = await apiClient.get<PaginatedDistributionResponse>(`/zakat/distributions?page=${page}&limit=${limit}`);
    return response.data;
  },

  async collectZakat(payload: CollectZakatPayload): Promise<ZakatDonation> {
    const response = await apiClient.post<ZakatDonation>('/zakat/donations', payload);
    return response.data;
  },

  async distributeZakat(payload: DistributeZakatPayload): Promise<ZakatDistribution> {
    const response = await apiClient.post<ZakatDistribution>('/zakat/distributions', payload);
    return response.data;
  }
};
