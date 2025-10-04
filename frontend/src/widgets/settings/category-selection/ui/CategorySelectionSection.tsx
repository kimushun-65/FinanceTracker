'use client';

import { Card, Label } from '@/shared/ui';
import type { CategoryWithMaster } from '@/entities/category';
import { useUpdateCategory } from '@/features';
import { useToast } from '@/shared/lib/hooks/use-toast';

interface CategorySelectionSectionProps {
  categories: CategoryWithMaster[];
}

export function CategorySelectionSection({
  categories,
}: CategorySelectionSectionProps) {
  const updateCategoryMutation = useUpdateCategory();
  const { toast } = useToast();

  const activeCategories = categories.filter((cat) => cat.isActive);

  const canDeactivate = (categoryId: string) => {
    // 最後の1つのアクティブカテゴリは無効化不可
    if (
      activeCategories.length === 1 &&
      activeCategories[0].id === categoryId
    ) {
      return false;
    }
    return true;
  };

  const handleToggle = async (categoryId: string, currentIsActive: boolean) => {
    if (currentIsActive && !canDeactivate(categoryId)) {
      toast({
        title: 'Cannot deactivate',
        description: 'At least one category must be active',
        variant: 'destructive',
      });
      return;
    }

    try {
      await updateCategoryMutation.mutateAsync({
        id: categoryId,
        isActive: !currentIsActive,
      });

      toast({
        title: 'Success',
        description: `Category ${!currentIsActive ? 'activated' : 'deactivated'} successfully`,
      });
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to update category',
        variant: 'destructive',
      });
    }
  };

  // カテゴリをタイプ別にグループ化
  const expenseCategories = categories.filter(
    (cat) => cat.master.type === 'expense',
  );
  const incomeCategories = categories.filter(
    (cat) => cat.master.type === 'income',
  );

  const renderCategoryCheckbox = (category: CategoryWithMaster) => {
    const isDisabled = category.isActive && !canDeactivate(category.id);

    return (
      <div
        key={category.id}
        className='flex items-center justify-between rounded-lg border p-3 hover:bg-gray-50'
      >
        <div className='flex items-center gap-3'>
          <input
            type='checkbox'
            id={`category-${category.id}`}
            checked={category.isActive}
            disabled={isDisabled || updateCategoryMutation.isPending}
            onChange={() => handleToggle(category.id, category.isActive)}
            className='size-4 rounded border-gray-300 text-blue-600 focus:ring-2 focus:ring-blue-500 disabled:cursor-not-allowed disabled:opacity-50'
          />
          <label
            htmlFor={`category-${category.id}`}
            className='flex cursor-pointer items-center gap-2'
          >
            <span className='text-2xl'>{category.master.icon}</span>
            <span className='text-sm font-medium text-gray-900'>
              {category.displayName}
            </span>
          </label>
        </div>
        {isDisabled && <span className='text-xs text-gray-500'>Required</span>}
      </div>
    );
  };

  return (
    <Card className='p-6'>
      <div className='mb-4'>
        <h2 className='text-xl font-semibold text-gray-900'>
          Category Selection
        </h2>
      </div>

      {/* Scrollable category list */}
      <div className='max-h-96 space-y-6 overflow-auto rounded-lg border bg-white p-4'>
        {/* Expense Categories */}
        {expenseCategories.length > 0 && (
          <div>
            <Label className='mb-3 block text-sm font-medium text-gray-700'>
              Expense Categories
            </Label>
            <div className='space-y-2'>
              {expenseCategories.map(renderCategoryCheckbox)}
            </div>
          </div>
        )}

        {/* Income Categories */}
        {incomeCategories.length > 0 && (
          <div>
            <Label className='mb-3 block text-sm font-medium text-gray-700'>
              Income Categories
            </Label>
            <div className='space-y-2'>
              {incomeCategories.map(renderCategoryCheckbox)}
            </div>
          </div>
        )}
      </div>
    </Card>
  );
}
