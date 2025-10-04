import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { formatMoney } from '@/shared/value-objects/money';
import type { BudgetSummary } from '@/entities/budget';

interface BudgetOverviewCardsProps {
  summary: BudgetSummary;
}

export function BudgetOverviewCards({ summary }: BudgetOverviewCardsProps) {
  const totalBudget = summary?.totalBudget || { amount: 0, currency: 'JPY' };
  const totalUsed = summary?.totalUsed || { amount: 0, currency: 'JPY' };
  const totalRemaining = summary?.totalRemaining || {
    amount: 0,
    currency: 'JPY',
  };

  const usagePercentage =
    totalBudget.amount > 0
      ? Math.round((totalUsed.amount / totalBudget.amount) * 100)
      : 0;

  return (
    <div className='mb-6 grid grid-cols-1 gap-6 md:grid-cols-3'>
      {/* Total Budget Card */}
      <Card>
        <CardHeader className='pb-2'>
          <CardTitle className='text-sm font-medium text-gray-600'>
            Total Budget
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-blue-600'>
            {formatMoney(totalBudget)}
          </div>
          <p className='mt-1 text-xs text-gray-500'>
            Active Budgets: {summary?.activeCount || 0}
          </p>
        </CardContent>
      </Card>

      {/* Total Used Card */}
      <Card>
        <CardHeader className='pb-2'>
          <CardTitle className='text-sm font-medium text-gray-600'>
            Total Used
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-orange-600'>
            {formatMoney(totalUsed)}
          </div>
          <p className='mt-1 text-xs text-gray-500'>
            {usagePercentage}% of budget
          </p>
        </CardContent>
      </Card>

      {/* Total Remaining Card */}
      <Card>
        <CardHeader className='pb-2'>
          <CardTitle className='text-sm font-medium text-gray-600'>
            Total Remaining
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div
            className={`text-2xl font-bold ${
              totalRemaining.amount >= 0 ? 'text-green-600' : 'text-red-600'
            }`}
          >
            {formatMoney(totalRemaining)}
          </div>
          <p className='mt-1 text-xs text-gray-500'>
            Over Budget: {summary?.overBudgetCount || 0}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
