import type { Transaction, TransactionType } from '../model';
import type { Money } from '@/shared/value-objects';
import { formatMoney as formatMoneyVO } from '@/shared/value-objects/money';

export const formatTransactionType = (type: TransactionType): string => {
  const typeMap = {
    income: '収入',
    expense: '支出',
  };
  return typeMap[type];
};

export const formatTransactionSummary = (transaction: Transaction): string => {
  const typeText = formatTransactionType(transaction.type);
  const amountText = formatMoney(transaction.amount);
  return `${typeText}: ${amountText} - ${transaction.description}`;
};

// Safe money formatting that tolerates undefined or string amounts
export const formatMoney = (money: Money): string => {
  const amount =
    typeof money?.amount === 'number'
      ? money.amount
      : Number((money as any)?.amount ?? 0);
  const currency = (money && (money as any).currency) || 'JPY';
  return formatMoneyVO({ amount, currency } as Money);
};

export const formatMoneyWithSign = (
  money: Money,
  type: TransactionType,
): string => {
  const sign = type === 'income' ? '+' : '-';
  return `${sign}${formatMoney(money)}`;
};

export const formatTransactionPeriod = (
  startDate: string,
  endDate: string,
): string => {
  const start = new Date(startDate);
  const end = new Date(endDate);

  const startText = start.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });

  const endText = end.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });

  return `${startText} 〜 ${endText}`;
};

export const formatMonthYear = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: 'long',
  });
};

export const formatRelativeDate = (dateString: string): string => {
  const date = new Date(dateString);
  const now = new Date();
  const diffInDays = Math.floor(
    (now.getTime() - date.getTime()) / (1000 * 60 * 60 * 24),
  );

  if (diffInDays === 0) return '今日';
  if (diffInDays === 1) return '昨日';
  if (diffInDays < 7) return `${diffInDays}日前`;
  if (diffInDays < 30) return `${Math.floor(diffInDays / 7)}週間前`;
  if (diffInDays < 365) return `${Math.floor(diffInDays / 30)}ヶ月前`;

  return `${Math.floor(diffInDays / 365)}年前`;
};

export type SelectOption = {
  value: string;
  label: string;
};

/**
 * トランザクションタイプのオプション
 */
export const transactionTypeOptions: SelectOption[] = [
  { value: 'expense', label: 'Expense' },
  { value: 'income', label: 'Income' },
];
