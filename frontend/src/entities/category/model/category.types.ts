import type { BaseEntity } from '@/shared/types';

export type CategoryType = 'income' | 'expense';

export type Category = BaseEntity & {
  userId: string;
  categoryMasterId: string;
  customName?: string;
  isActive: boolean;
};

export type CategoryMaster = BaseEntity & {
  name: string;
  type: CategoryType;
  icon: string;
  color: string;
  displayOrder: number;
};

export type CategoryWithMaster = Category & {
  master: CategoryMaster;
  displayName: string;
};

export type CreateCategoryPayload = {
  categoryMasterId: string;
  customName?: string;
};

export type UpdateCategoryPayload = {
  customName?: string;
  isActive?: boolean;
};

export type CategoryListResponse = {
  categories: CategoryWithMaster[];
  total: number;
};

export type CategoryMasterListResponse = {
  categoryMasters: CategoryMaster[];
  total: number;
};
