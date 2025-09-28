import type { Account, AccountWithDisplayName } from '../model';
import { ACCOUNT_TYPE_LABELS } from '../model';

export const transformToAccountWithDisplayName = (
  account: Account,
): AccountWithDisplayName => {
  return {
    ...account,
    displayTypeName: ACCOUNT_TYPE_LABELS[account.accountType],
  };
};

export const transformAccountListToDisplayList = (
  accounts: Account[],
): AccountWithDisplayName[] => {
  return accounts.map(transformToAccountWithDisplayName);
};

export const sortAccountsByType = (accounts: Account[]): Account[] => {
  const typeOrder = { checking: 0, investment: 1, cash: 2 };

  return [...accounts].sort((a, b) => {
    const orderDiff = typeOrder[a.accountType] - typeOrder[b.accountType];
    if (orderDiff !== 0) return orderDiff;

    return a.name.localeCompare(b.name, 'ja');
  });
};

export const sortAccountsByBalance = (
  accounts: Account[],
  descending = true,
): Account[] => {
  return [...accounts].sort((a, b) => {
    const diff = a.balance.current.amount - b.balance.current.amount;
    return descending ? -diff : diff;
  });
};

export const groupAccountsByType = (
  accounts: Account[],
): Map<Account['accountType'], Account[]> => {
  const grouped = new Map<Account['accountType'], Account[]>();

  accounts.forEach((account) => {
    const existing = grouped.get(account.accountType) || [];
    grouped.set(account.accountType, [...existing, account]);
  });

  return grouped;
};
