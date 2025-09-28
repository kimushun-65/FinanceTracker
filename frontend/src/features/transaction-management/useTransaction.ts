import { useQuery } from '@tanstack/react-query';
import { transactionApi, transactionKeys } from '@/entities/transaction';

export const useTransaction = (id: string) => {
  return useQuery({
    queryKey: transactionKeys.detail(id),
    queryFn: () => transactionApi.get(id),
    staleTime: 2 * 60 * 1000, // 2分
    enabled: !!id,
  });
};
