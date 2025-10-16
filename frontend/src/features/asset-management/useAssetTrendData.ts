import { useMemo } from 'react';
import { transformToTrendData, filterByPeriod } from '@/entities/asset';
import type { AssetSnapshot, AssetPeriod } from '@/entities/asset';

export const useAssetTrendData = (
  snapshots: AssetSnapshot[] | undefined,
  period: AssetPeriod
) => {
  return useMemo(() => {
    if (!snapshots) return [];

    const filtered = filterByPeriod(snapshots, period);
    return transformToTrendData(filtered);
  }, [snapshots, period]);
};
