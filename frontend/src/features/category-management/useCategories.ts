import { useQuery } from '@tanstack/react-query';
import { categoryApi, categoryKeys } from '@/entities/category';

export const useCategories = () => {
  return useQuery({
    queryKey: categoryKeys.lists(),
    queryFn: async () => {
      const response = await categoryApi.list();
      return response.categories; // categories 配列を直接返す
    },
    staleTime: 30 * 60 * 1000, // 30分
  });
};
