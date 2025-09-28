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
  const variantMap: Record<AccountType, 'default' | 'secondary' | 'outline'> = {
    checking: 'default',
    investment: 'secondary',
    cash: 'outline',
  };

  return (
    <Badge
      variant={variantMap[type]}
      className='inline-flex items-center gap-1'
    >
      {showIcon && <AccountTypeIcon type={type} size={14} />}
      <span>{formatAccountType(type)}</span>
    </Badge>
  );
};
