import type { Account, AccountWithDisplayName, BackendAccountResponse, BackendAccountType, AccountType } from '../model';
import { ACCOUNT_TYPE_LABELS } from '../model';
import type { Money } from '@/shared/value-objects';

// Backend account type to frontend type mapping (direct mapping since they're the same)
const BACKEND_TO_FRONTEND_TYPE_MAP: Record<BackendAccountType, AccountType> = {
  'checking': 'checking',
  'investment': 'investment',
  'cash': 'cash',
};

// Transform backend API response to frontend Account type
export const transformBackendAccountToFrontend = (
  backendAccount: BackendAccountResponse,
): Account => {
  const balanceAmount = parseFloat(backendAccount.balance);
  
  const currentBalance: Money = {
    amount: balanceAmount,
    currency: backendAccount.currency,
  };

  return {
    id: backendAccount.id,
    userId: backendAccount.user_id,
    name: backendAccount.name,
    accountType: BACKEND_TO_FRONTEND_TYPE_MAP[backendAccount.account_type],
    balance: {
      current: currentBalance,
      status: balanceAmount > 0 ? 'normal' : balanceAmount === 0 ? 'zero' : 'negative',
    },
    createdAt: backendAccount.created_at,
    updatedAt: backendAccount.updated_at,
  };
};

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
