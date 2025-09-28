import type { Money } from '@/shared/value-objects';
import type {
  CreateTransactionPayload,
  UpdateTransactionPayload,
} from '../model';

export const validateTransactionAmount = (amount: Money): boolean => {
  return amount.amount > 0;
};

export const validateTransactionDescription = (
  description: string,
): boolean => {
  return description.length >= 1 && description.length <= 500;
};

export const validateTransactionDate = (date: string): boolean => {
  const transactionDate = new Date(date);
  const now = new Date();
  const oneYearAgo = new Date();
  oneYearAgo.setFullYear(now.getFullYear() - 1);

  return transactionDate >= oneYearAgo && transactionDate <= now;
};

export const validateCreateTransactionPayload = (
  payload: CreateTransactionPayload,
): string[] => {
  const errors: string[] = [];

  if (!validateTransactionAmount(payload.amount)) {
    errors.push('金額は0より大きい値を入力してください');
  }

  if (!validateTransactionDescription(payload.description)) {
    errors.push('説明は1文字以上500文字以下で入力してください');
  }

  if (payload.date && !validateTransactionDate(payload.date)) {
    errors.push('取引日は1年前から今日までの範囲で入力してください');
  }

  return errors;
};

export const validateUpdateTransactionPayload = (
  payload: UpdateTransactionPayload,
): string[] => {
  const errors: string[] = [];

  if (payload.amount && !validateTransactionAmount(payload.amount)) {
    errors.push('金額は0より大きい値を入力してください');
  }

  if (
    payload.description &&
    !validateTransactionDescription(payload.description)
  ) {
    errors.push('説明は1文字以上500文字以下で入力してください');
  }

  if (payload.date && !validateTransactionDate(payload.date)) {
    errors.push('取引日は1年前から今日までの範囲で入力してください');
  }

  return errors;
};
