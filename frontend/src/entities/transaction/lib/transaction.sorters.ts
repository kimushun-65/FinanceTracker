import type { Transaction } from '../model';

export type SortDirection = 'asc' | 'desc';
export type TransactionSortField =
  | 'date'
  | 'amount'
  | 'description'
  | 'type'
  | 'createdAt';

export const sortTransactionsByDate = (
  transactions: Transaction[],
  direction: SortDirection = 'desc',
): Transaction[] => {
  return [...transactions].sort((a, b) => {
    const dateA = new Date(a.date).getTime();
    const dateB = new Date(b.date).getTime();
    return direction === 'asc' ? dateA - dateB : dateB - dateA;
  });
};

export const sortTransactionsByAmount = (
  transactions: Transaction[],
  direction: SortDirection = 'desc',
): Transaction[] => {
  return [...transactions].sort((a, b) => {
    const amountA = a.amount.amount;
    const amountB = b.amount.amount;
    return direction === 'asc' ? amountA - amountB : amountB - amountA;
  });
};

export const sortTransactionsByDescription = (
  transactions: Transaction[],
  direction: SortDirection = 'asc',
): Transaction[] => {
  return [...transactions].sort((a, b) => {
    const comparison = a.description.localeCompare(b.description, 'ja-JP');
    return direction === 'asc' ? comparison : -comparison;
  });
};

export const sortTransactionsByType = (
  transactions: Transaction[],
  direction: SortDirection = 'asc',
): Transaction[] => {
  return [...transactions].sort((a, b) => {
    const comparison = a.type.localeCompare(b.type);
    return direction === 'asc' ? comparison : -comparison;
  });
};

export const sortTransactionsByCreatedAt = (
  transactions: Transaction[],
  direction: SortDirection = 'desc',
): Transaction[] => {
  return [...transactions].sort((a, b) => {
    const dateA = new Date(a.createdAt).getTime();
    const dateB = new Date(b.createdAt).getTime();
    return direction === 'asc' ? dateA - dateB : dateB - dateA;
  });
};

export const createTransactionSorter = (
  field: TransactionSortField,
  direction: SortDirection = 'desc',
) => {
  return (transactions: Transaction[]): Transaction[] => {
    switch (field) {
      case 'date':
        return sortTransactionsByDate(transactions, direction);
      case 'amount':
        return sortTransactionsByAmount(transactions, direction);
      case 'description':
        return sortTransactionsByDescription(transactions, direction);
      case 'type':
        return sortTransactionsByType(transactions, direction);
      case 'createdAt':
        return sortTransactionsByCreatedAt(transactions, direction);
      default:
        return transactions;
    }
  };
};

export const multiSortTransactions = (
  transactions: Transaction[],
  sortConfig: Array<{ field: TransactionSortField; direction: SortDirection }>,
): Transaction[] => {
  return [...transactions].sort((a, b) => {
    for (const { field, direction } of sortConfig) {
      let comparison = 0;

      switch (field) {
        case 'date':
          comparison = new Date(a.date).getTime() - new Date(b.date).getTime();
          break;
        case 'amount':
          comparison = a.amount.amount - b.amount.amount;
          break;
        case 'description':
          comparison = a.description.localeCompare(b.description, 'ja-JP');
          break;
        case 'type':
          comparison = a.type.localeCompare(b.type);
          break;
        case 'createdAt':
          comparison =
            new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
          break;
      }

      if (comparison !== 0) {
        return direction === 'asc' ? comparison : -comparison;
      }
    }

    return 0;
  });
};
