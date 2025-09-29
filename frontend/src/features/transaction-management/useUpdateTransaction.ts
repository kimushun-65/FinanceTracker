import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  transactionApi,
  transactionKeys,
  type Transaction,
  type UpdateTransactionPayload,
} from '@/entities/transaction';
import { useToast } from '@/shared/ui';

export const useUpdateTransaction = () => {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string;
      data: UpdateTransactionPayload;
    }) => transactionApi.update(id, data),

    onMutate: async ({ id, data: newData }) => {
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

      // サマリーを無効化
      queryClient.invalidateQueries({ queryKey: transactionKeys.summaries() });

      // 成功時のトースト表示
      toast({
        title: '取引が更新されました',
        description: '取引情報が正常に更新されました。',
      });
    },
  });
};
