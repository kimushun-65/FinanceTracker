'use client';

import { useState } from 'react';
import {
  AppLayout,
  TransactionSummaryCards,
  TransactionListTable,
  TransactionFilters,
  CreateTransactionModal,
  EditTransactionModal,
} from '@/widgets';
import { Button, Loading } from '@/shared/ui';
import type {
  Transaction,
  TransactionListParams,
} from '@/entities/transaction';
import {
  useAccounts,
  useCategories,
  useTransactionMonthlySummary,
  useTransactions,
} from '@/features';

export default function TransactionsContainer() {
  const [filters, setFilters] = useState<TransactionListParams>({});
  const { data: transactionData, isLoading: transactionIsLoading } =
    useTransactions(filters);
  const { data: summary, isLoading } = useTransactionMonthlySummary();
  const { data: categories } = useCategories();
  const { data: accounts } = useAccounts();
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

  if (transactionIsLoading || isLoading) {
    return <Loading />;
  }

  if (!summary || !categories || !transactionData || !accounts) {
    return <Loading />;
  }

  return (
    <AppLayout>
      <div className='container mx-auto px-4 py-6'>
        {/* Transaction Summary Cards */}
        <TransactionSummaryCards summary={summary} />

        {/* Filters with Add Button */}
        <div className='mb-6 rounded-lg bg-gray-50 p-4'>
          <div className='flex flex-wrap items-center gap-4'>
            <TransactionFilters
              filters={filters}
              onFiltersChange={setFilters}
              categories={categories}
            />
            <div className='ml-auto'>
              <CreateTransactionModal
                trigger={<Button size='lg'>+ Add Transaction</Button>}
                accounts={accounts}
                categories={categories}
              />
            </div>
          </div>
        </div>

        {/* Transaction List */}
        <div className='rounded-lg bg-white shadow'>
          <div className='p-6'>
            <TransactionListTable
              onEditTransaction={handleEditTransaction}
              categories={categories}
              transactionData={transactionData}
            />
          </div>
        </div>

        {/* Edit Transaction Modal */}
        <EditTransactionModal
          transaction={editingTransaction}
          isOpen={isEditModalOpen}
          onClose={handleCloseEditModal}
          categories={categories}
          accounts={accounts}
        />
      </div>
    </AppLayout>
  );
}
