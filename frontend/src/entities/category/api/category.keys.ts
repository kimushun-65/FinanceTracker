import type { CategoryType } from '../model';

export const categoryKeys = {
  all: ['categories'] as const,
  lists: () => [...categoryKeys.all, 'list'] as const,
  list: (type?: CategoryType) => [...categoryKeys.lists(), type] as const,
  details: () => [...categoryKeys.all, 'detail'] as const,
  detail: (id: string) => [...categoryKeys.details(), id] as const,
} as const;

export const categoryMasterKeys = {
  all: ['categoryMasters'] as const,
  lists: () => [...categoryMasterKeys.all, 'list'] as const,
  list: (type?: CategoryType) => [...categoryMasterKeys.lists(), type] as const,
  details: () => [...categoryMasterKeys.all, 'detail'] as const,
  detail: (id: string) => [...categoryMasterKeys.details(), id] as const,
} as const;
