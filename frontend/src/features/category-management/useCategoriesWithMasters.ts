import { useQuery } from '@tanstack/react-query';
import {
  categoryApi,
  categoryMasterApi,
  categoryKeys,
  transformCategoryWithMaster,
  type CategoryWithMaster,
} from '@/entities/category';

export const useCategoriesWithMasters = () => {
  return useQuery({
    queryKey: [...categoryKeys.lists(), 'with-masters'],
    queryFn: async () => {
      // Fetch both categories and masters in parallel
      const [categoriesResponse, mastersResponse] = await Promise.all([
        categoryApi.list(),
        categoryMasterApi.list(),
      ]);

      // Create a map of masters for quick lookup
      const mastersMap = new Map(
        mastersResponse.categoryMasters.map((master) => [master.id, master]),
      );

      // Check if categoriesResponse has categories property or is an array
      const categoriesArray = Array.isArray(categoriesResponse)
        ? categoriesResponse
        : categoriesResponse.categories || [];

      // Transform categories with their masters
      const categoriesWithMasters: CategoryWithMaster[] = categoriesArray.map(
        (category) => {
          const master = mastersMap.get(
            category.categoryMasterId || category.category_master_id,
          );
          if (!master) {
            // Fallback if master not found
            return {
              ...category,
              master: {
                id: category.categoryMasterId,
                name: category.name || 'Unknown',
                type: 'expense' as const,
                icon: '📦',
                color: '#gray',
                displayOrder: 0,
                createdAt: category.createdAt,
                updatedAt: category.updatedAt,
              },
              displayName: category.customName || category.name || 'Unknown',
            };
          }
          return transformCategoryWithMaster(category, master);
        },
      );

      return categoriesWithMasters;
    },
    staleTime: 30 * 60 * 1000, // 30分
  });
};
