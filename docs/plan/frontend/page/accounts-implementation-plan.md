# 口座管理ページアーキテクチャドキュメント

## 概要
FinanceTrackerフロントエンドアプリケーションにおける口座管理ページ実装の完全なディレクトリ構造とコンポーネントアーキテクチャについて説明します。複数の金融口座を統合管理し、総資産の把握を行うページの実装計画です。

## 完全なディレクトリ構成

```
frontend/src/
├── app/
│   └── accounts/
│       └── page.tsx                          # エントリーポイント
├── page-components/
│   └── account/
│       └── ui/
│           └── accountContainer.tsx          # メインオーケストレーター
├── widgets/
│   ├── account/
│   │   ├── account-list/
│   │   │   └── ui/
│   │   │       └── AccountListTable.tsx     # 口座一覧テーブル
│   │   ├── total-assets/
│   │   │   └── ui/
│   │   │       └── TotalAssetsWidget.tsx    # 総資産表示
│   │   ├── create-account/
│   │   │   └── CreateAccountModal.tsx        # 新規作成モーダル
│   │   └── edit-account/
│   │       └── EditAccountModal.tsx          # 編集モーダル
│   └── layout/
│       ├── AppLayout.tsx                     # アプリレイアウト
│       ├── Header.tsx                        # ヘッダー
│       └── Sidebar.tsx                       # サイドバー
├── entities/
│   └── account/
│       ├── model/
│       │   ├── account.types.ts              # 型定義
│       │   ├── account.schema.ts             # バリデーションスキーマ
│       │   ├── account.constants.ts          # 定数
│       │   └── index.ts                      # モデルエクスポート
│       ├── api/
│       │   ├── account.client.ts             # APIクライアント
│       │   ├── account.endpoints.ts          # エンドポイント定義
│       │   ├── account.keys.ts               # React Queryキー
│       │   └── index.ts                      # APIエクスポート
│       ├── lib/
│       │   ├── account.transformers.ts       # データ変換
│       │   ├── account.formatters.ts         # フォーマッター
│       │   ├── account.validations.ts        # バリデーション
│       │   ├── account.aggregations.ts       # 集計処理
│       │   └── index.ts                      # ライブラリエクスポート
│       ├── ui/
│       │   ├── AccountTypeBadge/
│       │   │   ├── AccountTypeBadge.tsx      # アカウントタイプバッジ
│       │   │   └── index.ts
│       │   ├── AccountBalanceDisplay/
│       │   │   ├── AccountBalanceDisplay.tsx # 残高表示
│       │   │   └── index.ts
│       │   └── index.ts                      # UIエクスポート
│       └── index.ts                          # エンティティエクスポート
├── features/
│   └── account-management/
│       ├── useAccounts.ts                    # 口座一覧取得
│       ├── useCreateAccount.ts               # 新規作成
│       ├── useUpdateAccount.ts               # 更新
│       ├── useDeleteAccount.ts               # 削除
│       ├── useAccountAggregates.ts           # 集計処理
│       └── index.ts                          # 機能エクスポート
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
    │   └── index.ts
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

## 対応API
- `GET /accounts` - 口座一覧取得（ページング対応）
- `POST /accounts` - 新規口座作成
- `PUT /accounts/{id}` - 口座情報更新
- `DELETE /accounts/{id}` - 口座削除

## アーキテクチャ階層

### 1. ページ層
**エントリーポイント:**
```
/app/accounts/page.tsx
```
- メインルートのエントリーポイント
- `AccountContainer`コンポーネントをレンダリング

### 2. ページコンポーネント層
**メインコンテナ:**
```
/page-components/account/ui/accountContainer.tsx
```
- 中央オーケストレーター・コンポーネント
- グローバル状態を管理（フィルター、モーダル、総資産計算）
- すべての口座関連ウィジェットを統括
- データ取得とエラー状態を処理

**主要な依存関係:**
- `@/widgets` - すべての口座関連ウィジェット
- `@/shared/ui` - UIコンポーネント（Button、Loading）
- `@/entities/account` - 型定義
- `@/features` - データ管理用のReact Queryフック

### 3. ウィジェット層
**口座ウィジェットディレクトリ: `/widgets/account/`**

#### 核となるウィジェットコンポーネント:
```
├── account-list/
│   └── ui/AccountListTable.tsx               # CRUD操作付きデータテーブル
├── total-assets/
│   └── ui/TotalAssetsWidget.tsx              # 総資産表示ウィジェット
├── create-account/
│   └── CreateAccountModal.tsx                # 新規口座作成フォーム
└── edit-account/
    └── EditAccountModal.tsx                  # 既存口座編集フォーム
```

**ウィジェットの責任範囲:**
- **AccountListTable**: 口座の表示、編集、削除
- **TotalAssetsWidget**: 総資産とサマリー情報の表示
- **CreateAccountModal**: 新規口座の作成
- **EditAccountModal**: 口座情報の修正

### 4. エンティティ層
**口座エンティティディレクトリ: `/entities/account/`**

#### モデル層:
```
├── model/
│   ├── account.types.ts                      # TypeScript型定義
│   ├── account.schema.ts                     # バリデーションスキーマ
│   ├── account.constants.ts                  # 定数と列挙型
│   └── index.ts                             # モデルエクスポート
```

#### API層:
```
├── api/
│   ├── account.client.ts                     # APIクライアント関数
│   ├── account.endpoints.ts                  # APIエンドポイント定義
│   ├── account.keys.ts                      # React Queryキー
│   └── index.ts                             # APIエクスポート
```

#### ビジネスロジック層:
```
├── lib/
│   ├── account.transformers.ts               # データ変換ユーティリティ
│   ├── account.formatters.ts                # 表示フォーマット関数
│   ├── account.validations.ts               # フォームバリデーションロジック
│   ├── account.aggregations.ts              # データ集計関数（総資産計算等）
│   └── index.ts                             # ライブラリエクスポート
```

#### UIコンポーネント:
```
├── ui/
│   ├── AccountTypeBadge/
│   │   ├── AccountTypeBadge.tsx             # 口座タイプバッジ
│   │   └── index.ts
│   ├── AccountBalanceDisplay/
│   │   ├── AccountBalanceDisplay.tsx        # フォーマット済み残高表示
│   │   └── index.ts
│   └── index.ts                             # UIエクスポート
```

### 5. 機能層
**口座管理ディレクトリ: `/features/account-management/`**

#### React Queryフック:
```
├── useAccounts.ts                            # 口座一覧取得
├── useCreateAccount.ts                       # 新規口座作成
├── useUpdateAccount.ts                       # 既存口座更新
├── useDeleteAccount.ts                       # 口座削除
├── useAccountAggregates.ts                   # 集計計算（総資産等）
└── index.ts                                 # 機能エクスポート
```

**フックの責任範囲:**
- React Queryを使用したサーバー状態管理
- API呼び出しのオーケストレーション
- キャッシュ無効化戦略
- エラー処理とローディング状態

### 6. 共有層依存関係

#### UIコンポーネント (`/shared/ui/`):
```
├── button.tsx                               # プライマリアクションボタン
├── input.tsx                                # フォーム入力フィールド
├── label.tsx                                # フォームラベル
├── select.tsx                               # ドロップダウン選択
├── modal.tsx                                # モーダルラッパーコンポーネント
├── table.tsx                                # データテーブルコンポーネント
├── card.tsx                                 # カードレイアウトコンポーネント
├── loading.tsx                              # ローディング状態コンポーネント
└── toaster.tsx                             # トースト通知
```

#### 値オブジェクト (`/shared/value-objects/`):
```
├── money/                                   # 金額型とフォーマット
└── date/                                    # 日付ユーティリティとフォーマット
```

#### 共有ユーティリティ (`/shared/`):
```
├── api/client.ts                            # HTTPクライアント設定
├── types/                                   # 共有TypeScript型
├── config/env.ts                           # 環境設定
└── lib/hooks/use-toast.ts                  # トースト通知フック
```

## データフローアーキテクチャ

### 1. 状態管理フロー:
```
AccountContainer → 機能フック → APIクライアント → バックエンド
               ↑                                      ↓
         ウィジェットコンポーネント ← エンティティ変換器 ← APIレスポンス
```

### 2. コンポーネント階層:
```
AccountContainer
├── TotalAssetsWidget
├── AccountListTable
├── CreateAccountModal
└── EditAccountModal
```

### 3. データ依存関係:
- **口座**: ページのメインデータ
- **総資産**: 全口座の残高合計
- **構成比**: 各口座の資産比率

## 実装フェーズ

### Phase 1: 基本表示機能
1. **総資産表示**
   - 全口座の残高合計計算
   - 前月比計算（簡易実装）
   - 大型数値表示

2. **口座一覧表示**
   - GET /accounts API連携
   - テーブル形式表示
   - 口座種別アイコン表示

### Phase 2: CRUD機能  
3. **口座作成モーダル**
   - POST /accounts API連携
   - 口座名・種別・初期残高入力
   - バリデーション実装

4. **口座編集機能**
   - PUT /accounts/{id} API連携
   - インライン編集実装
   - 名称変更対応

5. **口座削除機能**
   - DELETE /accounts/{id} API連携
   - 確認ダイアログ実装
   - 関連データ警告表示

### Phase 3: 高度な機能（将来実装）
6. **資産推移グラフ**
   - 時系列データ表示
   - グラフライブラリ使用

7. **残高手動更新**
   - 取引以外での残高調整
   - 履歴記録機能

## コンポーネント詳細設計

### 1. AccountContainer (page-components/account/ui/accountContainer.tsx)
```tsx
export default function AccountContainer() {
  const { data: accounts, isLoading } = useAccounts()
  const { totalAssets, percentageChange } = useAccountAggregates(accounts)
  
  return (
    <div className="space-y-6">
      <TotalAssetsWidget 
        total={totalAssets} 
        change={percentageChange} 
      />
      <AccountListTable 
        accounts={accounts} 
        isLoading={isLoading} 
      />
      <CreateAccountModal />
    </div>
  )
}
```

### 2. TotalAssetsWidget
- 総資産の大型表示
- 前月比の簡易計算
- 更新日時表示

### 3. AccountListTable
- 口座データのテーブル表示
- 構成比率の計算と表示
- 各行にEdit/Deleteアクション

### 4. CreateAccountModal  
- モーダル形式での口座作成
- POST /accounts への送信
- 口座種別選択ドロップダウン

## 実装必要ファイル一覧

### 🔴 Phase 1: 新規作成必要ファイル

#### ページ層 (1ファイル)
- [ ] `frontend/src/app/accounts/page.tsx` - エントリーポイント

#### ページコンポーネント層 (1ファイル)
- [ ] `frontend/src/page-components/account/ui/accountContainer.tsx` - メインオーケストレーター

#### ウィジェット層 (4ファイル)
- [ ] `frontend/src/widgets/account/account-list/ui/AccountListTable.tsx` - 口座一覧テーブル
- [ ] `frontend/src/widgets/account/total-assets/ui/TotalAssetsWidget.tsx` - 総資産表示
- [ ] `frontend/src/widgets/account/create-account/CreateAccountModal.tsx` - 作成モーダル
- [ ] `frontend/src/widgets/account/edit-account/EditAccountModal.tsx` - 編集モーダル

#### エンティティ層 - モデル (4ファイル)
- [ ] `frontend/src/entities/account/model/account.types.ts` - 型定義
- [ ] `frontend/src/entities/account/model/account.schema.ts` - バリデーションスキーマ
- [ ] `frontend/src/entities/account/model/account.constants.ts` - 定数
- [ ] `frontend/src/entities/account/model/index.ts` - モデルエクスポート

#### エンティティ層 - API (4ファイル)
- [ ] `frontend/src/entities/account/api/account.client.ts` - APIクライアント
- [ ] `frontend/src/entities/account/api/account.endpoints.ts` - エンドポイント定義
- [ ] `frontend/src/entities/account/api/account.keys.ts` - React Queryキー
- [ ] `frontend/src/entities/account/api/index.ts` - APIエクスポート

#### エンティティ層 - ビジネスロジック (5ファイル)
- [ ] `frontend/src/entities/account/lib/account.transformers.ts` - データ変換
- [ ] `frontend/src/entities/account/lib/account.formatters.ts` - フォーマッター
- [ ] `frontend/src/entities/account/lib/account.validations.ts` - バリデーション
- [ ] `frontend/src/entities/account/lib/account.aggregations.ts` - 集計処理
- [ ] `frontend/src/entities/account/lib/index.ts` - ライブラリエクスポート

#### エンティティ層 - UI (4ファイル)
- [ ] `frontend/src/entities/account/ui/AccountTypeBadge/AccountTypeBadge.tsx` - タイプバッジ
- [ ] `frontend/src/entities/account/ui/AccountTypeBadge/index.ts` - バッジエクスポート
- [ ] `frontend/src/entities/account/ui/AccountBalanceDisplay/AccountBalanceDisplay.tsx` - 残高表示
- [ ] `frontend/src/entities/account/ui/AccountBalanceDisplay/index.ts` - 残高エクスポート
- [ ] `frontend/src/entities/account/ui/index.ts` - UIエクスポート
- [ ] `frontend/src/entities/account/index.ts` - エンティティエクスポート

#### 機能層 (6ファイル)
- [ ] `frontend/src/features/account-management/useAccounts.ts` - 口座一覧取得
- [ ] `frontend/src/features/account-management/useCreateAccount.ts` - 新規作成
- [ ] `frontend/src/features/account-management/useUpdateAccount.ts` - 更新
- [ ] `frontend/src/features/account-management/useDeleteAccount.ts` - 削除
- [ ] `frontend/src/features/account-management/useAccountAggregates.ts` - 集計処理
- [ ] `frontend/src/features/account-management/index.ts` - 機能エクスポート

### 🟢 既存利用可能ファイル

#### 共有UI層 (既存利用)
- [x] `frontend/src/shared/ui/button.tsx` - ボタンコンポーネント
- [x] `frontend/src/shared/ui/input.tsx` - 入力フィールド
- [x] `frontend/src/shared/ui/label.tsx` - ラベル
- [x] `frontend/src/shared/ui/select.tsx` - セレクトボックス
- [x] `frontend/src/shared/ui/modal.tsx` - モーダル
- [x] `frontend/src/shared/ui/table.tsx` - テーブル
- [x] `frontend/src/shared/ui/card.tsx` - カード
- [x] `frontend/src/shared/ui/loading.tsx` - ローディング
- [x] `frontend/src/shared/ui/toaster.tsx` - トースト通知

#### 共有ユーティリティ層 (既存利用)
- [x] `frontend/src/shared/value-objects/money/` - 金額型とフォーマット
- [x] `frontend/src/shared/value-objects/date/` - 日付ユーティリティ
- [x] `frontend/src/shared/api/client.ts` - HTTPクライアント
- [x] `frontend/src/shared/types/common.types.ts` - 共通型定義
- [x] `frontend/src/shared/config/env.ts` - 環境設定
- [x] `frontend/src/shared/lib/hooks/use-toast.ts` - トーストフック

#### レイアウトコンポーネント (既存利用)
- [x] `frontend/src/widgets/layout/AppLayout.tsx` - アプリレイアウト
- [x] `frontend/src/widgets/layout/Header.tsx` - ヘッダー
- [x] `frontend/src/widgets/layout/Sidebar.tsx` - サイドバー

### 📊 実装ファイル数サマリー

| カテゴリ | 新規作成 | 既存利用 | 合計 |
|----------|----------|----------|------|
| ページ層 | 1 | 0 | 1 |
| ページコンポーネント層 | 1 | 0 | 1 |
| ウィジェット層 | 4 | 3 | 7 |
| エンティティ層 | 17 | 0 | 17 |
| 機能層 | 6 | 0 | 6 |
| 共有層 | 0 | 15+ | 15+ |
| **合計** | **29** | **18+** | **47+** |

### 🚀 実装優先順序

#### Step 1: 基盤レイヤー (5ファイル)
1. `entities/account/model/account.types.ts` - 型定義
2. `entities/account/model/account.constants.ts` - 定数
3. `entities/account/api/account.endpoints.ts` - エンドポイント
4. `entities/account/api/account.keys.ts` - React Queryキー
5. `entities/account/model/index.ts` - エクスポート

#### Step 2: API・ビジネスロジック層 (6ファイル)
6. `entities/account/api/account.client.ts` - APIクライアント
7. `entities/account/lib/account.transformers.ts` - データ変換
8. `entities/account/lib/account.formatters.ts` - フォーマッター
9. `entities/account/lib/account.aggregations.ts` - 集計処理
10. `entities/account/api/index.ts` - APIエクスポート
11. `entities/account/lib/index.ts` - ライブラリエクスポート

#### Step 3: 機能フック層 (3ファイル)
12. `features/account-management/useAccounts.ts` - 一覧取得
13. `features/account-management/useAccountAggregates.ts` - 集計
14. `features/account-management/index.ts` - エクスポート

#### Step 4: UI表示レイヤー (8ファイル)
15. `entities/account/ui/AccountTypeBadge/AccountTypeBadge.tsx` - タイプバッジ
16. `entities/account/ui/AccountBalanceDisplay/AccountBalanceDisplay.tsx` - 残高表示
17. `entities/account/ui/AccountTypeBadge/index.ts`
18. `entities/account/ui/AccountBalanceDisplay/index.ts`
19. `entities/account/ui/index.ts`
20. `widgets/account/total-assets/ui/TotalAssetsWidget.tsx` - 総資産表示
21. `widgets/account/account-list/ui/AccountListTable.tsx` - 一覧テーブル
22. `entities/account/index.ts` - エンティティエクスポート

#### Step 5: ページコンポーネント層 (2ファイル)
23. `page-components/account/ui/accountContainer.tsx` - メインコンテナ
24. `app/accounts/page.tsx` - エントリーポイント

#### Step 6: CRUD機能層 (7ファイル)
25. `entities/account/model/account.schema.ts` - バリデーション
26. `entities/account/lib/account.validations.ts` - バリデーションロジック
27. `features/account-management/useCreateAccount.ts` - 作成フック
28. `features/account-management/useUpdateAccount.ts` - 更新フック
29. `features/account-management/useDeleteAccount.ts` - 削除フック
30. `widgets/account/create-account/CreateAccountModal.tsx` - 作成モーダル
31. `widgets/account/edit-account/EditAccountModal.tsx` - 編集モーダル

## データフロー

### 取得フロー
```
Page Load → GET /accounts → Calculate Total → Render Widgets
```

### 作成フロー
```
Form Submit → Validation → POST /accounts → Success → List Refresh → Modal Close
```

### 更新フロー
```
Edit Click → Populate Form → PUT /accounts/{id} → Success → List Refresh
```

### 削除フロー
```
Delete Click → Confirmation → DELETE /accounts/{id} → Success → List Refresh
```

## 状態管理

### Global State (Zustand)
- 口座リスト: `useAccountStore`
- 総資産情報: `useAssetStore`

### Local State
- モーダル開閉状態
- フォーム入力値
- ローディング状態

## 口座種別定義

### サポート口座タイプ
- **checking**: 普通預金 🏛️
- **savings**: 定期預金 💰  
- **investment**: 投資信託 📊
- **stock**: 株式 💹
- **other**: その他 📦

### 表示形式
- アイコン + 日本語名表示
- 種別ごとの色分け

## バリデーション

### フロントエンド
- 口座名: 必須、50文字以内
- 口座種別: 必須選択
- 初期残高: 0以上の数値

### エラー表示
- インラインエラーメッセージ
- 重複口座名の警告

## 計算ロジック

### 総資産計算
```typescript
const totalAssets = accounts.reduce((sum, account) => 
  sum + account.current_balance, 0
);
```

### 構成比率計算
```typescript
const percentage = (account.current_balance / totalAssets) * 100;
```

### 前月比計算（簡易）
```typescript
// 実装制約により簡易計算
const monthlyChange = account.current_balance - account.initial_balance;
```

## レスポンシブ対応

### Mobile (767px以下)
- テーブルをカード形式に変更
- 総資産を上部固定表示
- モーダルを全画面表示

### Tablet (768px-1279px)
- テーブル列幅調整
- アイコンサイズ調整

## パフォーマンス最適化

### メモ化
- 総資産計算のメモ化
- 構成比率計算のメモ化

### API最適化
- 口座一覧の適切なキャッシュ
- 更新後の差分更新

## 画面設計 (ASCII)

```
+--------------------------------------------------------------------------------------------------+
|  FinSight                                                               AI User | Premium Plan    |
+----------------+---------------------------------------------------------------------------------+
| ≡ Dashboard    |                              Assets Management                [+ Add Account]     |
|                |                                                                                 |
| 📊 Transactions |  +--------------------------------------------------------------------+          |
|                |  |                      Total Assets                                  |          |
| 💰 Budget       |  |                                                                  |          |
|                |  |                      ¥1,500,000                                   |          |
| 📈 Assets       |  |                                                                  |          |
|                |  |                   ▲ +¥50,000 (+3.4%)                              |          |
| 📑 Reports      |  +--------------------------------------------------------------------+          |
|                |                                                                                 |
| ⚙️ Settings     |  Account List                                                                    |
|                |  +--------------------------------------------------------------------+          |
|                |  | Account Name       | Type        | Balance     | % | Change | Actions  |          |
|                |  |-------------------|-------------|-------------|---|--------|----------|          |
|                |  | みずほ銀行          | 🏛️ 普通預金  | ¥1,200,000  | 80% | ▲ +¥20k | [Edit][Delete] |
|                |  | SBI証券           | 📊 投資信託  | ¥250,000    | 17% | ▲ +¥25k | [Edit][Delete] |
|                |  | 楽天証券          | 💹 株式      | ¥50,000     | 3%  | ▲ +¥5k  | [Edit][Delete] |
|                |  +--------------------------------------------------------------------+          |
|                |                                                                                 |
+----------------+---------------------------------------------------------------------------------+

口座追加モーダル表示時:
+--------------------------------------------------------------------------------------------------+
|  FinSight                                                               AI User | Premium Plan    |
+----------------+---------------------------------------------------------------------------------+
| ≡ Dashboard    |                              Assets Management                [+ Add Account]     |
|                |                                                                                 |
| 📊 Transactions |  +--------------------------------------------------------------------+          |
|                |  |                      Total Assets                                  |          |
| 💰 Budget       |  |                      ¥1,500,000                                   |          |
|                |  |                   ▲ +¥50,000 (+3.4%)                              |          |
| 📈 Assets       |  +--------------------------------------------------------------------+          |
|                |                          +---------------------------+                          |
| 📑 Reports      |                          | Add Account        [X]   |                          |
|                |                          +---------------------------+                          |
| ⚙️ Settings     |                          | Account Name:             |                          |
|                |                          | [_____________________]   |                          |
|                |                          |                           |                          |
|                |                          | Account Type:             |                          |
|                |                          | [Select Type         ▼]   |                          |
|                |                          |                           |                          |
|                |                          | Initial Balance:          |                          |
|                |                          | ¥ [___________________]   |                          |
|                |                          |                           |                          |
|                |                          | [Cancel]        [Add]     |                          |
|                |                          +---------------------------+                          |
|                |                                                                                 |
+----------------+---------------------------------------------------------------------------------+
```

### コンポーネント配置説明

#### ヘッダー部
- **ナビゲーション**: 左サイドバーにメニュー項目
- **Add Accountボタン**: 右上に配置

#### メインコンテンツエリア
1. **総資産表示 (上部)**
   - 大型数値で総資産表示
   - 前月比を増減アイコンで表示
   - 背景色で資産状況を表現

2. **口座一覧テーブル (中央)**
   - 列: 口座名・種別・残高・構成比・変動・アクション
   - 種別はアイコン+日本語表示
   - 構成比は%表示
   - アクションは編集・削除ボタン

#### 口座追加モーダル
- **中央配置**: 画面中央にオーバーレイ表示
- **フィールド**: 口座名・種別・初期残高
- **バリデーション**: リアルタイム入力検証

## テスト戦略

### Unit Tests
- 総資産計算ロジック
- 構成比率計算ロジック
- フォームバリデーション

### Integration Tests
- CRUD操作の統合テスト
- 口座一覧の更新確認

## 実装順序

1. **基盤実装** (1日)
   - ページレイアウト作成
   - API client設定
   - Store設定

2. **表示機能** (2日)
   - AccountListWidget実装
   - TotalAssetsWidget実装
   - 計算ロジック実装

3. **CRUD機能** (2-3日)
   - 作成モーダル実装
   - 編集機能実装
   - 削除機能実装

4. **仕上げ** (1日)
   - レスポンシブ対応
   - エラーハンドリング
   - パフォーマンス最適化

## 制限事項

### API制限による制約
- 資産推移履歴APIが未実装のため履歴グラフは将来実装
- 残高変動の詳細履歴は取引データから推測
- 資産予測機能は将来実装予定

### 代替実装
- 前月比は初期残高からの変動で代用
- 資産推移は口座残高の単純表示