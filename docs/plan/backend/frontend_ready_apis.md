# フロントエンド実装用 API対応表

## 概要
バックエンドの実装が完了した主要APIと、各フロントエンドページとの対応関係をまとめたドキュメント。
フロントエンド開発時の参照用。

## 完了済みAPI一覧

### 認証関連API ✅ **完全対応**
| エンドポイント | メソッド | 機能 | 対応ページ |
|---------------|---------|------|-----------|
| `/auth/login` | GET | Auth0ログインページリダイレクト | ログイン |
| `/auth/callback` | GET | Auth0認証コールバック | ログイン |
| `/auth/logout` | POST | ログアウト | 全ページ |
| `/auth/user` | GET | 認証済みユーザー情報取得 | 全ページ |
| `/auth/check` | GET | 認証状態チェック | 全ページ |

### ユーザー管理API ✅ **完全対応**
| エンドポイント | メソッド | 機能 | 対応ページ |
|---------------|---------|------|-----------|
| `/api/v1/users/me` | GET | ユーザー情報取得 | ダッシュボード、設定 |
| `/api/v1/users/me` | PUT | ユーザー情報更新 | 設定 |

### 口座管理API ✅ **完全対応**
| エンドポイント | メソッド | 機能 | 対応ページ |
|---------------|---------|------|-----------|
| `/api/v1/accounts` | GET | 口座一覧取得 | 資産管理 |
| `/api/v1/accounts` | POST | 口座作成 | 資産管理 |
| `/api/v1/accounts/:id` | GET | 口座詳細取得 | 資産管理 |
| `/api/v1/accounts/:id` | PUT | 口座更新 | 資産管理 |
| `/api/v1/accounts/:id` | DELETE | 口座削除 | 資産管理 |

### 取引管理API ✅ **完全対応**
| エンドポイント | メソッド | 機能 | 対応ページ |
|---------------|---------|------|-----------|
| `/api/v1/transactions` | GET | 取引一覧取得 | 取引管理、ダッシュボード |
| `/api/v1/transactions` | POST | 取引作成 | 取引管理 |
| `/api/v1/transactions/:id` | GET | 取引詳細取得 | 取引管理 |
| `/api/v1/transactions/:id` | PUT | 取引更新 | 取引管理 |
| `/api/v1/transactions/:id` | DELETE | 取引削除 | 取引管理 |
| `/api/v1/transactions/summary/monthly` | GET | 月次サマリー取得 | ダッシュボード、取引管理 |

#### 取引一覧API クエリパラメータ
- `account_id`: 口座IDでフィルタ
- `category_id`: カテゴリIDでフィルタ  
- `transaction_type`: 取引タイプ（income/expense）でフィルタ
- `date_from`: 開始日（YYYY-MM-DD形式）
- `date_to`: 終了日（YYYY-MM-DD形式）
- `limit`: 取得件数制限（デフォルト: 100）
- `offset`: 取得開始位置（デフォルト: 0）
- `order_by`: ソート順（デフォルト: transaction_date desc）

### カテゴリ管理API ✅ **完全対応** 
| エンドポイント | メソッド | 機能 | 対応ページ |
|---------------|---------|------|-----------|
| `/api/v1/categories` | GET | カテゴリ一覧取得 | 取引管理、ダッシュボード |
| `/api/v1/categories` | POST | カテゴリ作成 | 設定 |
| `/api/v1/categories/:id` | GET | カテゴリ詳細取得 | 設定 |
| `/api/v1/categories/:id` | PUT | カテゴリ更新 | 設定 |
| `/api/v1/categories/:id` | DELETE | カテゴリ削除 | 設定 |
| `/api/v1/categories/master` | GET | マスターカテゴリ一覧取得 | 設定 |

#### カテゴリ一覧API クエリパラメータ
- `category_master_id`: カテゴリマスターIDでフィルタ
- `is_active`: アクティブ状態でフィルタ（true/false）
- `order_by`: ソート順（デフォルト: created_at desc）

### 予算管理API ✅ **完全対応**
| エンドポイント | メソッド | 機能 | 対応ページ |
|---------------|---------|------|-----------|
| `/api/v1/budgets` | GET | 予算一覧取得 | 予算管理 |
| `/api/v1/budgets` | POST | 予算作成 | 予算管理 |
| `/api/v1/budgets/:id` | GET | 予算詳細取得 | 予算管理 |
| `/api/v1/budgets/:id` | PUT | 予算更新 | 予算管理 |
| `/api/v1/budgets/:id` | DELETE | 予算削除 | 予算管理 |
| `/api/v1/budgets/current` | GET | 現在期間の予算取得 | 予算管理、ダッシュボード |

#### 予算一覧API クエリパラメータ
- `category_id`: カテゴリIDでフィルタ
- `period`: 期間タイプ（monthly/yearly）でフィルタ
- `is_active`: アクティブ状態でフィルタ（true/false）
- `order_by`: ソート順（デフォルト: start_date desc）

### システムAPI ✅ **完全対応**
| エンドポイント | メソッド | 機能 | 対応ページ |
|---------------|---------|------|-----------|
| `/health` | GET | ヘルスチェック | 監視用 |
| `/api/v1/` | GET | API情報取得 | システム情報 |

## 各ページとAPI対応関係

### ダッシュボードページ ✅ **API完全対応**

#### 必要API
1. **月次サマリー取得** - `GET /api/v1/transactions/summary/monthly`
   - 収入・支出・純収入の計算
   - 期間指定パラメータ対応

2. **取引一覧取得** - `GET /api/v1/transactions` 
   - カテゴリ別集計用データ取得
   - 日付フィルタリングによる期間指定

3. **カテゴリ一覧取得** - `GET /api/v1/categories`
   - グラフ表示用カテゴリ情報
   - カテゴリ名とカラー情報を含む

4. **ユーザー情報取得** - `GET /api/v1/users/me`
   - ユーザー名、プラン情報表示

#### 実装可能機能
- ✅ 月別収支トレンドグラフ（6ヶ月）
- ✅ カテゴリ別支出ドーナツチャート
- ✅ サマリーカード（収入・支出・純収入）
- ✅ ユーザー情報表示

### 取引管理ページ ✅ **API完全対応**

#### 必要API
1. **取引CRUD操作**
   - `GET /api/v1/transactions` - 一覧取得（フィルタ・ページング対応）
   - `POST /api/v1/transactions` - 新規作成
   - `PUT /api/v1/transactions/:id` - 更新
   - `DELETE /api/v1/transactions/:id` - 削除

2. **月次サマリー取得** - `GET /api/v1/transactions/summary/monthly`
   - 上部カード表示用

3. **カテゴリ一覧取得** - `GET /api/v1/categories`
   - 取引作成・編集時の選択肢

#### 実装可能機能
- ✅ 取引一覧表示（ソート・フィルタ・ページング）
- ✅ 取引追加モーダル
- ✅ 取引編集・削除
- ✅ 月間サマリーカード

### 予算管理ページ ✅ **API完全対応**

#### 必要API
1. **予算CRUD操作**
   - `GET /api/v1/budgets` - 一覧取得
   - `POST /api/v1/budgets` - 新規作成
   - `PUT /api/v1/budgets/:id` - 更新
   - `DELETE /api/v1/budgets/:id` - 削除

2. **現在期間予算取得** - `GET /api/v1/budgets/current`
   - 進捗表示用

3. **カテゴリ一覧取得** - `GET /api/v1/categories`
   - 予算設定時の選択肢

#### 実装可能機能
- ✅ 予算一覧表示
- ✅ 予算作成・編集
- ✅ 予算進捗表示
- ✅ カテゴリ別予算管理

### 資産管理ページ ⚠️ **API未対応（今後実装予定）**
以下のAPIが未実装のため、資産管理ページは現在実装不可：

- `GET /api/v1/assets/snapshots` - 資産スナップショット一覧
- `POST /api/v1/assets/snapshots` - 資産スナップショット作成
- `GET /api/v1/assets/forecasts` - 資産予測一覧
- `GET /api/v1/assets/forecasts/latest` - 最新資産予測

**対応方針**: 口座一覧API（`GET /api/v1/accounts`）で基本的な資産表示は可能

### 設定ページ ✅ **基本機能API対応**

#### 対応済みAPI
1. **ユーザー情報管理**
   - `GET /api/v1/users/me` - ユーザー情報取得
   - `PUT /api/v1/users/me` - ユーザー情報更新

2. **カテゴリ管理**
   - `GET /api/v1/categories` - カテゴリ一覧
   - `POST /api/v1/categories` - カテゴリ作成
   - `PUT /api/v1/categories/:id` - カテゴリ更新
   - `DELETE /api/v1/categories/:id` - カテゴリ削除
   - `GET /api/v1/categories/master` - マスターカテゴリ一覧

#### 未対応API（今後実装予定）
- 通知設定API（notification-settings）

## 共通仕様

### 認証
すべてのAPI（ヘルスチェックとAPI情報を除く）は認証が必要：
- Authorization header: `Bearer {JWT_TOKEN}`
- Auth0のJWTトークンを使用

### エラーレスポンス形式
```json
{
    "error": "エラーメッセージ",
    "code": "エラーコード（オプション）",
    "details": {} // 詳細情報（オプション）
}
```

### 主要ステータスコード
- 200: 成功
- 201: 作成成功
- 400: バリデーションエラー
- 401: 認証エラー
- 403: アクセス権限なし
- 404: リソースが見つからない
- 500: 内部エラー

## フロントエンド実装推奨順序

### Phase 1: 基本機能（2週間）
1. **認証ページ** - Auth0統合
2. **ダッシュボードページ** - 基本レイアウト、API連携
3. **取引管理ページ** - CRUD操作、フィルタリング

### Phase 2: 拡張機能（1週間）
1. **予算管理ページ** - 予算設定、進捗表示
2. **設定ページ** - ユーザー情報、カテゴリ管理

### Phase 3: 高度な機能（将来実装）
1. **資産管理ページ** - 資産API実装後
2. **レポートページ** - レポートAPI実装後

## 技術的備考

### 日付形式
すべての日付は**YYYY-MM-DD形式**（文字列）で統一済み

### 金額形式  
金額はdecimal.Decimal型で管理、フロントエンドでは数値として扱う

### カテゴリ情報
カテゴリAPIレスポンスにはカテゴリマスターの名前情報が含まれるため、
グラフ表示やUI表示に直接利用可能

### 認証フロー
getUserIDヘルパー関数により認証処理が統一されているため、
フロントエンドは標準的なJWT認証フローで実装可能