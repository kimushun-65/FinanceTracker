import { formatMoney } from '@/shared/value-objects/money';
import type { Money } from '@/shared/value-objects';
import type { BalanceStatus } from '@/entities/account';

interface AccountBalanceDisplayProps {
  balance: Money;
  status: BalanceStatus;
  className?: string;
}

export const AccountBalanceDisplay = ({
  balance,
  status,
  className = '',
}: AccountBalanceDisplayProps) => {
  const getStatusColor = (status: BalanceStatus): string => {
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

  return (
    <span className={`font-medium ${getStatusColor(status)} ${className}`}>
      {formatMoney(balance)}
    </span>
  );
};