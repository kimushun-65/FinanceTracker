import { useMemo } from 'react';
import type { Account } from '@/entities/account';
import type { Money, Currency } from '@/shared/value-objects';

export type AccountAggregates = {
  totalAssets: Money;
  accountCount: number;
  accountsByType: Record<string, Account[]>;
};

export const useAccountAggregates = (accounts: Account[] | undefined) => {
  return useMemo((): AccountAggregates => {
    if (!accounts || accounts.length === 0) {
      return {
        totalAssets: { amount: 0, currency: 'JPY' as Currency },
        accountCount: 0,
        accountsByType: {},
      };
    }

    // 総資産計算
    const totalAssets = accounts.reduce(
      (sum, account) => ({
        amount: sum.amount + account.balance.current.amount,
        currency: 'JPY' as Currency,
      }),
      { amount: 0, currency: 'JPY' as Currency } as Money,
    );

    // タイプ別グループ化
    const accountsByType = accounts.reduce(
      (groups, account) => {
        const type = account.accountType;
        if (!groups[type]) {
          groups[type] = [];
        }
        groups[type].push(account);
        return groups;
      },
      {} as Record<string, Account[]>,
    );

    return {
      totalAssets,
      accountCount: accounts.length,
      accountsByType,
    };
  }, [accounts]);
};
