import type { PeriodType } from '../model';

export const budgetKeys = {
  all: ['budgets'] as const,
  lists: () => [...budgetKeys.all, 'list'] as const,
  list: (period?: PeriodType, categoryId?: string) =>
    [...budgetKeys.lists(), { period, categoryId }] as const,
  details: () => [...budgetKeys.all, 'detail'] as const,
  detail: (id: string) => [...budgetKeys.details(), id] as const,
  summary: () => [...budgetKeys.all, 'summary'] as const,
} as const;
