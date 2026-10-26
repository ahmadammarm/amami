import { request } from '../client';
import type { Page } from '../../types/api';
import type { ZakatDistribution, ZakatDonation } from '../../types/zakat';
import type { CollectZakatPayload, DistributeZakatPayload } from '../../schemas/zakat';

export const zakatService = {
  getDonations(page = 1, limit = 10): Promise<Page<ZakatDonation>> {
    return request<Page<ZakatDonation>>({ method: 'GET', url: '/zakat/donations', params: { page, limit } });
  },

  getDistributions(page = 1, limit = 10): Promise<Page<ZakatDistribution>> {
    return request<Page<ZakatDistribution>>({ method: 'GET', url: '/zakat/distributions', params: { page, limit } });
  },

  collectZakat(payload: CollectZakatPayload): Promise<ZakatDonation> {
    return request<ZakatDonation>({ method: 'POST', url: '/zakat/donations', data: payload });
  },

  distributeZakat(payload: DistributeZakatPayload): Promise<ZakatDistribution> {
    return request<ZakatDistribution>({ method: 'POST', url: '/zakat/distributions', data: payload });
  },
};
