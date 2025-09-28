import { useQuery } from '@tanstack/react-query';
import {
  transactionApi,
  transactionKeys,
  type TransactionListParams,
} from '@/entities/transaction';

export const useTransactions = (params: TransactionListParams = {}) => {
  return useQuery({
    queryKey: transactionKeys.list(params),
    queryFn: () => transactionApi.list(params),
    staleTime: 2 * 60 * 1000, // 2分
    enabled: true,
  });
};
