'use client';

import { useState } from 'react';
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
  Button,
  Pagination,
} from '@/shared/ui';
import { useTransactions, useDeleteTransaction } from '@/features/transaction-management';
import { useCategories } from '@/features/category-management';
import { TransactionAmountDisplay } from '@/entities/transaction';
import { CategoryIcon } from '@/entities/category';
import { formatMoney } from '@/shared/value-objects/money';
import type { Transaction, TransactionListParams } from '@/entities/transaction';

interface TransactionListTableProps {
  filters?: TransactionListParams;
  onEditTransaction?: (transaction: Transaction) => void;
}

export function TransactionListTable({
  filters = {},
  onEditTransaction,
}: TransactionListTableProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 20;

  const queryParams = {
    ...filters,
    limit: itemsPerPage,
    offset: (currentPage - 1) * itemsPerPage,
  };

  const {
    data: transactionData,
    isLoading,
    error,
  } = useTransactions(queryParams);

  const { data: categories } = useCategories();
  const deleteTransactionMutation = useDeleteTransaction();

  const handleDelete = async (transactionId: string) => {
    if (window.confirm('Are you sure you want to delete this transaction?')) {
      try {
        await deleteTransactionMutation.mutateAsync(transactionId);
      } catch (error) {
        console.error('Failed to delete transaction:', error);
      }
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
      month: '2-digit',
      day: '2-digit',
    });
  };

  const getCategoryInfo = (categoryId: string) => {
    const category = categories?.find((cat) => cat.id === categoryId);
    return {
      name: category?.name || 'Unknown',
      icon: category?.icon || '📦',
    };
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="animate-pulse">
          <div className="h-10 bg-gray-200 rounded mb-4"></div>
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="h-16 bg-gray-100 rounded mb-2"></div>
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6 text-center">
        <p className="text-red-600">Failed to load transactions</p>
        <Button
          onClick={() => window.location.reload()}
          className="mt-4"
          variant="outline"
        >
          Retry
        </Button>
      </div>
    );
  }

  const transactions = transactionData?.transactions || [];
  const totalPages = Math.ceil((transactionData?.total || 0) / itemsPerPage);

  return (
    <div className="space-y-4">
      <div className="bg-white rounded-lg border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Date</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {transactions.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center py-8 text-gray-500">
                  No transactions found
                </TableCell>
              </TableRow>
            ) : (
              transactions.map((transaction) => {
                const categoryInfo = getCategoryInfo(transaction.categoryId);
                return (
                  <TableRow key={transaction.id}>
                    <TableCell>{formatDate(transaction.date)}</TableCell>
                    <TableCell className="max-w-xs truncate">
                      {transaction.description}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <span>{categoryInfo.icon}</span>
                        <span className="text-sm">{categoryInfo.name}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        <span className={`text-sm ${
                          transaction.type === 'expense' ? 'text-red-600' : 'text-green-600'
                        }`}>
                          {transaction.type === 'expense' ? '↓' : '↑'}
                        </span>
                        <TransactionAmountDisplay
                          amount={transaction.amount}
                          type={transaction.type}
                        />
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Button
                          onClick={() => onEditTransaction?.(transaction)}
                          size="sm"
                          variant="outline"
                          className="h-8 w-8 p-0"
                        >
                          ✏️
                        </Button>
                        <Button
                          onClick={() => handleDelete(transaction.id)}
                          size="sm"
                          variant="outline"
                          className="h-8 w-8 p-0 hover:bg-red-50 hover:border-red-300"
                          disabled={deleteTransactionMutation.isPending}
                        >
                          🗑️
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                );
              })
            )}
          </TableBody>
        </Table>
      </div>

      {totalPages > 1 && (
        <Pagination
          currentPage={currentPage}
          totalPages={totalPages}
          onPageChange={setCurrentPage}
        />
      )}
    </div>
  );
}