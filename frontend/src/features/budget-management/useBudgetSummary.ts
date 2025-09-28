import { useQuery } from '@tanstack/react-query';
import { budgetApi, budgetKeys } from '@/entities/budget';

export const useBudgetSummary = () => {
  return useQuery({
    queryKey: budgetKeys.summary(),
    queryFn: () => budgetApi.getSummary(),
    staleTime: 5 * 60 * 1000, // 5分
  });
};
