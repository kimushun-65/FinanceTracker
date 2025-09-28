import type { Category, CategoryMaster, CategoryWithMaster } from '../model';

export const transformCategoryWithMaster = (
  category: Category,
  master: CategoryMaster,
): CategoryWithMaster => {
  return {
    ...category,
    master,
    displayName: category.customName || master.name,
  };
};

export const transformApiCategoryResponse = (apiCategory: any): Category => {
  return {
    id: apiCategory.id,
    createdAt: apiCategory.createdAt,
    updatedAt: apiCategory.updatedAt,
    userId: apiCategory.userId,
    categoryMasterId: apiCategory.categoryMasterId,
    customName: apiCategory.customName || undefined,
    isActive: apiCategory.isActive,
  };
};

export const transformApiCategoryMasterResponse = (
  apiCategoryMaster: any,
): CategoryMaster => {
  return {
    id: apiCategoryMaster.id,
    createdAt: apiCategoryMaster.createdAt,
    updatedAt: apiCategoryMaster.updatedAt,
    name: apiCategoryMaster.name,
    type: apiCategoryMaster.type,
    icon: apiCategoryMaster.icon,
    color: apiCategoryMaster.color,
    displayOrder: apiCategoryMaster.displayOrder,
  };
};

export const transformApiCategoryWithMasterResponse = (
  apiResponse: any,
): CategoryWithMaster => {
  const category = transformApiCategoryResponse(apiResponse);
  const master = transformApiCategoryMasterResponse(
    apiResponse.master || apiResponse.categoryMaster,
  );

  return transformCategoryWithMaster(category, master);
};
