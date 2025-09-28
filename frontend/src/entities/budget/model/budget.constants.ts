import type { PeriodType, BudgetStatus } from './budget.types';

export const PERIOD_TYPES: Record<PeriodType, string> = {
  monthly: '月間',
  yearly: '年間',
} as const;

export const BUDGET_STATUS: Record<BudgetStatus, string> = {
  normal: '正常',
  warning: '注意',
  exceeded: '超過',
} as const;

export const BUDGET_STATUS_COLORS: Record<BudgetStatus, string> = {
  normal: '#10b981',
  warning: '#f59e0b',
  exceeded: '#ef4444',
} as const;

export const BUDGET_THRESHOLDS = {
  WARNING_PERCENTAGE: 80,
  EXCEEDED_PERCENTAGE: 100,
} as const;

export const BUDGET_VALIDATION = {
  maxAmount: 100000000, // 1億円
  minAmount: 1,
  maxPeriodDays: 365,
  minPeriodDays: 1,
} as const;
