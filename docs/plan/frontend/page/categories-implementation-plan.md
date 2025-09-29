# カテゴリ管理ページ 実装計画

## 概要
取引で使用するカテゴリの管理ページの実装計画。マスターカテゴリの選択・カスタムカテゴリの作成・編集・削除機能を提供する。

## 対応API
- `GET /categories/master` - マスターカテゴリ一覧取得
- `GET /categories` - ユーザーカテゴリ一覧取得  
- `POST /categories` - カテゴリ新規作成
- `PUT /categories/{id}` - カテゴリ更新（名前・有効状態）
- `DELETE /categories/{id}` - カテゴリ削除（無効化）

## 実装対象コンポーネント

### Pages
- `frontend/src/pages/categories/index.tsx` - カテゴリ管理ページ

### Features
- `frontend/src/features/category/create-category/` - カテゴリ作成機能
- `frontend/src/features/category/edit-category/` - カテゴリ編集機能
- `frontend/src/features/category/toggle-category/` - カテゴリ有効化切り替え

### Widgets
- `frontend/src/widgets/master-categories/` - マスターカテゴリ選択
- `frontend/src/widgets/user-categories/` - ユーザーカテゴリ一覧
- `frontend/src/widgets/category-usage-stats/` - カテゴリ使用統計（将来実装）

### Entities
- `frontend/src/entities/category/` - カテゴリエンティティ（既存利用）

### Shared
- `frontend/src/shared/ui/toggle/` - トグルスイッチ
- `frontend/src/shared/ui/modal/` - モーダルコンポーネント
- `frontend/src/shared/ui/color-picker/` - カラーピッカー（将来実装）

## 実装フェーズ

### Phase 1: 基本表示機能
1. **マスターカテゴリ表示**
   - GET /categories/master API連携
   - チェックボックス形式での表示
   - カテゴリアイコン・名前表示

2. **ユーザーカテゴリ表示**
   - GET /categories API連携
   - 有効/無効状態の表示
   - カスタム名表示

### Phase 2: カテゴリ選択機能
3. **マスターカテゴリ選択**
   - POST /categories でマスターカテゴリを有効化
   - チェックボックスでの選択・解除
   - リアルタイム状態更新

4. **カテゴリ有効化切り替え**
   - PUT /categories/{id} でis_active切り替え
   - トグルスイッチでの操作
   - 即座のUI反映

### Phase 3: カスタムカテゴリ機能  
5. **カスタムカテゴリ作成**
   - POST /categories でカスタムカテゴリ作成
   - カスタム名設定機能
   - マスターカテゴリベースでの作成

6. **カテゴリ編集機能**
   - PUT /categories/{id} でカスタム名変更
   - インライン編集実装

## コンポーネント詳細設計

### 1. CategoryManagePage (pages/categories/index.tsx)
```tsx
export default function CategoryManagePage() {
  return (
    <Layout>
      <MasterCategoriesWidget />
      <UserCategoriesWidget />
      <CreateCategoryModal />
    </Layout>
  )
}
```

### 2. MasterCategoriesWidget
- マスターカテゴリの一覧表示
- チェックボックスでの選択状態管理
- カテゴリ有効化API呼び出し

### 3. UserCategoriesWidget
- ユーザーカテゴリの一覧表示
- 有効/無効のトグル切り替え
- カスタム名の編集機能

### 4. CreateCategoryModal（将来実装）
- マスターカテゴリベースでのカスタムカテゴリ作成
- カスタム名・色の設定

## データフロー

### カテゴリ選択フロー
```
Checkbox Change → POST /categories → Success → Refetch Categories → Update UI
```

### 有効化切り替えフロー
```
Toggle Switch → PUT /categories/{id} → Success → Update Local State
```

### カスタムカテゴリ作成フロー
```
Form Submit → POST /categories → Success → Refetch Categories → Modal Close
```

## 状態管理

### Global State (Zustand)
- マスターカテゴリ: `useMasterCategoryStore`
- ユーザーカテゴリ: `useCategoryStore`

### Local State
- 選択状態（チェックボックス）
- 編集モード状態
- モーダル開閉状態

## カテゴリタイプ定義

### マスターカテゴリ
- **Food**: 食費 🍎
- **Household**: 家賃・住居費 🏠  
- **Utilities**: 光熱費 ⚡
- **Transportation**: 交通費 🚃
- **Entertainment**: エンターテインメント 🎬
- **Others**: その他 📦
- **Income**: 収入 💰

### 表示形式
- アイコン + 日本語名 + 英語名
- カテゴリ種別での色分け

## ビジネスルール

### カテゴリ選択ルール
- 最低1つのカテゴリは有効状態を維持
- マスターカテゴリの選択でユーザーカテゴリが自動作成
- カスタム名が設定されていない場合はマスター名を使用

### 削除制限
- 取引で使用中のカテゴリは削除不可（無効化のみ）
- 削除時の確認ダイアログ表示

## バリデーション

### フロントエンド
- カスタム名: 1-50文字
- 最低1カテゴリ選択必須

### API制約
- 重複カテゴリ名チェック
- 必須カテゴリの削除防止

## 画面設計 (ASCII)

```
+--------------------------------------------------------------------------------------------------+
|  FinSight                                                               AI User | Premium Plan    |
+----------------+---------------------------------------------------------------------------------+
| ≡ Dashboard    |                              Category Management                               |
|                |                                                                                 |
| 📊 Transactions | Default Categories                                                             |
|                | Select categories to use in your transactions:                                 |
| 💰 Budget       |                                                                                 |
|                | +-----------------------------------------------------------------------------+ |
| 📈 Assets       | | [✓] 🍎 Food (食費)                                                          | |
|                | | [✓] 🏠 Household (家賃・住居費)                                                | |
| 📑 Reports      | | [✓] ⚡ Utilities (光熱費)                                                    | |
|                | | [✓] 🚃 Transportation (交通費)                                               | |
| ⚙️ Settings     | | [✓] 🎬 Entertainment (エンターテインメント)                                      | |
|                | | [✓] 📦 Others (その他)                                                      | |
|                | | [ ] 💰 Income (収入)                                                        | |
|                | +-----------------------------------------------------------------------------+ |
|                |                                                                                 |
|                | Active Categories                                                               |
|                | Manage your currently active categories:                                       |
|                |                                                                                 |
|                | +-----------------------------------------------------------------------------+ |
|                | | Category                    | Custom Name      | Active    | Actions        | |
|                | +-----------------------------------------------------------------------------+ |
|                | | 🍎 Food                     | 食費             | [ON] OFF  | [Edit]         | |
|                | | 🏠 Household                | 家賃・住居費      | [ON] OFF  | [Edit]         | |
|                | | ⚡ Utilities                | 光熱費           | [ON] OFF  | [Edit]         | |
|                | | 🚃 Transportation           | 交通費           | [ON] OFF  | [Edit]         | |
|                | | 🎬 Entertainment            | エンターテインメント | [ON] OFF  | [Edit]         | |
|                | | 📦 Others                   | その他           | [ON] OFF  | [Edit]         | |
|                | +-----------------------------------------------------------------------------+ |
|                |                                                                                 |
|                | [Save Changes]                                                                  |
|                |                                                                                 |
+----------------+---------------------------------------------------------------------------------+

編集モード時:
+--------------------------------------------------------------------------------------------------+
|  FinSight                                                               AI User | Premium Plan    |
+----------------+---------------------------------------------------------------------------------+
| ≡ Dashboard    |                              Category Management                               |
|                |                                                                                 |
| 📊 Transactions | Default Categories                                                             |
|                | Select categories to use in your transactions:                                 |
| 💰 Budget       |                                                                                 |
|                | +-----------------------------------------------------------------------------+ |
| 📈 Assets       | | [✓] 🍎 Food (食費)                                                          | |
|                | | [✓] 🏠 Household (家賃・住居費)                                                | |
| 📑 Reports      | | [✓] ⚡ Utilities (光熱費)                                                    | |
|                | | [✓] 🚃 Transportation (交通費)                                               | |
| ⚙️ Settings     | | [✓] 🎬 Entertainment (エンターテインメント)                                      | |
|                | | [✓] 📦 Others (その他)                                                      | |
|                | | [ ] 💰 Income (収入)                                                        | |
|                | +-----------------------------------------------------------------------------+ |
|                |                                                                                 |
|                | Active Categories                                                               |
|                | Manage your currently active categories:                                       |
|                |                                                                                 |
|                | +-----------------------------------------------------------------------------+ |
|                | | Category                    | Custom Name      | Active    | Actions        | |
|                | +-----------------------------------------------------------------------------+ |
|                | | 🍎 Food                     | [外食費_______]  | [ON] OFF  | [Save][Cancel] | |
|                | | 🏠 Household                | 家賃・住居費      | [ON] OFF  | [Edit]         | |
|                | | ⚡ Utilities                | 光熱費           | [ON] OFF  | [Edit]         | |
|                | | 🚃 Transportation           | 交通費           | [ON] OFF  | [Edit]         | |
|                | | 🎬 Entertainment            | エンターテインメント | [ON] OFF  | [Edit]         | |
|                | | 📦 Others                   | その他           | [ON] OFF  | [Edit]         | |
|                | +-----------------------------------------------------------------------------+ |
|                |                                                                                 |
|                | [Save Changes]                                                                  |
|                |                                                                                 |
+----------------+---------------------------------------------------------------------------------+
```

### コンポーネント配置説明

#### ヘッダー部
- **ページタイトル**: Category Management

#### メインコンテンツエリア
1. **デフォルトカテゴリ選択 (上部)**
   - チェックボックスリスト
   - アイコン + 英語名 + 日本語名表示
   - 選択状態でユーザーカテゴリが自動作成

2. **アクティブカテゴリ管理 (下部)**
   - テーブル形式での表示
   - カスタム名の編集機能
   - 有効/無効のトグル切り替え
   - 編集・保存・キャンセルボタン

#### インタラクション
- **チェックボックス**: マスターカテゴリの選択・解除
- **トグルスイッチ**: カテゴリの有効・無効切り替え
- **インライン編集**: カスタム名の変更
- **保存ボタン**: 変更内容の一括保存

## レスポンシブ対応

### Mobile (767px以下)
- テーブルをカード形式に変更
- チェックボックスリストを1列表示
- 編集フィールドを全幅表示

### Tablet (768px-1279px)
- テーブル列幅調整
- アイコンサイズ調整

## 実装順序

1. **基盤実装** (1日)
   - ページレイアウト作成
   - API client設定
   - Store設定

2. **表示機能** (1-2日)
   - MasterCategoriesWidget実装
   - UserCategoriesWidget実装
   - チェックボックス状態管理

3. **操作機能** (2日)
   - カテゴリ選択機能実装
   - トグル切り替え実装
   - インライン編集実装

4. **仕上げ** (1日)
   - バリデーション実装
   - エラーハンドリング
   - レスポンシブ対応

## テスト戦略

### Unit Tests
- カテゴリ選択ロジック
- 有効化切り替えロジック
- カスタム名バリデーション

### Integration Tests
- マスターカテゴリ選択→ユーザーカテゴリ作成フロー
- カテゴリ有効化→取引での利用可能性確認

## エラーハンドリング

### API エラー
- ネットワークエラー時の再試行
- バリデーションエラーの表示
- 重複エラーの対応

### UI フィードバック
- 保存成功のトースト通知
- 編集中の状態表示
- 必須カテゴリ削除防止メッセージ

## 今後の拡張予定

### Phase 4: 高度な機能（将来実装）
- カスタムアイコン選択
- カテゴリ色のカスタマイズ  
- カテゴリ使用統計の表示
- カテゴリのインポート/エクスポート