import { useQuery } from '@tanstack/react-query';
import { assetApi, assetKeys } from '@/entities/asset';
import type { AssetSnapshotListParams } from '@/entities/asset';

export const useAssetSnapshots = (params: AssetSnapshotListParams = {}) => {
  return useQuery({
    queryKey: assetKeys.snapshotList(params),
    queryFn: () => assetApi.getSnapshots(params),
    staleTime: 5 * 60 * 1000, // 5分間キャッシュ
  });
};
