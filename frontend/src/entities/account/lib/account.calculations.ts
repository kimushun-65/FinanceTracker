import type { Account, AccountType } from '../model';
import type { Money } from '@/shared/value-objects';

export const canWithdraw = (account: Account, amount: Money): boolean => {
  if (amount.currency !== account.balance.current.currency) {
    return false;
  }

  if (account.accountType === 'cash') {
    return account.balance.current.amount >= amount.amount;
  }

  return true;
};

export const calculateTotalBalance = (accounts: Account[]): Money => {
  if (accounts.length === 0) {
    return { amount: 0, currency: 'JPY' };
  }

  const currency = accounts[0]?.balance.current.currency || 'JPY';
  const total = accounts.reduce((sum, account) => {
    if (account.balance.current.currency !== currency) {
      throw new Error('通貨が異なる口座の合計は計算できません');
    }
    return sum + account.balance.current.amount;
  }, 0);

  return { amount: total, currency };
};

export const calculateBalancesByType = (
  accounts: Account[],
): Map<AccountType, Money> => {
  const balancesByType = new Map<AccountType, Money>();

  accounts.forEach((account) => {
    const currentBalance = balancesByType.get(account.accountType);

    if (!currentBalance) {
      balancesByType.set(account.accountType, {
        amount: account.balance.current.amount,
        currency: account.balance.current.currency,
      });
    } else {
      if (currentBalance.currency !== account.balance.current.currency) {
        throw new Error('通貨が異なる口座の合計は計算できません');
      }
      balancesByType.set(account.accountType, {
        amount: currentBalance.amount + account.balance.current.amount,
        currency: currentBalance.currency,
      });
    }
  });

  return balancesByType;
};

export const calculateGainLossPercentage = (account: Account): number => {
  if (account.balance.initial.amount === 0) {
    return 0;
  }

  return (account.balance.gain.amount / account.balance.initial.amount) * 100;
};
