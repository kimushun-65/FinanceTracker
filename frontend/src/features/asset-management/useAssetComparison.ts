import { useMemo } from 'react';
import { calculateAssetComparison } from '@/entities/asset';
import type { AssetSnapshot } from '@/entities/asset';

export const useAssetComparison = (
  current: AssetSnapshot | undefined | null,
  previous: AssetSnapshot | undefined | null
) => {
  return useMemo(
    () => calculateAssetComparison(current, previous),
    [current, previous]
  );
};
