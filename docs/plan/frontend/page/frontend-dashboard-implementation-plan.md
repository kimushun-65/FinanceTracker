# ダッシュボードページ実装計画

## 概要
この資料では、FinanceTrackerフロントエンドアプリケーションにおけるダッシュボードページの完全なAPI統合と実装アーキテクチャについて説明します。

## 現状と課題

### 現在の実装状況
- ✅ UI構造は完成（モックデータで表示）
- ✅ レイアウトコンポーネント実装済み
- ❌ バックエンドAPIとの接続が未実装
- ❌ リアルタイムデータ表示が未対応

### 必要な作業
1. 既存のFeatureフックを使用したAPI接続
2. モックデータの削除と実データへの置き換え
3. ローディング状態とエラーハンドリングの実装
4. リアルタイムデータ更新の実装

## 完全なディレクトリ構成

```
frontend/src/
├── app/
│   └── dashboard/
│       └── page.tsx                          # エントリーポイント
├── page-components/
│   └── dashboard/
│       └── ui/
│           └── DashboardContainer.tsx        # メインオーケストレーター（要改修）
├── widgets/
│   ├── dashboard/
│   │   ├── financial-summary/
│   │   │   └── ui/
│   │   │       └── FinancialSummaryCards.tsx # 収支サマリーカード（新規）
│   │   ├── recent-transactions/
│   │   │   └── ui/
│   │   │       └── RecentTransactionsList.tsx # 最近の取引リスト（新規）
│   │   ├── budget-progress/
│   │   │   └── ui/
│   │   │       └── BudgetProgressWidget.tsx  # 予算進捗ウィジェット（新規）
│   │   ├── income-expense-chart/
│   │   │   └── ui/
│   │   │       └── IncomeExpenseChart.tsx    # 収支トレンドグラフ（新規）
│   │   └── category-breakdown/
│   │       └── ui/
│   │           └── CategoryBreakdownChart.tsx # カテゴリ別内訳（新規）
│   └── layout/
│       ├── AppLayout.tsx                     # アプリレイアウト
│       ├── Header.tsx                        # ヘッダー
│       └── Sidebar.tsx                       # サイドバー
├── entities/
│   ├── transaction/
│   │   ├── model/
│   │   │   ├── transaction.types.ts          # 型定義（既存）
│   │   │   └── index.ts
│   │   ├── api/
│   │   │   ├── transaction.client.ts         # APIクライアント（既存）
│   │   │   └── index.ts
│   │   └── lib/
│   │       ├── transaction.aggregations.ts   # 集計処理（既存）
│   │       └── index.ts
│   ├── account/
│   │   ├── model/
│   │   │   ├── account.types.ts              # 型定義（既存）
│   │   │   └── index.ts
│   │   ├── api/
│   │   │   ├── account.client.ts             # APIクライアント（既存）
│   │   │   └── index.ts
│   │   └── index.ts
│   ├── budget/
│   │   ├── model/
│   │   │   ├── budget.types.ts               # 型定義（既存）
│   │   │   └── index.ts
│   │   ├── api/
│   │   │   ├── budget.client.ts              # APIクライアント（既存）
│   │   │   └── index.ts
│   │   └── lib/
│   │       ├── budget.calculations.ts        # 計算処理（既存）
│   │       └── index.ts
│   └── category/
│       ├── model/
│       │   ├── category.types.ts             # 型定義（既存）
│       │   └── index.ts
│       └── index.ts
├── features/
│   ├── dashboard-management/                 # 新規ディレクトリ
│   │   ├── useDashboardSummary.ts            # ダッシュボードサマリー（新規）
│   │   └── index.ts
│   ├── transaction-management/               # 既存
│   │   ├── useTransactions.ts                # 取引一覧（既存）
│   │   ├── useTransactionMonthlySummary.ts   # 月次サマリー（既存）
│   │   ├── useTransactionCategorySummary.ts  # カテゴリ別サマリー（新規）
│   │   └── index.ts
│   ├── account-management/                   # 既存
│   │   ├── useAccounts.ts                    # 口座一覧（既存）
│   │   ├── useAccountAggregates.ts           # 口座集計（既存）
│   │   └── index.ts
│   └── budget-management/                    # 既存
│       ├── useCurrentBudgets.ts              # 現在の予算（既存）
│       ├── useBudgetSummary.ts               # 予算サマリー（既存）
│       └── index.ts
└── shared/
    ├── ui/
    │   ├── card.tsx                          # カードコンポーネント（既存）
    │   ├── loading.tsx                       # ローディング（既存）
    │   └── index.ts
    └── lib/
        └── chart/                            # 新規ディレクトリ
            ├── chartConfig.ts                # チャート設定（新規）
            └── index.ts
```

## アーキテクチャ階層

### 1. ページ層
**エントリーポイント:**
```
/app/dashboard/page.tsx
```
- メインルートのエントリーポイント
- `DashboardContainer`コンポーネントをレンダリング
- ProtectedRouteでラップ

### 2. ページコンポーネント層
**メインコンテナ:**
```
/page-components/dashboard/ui/DashboardContainer.tsx
```
**現在の状態:**
- モックデータを直接記述
- Auth同期のみ実装済み

**必要な改修:**
1. モックデータの削除
2. 以下のFeatureフックを統合:
   - `useAccountAggregates` - 総資産計算
   - `useTransactionMonthlySummary` - 月次サマリー
   - `useCurrentBudgets` + `useBudgetSummary` - 予算情報
   - `useTransactions` - 最近の取引（limit=5）
   - `useTransactionCategorySummary` - カテゴリ別支出
3. ローディング状態の統合
4. エラーハンドリングの追加
5. 新規ウィジェットコンポーネントへの切り替え

**主要な依存関係:**
- `@/widgets/dashboard/*` - 新規作成するダッシュボードウィジェット
- `@/shared/ui` - UIコンポーネント
- `@/features` - データ管理用フック
- `@/entities` - 型定義

### 3. ウィジェット層
**ダッシュボードウィジェットディレクトリ: `/widgets/dashboard/`** (新規作成)

#### 核となるウィジェットコンポーネント:

#### 3.1 FinancialSummaryCards（新規）
```typescript
// /widgets/dashboard/financial-summary/ui/FinancialSummaryCards.tsx

interface Props {
  totalAssets: number;
  monthlyIncome: number;
  monthlyExpenses: number;
  netIncome: number;
  budgetRemaining: number;
  budgetProgress: number;
}

export const FinancialSummaryCards: React.FC<Props>
```

**責任範囲:**
- 4つのサマリーカードを表示
  1. Total Balance（総資産）
  2. Monthly Expenses（月次支出）
  3. Budget Remaining（予算残高）
  4. Savings Goal（貯蓄目標）※将来実装
- 前月比の表示
- 色分けされた増減表示

**使用するAPI:**
- `useAccounts` → 総資産
- `useTransactionMonthlySummary` → 月次収支
- `useBudgetSummary` → 予算残高

#### 3.2 RecentTransactionsList（新規）
```typescript
// /widgets/dashboard/recent-transactions/ui/RecentTransactionsList.tsx

interface Props {
  transactions: Transaction[];
  categories: Category[];
  isLoading: boolean;
}

export const RecentTransactionsList: React.FC<Props>
```

**責任範囲:**
- 最新5件の取引を表示
- カテゴリ名とアイコン表示
- 金額の色分け（収入=緑、支出=赤）
- 「すべて見る」リンク

**使用するAPI:**
- `useTransactions({ limit: 5, sort: 'date:desc' })`
- `useCategories`

#### 3.3 BudgetProgressWidget（新規）
```typescript
// /widgets/dashboard/budget-progress/ui/BudgetProgressWidget.tsx

interface Props {
  budgets: Budget[];
  categories: Category[];
}

export const BudgetProgressWidget: React.FC<Props>
```

**責任範囲:**
- 主要カテゴリの予算進捗を3つ表示
- プログレスバー表示
- 使用金額 / 予算額の表示
- 進捗率に応じた色変更（緑→黄→赤）

**使用するAPI:**
- `useCurrentBudgets`
- `useCategories`

#### 3.4 IncomeExpenseChart（新規）
```typescript
// /widgets/dashboard/income-expense-chart/ui/IncomeExpenseChart.tsx

interface Props {
  monthlySummary: MonthlySummary[];
  period: '1m' | '6m' | '1y';
  onPeriodChange: (period: string) => void;
}

export const IncomeExpenseChart: React.FC<Props>
```

**責任範囲:**
- 収支トレンドの棒グラフ表示
- 期間切り替え（1ヶ月、6ヶ月、1年）
- 収入と支出の2系列グラフ
- レスポンシブ対応

**使用するAPI:**
- `useTransactionMonthlySummary({ period })`

**使用ライブラリ:**
- Recharts または Chart.js

#### 3.5 CategoryBreakdownChart（新規）
```typescript
// /widgets/dashboard/category-breakdown/ui/CategoryBreakdownChart.tsx

interface Props {
  categorySummary: CategorySummary[];
}

export const CategoryBreakdownChart: React.FC<Props>
```

**責任範囲:**
- カテゴリ別支出のドーナツチャート
- カテゴリごとの色分け
- パーセンテージ表示
- ホバー時の詳細表示

**使用するAPI:**
- `useTransactionCategorySummary({ type: 'expense', period: '1m' })`

### 4. エンティティ層
**既存エンティティを活用:**
- `/entities/transaction/` - 取引関連の型とロジック
- `/entities/account/` - 口座関連の型とロジック
- `/entities/budget/` - 予算関連の型とロジック
- `/entities/category/` - カテゴリ関連の型とロジック

**新規追加不要** - 既存エンティティで対応可能

### 5. 機能層（Features）

#### 5.1 既存Featureフックの活用
```
/features/transaction-management/
├── useTransactions.ts                  # 取引一覧（既存）
├── useTransactionMonthlySummary.ts     # 月次サマリー（既存）
└── useTransactionCategorySummary.ts    # カテゴリ別サマリー（新規作成）

/features/account-management/
├── useAccounts.ts                      # 口座一覧（既存）
└── useAccountAggregates.ts             # 総資産計算（既存）

/features/budget-management/
├── useCurrentBudgets.ts                # 現在の予算（既存）
└── useBudgetSummary.ts                 # 予算サマリー（既存）

/features/category-management/
└── useCategories.ts                    # カテゴリ一覧（既存）
```

#### 5.2 新規Featureフック

**useTransactionCategorySummary.ts（新規作成）**
```typescript
// /features/transaction-management/useTransactionCategorySummary.ts

export const useTransactionCategorySummary = (params?: {
  type?: 'income' | 'expense';
  from?: string;
  to?: string;
}) => {
  return useQuery({
    queryKey: transactionKeys.categorySummary(params),
    queryFn: () => transactionApi.getCategorySummary(params),
  });
};
```

**バックエンドAPI:**
- `GET /api/v1/transactions/summary/by-category?type=expense&from=2024-01-01&to=2024-01-31`

**useDashboardSummary.ts（新規作成・オプション）**
```typescript
// /features/dashboard-management/useDashboardSummary.ts

/**
 * ダッシュボード用の複数データを一度に取得
 * パフォーマンス最適化のための統合フック
 */
export const useDashboardSummary = () => {
  const accounts = useAccounts();
  const { totalAssets } = useAccountAggregates(accounts.data);
  const monthlySummary = useTransactionMonthlySummary();
  const budgetSummary = useBudgetSummary('monthly');
  const recentTransactions = useTransactions({ limit: 5, sort: 'date:desc' });
  const categorySummary = useTransactionCategorySummary({ type: 'expense' });

  return {
    totalAssets,
    monthlySummary: monthlySummary.data,
    budgetSummary: budgetSummary.data,
    recentTransactions: recentTransactions.data,
    categorySummary: categorySummary.data,
    isLoading: accounts.isLoading ||
               monthlySummary.isLoading ||
               budgetSummary.isLoading,
  };
};
```

## データフローアーキテクチャ

### 1. 状態管理フロー:
```
DashboardContainer
    ↓
複数のFeatureフック（並列実行）
    ├── useAccounts → 総資産計算
    ├── useTransactionMonthlySummary → 月次サマリー
    ├── useCurrentBudgets → 予算一覧
    ├── useBudgetSummary → 予算サマリー
    ├── useTransactions → 最近の取引
    └── useTransactionCategorySummary → カテゴリ別集計
         ↓
    APIクライアント → バックエンド
         ↓
    React Query キャッシュ
         ↓
    ウィジェットコンポーネント
```

### 2. コンポーネント階層:
```
DashboardContainer
├── FinancialSummaryCards
│   ├── TotalBalanceCard
│   ├── MonthlyExpensesCard
│   ├── BudgetRemainingCard
│   └── SavingsGoalCard
├── IncomeExpenseChart
├── CategoryBreakdownChart
├── RecentTransactionsList
└── BudgetProgressWidget
```

### 3. データ依存関係:
```
FinancialSummaryCards
    ← useAccountAggregates (総資産)
    ← useTransactionMonthlySummary (月次収支)
    ← useBudgetSummary (予算残高)

IncomeExpenseChart
    ← useTransactionMonthlySummary (期間指定)

CategoryBreakdownChart
    ← useTransactionCategorySummary (カテゴリ別支出)

RecentTransactionsList
    ← useTransactions (limit=5)
    ← useCategories (カテゴリ名表示)

BudgetProgressWidget
    ← useCurrentBudgets (現在の予算)
    ← useCategories (カテゴリ名表示)
```

## 実装ステップ

### Phase 1: 新規Featureフック作成（1日）
1. ✅ `useTransactionCategorySummary.ts` 作成
2. ✅ バックエンドAPIとの接続確認
3. ✅ 型定義の追加（必要に応じて）
4. ✅ React Queryキー定義

### Phase 2: ウィジェットコンポーネント作成（3-4日）
1. ✅ `FinancialSummaryCards.tsx` - 1日
   - 4つのカードコンポーネント
   - 前月比表示ロジック
2. ✅ `RecentTransactionsList.tsx` - 0.5日
   - 取引リスト表示
   - カテゴリアイコン統合
3. ✅ `BudgetProgressWidget.tsx` - 0.5日
   - プログレスバー実装
   - 進捗率計算
4. ✅ `IncomeExpenseChart.tsx` - 1日
   - グラフライブラリ選定（Recharts推奨）
   - 棒グラフ実装
   - 期間切り替え機能
5. ✅ `CategoryBreakdownChart.tsx` - 1日
   - ドーナツチャート実装
   - カテゴリ色分け

### Phase 3: DashboardContainer改修（1日）
1. ✅ モックデータ削除
2. ✅ Featureフック統合
3. ✅ 新規ウィジェットに置き換え
4. ✅ ローディング状態の実装
5. ✅ エラーハンドリング追加

### Phase 4: テストと最適化（1日）
1. ✅ ユニットテスト作成
2. ✅ パフォーマンス最適化
3. ✅ レスポンシブ対応確認
4. ✅ アクセシビリティチェック

**総工数見積もり: 6-7日**

## 主要な設計原則

### 1. Feature-Sliced Design (FSD):
- **App**: ルートレベルコンポーネント
- **Pages**: ページオーケストレーションコンポーネント
- **Widgets**: 複雑なUIウィジェット（新規作成）
- **Features**: ビジネスロジックとデータ管理（既存活用+一部新規）
- **Entities**: ドメインモデル（既存活用）
- **Shared**: 再利用可能なユーティリティ

### 2. パフォーマンス最適化:
- React Query によるデータキャッシング
- 複数APIの並列フェッチ
- メモ化による不要な再レンダリング防止
- 遅延ローディングの活用

### 3. データ整合性:
- すべてのウィジェットが同じデータソースを参照
- React Query のキャッシュ戦略による一貫性確保
- リアルタイム更新の自動反映

## ファイル数サマリー

| 層 | ディレクトリ | 新規ファイル数 | 改修ファイル数 |
|-------|-----------|------------|------------|
| App | `/app/dashboard/` | 0 | 0 |
| ページコンポーネント | `/page-components/dashboard/` | 0 | 1 |
| ウィジェット | `/widgets/dashboard/` | 5 | 0 |
| 機能 | `/features/transaction-management/` | 1 | 0 |
| 機能 | `/features/dashboard-management/` | 1 | 0 |
| 共有 | `/shared/lib/chart/` | 1 | 0 |
| **合計** | | **8** | **1** |

## API対応表

| ウィジェット | バックエンドAPI | 実装状況 |
|-----------|--------------|---------|
| FinancialSummaryCards | `GET /accounts` | ✅ 実装済み |
| | `GET /transactions/summary/monthly` | ✅ 実装済み |
| | `GET /budgets/summary` | ✅ 実装済み |
| RecentTransactionsList | `GET /transactions?limit=5` | ✅ 実装済み |
| | `GET /categories` | ✅ 実装済み |
| BudgetProgressWidget | `GET /budgets/current` | ✅ 実装済み |
| | `GET /categories` | ✅ 実装済み |
| IncomeExpenseChart | `GET /transactions/summary/monthly` | ✅ 実装済み |
| CategoryBreakdownChart | `GET /transactions/summary/by-category` | ✅ 実装済み |

**結論: すべての必要なAPIが実装済み！**

## このアーキテクチャの利点

1. **既存コードの再利用**: 既存のFeatureフックとエンティティを最大活用
2. **段階的実装**: ウィジェット単位で独立して実装可能
3. **パフォーマンス**: React Queryによる効率的なデータ取得
4. **保守性**: 各ウィジェットが独立しており、変更が局所化される
5. **テスト性**: 各ウィジェットを独立してテスト可能
6. **拡張性**: 新しいウィジェットの追加が容易

## 実装のベストプラクティス

### 1. データフェッチング:
```typescript
// ✅ 良い例: 並列フェッチ
const DashboardContainer = () => {
  const { data: accounts } = useAccounts();
  const { data: summary } = useTransactionMonthlySummary();
  const { data: budgets } = useCurrentBudgets();
  // すべて並列で実行される
};

// ❌ 悪い例: 直列フェッチ
const DashboardContainer = () => {
  const { data: accounts } = useAccounts();
  if (accounts) {
    const { data: summary } = useTransactionMonthlySummary();
    // accountsの取得完了を待ってしまう
  }
};
```

### 2. エラーハンドリング:
```typescript
const DashboardContainer = () => {
  const accounts = useAccounts();
  const summary = useTransactionMonthlySummary();

  // 個別のエラーチェック
  if (accounts.error) return <ErrorDisplay error={accounts.error} />;
  if (summary.error) return <ErrorDisplay error={summary.error} />;

  // ローディング状態
  if (accounts.isLoading || summary.isLoading) return <Loading />;

  return <DashboardContent />;
};
```

### 3. メモ化:
```typescript
const FinancialSummaryCards = ({ accounts, summary, budgets }) => {
  // 計算コストの高い処理はメモ化
  const totalAssets = useMemo(
    () => calculateTotalAssets(accounts),
    [accounts]
  );

  return <div>{totalAssets}</div>;
};
```

### 4. チャートのパフォーマンス:
```typescript
const IncomeExpenseChart = ({ data }) => {
  // データ変換をメモ化
  const chartData = useMemo(
    () => transformDataForChart(data),
    [data]
  );

  return <ResponsiveContainer>
    <BarChart data={chartData}>
      {/* ... */}
    </BarChart>
  </ResponsiveContainer>;
};
```

## 推奨ライブラリ

### グラフ表示:
- **Recharts** (推奨)
  - React専用の宣言的なチャートライブラリ
  - レスポンシブ対応
  - TypeScript完全サポート
  - ```bash
    npm install recharts
    ```

### 代替案:
- **Chart.js + react-chartjs-2**
  - より軽量
  - カスタマイズ性が高い

## 実装チェックリスト

### Phase 1: 準備
- [ ] Rechartsライブラリのインストール
- [ ] `useTransactionCategorySummary` フック作成
- [ ] 型定義の追加・確認

### Phase 2: ウィジェット作成
- [ ] `FinancialSummaryCards` 実装
- [ ] `RecentTransactionsList` 実装
- [ ] `BudgetProgressWidget` 実装
- [ ] `IncomeExpenseChart` 実装
- [ ] `CategoryBreakdownChart` 実装

### Phase 3: 統合
- [ ] `DashboardContainer` 改修
- [ ] モックデータ削除
- [ ] API接続確認
- [ ] ローディング状態実装
- [ ] エラーハンドリング実装

### Phase 4: 仕上げ
- [ ] レスポンシブデザイン確認
- [ ] アクセシビリティチェック
- [ ] パフォーマンス最適化
- [ ] ユニットテスト作成
- [ ] E2Eテスト作成

## まとめ

ダッシュボードページは**既存のインフラを最大限活用**することで、効率的に実装できます。

**主な利点:**
- ✅ バックエンドAPIは完全実装済み
- ✅ 既存のFeatureフックをそのまま活用可能
- ✅ UI構造は既に完成
- ✅ 必要なのはAPI接続と新規ウィジェットのみ

**工数見積もり:**
- 最短: 5-6日（経験豊富な開発者）
- 標準: 6-7日（通常の開発者）
- 最長: 8-9日（初めてのプロジェクト参加者）

この実装により、ユーザーに**リアルタイムで正確な財務情報**を提供する、完全に機能するダッシュボードが完成します。
