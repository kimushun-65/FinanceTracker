import type { TransactionListParams } from '../model';

export const transactionKeys = {
  all: ['transactions'] as const,
  lists: () => [...transactionKeys.all, 'list'] as const,
  list: (params: TransactionListParams = {}) =>
    [...transactionKeys.lists(), params] as const,
  details: () => [...transactionKeys.all, 'detail'] as const,
  detail: (id: string) => [...transactionKeys.details(), id] as const,
} as const;
