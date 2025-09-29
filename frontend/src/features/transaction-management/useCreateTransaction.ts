import { useMutation, useQueryClient } from '@tanstack/react-query';
import { transactionApi, transactionKeys } from '@/entities/transaction';
import { accountKeys } from '@/entities/account';
import { budgetKeys } from '@/entities/budget';
import { useToast } from '@/shared/ui';

export const useCreateTransaction = () => {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: transactionApi.create,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: transactionKeys.lists() });

      // サマリーを無効化
      queryClient.invalidateQueries({ queryKey: transactionKeys.summaries() });

      // 関連アカウントの残高を無効化
      if (data.accountId) {
        queryClient.invalidateQueries({
          queryKey: accountKeys.detail(data.accountId),
        });
      }

      // 予算の使用状況を無効化
      queryClient.invalidateQueries({
        queryKey: budgetKeys.lists(),
      });

      // 成功時のトースト表示
      toast({
        title: '取引が作成されました',
        description: '新しい取引が正常に作成されました。',
      });
    },
  });
};
