export const accountKeys = {
  all: ['accounts'] as const,
  lists: () => [...accountKeys.all, 'list'] as const,
  list: (filters?: Record<string, unknown>) =>
    [...accountKeys.lists(), filters] as const,
  details: () => [...accountKeys.all, 'detail'] as const,
  detail: (id: string) => [...accountKeys.details(), id] as const,
  balance: (id: string) => [...accountKeys.detail(id), 'balance'] as const,
  transactions: (id: string, filters?: Record<string, unknown>) =>
    [...accountKeys.detail(id), 'transactions', filters] as const,
} as const;

export type AccountKeys = typeof accountKeys;
