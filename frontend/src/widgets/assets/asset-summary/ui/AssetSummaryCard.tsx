'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { formatMoney } from '@/shared/value-objects/money';
import type { Money } from '@/shared/value-objects';
import { format } from 'date-fns';

interface AssetSummaryCardProps {
  totalAssets: Money;
  monthOverMonthChange: Money;
  percentageChange: number;
  lastUpdated?: Date;
  isLoading?: boolean;
}

export const AssetSummaryCard: React.FC<AssetSummaryCardProps> = ({
  totalAssets,
  monthOverMonthChange,
  percentageChange,
  lastUpdated,
  isLoading = false,
}) => {
  const isPositive = monthOverMonthChange.amount >= 0;

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>総資産</CardTitle>
        </CardHeader>
        <CardContent>
          <div className='animate-pulse space-y-2'>
            <div className='h-12 w-48 rounded bg-muted' />
            <div className='h-6 w-32 rounded bg-muted' />
            <div className='h-4 w-40 rounded bg-muted' />
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>総資産</CardTitle>
      </CardHeader>
      <CardContent>
        <div className='space-y-2'>
          <p className='text-4xl font-bold'>{formatMoney(totalAssets)}</p>
          <div className='flex items-center gap-2'>
            <span
              className={`font-medium ${
                isPositive ? 'text-green-600' : 'text-red-600'
              }`}
            >
              {isPositive ? '▲' : '▼'} {formatMoney(monthOverMonthChange)} (
              {percentageChange.toFixed(1)}%)
            </span>
            <span className='text-sm text-muted-foreground'>vs 前月</span>
          </div>
          {lastUpdated && (
            <p className='text-sm text-muted-foreground'>
              最終更新: {format(lastUpdated, 'yyyy-MM-dd HH:mm')}
            </p>
          )}
        </div>
      </CardContent>
    </Card>
  );
};
