import { useQuery } from '@tanstack/react-query';
import { budgetApi, budgetKeys } from '@/entities/budget';

export const useBudget = (id: string) => {
  return useQuery({
    queryKey: budgetKeys.detail(id),
    queryFn: () => budgetApi.get(id),
    staleTime: 10 * 60 * 1000, // 10分
    enabled: !!id,
  });
};
