# 実装完了済みAPI一覧

## 概要
このドキュメントは、FinanceTrackerバックエンドで実装・デバッグが完了したAPIエンドポイントの一覧です。
すべてのAPIは認証が必要で、Auth0のJWTトークンを使用します。

## 認証API

### Auth Handler
| メソッド | エンドポイント | 説明 | ステータス |
|---------|---------------|------|-----------|
| GET | `/auth/login` | Auth0ログインページへリダイレクト | ✅ 実装済み |
| GET | `/auth/callback` | Auth0認証コールバック処理 | ✅ 実装済み |
| POST | `/auth/logout` | ログアウト処理 | ✅ 実装済み |
| GET | `/auth/user` | 現在の認証済みユーザー情報取得 | ✅ 実装済み |
| POST | `/auth/token` | HttpOnlyクッキーにトークン設定 | ✅ 実装済み |
| DELETE | `/auth/token` | HttpOnlyクッキーからトークン削除 | ✅ 実装済み |
| GET | `/auth/check` | 認証状態チェック | ✅ 実装済み |

## ユーザーAPI

### User Handler
| メソッド | エンドポイント | 説明 | ステータス |
|---------|---------------|------|-----------|
| GET | `/api/v1/users/me` | 現在のユーザー情報取得 | ✅ 実装・デバッグ済み |
| PUT | `/api/v1/users/me` | 現在のユーザー情報更新 | ✅ 実装・デバッグ済み |

## 口座API

### Account Handler
| メソッド | エンドポイント | 説明 | ステータス |
|---------|---------------|------|-----------|
| GET | `/api/v1/accounts` | 口座一覧取得 | ✅ 実装・デバッグ済み |
| POST | `/api/v1/accounts` | 口座作成 | ✅ 実装・デバッグ済み |
| GET | `/api/v1/accounts/:id` | 口座詳細取得 | ✅ 実装・デバッグ済み |
| PUT | `/api/v1/accounts/:id` | 口座更新 | ✅ 実装・デバッグ済み |
| DELETE | `/api/v1/accounts/:id` | 口座削除 | ✅ 実装・デバッグ済み |

## 取引API

### Transaction Handler
| メソッド | エンドポイント | 説明 | ステータス | 備考 |
|---------|---------------|------|-----------|------|
| GET | `/api/v1/transactions` | 取引一覧取得 | ✅ 実装・デバッグ済み | 日付フィルタリング修正済み |
| POST | `/api/v1/transactions` | 取引作成 | ✅ 実装・デバッグ済み | |
| GET | `/api/v1/transactions/:id` | 取引詳細取得 | ✅ 実装・デバッグ済み | |
| PUT | `/api/v1/transactions/:id` | 取引更新 | ✅ 実装・デバッグ済み | 認証エラー修正済み |
| DELETE | `/api/v1/transactions/:id` | 取引削除 | ✅ 実装・デバッグ済み | 認証エラー修正済み |
| GET | `/api/v1/transactions/summary/monthly` | 月次サマリー取得 | ✅ 実装・デバッグ済み | 認証エラー修正済み |

### 取引一覧のクエリパラメータ
- `account_id`: 口座IDでフィルタ
- `category_id`: カテゴリIDでフィルタ
- `transaction_type`: 取引タイプ（income/expense）でフィルタ
- `date_from`: 開始日（YYYY-MM-DD形式）
- `date_to`: 終了日（YYYY-MM-DD形式）
- `limit`: 取得件数制限（デフォルト: 100）
- `offset`: 取得開始位置（デフォルト: 0）
- `order_by`: ソート順（デフォルト: transaction_date desc）

## カテゴリAPI

### Category Handler
| メソッド | エンドポイント | 説明 | ステータス | 備考 |
|---------|---------------|------|-----------|------|
| GET | `/api/v1/categories` | カテゴリ一覧取得 | ✅ 実装・デバッグ済み | name フィールド追加済み |
| POST | `/api/v1/categories` | カテゴリ作成 | ✅ 実装・デバッグ済み | name フィールド追加済み |
| GET | `/api/v1/categories/:id` | カテゴリ詳細取得 | ✅ 実装・デバッグ済み | name フィールド追加済み |
| PUT | `/api/v1/categories/:id` | カテゴリ更新 | ✅ 実装・デバッグ済み | name フィールド追加済み |
| DELETE | `/api/v1/categories/:id` | カテゴリ削除（論理削除） | ✅ 実装・デバッグ済み | |
| GET | `/api/v1/categories/master` | マスターカテゴリ一覧取得 | ✅ 実装・デバッグ済み | 認証エラー修正済み |

### カテゴリ一覧のクエリパラメータ
- `category_master_id`: カテゴリマスターIDでフィルタ
- `is_active`: アクティブ状態でフィルタ（true/false）
- `order_by`: ソート順（デフォルト: created_at desc）

### カテゴリの仕組み
- **カテゴリマスター**: システム全体で共有される カテゴリのテンプレート（食費、交通費など）
- **ユーザーカテゴリ**: 各ユーザーが実際に使用するカテゴリ（カスタム名の設定可能）
- カテゴリ削除は論理削除（`is_active`をfalseに設定）

## 予算API

### Budget Handler
| メソッド | エンドポイント | 説明 | ステータス | 備考 |
|---------|---------------|------|-----------|------|
| GET | `/api/v1/budgets` | 予算一覧取得 | ✅ 実装・デバッグ済み | 認証エラー修正済み |
| POST | `/api/v1/budgets` | 予算作成 | ✅ 実装・デバッグ済み | 日付パースエラー修正済み |
| GET | `/api/v1/budgets/:id` | 予算詳細取得 | ✅ 実装・デバッグ済み | 認証エラー修正済み |
| PUT | `/api/v1/budgets/:id` | 予算更新 | ✅ 実装・デバッグ済み | 認証エラー修正済み |
| DELETE | `/api/v1/budgets/:id` | 予算削除 | ✅ 実装・デバッグ済み | 認証エラー修正済み |
| GET | `/api/v1/budgets/current` | 現在の期間の予算取得 | ✅ 実装・デバッグ済み | 認証エラー修正済み |

### 予算一覧のクエリパラメータ
- `category_id`: カテゴリIDでフィルタ
- `period`: 期間タイプ（monthly/yearly）でフィルタ
- `is_active`: アクティブ状態でフィルタ（true/false）
- `order_by`: ソート順（デフォルト: start_date desc）

### 予算作成/更新の日付形式
- `start_date`: YYYY-MM-DD形式の文字列
- `end_date`: YYYY-MM-DD形式の文字列（オプション）

## システムエンドポイント

| メソッド | エンドポイント | 説明 | ステータス |
|---------|---------------|------|-----------|
| GET | `/health` | ヘルスチェック | ✅ 実装済み |
| GET | `/api/v1/` | API情報取得 | ✅ 実装済み |

## 共通仕様

### 認証
すべてのAPI（ヘルスチェックとAPI情報を除く）は認証が必要です。
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

### 主なステータスコード
- 200: 成功
- 201: 作成成功
- 400: バリデーションエラー
- 401: 認証エラー
- 403: アクセス権限なし
- 404: リソースが見つからない
- 409: 競合（重複など）
- 500: 内部エラー

## 修正履歴

### 2025-09-26
- **認証エラー修正**: `getUserID`ヘルパー関数を使用して、全ハンドラーの認証フローを統一
- **日付フィルタリング修正**: トランザクションAPIの日付範囲クエリを修正
- **カテゴリ名フィールド追加**: カテゴリAPIのレスポンスにカテゴリマスターの名前を含めるよう修正
- **日付パースエラー修正**: 予算APIの日付フィールドを文字列形式（YYYY-MM-DD）で受け取るよう修正
- **decimal.Decimal バリデーション修正**: Goのバリデーターがdecimal型をサポートしないため、gt/gte制約を削除

## 未実装API

以下のAPIは今後実装予定です：
- 資産管理API（/api/v1/assets/*）
- レポートAPI（/api/v1/reports/*）
- 通知設定API（/api/v1/notification-settings）