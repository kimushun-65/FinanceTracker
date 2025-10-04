import { useMutation, useQueryClient } from '@tanstack/react-query';
import { budgetApi, budgetKeys } from '@/entities/budget';

export const useDeleteBudget = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: budgetApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.lists() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.current() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.all });
    },
  });
};
