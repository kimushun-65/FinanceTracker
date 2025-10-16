import { apiClient } from '@/shared/api/client';
import type {
  AssetSnapshot,
  AssetSnapshotListParams,
  AssetSnapshotListResponse,
  CreateSnapshotPayload,
} from '../model';
import { assetEndpoints } from './asset.endpoints';

export const assetApi = {
  // スナップショット一覧取得
  getSnapshots: async (
    params: AssetSnapshotListParams,
  ): Promise<AssetSnapshotListResponse> => {
    const queryParams = new URLSearchParams();
    if (params.from) queryParams.append('from', params.from);
    if (params.to) queryParams.append('to', params.to);

    const url = queryParams.toString()
      ? `${assetEndpoints.snapshots}?${queryParams.toString()}`
      : assetEndpoints.snapshots;

    const response = await apiClient.get<any>(url);

    const snapshots: AssetSnapshot[] = Array.isArray(response.data?.snapshots)
      ? response.data.snapshots.map((snap: any) => ({
          id: snap.id || snap.snapshot_id,
          userId: snap.user_id,
          snapshotDate: snap.snapshot_date,
          totalAssets: {
            amount: Number(snap.total_assets || 0),
            currency: 'JPY',
          },
          accounts: Array.isArray(snap.accounts)
            ? snap.accounts.map((acc: any) => ({
                accountId: acc.account_id,
                accountName: acc.account_name,
                balance: {
                  amount: Number(acc.balance || 0),
                  currency: 'JPY',
                },
              }))
            : [],
          createdAt: snap.created_at,
        }))
      : [];

    return {
      snapshots,
      totalCount: response.data?.total_count || snapshots.length,
    };
  },

  // 最新スナップショット取得
  getLatestSnapshot: async (): Promise<AssetSnapshot | null> => {
    try {
      const response = await apiClient.get<any>(assetEndpoints.latestSnapshot);
      const snap = response.data;

      if (!snap) return null;

      return {
        id: snap.id || snap.snapshot_id,
        userId: snap.user_id,
        snapshotDate: snap.snapshot_date,
        totalAssets: {
          amount: Number(snap.total_assets || 0),
          currency: 'JPY',
        },
        accounts: Array.isArray(snap.accounts)
          ? snap.accounts.map((acc: any) => ({
              accountId: acc.account_id,
              accountName: acc.account_name,
              balance: {
                amount: Number(acc.balance || 0),
                currency: 'JPY',
              },
            }))
          : [],
        createdAt: snap.created_at,
      };
    } catch (error) {
      return null;
    }
  },

  // 現在の資産状況計算
  getCurrentStatus: async (): Promise<AssetSnapshot> => {
    const response = await apiClient.get<any>(assetEndpoints.currentStatus);
    const snap = response.data;

    return {
      id: snap.id || 'current',
      userId: snap.user_id,
      snapshotDate: snap.snapshot_date || new Date().toISOString(),
      totalAssets: {
        amount: Number(snap.total_assets || 0),
        currency: 'JPY',
      },
      accounts: Array.isArray(snap.accounts)
        ? snap.accounts.map((acc: any) => ({
            accountId: acc.account_id,
            accountName: acc.account_name,
            balance: {
              amount: Number(acc.balance || 0),
              currency: 'JPY',
            },
          }))
        : [],
      createdAt: snap.created_at || new Date().toISOString(),
    };
  },

  // スナップショット作成
  createSnapshot: async (
    payload: CreateSnapshotPayload,
  ): Promise<AssetSnapshot> => {
    // まず現在の資産状況を取得
    const currentStatusResponse = await apiClient.get<any>(
      assetEndpoints.currentStatus,
    );
    const currentStatus = currentStatusResponse.data;

    // 口座データを整形
    const accounts = Array.isArray(currentStatus.accounts)
      ? currentStatus.accounts.map((acc: any) => ({
          account_id: acc.account_id,
          balance: Number(acc.balance || 0),
        }))
      : [];

    // スナップショット作成リクエスト
    const response = await apiClient.post<any>(assetEndpoints.createSnapshot, {
      snapshot_date: payload.snapshotDate,
      accounts: accounts,
    });

    const snap = response.data;

    return {
      id: snap.id || snap.snapshot_id,
      userId: snap.user_id,
      snapshotDate: snap.snapshot_date,
      totalAssets: {
        amount: Number(snap.total_assets || 0),
        currency: 'JPY',
      },
      accounts: Array.isArray(snap.accounts)
        ? snap.accounts.map((acc: any) => ({
            accountId: acc.account_id,
            accountName: acc.account_name,
            balance: {
              amount: Number(acc.balance || 0),
              currency: 'JPY',
            },
          }))
        : [],
      createdAt: snap.created_at,
    };
  },
} as const;

export type AssetApi = typeof assetApi;
