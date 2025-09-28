import { useMutation, useQueryClient } from '@tanstack/react-query';
import { transactionApi, transactionKeys } from '@/entities/transaction';

export const useDeleteTransaction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: transactionApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: transactionKeys.lists() });
    },
  });
};
