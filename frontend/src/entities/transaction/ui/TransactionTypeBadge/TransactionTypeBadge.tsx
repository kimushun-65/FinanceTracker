import type { FC } from 'react';
import { Badge } from '@/shared/ui/badge';
import type { TransactionType } from '../../model';
import { formatTransactionType } from '../../lib';

export type TransactionTypeBadgeProps = {
  type: TransactionType;
  size?: 'sm' | 'default';
};

export const TransactionTypeBadge: FC<TransactionTypeBadgeProps> = ({
  type,
  size = 'default',
}) => {
  const variantMap: Record<TransactionType, 'default' | 'destructive'> = {
    income: 'default',
    expense: 'destructive',
  };

  const sizeClass = size === 'sm' ? 'text-xs px-2 py-0.5' : '';

  return (
    <Badge variant={variantMap[type]} className={sizeClass}>
      {formatTransactionType(type)}
    </Badge>
  );
};
