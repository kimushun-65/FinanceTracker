import { useQuery } from '@tanstack/react-query';
import { categoryApi, categoryKeys } from '@/entities/category';

export const useCategory = (id: string) => {
  return useQuery({
    queryKey: categoryKeys.detail(id),
    queryFn: () => categoryApi.get(id),
    staleTime: 30 * 60 * 1000, // 30分
    enabled: !!id,
  });
};