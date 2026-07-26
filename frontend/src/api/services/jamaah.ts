import apiClient from '../client';
import type { Jamaah, CreateJamaahPayload, UpdateJamaahPayload, PaginatedJamaahResponse } from '../../types/jamaah';

export const jamaahService = {
  async getJamaahs(page: number = 1, limit: number = 10, query: string = ''): Promise<PaginatedJamaahResponse> {
    const response = await apiClient.get<PaginatedJamaahResponse>(`/jamaah?page=${page}&limit=${limit}&query=${encodeURIComponent(query)}`);
    return response.data;
  },

  async getJamaahByID(id: string): Promise<Jamaah> {
    const response = await apiClient.get<Jamaah>(`/jamaah/${id}`);
    return response.data;
  },

  async createJamaah(payload: CreateJamaahPayload): Promise<Jamaah> {
    const response = await apiClient.post<Jamaah>('/jamaah', payload);
    return response.data;
  },

  async updateJamaah(id: string, payload: UpdateJamaahPayload): Promise<Jamaah> {
    const response = await apiClient.put<Jamaah>(`/jamaah/${id}`, payload);
    return response.data;
  }
};
