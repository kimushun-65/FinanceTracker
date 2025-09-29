import { useQuery } from '@tanstack/react-query';
import {
  transactionApi,
  transactionKeys,
  type MonthlySummaryParams,
} from '@/entities/transaction';

export const useTransactionMonthlySummary = (
  params: MonthlySummaryParams = {},
) => {
  // Default to current year/month when not provided
  const now = new Date();
  const defaultYear = now.getFullYear();
  const defaultMonth = now.getMonth() + 1; // 1-12

  const effectiveParams: MonthlySummaryParams = {
    year: params.year ?? defaultYear,
    month: params.month ?? defaultMonth,
  };

  return useQuery({
    queryKey: transactionKeys.monthlySummary(effectiveParams),
    queryFn: () => transactionApi.getMonthlySummary(effectiveParams),
    staleTime: 5 * 60 * 1000, // 5分
    enabled: true,
  });
};
