import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  budgetApi,
  budgetKeys,
  type CreateBudgetPayload,
} from '@/entities/budget';

export const useCreateBudget = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: budgetApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.lists() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.summary() });
    },
  });
};
