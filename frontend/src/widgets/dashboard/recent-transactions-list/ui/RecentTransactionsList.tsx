'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { formatMoney } from '@/shared/value-objects/money';
import { formatDateShort } from '@/shared/value-objects/date';
import { useTransactions } from '@/features/transaction-management';
import type { Transaction } from '@/entities/transaction';

interface TransactionItemProps {
  transaction: Transaction;
}

const TransactionItem = ({ transaction }: TransactionItemProps) => {
  const isIncome = transaction.type === 'income';

  return (
    <div className='flex items-center justify-between border-b py-3 last:border-b-0'>
      <div className='flex-1'>
        <p className='font-medium'>{transaction.description || '(説明なし)'}</p>
        <p className='text-sm text-muted-foreground'>
          {formatDateShort(transaction.date)}
        </p>
      </div>
      <div
        className={`font-semibold ${isIncome ? 'text-green-600' : 'text-red-600'}`}
      >
        {isIncome ? '+' : '-'}
        {formatMoney(transaction.amount)}
      </div>
    </div>
  );
};

const LoadingSkeleton = () => (
  <div className='space-y-3'>
    {[1, 2, 3, 4, 5].map((i) => (
      <div key={i} className='flex items-center justify-between py-3'>
        <div className='flex-1 space-y-2'>
          <div className='h-5 w-32 animate-pulse rounded bg-muted' />
          <div className='h-4 w-20 animate-pulse rounded bg-muted' />
        </div>
        <div className='h-6 w-24 animate-pulse rounded bg-muted' />
      </div>
    ))}
  </div>
);

const EmptyState = () => (
  <div className='py-8 text-center text-muted-foreground'>
    <p>最近の取引はありません</p>
  </div>
);

export const RecentTransactionsList = () => {
  const { data, isLoading } = useTransactions({
    limit: 5,
    sortBy: 'date',
    sortOrder: 'desc',
  });

  const transactions = data?.transactions ?? [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>最近の取引</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <LoadingSkeleton />
        ) : transactions.length === 0 ? (
          <EmptyState />
        ) : (
          <div>
            {transactions.map((transaction) => (
              <TransactionItem key={transaction.id} transaction={transaction} />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
};
