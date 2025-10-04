import type { FC } from 'react';
import { Badge } from '@/shared/ui/badge';
import type { AccountType } from '../../model';
import { formatAccountType } from '../../lib';
import { AccountTypeIcon } from '../AccountTypeIcon';

export type AccountTypeBadgeProps = {
  type: AccountType;
  showIcon?: boolean;
};

export const AccountTypeBadge: FC<AccountTypeBadgeProps> = ({
  type,
  showIcon = true,
}) => {
  const colorMap: Record<AccountType, string> = {
    checking: 'bg-red-100 text-red-800 border-red-200',
    investment: 'bg-blue-100 text-blue-800 border-blue-200',
    cash: 'bg-green-100 text-green-800 border-green-200',
    credit_card: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  };

  return (
    <Badge
      variant='outline'
      className={`inline-flex items-center gap-1 ${colorMap[type]}`}
    >
      {showIcon && <AccountTypeIcon type={type} size={14} />}
      <span>{formatAccountType(type)}</span>
    </Badge>
  );
};
