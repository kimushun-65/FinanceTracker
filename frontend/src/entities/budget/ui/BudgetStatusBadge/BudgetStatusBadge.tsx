import type { FC } from 'react';
import { Badge } from '@/shared/ui/badge';
import type { BudgetStatus } from '../../model';
import { BUDGET_STATUS } from '../../model';

export type BudgetStatusBadgeProps = {
  status: BudgetStatus;
  size?: 'sm' | 'default';
  className?: string;
};

export const BudgetStatusBadge: FC<BudgetStatusBadgeProps> = ({
  status,
  size = 'default',
  className = '',
}) => {
  const getVariant = (budgetStatus: BudgetStatus) => {
    switch (budgetStatus) {
      case 'normal':
        return 'default' as const;
      case 'warning':
        return 'secondary' as const;
      case 'exceeded':
        return 'destructive' as const;
      default:
        return 'default' as const;
    }
  };

  const sizeClass = size === 'sm' ? 'text-xs px-2 py-0.5' : '';

  return (
    <Badge variant={getVariant(status)} className={`${sizeClass} ${className}`}>
      {BUDGET_STATUS[status]}
    </Badge>
  );
};
