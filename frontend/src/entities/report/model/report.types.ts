import type { Money } from '@/shared/value-objects';

export type TransactionType = 'income' | 'expense';

export interface CategorySummary {
  categoryId: string;
  categoryName: string;
  totalAmount: Money;
  transactionCount: number;
  percentage: number;
  type: TransactionType;
}

export interface MonthlyReport {
  totalIncome: Money;
  totalExpenses: Money;
  netIncome: Money;
  transactionCount: number;
  averageTransactionAmount: Money;
  topCategories: CategorySummary[];
}

export interface PeriodData {
  income: Money;
  expenses: Money;
  netIncome: Money;
  transactionCount: number;
}

export interface PeriodComparison {
  current: PeriodData;
  previous: PeriodData;
  change: {
    income: Money;
    expenses: Money;
    netIncome: Money;
  };
  percentageChange: {
    income: number;
    expenses: number;
    netIncome: number;
  };
}

export interface TrendData {
  date: string;
  [categoryId: string]: number | string;
}

export interface DateRange {
  from: Date;
  to: Date;
}

export type ComparisonPeriod = 'month' | 'quarter' | 'year';
