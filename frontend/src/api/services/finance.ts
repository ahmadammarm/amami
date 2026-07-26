import apiClient from '../client';
import type { Fund, Transaction, CreateFundPayload, CreateTransactionPayload, PaginatedTransactionResponse } from '../../types/finance';

export const financeService = {
  async getFunds(): Promise<Fund[]> {
    const response = await apiClient.get<Fund[]>('/finance/funds');
    return response.data;
  },

  async createFund(payload: CreateFundPayload): Promise<Fund> {
    const response = await apiClient.post<Fund>('/finance/funds', payload);
    return response.data;
  },

  async getTransactions(page: number = 1, limit: number = 10): Promise<PaginatedTransactionResponse> {
    const response = await apiClient.get<PaginatedTransactionResponse>(`/finance/transactions?page=${page}&limit=${limit}`);
    return response.data;
  },

  async createTransaction(payload: CreateTransactionPayload): Promise<Transaction> {
    const response = await apiClient.post<Transaction>('/finance/transactions', payload);
    return response.data;
  }
};
