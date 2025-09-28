import { useQuery } from '@tanstack/react-query';
import { budgetApi, budgetKeys } from '@/entities/budget';

export const useBudgets = () => {
  return useQuery({
    queryKey: budgetKeys.lists(),
    queryFn: () => budgetApi.list(),
    staleTime: 10 * 60 * 1000, // 10分
  });
};
