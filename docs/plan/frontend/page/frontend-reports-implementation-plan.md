# レポートページ実装計画

## 概要
この資料では、FinanceTrackerフロントエンドアプリケーションにおけるレポートページの完全な実装アーキテクチャについて説明します。このページは**完全に新規作成**が必要です。

## ページ要件（要件定義より）

### 主要機能
1. **月次レポート**: 月別の収支・カテゴリ別サマリー
2. **期間比較**: 月次・四半期・年次の比較
3. **カテゴリ別分析**: 詳細なカテゴリ別支出分析
4. **トレンド分析**: 長期的な支出・収入トレンド
5. **レポート出力**: PDF/CSVエクスポート（将来実装予定）

### 対応API
- ✅ `GET /api/v1/transactions/summary/monthly` - 月次サマリー
- ✅ `GET /api/v1/transactions/summary/by-category` - カテゴリ別サマリー
- ✅ `GET /api/v1/transactions?from=&to=` - 期間別取引一覧
- ✅ `GET /api/v1/budgets/summary` - 予算サマリー
- ❌ `GET /api/v1/reports/summary/monthly` - 統合月次レポート（未実装）
- ❌ `POST /api/v1/reports/export` - レポート出力（未実装）

## 完全なディレクトリ構成

```
frontend/src/
├── app/
│   └── reports/                              # 新規ディレクトリ
│       └── page.tsx                          # エントリーポイント（新規）
├── page-components/
│   └── reports/                              # 新規ディレクトリ
│       └── ui/
│           └── ReportsContainer.tsx          # メインオーケストレーター（新規）
├── widgets/
│   ├── reports/                              # 新規ディレクトリ
│   │   ├── monthly-report/
│   │   │   └── ui/
│   │   │       └── MonthlyReportCard.tsx     # 月次レポートカード（新規）
│   │   ├── period-comparison/
│   │   │   └── ui/
│   │   │       └── PeriodComparisonWidget.tsx # 期間比較ウィジェット（新規）
│   │   ├── category-analysis/
│   │   │   └── ui/
│   │   │       └── CategoryAnalysisChart.tsx  # カテゴリ分析グラフ（新規）
│   │   ├── trend-analysis/
│   │   │   └── ui/
│   │   │       └── TrendAnalysisChart.tsx    # トレンド分析グラフ（新規）
│   │   ├── budget-vs-actual/
│   │   │   └── ui/
│   │   │       └── BudgetVsActualWidget.tsx  # 予算vs実績（新規）
│   │   ├── report-filters/
│   │   │   └── ui/
│   │   │       └── ReportFilters.tsx         # レポートフィルター（新規）
│   │   └── report-export/                    # 将来実装予定
│   │       └── ui/
│   │           └── ReportExportButton.tsx    # エクスポートボタン（新規）
│   └── layout/
│       ├── AppLayout.tsx                     # 既存
│       └── index.ts
├── entities/
│   ├── transaction/                          # 既存（流用）
│   │   ├── model/
│   │   │   ├── transaction.types.ts          # 既存
│   │   │   └── index.ts
│   │   ├── api/
│   │   │   ├── transaction.client.ts         # 既存
│   │   │   └── index.ts
│   │   ├── lib/
│   │   │   ├── transaction.aggregations.ts   # 既存
│   │   │   └── index.ts
│   │   └── index.ts
│   ├── category/                             # 既存（流用）
│   │   ├── model/
│   │   │   └── category.types.ts             # 既存
│   │   └── index.ts
│   ├── budget/                               # 既存（流用）
│   │   └── model/
│   │       └── budget.types.ts               # 既存
│   └── report/                               # 新規ディレクトリ
│       ├── model/
│       │   ├── report.types.ts               # 型定義（新規）
│       │   ├── report.constants.ts           # 定数（新規）
│       │   └── index.ts
│       ├── lib/
│       │   ├── report.calculations.ts        # 計算処理（新規）
│       │   ├── report.comparisons.ts         # 比較処理（新規）
│       │   ├── report.formatters.ts          # フォーマッター（新規）
│       │   └── index.ts
│       └── index.ts
├── features/
│   ├── transaction-management/               # 既存（流用）
│   │   ├── useTransactions.ts                # 既存
│   │   ├── useTransactionMonthlySummary.ts   # 既存
│   │   ├── useTransactionCategorySummary.ts  # 既存（ダッシュボードで作成）
│   │   └── index.ts
│   ├── budget-management/                    # 既存（流用）
│   │   ├── useBudgetSummary.ts               # 既存
│   │   └── index.ts
│   ├── category-management/                  # 既存（流用）
│   │   ├── useCategories.ts                  # 既存
│   │   └── index.ts
│   └── report-management/                    # 新規ディレクトリ
│       ├── useMonthlyReport.ts               # 月次レポート（新規）
│       ├── usePeriodComparison.ts            # 期間比較（新規）
│       ├── useCategoryTrend.ts               # カテゴリトレンド（新規）
│       └── index.ts
└── shared/
    ├── ui/
    │   ├── card.tsx                          # 既存
    │   ├── select.tsx                        # 既存
    │   ├── button.tsx                        # 既存
    │   └── index.ts
    ├── lib/
    │   └── chart/
    │       ├── chartConfig.ts                # 既存
    │       └── index.ts
    └── value-objects/
        └── date/
            ├── date.utils.ts                 # 既存
            └── index.ts
```

## アーキテクチャ階層

### 1. ページ層
**エントリーポイント（新規作成）:**
```typescript
// /app/reports/page.tsx

import { ReportsContainer } from '../../page-components/reports';
import { ProtectedRoute } from '../../shared/ui/components/ProtectedRoute';

export default function ReportsPage() {
  return (
    <ProtectedRoute>
      <ReportsContainer />
    </ProtectedRoute>
  );
}
```

### 2. ページコンポーネント層
**メインコンテナ（新規作成）:**
```typescript
// /page-components/reports/ui/ReportsContainer.tsx

export const ReportsContainer: React.FC = () => {
  // フィルター状態
  const [dateRange, setDateRange] = useState<DateRange>({
    from: startOfMonth(new Date()),
    to: endOfMonth(new Date()),
  });
  const [selectedCategories, setSelectedCategories] = useState<string[]>([]);
  const [comparisonPeriod, setComparisonPeriod] = useState<'month' | 'quarter' | 'year'>('month');

  // データ取得
  const { data: monthlySummary } = useTransactionMonthlySummary({
    from: dateRange.from,
    to: dateRange.to,
  });
  const { data: categorySummary } = useTransactionCategorySummary({
    from: format(dateRange.from, 'yyyy-MM-dd'),
    to: format(dateRange.to, 'yyyy-MM-dd'),
    type: 'expense',
  });
  const { data: budgetSummary } = useBudgetSummary('monthly');
  const { data: categories } = useCategories();
  const { data: periodComparison } = usePeriodComparison(dateRange, comparisonPeriod);

  // 計算処理
  const monthlyReport = useMonthlyReport(monthlySummary, categorySummary);
  const categoryTrend = useCategoryTrend(categorySummary, selectedCategories);

  return (
    <AppLayout title='Reports'>
      <div className='space-y-6'>
        {/* レポートフィルター */}
        <ReportFilters
          dateRange={dateRange}
          onDateRangeChange={setDateRange}
          selectedCategories={selectedCategories}
          onCategoriesChange={setSelectedCategories}
          categories={categories || []}
          comparisonPeriod={comparisonPeriod}
          onComparisonPeriodChange={setComparisonPeriod}
        />

        {/* 月次レポートカード */}
        <MonthlyReportCard
          report={monthlyReport}
          dateRange={dateRange}
        />

        {/* 期間比較 */}
        <PeriodComparisonWidget
          comparison={periodComparison}
          period={comparisonPeriod}
        />

        {/* 予算 vs 実績 */}
        <BudgetVsActualWidget
          budgetSummary={budgetSummary}
          actualSummary={categorySummary}
        />

        <div className='grid gap-6 md:grid-cols-2'>
          {/* カテゴリ分析グラフ */}
          <CategoryAnalysisChart
            data={categorySummary}
            categories={categories || []}
          />

          {/* トレンド分析グラフ */}
          <TrendAnalysisChart
            data={categoryTrend}
            selectedCategories={selectedCategories}
          />
        </div>

        {/* エクスポートボタン（将来実装）*/}
        {/* <ReportExportButton /> */}
      </div>
    </AppLayout>
  );
};
```

**主要な責任範囲:**
- レポートフィルター状態管理
- 複数データソースの統合
- ウィジェットのオーケストレーション
- ローディング・エラーハンドリング

### 3. ウィジェット層

#### 3.1 ReportFilters（新規）
```typescript
// /widgets/reports/report-filters/ui/ReportFilters.tsx

interface Props {
  dateRange: DateRange;
  onDateRangeChange: (range: DateRange) => void;
  selectedCategories: string[];
  onCategoriesChange: (categories: string[]) => void;
  categories: Category[];
  comparisonPeriod: 'month' | 'quarter' | 'year';
  onComparisonPeriodChange: (period: string) => void;
}

export const ReportFilters: React.FC<Props>
```

**責任範囲:**
- 日付範囲選択（From/To）
- カテゴリー複数選択
- 比較期間選択（月次・四半期・年次）
- クイック期間選択（今月、先月、今四半期、今年）
- フィルターリセットボタン

**UI:**
```
┌─────────────────────────────────────────────────────┐
│  Report Filters                                     │
│                                                     │
│  Period: [From: 2024-01-01] [To: 2024-01-31]      │
│  Quick: [This Month] [Last Month] [This Quarter]   │
│                                                     │
│  Categories: [Select Categories ▼]                 │
│  Comparison: ● Month ○ Quarter ○ Year              │
│                                                     │
│  [Reset Filters]                                   │
└─────────────────────────────────────────────────────┘
```

#### 3.2 MonthlyReportCard（新規）
```typescript
// /widgets/reports/monthly-report/ui/MonthlyReportCard.tsx

interface Props {
  report: MonthlyReport;
  dateRange: DateRange;
}

interface MonthlyReport {
  totalIncome: number;
  totalExpenses: number;
  netIncome: number;
  transactionCount: number;
  averageTransactionAmount: number;
  topCategories: CategorySummary[];
}

export const MonthlyReportCard: React.FC<Props>
```

**責任範囲:**
- 選択期間のサマリー表示
- 総収入・総支出・純収入
- 取引件数・平均取引額
- トップカテゴリー表示（支出上位5つ）

**UI:**
```
┌─────────────────────────────────────────────────────┐
│  Monthly Report: 2024年1月                          │
│                                                     │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐           │
│  │ Income  │  │ Expense │  │Net Income│          │
│  │¥350,000 │  │¥280,000 │  │+¥70,000 │          │
│  └─────────┘  └─────────┘  └─────────┘           │
│                                                     │
│  Transactions: 145 件                              │
│  Average: ¥1,931 / transaction                     │
│                                                     │
│  Top Expense Categories:                           │
│  1. 食費        ¥80,000 (28.6%)                    │
│  2. 交通費      ¥50,000 (17.9%)                    │
│  3. 娯楽        ¥45,000 (16.1%)                    │
└─────────────────────────────────────────────────────┘
```

#### 3.3 PeriodComparisonWidget（新規）
```typescript
// /widgets/reports/period-comparison/ui/PeriodComparisonWidget.tsx

interface Props {
  comparison: PeriodComparison;
  period: 'month' | 'quarter' | 'year';
}

interface PeriodComparison {
  current: PeriodData;
  previous: PeriodData;
  change: {
    income: number;
    expenses: number;
    netIncome: number;
  };
  percentageChange: {
    income: number;
    expenses: number;
    netIncome: number;
  };
}

export const PeriodComparisonWidget: React.FC<Props>
```

**責任範囲:**
- 現在期間と前期間の比較
- 増減額・増減率の表示
- 視覚的な比較表示（棒グラフ）
- トレンドアイコン（増加・減少）

**UI:**
```
┌─────────────────────────────────────────────────────┐
│  Period Comparison (Month over Month)              │
│                                                     │
│           Current      Previous      Change        │
│  Income   ¥350,000    ¥330,000     ▲+¥20k (+6.1%)│
│  Expense  ¥280,000    ¥300,000     ▼-¥20k (-6.7%)│
│  Net      ¥70,000     ¥30,000      ▲+¥40k (+133%)│
│                                                     │
│  [Comparison Chart]                                │
│  ████████████  Current                             │
│  ██████████    Previous                            │
└─────────────────────────────────────────────────────┘
```

#### 3.4 CategoryAnalysisChart（新規）
```typescript
// /widgets/reports/category-analysis/ui/CategoryAnalysisChart.tsx

interface Props {
  data: CategorySummary[];
  categories: Category[];
}

export const CategoryAnalysisChart: React.FC<Props>
```

**責任範囲:**
- カテゴリ別支出の円グラフ表示
- カテゴリごとの金額・割合表示
- ホバー時の詳細情報
- カテゴリ色分け
- レジェンド表示

**使用ライブラリ:** Recharts（PieChart）

#### 3.5 TrendAnalysisChart（新規）
```typescript
// /widgets/reports/trend-analysis/ui/TrendAnalysisChart.tsx

interface Props {
  data: TrendData[];
  selectedCategories: string[];
}

interface TrendData {
  date: string;
  [categoryId: string]: number | string;
}

export const TrendAnalysisChart: React.FC<Props>
```

**責任範囲:**
- 選択カテゴリのトレンド表示（折れ線グラフ）
- 複数カテゴリの重ね合わせ表示
- 期間別の推移可視化
- データポイントクリックで詳細表示

**使用ライブラリ:** Recharts（LineChart）

#### 3.6 BudgetVsActualWidget（新規）
```typescript
// /widgets/reports/budget-vs-actual/ui/BudgetVsActualWidget.tsx

interface Props {
  budgetSummary: BudgetSummary;
  actualSummary: CategorySummary[];
}

export const BudgetVsActualWidget: React.FC<Props>
```

**責任範囲:**
- 予算と実績の比較表示
- カテゴリ別の予算達成率
- 予算超過の警告表示
- プログレスバー表示

**UI:**
```
┌─────────────────────────────────────────────────────┐
│  Budget vs Actual                                   │
│                                                     │
│  Category     Budget    Actual    Progress         │
│  食費         ¥60,000   ¥50,000   ████████░░ 83%   │
│  交通費       ¥40,000   ¥45,000   ████████████ 113% │ ⚠️
│  娯楽         ¥30,000   ¥25,000   ████████░░ 83%   │
│                                                     │
│  Total Budget: ¥250,000                            │
│  Total Actual: ¥240,000                            │
│  Remaining: ¥10,000 (4%)                           │
└─────────────────────────────────────────────────────┘
```

### 4. エンティティ層

#### 4.1 Reportエンティティ（新規作成）

**model/report.types.ts（新規）**
```typescript
// /entities/report/model/report.types.ts

export interface MonthlyReport {
  totalIncome: number;
  totalExpenses: number;
  netIncome: number;
  transactionCount: number;
  averageTransactionAmount: number;
  topCategories: CategorySummary[];
}

export interface PeriodData {
  income: number;
  expenses: number;
  netIncome: number;
  transactionCount: number;
}

export interface PeriodComparison {
  current: PeriodData;
  previous: PeriodData;
  change: {
    income: number;
    expenses: number;
    netIncome: number;
  };
  percentageChange: {
    income: number;
    expenses: number;
    netIncome: number;
  };
}

export interface TrendData {
  date: string;
  [categoryId: string]: number | string;
}

export interface DateRange {
  from: Date;
  to: Date;
}

export type ComparisonPeriod = 'month' | 'quarter' | 'year';
```

**lib/report.calculations.ts（新規）**
```typescript
// /entities/report/lib/report.calculations.ts

import type {
  MonthlyReport,
  PeriodComparison,
  PeriodData,
} from '../model';

// 月次レポート計算
export const calculateMonthlyReport = (
  transactions: Transaction[],
  categorySummary: CategorySummary[]
): MonthlyReport => {
  const income = transactions
    .filter((t) => t.type === 'income')
    .reduce((sum, t) => sum + t.amount, 0);

  const expenses = transactions
    .filter((t) => t.type === 'expense')
    .reduce((sum, t) => sum + t.amount, 0);

  return {
    totalIncome: income,
    totalExpenses: expenses,
    netIncome: income - expenses,
    transactionCount: transactions.length,
    averageTransactionAmount: transactions.length
      ? (income + expenses) / transactions.length
      : 0,
    topCategories: categorySummary
      .sort((a, b) => b.total_amount - a.total_amount)
      .slice(0, 5),
  };
};

// 期間比較計算
export const calculatePeriodComparison = (
  current: PeriodData,
  previous: PeriodData
): PeriodComparison => {
  const change = {
    income: current.income - previous.income,
    expenses: current.expenses - previous.expenses,
    netIncome: current.netIncome - previous.netIncome,
  };

  const percentageChange = {
    income: previous.income
      ? (change.income / previous.income) * 100
      : 0,
    expenses: previous.expenses
      ? (change.expenses / previous.expenses) * 100
      : 0,
    netIncome: previous.netIncome
      ? (change.netIncome / previous.netIncome) * 100
      : 0,
  };

  return {
    current,
    previous,
    change,
    percentageChange,
  };
};
```

**lib/report.comparisons.ts（新規）**
```typescript
// /entities/report/lib/report.comparisons.ts

import { subMonths, subQuarters, subYears } from 'date-fns';
import type { DateRange, ComparisonPeriod } from '../model';

// 比較期間の計算
export const calculateComparisonPeriod = (
  currentRange: DateRange,
  period: ComparisonPeriod
): DateRange => {
  switch (period) {
    case 'month':
      return {
        from: subMonths(currentRange.from, 1),
        to: subMonths(currentRange.to, 1),
      };
    case 'quarter':
      return {
        from: subQuarters(currentRange.from, 1),
        to: subQuarters(currentRange.to, 1),
      };
    case 'year':
      return {
        from: subYears(currentRange.from, 1),
        to: subYears(currentRange.to, 1),
      };
  }
};

// トレンドデータ変換
export const transformToTrendData = (
  transactions: Transaction[],
  categories: Category[]
): TrendData[] => {
  const grouped = groupByDate(transactions);

  return Object.entries(grouped).map(([date, txs]) => {
    const data: TrendData = { date };

    categories.forEach((category) => {
      const categoryTotal = txs
        .filter((t) => t.category_id === category.id)
        .reduce((sum, t) => sum + t.amount, 0);

      data[category.id] = categoryTotal;
    });

    return data;
  });
};

function groupByDate(
  transactions: Transaction[]
): Record<string, Transaction[]> {
  return transactions.reduce((acc, tx) => {
    const date = tx.date.split('T')[0];
    if (!acc[date]) acc[date] = [];
    acc[date].push(tx);
    return acc;
  }, {} as Record<string, Transaction[]>);
}
```

### 5. 機能層（Features）

#### 5.1 Report Management Features（新規作成）

**useMonthlyReport.ts（新規）**
```typescript
// /features/report-management/useMonthlyReport.ts

import { useMemo } from 'react';
import { calculateMonthlyReport } from '@/entities/report';
import type { Transaction, CategorySummary } from '@/entities';

export const useMonthlyReport = (
  transactions: Transaction[] | undefined,
  categorySummary: CategorySummary[] | undefined
) => {
  return useMemo(() => {
    if (!transactions || !categorySummary) return null;
    return calculateMonthlyReport(transactions, categorySummary);
  }, [transactions, categorySummary]);
};
```

**usePeriodComparison.ts（新規）**
```typescript
// /features/report-management/usePeriodComparison.ts

import { useMemo } from 'react';
import { useTransactions } from '@/features/transaction-management';
import {
  calculateComparisonPeriod,
  calculatePeriodComparison,
} from '@/entities/report';
import type { DateRange, ComparisonPeriod } from '@/entities/report';

export const usePeriodComparison = (
  currentRange: DateRange,
  period: ComparisonPeriod
) => {
  // 前期間の計算
  const previousRange = useMemo(
    () => calculateComparisonPeriod(currentRange, period),
    [currentRange, period]
  );

  // 現在期間のデータ取得
  const { data: currentTransactions } = useTransactions({
    from: format(currentRange.from, 'yyyy-MM-dd'),
    to: format(currentRange.to, 'yyyy-MM-dd'),
  });

  // 前期間のデータ取得
  const { data: previousTransactions } = useTransactions({
    from: format(previousRange.from, 'yyyy-MM-dd'),
    to: format(previousRange.to, 'yyyy-MM-dd'),
  });

  return useMemo(() => {
    if (!currentTransactions || !previousTransactions) return null;

    const current = calculatePeriodData(currentTransactions);
    const previous = calculatePeriodData(previousTransactions);

    return calculatePeriodComparison(current, previous);
  }, [currentTransactions, previousTransactions]);
};

function calculatePeriodData(transactions: Transaction[]): PeriodData {
  const income = transactions
    .filter((t) => t.type === 'income')
    .reduce((sum, t) => sum + t.amount, 0);

  const expenses = transactions
    .filter((t) => t.type === 'expense')
    .reduce((sum, t) => sum + t.amount, 0);

  return {
    income,
    expenses,
    netIncome: income - expenses,
    transactionCount: transactions.length,
  };
}
```

**useCategoryTrend.ts（新規）**
```typescript
// /features/report-management/useCategoryTrend.ts

import { useMemo } from 'react';
import { transformToTrendData } from '@/entities/report';
import type { CategorySummary, Category } from '@/entities';

export const useCategoryTrend = (
  categorySummary: CategorySummary[] | undefined,
  selectedCategories: string[]
) => {
  return useMemo(() => {
    if (!categorySummary) return [];

    // 選択されたカテゴリのみフィルタリング
    const filtered = selectedCategories.length
      ? categorySummary.filter((c) =>
          selectedCategories.includes(c.category_id)
        )
      : categorySummary;

    return transformToTrendData(filtered);
  }, [categorySummary, selectedCategories]);
};
```

## データフローアーキテクチャ

### 1. 状態管理フロー:
```
ReportsContainer
    ↓
フィルター状態（ローカル）
    ├── dateRange
    ├── selectedCategories
    └── comparisonPeriod
         ↓
並列データフェッチ
    ├── useTransactions (current period)
    ├── useTransactions (previous period)
    ├── useTransactionMonthlySummary
    ├── useTransactionCategorySummary
    ├── useBudgetSummary
    └── useCategories
         ↓
    APIクライアント → バックエンド
         ↓
    計算・変換フック
    ├── useMonthlyReport
    ├── usePeriodComparison
    └── useCategoryTrend
         ↓
    ウィジェットコンポーネント
```

### 2. コンポーネント階層:
```
ReportsContainer
├── ReportFilters
├── MonthlyReportCard
├── PeriodComparisonWidget
├── BudgetVsActualWidget
└── Charts Grid
    ├── CategoryAnalysisChart
    └── TrendAnalysisChart
```

### 3. データ依存関係:
```
MonthlyReportCard
    ← useTransactions (取引一覧)
    ← useTransactionCategorySummary (カテゴリサマリー)
    ← useMonthlyReport (計算処理)

PeriodComparisonWidget
    ← useTransactions (現在期間)
    ← useTransactions (前期間)
    ← usePeriodComparison (比較計算)

BudgetVsActualWidget
    ← useBudgetSummary (予算)
    ← useTransactionCategorySummary (実績)

CategoryAnalysisChart
    ← useTransactionCategorySummary
    ← useCategories

TrendAnalysisChart
    ← useTransactionCategorySummary
    ← useCategoryTrend (トレンド変換)
```

## 実装ステップ

### Phase 1: エンティティ層作成（1日）
1. ✅ `/entities/report/` ディレクトリ作成
2. ✅ 型定義作成（`report.types.ts`）
3. ✅ 計算ロジック作成（`report.calculations.ts`）
4. ✅ 比較ロジック作成（`report.comparisons.ts`）
5. ✅ フォーマッター作成（`report.formatters.ts`）

### Phase 2: 機能層作成（1日）
1. ✅ `/features/report-management/` ディレクトリ作成
2. ✅ `useMonthlyReport` フック作成
3. ✅ `usePeriodComparison` フック作成
4. ✅ `useCategoryTrend` フック作成

### Phase 3: ウィジェット作成（3日）
1. ✅ `ReportFilters` 実装 - 0.5日
2. ✅ `MonthlyReportCard` 実装 - 0.5日
3. ✅ `PeriodComparisonWidget` 実装 - 0.5日
4. ✅ `BudgetVsActualWidget` 実装 - 0.5日
5. ✅ `CategoryAnalysisChart` 実装 - 0.5日
6. ✅ `TrendAnalysisChart` 実装 - 0.5日

### Phase 4: ページ統合（1日）
1. ✅ `/app/reports/page.tsx` 作成
2. ✅ `ReportsContainer` 作成
3. ✅ ウィジェット統合
4. ✅ フィルター機能実装
5. ✅ ローディング・エラーハンドリング

### Phase 5: テストと最適化（1日）
1. ✅ ユニットテスト作成
2. ✅ E2Eテスト作成
3. ✅ パフォーマンス最適化
4. ✅ レスポンシブ確認
5. ✅ アクセシビリティチェック

**総工数見積もり: 7日**

## ファイル数サマリー

| 層 | ディレクトリ | 新規ファイル数 | 流用ファイル数 |
|-------|-----------|------------|------------|
| App | `/app/reports/` | 1 | 0 |
| ページコンポーネント | `/page-components/reports/` | 1 | 0 |
| ウィジェット | `/widgets/reports/` | 6 | 0 |
| エンティティ | `/entities/report/` | 6 | 0 |
| 機能 | `/features/report-management/` | 3 | 0 |
| 機能 | `/features/transaction-management/` | 0 | 3 |
| 機能 | `/features/budget-management/` | 0 | 1 |
| 機能 | `/features/category-management/` | 0 | 1 |
| **合計** | | **17** | **5** |

## API対応表

| 機能 | バックエンドAPI | 実装状況 |
|------|--------------|---------|
| 月次サマリー | `GET /transactions/summary/monthly` | ✅ 実装済み |
| カテゴリ別サマリー | `GET /transactions/summary/by-category` | ✅ 実装済み |
| 期間別取引一覧 | `GET /transactions?from=&to=` | ✅ 実装済み |
| 予算サマリー | `GET /budgets/summary` | ✅ 実装済み |
| カテゴリ一覧 | `GET /categories` | ✅ 実装済み |
| 統合月次レポート | `GET /reports/summary/monthly` | ❌ 未実装（代替可能） |
| レポート出力 | `POST /reports/export` | ❌ 未実装（将来対応） |

**結論: レポート出力以外は既存APIで対応可能！**

## このアーキテクチャの利点

1. **既存APIの最大活用**: 新規APIなしで実装可能
2. **段階的実装**: ウィジェット単位で独立して実装可能
3. **柔軟なフィルタリング**: ユーザーが自由に期間・カテゴリを選択可能
4. **リアルタイム分析**: フィルター変更時に即座に再計算
5. **拡張性**: 将来のレポート出力機能追加が容易
6. **再利用性**: 既存のトランザクション・予算機能を活用

## 実装のベストプラクティス

### 1. データキャッシング:
```typescript
// 現在期間と前期間のデータを独立してキャッシュ
useTransactions({ from: '2024-01-01', to: '2024-01-31' }); // Current
useTransactions({ from: '2023-12-01', to: '2023-12-31' }); // Previous
```

### 2. 計算のメモ化:
```typescript
// 重い計算処理はメモ化
const monthlyReport = useMonthlyReport(transactions, categorySummary);
const periodComparison = usePeriodComparison(currentRange, period);
```

### 3. フィルター最適化:
```typescript
// デバウンスで不要なAPI呼び出しを削減
const debouncedDateRange = useDebounce(dateRange, 500);
```

### 4. チャートのパフォーマンス:
```typescript
// データ量が多い場合はサンプリング
const sampledData = useMemo(
  () => sampleData(trendData, 100),
  [trendData]
);
```

## まとめ

レポートページは**既存APIを組み合わせる**ことで実装可能です。

**主な利点:**
- ✅ 新規バックエンドAPI開発は不要
- ✅ 既存のトランザクション・予算APIで対応可能
- ✅ Feature-Sliced Designに準拠
- ✅ 柔軟なフィルタリング機能

**工数見積もり:**
- 最短: 6日（経験豊富な開発者）
- 標準: 7日（通常の開発者）
- 最長: 9日（初めてのプロジェクト参加者）

この実装により、ユーザーに**詳細な財務分析機能**と**柔軟なレポート機能**を提供する、完全に機能するレポートページが完成します。

## 将来の拡張機能

### Phase 2（将来実装）:
1. **レポート出力機能**
   - PDF出力
   - CSV出力
   - メール送信

2. **高度な分析**
   - 年次比較
   - 予測分析
   - 異常検知

3. **カスタムレポート**
   - ユーザー定義レポート
   - レポートテンプレート
   - 自動レポート配信

これらの機能は、バックエンドAPIが実装され次第、段階的に追加できる設計になっています。
