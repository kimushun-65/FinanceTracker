import type { Transaction, TransactionWithDetails } from '../model';
import type { Money } from '@/shared/value-objects';

export const transformTransactionForDisplay = (
  transaction: Transaction,
): TransactionWithDetails => {
  return {
    ...transaction,
    displayAmount: formatTransactionAmount(
      transaction.amount,
      transaction.type,
    ),
    displayDate: formatTransactionDate(transaction.date),
  };
};

export const formatTransactionAmount = (
  amount: Money,
  type: 'income' | 'expense',
): string => {
  const sign = type === 'income' ? '+' : '-';
  const formattedAmount =
    amount.currency === 'JPY'
      ? `¥${amount.amount.toLocaleString('ja-JP')}`
      : `${amount.currency} ${amount.amount}`;

  return `${sign}${formattedAmount}`;
};

export const formatTransactionDate = (date: string): string => {
  const transactionDate = new Date(date);
  return transactionDate.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
};

export const formatTransactionDateTime = (date: string): string => {
  const transactionDate = new Date(date);
  return transactionDate.toLocaleString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
};

export const transformApiTransactionResponse = (
  apiTransaction: any,
): Transaction => {
  return {
    id: apiTransaction.id,
    createdAt: apiTransaction.createdAt,
    updatedAt: apiTransaction.updatedAt,
    userId: apiTransaction.userId,
    accountId: apiTransaction.accountId,
    categoryId: apiTransaction.categoryId,
    amount: {
      amount: apiTransaction.amount,
      currency: apiTransaction.currency || 'JPY',
    },
    type: apiTransaction.type,
    description: apiTransaction.description,
    date: apiTransaction.date,
  };
};
