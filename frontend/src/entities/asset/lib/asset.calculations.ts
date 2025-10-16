import type { AssetSnapshot, AssetComparison, AssetTrendData, AssetPeriod } from '../model';

// 前月比計算
export const calculateAssetComparison = (
  current: AssetSnapshot | undefined | null,
  previous: AssetSnapshot | undefined | null
): AssetComparison => {
  if (!current) {
    return {
      current: { amount: 0, currency: 'JPY' },
      previous: { amount: 0, currency: 'JPY' },
      change: { amount: 0, currency: 'JPY' },
      percentageChange: 0,
    };
  }

  if (!previous) {
    return {
      current: current.totalAssets,
      previous: { amount: 0, currency: 'JPY' },
      change: current.totalAssets,
      percentageChange: 0,
    };
  }

  const changeAmount = current.totalAssets.amount - previous.totalAssets.amount;
  const percentageChange = previous.totalAssets.amount
    ? (changeAmount / previous.totalAssets.amount) * 100
    : 0;

  return {
    current: current.totalAssets,
    previous: previous.totalAssets,
    change: { amount: changeAmount, currency: 'JPY' },
    percentageChange,
  };
};

// トレンドデータ変換
export const transformToTrendData = (
  snapshots: AssetSnapshot[]
): AssetTrendData[] => {
  return snapshots
    .sort((a, b) => new Date(a.snapshotDate).getTime() - new Date(b.snapshotDate).getTime())
    .map((snapshot) => ({
      date: snapshot.snapshotDate,
      totalAssets: snapshot.totalAssets.amount,
      accounts: snapshot.accounts.reduce((acc, account) => {
        acc[account.accountId] = account.balance.amount;
        return acc;
      }, {} as { [key: string]: number }),
    }));
};

// 期間フィルタリング
export const filterByPeriod = (
  snapshots: AssetSnapshot[],
  period: AssetPeriod
): AssetSnapshot[] => {
  if (period === 'all') return snapshots;

  const now = new Date();
  const cutoff = new Date(now);

  switch (period) {
    case '1m':
      cutoff.setMonth(now.getMonth() - 1);
      break;
    case '3m':
      cutoff.setMonth(now.getMonth() - 3);
      break;
    case '6m':
      cutoff.setMonth(now.getMonth() - 6);
      break;
    case '1y':
      cutoff.setFullYear(now.getFullYear() - 1);
      break;
  }

  return snapshots.filter(
    (s) => new Date(s.snapshotDate) >= cutoff
  );
};
