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
  Input,
} from '@/shared/ui';
import { BudgetProgressBar } from '@/entities/budget';
import type { CategoryWithMaster } from '@/entities/category';
import type { BudgetWithUsage } from '@/entities/budget';
import { useUpdateBudget } from '@/features';
import { formatMoney, type Currency } from '@/shared/value-objects/money';
import { DeleteBudgetModal } from '../../delete-budget';

interface BudgetListTableProps {
  budgets: BudgetWithUsage[];
  categories: CategoryWithMaster[];
  onEditBudget?: (budget: BudgetWithUsage) => void;
}

export function BudgetListTable({
  budgets,
  categories,
  onEditBudget,
}: BudgetListTableProps) {
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editAmount, setEditAmount] = useState<number>(0);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [deletingBudget, setDeletingBudget] = useState<{
    id: string;
    categoryName: string;
  } | null>(null);
  const updateBudgetMutation = useUpdateBudget();

  const handleEdit = (budget: BudgetWithUsage) => {
    setEditingId(budget.id);
    setEditAmount(budget.amount.amount);
  };

  const handleSave = async (budgetId: string, currency: Currency) => {
    try {
      await updateBudgetMutation.mutateAsync({
        id: budgetId,
        amount: { amount: editAmount, currency },
      });
      setEditingId(null);
    } catch (error) {
      console.error('Failed to update budget:', error);
    }
  };

  const handleCancel = () => {
    setEditingId(null);
    setEditAmount(0);
  };

  const handleDelete = (budgetId: string, categoryName: string) => {
    setDeletingBudget({ id: budgetId, categoryName });
    setDeleteModalOpen(true);
  };

  const handleCloseDeleteModal = () => {
    setDeleteModalOpen(false);
    setDeletingBudget(null);
  };

  const getCategoryInfo = (categoryId: string) => {
    const category: any = categories?.find((cat: any) => cat.id === categoryId);
    return {
      name: category?.name || 'Unknown',
      icon: category?.icon || '📦',
    };
  };

  return (
    <div className='space-y-4'>
      <div className='max-h-96 overflow-auto rounded-lg border bg-white'>
        <Table>
          <TableHeader className='sticky top-0 z-10 bg-white'>
            <TableRow>
              <TableHead>Category</TableHead>
              <TableHead>Budget</TableHead>
              <TableHead>Spent</TableHead>
              <TableHead>Remaining</TableHead>
              <TableHead className='w-64'>Progress</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {budgets.length === 0 ? (
              <TableRow>
                <td colSpan={6} className='py-8 text-center text-gray-500'>
                  No budgets found. Create your first budget!
                </td>
              </TableRow>
            ) : (
              budgets.map((budget) => {
                const categoryInfo = getCategoryInfo(budget.categoryId);
                const isEditing = editingId === budget.id;

                return (
                  <TableRow key={budget.id}>
                    <TableCell>
                      <div className='flex items-center gap-2'>
                        <span>{categoryInfo.icon}</span>
                        <span className='text-sm font-medium'>
                          {categoryInfo.name}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      {isEditing ? (
                        <Input
                          type='number'
                          value={editAmount}
                          onChange={(e) =>
                            setEditAmount(Number(e.target.value))
                          }
                          className='w-32'
                          min={0}
                        />
                      ) : (
                        <span className='font-semibold'>
                          {formatMoney(budget.amount)}
                        </span>
                      )}
                    </TableCell>
                    <TableCell>
                      <span className='text-orange-600'>
                        {formatMoney(budget.used)}
                      </span>
                    </TableCell>
                    <TableCell>
                      <span
                        className={
                          budget.remaining.amount >= 0
                            ? 'text-green-600'
                            : 'text-red-600'
                        }
                      >
                        {formatMoney(budget.remaining)}
                      </span>
                    </TableCell>
                    <TableCell>
                      <BudgetProgressBar
                        budget={budget}
                        size='sm'
                        showPercentage={false}
                      />
                      <div className='mt-1 text-xs text-gray-500'>
                        {budget.usagePercentage.toFixed(1)}%
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className='flex items-center gap-2'>
                        {isEditing ? (
                          <>
                            <Button
                              onClick={() =>
                                handleSave(budget.id, budget.amount.currency)
                              }
                              size='sm'
                              variant='outline'
                              className='h-8 px-2'
                              disabled={updateBudgetMutation.isPending}
                            >
                              Save
                            </Button>
                            <Button
                              onClick={handleCancel}
                              size='sm'
                              variant='outline'
                              className='h-8 px-2'
                            >
                              Cancel
                            </Button>
                          </>
                        ) : (
                          <>
                            <Button
                              onClick={() => handleEdit(budget)}
                              size='sm'
                              variant='outline'
                              className='size-8 p-0'
                            >
                              ✏️
                            </Button>
                            <Button
                              onClick={() =>
                                handleDelete(budget.id, categoryInfo.name)
                              }
                              size='sm'
                              variant='outline'
                              className='size-8 p-0 hover:border-red-300 hover:bg-red-50'
                            >
                              🗑️
                            </Button>
                          </>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                );
              })
            )}
          </TableBody>
        </Table>
      </div>

      {deletingBudget && (
        <DeleteBudgetModal
          isOpen={deleteModalOpen}
          onClose={handleCloseDeleteModal}
          budgetId={deletingBudget.id}
          categoryName={deletingBudget.categoryName}
        />
      )}
    </div>
  );
}
