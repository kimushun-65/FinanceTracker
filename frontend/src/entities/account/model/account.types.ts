import type { BaseEntity } from '@/shared/types';
import type { Money } from '@/shared/value-objects';

export type AccountType = 'checking' | 'investment' | 'cash';

export type BalanceStatus = 'normal' | 'zero' | 'negative';

export type Balance = {
  current: Money;
  initial: Money;
  gain: Money;
  status: BalanceStatus;
};

export type Account = BaseEntity & {
  userId: string;
  name: string;
  accountType: AccountType;
  balance: Balance;
};

export type CreateAccountPayload = {
  name: string;
  accountType: AccountType;
  initialBalance?: Money;
};

export type UpdateAccountPayload = {
  name?: string;
  accountType?: AccountType;
};

export type AccountWithDisplayName = Account & {
  displayTypeName: string;
};

export type AccountListResponse = {
  accounts: Account[];
  total: number;
};
