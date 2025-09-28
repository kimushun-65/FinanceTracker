import { useQuery } from '@tanstack/react-query';
import { categoryMasterApi, categoryMasterKeys } from '@/entities/category';

export const useCategoryMasters = () => {
  return useQuery({
    queryKey: categoryMasterKeys.lists(),
    queryFn: () => categoryMasterApi.list(),
    staleTime: 60 * 60 * 1000, // 1時間
  });
};
