'use client';

import { Button } from '@/shared/ui';
import { useDeleteBudget } from '@/features/budget-management';

interface DeleteBudgetModalProps {
  isOpen: boolean;
  onClose: () => void;
  budgetId: string;
  categoryName: string;
}

export function DeleteBudgetModal({
  isOpen,
  onClose,
  budgetId,
  categoryName,
}: DeleteBudgetModalProps) {
  const deleteBudgetMutation = useDeleteBudget();

  const handleDelete = async () => {
    try {
      await deleteBudgetMutation.mutateAsync(budgetId);
      onClose();
    } catch (error) {
      console.error('Failed to delete budget:', error);
    }
  };

  if (!isOpen) return null;

  return (
    <div className='fixed inset-0 z-50 flex items-center justify-center bg-black/50'>
      <div className='relative mx-4 w-full max-w-md rounded-lg bg-white shadow-xl'>
        {/* Header */}
        <div className='flex items-center justify-between border-b p-6'>
          <h2 className='text-lg font-semibold'>Delete Budget</h2>
          <button
            onClick={onClose}
            className='text-xl font-bold text-gray-400 hover:text-gray-600'
            aria-label='Close modal'
          >
            ×
          </button>
        </div>

        {/* Body */}
        <div className='p-6'>
          <p className='text-gray-700'>
            Are you sure you want to delete the budget for{' '}
            <span className='font-semibold'>{categoryName}</span>?
          </p>
          <p className='mt-2 text-sm text-gray-500'>
            This action cannot be undone.
          </p>
        </div>

        {/* Footer */}
        <div className='flex justify-end gap-2 border-t p-6'>
          <Button
            onClick={onClose}
            variant='outline'
            disabled={deleteBudgetMutation.isPending}
          >
            Cancel
          </Button>
          <Button
            onClick={handleDelete}
            className='bg-red-600 text-white hover:bg-red-700'
            disabled={deleteBudgetMutation.isPending}
          >
            {deleteBudgetMutation.isPending ? 'Deleting...' : 'Delete'}
          </Button>
        </div>
      </div>
    </div>
  );
}
