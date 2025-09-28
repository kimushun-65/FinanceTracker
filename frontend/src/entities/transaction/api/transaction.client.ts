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
import { transformApiTransactionResponse } from '../lib';

export const transactionApi = {
  list: async (
    params: TransactionListParams = {},
  ): Promise<TransactionListResponse> => {
    const queryString = createTransactionFilters(params);
    const url = queryString ? `${endpoints.list}?${queryString}` : endpoints.list;
    const res = await apiClient.get<any>(url);
    const api = res.data;

    const perPage = Number(api?.per_page ?? params.limit ?? 20);
    const page = Number(api?.page ?? (params.offset ? Math.floor(params.offset / perPage) + 1 : 1));
    const totalCount = Number(api?.total_count ?? 0);

    const transactions: Transaction[] = Array.isArray(api?.transactions)
      ? api.transactions.map((t: any) => transformApiTransactionResponse(t))
      : [];

    const hasMore = page * perPage < totalCount;

    const result: TransactionListResponse = {
      transactions,
      total: totalCount,
      hasMore,
    };
    return result;
  },

  get: async (id: string): Promise<Transaction> => {
    const response = await apiClient.get<any>(endpoints.get(id));
    return transformApiTransactionResponse(response.data);
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
    
    const response = await apiClient.get<any>(url);
    const api = response.data;

    const toNumber = (v: any): number => {
      if (typeof v === 'number') return Math.round(v);
      if (v == null) return 0;
      const n = Number(String(v));
      return Number.isFinite(n) ? Math.round(n) : 0;
    };

    const transactionCount = Array.isArray(api?.daily_data)
      ? api.daily_data.reduce((sum: number, d: any) => sum + (d?.count ?? 0), 0)
      : 0;

    const summary: MonthlySummary = {
      totalIncome: { amount: toNumber(api?.total_income), currency: 'JPY' },
      totalExpense: { amount: toNumber(api?.total_expense), currency: 'JPY' },
      netIncome: { amount: toNumber(api?.net_amount), currency: 'JPY' },
      transactionCount,
      period: {
        year: api?.year ?? (params.year ?? new Date().getFullYear()),
        month: api?.month ?? (params.month ?? new Date().getMonth() + 1),
      },
    };

    return summary;
  },
};
