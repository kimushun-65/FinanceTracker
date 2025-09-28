import type {
  Transaction,
  MonthlyAggregate,
  CategoryAggregate,
  TransactionType,
} from '../model';
import type { Money } from '@/shared/value-objects';
import { addMoney } from '@/shared/value-objects';

export const groupTransactionsByMonth = (
  transactions: Transaction[],
): Map<string, Transaction[]> => {
  const grouped = new Map<string, Transaction[]>();

  transactions.forEach((transaction) => {
    const month = transaction.date.substring(0, 7); // YYYY-MM
    if (!grouped.has(month)) {
      grouped.set(month, []);
    }
    grouped.get(month)!.push(transaction);
  });

  return grouped;
};

export const calculateMonthlyAggregates = (
  transactions: Transaction[],
): MonthlyAggregate[] => {
  const grouped = groupTransactionsByMonth(transactions);
  const aggregates: MonthlyAggregate[] = [];

  grouped.forEach((monthTransactions, month) => {
    const income = calculateTotalByType(monthTransactions, 'income');
    const expense = calculateTotalByType(monthTransactions, 'expense');
    const net = subtractMoney(income, expense);

    aggregates.push({
      month,
      income,
      expense,
      net,
      transactionCount: monthTransactions.length,
    });
  });

  return aggregates.sort((a, b) => b.month.localeCompare(a.month));
};

export const calculateCategoryAggregates = (
  transactions: Transaction[],
): CategoryAggregate[] => {
  const categoryMap = new Map<string, Transaction[]>();

  transactions.forEach((transaction) => {
    if (!categoryMap.has(transaction.categoryId)) {
      categoryMap.set(transaction.categoryId, []);
    }
    categoryMap.get(transaction.categoryId)!.push(transaction);
  });

  const aggregates: CategoryAggregate[] = [];

  categoryMap.forEach((categoryTransactions, categoryId) => {
    const income = calculateTotalByType(categoryTransactions, 'income');
    const expense = calculateTotalByType(categoryTransactions, 'expense');
    const net = subtractMoney(income, expense);

    aggregates.push({
      categoryId,
      income,
      expense,
      net,
      transactionCount: categoryTransactions.length,
    });
  });

  return aggregates.sort(
    (a, b) => Math.abs(b.net.amount) - Math.abs(a.net.amount),
  );
};

export const calculateTotalByType = (
  transactions: Transaction[],
  type: TransactionType,
): Money => {
  const filtered = transactions.filter((t) => t.type === type);
  if (filtered.length === 0) {
    return { amount: 0, currency: 'JPY' };
  }

  return filtered.reduce((total, transaction) => {
    if (total.currency === transaction.amount.currency) {
      return addMoney(total, transaction.amount);
    }
    return total;
  }, filtered[0].amount);
};

export const calculateTotalAmount = (transactions: Transaction[]): Money => {
  if (transactions.length === 0) {
    return { amount: 0, currency: 'JPY' };
  }

  let total = { amount: 0, currency: transactions[0].amount.currency };

  transactions.forEach((transaction) => {
    if (transaction.type === 'income') {
      total = addMoney(total, transaction.amount);
    } else {
      total = subtractMoney(total, transaction.amount);
    }
  });

  return total;
};

const subtractMoney = (a: Money, b: Money): Money => {
  if (a.currency !== b.currency) {
    throw new Error('Currency mismatch');
  }
  return { amount: a.amount - b.amount, currency: a.currency };
};
