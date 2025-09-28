import { apiClient } from '@/shared/api/client';
import type {
  Budget,
  BudgetWithUsage,
  CreateBudgetPayload,
  UpdateBudgetPayload,
  BudgetListResponse,
  BudgetSummary,
  PeriodType,
} from '../model';
import { endpoints } from './budget.endpoints';

export const budgetApi = {
  list: async (
    period?: PeriodType,
    categoryId?: string,
  ): Promise<BudgetListResponse> => {
    const params = new URLSearchParams();
    if (period) params.append('period', period);
    if (categoryId) params.append('categoryId', categoryId);

    const url = params.toString()
      ? `${endpoints.list}?${params.toString()}`
      : endpoints.list;
    const response = await apiClient.get<BudgetListResponse>(url);
    return response.data;
  },

  get: async (id: string): Promise<BudgetWithUsage> => {
    const response = await apiClient.get<BudgetWithUsage>(endpoints.get(id));
    return response.data;
  },

  create: async (payload: CreateBudgetPayload): Promise<Budget> => {
    const response = await apiClient.post<Budget>(endpoints.create, payload);
    return response.data;
  },

  update: async (id: string, payload: UpdateBudgetPayload): Promise<Budget> => {
    const response = await apiClient.put<Budget>(endpoints.update(id), payload);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(endpoints.delete(id));
  },

  getSummary: async (): Promise<BudgetSummary> => {
    const response = await apiClient.get<BudgetSummary>(endpoints.summary);
    return response.data;
  },
};
