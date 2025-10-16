import { useMutation, useQueryClient } from '@tanstack/react-query';
import { assetApi, assetKeys } from '@/entities/asset';
import type { CreateSnapshotPayload } from '@/entities/asset';

export const useCreateAssetSnapshot = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: CreateSnapshotPayload) =>
      assetApi.createSnapshot(payload),
    onSuccess: () => {
      // キャッシュ無効化
      queryClient.invalidateQueries({ queryKey: assetKeys.snapshots() });
      queryClient.invalidateQueries({ queryKey: assetKeys.currentStatus() });
    },
  });
};
