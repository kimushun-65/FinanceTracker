import type { CreateCategoryPayload, UpdateCategoryPayload } from '../model';
import { CATEGORY_VALIDATION } from '../model';

export const validateCategoryName = (name: string): boolean => {
  return name.length >= 1 && name.length <= CATEGORY_VALIDATION.maxNameLength;
};

export const validateCustomName = (customName?: string): boolean => {
  if (!customName) return true;
  return (
    customName.length >= 1 &&
    customName.length <= CATEGORY_VALIDATION.maxCustomNameLength
  );
};

export const validateHexColor = (color: string): boolean => {
  return CATEGORY_VALIDATION.colorRegex.test(color);
};

export const validateCreateCategoryPayload = (
  payload: CreateCategoryPayload,
): string[] => {
  const errors: string[] = [];

  if (!payload.categoryMasterId) {
    errors.push('カテゴリマスターIDは必須です');
  }

  if (payload.customName && !validateCustomName(payload.customName)) {
    errors.push(
      `カスタム名は1文字以上${CATEGORY_VALIDATION.maxCustomNameLength}文字以下で入力してください`,
    );
  }

  return errors;
};

export const validateUpdateCategoryPayload = (
  payload: UpdateCategoryPayload,
): string[] => {
  const errors: string[] = [];

  if (
    payload.customName !== undefined &&
    !validateCustomName(payload.customName)
  ) {
    errors.push(
      `カスタム名は1文字以上${CATEGORY_VALIDATION.maxCustomNameLength}文字以下で入力してください`,
    );
  }

  return errors;
};
