import { useMutation, useQueryClient } from '@tanstack/react-query';
import { transactionApi, transactionKeys, accountKeys, budgetKeys } from '@/entities/transaction';

export const useCreateTransaction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: transactionApi.create,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: transactionKeys.lists() });
      
      // 関連アカウントの残高を無効化
      if (data.accountId) {
        queryClient.invalidateQueries({ 
          queryKey: accountKeys.detail(data.accountId) 
        });
      }
      
      // 予算の使用状況を無効化
      queryClient.invalidateQueries({ 
        queryKey: budgetKeys.lists() 
      });
    },
  });
};