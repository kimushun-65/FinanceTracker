import { useQuery } from '@tanstack/react-query';
import { accountApi, accountKeys } from '@/entities/account';

export const useAccounts = () => {
  return useQuery({
    queryKey: accountKeys.lists(),
    queryFn: accountApi.list,
    staleTime: 5 * 60 * 1000, // 5分
  });
};