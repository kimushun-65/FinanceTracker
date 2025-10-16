'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { formatMoney } from '@/shared/value-objects/money';
import { useTransactionMonthlySummary } from '@/features/transaction-management';
import type { Money } from '@/shared/value-objects';

interface SummaryCardProps {
  title: string;
  amount: Money;
  trend?: {
    value: number;
    isPositive: boolean;
  };
  isLoading?: boolean;
}

const SummaryCard = ({ title, amount, trend, isLoading }: SummaryCardProps) => {
  if (isLoading) {
    return (
      <Card>
        <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
          <CardTitle className='text-sm font-medium'>{title}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className='h-8 w-32 animate-pulse rounded bg-muted' />
          {trend && (
            <div className='mt-1 h-4 w-24 animate-pulse rounded bg-muted' />
          )}
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
        <CardTitle className='text-sm font-medium'>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className='text-2xl font-bold'>{formatMoney(amount)}</div>
        {trend && (
          <p
            className={`text-xs ${trend.isPositive ? 'text-green-600' : 'text-red-600'}`}
          >
            {trend.isPositive ? '+' : ''}
            {trend.value}% from last month
          </p>
        )}
      </CardContent>
    </Card>
  );
};

export const FinancialSummaryCards = () => {
  const currentDate = new Date();
  const currentYear = currentDate.getFullYear();
  const currentMonth = currentDate.getMonth() + 1;

  const { data: currentSummary, isLoading } = useTransactionMonthlySummary({
    year: currentYear,
    month: currentMonth,
  });

  // 前月のデータ取得（トレンド計算用）
  const previousMonth = currentMonth === 1 ? 12 : currentMonth - 1;
  const previousYear = currentMonth === 1 ? currentYear - 1 : currentYear;

  const { data: previousSummary } = useTransactionMonthlySummary({
    year: previousYear,
    month: previousMonth,
  });

  // トレンドを計算
  const calculateTrend = (current: number, previous: number | undefined) => {
    if (!previous || previous === 0) return undefined;
    const change = ((current - previous) / previous) * 100;
    return {
      value: Math.round(change * 10) / 10, // 小数点1桁まで
      isPositive: change >= 0,
    };
  };

  const incomeTrend = previousSummary
    ? calculateTrend(
        currentSummary?.totalIncome.amount ?? 0,
        previousSummary.totalIncome.amount,
      )
    : undefined;

  const expenseTrend = previousSummary
    ? calculateTrend(
        currentSummary?.totalExpense.amount ?? 0,
        previousSummary.totalExpense.amount,
      )
    : undefined;

  const netTrend = previousSummary
    ? calculateTrend(
        currentSummary?.netIncome.amount ?? 0,
        previousSummary.netIncome.amount,
      )
    : undefined;

  return (
    <div className='grid gap-4 md:grid-cols-2 lg:grid-cols-3'>
      <SummaryCard
        title='今月の収入'
        amount={currentSummary?.totalIncome ?? { amount: 0, currency: 'JPY' }}
        trend={incomeTrend}
        isLoading={isLoading}
      />
      <SummaryCard
        title='今月の支出'
        amount={currentSummary?.totalExpense ?? { amount: 0, currency: 'JPY' }}
        trend={expenseTrend}
        isLoading={isLoading}
      />
      <SummaryCard
        title='今月の収支'
        amount={currentSummary?.netIncome ?? { amount: 0, currency: 'JPY' }}
        trend={netTrend}
        isLoading={isLoading}
      />
    </div>
  );
};
