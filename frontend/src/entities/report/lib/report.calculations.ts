import type {
  MonthlyReport,
  PeriodData,
  PeriodComparison,
  CategorySummary,
} from '../model';

// Transactionの最小限の型定義
interface Transaction {
  type: 'income' | 'expense';
  amount: {
    amount: number;
    currency: string;
  };
}

/**
 * 月次レポート計算
 */
export const calculateMonthlyReport = (
  transactions: Transaction[],
  categorySummary: CategorySummary[],
): MonthlyReport => {
  const income = transactions
    .filter((t) => t.type === 'income')
    .reduce((sum, t) => sum + t.amount.amount, 0);

  const expenses = transactions
    .filter((t) => t.type === 'expense')
    .reduce((sum, t) => sum + t.amount.amount, 0);

  const totalTransactionAmount = income + expenses;
  const avgAmount = transactions.length
    ? totalTransactionAmount / transactions.length
    : 0;

  return {
    totalIncome: { amount: income, currency: 'JPY' },
    totalExpenses: { amount: expenses, currency: 'JPY' },
    netIncome: { amount: income - expenses, currency: 'JPY' },
    transactionCount: transactions.length,
    averageTransactionAmount: {
      amount: Math.round(avgAmount),
      currency: 'JPY',
    },
    topCategories: categorySummary
      .sort((a, b) => b.totalAmount.amount - a.totalAmount.amount)
      .slice(0, 5),
  };
};

/**
 * 期間比較計算
 */
export const calculatePeriodComparison = (
  current: PeriodData,
  previous: PeriodData,
): PeriodComparison => {
  const change = {
    income: {
      amount: current.income.amount - previous.income.amount,
      currency: 'JPY' as const,
    },
    expenses: {
      amount: current.expenses.amount - previous.expenses.amount,
      currency: 'JPY' as const,
    },
    netIncome: {
      amount: current.netIncome.amount - previous.netIncome.amount,
      currency: 'JPY' as const,
    },
  };

  const percentageChange = {
    income: previous.income.amount
      ? (change.income.amount / previous.income.amount) * 100
      : 0,
    expenses: previous.expenses.amount
      ? (change.expenses.amount / previous.expenses.amount) * 100
      : 0,
    netIncome: previous.netIncome.amount
      ? (change.netIncome.amount / previous.netIncome.amount) * 100
      : 0,
  };

  return {
    current,
    previous,
    change,
    percentageChange,
  };
};

/**
 * PeriodData計算ヘルパー
 */
export const calculatePeriodData = (
  transactions: Transaction[],
): PeriodData => {
  const income = transactions
    .filter((t) => t.type === 'income')
    .reduce((sum, t) => sum + t.amount.amount, 0);

  const expenses = transactions
    .filter((t) => t.type === 'expense')
    .reduce((sum, t) => sum + t.amount.amount, 0);

  return {
    income: { amount: income, currency: 'JPY' },
    expenses: { amount: expenses, currency: 'JPY' },
    netIncome: { amount: income - expenses, currency: 'JPY' },
    transactionCount: transactions.length,
  };
};
