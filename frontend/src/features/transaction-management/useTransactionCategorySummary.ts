import { useQuery } from '@tanstack/react-query';
import {
  transactionApi,
  transactionKeys,
  type CategorySummaryParams,
} from '@/entities/transaction';

export const useTransactionCategorySummary = (
  params: CategorySummaryParams = {},
) => {
  return useQuery({
    queryKey: transactionKeys.categorySummary(params),
    queryFn: () => transactionApi.getCategorySummary(params),
    staleTime: 5 * 60 * 1000, // 5分
    enabled: true,
  });
};
