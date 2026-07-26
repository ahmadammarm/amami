import apiClient from '../client';
import type { ApiResponse } from '../../types/api';
import type { 
  QurbanPackage, CreatePackagePayload, 
  QurbanBooking, CreateBookingPayload,
  QurbanAnimal, CreateAnimalPayload 
} from '../../types/qurban';

export const qurbanService = {
  // Packages
  async getPackages(page = 1, limit = 100): Promise<{ items: QurbanPackage[], total: number }> {
    const response = await apiClient.get(`/qurban/packages?page=${page}&limit=${limit}`);
    return response.data.data;
  },
  
  async createPackage(payload: CreatePackagePayload): Promise<QurbanPackage> {
    const response = await apiClient.post<ApiResponse<QurbanPackage>>('/qurban/packages', payload);
    return response.data.data;
  },

  async updatePackage(id: number, payload: CreatePackagePayload): Promise<QurbanPackage> {
    const response = await apiClient.put<ApiResponse<QurbanPackage>>(`/qurban/packages/${id}`, payload);
    return response.data.data;
  },

  async deletePackage(id: number): Promise<void> {
    await apiClient.delete(`/qurban/packages/${id}`);
  },

  // Bookings
  async getBookings(): Promise<QurbanBooking[]> {
    const response = await apiClient.get<ApiResponse<QurbanBooking[]>>('/qurban/bookings');
    return response.data.data;
  },

  async createBooking(payload: CreateBookingPayload): Promise<QurbanBooking> {
    const response = await apiClient.post<ApiResponse<QurbanBooking>>('/qurban/bookings', payload);
    return response.data.data;
  },

  // Animals
  async getAnimals(): Promise<QurbanAnimal[]> {
    const response = await apiClient.get<ApiResponse<QurbanAnimal[]>>('/qurban/animals');
    return response.data.data;
  },

  async createAnimal(payload: CreateAnimalPayload): Promise<QurbanAnimal> {
    const response = await apiClient.post<ApiResponse<QurbanAnimal>>('/qurban/animals', payload);
    return response.data.data;
  },

  async updateAnimal(id: string, payload: CreateAnimalPayload): Promise<QurbanAnimal> {
    const response = await apiClient.put<ApiResponse<QurbanAnimal>>(`/qurban/animals/${id}`, payload);
    return response.data.data;
  },

  async deleteAnimal(id: string): Promise<void> {
    await apiClient.delete(`/qurban/animals/${id}`);
  }
};
