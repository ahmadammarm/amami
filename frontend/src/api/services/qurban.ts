import { request } from '../client';
import type { Page } from '../../types/api';
import type { QurbanAnimal, QurbanBooking, QurbanPackage } from '../../types/qurban';
import type { AnimalPayload, BookingPayload, PackagePayload } from '../../schemas/qurban';

export const qurbanService = {
  getPackages(page = 1, limit = 100): Promise<Page<QurbanPackage>> {
    return request<Page<QurbanPackage>>({ method: 'GET', url: '/qurban/packages', params: { page, limit } });
  },

  createPackage(payload: PackagePayload): Promise<QurbanPackage> {
    return request<QurbanPackage>({ method: 'POST', url: '/qurban/packages', data: payload });
  },

  updatePackage(id: number, payload: PackagePayload): Promise<QurbanPackage> {
    return request<QurbanPackage>({ method: 'PUT', url: `/qurban/packages/${id}`, data: payload });
  },

  deletePackage(id: number): Promise<void> {
    return request<void>({ method: 'DELETE', url: `/qurban/packages/${id}` });
  },

  getBookings(): Promise<QurbanBooking[]> {
    return request<QurbanBooking[]>({ method: 'GET', url: '/qurban/bookings' });
  },

  createBooking(payload: BookingPayload): Promise<QurbanBooking> {
    return request<QurbanBooking>({ method: 'POST', url: '/qurban/bookings', data: payload });
  },

  getAnimals(): Promise<QurbanAnimal[]> {
    return request<QurbanAnimal[]>({ method: 'GET', url: '/qurban/animals' });
  },

  createAnimal(payload: AnimalPayload): Promise<QurbanAnimal> {
    return request<QurbanAnimal>({ method: 'POST', url: '/qurban/animals', data: payload });
  },

  updateAnimal(id: string, payload: AnimalPayload): Promise<QurbanAnimal> {
    return request<QurbanAnimal>({ method: 'PUT', url: `/qurban/animals/${id}`, data: payload });
  },

  deleteAnimal(id: string): Promise<void> {
    return request<void>({ method: 'DELETE', url: `/qurban/animals/${id}` });
  },
};
