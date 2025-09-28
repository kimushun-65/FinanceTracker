import { apiClient } from '@/shared/api/client';
import type {
  Transaction,
  CreateTransactionPayload,
  UpdateTransactionPayload,
  TransactionListParams,
  TransactionListResponse,
  MonthlySummaryParams,
  MonthlySummary,
} from '../model';
import { endpoints } from './transaction.endpoints';
import { createTransactionFilters } from './transaction.filters';

export const transactionApi = {
  list: async (
    params: TransactionListParams = {},
  ): Promise<TransactionListResponse> => {
    const queryString = createTransactionFilters(params);
    const url = queryString
      ? `${endpoints.list}?${queryString}`
      : endpoints.list;
    const response = await apiClient.get<TransactionListResponse>(url);
    return response.data;
  },

  get: async (id: string): Promise<Transaction> => {
    const response = await apiClient.get<Transaction>(endpoints.get(id));
    return response.data;
  },

  create: async (payload: CreateTransactionPayload): Promise<Transaction> => {
    const response = await apiClient.post<Transaction>(
      endpoints.create,
      payload,
    );
    return response.data;
  },

  update: async (
    id: string,
    payload: UpdateTransactionPayload,
  ): Promise<Transaction> => {
    const response = await apiClient.put<Transaction>(
      endpoints.update(id),
      payload,
    );
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(endpoints.delete(id));
  },

  getMonthlySummary: async (
    params: MonthlySummaryParams = {},
  ): Promise<MonthlySummary> => {
    const queryParams = new URLSearchParams();
    if (params.year) queryParams.append('year', params.year.toString());
    if (params.month) queryParams.append('month', params.month.toString());
    
    const url = queryParams.toString()
      ? `${endpoints.monthlySummary}?${queryParams.toString()}`
      : endpoints.monthlySummary;
    
    const response = await apiClient.get<MonthlySummary>(url);
    return response.data;
  },
};
