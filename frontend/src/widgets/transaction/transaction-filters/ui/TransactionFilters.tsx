'use client';

import { Select } from '@/shared/ui';
import type { TransactionListParams } from '@/entities/transaction';
import { Category } from '@/entities/category';
import { periodOptions } from '../config/periodOptions';
import { getPeriodDateRange } from '@/shared/value-objects/date';

interface TransactionFiltersProps {
  filters: TransactionListParams;
  onFiltersChange: (filters: TransactionListParams) => void;
  categories: Category[];
}

export function TransactionFilters({
  filters,
  onFiltersChange,
  categories,
}: TransactionFiltersProps) {
  const categoryOptions = [
    { value: '', label: 'All Categories' },
    ...(categories?.map((category: any) => ({
      value: category.id,
      label: category.name || 'Unknown',
    })) || []),
  ];

  const handleCategoryChange = (value: string) => {
    onFiltersChange({
      ...filters,
      categoryId: value || undefined,
    });
  };

  const handlePeriodChange = (value: string) => {
    const [startDate, endDate] = getPeriodDateRange(value);

    onFiltersChange({
      ...filters,
      startDate,
      endDate,
    });
  };

  return (
    <>
      <div className='min-w-[200px]'>
        <Select
          value={filters.categoryId || ''}
          onValueChange={handleCategoryChange}
          options={categoryOptions}
          placeholder='All Categories'
        />
      </div>

      <div className='min-w-[150px]'>
        <Select
          value='current' // Default to current month
          onValueChange={handlePeriodChange}
          options={periodOptions}
          placeholder='This Month'
        />
      </div>
    </>
  );
}
