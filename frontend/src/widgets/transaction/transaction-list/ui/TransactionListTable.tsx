'use client';
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
  Button,
} from '@/shared/ui';
import { TransactionAmountDisplay } from '@/entities/transaction';
import { Category } from '@/entities/category';
import type {
  Transaction,
  TransactionListResponse,
} from '@/entities/transaction';
import { useDeleteTransaction } from '@/features';
import { formatDateShort } from '@/shared/value-objects/date';

interface TransactionListTableProps {
  onEditTransaction?: (transaction: Transaction) => void;
  categories: Category[];
  transactionData: TransactionListResponse;
}

export function TransactionListTable({
  onEditTransaction,
  categories,
  transactionData,
}: TransactionListTableProps) {
  const deleteTransactionMutation = useDeleteTransaction();

  const handleDelete = async (transactionId: string) => {
    try {
      await deleteTransactionMutation.mutateAsync(transactionId);
    } catch (error) {
      console.error('Failed to delete transaction:', error);
    }
  };

  const getCategoryInfo = (categoryId: string) => {
    const category: any = categories?.find((cat: any) => cat.id === categoryId);
    return {
      name: category?.name || 'Unknown',
      icon: category?.icon || '📦',
    };
  };

  const transactions = transactionData?.transactions || [];

  return (
    <div className='space-y-4'>
      <div className='max-h-96 overflow-auto rounded-lg border bg-white'>
        <Table>
          <TableHeader className='sticky top-0 z-10 bg-white'>
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
                <td colSpan={5} className='py-8 text-center text-gray-500'>
                  No transactions found
                </td>
              </TableRow>
            ) : (
              transactions.map((transaction) => {
                const categoryInfo = getCategoryInfo(transaction.categoryId);
                return (
                  <TableRow key={transaction.id}>
                    <TableCell>{formatDateShort(transaction.date)}</TableCell>
                    <TableCell className='max-w-xs truncate'>
                      {transaction.description}
                    </TableCell>
                    <TableCell>
                      <div className='flex items-center gap-2'>
                        <span>{categoryInfo.icon}</span>
                        <span className='text-sm'>{categoryInfo.name}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className='flex items-center gap-1'>
                        <span
                          className={`text-sm ${
                            transaction.type === 'expense'
                              ? 'text-red-600'
                              : 'text-green-600'
                          }`}
                        >
                          {transaction.type === 'expense' ? '↓' : '↑'}
                        </span>
                        <TransactionAmountDisplay
                          amount={transaction.amount}
                          type={transaction.type}
                        />
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className='flex items-center gap-2'>
                        <Button
                          onClick={() => onEditTransaction?.(transaction)}
                          size='sm'
                          variant='outline'
                          className='h-8 w-8 p-0'
                        >
                          ✏️
                        </Button>
                        <Button
                          onClick={() => handleDelete(transaction.id)}
                          size='sm'
                          variant='outline'
                          className='h-8 w-8 p-0 hover:border-red-300 hover:bg-red-50'
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
    </div>
  );
}
