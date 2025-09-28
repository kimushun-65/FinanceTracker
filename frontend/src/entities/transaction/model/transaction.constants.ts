import type { TransactionType } from './transaction.types';

export const TRANSACTION_TYPE_LABELS: Record<TransactionType, string> = {
  income: '収入',
  expense: '支出',
} as const;

export const TRANSACTION_TYPE_COLORS: Record<TransactionType, string> = {
  income: 'text-green-600',
  expense: 'text-red-600',
} as const;

export const TRANSACTION_TYPE_ICONS: Record<TransactionType, string> = {
  income: '💰',
  expense: '💸',
} as const;

export const TRANSACTION_LIMITS = {
  DESCRIPTION_MIN_LENGTH: 0,
  DESCRIPTION_MAX_LENGTH: 500,
  AMOUNT_MIN: 1,
  AMOUNT_MAX: 999999999999,
  DEFAULT_PAGE_SIZE: 20,
  MAX_PAGE_SIZE: 100,
} as const;

export const TRANSACTION_TYPE_OPTIONS = [
  { value: 'income', label: TRANSACTION_TYPE_LABELS.income },
  { value: 'expense', label: TRANSACTION_TYPE_LABELS.expense },
] as const;

export const TRANSACTION_SORT_OPTIONS = [
  { value: 'date', label: '日付' },
  { value: 'amount', label: '金額' },
] as const;

export const TRANSACTION_SORT_ORDER_OPTIONS = [
  { value: 'desc', label: '降順' },
  { value: 'asc', label: '昇順' },
] as const;
