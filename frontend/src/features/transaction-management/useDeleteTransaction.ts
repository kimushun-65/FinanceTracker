import { useMutation, useQueryClient } from '@tanstack/react-query';
import { transactionApi, transactionKeys } from '@/entities/transaction';
import { useToast } from '@/shared/ui';

export const useDeleteTransaction = () => {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: transactionApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: transactionKeys.lists() });

      // サマリーを無効化
      queryClient.invalidateQueries({ queryKey: transactionKeys.summaries() });

      // 成功時のトースト表示
      toast({
        title: '取引が削除されました',
        description: '取引が正常に削除されました。',
      });
    },
  });
};
