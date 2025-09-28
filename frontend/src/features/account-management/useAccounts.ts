import { useQuery } from '@tanstack/react-query';
import { accountApi, accountKeys } from '@/entities/account';

export const useAccounts = () => {
  return useQuery({
    queryKey: accountKeys.lists(),
    queryFn: async () => {
      const response = await accountApi.list();
      return response.accounts; // accounts 配列を直接返す
    },
    staleTime: 5 * 60 * 1000, // 5分
  });
};
