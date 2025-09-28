# 口座管理ページ 実装計画

## 概要  
複数の金融口座を統合管理し、総資産の把握を行うページの実装計画。預金、投資信託、株式などの資産を一元管理する。

## 対応API
- `GET /accounts` - 口座一覧取得（ページング対応）
- `POST /accounts` - 新規口座作成
- `PUT /accounts/{id}` - 口座情報更新
- `DELETE /accounts/{id}` - 口座削除

## 実装対象コンポーネント

### Pages
- `frontend/src/pages/accounts/index.tsx` - 口座管理ページ

### Features  
- `frontend/src/features/account/create-account/` - 口座作成機能
- `frontend/src/features/account/edit-account/` - 口座編集機能
- `frontend/src/features/account/delete-account/` - 口座削除機能

### Widgets
- `frontend/src/widgets/total-assets/` - 総資産表示
- `frontend/src/widgets/account-list/` - 口座一覧テーブル
- `frontend/src/widgets/asset-chart/` - 資産構成チャート（将来実装）

### Entities
- `frontend/src/entities/account/` - 口座エンティティ（既存利用）

### Shared  
- `frontend/src/shared/ui/modal/` - モーダルコンポーネント
- `frontend/src/shared/ui/form/` - フォームコンポーネント
- `frontend/src/shared/ui/confirmation/` - 確認ダイアログ

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

### 1. AccountListPage (pages/accounts/index.tsx)
```tsx
export default function AccountListPage() {
  return (
    <Layout>
      <TotalAssetsWidget />
      <AccountListWidget />
      <CreateAccountModal />
    </Layout>
  )
}
```

### 2. TotalAssetsWidget
- 総資産の大型表示
- 前月比の簡易計算
- 更新日時表示

### 3. AccountListWidget
- 口座データのテーブル表示
- 構成比率の計算と表示
- 各行にEdit/Deleteアクション

### 4. CreateAccountModal  
- モーダル形式での口座作成
- POST /accounts への送信
- 口座種別選択ドロップダウン

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