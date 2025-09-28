import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  transactionApi,
  transactionKeys,
  type Transaction,
  type UpdateTransactionPayload,
} from '@/entities/transaction';

export const useUpdateTransaction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      ...payload
    }: { id: string } & UpdateTransactionPayload) =>
      transactionApi.update(id, payload),

    onMutate: async ({ id, ...newData }) => {
      await queryClient.cancelQueries({ queryKey: transactionKeys.detail(id) });

      const previousTransaction = queryClient.getQueryData<Transaction>(
        transactionKeys.detail(id),
      );

      queryClient.setQueryData<Transaction>(
        transactionKeys.detail(id),
        (old) => ({
          ...old!,
          ...newData,
          updatedAt: new Date().toISOString(),
        }),
      );

      return { previousTransaction };
    },

    onError: (_, { id }, context) => {
      if (context?.previousTransaction) {
        queryClient.setQueryData(
          transactionKeys.detail(id),
          context.previousTransaction,
        );
      }
    },

    onSuccess: (data, { id }) => {
      queryClient.setQueryData(transactionKeys.detail(id), data);
      queryClient.invalidateQueries({ queryKey: transactionKeys.lists() });
    },
  });
};
