import { useQuery } from '@tanstack/react-query';
import { accountApi, accountKeys } from '@/entities/account';

export const useAccount = (id: string) => {
  return useQuery({
    queryKey: accountKeys.detail(id),
    queryFn: () => accountApi.get(id),
    staleTime: 5 * 60 * 1000, // 5分
    enabled: !!id,
  });
};
