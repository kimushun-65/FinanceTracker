import type { BaseEntity } from '@/shared/types';
import type { Money } from '@/shared/value-objects';

export type TransactionType = 'income' | 'expense';

export type Transaction = BaseEntity & {
  userId: string;
  accountId: string;
  categoryId: string;
  amount: Money;
  type: TransactionType;
  description: string;
  date: string;
};

export type CreateTransactionPayload = {
  accountId: string;
  categoryId: string;
  amount: Money;
  type: TransactionType;
  description: string;
  date?: string;
};

export type UpdateTransactionPayload = {
  amount?: Money;
  description?: string;
  date?: string;
  categoryId?: string;
  accountId?: string;
};

export type TransactionListParams = {
  accountId?: string;
  categoryId?: string;
  type?: TransactionType;
  startDate?: string;
  endDate?: string;
  limit?: number;
  offset?: number;
  sortBy?: 'date' | 'amount';
  sortOrder?: 'asc' | 'desc';
};

export type TransactionListResponse = {
  transactions: Transaction[];
  total: number;
  hasMore: boolean;
};

export type TransactionAggregate = {
  totalIncome: Money;
  totalExpense: Money;
  netAmount: Money;
  transactionCount: number;
};

export type TransactionCategoryAggregate = {
  categoryId: string;
  categoryName: string;
  totalAmount: Money;
  transactionCount: number;
  percentage: number;
};

export type MonthlyAggregate = {
  month: string;
  income: Money;
  expense: Money;
  net: Money;
  transactionCount: number;
};

export type CategoryAggregate = {
  categoryId: string;
  income: Money;
  expense: Money;
  net: Money;
  transactionCount: number;
};

export type TransactionWithDetails = Transaction & {
  displayAmount: string;
  displayDate: string;
};
