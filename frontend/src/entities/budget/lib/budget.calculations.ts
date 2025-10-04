import type { Money } from '@/shared/value-objects';
import type { Budget, BudgetWithUsage, BudgetStatus } from '../model';
import { BUDGET_THRESHOLDS } from '../model';
import { addMoney, subtractMoney } from '@/shared/value-objects';

export const calculateBudgetUsage = (
  budget: Budget,
  usedAmount: Money,
): BudgetWithUsage => {
  const usagePercentage = calculateUsagePercentage(budget.amount, usedAmount);
  const remaining = calculateRemainingAmount(budget.amount, usedAmount);
  const status = determineBudgetStatus(usagePercentage);

  return {
    ...budget,
    used: usedAmount,
    remaining,
    usagePercentage,
    status,
  };
};

export const calculateUsagePercentage = (
  budgetAmount: Money,
  usedAmount: Money,
): number => {
  if (budgetAmount.amount === 0) return 0;

  const percentage = (usedAmount.amount / budgetAmount.amount) * 100;
  return Math.round(percentage * 100) / 100; // 小数点第2位まで
};

export const calculateRemainingAmount = (
  budgetAmount: Money,
  usedAmount: Money,
): Money => {
  try {
    return subtractMoney(budgetAmount, usedAmount);
  } catch {
    // 通貨が一致しない場合は0を返す
    return { amount: 0, currency: budgetAmount.currency };
  }
};

export const determineBudgetStatus = (
  usagePercentage: number,
): BudgetStatus => {
  if (usagePercentage >= BUDGET_THRESHOLDS.EXCEEDED_PERCENTAGE) {
    return 'exceeded';
  } else if (usagePercentage >= BUDGET_THRESHOLDS.WARNING_PERCENTAGE) {
    return 'warning';
  } else {
    return 'normal';
  }
};

export const calculateTotalBudgetAmount = (budgets: Budget[]): Money => {
  if (budgets.length === 0) {
    return { amount: 0, currency: 'JPY' };
  }

  return budgets.reduce(
    (total, budget) => {
      try {
        return addMoney(total, budget.amount);
      } catch {
        // 通貨が一致しない場合はスキップ
        return total;
      }
    },
    { amount: 0, currency: budgets[0].amount.currency },
  );
};

export const calculateTotalUsedAmount = (
  budgetsWithUsage: BudgetWithUsage[],
): Money => {
  if (budgetsWithUsage.length === 0) {
    return { amount: 0, currency: 'JPY' };
  }

  // 最初の有効な通貨を見つける
  const firstValidCurrency =
    budgetsWithUsage.find((b) => b.used?.currency)?.used?.currency || 'JPY';

  return budgetsWithUsage.reduce(
    (total, budget) => {
      // budget.usedが存在し、currencyが設定されている場合のみ加算
      if (!budget.used || !budget.used.currency) {
        return total;
      }

      try {
        return addMoney(total, budget.used);
      } catch {
        // 通貨が一致しない場合はスキップ
        return total;
      }
    },
    { amount: 0, currency: firstValidCurrency },
  );
};

export const calculateBudgetProgress = (
  budgetAmount: Money,
  usedAmount: Money,
): number => {
  if (!budgetAmount || budgetAmount.amount === 0) return 0;
  if (!usedAmount || usedAmount.amount === undefined) return 0;

  const progress = Math.min(usedAmount.amount / budgetAmount.amount, 1);
  return Math.round(progress * 100);
};

export const isOverBudget = (budget: BudgetWithUsage): boolean => {
  return budget.status === 'exceeded';
};

export const isNearBudgetLimit = (budget: BudgetWithUsage): boolean => {
  return budget.status === 'warning' || budget.status === 'exceeded';
};

export const getDaysRemainingInPeriod = (endDate?: string): number | null => {
  if (!endDate) return null;

  const end = new Date(endDate);
  const now = new Date();
  const diffInTime = end.getTime() - now.getTime();
  const diffInDays = Math.ceil(diffInTime / (1000 * 60 * 60 * 24));

  return Math.max(0, diffInDays);
};
