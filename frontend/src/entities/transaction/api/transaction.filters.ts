import type { TransactionListParams, TransactionType } from '../model';

export const createTransactionFilters = (params: TransactionListParams) => {
  const searchParams = new URLSearchParams();

  if (params.accountId) searchParams.append('accountId', params.accountId);
  if (params.categoryId) searchParams.append('categoryId', params.categoryId);
  if (params.type) searchParams.append('type', params.type);
  if (params.startDate) searchParams.append('startDate', params.startDate);
  if (params.endDate) searchParams.append('endDate', params.endDate);
  if (params.limit) searchParams.append('limit', params.limit.toString());
  if (params.offset) searchParams.append('offset', params.offset.toString());

  return searchParams.toString();
};

export const validateDateRange = (
  startDate?: string,
  endDate?: string,
): boolean => {
  if (!startDate || !endDate) return true;
  return new Date(startDate) <= new Date(endDate);
};

export const filterTransactionsByType =
  (type?: TransactionType) =>
  (params: TransactionListParams): TransactionListParams => {
    return { ...params, type };
  };

export const filterTransactionsByAccount =
  (accountId: string) =>
  (params: TransactionListParams): TransactionListParams => {
    return { ...params, accountId };
  };

export const filterTransactionsByCategory =
  (categoryId: string) =>
  (params: TransactionListParams): TransactionListParams => {
    return { ...params, categoryId };
  };

export const filterTransactionsByDateRange =
  (startDate: string, endDate: string) =>
  (params: TransactionListParams): TransactionListParams => {
    return { ...params, startDate, endDate };
  };
