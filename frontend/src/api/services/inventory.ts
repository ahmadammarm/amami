import { request } from '../client';
import type { Page } from '../../types/api';
import type { Asset, AssetLoan } from '../../types/inventory';
import type { AssetPayload } from '../../schemas/inventory';

export interface CreateAssetLoanPayload {
  jamaah_id: string;
  loan_date: string;
  due_date: string | null;
  condition_notes: string;
}

export interface ReturnAssetLoanPayload {
  return_date: string;
  condition_notes: string;
}

export const inventoryService = {
  getAssets(page = 1, limit = 100): Promise<Page<Asset>> {
    return request<Page<Asset>>({ method: 'GET', url: '/inventory/assets', params: { page, limit } });
  },

  createAsset(payload: AssetPayload): Promise<Asset> {
    return request<Asset>({ method: 'POST', url: '/inventory/assets', data: payload });
  },

  updateAsset(id: string, payload: AssetPayload): Promise<Asset> {
    return request<Asset>({ method: 'PUT', url: `/inventory/assets/${id}`, data: payload });
  },

  deleteAsset(id: string): Promise<void> {
    return request<void>({ method: 'DELETE', url: `/inventory/assets/${id}` });
  },

  createLoan(assetId: string, payload: CreateAssetLoanPayload): Promise<AssetLoan> {
    return request<AssetLoan>({ method: 'POST', url: `/inventory/assets/${assetId}/loans`, data: payload });
  },

  returnLoan(loanId: string, payload: ReturnAssetLoanPayload): Promise<AssetLoan> {
    return request<AssetLoan>({ method: 'PUT', url: `/inventory/assets/loans/${loanId}/return`, data: payload });
  },
};
