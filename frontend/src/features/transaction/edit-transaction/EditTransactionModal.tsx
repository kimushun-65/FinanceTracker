'use client';

import { useState, useEffect } from 'react';
import {
  Modal,
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
import { useUpdateTransaction } from '@/features/transaction-management';
import { useCategories } from '@/features/category-management';
import { useAccounts } from '@/features/account-management';
import type { Transaction, UpdateTransactionPayload, TransactionType } from '@/entities/transaction';

interface EditTransactionModalProps {
  transaction: Transaction | null;
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

export function EditTransactionModal({
  transaction,
  isOpen,
  onClose,
  onSuccess,
}: EditTransactionModalProps) {
  const updateTransactionMutation = useUpdateTransaction();
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

  // Update form data when transaction changes
  useEffect(() => {
    if (transaction) {
      setFormData({
        type: transaction.type,
        amount: transaction.amount.amount,
        categoryId: transaction.categoryId,
        accountId: transaction.accountId,
        description: transaction.description,
        date: transaction.date.split('T')[0],
      });
    }
  }, [transaction]);

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
    
    if (!transaction || !validateForm()) {
      return;
    }

    try {
      const payload: UpdateTransactionPayload = {
        accountId: formData.accountId,
        categoryId: formData.categoryId,
        amount: {
          amount: formData.amount,
          currency: 'JPY',
        },
        description: formData.description,
        date: formData.date,
      };

      await updateTransactionMutation.mutateAsync({
        id: transaction.id,
        data: payload,
      });
      
      onClose();
      onSuccess?.();
    } catch (error) {
      console.error('Failed to update transaction:', error);
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

  if (!isOpen || !transaction) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div
        className="fixed inset-0 bg-black bg-opacity-50"
        onClick={onClose}
        aria-hidden="true"
      />
      
      <div className="relative w-full max-w-lg mx-4 bg-white rounded-lg shadow-xl">
        <div className="flex items-center justify-between p-6 border-b">
          <h2 className="text-lg font-semibold">Edit Transaction</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 text-xl font-bold"
            aria-label="Close modal"
          >
            ×
          </button>
        </div>
        
        <div className="p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Transaction Type - Read only for edit */}
            <div>
              <Label htmlFor="type">Type</Label>
              <Select
                value={formData.type}
                onValueChange={(value) =>
                  setFormData({ ...formData, type: value as TransactionType })
                }
                options={typeOptions}
                disabled={true} // Transaction type shouldn't be changeable
              />
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

            <div className="flex justify-end gap-3 pt-4 border-t">
              <Button
                type="button"
                variant="outline"
                onClick={onClose}
                disabled={updateTransactionMutation.isPending}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={updateTransactionMutation.isPending}
              >
                {updateTransactionMutation.isPending ? 'Saving...' : 'Save Changes'}
              </Button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}