import type { Budget } from '../model';

export const isBudgetActive = (budget: Budget): boolean => {
  return budget.isActive;
};

export const isBudgetInPeriod = (
  budget: Budget,
  date: Date = new Date(),
): boolean => {
  const startDate = new Date(budget.startDate);
  const endDate = budget.endDate ? new Date(budget.endDate) : null;

  if (date < startDate) {
    return false;
  }

  if (endDate && date > endDate) {
    return false;
  }

  return true;
};

export const isBudgetValid = (budget: Budget): boolean => {
  return isBudgetActive(budget) && isBudgetInPeriod(budget);
};

export const isBudgetExpired = (budget: Budget): boolean => {
  if (!budget.endDate) return false;

  const endDate = new Date(budget.endDate);
  const now = new Date();

  return now > endDate;
};

export const isBudgetUpcoming = (budget: Budget): boolean => {
  const startDate = new Date(budget.startDate);
  const now = new Date();

  return now < startDate;
};

export const getBudgetPeriodStatus = (
  budget: Budget,
): 'upcoming' | 'active' | 'expired' => {
  if (isBudgetExpired(budget)) {
    return 'expired';
  } else if (isBudgetUpcoming(budget)) {
    return 'upcoming';
  } else {
    return 'active';
  }
};

export const canEditBudget = (budget: Budget): boolean => {
  // 期限切れの予算は編集不可
  if (isBudgetExpired(budget)) {
    return false;
  }

  return true;
};

export const canDeleteBudget = (budget: Budget): boolean => {
  // 予算は常に削除可能（履歴保持のため）
  return true;
};

export const shouldShowBudgetWarning = (budget: Budget): boolean => {
  // 期限が近い（7日以内）場合に警告を表示
  if (!budget.endDate) return false;

  const endDate = new Date(budget.endDate);
  const now = new Date();
  const daysRemaining = Math.ceil(
    (endDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24),
  );

  return daysRemaining <= 7 && daysRemaining > 0;
};
