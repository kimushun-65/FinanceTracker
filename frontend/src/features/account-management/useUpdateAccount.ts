import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  accountApi,
  accountKeys,
  type Account,
  type UpdateAccountPayload,
} from '@/entities/account';

export const useUpdateAccount = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, ...payload }: { id: string } & UpdateAccountPayload) =>
      accountApi.update(id, payload),

    onMutate: async ({ id }) => {
      await queryClient.cancelQueries({ queryKey: accountKeys.detail(id) });

      const previousAccount = queryClient.getQueryData<Account>(
        accountKeys.detail(id),
      );

      return { previousAccount };
    },

    onError: (_, { id }, context) => {
      if (context?.previousAccount) {
        queryClient.setQueryData(
          accountKeys.detail(id),
          context.previousAccount,
        );
      }
    },

    onSuccess: (data, { id }) => {
      queryClient.setQueryData(accountKeys.detail(id), data);
      queryClient.invalidateQueries({ queryKey: accountKeys.lists() });
    },
  });
};
