import { ACCOUNT_LIMITS } from '../model';

export const validateAccountName = (name: string): boolean => {
  return (
    name.length >= ACCOUNT_LIMITS.NAME_MIN_LENGTH &&
    name.length <= ACCOUNT_LIMITS.NAME_MAX_LENGTH
  );
};

export const validateAccountNameError = (name: string): string | null => {
  if (name.length < ACCOUNT_LIMITS.NAME_MIN_LENGTH) {
    return `口座名は${ACCOUNT_LIMITS.NAME_MIN_LENGTH}文字以上で入力してください`;
  }
  if (name.length > ACCOUNT_LIMITS.NAME_MAX_LENGTH) {
    return `口座名は${ACCOUNT_LIMITS.NAME_MAX_LENGTH}文字以内で入力してください`;
  }
  return null;
};

export const validateInitialBalance = (amount: number): boolean => {
  return (
    amount >= ACCOUNT_LIMITS.INITIAL_BALANCE_MIN &&
    amount <= ACCOUNT_LIMITS.INITIAL_BALANCE_MAX
  );
};

export const validateInitialBalanceError = (amount: number): string | null => {
  if (amount < ACCOUNT_LIMITS.INITIAL_BALANCE_MIN) {
    return '初期残高は0以上で入力してください';
  }
  if (amount > ACCOUNT_LIMITS.INITIAL_BALANCE_MAX) {
    return '初期残高が大きすぎます';
  }
  return null;
};
