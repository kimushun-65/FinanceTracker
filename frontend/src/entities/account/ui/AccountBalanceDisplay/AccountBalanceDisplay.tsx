import { formatMoney } from '@/shared/value-objects/money';
import type { Money } from '@/shared/value-objects';
import type { BalanceStatus, AccountType } from '@/entities/account';

interface AccountBalanceDisplayProps {
  balance: Money;
  status: BalanceStatus;
  accountType?: AccountType;
  className?: string;
}

export const AccountBalanceDisplay = ({
  balance,
  status,
  accountType,
  className = '',
}: AccountBalanceDisplayProps) => {
  const getStatusColor = (
    status: BalanceStatus,
    accountType?: AccountType,
  ): string => {
    // For credit cards, negative balance is normal (represents debt)
    if (accountType === 'credit_card') {
      if (balance.amount < 0) {
        return 'text-orange-600'; // Debt amount
      } else if (balance.amount === 0) {
        return 'text-green-600'; // Paid off
      } else {
        return 'text-green-600'; // Overpaid (credit available)
      }
    }

    // For regular accounts
    switch (status) {
      case 'normal':
        return 'text-gray-900';
      case 'zero':
        return 'text-gray-500';
      case 'negative':
        return 'text-red-600';
      default:
        return 'text-gray-900';
    }
  };

  const displayBalance =
    accountType === 'credit_card' && balance.amount < 0
      ? { ...balance, amount: Math.abs(balance.amount) }
      : balance;

  const prefix = accountType === 'credit_card' && balance.amount < 0 ? '-' : '';

  return (
    <span
      className={`font-medium ${getStatusColor(status, accountType)} ${className}`}
    >
      {prefix}
      {formatMoney(displayBalance)}
    </span>
  );
};
