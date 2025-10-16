'use client';

import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { formatMoney } from '@/shared/value-objects/money';
import type { Account } from '@/entities/account';

interface AssetBreakdownListProps {
  accounts: Account[];
  totalAssets: number;
  isLoading?: boolean;
}

export const AssetBreakdownList: React.FC<AssetBreakdownListProps> = ({
  accounts,
  totalAssets,
  isLoading = false,
}) => {
  const [sortBy, setSortBy] = useState<'name' | 'balance' | 'type'>('balance');

  const calculatePercentage = (amount: number): number => {
    return totalAssets > 0 ? (amount / totalAssets) * 100 : 0;
  };

  const sortedAccounts = [...accounts].sort((a, b) => {
    switch (sortBy) {
      case 'name':
        return a.name.localeCompare(b.name);
      case 'balance':
        return b.balance.current.amount - a.balance.current.amount;
      case 'type':
        return a.accountType.localeCompare(b.accountType);
      default:
        return 0;
    }
  });

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>資産内訳</CardTitle>
        </CardHeader>
        <CardContent>
          <div className='animate-pulse space-y-4'>
            {[1, 2, 3].map((i) => (
              <div key={i} className='h-16 rounded bg-muted' />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
        <CardTitle>資産内訳</CardTitle>
        <div className='flex gap-2'>
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value as any)}
            className='rounded border px-2 py-1 text-sm'
          >
            <option value='balance'>残高順</option>
            <option value='name'>名前順</option>
            <option value='type'>種類順</option>
          </select>
        </div>
      </CardHeader>
      <CardContent>
        {sortedAccounts.length === 0 ? (
          <p className='text-muted-foreground'>口座がまだありません</p>
        ) : (
          <div className='space-y-4'>
            {sortedAccounts.map((account) => (
              <div
                key={account.id}
                className='flex items-center justify-between rounded-lg border p-4'
              >
                <div className='flex-1'>
                  <div className='flex items-center gap-2'>
                    <h4 className='font-medium'>{account.name}</h4>
                    <span className='rounded bg-muted px-2 py-1 text-xs'>
                      {account.accountType}
                    </span>
                  </div>
                  <div className='mt-1 flex items-center gap-4 text-sm text-muted-foreground'>
                    <span>{formatMoney(account.balance.current)}</span>
                    <span>
                      {calculatePercentage(
                        account.balance.current.amount,
                      ).toFixed(1)}
                      %
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
};
