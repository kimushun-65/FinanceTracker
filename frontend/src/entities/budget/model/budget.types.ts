import type { BaseEntity } from '@/shared/types';
import type { Money } from '@/shared/value-objects';

export type PeriodType = 'monthly' | 'yearly';

export type Budget = BaseEntity & {
  userId: string;
  categoryId: string;
  amount: Money;
  periodType: PeriodType;
  startDate: string;
  endDate?: string;
  isActive: boolean;
};

export type CreateBudgetPayload = {
  categoryId: string;
  amount: Money;
  periodType: PeriodType;
  startDate: string;
  endDate?: string;
};

export type UpdateBudgetPayload = {
  amount?: Money;
  startDate?: string;
  endDate?: string;
  isActive?: boolean;
};

export type BudgetWithUsage = Budget & {
  used: Money;
  remaining: Money;
  usagePercentage: number;
  status: 'normal' | 'warning' | 'exceeded';
};

export type BudgetListResponse = {
  budgets: BudgetWithUsage[];
  total: number;
};

export type BudgetStatus = 'normal' | 'warning' | 'exceeded';

export type BudgetPeriod = {
  startDate: string;
  endDate: string;
  isActive: boolean;
};

export type BudgetSummary = {
  totalBudget: Money;
  totalUsed: Money;
  totalRemaining: Money;
  overBudgetCount: number;
  activeCount: number;
};
