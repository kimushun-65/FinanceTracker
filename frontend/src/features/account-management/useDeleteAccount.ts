import { useMutation, useQueryClient } from '@tanstack/react-query';
import { accountApi, accountKeys } from '@/entities/account';

export const useDeleteAccount = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: accountApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: accountKeys.lists() });
    },
  });
};