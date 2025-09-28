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
import { useCategories } from '@/features/category-management';
import { useAccounts } from '@/features/account-management';
import type { CreateTransactionPayload, TransactionType } from '@/entities/transaction';

interface CreateTransactionModalProps {
  trigger?: React.ReactNode;
  onSuccess?: () => void;
}

function CreateTransactionForm({ onSuccess }: { onSuccess?: () => void }) {
  const { closeModal } = useModalContext();
  const createTransactionMutation = useCreateTransaction();
  const { data: categories } = useCategories();
  const { data: accounts } = useAccounts();

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
    const newErrors: Record<string, string> = {};

    if (!formData.amount || formData.amount <= 0) {
      newErrors.amount = 'Amount must be greater than 0';
    }
    if (!formData.categoryId) {
      newErrors.categoryId = 'Category is required';
    }
    if (!formData.accountId) {
      newErrors.accountId = 'Account is required';
    }
    if (!formData.description.trim()) {
      newErrors.description = 'Description is required';
    }

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

  const typeOptions = [
    { value: 'expense', label: 'Expense' },
    { value: 'income', label: 'Income' },
  ];

  const categoryOptions = categories?.map((category) => ({
    value: category.id,
    label: category.name,
  })) || [];

  const accountOptions = accounts?.map((account) => ({
    value: account.id,
    label: account.name,
  })) || [];

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Transaction Type */}
      <div>
        <Label htmlFor="type">Type</Label>
        <Select
          value={formData.type}
          onValueChange={(value) =>
            setFormData({ ...formData, type: value as TransactionType })
          }
          options={typeOptions}
        />
        {errors.type && (
          <p className="text-sm text-red-600 mt-1">{errors.type}</p>
        )}
      </div>

      {/* Amount */}
      <div>
        <Label htmlFor="amount">Amount</Label>
        <Input
          id="amount"
          type="number"
          min="0"
          step="1"
          value={formData.amount || ''}
          onChange={(e) =>
            setFormData({ ...formData, amount: Number(e.target.value) })
          }
          placeholder="Enter amount"
          className={errors.amount ? 'border-red-500' : ''}
        />
        {errors.amount && (
          <p className="text-sm text-red-600 mt-1">{errors.amount}</p>
        )}
      </div>

      {/* Category */}
      <div>
        <Label htmlFor="category">Category</Label>
        <Select
          value={formData.categoryId}
          onValueChange={(value) =>
            setFormData({ ...formData, categoryId: value })
          }
          options={categoryOptions}
          placeholder="Select Category"
        />
        {errors.categoryId && (
          <p className="text-sm text-red-600 mt-1">{errors.categoryId}</p>
        )}
      </div>

      {/* Account */}
      <div>
        <Label htmlFor="account">Account</Label>
        <Select
          value={formData.accountId}
          onValueChange={(value) =>
            setFormData({ ...formData, accountId: value })
          }
          options={accountOptions}
          placeholder="Select Account"
        />
        {errors.accountId && (
          <p className="text-sm text-red-600 mt-1">{errors.accountId}</p>
        )}
      </div>

      {/* Date */}
      <div>
        <Label htmlFor="date">Date</Label>
        <Input
          id="date"
          type="date"
          value={formData.date}
          onChange={(e) => setFormData({ ...formData, date: e.target.value })}
          className={errors.date ? 'border-red-500' : ''}
        />
        {errors.date && (
          <p className="text-sm text-red-600 mt-1">{errors.date}</p>
        )}
      </div>

      {/* Description */}
      <div>
        <Label htmlFor="description">Description</Label>
        <Input
          id="description"
          type="text"
          value={formData.description}
          onChange={(e) =>
            setFormData({ ...formData, description: e.target.value })
          }
          placeholder="Enter description"
          className={errors.description ? 'border-red-500' : ''}
        />
        {errors.description && (
          <p className="text-sm text-red-600 mt-1">{errors.description}</p>
        )}
      </div>

      <ModalFooter>
        <Button
          type="button"
          variant="outline"
          onClick={closeModal}
          disabled={createTransactionMutation.isPending}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          disabled={createTransactionMutation.isPending}
        >
          {createTransactionMutation.isPending ? 'Adding...' : 'Add Transaction'}
        </Button>
      </ModalFooter>
    </form>
  );
}

export function CreateTransactionModal({
  trigger,
  onSuccess,
}: CreateTransactionModalProps) {
  return (
    <Modal>
      <ModalTrigger asChild>
        {trigger || (
          <Button>
            + Add Transaction
          </Button>
        )}
      </ModalTrigger>
      <ModalContent>
        <ModalHeader>Add Transaction</ModalHeader>
        <ModalBody>
          <CreateTransactionForm onSuccess={onSuccess} />
        </ModalBody>
      </ModalContent>
    </Modal>
  );
}