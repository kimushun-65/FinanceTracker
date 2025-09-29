import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { formatMoney } from '@/shared/value-objects/money';
import type { MonthlySummary } from '@/entities/transaction';

interface TransactionSummaryCardsProps {
  summary: MonthlySummary;
}

export function TransactionSummaryCards({
  summary,
}: TransactionSummaryCardsProps) {
  const totalIncome = summary?.totalIncome || { amount: 0, currency: 'JPY' };
  const totalExpense = summary?.totalExpense || { amount: 0, currency: 'JPY' };
  const netIncome = summary?.netIncome || { amount: 0, currency: 'JPY' };

  return (
    <div className='mb-6 grid grid-cols-1 gap-6 md:grid-cols-3'>
      {/* Monthly Income Card */}
      <Card>
        <CardHeader className='pb-2'>
          <CardTitle className='text-sm font-medium text-gray-600'>
            Monthly Income
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-green-600'>
            {formatMoney(totalIncome)}
          </div>
        </CardContent>
      </Card>

      {/* Monthly Expense Card */}
      <Card>
        <CardHeader className='pb-2'>
          <CardTitle className='text-sm font-medium text-gray-600'>
            Monthly Expense
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-red-600'>
            {formatMoney(totalExpense)}
          </div>
        </CardContent>
      </Card>

      {/* Net Income Card */}
      <Card>
        <CardHeader className='pb-2'>
          <CardTitle className='text-sm font-medium text-gray-600'>
            Net Income
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div
            className={`text-2xl font-bold ${
              netIncome.amount >= 0 ? 'text-green-600' : 'text-red-600'
            }`}
          >
            {netIncome.amount >= 0 ? '+' : ''}
            {formatMoney(netIncome)}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
