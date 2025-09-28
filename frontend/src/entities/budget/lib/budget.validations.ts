import type { Money } from '@/shared/value-objects';
import type {
  CreateBudgetPayload,
  UpdateBudgetPayload,
  PeriodType,
} from '../model';
import { BUDGET_VALIDATION } from '../model';

export const validateBudgetAmount = (amount: Money): boolean => {
  return (
    amount.amount >= BUDGET_VALIDATION.minAmount &&
    amount.amount <= BUDGET_VALIDATION.maxAmount
  );
};

export const validateBudgetPeriod = (
  startDate: string,
  endDate?: string,
): boolean => {
  const start = new Date(startDate);
  const end = endDate ? new Date(endDate) : null;

  if (end && start >= end) {
    return false;
  }

  if (end) {
    const diffInDays = Math.ceil(
      (end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24),
    );
    return (
      diffInDays >= BUDGET_VALIDATION.minPeriodDays &&
      diffInDays <= BUDGET_VALIDATION.maxPeriodDays
    );
  }

  return true;
};

export const validatePeriodType = (
  periodType: PeriodType,
  startDate: string,
  endDate?: string,
): boolean => {
  if (!endDate) return true;

  const start = new Date(startDate);
  const end = new Date(endDate);
  const diffInDays = Math.ceil(
    (end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24),
  );

  if (periodType === 'monthly') {
    return diffInDays <= 31;
  } else if (periodType === 'yearly') {
    return diffInDays <= 366;
  }

  return true;
};

export const validateCreateBudgetPayload = (
  payload: CreateBudgetPayload,
): string[] => {
  const errors: string[] = [];

  if (!payload.categoryId) {
    errors.push('カテゴリIDは必須です');
  }

  if (!validateBudgetAmount(payload.amount)) {
    errors.push(
      `予算額は${BUDGET_VALIDATION.minAmount}円以上${BUDGET_VALIDATION.maxAmount.toLocaleString()}円以下で入力してください`,
    );
  }

  if (!validateBudgetPeriod(payload.startDate, payload.endDate)) {
    errors.push('有効な期間を設定してください');
  }

  if (
    !validatePeriodType(payload.periodType, payload.startDate, payload.endDate)
  ) {
    errors.push('期間タイプと期間が一致しません');
  }

  return errors;
};

export const validateUpdateBudgetPayload = (
  payload: UpdateBudgetPayload,
): string[] => {
  const errors: string[] = [];

  if (payload.amount && !validateBudgetAmount(payload.amount)) {
    errors.push(
      `予算額は${BUDGET_VALIDATION.minAmount}円以上${BUDGET_VALIDATION.maxAmount.toLocaleString()}円以下で入力してください`,
    );
  }

  if (
    payload.startDate &&
    payload.endDate &&
    !validateBudgetPeriod(payload.startDate, payload.endDate)
  ) {
    errors.push('有効な期間を設定してください');
  }

  return errors;
};
