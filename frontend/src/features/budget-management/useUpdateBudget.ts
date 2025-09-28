import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  budgetApi,
  budgetKeys,
  type Budget,
  type UpdateBudgetPayload,
} from '@/entities/budget';

export const useUpdateBudget = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, ...payload }: { id: string } & UpdateBudgetPayload) =>
      budgetApi.update(id, payload),

    onMutate: async ({ id, ...newData }) => {
      await queryClient.cancelQueries({ queryKey: budgetKeys.detail(id) });

      const previousBudget = queryClient.getQueryData<Budget>(
        budgetKeys.detail(id),
      );

      queryClient.setQueryData<Budget>(budgetKeys.detail(id), (old) => ({
        ...old!,
        ...newData,
        updatedAt: new Date().toISOString(),
      }));

      return { previousBudget };
    },

    onError: (_, { id }, context) => {
      if (context?.previousBudget) {
        queryClient.setQueryData(budgetKeys.detail(id), context.previousBudget);
      }
    },

    onSuccess: (data, { id }) => {
      queryClient.setQueryData(budgetKeys.detail(id), data);
      queryClient.invalidateQueries({ queryKey: budgetKeys.lists() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.summary() });
    },
  });
};
