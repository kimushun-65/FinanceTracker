import { useQuery } from '@tanstack/react-query';
import { categoryApi, categoryKeys } from '@/entities/category';

export const useCategoryMasters = () => {
  return useQuery({
    queryKey: categoryKeys.masters(),
    queryFn: categoryApi.listMasters,
    staleTime: 60 * 60 * 1000, // 1時間
  });
};