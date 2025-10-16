# 資産管理ページ実装計画

## 概要
この資料では、FinanceTrackerフロントエンドアプリケーションにおける資産管理ページの完全な実装アーキテクチャについて説明します。このページは**完全に新規作成**が必要です。

## ページ要件（要件定義より）

### 主要機能
1. **総資産表示**: 全口座の合計資産額と前月比
2. **資産内訳リスト**: 口座別の残高・構成比・前月比増減
3. **口座管理機能**: 口座のCRUD操作
4. **資産推移グラフ**: 時系列での資産変動可視化
5. **資産予測機能**: 1年後の資産予測（将来実装予定）

### 対応API
- ✅ `GET /api/v1/accounts` - 口座一覧
- ✅ `POST /api/v1/accounts` - 口座作成
- ✅ `PUT /api/v1/accounts/:id` - 口座更新
- ✅ `DELETE /api/v1/accounts/:id` - 口座削除
- ✅ `GET /api/v1/reports/assets/snapshots` - 資産スナップショット一覧
- ✅ `GET /api/v1/reports/assets/snapshots/latest` - 最新スナップショット
- ✅ `GET /api/v1/reports/assets/snapshots/current` - 現在の資産状況
- ✅ `POST /api/v1/reports/assets/snapshots` - スナップショット作成
- ❌ `POST /api/v1/assets/forecasts` - 資産予測（未実装）

## 完全なディレクトリ構成

```
frontend/src/
├── app/
│   └── assets/                               # 新規ディレクトリ
│       └── page.tsx                          # エントリーポイント（新規）
├── page-components/
│   └── assets/                               # 新規ディレクトリ
│       └── ui/
│           └── AssetsContainer.tsx           # メインオーケストレーター（新規）
├── widgets/
│   ├── assets/                               # 新規ディレクトリ
│   │   ├── asset-summary/
│   │   │   └── ui/
│   │   │       └── AssetSummaryCard.tsx      # 総資産表示カード（新規）
│   │   ├── asset-breakdown/
│   │   │   └── ui/
│   │   │       └── AssetBreakdownList.tsx    # 資産内訳リスト（新規）
│   │   ├── asset-trend-chart/
│   │   │   └── ui/
│   │   │       └── AssetTrendChart.tsx       # 資産推移グラフ（新規）
│   │   ├── snapshot-management/
│   │   │   └── ui/
│   │   │       └── SnapshotManagementWidget.tsx # スナップショット管理（新規）
│   │   └── asset-forecast/                   # 将来実装予定
│   │       └── ui/
│   │           └── AssetForecastWidget.tsx   # 資産予測ウィジェット（新規）
│   └── account/                              # 既存ディレクトリ（一部流用）
│       ├── create-account/
│       │   └── CreateAccountModal.tsx        # 既存流用
│       ├── edit-account/
│       │   └── EditAccountModal.tsx          # 既存流用
│       └── delete-account/
│           └── DeleteAccountModal.tsx        # 既存流用
├── entities/
│   ├── account/                              # 既存（流用）
│   │   ├── model/
│   │   │   ├── account.types.ts              # 既存
│   │   │   └── index.ts
│   │   ├── api/
│   │   │   ├── account.client.ts             # 既存
│   │   │   └── index.ts
│   │   └── index.ts
│   └── asset/                                # 新規ディレクトリ
│       ├── model/
│       │   ├── asset.types.ts                # 型定義（新規）
│       │   ├── asset.schema.ts               # バリデーションスキーマ（新規）
│       │   ├── asset.constants.ts            # 定数（新規）
│       │   └── index.ts
│       ├── api/
│       │   ├── asset.client.ts               # APIクライアント（新規）
│       │   ├── asset.endpoints.ts            # エンドポイント定義（新規）
│       │   ├── asset.keys.ts                 # React Queryキー（新規）
│       │   └── index.ts
│       ├── lib/
│       │   ├── asset.transformers.ts         # データ変換（新規）
│       │   ├── asset.formatters.ts           # フォーマッター（新規）
│       │   ├── asset.calculations.ts         # 計算処理（新規）
│       │   └── index.ts
│       ├── ui/
│       │   ├── AssetSnapshotBadge/
│       │   │   ├── AssetSnapshotBadge.tsx    # スナップショットバッジ（新規）
│       │   │   └── index.ts
│       │   ├── AssetTrendIndicator/
│       │   │   ├── AssetTrendIndicator.tsx   # トレンド表示（新規）
│       │   │   └── index.ts
│       │   └── index.ts
│       └── index.ts
├── features/
│   ├── account-management/                   # 既存（流用）
│   │   ├── useAccounts.ts                    # 既存
│   │   ├── useCreateAccount.ts               # 既存
│   │   ├── useUpdateAccount.ts               # 既存
│   │   ├── useDeleteAccount.ts               # 既存
│   │   ├── useAccountAggregates.ts           # 既存
│   │   └── index.ts
│   └── asset-management/                     # 新規ディレクトリ
│       ├── useAssetSnapshots.ts              # スナップショット一覧（新規）
│       ├── useLatestAssetSnapshot.ts         # 最新スナップショット（新規）
│       ├── useCurrentAssetStatus.ts          # 現在の資産状況（新規）
│       ├── useCreateAssetSnapshot.ts         # スナップショット作成（新規）
│       ├── useAssetTrendData.ts              # 推移データ処理（新規）
│       ├── useAssetComparison.ts             # 前月比計算（新規）
│       └── index.ts
└── shared/
    ├── ui/
    │   ├── card.tsx                          # 既存
    │   ├── table.tsx                         # 既存
    │   ├── loading.tsx                       # 既存
    │   └── index.ts
    └── lib/
        └── chart/
            ├── chartConfig.ts                # 既存
            └── index.ts
```

## アーキテクチャ階層

### 1. ページ層
**エントリーポイント（新規作成）:**
```typescript
// /app/assets/page.tsx

import { AssetsContainer } from '../../page-components/assets';
import { ProtectedRoute } from '../../shared/ui/components/ProtectedRoute';

export default function AssetsPage() {
  return (
    <ProtectedRoute>
      <AssetsContainer />
    </ProtectedRoute>
  );
}
```

### 2. ページコンポーネント層
**メインコンテナ（新規作成）:**
```typescript
// /page-components/assets/ui/AssetsContainer.tsx

export const AssetsContainer: React.FC = () => {
  // データ取得
  const { data: accounts } = useAccounts();
  const { data: latestSnapshot } = useLatestAssetSnapshot();
  const { data: currentStatus } = useCurrentAssetStatus();
  const { data: snapshots } = useAssetSnapshots({
    from: sixMonthsAgo,
    to: today,
  });

  // 状態管理
  const [selectedPeriod, setSelectedPeriod] = useState('6m');
  const [editingAccount, setEditingAccount] = useState<Account | null>(null);

  // 計算処理
  const { totalAssets, accountCount } = useAccountAggregates(accounts);
  const { monthOverMonthChange, percentageChange } = useAssetComparison(
    latestSnapshot,
    currentStatus
  );
  const trendData = useAssetTrendData(snapshots, selectedPeriod);

  return (
    <AppLayout title='Assets Management'>
      <div className='space-y-6'>
        {/* 総資産サマリー */}
        <AssetSummaryCard
          totalAssets={totalAssets}
          monthOverMonthChange={monthOverMonthChange}
          percentageChange={percentageChange}
          lastUpdated={currentStatus?.createdAt}
        />

        {/* 資産内訳リスト */}
        <AssetBreakdownList
          accounts={accounts || []}
          totalAssets={totalAssets}
          onEditAccount={setEditingAccount}
          onDeleteAccount={handleDeleteAccount}
        />

        {/* 資産推移グラフ */}
        <AssetTrendChart
          data={trendData}
          period={selectedPeriod}
          onPeriodChange={setSelectedPeriod}
        />

        {/* スナップショット管理 */}
        <SnapshotManagementWidget
          latestSnapshot={latestSnapshot}
          onCreateSnapshot={handleCreateSnapshot}
        />

        {/* 資産予測（将来実装）*/}
        {/* <AssetForecastWidget /> */}
      </div>

      {/* Modals */}
      <CreateAccountModal ... />
      <EditAccountModal ... />
      <DeleteAccountModal ... />
    </AppLayout>
  );
};
```

**主要な責任範囲:**
- 全ウィジェットのオーケストレーション
- データ取得とキャッシュ管理
- モーダル状態管理
- エラーハンドリング
- ローディング状態管理

### 3. ウィジェット層

#### 3.1 AssetSummaryCard（新規）
```typescript
// /widgets/assets/asset-summary/ui/AssetSummaryCard.tsx

interface Props {
  totalAssets: number;
  monthOverMonthChange: number;
  percentageChange: number;
  lastUpdated?: Date;
}

export const AssetSummaryCard: React.FC<Props>
```

**責任範囲:**
- 総資産額の大型表示
- 前月比の増減表示（矢印アイコン + 金額 + パーセント）
- 最終更新日時の表示
- トレンド色分け（増加=緑、減少=赤）

**UIデザイン:**
```
┌────────────────────────────────────┐
│  Total Assets                      │
│                                    │
│  ¥ 1,500,000                      │
│                                    │
│  ▲ +¥50,000 (+3.4%)               │
│  vs last month                     │
│                                    │
│  Last updated: 2025-01-09 14:30   │
└────────────────────────────────────┘
```

#### 3.2 AssetBreakdownList（新規）
```typescript
// /widgets/assets/asset-breakdown/ui/AssetBreakdownList.tsx

interface Props {
  accounts: Account[];
  totalAssets: number;
  onEditAccount: (account: Account) => void;
  onDeleteAccount: (account: Account) => void;
}

export const AssetBreakdownList: React.FC<Props>
```

**責任範囲:**
- 口座一覧のテーブル表示
- 口座種別アイコン表示
- 現在残高・構成比率・前月比の表示
- 編集・削除アクションボタン
- 口座追加ボタン
- ソート機能（残高順、名前順、種別順）
- グルーピング表示（種別ごと）

**テーブル列:**
| 口座名 | 種別 | 残高 | 構成比 | 前月比 | アクション |
|--------|------|------|--------|--------|-----------|
| みずほ銀行 | 🏛️ 普通預金 | ¥1,200,000 | 80% | ▲ +¥20k | [Edit][Delete] |

#### 3.3 AssetTrendChart（新規）
```typescript
// /widgets/assets/asset-trend-chart/ui/AssetTrendChart.tsx

interface Props {
  data: AssetTrendData[];
  period: '1m' | '3m' | '6m' | '1y' | 'all';
  onPeriodChange: (period: string) => void;
}

export const AssetTrendChart: React.FC<Props>
```

**責任範囲:**
- 資産推移の折れ線グラフ表示
- 期間切り替え（1ヶ月、3ヶ月、6ヶ月、1年、全期間）
- ホバー時の詳細表示
- レスポンシブ対応
- データポイントクリックで詳細モーダル表示

**使用ライブラリ:**
- Recharts（LineChart）

**グラフ表示:**
```
Total Asset Trend          [1M] [3M] [6M] [1Y] [All]
┌────────────────────────────────────────────┐
│ 1.6M│                            ●         │
│     │                        ●───◦         │
│ 1.5M│                   ●───◦              │
│     │              ●───◦                   │
│ 1.4M│         ●───◦                        │
│     │    ●───◦                             │
│ 1.3M│───◦                                  │
│     └──────────────────────────────────────│
│      Jun  Jul  Aug  Sep  Oct  Nov  Dec    │
└────────────────────────────────────────────┘
Last updated: 2025-12-15 14:30
```

#### 3.4 SnapshotManagementWidget（新規）
```typescript
// /widgets/assets/snapshot-management/ui/SnapshotManagementWidget.tsx

interface Props {
  latestSnapshot?: AssetSnapshot;
  onCreateSnapshot: () => void;
  isCreating: boolean;
}

export const SnapshotManagementWidget: React.FC<Props>
```

**責任範囲:**
- 最新スナップショット情報表示
- スナップショット手動作成ボタン
- 自動作成スケジュール表示（「毎日0時に自動作成」）
- スナップショット履歴へのリンク

**UI:**
```
┌──────────────────────────────────────────┐
│  Asset Snapshot Management               │
│                                          │
│  Latest Snapshot                         │
│  Date: 2025-01-08 00:00                 │
│  Total Assets: ¥1,480,000               │
│                                          │
│  [📸 Create Snapshot Now]               │
│                                          │
│  ℹ️ Automatic snapshots created daily    │
│     at 00:00 JST                        │
└──────────────────────────────────────────┘
```

#### 3.5 AssetForecastWidget（将来実装）
```typescript
// /widgets/assets/asset-forecast/ui/AssetForecastWidget.tsx

interface Props {
  currentAssets: number;
  forecastPeriod: '1y' | '3y' | '5y';
}

export const AssetForecastWidget: React.FC<Props>
```

**責任範囲（将来実装）:**
- 予測期間選択
- 予測シナリオ選択（楽観的・現実的・悲観的）
- 予測結果表示
- 予測グラフ表示

**実装状況:** バックエンドAPI未実装のため保留

### 4. エンティティ層

#### 4.1 Assetエンティティ（新規作成）

**model/asset.types.ts（新規）**
```typescript
// /entities/asset/model/asset.types.ts

export interface AssetSnapshot {
  id: string;
  user_id: string;
  snapshot_date: string; // ISO 8601
  total_assets: number;
  accounts: AccountBalance[];
  created_at: string;
}

export interface AccountBalance {
  account_id: string;
  account_name: string;
  balance: number;
  currency: string;
}

export interface AssetSnapshotListParams {
  from?: string; // YYYY-MM-DD
  to?: string;   // YYYY-MM-DD
}

export interface AssetSnapshotListResponse {
  snapshots: AssetSnapshot[];
  total_count: number;
}

export interface AssetTrendData {
  date: string;
  totalAssets: number;
  accounts: {
    [accountId: string]: number;
  };
}

export interface AssetComparison {
  current: number;
  previous: number;
  change: number;
  percentageChange: number;
}
```

**api/asset.client.ts（新規）**
```typescript
// /entities/asset/api/asset.client.ts

import { apiClient } from '@/shared/api';
import type {
  AssetSnapshot,
  AssetSnapshotListParams,
  AssetSnapshotListResponse,
} from '../model';

export const assetApi = {
  // スナップショット一覧取得
  getSnapshots: async (params: AssetSnapshotListParams) => {
    const { data } = await apiClient.get<AssetSnapshotListResponse>(
      '/reports/assets/snapshots',
      { params }
    );
    return data;
  },

  // 最新スナップショット取得
  getLatestSnapshot: async () => {
    const { data } = await apiClient.get<AssetSnapshot>(
      '/reports/assets/snapshots/latest'
    );
    return data;
  },

  // 現在の資産状況計算
  getCurrentStatus: async () => {
    const { data } = await apiClient.get<AssetSnapshot>(
      '/reports/assets/snapshots/current'
    );
    return data;
  },

  // スナップショット作成
  createSnapshot: async (snapshotDate: string) => {
    const { data } = await apiClient.post<AssetSnapshot>(
      '/reports/assets/snapshots',
      { snapshot_date: snapshotDate }
    );
    return data;
  },
};
```

**lib/asset.calculations.ts（新規）**
```typescript
// /entities/asset/lib/asset.calculations.ts

import type { AssetSnapshot, AssetComparison, AssetTrendData } from '../model';

// 前月比計算
export const calculateAssetComparison = (
  current: AssetSnapshot | undefined,
  previous: AssetSnapshot | undefined
): AssetComparison => {
  if (!current || !previous) {
    return {
      current: current?.total_assets || 0,
      previous: 0,
      change: 0,
      percentageChange: 0,
    };
  }

  const change = current.total_assets - previous.total_assets;
  const percentageChange = (change / previous.total_assets) * 100;

  return {
    current: current.total_assets,
    previous: previous.total_assets,
    change,
    percentageChange,
  };
};

// トレンドデータ変換
export const transformToTrendData = (
  snapshots: AssetSnapshot[]
): AssetTrendData[] => {
  return snapshots.map((snapshot) => ({
    date: snapshot.snapshot_date,
    totalAssets: snapshot.total_assets,
    accounts: snapshot.accounts.reduce((acc, account) => {
      acc[account.account_id] = account.balance;
      return acc;
    }, {} as { [key: string]: number }),
  }));
};

// 口座種別ごとの集計
export const aggregateByAccountType = (
  accounts: Account[]
): Record<string, number> => {
  return accounts.reduce((acc, account) => {
    const type = account.account_type;
    acc[type] = (acc[type] || 0) + account.current_balance;
    return acc;
  }, {} as Record<string, number>);
};
```

### 5. 機能層（Features）

#### 5.1 Asset Management Features（新規作成）

**useAssetSnapshots.ts（新規）**
```typescript
// /features/asset-management/useAssetSnapshots.ts

import { useQuery } from '@tanstack/react-query';
import { assetApi } from '@/entities/asset';
import type { AssetSnapshotListParams } from '@/entities/asset';

export const useAssetSnapshots = (params: AssetSnapshotListParams) => {
  return useQuery({
    queryKey: ['assets', 'snapshots', params],
    queryFn: () => assetApi.getSnapshots(params),
    staleTime: 5 * 60 * 1000, // 5分間キャッシュ
  });
};
```

**useLatestAssetSnapshot.ts（新規）**
```typescript
// /features/asset-management/useLatestAssetSnapshot.ts

import { useQuery } from '@tanstack/react-query';
import { assetApi } from '@/entities/asset';

export const useLatestAssetSnapshot = () => {
  return useQuery({
    queryKey: ['assets', 'snapshots', 'latest'],
    queryFn: () => assetApi.getLatestSnapshot(),
    staleTime: 5 * 60 * 1000,
  });
};
```

**useCurrentAssetStatus.ts（新規）**
```typescript
// /features/asset-management/useCurrentAssetStatus.ts

import { useQuery } from '@tanstack/react-query';
import { assetApi } from '@/entities/asset';

export const useCurrentAssetStatus = () => {
  return useQuery({
    queryKey: ['assets', 'current'],
    queryFn: () => assetApi.getCurrentStatus(),
    staleTime: 1 * 60 * 1000, // 1分間キャッシュ（リアルタイム性重視）
  });
};
```

**useCreateAssetSnapshot.ts（新規）**
```typescript
// /features/asset-management/useCreateAssetSnapshot.ts

import { useMutation, useQueryClient } from '@tanstack/react-query';
import { assetApi } from '@/entities/asset';

export const useCreateAssetSnapshot = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (snapshotDate: string) =>
      assetApi.createSnapshot(snapshotDate),
    onSuccess: () => {
      // キャッシュ無効化
      queryClient.invalidateQueries({ queryKey: ['assets', 'snapshots'] });
      queryClient.invalidateQueries({ queryKey: ['assets', 'current'] });
    },
  });
};
```

**useAssetComparison.ts（新規）**
```typescript
// /features/asset-management/useAssetComparison.ts

import { useMemo } from 'react';
import { calculateAssetComparison } from '@/entities/asset';
import type { AssetSnapshot } from '@/entities/asset';

export const useAssetComparison = (
  current: AssetSnapshot | undefined,
  previous: AssetSnapshot | undefined
) => {
  return useMemo(
    () => calculateAssetComparison(current, previous),
    [current, previous]
  );
};
```

**useAssetTrendData.ts（新規）**
```typescript
// /features/asset-management/useAssetTrendData.ts

import { useMemo } from 'react';
import { transformToTrendData } from '@/entities/asset';
import type { AssetSnapshot } from '@/entities/asset';

export const useAssetTrendData = (
  snapshots: AssetSnapshot[] | undefined,
  period: string
) => {
  return useMemo(() => {
    if (!snapshots) return [];

    // 期間に応じたフィルタリング
    const filtered = filterByPeriod(snapshots, period);

    return transformToTrendData(filtered);
  }, [snapshots, period]);
};

function filterByPeriod(
  snapshots: AssetSnapshot[],
  period: string
): AssetSnapshot[] {
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
    case 'all':
      return snapshots;
  }

  return snapshots.filter(
    (s) => new Date(s.snapshot_date) >= cutoff
  );
}
```

## データフローアーキテクチャ

### 1. 状態管理フロー:
```
AssetsContainer
    ↓
並列データフェッチ
    ├── useAccounts → 口座一覧
    ├── useLatestAssetSnapshot → 最新スナップショット
    ├── useCurrentAssetStatus → 現在の資産状況
    └── useAssetSnapshots → 期間別スナップショット一覧
         ↓
    APIクライアント → バックエンド
         ↓
    計算・変換フック
    ├── useAccountAggregates → 総資産計算
    ├── useAssetComparison → 前月比計算
    └── useAssetTrendData → グラフデータ変換
         ↓
    ウィジェットコンポーネント
```

### 2. コンポーネント階層:
```
AssetsContainer
├── AssetSummaryCard
├── AssetBreakdownList
│   ├── AccountTableRow (x N)
│   └── CreateAccountButton
├── AssetTrendChart
├── SnapshotManagementWidget
└── Modals
    ├── CreateAccountModal
    ├── EditAccountModal
    └── DeleteAccountModal
```

### 3. データ依存関係:
```
AssetSummaryCard
    ← useCurrentAssetStatus (総資産)
    ← useLatestAssetSnapshot (前月データ)
    ← useAssetComparison (前月比計算)

AssetBreakdownList
    ← useAccounts (口座一覧)
    ← useAccountAggregates (構成比計算)

AssetTrendChart
    ← useAssetSnapshots (期間別データ)
    ← useAssetTrendData (グラフ用変換)

SnapshotManagementWidget
    ← useLatestAssetSnapshot (最新情報)
    ← useCreateAssetSnapshot (作成処理)
```

## 実装ステップ

### Phase 1: エンティティ層作成（2日）
1. ✅ `/entities/asset/` ディレクトリ作成
2. ✅ 型定義作成（`asset.types.ts`）
3. ✅ APIクライアント作成（`asset.client.ts`）
4. ✅ エンドポイント定義（`asset.endpoints.ts`）
5. ✅ React Queryキー定義（`asset.keys.ts`）
6. ✅ 計算ロジック作成（`asset.calculations.ts`）
7. ✅ データ変換ロジック作成（`asset.transformers.ts`）
8. ✅ UIコンポーネント作成（バッジ、インジケーター）

### Phase 2: 機能層作成（1日）
1. ✅ `/features/asset-management/` ディレクトリ作成
2. ✅ `useAssetSnapshots` フック作成
3. ✅ `useLatestAssetSnapshot` フック作成
4. ✅ `useCurrentAssetStatus` フック作成
5. ✅ `useCreateAssetSnapshot` フック作成
6. ✅ `useAssetComparison` フック作成
7. ✅ `useAssetTrendData` フック作成

### Phase 3: ウィジェット作成（3日）
1. ✅ `AssetSummaryCard` 実装 - 0.5日
2. ✅ `AssetBreakdownList` 実装 - 1日
   - テーブル構造
   - ソート機能
   - グルーピング機能
3. ✅ `AssetTrendChart` 実装 - 1日
   - Rechartsによるグラフ実装
   - 期間切り替え
   - レスポンシブ対応
4. ✅ `SnapshotManagementWidget` 実装 - 0.5日

### Phase 4: ページ統合（1日）
1. ✅ `/app/assets/page.tsx` 作成
2. ✅ `AssetsContainer` 作成
3. ✅ ウィジェット統合
4. ✅ 既存モーダルの流用（口座CRUD）
5. ✅ ローディング・エラーハンドリング

### Phase 5: テストと最適化（1日）
1. ✅ ユニットテスト作成
2. ✅ E2Eテスト作成
3. ✅ パフォーマンス最適化
4. ✅ レスポンシブ確認
5. ✅ アクセシビリティチェック

**総工数見積もり: 8日**

## ファイル数サマリー

| 層 | ディレクトリ | 新規ファイル数 | 流用ファイル数 |
|-------|-----------|------------|------------|
| App | `/app/assets/` | 1 | 0 |
| ページコンポーネント | `/page-components/assets/` | 1 | 0 |
| ウィジェット | `/widgets/assets/` | 4 | 0 |
| ウィジェット | `/widgets/account/` | 0 | 3 |
| エンティティ | `/entities/asset/` | 13 | 0 |
| 機能 | `/features/asset-management/` | 7 | 0 |
| 機能 | `/features/account-management/` | 0 | 5 |
| **合計** | | **26** | **8** |

## API対応表

| 機能 | バックエンドAPI | 実装状況 |
|------|--------------|---------|
| 総資産表示 | `GET /accounts` | ✅ 実装済み |
| 口座一覧 | `GET /accounts` | ✅ 実装済み |
| 口座作成 | `POST /accounts` | ✅ 実装済み |
| 口座編集 | `PUT /accounts/:id` | ✅ 実装済み |
| 口座削除 | `DELETE /accounts/:id` | ✅ 実装済み |
| スナップショット一覧 | `GET /reports/assets/snapshots` | ✅ 実装済み |
| 最新スナップショット | `GET /reports/assets/snapshots/latest` | ✅ 実装済み |
| 現在の資産状況 | `GET /reports/assets/snapshots/current` | ✅ 実装済み |
| スナップショット作成 | `POST /reports/assets/snapshots` | ✅ 実装済み |
| 資産予測 | `POST /assets/forecasts` | ❌ 未実装（将来対応） |

**結論: 資産予測以外のすべてのAPIが実装済み！**

## このアーキテクチャの利点

1. **既存コードの再利用**: 口座管理の既存ウィジェットを流用
2. **段階的実装**: 各ウィジェットを独立して実装可能
3. **データ一貫性**: React Queryキャッシュによる整合性確保
4. **パフォーマンス**: 並列データフェッチと効率的なキャッシング
5. **拡張性**: 将来の資産予測機能追加が容易
6. **型安全性**: TypeScriptによる完全な型付け

## 実装のベストプラクティス

### 1. データキャッシング戦略:
```typescript
// 現在の資産状況: 1分キャッシュ（リアルタイム性重視）
useCurrentAssetStatus(); // staleTime: 1min

// スナップショット: 5分キャッシュ（安定性重視）
useAssetSnapshots(); // staleTime: 5min
```

### 2. エラーハンドリング:
```typescript
const AssetsContainer = () => {
  const snapshots = useAssetSnapshots(params);
  const current = useCurrentAssetStatus();

  if (snapshots.error || current.error) {
    return <ErrorDisplay />;
  }

  if (snapshots.isLoading || current.isLoading) {
    return <Loading />;
  }

  return <AssetsContent />;
};
```

### 3. 前月比計算の最適化:
```typescript
// メモ化で計算コストを削減
const comparison = useAssetComparison(current, latest);
```

### 4. グラフデータの変換:
```typescript
// 大量データの変換はメモ化
const trendData = useAssetTrendData(snapshots, period);
```

## まとめ

資産管理ページは**完全新規作成**が必要ですが、**既存のインフラとコンポーネントを活用**することで効率的に実装できます。

**主な利点:**
- ✅ バックエンドAPIは完全実装済み
- ✅ 口座管理の既存ウィジェットを流用可能
- ✅ Feature-Sliced Designに準拠
- ✅ 段階的な実装が可能

**工数見積もり:**
- 最短: 7日（経験豊富な開発者 + Recharts経験者）
- 標準: 8日（通常の開発者）
- 最長: 10日（初めてのプロジェクト参加者）

この実装により、ユーザーに**包括的な資産管理機能**と**視覚的にわかりやすい資産推移**を提供する、完全に機能する資産管理ページが完成します。
