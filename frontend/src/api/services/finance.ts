import { request } from '../client';
import type { Page } from '../../types/api';
import type { Fund, Transaction } from '../../types/finance';
import type { CreateFundPayload, CreateTransactionPayload, ReverseTransactionPayload } from '../../schemas/finance';

export const financeService = {
  getFunds(): Promise<Fund[]> {
    return request<Fund[]>({ method: 'GET', url: '/finance/funds' });
  },

  createFund(payload: CreateFundPayload): Promise<Fund> {
    return request<Fund>({ method: 'POST', url: '/finance/funds', data: payload });
  },

  getTransactions(page = 1, limit = 10): Promise<Page<Transaction>> {
    return request<Page<Transaction>>({ method: 'GET', url: '/finance/transactions', params: { page, limit } });
  },

  createTransaction(payload: CreateTransactionPayload): Promise<Transaction> {
    return request<Transaction>({ method: 'POST', url: '/finance/transactions', data: payload });
  },

  reverseTransaction(id: string, payload: ReverseTransactionPayload): Promise<Transaction> {
    return request<Transaction>({ method: 'POST', url: `/finance/transactions/${id}/reverse`, data: payload });
  },
};
