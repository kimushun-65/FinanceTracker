'use client';

import { useState } from 'react';
import {
  Modal,
  ModalTrigger,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  Input,
  Label,
  Select,
  useModalContext,
} from '@/shared/ui';
import { useCreateTransaction } from '@/features/transaction-management';
import type {
  CreateTransactionPayload,
  TransactionType,
} from '@/entities/transaction';
import {
  transactionTypeOptions,
  validateTransactionForm,
} from '@/entities/transaction';
import { Account, accountsToSelectOptions } from '@/entities/account';
import { Category, categoriesToSelectOptions } from '@/entities/category';

interface CreateTransactionModalProps {
  trigger?: React.ReactNode;
  onSuccess?: () => void;
  accounts: Account[];
  categories: Category[];
}

function CreateTransactionForm({
  onSuccess,
  accounts,
  categories,
}: {
  onSuccess?: () => void;
  accounts: Account[];
  categories: Category[];
}) {
  const { closeModal } = useModalContext();
  const createTransactionMutation = useCreateTransaction();

  const [formData, setFormData] = useState<{
    type: TransactionType;
    amount: number;
    categoryId: string;
    accountId: string;
    description: string;
    date: string;
  }>({
    type: 'expense',
    amount: 0,
    categoryId: '',
    accountId: '',
    description: '',
    date: new Date().toISOString().split('T')[0],
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const newErrors = validateTransactionForm(formData);
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    try {
      const payload: CreateTransactionPayload = {
        accountId: formData.accountId,
        categoryId: formData.categoryId,
        amount: {
          amount: formData.amount,
          currency: 'JPY',
        },
        type: formData.type,
        description: formData.description,
        date: formData.date,
      };

      await createTransactionMutation.mutateAsync(payload);
      closeModal();
      onSuccess?.();
    } catch (error) {
      console.error('Failed to create transaction:', error);
    }
  };

  const typeOptions = transactionTypeOptions;
  // Only show active categories
  const activeCategories = categories.filter((cat) => cat.isActive);
  const categoryOptions = categoriesToSelectOptions(activeCategories);
  const accountOptions = accountsToSelectOptions(accounts);

  return (
    <form onSubmit={handleSubmit} className='space-y-4'>
      {/* Transaction Type */}
      <div>
        <Label htmlFor='type'>Type</Label>
        <Select
          value={formData.type}
          onValueChange={(value) =>
            setFormData({ ...formData, type: value as TransactionType })
          }
          options={typeOptions}
        />
        {errors.type && (
          <p className='mt-1 text-sm text-red-600'>{errors.type}</p>
        )}
      </div>

      {/* Amount */}
      <div>
        <Label htmlFor='amount'>Amount</Label>
        <Input
          id='amount'
          type='number'
          min='0'
          step='1'
          value={formData.amount || ''}
          onChange={(e) =>
            setFormData({ ...formData, amount: Number(e.target.value) })
          }
          placeholder='Enter amount'
          className={errors.amount ? 'border-red-500' : ''}
        />
        {errors.amount && (
          <p className='mt-1 text-sm text-red-600'>{errors.amount}</p>
        )}
      </div>

      {/* Category */}
      <div>
        <Label htmlFor='category'>Category</Label>
        <Select
          value={formData.categoryId}
          onValueChange={(value) =>
            setFormData({ ...formData, categoryId: value })
          }
          options={categoryOptions}
          placeholder='Select Category'
        />
        {errors.categoryId && (
          <p className='mt-1 text-sm text-red-600'>{errors.categoryId}</p>
        )}
      </div>

      {/* Account */}
      <div>
        <Label htmlFor='account'>Account</Label>
        <Select
          value={formData.accountId}
          onValueChange={(value) =>
            setFormData({ ...formData, accountId: value })
          }
          options={accountOptions}
          placeholder='Select Account'
        />
        {errors.accountId && (
          <p className='mt-1 text-sm text-red-600'>{errors.accountId}</p>
        )}
      </div>

      {/* Date */}
      <div>
        <Label htmlFor='date'>Date</Label>
        <Input
          id='date'
          type='date'
          value={formData.date}
          onChange={(e) => setFormData({ ...formData, date: e.target.value })}
          className={errors.date ? 'border-red-500' : ''}
        />
        {errors.date && (
          <p className='mt-1 text-sm text-red-600'>{errors.date}</p>
        )}
      </div>

      {/* Description */}
      <div>
        <Label htmlFor='description'>Description</Label>
        <Input
          id='description'
          type='text'
          value={formData.description}
          onChange={(e) =>
            setFormData({ ...formData, description: e.target.value })
          }
          placeholder='Enter description'
          className={errors.description ? 'border-red-500' : ''}
        />
        {errors.description && (
          <p className='mt-1 text-sm text-red-600'>{errors.description}</p>
        )}
      </div>

      <ModalFooter>
        <Button
          type='button'
          variant='outline'
          onClick={closeModal}
          disabled={createTransactionMutation.isPending}
        >
          Cancel
        </Button>
        <Button type='submit' disabled={createTransactionMutation.isPending}>
          {createTransactionMutation.isPending
            ? 'Adding...'
            : 'Add Transaction'}
        </Button>
      </ModalFooter>
    </form>
  );
}

export function CreateTransactionModal({
  trigger,
  onSuccess,
  accounts,
  categories,
}: CreateTransactionModalProps) {
  return (
    <Modal>
      <ModalTrigger asChild>
        {trigger || <Button>+ Add Transaction</Button>}
      </ModalTrigger>
      <ModalContent>
        <ModalHeader>Add Transaction</ModalHeader>
        <ModalBody>
          <CreateTransactionForm
            onSuccess={onSuccess}
            accounts={accounts}
            categories={categories}
          />
        </ModalBody>
      </ModalContent>
    </Modal>
  );
}
