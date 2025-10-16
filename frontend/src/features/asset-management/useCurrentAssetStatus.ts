import { useQuery } from '@tanstack/react-query';
import { assetApi, assetKeys } from '@/entities/asset';

export const useCurrentAssetStatus = () => {
  return useQuery({
    queryKey: assetKeys.currentStatus(),
    queryFn: () => assetApi.getCurrentStatus(),
    staleTime: 1 * 60 * 1000, // 1分間キャッシュ（リアルタイム性重視）
  });
};
