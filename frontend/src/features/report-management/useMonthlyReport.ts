import { useMemo } from 'react';
import { calculateMonthlyReport } from '@/entities/report';
import type { Transaction } from '@/entities/transaction';
import type { CategorySummary } from '@/entities/report';

export const useMonthlyReport = (
  transactions: Transaction[] | undefined,
  categorySummary: CategorySummary[] | undefined,
) => {
  return useMemo(() => {
    if (!transactions || !categorySummary) return null;
    // transactionとcategorySummaryは構造互換性があるため、型アサーションで変換
    return calculateMonthlyReport(transactions as any, categorySummary);
  }, [transactions, categorySummary]);
};
