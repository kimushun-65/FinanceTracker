import { apiClient } from '@/shared/api/client';
import type {
  Transaction,
  CreateTransactionPayload,
  UpdateTransactionPayload,
  TransactionListParams,
  TransactionListResponse,
  MonthlySummaryParams,
  MonthlySummary,
  CategorySummaryParams,
  CategorySummary,
} from '../model';
import { endpoints } from './transaction.endpoints';
import { createTransactionFilters } from './transaction.filters';
import { transformApiTransactionResponse } from '../lib';

export const transactionApi = {
  list: async (
    params: TransactionListParams = {},
  ): Promise<TransactionListResponse> => {
    const queryString = createTransactionFilters(params);
    const url = queryString
      ? `${endpoints.list}?${queryString}`
      : endpoints.list;
    const res = await apiClient.get<any>(url);
    const api = res.data;

    const perPage = Number(api?.per_page ?? params.limit ?? 20);
    const page = Number(
      api?.page ??
        (params.offset ? Math.floor(params.offset / perPage) + 1 : 1),
    );
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
    // 日付をISO文字列に変換（バックエンドが期待するRFC3339形式）
    const formatDate = (dateString?: string): string | undefined => {
      if (!dateString) return undefined;
      const date = new Date(dateString);
      return date.toISOString();
    };

    // バックエンドが期待する形式に変換
    const backendPayload = {
      account_id: payload.accountId,
      category_id: payload.categoryId,
      amount: payload.amount.amount, // Moneyオブジェクトから数値を取り出す
      transaction_type: payload.type,
      description: payload.description,
      transaction_date: formatDate(payload.date),
    };

    const response = await apiClient.post<Transaction>(
      endpoints.create,
      backendPayload,
    );
    return response.data;
  },

  update: async (
    id: string,
    payload: UpdateTransactionPayload,
  ): Promise<Transaction> => {
    // 日付をISO文字列に変換（バックエンドが期待するRFC3339形式）
    const formatDate = (dateString?: string): string | undefined => {
      if (!dateString) return undefined;
      const date = new Date(dateString);
      return date.toISOString();
    };

    // バックエンドが期待する形式に変換
    const backendPayload: any = {};

    if (payload.accountId !== undefined) {
      backendPayload.account_id = payload.accountId;
    }
    if (payload.categoryId !== undefined) {
      backendPayload.category_id = payload.categoryId;
    }
    if (payload.amount !== undefined) {
      backendPayload.amount = payload.amount.amount; // Moneyオブジェクトから数値を取り出す
    }
    if (payload.description !== undefined) {
      backendPayload.description = payload.description;
    }
    if (payload.date !== undefined) {
      backendPayload.transaction_date = formatDate(payload.date);
    }

    const response = await apiClient.put<Transaction>(
      endpoints.update(id),
      backendPayload,
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
        year: api?.year ?? params.year ?? new Date().getFullYear(),
        month: api?.month ?? params.month ?? new Date().getMonth() + 1,
      },
    };

    return summary;
  },

  getCategorySummary: async (
    params: CategorySummaryParams = {},
  ): Promise<CategorySummary[]> => {
    const queryParams = new URLSearchParams();
    if (params.type) queryParams.append('type', params.type);
    if (params.from) queryParams.append('from', params.from);
    if (params.to) queryParams.append('to', params.to);

    const url = queryParams.toString()
      ? `${endpoints.categorySummary}?${queryParams.toString()}`
      : endpoints.categorySummary;

    const response = await apiClient.get<any>(url);
    const api = response.data;

    const toNumber = (v: any): number => {
      if (typeof v === 'number') return Math.round(v);
      if (v == null) return 0;
      const n = Number(String(v));
      return Number.isFinite(n) ? Math.round(n) : 0;
    };

    const summaries: CategorySummary[] = Array.isArray(api?.categories)
      ? api.categories.map((cat: any) => ({
          categoryId: cat?.category_id ?? '',
          categoryName: cat?.category_name ?? '',
          totalAmount: { amount: toNumber(cat?.total_amount), currency: 'JPY' },
          transactionCount: cat?.transaction_count ?? 0,
          percentage: cat?.percentage ?? 0,
          type: (cat?.type ?? 'expense') as 'income' | 'expense',
        }))
      : [];

    return summaries;
  },
};
