import { useQuery } from '@tanstack/react-query';
import { budgetApi, budgetKeys } from '@/entities/budget';

export const useCurrentBudgets = () => {
  return useQuery({
    queryKey: budgetKeys.current(),
    queryFn: () => budgetApi.getCurrent(),
    staleTime: 5 * 60 * 1000, // 5分
  });
};
