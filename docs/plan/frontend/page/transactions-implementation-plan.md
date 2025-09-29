# 取引管理ページ 実装計画

## 概要
収入・支出の取引記録・管理を行うページの実装計画。バックエンドAPIとの連携により、CRUD操作とリアルタイム集計を実現する。

## 対応API
- `GET /transactions` - 取引一覧取得（フィルタ・ページング対応）
- `POST /transactions` - 取引新規作成
- `PUT /transactions/{id}` - 取引情報更新
- `DELETE /transactions/{id}` - 取引削除
- `GET /transactions/summary/monthly` - 月次サマリー取得
- `GET /categories` - カテゴリ一覧取得

## 実装対象コンポーネント

### Pages
- `frontend/src/pages/transactions/index.tsx` - 取引一覧ページ

### Features
- `frontend/src/features/transaction/create-transaction/` - 取引作成機能
- `frontend/src/features/transaction/edit-transaction/` - 取引編集機能
- `frontend/src/features/transaction/delete-transaction/` - 取引削除機能
- `frontend/src/features/transaction/filter-transactions/` - 取引フィルタ機能

### Widgets
- `frontend/src/widgets/transaction-summary/` - 月次サマリーカード
- `frontend/src/widgets/transaction-list/` - 取引一覧テーブル
- `frontend/src/widgets/transaction-filters/` - フィルタコンポーネント

### Entities
- `frontend/src/entities/transaction/` - 取引エンティティ（既存利用）
- `frontend/src/entities/category/` - カテゴリエンティティ（既存利用）

### Shared
- `frontend/src/shared/ui/modal/` - モーダルコンポーネント
- `frontend/src/shared/ui/form/` - フォームコンポーネント
- `frontend/src/shared/ui/pagination/` - ページネーションコンポーネント

## 実装フェーズ

### Phase 1: 基本表示機能
1. **取引一覧表示**
   - GET /transactions API連携
   - 基本的なテーブル表示
   - ページネーション実装
   
2. **月次サマリー表示**
   - GET /transactions/summary/monthly API連携
   - サマリーカード表示（収入・支出・純収入）

### Phase 2: CRUD機能
3. **取引作成モーダル**
   - POST /transactions API連携
   - バリデーション実装
   - カテゴリ選択ドロップダウン
   
4. **取引編集機能**
   - PUT /transactions/{id} API連携
   - インライン編集またはモーダル編集
   
5. **取引削除機能**
   - DELETE /transactions/{id} API連携
   - 確認ダイアログ実装

### Phase 3: 高度な機能
6. **フィルタ機能**
   - 期間指定フィルタ
   - カテゴリ別フィルタ
   - 収入/支出フィルタ
   
7. **ソート機能**
   - 日付順、金額順、カテゴリ順

## コンポーネント詳細設計

### 1. TransactionListPage (pages/transactions/index.tsx)
```tsx
export default function TransactionListPage() {
  // ページ全体のレイアウトとstate管理
  return (
    <Layout>
      <TransactionSummaryWidget />
      <TransactionFiltersWidget />
      <TransactionListWidget />
      <CreateTransactionModal />
    </Layout>
  )
}
```

### 2. TransactionSummaryWidget
- 月次収入・支出・純収入の表示
- GET /transactions/summary/monthly から取得
- カード形式での表示

### 3. TransactionListWidget  
- 取引データのテーブル表示
- ページネーション対応
- 各行にEdit/Deleteアクション

### 4. CreateTransactionModal
- モーダル形式での取引作成
- POST /transactions への送信
- バリデーション機能

### 5. TransactionFiltersWidget
- 期間、カテゴリ、タイプでのフィルタ
- クエリパラメータとの同期

## データフロー

### 取得フロー
```
Page Load → API Call → Entity Update → Widget Re-render
```

### 作成フロー  
```
Form Submit → Validation → API Call → Success → List Refresh → Modal Close
```

### 更新フロー
```
Edit Click → Populate Form → API Call → Success → List Refresh
```

### 削除フロー
```
Delete Click → Confirmation → API Call → Success → List Refresh
```

## 状態管理

### Global State (Zustand)
- 取引リスト: `useTransactionStore`
- カテゴリリスト: `useCategoryStore`
- フィルタ状態: `useTransactionFilterStore`

### Local State
- モーダル開閉状態
- フォーム入力値
- ローディング状態

## API連携

### Error Handling
- ネットワークエラー時の再試行機能
- バリデーションエラーの表示
- 楽観的UI更新とrollback

### Caching Strategy
- React Query使用
- 取引リストの5分キャッシュ
- 作成・更新後のinvalidation

## バリデーション

### フロントエンド
- 金額: 必須、正の数値
- カテゴリ: 必須選択
- 日付: 必須、フォーマット検証

### エラー表示
- インラインエラーメッセージ
- トースト通知でAPI エラー

## レスポンシブ対応

### Mobile (767px以下)
- テーブルをカード形式に変更
- モーダルを全画面表示
- フィルタを折りたたみ表示

### Tablet (768px-1279px)  
- テーブル列数を調整
- モーダルサイズ調整

## パフォーマンス最適化

### 仮想化
- 大量データ対応のため、react-window使用検討

### レンダリング最適化
- React.memo使用
- useMemoでの計算結果キャッシュ
- useCallbackでの関数メモ化

## テスト戦略

### Unit Tests
- コンポーネントのレンダリングテスト
- ユーザーインタラクションテスト
- API連携のモックテスト

### Integration Tests
- ページ全体のユーザーフローテスト
- フィルタ機能の統合テスト

## 実装順序

1. **基盤実装** (1-2日)
   - ページレイアウト作成
   - API client設定
   - Store設定

2. **一覧表示** (2-3日)  
   - TransactionListWidget実装
   - ページネーション実装
   - サマリーカード実装

3. **CRUD機能** (3-4日)
   - 作成モーダル実装
   - 編集機能実装  
   - 削除機能実装

4. **フィルタ・ソート** (2-3日)
   - フィルタUI実装
   - APIクエリパラメータ連携

5. **仕上げ** (1-2日)
   - レスポンシブ対応
   - エラーハンドリング強化
   - パフォーマンス最適化

## 画面設計 (ASCII)

```
+--------------------------------------------------------------------------------------------------+
|  FinSight                                                               AI User | Premium Plan    |
+----------------+---------------------------------------------------------------------------------+
| ≡ Dashboard    |                              Transactions                      [+ Add Transaction]|
|                |                                                                                 |
| 📊 Transactions | +----------------+  +----------------+  +----------------+                     |
|                | | Monthly Income |  | Monthly Expense|  | Net Income     |                     |
| 💰 Budget       | |                |  |                |  |                |                     |
|                | |      ¥0        |  |   ¥106,672    |  |   -¥106,672    |                     |
| 📈 Assets       | +----------------+  +----------------+  +----------------+                     |
|                |                                                                                 |
| 📑 Reports      | Filter: [All Types ▼] [All Categories ▼] [This Month ▼]                       |
|                |                                                                                 |
| ⚙️ Settings     | Transaction List                                                               |
|                | +-----------------------------------------------------------------------------+ |
|                | | Date    | Description      | Category      | Amount        | Actions     | |
|                | +-----------------------------------------------------------------------------+ |
|                | | 09/30   | スーパー          | 🍎 Food       | ↓¥7,450      | ✏️ 🗑️       | |
|                | | 09/29   | Amazon           | 🛍️ Shopping   | ↓¥3,280      | ✏️ 🗑️       | |
|                | | 09/29   | 電気代           | ⚡ Utilities  | ↓¥5,073      | ✏️ 🗑️       | |
|                | | 09/28   | カフェ           | 🍎 Food       | ↓¥1,200      | ✏️ 🗑️       | |
|                | | 09/27   | 電車賃           | 🚃 Transport  | ↓¥440        | ✏️ 🗑️       | |
|                | | 09/27   | ランチ           | 🍎 Food       | ↓¥890        | ✏️ 🗑️       | |
|                | | 09/26   | Netflix          | 🎬 Entertainment | ↓¥1,490     | ✏️ 🗑️       | |
|                | | 09/25   | コンビニ         | 🍎 Food       | ↓¥523        | ✏️ 🗑️       | |
|                | | 09/25   | 書籍             | 📚 Education  | ↓¥2,420      | ✏️ 🗑️       | |
|                | | 09/24   | ガソリン         | 🚗 Transport  | ↓¥5,230      | ✏️ 🗑️       | |
|                | +-----------------------------------------------------------------------------+ |
|                |                                                                                 |
|                | [◀ Previous] Page 1 of 5 [Next ▶]                                              |
|                |                                                                                 |
+----------------+---------------------------------------------------------------------------------+

モーダル表示時:
+--------------------------------------------------------------------------------------------------+
|  FinSight                                                               AI User | Premium Plan    |
+----------------+---------------------------------------------------------------------------------+
| ≡ Dashboard    |                              Transactions                      [+ Add Transaction]|
|                |                                                                                 |
| 📊 Transactions | +----------------+  +----------------+  +----------------+                     |
|                | | Monthly Income |  | Monthly Expense|  | Net Income     |                     |
| 💰 Budget       | |                |  |                |  |                |                     |
|                | |      ¥0        |  |   ¥106,672    |  |   -¥106,672    |                     |
| 📈 Assets       | +----------------+  +----------------+  +----------------+                     |
|                |                          +---------------------------+                          |
| 📑 Reports      |                          | Add Transaction      [X] |                          |
|                |                          +---------------------------+                          |
| ⚙️ Settings     |                          | Type: [Expense] Income   |                          |
|                |                          |                           |                          |
|                |                          | Amount:                   |                          |
|                |                          | ¥ [___________________]   |                          |
|                |                          |                           |                          |
|                |                          | Category:                 |                          |
|                |                          | [Select Category     ▼]   |                          |
|                |                          |                           |                          |
|                |                          | Date:                     |                          |
|                |                          | [2025/09/30        📅]   |                          |
|                |                          |                           |                          |
|                |                          | Description:              |                          |
|                |                          | [_____________________]   |                          |
|                |                          |                           |                          |
|                |                          | [Cancel]        [Add]     |                          |
|                |                          +---------------------------+                          |
+----------------+---------------------------------------------------------------------------------+
```

### コンポーネント配置説明

#### ヘッダー部
- **ナビゲーション**: 左サイドバーにメニュー項目
- **ユーザー情報**: 右上にユーザー名とプラン表示
- **Add Transactionボタン**: 右上に配置

#### メインコンテンツエリア
1. **サマリーカード (上部)**
   - 3つの等幅カード配置
   - 月次収入・支出・純収入表示
   - 金額は大きなフォントで強調

2. **フィルタ部 (中上部)**
   - 横並びドロップダウン配置
   - タイプ・カテゴリ・期間フィルタ

3. **取引一覧テーブル (中央)**
   - 列: 日付・説明・カテゴリ・金額・アクション
   - カテゴリはアイコン+名前表示
   - 金額は収入/支出で色分け
   - アクションは編集・削除アイコン

4. **ページネーション (下部)**
   - 中央配置でナビゲーションボタン

#### 取引追加モーダル
- **オーバーレイ**: 背景を半透明マスク
- **中央配置**: 画面中央に配置
- **フィールド**: 縦並びレイアウト
- **ボタン**: 下部右寄せ配置

## 依存関係
- Categories実装完了後にカテゴリ選択機能実装
- Account実装完了後に口座選択機能実装