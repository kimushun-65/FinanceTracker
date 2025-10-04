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
import { calculateBudgetUsage } from '../lib';
import { endpoints } from './budget.endpoints';

// Backend response type
type BackendBudgetResponse = {
  id: string;
  user_id: string;
  category_id: string;
  amount: number;
  period: string;
  start_date: string;
  end_date?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

type BackendBudgetListResponse = {
  budgets: BackendBudgetResponse[];
  total_count: number;
};

type BackendBudgetSummaryResponse = {
  period: 'monthly' | 'yearly';
  start_date: string;
  end_date: string;
  total_budget: number;
  total_used: number;
  total_remaining: number;
};

// Transform backend response to frontend type
const transformBudgetResponse = (
  backendBudget: BackendBudgetResponse,
): BudgetWithUsage => {
  const budget: Budget = {
    id: backendBudget.id,
    userId: backendBudget.user_id,
    categoryId: backendBudget.category_id,
    amount: {
      amount: backendBudget.amount,
      currency: 'JPY',
    },
    periodType: backendBudget.period as PeriodType,
    startDate: backendBudget.start_date,
    endDate: backendBudget.end_date,
    isActive: backendBudget.is_active,
    createdAt: backendBudget.created_at,
    updatedAt: backendBudget.updated_at,
  };

  // Calculate usage (for now, used amount is 0)
  // TODO: Fetch actual transaction data to calculate used amount
  return calculateBudgetUsage(budget, { amount: 0, currency: 'JPY' });
};

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
    const response = await apiClient.get<BackendBudgetListResponse>(url);

    // Transform backend response to frontend type
    const budgets = response.data.budgets.map(transformBudgetResponse);

    return {
      budgets,
      total: response.data.total_count,
    };
  },

  get: async (id: string): Promise<BudgetWithUsage> => {
    const response = await apiClient.get<BudgetWithUsage>(endpoints.get(id));
    return response.data;
  },

  create: async (payload: CreateBudgetPayload): Promise<Budget> => {
    const backendPayload = {
      category_id: payload.categoryId,
      amount: payload.amount.amount,
      period: payload.periodType,
      start_date: payload.startDate,
      end_date: payload.endDate,
    };
    const response = await apiClient.post<BackendBudgetResponse>(
      endpoints.create,
      backendPayload,
    );
    return {
      id: response.data.id,
      userId: response.data.user_id,
      categoryId: response.data.category_id,
      amount: {
        amount: response.data.amount,
        currency: 'JPY',
      },
      periodType: response.data.period as PeriodType,
      startDate: response.data.start_date,
      endDate: response.data.end_date,
      isActive: response.data.is_active,
      createdAt: response.data.created_at,
      updatedAt: response.data.updated_at,
    };
  },

  update: async (id: string, payload: UpdateBudgetPayload): Promise<Budget> => {
    const backendPayload: any = {};
    if (payload.amount !== undefined) {
      backendPayload.amount = payload.amount.amount;
    }
    if (payload.startDate !== undefined) {
      backendPayload.start_date = payload.startDate;
    }
    if (payload.endDate !== undefined) {
      backendPayload.end_date = payload.endDate;
    }
    if (payload.isActive !== undefined) {
      backendPayload.is_active = payload.isActive;
    }

    const response = await apiClient.put<BackendBudgetResponse>(
      endpoints.update(id),
      backendPayload,
    );
    return {
      id: response.data.id,
      userId: response.data.user_id,
      categoryId: response.data.category_id,
      amount: {
        amount: response.data.amount,
        currency: 'JPY',
      },
      periodType: response.data.period as PeriodType,
      startDate: response.data.start_date,
      endDate: response.data.end_date,
      isActive: response.data.is_active,
      createdAt: response.data.created_at,
      updatedAt: response.data.updated_at,
    };
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(endpoints.delete(id));
  },

  getCurrent: async (): Promise<BudgetListResponse> => {
    const response = await apiClient.get<BackendBudgetListResponse>(
      endpoints.current,
    );
    const budgets = response.data.budgets.map(transformBudgetResponse);
    return { budgets, total: response.data.total_count };
  },

  getSummary: async (
    period: 'monthly' | 'yearly' = 'monthly',
  ): Promise<BudgetSummary> => {
    const params = new URLSearchParams();
    params.append('period', period);

    const url = `${endpoints.summary}?${params.toString()}`;
    const response = await apiClient.get<BackendBudgetSummaryResponse>(url);

    return {
      totalBudget: { amount: response.data.total_budget, currency: 'JPY' },
      totalUsed: { amount: response.data.total_used, currency: 'JPY' },
      totalRemaining: {
        amount: response.data.total_remaining,
        currency: 'JPY',
      },
      overBudgetCount: 0, // TODO: バックエンドから取得
      activeCount: 0, // TODO: バックエンドから取得
    };
  },
};
