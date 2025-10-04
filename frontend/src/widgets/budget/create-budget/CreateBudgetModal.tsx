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
import { useCreateBudget } from '@/features/budget-management';
import type { CreateBudgetPayload, PeriodType } from '@/entities/budget';
import { PERIOD_TYPES } from '@/entities/budget';
import {
  type CategoryWithMaster,
  categoriesToSelectOptions,
} from '@/entities/category';

interface CreateBudgetModalProps {
  trigger?: React.ReactNode;
  onSuccess?: () => void;
  categories: CategoryWithMaster[];
}

function CreateBudgetForm({
  onSuccess,
  categories,
}: {
  onSuccess?: () => void;
  categories: CategoryWithMaster[];
}) {
  const { closeModal } = useModalContext();
  const createBudgetMutation = useCreateBudget();

  const [formData, setFormData] = useState<{
    categoryId: string;
    amount: number;
    periodType: PeriodType;
    startDate: string;
    endDate?: string;
  }>({
    categoryId: '',
    amount: 0,
    periodType: 'monthly',
    startDate: new Date().toISOString().split('T')[0],
    endDate: undefined,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.categoryId) {
      newErrors.categoryId = 'Category is required';
    }
    if (!formData.amount || formData.amount <= 0) {
      newErrors.amount = 'Amount must be greater than 0';
    }
    if (!formData.startDate) {
      newErrors.startDate = 'Start date is required';
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
      const payload: CreateBudgetPayload = {
        categoryId: formData.categoryId,
        amount: {
          amount: formData.amount,
          currency: 'JPY',
        },
        periodType: formData.periodType,
        startDate: formData.startDate,
        endDate: formData.endDate,
      };

      await createBudgetMutation.mutateAsync(payload);
      closeModal();
      onSuccess?.();
    } catch (error) {
      console.error('Failed to create budget:', error);
    }
  };

  const periodTypeOptions = Object.entries(PERIOD_TYPES).map(
    ([value, label]) => ({
      value,
      label,
    }),
  );

  const categoryOptions = categoriesToSelectOptions(categories);

  return (
    <form onSubmit={handleSubmit} className='space-y-4'>
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

      {/* Amount */}
      <div>
        <Label htmlFor='amount'>Budget Amount (¥)</Label>
        <Input
          id='amount'
          type='number'
          min='0'
          step='1'
          value={formData.amount || ''}
          onChange={(e) =>
            setFormData({ ...formData, amount: Number(e.target.value) })
          }
          placeholder='Enter budget amount'
          className={errors.amount ? 'border-red-500' : ''}
        />
        {errors.amount && (
          <p className='mt-1 text-sm text-red-600'>{errors.amount}</p>
        )}
      </div>

      {/* Period Type */}
      <div>
        <Label htmlFor='periodType'>Period Type</Label>
        <Select
          value={formData.periodType}
          onValueChange={(value) =>
            setFormData({ ...formData, periodType: value as PeriodType })
          }
          options={periodTypeOptions}
        />
      </div>

      {/* Start Date */}
      <div>
        <Label htmlFor='startDate'>Start Date</Label>
        <Input
          id='startDate'
          type='date'
          value={formData.startDate}
          onChange={(e) =>
            setFormData({ ...formData, startDate: e.target.value })
          }
          className={errors.startDate ? 'border-red-500' : ''}
        />
        {errors.startDate && (
          <p className='mt-1 text-sm text-red-600'>{errors.startDate}</p>
        )}
      </div>

      {/* End Date (Optional) */}
      <div>
        <Label htmlFor='endDate'>End Date (Optional)</Label>
        <Input
          id='endDate'
          type='date'
          value={formData.endDate || ''}
          onChange={(e) =>
            setFormData({ ...formData, endDate: e.target.value || undefined })
          }
        />
        <p className='mt-1 text-xs text-gray-500'>
          Leave empty for ongoing budget
        </p>
      </div>

      <ModalFooter>
        <Button
          type='button'
          variant='outline'
          onClick={closeModal}
          disabled={createBudgetMutation.isPending}
        >
          Cancel
        </Button>
        <Button type='submit' disabled={createBudgetMutation.isPending}>
          {createBudgetMutation.isPending ? 'Creating...' : 'Create Budget'}
        </Button>
      </ModalFooter>
    </form>
  );
}

export function CreateBudgetModal({
  trigger,
  onSuccess,
  categories,
}: CreateBudgetModalProps) {
  return (
    <Modal>
      <ModalTrigger asChild>
        {trigger || <Button>+ Create Budget</Button>}
      </ModalTrigger>
      <ModalContent>
        <ModalHeader>Create Budget</ModalHeader>
        <ModalBody>
          <CreateBudgetForm onSuccess={onSuccess} categories={categories} />
        </ModalBody>
      </ModalContent>
    </Modal>
  );
}
