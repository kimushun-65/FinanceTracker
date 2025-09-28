'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { useTransactionMonthlySummary } from '@/features/transaction-management';
import { formatMoney } from '@/shared/value-objects/money';

export function TransactionSummaryCards() {
  const { data: summary, isLoading, error } = useTransactionMonthlySummary();

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        {[1, 2, 3].map((i) => (
          <Card key={i} className="animate-pulse">
            <CardHeader className="pb-2">
              <CardTitle className="h-4 bg-gray-200 rounded"></CardTitle>
            </CardHeader>
            <CardContent>
              <div className="h-8 bg-gray-200 rounded"></div>
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        <Card className="border-red-200">
          <CardContent className="p-6">
            <p className="text-red-600 text-sm">Failed to load summary data</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  const totalIncome = summary?.totalIncome || { amount: 0, currency: 'JPY' };
  const totalExpense = summary?.totalExpense || { amount: 0, currency: 'JPY' };
  const netIncome = summary?.netIncome || { amount: 0, currency: 'JPY' };

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
      {/* Monthly Income Card */}
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-gray-600">
            Monthly Income
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-green-600">
            {formatMoney(totalIncome)}
          </div>
        </CardContent>
      </Card>

      {/* Monthly Expense Card */}
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-gray-600">
            Monthly Expense
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-red-600">
            {formatMoney(totalExpense)}
          </div>
        </CardContent>
      </Card>

      {/* Net Income Card */}
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-gray-600">
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