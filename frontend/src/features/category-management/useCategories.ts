import { useQuery } from '@tanstack/react-query';
import { categoryApi, categoryKeys } from '@/entities/category';

export const useCategories = () => {
  return useQuery({
    queryKey: categoryKeys.lists(),
    queryFn: categoryApi.list,
    staleTime: 30 * 60 * 1000, // 30分
  });
};