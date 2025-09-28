import { useState, useCallback } from 'react';
import type {
  TransactionListParams,
  TransactionType,
} from '@/entities/transaction';

export const useTransactionFilters = () => {
  const [filters, setFilters] = useState<TransactionListParams>({});

  const setAccountFilter = useCallback((accountId?: string) => {
    setFilters((prev) => ({ ...prev, accountId }));
  }, []);

  const setTypeFilter = useCallback((type?: TransactionType) => {
    setFilters((prev) => ({ ...prev, type }));
  }, []);

  const setDateRange = useCallback((startDate?: string, endDate?: string) => {
    setFilters((prev) => ({ ...prev, startDate, endDate }));
  }, []);

  const clearFilters = useCallback(() => {
    setFilters({});
  }, []);

  return {
    filters,
    setAccountFilter,
    setTypeFilter,
    setDateRange,
    clearFilters,
  };
};
