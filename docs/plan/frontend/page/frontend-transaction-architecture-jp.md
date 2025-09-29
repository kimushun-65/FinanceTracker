# トランザクションページアーキテクチャドキュメント

## 概要
この資料では、FinanceTrackerフロントエンドアプリケーションにおけるトランザクションページ実装の完全なディレクトリ構造とコンポーネントアーキテクチャについて説明します。

## 完全なディレクトリ構成

```
frontend/src/
├── app/
│   └── transactions/
│       └── page.tsx                          # エントリーポイント
├── page-components/
│   └── transaction/
│       └── ui/
│           └── transactionContainer.tsx      # メインオーケストレーター
├── widgets/
│   ├── transaction/
│   │   ├── transaction-list/
│   │   │   └── ui/
│   │   │       └── TransactionListTable.tsx # データテーブル
│   │   ├── transaction-filters/
│   │   │   ├── ui/
│   │   │   │   └── TransactionFilters.tsx   # フィルターUI
│   │   │   └── config/
│   │   │       └── periodOptions.ts         # 期間オプション設定
│   │   ├── transaction-summary/
│   │   │   └── ui/
│   │   │       └── TransactionSummaryCards.tsx # サマリーカード
│   │   ├── create-transaction/
│   │   │   └── CreateTransactionModal.tsx    # 新規作成モーダル
│   │   └── edit-transaction/
│   │       └── EditTransactionModal.tsx      # 編集モーダル
│   └── layout/
│       ├── AppLayout.tsx                     # アプリレイアウト
│       ├── Header.tsx                        # ヘッダー
│       └── Sidebar.tsx                       # サイドバー
├── entities/
│   ├── transaction/
│   │   ├── model/
│   │   │   ├── transaction.types.ts          # 型定義
│   │   │   ├── transaction.schema.ts         # バリデーションスキーマ
│   │   │   ├── transaction.constants.ts      # 定数
│   │   │   └── index.ts                      # モデルエクスポート
│   │   ├── api/
│   │   │   ├── transaction.client.ts         # APIクライアント
│   │   │   ├── transaction.endpoints.ts      # エンドポイント定義
│   │   │   ├── transaction.keys.ts           # React Queryキー
│   │   │   ├── transaction.filters.ts        # フィルター処理
│   │   │   └── index.ts                      # APIエクスポート
│   │   ├── lib/
│   │   │   ├── transaction.transformers.ts   # データ変換
│   │   │   ├── transaction.formatters.ts     # フォーマッター
│   │   │   ├── transaction.validations.ts    # バリデーション
│   │   │   ├── transaction.sorters.ts        # ソート処理
│   │   │   ├── transaction.aggregations.ts   # 集計処理
│   │   │   └── index.ts                      # ライブラリエクスポート
│   │   ├── ui/
│   │   │   ├── TransactionTypeBadge/
│   │   │   │   ├── TransactionTypeBadge.tsx  # タイプバッジ
│   │   │   │   └── index.ts
│   │   │   ├── TransactionAmountDisplay/
│   │   │   │   ├── TransactionAmountDisplay.tsx # 金額表示
│   │   │   │   └── index.ts
│   │   │   └── index.ts                      # UIエクスポート
│   │   └── index.ts                          # エンティティエクスポート
│   ├── account/                              # アカウントエンティティ
│   │   ├── model/
│   │   │   ├── account.types.ts              # アカウント型定義
│   │   │   └── index.ts
│   │   ├── lib/
│   │   │   ├── account.transformers.ts       # アカウント変換
│   │   │   └── index.ts
│   │   ├── ui/
│   │   │   └── AccountTypeBadge/
│   │   │       ├── AccountTypeBadge.tsx      # アカウントタイプ表示
│   │   │       └── index.ts
│   │   └── index.ts
│   └── category/                             # カテゴリーエンティティ
│       ├── model/
│       │   ├── category.types.ts             # カテゴリー型定義
│       │   └── index.ts
│       ├── lib/
│       │   ├── category.transformers.ts      # カテゴリー変換
│       │   └── index.ts
│       ├── ui/
│       │   ├── CategoryColorBadge/
│       │   │   ├── CategoryColorBadge.tsx    # カラーバッジ
│       │   │   └── index.ts
│       │   └── CategoryIcon/
│       │       ├── CategoryIcon.tsx          # アイコン表示
│       │       └── index.ts
│       └── index.ts
├── features/
│   ├── transaction-management/
│   │   ├── useTransactions.ts                # トランザクション一覧取得
│   │   ├── useTransactionMonthlySummary.ts   # 月次サマリー取得
│   │   ├── useCreateTransaction.ts           # 新規作成
│   │   ├── useUpdateTransaction.ts           # 更新
│   │   ├── useDeleteTransaction.ts           # 削除
│   │   ├── useTransactionFilters.ts          # フィルター管理
│   │   ├── useTransactionAggregates.ts       # 集計処理
│   │   └── index.ts                          # 機能エクスポート
│   ├── account-management/
│   │   ├── useAccounts.ts                    # アカウント一覧取得
│   │   └── index.ts
│   └── category-management/
│       ├── useCategories.ts                  # カテゴリー一覧取得
│       └── index.ts
└── shared/
    ├── ui/
    │   ├── button.tsx                        # ボタンコンポーネント
    │   ├── input.tsx                         # 入力フィールド
    │   ├── label.tsx                         # ラベル
    │   ├── select.tsx                        # セレクトボックス
    │   ├── modal.tsx                         # モーダル
    │   ├── table.tsx                         # テーブル
    │   ├── card.tsx                          # カード
    │   ├── loading.tsx                       # ローディング
    │   ├── toaster.tsx                       # トースト通知
    │   └── index.ts                          # UIエクスポート
    ├── value-objects/
    │   ├── money/
    │   │   ├── money.types.ts                # 金額型定義
    │   │   ├── money.formatters.ts           # 金額フォーマッター
    │   │   └── index.ts
    │   ├── date/
    │   │   ├── date.utils.ts                 # 日付ユーティリティ
    │   │   ├── date.formatters.ts            # 日付フォーマッター
    │   │   └── index.ts
    │   └── email/
    │       ├── email.types.ts                # メール型定義
    │       ├── email.validators.ts           # メールバリデーター
    │       └── index.ts
    ├── api/
    │   ├── client.ts                         # HTTPクライアント
    │   └── index.ts
    ├── types/
    │   ├── common.types.ts                   # 共通型定義
    │   └── index.ts
    ├── config/
    │   ├── env.ts                           # 環境設定
    │   └── index.ts
    ├── lib/
    │   └── hooks/
    │       ├── use-toast.ts                  # トーストフック
    │       └── index.ts
    └── index.ts                              # 共有エクスポート
```

## アーキテクチャ階層

### 1. ページ層
**エントリーポイント:**
```
/app/transactions/page.tsx
```
- メインルートのエントリーポイント
- `TransactionContainer`コンポーネントをレンダリング

### 2. ページコンポーネント層
**メインコンテナ:**
```
/page-components/transaction/ui/transactionContainer.tsx
```
- 中央オーケストレーター・コンポーネント
- グローバル状態を管理（フィルター、ページネーション、モーダル）
- すべてのトランザクションウィジェットを統括
- データ取得とエラー状態を処理

**主要な依存関係:**
- `@/widgets` - すべてのトランザクション関連ウィジェット
- `@/shared/ui` - UIコンポーネント（Button、Loading）
- `@/entities/transaction` - 型定義
- `@/features` - データ管理用のReact Queryフック

### 3. ウィジェット層
**トランザクションウィジェットディレクトリ: `/widgets/transaction/`**

#### 核となるウィジェットコンポーネント:
```
├── transaction-list/
│   └── ui/TransactionListTable.tsx          # CRUD操作付きデータテーブル
├── transaction-filters/
│   ├── ui/TransactionFilters.tsx            # カテゴリーと期間フィルター
│   └── config/periodOptions.ts              # フィルター設定
├── transaction-summary/
│   └── ui/TransactionSummaryCards.tsx       # 財務サマリーカード
├── create-transaction/
│   └── CreateTransactionModal.tsx           # 新規トランザクション作成フォーム
└── edit-transaction/
    └── EditTransactionModal.tsx             # 既存トランザクション編集フォーム
```

**ウィジェットの責任範囲:**
- **TransactionListTable**: トランザクションの表示、編集、削除
- **TransactionFilters**: カテゴリーと日付範囲でのフィルタリング
- **TransactionSummaryCards**: 収入/支出のサマリー
- **CreateTransactionModal**: 新規トランザクションの作成
- **EditTransactionModal**: トランザクションの修正

### 4. エンティティ層
**トランザクションエンティティディレクトリ: `/entities/transaction/`**

#### モデル層:
```
├── model/
│   ├── transaction.types.ts                 # TypeScript型定義
│   ├── transaction.schema.ts                # バリデーションスキーマ
│   ├── transaction.constants.ts             # 定数と列挙型
│   └── index.ts                            # モデルエクスポート
```

#### API層:
```
├── api/
│   ├── transaction.client.ts                # APIクライアント関数
│   ├── transaction.endpoints.ts             # APIエンドポイント定義
│   ├── transaction.keys.ts                 # React Queryキー
│   ├── transaction.filters.ts              # クエリパラメータ処理
│   └── index.ts                            # APIエクスポート
```

#### ビジネスロジック層:
```
├── lib/
│   ├── transaction.transformers.ts          # データ変換ユーティリティ
│   ├── transaction.formatters.ts           # 表示フォーマット関数
│   ├── transaction.validations.ts          # フォームバリデーションロジック
│   ├── transaction.sorters.ts              # ソートユーティリティ
│   ├── transaction.aggregations.ts         # データ集計関数
│   └── index.ts                            # ライブラリエクスポート
```

#### UIコンポーネント:
```
├── ui/
│   ├── TransactionTypeBadge/
│   │   ├── TransactionTypeBadge.tsx        # 収入/支出タイプバッジ
│   │   └── index.ts
│   ├── TransactionAmountDisplay/
│   │   ├── TransactionAmountDisplay.tsx    # フォーマット済み金額表示
│   │   └── index.ts
│   └── index.ts                            # UIエクスポート
```

### 5. 機能層
**トランザクション管理ディレクトリ: `/features/transaction-management/`**

#### React Queryフック:
```
├── useTransactions.ts                       # フィルター付きトランザクション一覧
├── useTransactionMonthlySummary.ts         # 月次サマリーデータ
├── useCreateTransaction.ts                  # 新規トランザクション作成
├── useUpdateTransaction.ts                  # 既存トランザクション更新
├── useDeleteTransaction.ts                  # トランザクション削除
├── useTransactionFilters.ts                # フィルター状態管理
├── useTransactionAggregates.ts             # 集計計算
└── index.ts                                # 機能エクスポート
```

**フックの責任範囲:**
- React Queryを使用したサーバー状態管理
- API呼び出しのオーケストレーション
- キャッシュ無効化戦略
- エラー処理とローディング状態

### 6. 関連エンティティ依存関係

#### アカウントエンティティ (`/entities/account/`):
- トランザクションフォームでのアカウント選択に使用
- アカウント型定義と表示コンポーネントを提供
- アカウントバリデーションと変換ユーティリティ

#### カテゴリーエンティティ (`/entities/category/`):
- トランザクションの分類に使用
- カテゴリー型定義とバリデーション
- カテゴリー表示コンポーネント（アイコン、カラーバッジ）
- カテゴリー変換ユーティリティ

### 7. 共有層依存関係

#### UIコンポーネント (`/shared/ui/`):
```
├── button.tsx                              # プライマリアクションボタン
├── input.tsx                               # フォーム入力フィールド
├── label.tsx                               # フォームラベル
├── select.tsx                              # ドロップダウン選択
├── modal.tsx                               # モーダルラッパーコンポーネント
├── table.tsx                               # データテーブルコンポーネント
├── card.tsx                                # カードレイアウトコンポーネント
├── loading.tsx                             # ローディング状態コンポーネント
└── toaster.tsx                            # トースト通知
```

#### 値オブジェクト (`/shared/value-objects/`):
```
├── money/                                  # 金額型とフォーマット
├── date/                                   # 日付ユーティリティとフォーマット
└── email/                                  # メールバリデーションユーティリティ
```

#### 共有ユーティリティ (`/shared/`):
```
├── api/client.ts                           # HTTPクライアント設定
├── types/                                  # 共有TypeScript型
├── config/env.ts                          # 環境設定
└── lib/hooks/use-toast.ts                 # トースト通知フック
```

### 8. レイアウトコンポーネント

#### アプリレイアウト (`/widgets/layout/`):
```
├── AppLayout.tsx                           # メインアプリケーションレイアウト
├── Header.tsx                              # アプリケーションヘッダー
└── Sidebar.tsx                            # ナビゲーションサイドバー
```

## データフローアーキテクチャ

### 1. 状態管理フロー:
```
TransactionContainer → 機能フック → APIクライアント → バックエンド
                  ↑                                      ↓
            ウィジェットコンポーネント ← エンティティ変換器 ← APIレスポンス
```

### 2. コンポーネント階層:
```
TransactionContainer
├── TransactionSummaryCards
├── TransactionFilters
├── TransactionListTable
├── CreateTransactionModal
└── EditTransactionModal
```

### 3. データ依存関係:
- **トランザクション**: ページのメインデータ
- **カテゴリー**: 分類とフィルタリングに必要
- **アカウント**: トランザクション作成/編集に必要
- **月次サマリー**: サマリーカードに必要

## 主要な設計原則

### 1. Feature-Sliced Design (FSD):
- **App**: ルートレベルコンポーネント
- **Pages**: ページオーケストレーションコンポーネント
- **Widgets**: 複雑なUIウィジェット
- **Features**: ビジネスロジックとデータ管理
- **Entities**: ドメインモデルとビジネスルール
- **Shared**: 再利用可能なユーティリティとコンポーネント

### 2. 関心の分離:
- **UIコンポーネント**: 純粋なプレゼンテーションロジック
- **機能フック**: サーバー状態とビジネスロジック
- **エンティティ層**: ドメインモデルと変換
- **共有層**: 横断的関心事

### 3. コンポーネント合成:
- コンテナコンポーネントは状態とオーケストレーションを管理
- ウィジェットコンポーネントは特定のUI機能を処理
- エンティティコンポーネントはドメイン固有のUI要素を提供
- 共有コンポーネントは汎用的なUIビルディングブロックを提供

## ファイル数サマリー

| 層 | ディレクトリ | ファイル数 |
|-------|-----------|--------|
| App | `/app/transactions/` | 1 |
| ページコンポーネント | `/page-components/transaction/` | 1 |
| ウィジェット | `/widgets/transaction/` | 6 |
| エンティティ | `/entities/transaction/` | 15 |
| 機能 | `/features/transaction-management/` | 8 |
| 関連エンティティ | `/entities/account/`, `/entities/category/` | 10 |
| 共有 | `/shared/` | 15+ |
| **合計** | | **56+** |

## このアーキテクチャの利点

1. **モジュラリティ**: 各層が明確な責任を持つ
2. **再利用性**: 共有コンポーネントはページ間で使用可能
3. **テスト性**: 各層を独立してテスト可能
4. **保守性**: 変更が特定の層に局所化される
5. **スケーラビリティ**: 既存コードに影響を与えずに新機能を追加可能
6. **型安全性**: 全層にわたる強力なTypeScript型付け
7. **パフォーマンス**: React Queryによる最適化されたデータ取得とキャッシング

## 今後の機能拡張

1. **リアルタイム更新**: ライブトランザクション更新のためのWebSocket統合
2. **高度なフィルタリング**: より洗練されたフィルターオプション
3. **一括操作**: 複数トランザクション操作
4. **エクスポート機能**: CSV/PDFエクスポート機能
5. **高度な分析**: チャートと詳細レポート
6. **オフラインサポート**: オフライントランザクション入力のためのPWA機能

## 実装のベストプラクティス

### 1. 状態管理:
- React Queryを使用したサーバー状態管理
- ローカル状態は必要最小限に抑制
- フィルター状態はURLパラメータと同期

### 2. エラー処理:
- すべてのAPI呼び出しでエラー境界を使用
- ユーザーフレンドリーなエラーメッセージ
- ネットワークエラーの適切な処理

### 3. パフォーマンス最適化:
- React.memoによる不要な再レンダリング防止
- 大きなリストの仮想化
- 画像の遅延読み込み
- バンドルサイズの最適化

### 4. アクセシビリティ:
- キーボードナビゲーション対応
- スクリーンリーダー対応
- 適切なARIAラベル
- カラーコントラストの確保

この包括的なアーキテクチャにより、maintainable（保守しやすい）、scalable（拡張可能）、そしてuser-friendly（使いやすい）なトランザクション管理システムが実現されています。