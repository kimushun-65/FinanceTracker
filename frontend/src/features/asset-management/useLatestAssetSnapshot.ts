import { useQuery } from '@tanstack/react-query';
import { assetApi, assetKeys } from '@/entities/asset';

export const useLatestAssetSnapshot = () => {
  return useQuery({
    queryKey: assetKeys.latestSnapshot(),
    queryFn: () => assetApi.getLatestSnapshot(),
    staleTime: 5 * 60 * 1000,
  });
};
