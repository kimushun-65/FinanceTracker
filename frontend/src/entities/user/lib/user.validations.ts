import type { UpdateUserPayload, CreateUserPayload } from '../model';

export const validateUserName = (name: string): boolean => {
  return name.length >= 1 && name.length <= 100;
};

export const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

export const validateAuth0UserId = (auth0UserId: string): boolean => {
  return auth0UserId.length > 0 && auth0UserId.startsWith('auth0|');
};

export const validateUpdateUserPayload = (payload: UpdateUserPayload): string[] => {
  const errors: string[] = [];

  if (payload.name !== undefined && !validateUserName(payload.name)) {
    errors.push('名前は1文字以上100文字以下で入力してください');
  }

  if (payload.email !== undefined && !validateEmail(payload.email)) {
    errors.push('有効なメールアドレスを入力してください');
  }

  return errors;
};

export const validateCreateUserPayload = (payload: CreateUserPayload): string[] => {
  const errors: string[] = [];

  if (!validateAuth0UserId(payload.auth0UserId)) {
    errors.push('有効なAuth0ユーザーIDが必要です');
  }

  if (!validateEmail(payload.email)) {
    errors.push('有効なメールアドレスを入力してください');
  }

  if (!validateUserName(payload.name)) {
    errors.push('名前は1文字以上100文字以下で入力してください');
  }

  return errors;
};