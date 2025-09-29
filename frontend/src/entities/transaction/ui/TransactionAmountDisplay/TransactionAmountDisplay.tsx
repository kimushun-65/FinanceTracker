import type { FC } from 'react';
import type { Money } from '@/shared/value-objects';
import type { TransactionType } from '../../model';
import { formatMoneyWithSign } from '../../lib';
import { formatMoney as formatMoneyVO } from '@/shared/value-objects/money';

export type TransactionAmountDisplayProps = {
  amount: Money;
  type: TransactionType;
  size?: 'sm' | 'default' | 'lg';
  showSign?: boolean;
};

export const TransactionAmountDisplay: FC<TransactionAmountDisplayProps> = ({
  amount,
  type,
  size = 'default',
  showSign = true,
}) => {
  const getColorClass = (transactionType: TransactionType): string => {
    return transactionType === 'income'
      ? 'text-green-600 dark:text-green-400'
      : 'text-red-600 dark:text-red-400';
  };

  const getSizeClass = (displaySize: string): string => {
    const sizeMap = {
      sm: 'text-sm',
      default: 'text-base',
      lg: 'text-lg font-semibold',
    };
    return sizeMap[displaySize as keyof typeof sizeMap] || sizeMap.default;
  };

  const formattedAmount = showSign
    ? formatMoneyWithSign(amount, type)
    : formatMoneyVO({
        amount:
          typeof amount?.amount === 'number'
            ? amount.amount
            : Number((amount as any)?.amount ?? 0),
        currency: (amount as any)?.currency ?? 'JPY',
      });

  return (
    <span className={`${getColorClass(type)} ${getSizeClass(size)} font-mono`}>
      {formattedAmount}
    </span>
  );
};
