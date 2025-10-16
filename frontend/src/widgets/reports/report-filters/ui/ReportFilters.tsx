'use client';

import React from 'react';
import { Card, CardContent } from '@/shared/ui';
import type { DateRange, ComparisonPeriod } from '@/entities/report';
import type { CategoryWithMaster } from '@/entities/category';

interface ReportFiltersProps {
  dateRange: DateRange;
  onDateRangeChange: (range: DateRange) => void;
  selectedCategories: string[];
  onCategoriesChange: (categories: string[]) => void;
  categories: CategoryWithMaster[];
  comparisonPeriod: ComparisonPeriod;
  onComparisonPeriodChange: (period: ComparisonPeriod) => void;
}

export const ReportFilters: React.FC<ReportFiltersProps> = ({
  dateRange,
  onDateRangeChange,
  selectedCategories,
  onCategoriesChange,
  categories,
  comparisonPeriod,
  onComparisonPeriodChange,
}) => {
  const handleQuickPeriod = (
    period: 'thisMonth' | 'lastMonth' | 'thisQuarter' | 'thisYear',
  ) => {
    const now = new Date();
    let from: Date;
    let to: Date;

    switch (period) {
      case 'thisMonth':
        from = new Date(now.getFullYear(), now.getMonth(), 1);
        to = new Date(now.getFullYear(), now.getMonth() + 1, 0);
        break;
      case 'lastMonth':
        from = new Date(now.getFullYear(), now.getMonth() - 1, 1);
        to = new Date(now.getFullYear(), now.getMonth(), 0);
        break;
      case 'thisQuarter':
        const quarter = Math.floor(now.getMonth() / 3);
        from = new Date(now.getFullYear(), quarter * 3, 1);
        to = new Date(now.getFullYear(), quarter * 3 + 3, 0);
        break;
      case 'thisYear':
        from = new Date(now.getFullYear(), 0, 1);
        to = new Date(now.getFullYear(), 11, 31);
        break;
    }

    onDateRangeChange({ from, to });
  };

  return (
    <Card>
      <CardContent className='pt-6'>
        <div className='space-y-4'>
          {/* クイック期間選択 */}
          <div>
            <label className='mb-2 block text-sm font-medium'>
              Quick Period
            </label>
            <div className='flex flex-wrap gap-2'>
              <button
                onClick={() => handleQuickPeriod('thisMonth')}
                className='rounded-md border px-3 py-1 text-sm hover:bg-gray-100'
              >
                This Month
              </button>
              <button
                onClick={() => handleQuickPeriod('lastMonth')}
                className='rounded-md border px-3 py-1 text-sm hover:bg-gray-100'
              >
                Last Month
              </button>
              <button
                onClick={() => handleQuickPeriod('thisQuarter')}
                className='rounded-md border px-3 py-1 text-sm hover:bg-gray-100'
              >
                This Quarter
              </button>
              <button
                onClick={() => handleQuickPeriod('thisYear')}
                className='rounded-md border px-3 py-1 text-sm hover:bg-gray-100'
              >
                This Year
              </button>
            </div>
          </div>

          {/* 比較期間選択 */}
          <div>
            <label className='mb-2 block text-sm font-medium'>
              Comparison Period
            </label>
            <div className='flex gap-4'>
              {['month', 'quarter', 'year'].map((period) => (
                <label key={period} className='flex items-center'>
                  <input
                    type='radio'
                    name='comparisonPeriod'
                    value={period}
                    checked={comparisonPeriod === period}
                    onChange={(e) =>
                      onComparisonPeriodChange(
                        e.target.value as ComparisonPeriod,
                      )
                    }
                    className='mr-2'
                  />
                  <span className='capitalize'>{period}</span>
                </label>
              ))}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
