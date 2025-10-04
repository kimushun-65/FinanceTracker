import type { BaseEntity } from '@/shared/types';
import type { Money } from '@/shared/value-objects';

// Backend API actual response types (from API DTOs)
export type BackendAccountType =
  | 'checking'
  | 'investment'
  | 'cash'
  | 'credit_card';

// Frontend display types (same as backend for now)
export type AccountType = 'checking' | 'investment' | 'cash' | 'credit_card';

// Backend API response format
export type BackendAccountResponse = {
  id: string;
  user_id: string;
  name: string;
  account_type: BackendAccountType;
  balance: string; // Decimal string like "1500000.00"
  currency: string;
  created_at: string;
  updated_at: string;
};

export type BalanceStatus = 'normal' | 'zero' | 'negative';

export type Balance = {
  current: Money;
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
  balance?: Money;
};

export type AccountWithDisplayName = Account & {
  displayTypeName: string;
};

// Backend API response format
export type BackendAccountListResponse = {
  accounts: BackendAccountResponse[];
  total_count: number;
  total_balance: string;
  total_debt?: string;
  net_worth?: string;
};

// Frontend format
export type AccountListResponse = {
  accounts: Account[];
  total: number;
  totalAssets?: number;
  totalDebt?: number;
  netWorth?: number;
};
