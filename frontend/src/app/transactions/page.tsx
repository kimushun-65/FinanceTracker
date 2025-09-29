'use client';

import { useState } from 'react';
import { AppLayout } from '@/widgets/layout';
import { TransactionSummaryCards } from '@/widgets/transaction-summary';
import { TransactionListTable } from '@/widgets/transaction-list';
import { TransactionFilters } from '@/widgets/transaction-filters';
import {
  CreateTransactionModal,
  EditTransactionModal,
} from '@/features/transaction';
import { Button } from '@/shared/ui';
import type {
  Transaction,
  TransactionListParams,
} from '@/entities/transaction';

export default function TransactionsPage() {
  const [filters, setFilters] = useState<TransactionListParams>({});
  const [editingTransaction, setEditingTransaction] =
    useState<Transaction | null>(null);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  const handleEditTransaction = (transaction: Transaction) => {
    setEditingTransaction(transaction);
    setIsEditModalOpen(true);
  };

  const handleCloseEditModal = () => {
    setEditingTransaction(null);
    setIsEditModalOpen(false);
  };

  const handleTransactionSuccess = () => {
    // This will trigger a refetch of the transaction list
    // The useQuery hooks will automatically refetch when mutations succeed
  };

  return (
    <AppLayout>
      <div className='container mx-auto px-4 py-6'>
        {/* Transaction Summary Cards */}
        <TransactionSummaryCards />

        {/* Filters with Add Button */}
        <div className='mb-6 rounded-lg bg-gray-50 p-4'>
          <div className='flex flex-wrap items-center gap-4'>
            <TransactionFilters
              filters={filters}
              onFiltersChange={setFilters}
            />
            <div className='ml-auto'>
              <CreateTransactionModal
                onSuccess={handleTransactionSuccess}
                trigger={<Button size='lg'>+ Add Transaction</Button>}
              />
            </div>
          </div>
        </div>

        {/* Transaction List */}
        <div className='rounded-lg bg-white shadow'>
          <div className='border-b p-6'>
            <h2 className='text-lg font-semibold text-gray-900'>
              Transaction List
            </h2>
          </div>
          <div className='p-6'>
            <TransactionListTable
              filters={filters}
              onEditTransaction={handleEditTransaction}
            />
          </div>
        </div>

        {/* Edit Transaction Modal */}
        <EditTransactionModal
          transaction={editingTransaction}
          isOpen={isEditModalOpen}
          onClose={handleCloseEditModal}
          onSuccess={handleTransactionSuccess}
        />
      </div>
    </AppLayout>
  );
}
