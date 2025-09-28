import { useMutation, useQueryClient } from '@tanstack/react-query';
import { accountApi, accountKeys, type CreateAccountPayload } from '@/entities/account';

export const useCreateAccount = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: accountApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: accountKeys.lists() });
    },
  });
};