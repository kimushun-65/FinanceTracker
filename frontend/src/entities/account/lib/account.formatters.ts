import type { Account, AccountType, BalanceStatus } from '../model';
import { ACCOUNT_TYPE_LABELS, ACCOUNT_TYPE_ICONS } from '../model';
import { formatMoney } from '@/shared/value-objects';

export const formatAccountType = (type: AccountType): string => {
  return ACCOUNT_TYPE_LABELS[type];
};

export const getAccountTypeIcon = (type: AccountType): string => {
  return ACCOUNT_TYPE_ICONS[type];
};

export const formatAccountName = (account: Account): string => {
  return `${getAccountTypeIcon(account.accountType)} ${account.name}`;
};

export const formatAccountBalance = (account: Account): string => {
  return formatMoney(account.balance.current);
};

export const formatBalanceStatus = (status: BalanceStatus): string => {
  const statusLabels: Record<BalanceStatus, string> = {
    normal: '正常',
    zero: 'ゼロ',
    negative: 'マイナス',
  };

  return statusLabels[status];
};

export const getBalanceStatusColor = (status: BalanceStatus): string => {
  const statusColors: Record<BalanceStatus, string> = {
    normal: 'text-green-600',
    zero: 'text-gray-500',
    negative: 'text-red-600',
  };

  return statusColors[status];
};

export const formatAccountSummary = (account: Account): string => {
  return `${formatAccountName(account)}: ${formatAccountBalance(account)}`;
};

export type SelectOption = {
  value: string;
  label: string;
};

/**
 * アカウントをSelectオプションに変換
 */
export const accountsToSelectOptions = (
  accounts: Account[],
): SelectOption[] => {
  return (
    accounts?.map((account) => ({
      value: account.id,
      label: account.name,
    })) || []
  );
};
