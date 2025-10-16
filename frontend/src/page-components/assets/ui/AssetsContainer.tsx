'use client';

import React, { useState } from 'react';
import { subMonths, format } from 'date-fns';
import { AppLayout } from '@/widgets/layout';
import { AssetSummaryCard } from '@/widgets/assets/asset-summary/ui';
import { AssetBreakdownList } from '@/widgets/assets/asset-breakdown/ui';
import { AssetTrendChart } from '@/widgets/assets/asset-trend-chart/ui';
import {
  useAccounts,
  useAccountAggregates,
} from '@/features/account-management';
import {
  useAssetSnapshots,
  useLatestAssetSnapshot,
  useCurrentAssetStatus,
  useAssetComparison,
  useAssetTrendData,
} from '@/features/asset-management';
import type { Account } from '@/entities/account';
import type { AssetPeriod } from '@/entities/asset';

export const AssetsContainer: React.FC = () => {
  // データ取得
  const { data: accounts, isLoading: accountsLoading } = useAccounts();
  const { data: latestSnapshot } = useLatestAssetSnapshot();
  const { data: currentStatus, isLoading: currentStatusLoading } =
    useCurrentAssetStatus();

  const sixMonthsAgo = format(subMonths(new Date(), 6), 'yyyy-MM-dd');
  const today = format(new Date(), 'yyyy-MM-dd');

  const { data: snapshotsData } = useAssetSnapshots({
    from: sixMonthsAgo,
    to: today,
  });

  // 状態管理
  const [editingAccount, setEditingAccount] = useState<Account | null>(null);
  const [deletingAccount, setDeletingAccount] = useState<Account | null>(null);
  const [activeTab, setActiveTab] = useState<'list' | 'chart'>('list');
  const [selectedPeriod, setSelectedPeriod] = useState<AssetPeriod>('6m');

  // 計算処理
  const { totalAssets } = useAccountAggregates(accounts || []);
  const comparison = useAssetComparison(currentStatus, latestSnapshot);
  const trendData = useAssetTrendData(snapshotsData?.snapshots, selectedPeriod);

  return (
    <AppLayout>
      <div className='mb-8'>
        <div className='flex items-center justify-between'>
          <div>
            <h2 className='mb-2 text-3xl font-bold text-gray-900'>資産管理</h2>
          </div>
        </div>
      </div>

      <div className='space-y-6'>
        {/* 総資産サマリー */}
        <AssetSummaryCard
          totalAssets={comparison.current}
          monthOverMonthChange={comparison.change}
          percentageChange={comparison.percentageChange}
          lastUpdated={
            currentStatus?.createdAt
              ? new Date(currentStatus.createdAt)
              : undefined
          }
          isLoading={currentStatusLoading}
        />

        {/* タブ切り替え */}
        <div className='flex gap-2 border-b'>
          <button
            onClick={() => setActiveTab('list')}
            className={`px-4 py-2 font-medium transition-colors ${
              activeTab === 'list'
                ? 'border-b-2 border-blue-600 text-blue-600'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            リスト表示
          </button>
          <button
            onClick={() => setActiveTab('chart')}
            className={`px-4 py-2 font-medium transition-colors ${
              activeTab === 'chart'
                ? 'border-b-2 border-blue-600 text-blue-600'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            グラフ表示
          </button>
        </div>

        {/* コンテンツ表示 */}
        {activeTab === 'list' ? (
          <AssetBreakdownList
            accounts={accounts || []}
            totalAssets={totalAssets.amount}
            isLoading={accountsLoading}
          />
        ) : (
          <AssetTrendChart
            data={trendData}
            period={selectedPeriod}
            onPeriodChange={setSelectedPeriod}
            isLoading={accountsLoading}
          />
        )}
      </div>
    </AppLayout>
  );
};
