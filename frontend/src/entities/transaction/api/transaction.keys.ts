import type { TransactionListParams, MonthlySummaryParams } from '../model';

export const transactionKeys = {
  all: ['transactions'] as const,
  lists: () => [...transactionKeys.all, 'list'] as const,
  list: (params: TransactionListParams = {}) =>
    [...transactionKeys.lists(), params] as const,
  details: () => [...transactionKeys.all, 'detail'] as const,
  detail: (id: string) => [...transactionKeys.details(), id] as const,
  summaries: () => [...transactionKeys.all, 'summary'] as const,
  monthlySummary: (params: MonthlySummaryParams = {}) =>
    [...transactionKeys.summaries(), 'monthly', params] as const,
} as const;
