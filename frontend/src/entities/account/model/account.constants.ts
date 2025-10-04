import type { AccountType } from './account.types';

export const ACCOUNT_TYPE_LABELS: Record<AccountType, string> = {
  checking: '普通預金',
  investment: '投資',
  cash: '現金',
  credit_card: 'クレジットカード',
} as const;

export const ACCOUNT_TYPE_ICONS: Record<AccountType, string> = {
  checking: '🏦',
  investment: '📈',
  cash: '💵',
  credit_card: '💳',
} as const;

export const ACCOUNT_LIMITS = {
  NAME_MIN_LENGTH: 1,
  NAME_MAX_LENGTH: 100,
  INITIAL_BALANCE_MIN: 0,
  INITIAL_BALANCE_MAX: 999999999999,
} as const;

export const ACCOUNT_TYPE_OPTIONS = [
  { value: 'checking', label: ACCOUNT_TYPE_LABELS.checking },
  { value: 'investment', label: ACCOUNT_TYPE_LABELS.investment },
  { value: 'cash', label: ACCOUNT_TYPE_LABELS.cash },
  { value: 'credit_card', label: ACCOUNT_TYPE_LABELS.credit_card },
] as const;
