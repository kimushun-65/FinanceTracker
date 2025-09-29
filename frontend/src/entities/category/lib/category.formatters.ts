import type { Category } from '../model';

export type SelectOption = {
  value: string;
  label: string;
};

/**
 * カテゴリーをSelectオプションに変換
 */
export const categoriesToSelectOptions = (
  categories: Category[],
): SelectOption[] => {
  return (
    categories?.map((category: any) => ({
      value: category.id,
      label: category.name || category.displayName || 'Unknown',
    })) || []
  );
};
