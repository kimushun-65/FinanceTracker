import type { Transaction, TransactionWithDetails } from '../model';
import type { Money } from '@/shared/value-objects';
import { formatMoney as formatMoneyVO } from '@/shared/value-objects/money';

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
  const formattedAmount = formatMoneyVO({
    amount: typeof amount?.amount === 'number' ? amount.amount : Number((amount as any)?.amount ?? 0),
    currency: (amount as any)?.currency ?? 'JPY',
  });

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
  // Accepts both camelCase and snake_case from backend
  const amountRaw = apiTransaction.amount;
  const amountNum = typeof amountRaw === 'number' ? amountRaw : Number(String(amountRaw));
  return {
    id: apiTransaction.id,
    createdAt: apiTransaction.created_at ?? apiTransaction.createdAt,
    updatedAt: apiTransaction.updated_at ?? apiTransaction.updatedAt,
    userId: apiTransaction.user_id ?? apiTransaction.userId,
    accountId: apiTransaction.account_id ?? apiTransaction.accountId,
    categoryId: apiTransaction.category_id ?? apiTransaction.categoryId,
    amount: {
      amount: Number.isFinite(amountNum) ? Math.round(amountNum) : 0,
      currency: apiTransaction.currency || 'JPY',
    },
    type: apiTransaction.transaction_type ?? apiTransaction.type,
    description: apiTransaction.description,
    date: apiTransaction.transaction_date ?? apiTransaction.date,
  };
};
