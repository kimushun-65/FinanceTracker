# 予算管理ページ 実装計画

## 概要
月別カテゴリ別予算の設定・編集・進捗モニタリングを行うページの実装計画。予算設定、実績との対比、リアルタイム進捗表示機能を提供する。

## 対応API
- `GET /budgets` - 予算一覧取得（月指定対応）
- `GET /budgets/current` - 現在月の予算情報取得（進捗計算含む）
- `POST /budgets` - 新規予算作成
- `PUT /budgets/{id}` - 予算情報更新
- `DELETE /budgets/{id}` - 予算削除
- `GET /categories` - カテゴリ一覧取得
- `GET /transactions/summary/monthly` - 実績データ取得

## 実装対象コンポーネント

### Pages
- `frontend/src/pages/budgets/index.tsx` - 予算管理ページ

### Features
- `frontend/src/features/budget/create-budget/` - 予算作成機能
- `frontend/src/features/budget/edit-budget/` - 予算編集機能
- `frontend/src/features/budget/delete-budget/` - 予算削除機能
- `frontend/src/features/budget/budget-progress/` - 進捗計算機能

### Widgets
- `frontend/src/widgets/budget-overview/` - 予算概要表示
- `frontend/src/widgets/budget-list/` - 予算一覧・編集テーブル
- `frontend/src/widgets/budget-progress/` - 進捗バー表示
- `frontend/src/widgets/last-month-review/` - 先月実績レビュー

### Entities
- `frontend/src/entities/budget/` - 予算エンティティ（既存利用）
- `frontend/src/entities/category/` - カテゴリエンティティ（既存利用）
- `frontend/src/entities/transaction/` - 取引エンティティ（実績計算用）

### Shared
- `frontend/src/shared/ui/progress-bar/` - プログレスバーコンポーネント
- `frontend/src/shared/ui/form/` - フォームコンポーネント
- `frontend/src/shared/ui/tabs/` - タブコンポーネント

## 実装フェーズ

### Phase 1: 基本表示機能
1. **現在月予算表示**
   - GET /budgets/current API連携
   - カテゴリ別予算額表示
   - 基本的なテーブル表示

2. **進捗表示**
   - 予算消化率計算
   - プログレスバー表示
   - 色分け（安全圏・注意・警告・超過）

### Phase 2: 編集機能
3. **予算編集機能**
   - PUT /budgets/{id} API連携
   - インライン編集実装
   - クイック調整ボタン（±1000円）

4. **予算作成機能**
   - POST /budgets API連携
   - 新月度の予算設定
   - 前月コピー機能

### Phase 3: 高度な機能
5. **先月実績レビュー**
   - 前月の実績vs予算比較
   - トレンド表示
   - 特筆すべき支出のハイライト

6. **予算テンプレート**
   - 前月と同じ予算設定
   - ゼロベース設定
   - 一括調整機能

## コンポーネント詳細設計

### 1. BudgetManagePage (pages/budgets/index.tsx)
```tsx
export default function BudgetManagePage() {
  return (
    <Layout>
      <BudgetOverviewWidget />
      <BudgetTabsWidget>
        <BudgetListWidget />
        <LastMonthReviewWidget />
      </BudgetTabsWidget>
    </Layout>
  )
}
```

### 2. BudgetOverviewWidget
- 今月の総予算額表示
- 全体進捗表示
- 月度切り替え機能

### 3. BudgetListWidget
- カテゴリ別予算一覧
- インライン編集機能
- 進捗バー表示

### 4. BudgetProgressWidget
- リアルタイム進捗計算
- 色分けプログレスバー
- 警告アラート表示

### 5. LastMonthReviewWidget
- 前月実績の詳細表示
- 予算vs実績の比較
- 週別支出パターン

## データフロー

### 予算取得フロー
```
Page Load → GET /budgets/current → Calculate Progress → Render Widgets
```

### 予算更新フロー
```
Inline Edit → Validation → PUT /budgets/{id} → Success → Refresh Progress
```

### 進捗計算フロー
```
Budget Data + Transaction Data → Calculate Usage → Update Progress Bars
```

## 状態管理

### Global State (Zustand)
- 予算リスト: `useBudgetStore`
- 現在月予算: `useCurrentBudgetStore`
- カテゴリリスト: `useCategoryStore`

### Local State
- 編集モード状態
- タブ選択状態
- 月度選択状態

## 予算進捗計算ロジック

### 進捗率計算
```typescript
const progressPercentage = (spent / budgetAmount) * 100;
```

### 色分けルール
- **0-70%**: 緑色（安全圏）
- **70-90%**: 黄色（注意）
- **90-100%**: オレンジ（警告）
- **100%超**: 赤色（超過）

### 期間考慮計算
```typescript
const daysInMonth = new Date(year, month, 0).getDate();
const daysPassed = new Date().getDate();
const expectedSpending = (budgetAmount * daysPassed) / daysInMonth;
const paceIndicator = spent / expectedSpending;
```

## バリデーション

### フロントエンド
- 予算額: 0以上の数値
- 必須カテゴリ: 最低1カテゴリは予算設定

### ビジネスルール
- 月度重複チェック
- 予算額の妥当性検証

## 画面設計 (ASCII)

```
+--------------------------------------------------------------------------------------------------+
|  FinSight                                                               AI User | Premium Plan    |
+----------------+---------------------------------------------------------------------------------+
| ≡ Dashboard    |                         Budget Management - December 2025                      |
|                |                                                                                 |
| 📊 Transactions | [ This Month Budget ] [ Last Month Review ]                                     |
|                | ────────────────────────────────────────────────────────────────────────────── |
| 💰 Budget       |                                                                                 |
|                | December 2025 Budget Setting                                                    |
| 📈 Assets       | +-----------------------------------------------------------------------------+ |
|                | | Category          Budget        Spent        Progress              Actions  | |
| 📑 Reports      | +-----------------------------------------------------------------------------+ |
|                | | 🍎 Food           ¥50,000       ¥5,432      [██░░░░░░░░░░] 11%      [±1k]   | |
| ⚙️ Settings     | | 🚃 Transportation ¥12,000       ¥1,200      [██░░░░░░░░░░] 10%      [±1k]   | |
|                | | 🏠 Household      ¥90,000       ¥90,000     [████████████] 100%     [±1k]   | |
|                | | 🎬 Entertainment  ¥20,000       ¥2,100      [██░░░░░░░░░░] 11%      [±1k]   | |
|                | | ⚡ Utilities      ¥7,000        ¥0          [░░░░░░░░░░░░] 0%       [±1k]   | |
|                | | 📦 Others         ¥15,000       ¥1,340      [█░░░░░░░░░░░] 9%       [±1k]   | |
|                | +-----------------------------------------------------------------------------+ |
|                | | Total:            ¥194,000      ¥100,072    [████████░░░░] 52%              | |
|                | +-----------------------------------------------------------------------------+ |
|                |                                                                                 |
|                | Quick Actions:                                                                  |
|                | [Copy Last Month] [Reset All] [Save Changes]                                   |
|                |                                                                                 |
+----------------+---------------------------------------------------------------------------------+

先月レビュータブ選択時:
+--------------------------------------------------------------------------------------------------+
|  FinSight                                                               AI User | Premium Plan    |
+----------------+---------------------------------------------------------------------------------+
| ≡ Dashboard    |                         Budget Management - December 2025                      |
|                |                                                                                 |
| 📊 Transactions | [ This Month Budget ] [ Last Month Review ]                                     |
|                | ────────────────────────────────────────────────────────────────────────────── |
| 💰 Budget       |                                                                                 |
|                | November 2025 Spending Review                                                   |
| 📈 Assets       | +-----------------------------------------------------------------------------+ |
|                | | Category          Spending     Budget      Difference   Status              | |
| 📑 Reports      | +-----------------------------------------------------------------------------+ |
|                | | 🍎 Food           ¥35,220     ¥40,000     -¥4,780      ✅ Under Budget      | |
| ⚙️ Settings     | | 🚃 Transportation ¥9,850      ¥10,000     -¥150        ✅ Under Budget      | |
|                | | 🏠 Household      ¥90,000     ¥90,000     ¥0           ⚠️ On Target        | |
|                | | 🎬 Entertainment  ¥18,930     ¥18,000     +¥930        ⚠️ Over Budget       | |
|                | | ⚡ Utilities      ¥4,890      ¥6,000      -¥1,110      ✅ Under Budget      | |
|                | | 📦 Others         ¥12,340     ¥10,000     +¥2,340      ❌ Over Budget       | |
|                | +-----------------------------------------------------------------------------+ |
|                | | Total:            ¥171,230    ¥174,000    -¥2,770      ✅ Under Budget      | |
|                | +-----------------------------------------------------------------------------+ |
|                |                                                                                 |
|                | Notable Expenses in November                                                    |
|                | +-----------------------------------------------------------------------------+ |
|                | | 📅 Nov 25  Year-end party preparation                           ¥8,500       | |
|                | | 📅 Nov 18  Winter clothes shopping                             ¥5,200       | |
|                | | 📅 Nov 12  Birthday dinner                                     ¥3,800       | |
|                | +-----------------------------------------------------------------------------+ |
|                |                                                                                 |
|                | [Apply Last Month Budget to This Month]                                        |
|                |                                                                                 |
+----------------+---------------------------------------------------------------------------------+

編集モード時:
+--------------------------------------------------------------------------------------------------+
|  FinSight                                                               AI User | Premium Plan    |
+----------------+---------------------------------------------------------------------------------+
| ≡ Dashboard    |                         Budget Management - December 2025                      |
|                |                                                                                 |
| 📊 Transactions | [ This Month Budget ] [ Last Month Review ]                                     |
|                | ────────────────────────────────────────────────────────────────────────────── |
| 💰 Budget       |                                                                                 |
|                | December 2025 Budget Setting                                                    |
| 📈 Assets       | +-----------------------------------------------------------------------------+ |
|                | | Category          Budget        Spent        Progress              Actions  | |
| 📑 Reports      | +-----------------------------------------------------------------------------+ |
|                | | 🍎 Food           [¥52,000___]  ¥5,432      [██░░░░░░░░░░] 10%   [Save][Cancel]| |
| ⚙️ Settings     | | 🚃 Transportation ¥12,000       ¥1,200      [██░░░░░░░░░░] 10%      [±1k]   | |
|                | | 🏠 Household      ¥90,000       ¥90,000     [████████████] 100%     [±1k]   | |
|                | | 🎬 Entertainment  ¥20,000       ¥2,100      [██░░░░░░░░░░] 11%      [±1k]   | |
|                | | ⚡ Utilities      ¥7,000        ¥0          [░░░░░░░░░░░░] 0%       [±1k]   | |
|                | | 📦 Others         ¥15,000       ¥1,340      [█░░░░░░░░░░░] 9%       [±1k]   | |
|                | +-----------------------------------------------------------------------------+ |
|                | | Total:            ¥196,000      ¥100,072    [████████░░░░] 51%              | |
|                | +-----------------------------------------------------------------------------+ |
|                |                                                                                 |
|                | Quick Actions:                                                                  |
|                | [Copy Last Month] [Reset All] [Save Changes]                                   |
|                |                                                                                 |
+----------------+---------------------------------------------------------------------------------+
```

### コンポーネント配置説明

#### ヘッダー部
- **ページタイトル**: Budget Management + 現在月表示
- **タブナビゲーション**: This Month Budget / Last Month Review

#### メインコンテンツエリア
1. **予算設定テーブル (今月予算タブ)**
   - カテゴリアイコン + 名前
   - 予算額（編集可能）
   - 支出実績
   - 進捗バー（色分け）
   - クイック調整ボタン（±1k）

2. **実績レビューテーブル (先月レビュータブ)**
   - 先月の支出実績
   - 予算との差分
   - ステータス表示（Under/On/Over Budget）
   - 特筆すべき支出のリスト

#### インタラクション
- **インライン編集**: 予算額のクリック編集
- **クイック調整**: ±1000円ボタン
- **一括操作**: 前月コピー・リセット・保存
- **タブ切り替え**: 今月予算⇔先月レビュー

## レスポンシブ対応

### Mobile (767px以下)
- テーブルをカード形式に変更
- 進捗バーを縦向きに調整
- ボタンサイズ拡大

### Tablet (768px-1279px)
- テーブル列幅調整
- 進捗バー幅調整

## パフォーマンス最適化

### 計算の最適化
- 進捗率計算のメモ化
- リアルタイム更新の頻度制限

### API最適化
- 予算・取引データの適切なキャッシュ
- 差分更新の実装

## 実装順序

1. **基盤実装** (1日)
   - ページレイアウト作成
   - API client設定
   - Store設定

2. **表示機能** (2日)
   - BudgetListWidget実装
   - 進捗計算ロジック実装
   - プログレスバー実装

3. **編集機能** (2日)
   - インライン編集実装
   - クイック調整実装
   - 一括操作実装

4. **レビュー機能** (1-2日)
   - 先月レビュータブ実装
   - 実績比較表示
   - 特筆支出リスト

5. **仕上げ** (1日)
   - レスポンシブ対応
   - エラーハンドリング
   - パフォーマンス最適化

## エラーハンドリング

### API エラー
- 予算データ取得失敗
- 更新処理失敗
- バリデーションエラー

### UI フィードバック
- 保存成功のトースト通知
- 予算超過の警告表示
- ローディング状態の表示

## アラート機能

### 予算超過アラート
- 90%到達時の警告表示
- 100%超過時の危険表示
- 色とアイコンでの視覚的警告

### 月初リマインダー
- 新月度の予算設定促進
- 前月実績のレビュー案内

## テスト戦略

### Unit Tests
- 進捗率計算ロジック
- 色分け判定ロジック
- インライン編集機能

### Integration Tests
- 予算CRUD操作フロー
- 実績データとの連携
- タブ切り替え動作

## 将来拡張予定

### AI予算提案機能（将来実装）
- 過去データ分析による予算最適化
- 季節性を考慮した提案
- 類似ユーザーとの比較

### 高度な分析機能
- 予算達成率トレンド
- カテゴリ別効率性分析
- 予算vs実績の詳細レポート