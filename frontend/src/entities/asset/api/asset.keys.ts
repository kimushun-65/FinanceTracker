import type { AssetSnapshotListParams } from '../model';

export const assetKeys = {
  all: ['assets'] as const,
  snapshots: () => [...assetKeys.all, 'snapshots'] as const,
  snapshotList: (params: AssetSnapshotListParams) =>
    [...assetKeys.snapshots(), params] as const,
  latestSnapshot: () => [...assetKeys.snapshots(), 'latest'] as const,
  currentStatus: () => [...assetKeys.all, 'current'] as const,
} as const;
