'use client';

import { Select } from '@/shared/ui';
import { useCategories } from '@/features/category-management';
import type { TransactionListParams, TransactionType } from '@/entities/transaction';

interface TransactionFiltersProps {
  filters: TransactionListParams;
  onFiltersChange: (filters: TransactionListParams) => void;
}

export function TransactionFilters({
  filters,
  onFiltersChange,
}: TransactionFiltersProps) {
  const { data: categories } = useCategories();

  const typeOptions = [
    { value: '', label: 'All Types' },
    { value: 'income', label: 'Income' },
    { value: 'expense', label: 'Expense' },
  ];

  const categoryOptions = [
    { value: '', label: 'All Categories' },
    ...(categories?.map((category) => ({
      value: category.id,
      label: category.name,
    })) || []),
  ];

  const periodOptions = [
    { value: 'current', label: 'This Month' },
    { value: 'last', label: 'Last Month' },
    { value: 'last3', label: 'Last 3 Months' },
    { value: 'last6', label: 'Last 6 Months' },
    { value: 'year', label: 'This Year' },
  ];

  const handleTypeChange = (value: string) => {
    onFiltersChange({
      ...filters,
      type: value as TransactionType | undefined,
    });
  };

  const handleCategoryChange = (value: string) => {
    onFiltersChange({
      ...filters,
      categoryId: value || undefined,
    });
  };

  const handlePeriodChange = (value: string) => {
    const now = new Date();
    let startDate: string | undefined;
    let endDate: string | undefined;

    switch (value) {
      case 'current':
        startDate = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().split('T')[0];
        endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().split('T')[0];
        break;
      case 'last':
        startDate = new Date(now.getFullYear(), now.getMonth() - 1, 1).toISOString().split('T')[0];
        endDate = new Date(now.getFullYear(), now.getMonth(), 0).toISOString().split('T')[0];
        break;
      case 'last3':
        startDate = new Date(now.getFullYear(), now.getMonth() - 3, 1).toISOString().split('T')[0];
        endDate = now.toISOString().split('T')[0];
        break;
      case 'last6':
        startDate = new Date(now.getFullYear(), now.getMonth() - 6, 1).toISOString().split('T')[0];
        endDate = now.toISOString().split('T')[0];
        break;
      case 'year':
        startDate = new Date(now.getFullYear(), 0, 1).toISOString().split('T')[0];
        endDate = now.toISOString().split('T')[0];
        break;
      default:
        startDate = undefined;
        endDate = undefined;
    }

    onFiltersChange({
      ...filters,
      startDate,
      endDate,
    });
  };

  return (
    <div className="mb-6 p-4 bg-gray-50 rounded-lg">
      <div className="flex flex-wrap gap-4">
        <div className="min-w-[150px]">
          <Select
            value={filters.type || ''}
            onValueChange={handleTypeChange}
            options={typeOptions}
            placeholder="All Types"
          />
        </div>
        
        <div className="min-w-[200px]">
          <Select
            value={filters.categoryId || ''}
            onValueChange={handleCategoryChange}
            options={categoryOptions}
            placeholder="All Categories"
          />
        </div>
        
        <div className="min-w-[150px]">
          <Select
            value="current" // Default to current month
            onValueChange={handlePeriodChange}
            options={periodOptions}
            placeholder="This Month"
          />
        </div>
      </div>
    </div>
  );
}