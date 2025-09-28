import type { CategoryWithMaster, CategoryType } from '../model';
import { CATEGORY_TYPES } from '../model';

export const formatCategoryType = (type: CategoryType): string => {
  return CATEGORY_TYPES[type];
};

export const getCategoryDisplayName = (
  category: CategoryWithMaster,
): string => {
  return category.customName || category.master.name;
};

export const sortCategoriesByDisplayOrder = (
  categories: CategoryWithMaster[],
): CategoryWithMaster[] => {
  return [...categories].sort(
    (a, b) => a.master.displayOrder - b.master.displayOrder,
  );
};

export const sortCategoriesByName = (
  categories: CategoryWithMaster[],
): CategoryWithMaster[] => {
  return [...categories].sort((a, b) => {
    const nameA = getCategoryDisplayName(a);
    const nameB = getCategoryDisplayName(b);
    return nameA.localeCompare(nameB, 'ja-JP');
  });
};

export const filterCategoriesByType = (
  categories: CategoryWithMaster[],
  type: CategoryType,
): CategoryWithMaster[] => {
  return categories.filter((category) => category.master.type === type);
};

export const filterActiveCategories = (
  categories: CategoryWithMaster[],
): CategoryWithMaster[] => {
  return categories.filter((category) => category.isActive);
};

export const groupCategoriesByType = (
  categories: CategoryWithMaster[],
): Record<CategoryType, CategoryWithMaster[]> => {
  return categories.reduce(
    (groups, category) => {
      const type = category.master.type;
      if (!groups[type]) {
        groups[type] = [];
      }
      groups[type].push(category);
      return groups;
    },
    {} as Record<CategoryType, CategoryWithMaster[]>,
  );
};

export const findCategoryById = (
  categories: CategoryWithMaster[],
  id: string,
): CategoryWithMaster | undefined => {
  return categories.find((category) => category.id === id);
};

export const getCategoryIcon = (category: CategoryWithMaster): string => {
  return category.master.icon;
};

export const getCategoryColor = (category: CategoryWithMaster): string => {
  return category.master.color;
};
