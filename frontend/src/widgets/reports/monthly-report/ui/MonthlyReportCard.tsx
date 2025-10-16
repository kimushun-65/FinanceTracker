'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import type { MonthlyReport, DateRange } from '@/entities/report';
import { format } from 'date-fns';

interface MonthlyReportCardProps {
  report: MonthlyReport | null;
  dateRange: DateRange;
  isLoading?: boolean;
}

export const MonthlyReportCard: React.FC<MonthlyReportCardProps> = ({
  report,
  dateRange,
  isLoading = false,
}) => {
  const formatMoney = (amount: number): string => {
    return `¥${amount.toLocaleString()}`;
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Monthly Report</CardTitle>
        </CardHeader>
        <CardContent>
          <div className='animate-pulse space-y-4'>
            <div className='h-20 rounded bg-gray-200' />
            <div className='h-32 rounded bg-gray-200' />
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!report) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Monthly Report</CardTitle>
        </CardHeader>
        <CardContent>
          <p className='text-muted-foreground'>No data available</p>
        </CardContent>
      </Card>
    );
  }

  const periodLabel = `${format(dateRange.from, 'yyyy年M月')}`;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Monthly Report: {periodLabel}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className='space-y-6'>
          {/* サマリーカード */}
          <div className='grid gap-4 md:grid-cols-3'>
            <div className='rounded-lg border p-4'>
              <p className='text-sm text-muted-foreground'>Income</p>
              <p className='text-2xl font-bold text-green-600'>
                {formatMoney(report.totalIncome.amount)}
              </p>
            </div>
            <div className='rounded-lg border p-4'>
              <p className='text-sm text-muted-foreground'>Expense</p>
              <p className='text-2xl font-bold text-red-600'>
                {formatMoney(report.totalExpenses.amount)}
              </p>
            </div>
            <div className='rounded-lg border p-4'>
              <p className='text-sm text-muted-foreground'>Net Income</p>
              <p className='text-2xl font-bold'>
                {formatMoney(report.netIncome.amount)}
              </p>
            </div>
          </div>

          {/* 詳細情報 */}
          <div className='space-y-2'>
            <p className='text-sm'>
              <span className='font-medium'>Transactions:</span>{' '}
              {report.transactionCount} 件
            </p>
            <p className='text-sm'>
              <span className='font-medium'>Average:</span>{' '}
              {formatMoney(report.averageTransactionAmount.amount)} /
              transaction
            </p>
          </div>

          {/* トップカテゴリ */}
          <div>
            <h4 className='mb-2 font-medium'>Top Expense Categories:</h4>
            <div className='space-y-2'>
              {report.topCategories.slice(0, 5).map((cat, index) => (
                <div
                  key={cat.categoryId}
                  className='flex items-center justify-between text-sm'
                >
                  <span>
                    {index + 1}. {cat.categoryName}
                  </span>
                  <span className='font-medium'>
                    {formatMoney(cat.totalAmount.amount)} (
                    {cat.percentage.toFixed(1)}%)
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
