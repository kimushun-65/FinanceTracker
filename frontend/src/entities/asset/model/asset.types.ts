import type { Money } from '@/shared/value-objects';

export interface AssetSnapshot {
  id: string;
  userId: string;
  snapshotDate: string; // ISO 8601
  totalAssets: Money;
  accounts: AccountBalance[];
  createdAt: string;
}

export interface AccountBalance {
  accountId: string;
  accountName: string;
  balance: Money;
}

export interface AssetSnapshotListParams {
  from?: string; // YYYY-MM-DD
  to?: string; // YYYY-MM-DD
}

export interface AssetSnapshotListResponse {
  snapshots: AssetSnapshot[];
  totalCount: number;
}

export interface AssetTrendData {
  date: string;
  totalAssets: number;
  accounts: {
    [accountId: string]: number;
  };
}

export interface AssetComparison {
  current: Money;
  previous: Money;
  change: Money;
  percentageChange: number;
}

export interface CreateSnapshotPayload {
  snapshotDate: string; // YYYY-MM-DD
}
