import { useQuery } from '@tanstack/react-query';
import {
  transactionApi,
  transactionKeys,
  type TransactionListParams,
} from '@/entities/transaction';
import { calculateMonthlyAggregates } from '@/entities/transaction';

export const useTransactionAggregates = (
  params: TransactionListParams = {},
) => {
  return useQuery({
    queryKey: [...transactionKeys.list(params), 'aggregates'],
    queryFn: async () => {
      const response = await transactionApi.list(params);
      return calculateMonthlyAggregates(response.transactions);
    },
    staleTime: 5 * 60 * 1000, // 5分
  });
};
