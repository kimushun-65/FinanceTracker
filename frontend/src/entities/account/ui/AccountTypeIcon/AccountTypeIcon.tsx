import type { FC } from 'react';
import { Wallet, TrendingUp, Banknote, CreditCard } from 'lucide-react';
import type { AccountType } from '../../model';

export type AccountTypeIconProps = {
  type: AccountType;
  className?: string;
  size?: number;
};

export const AccountTypeIcon: FC<AccountTypeIconProps> = ({
  type,
  className = '',
  size = 20,
}) => {
  const iconProps = {
    className,
    size,
  };

  const icons: Record<AccountType, React.ReactNode> = {
    checking: <Banknote {...iconProps} />,
    investment: <TrendingUp {...iconProps} />,
    cash: <Wallet {...iconProps} />,
    credit_card: <CreditCard {...iconProps} />,
  };

  return icons[type] || null;
};
