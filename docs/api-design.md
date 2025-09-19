# REST API 設計書

## 概要

このドキュメントは、家計簿アプリケーション「FinSight」の REST API 設計を定義します。
各エンドポイントは RESTful な設計原則に従い、適切な HTTP メソッドとステータスコードを使用します。

## 共通仕様

### ベース URL

```
https://api.finsight.com/v1
```

### 認証

すべての API は Auth0 が発行した JWT アクセストークンによる認証が必要です（Auth0 コールバック・同期エンドポイントを除く）。

```
Authorization: Bearer <auth0_access_token>
```

**トークンの取得方法:**

1. クライアントアプリケーションで Auth0 SDK を使用してログイン
2. Auth0 から取得したアクセストークンを API リクエストに含める
3. トークンの有効期限が切れた場合は、Auth0 のリフレッシュトークンを使用して新しいアクセストークンを取得

### レスポンス形式

- Content-Type: application/json
- 日時フォーマット: ISO 8601 (YYYY-MM-DDTHH:mm:ssZ)
- 金額: 整数（円）
- UUID: RFC 4122 形式

### エラーレスポンス

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "エラーの詳細説明",
    "details": {}
  }
}
```

### ページネーション

リスト取得 API では以下のクエリパラメータをサポート：

| パラメータ | 型      | デフォルト  | 説明                     |
| ---------- | ------- | ----------- | ------------------------ |
| page       | integer | 1           | ページ番号（1 から開始） |
| limit      | integer | 20          | 1 ページあたりの件数     |
| sort       | string  | -created_at | ソート順（-で降順）      |

レスポンス形式：

```json
{
  "data": [...],
  "pagination": {
    "total": 100,
    "page": 1,
    "limit": 20,
    "total_pages": 5
  }
}
```

---

## 0. 認証 API（Auth0 統合）

### Auth0 認証フロー概要

本システムでは Auth0 を使用した Google 認証（ソーシャルログイン）のみを提供します。パスワードベースの認証は使用しません。

#### 認証フロー

1. **Google ログイン**: クライアントアプリが Auth0 経由で Google の認証画面にリダイレクト
2. **認証完了**: Google での認証後、Auth0 からアクセストークン・ID トークンを取得
3. **API 呼び出し**: アクセストークンを Authorization ヘッダーに含めて API を呼び出し
4. **トークン検証**: API 側で Auth0 の JWKS を使用してトークンを検証

#### サポートする認証プロバイダー

- Google（必須）
- その他のソーシャルログイン（将来的な拡張）

### Auth0 コールバック処理

```
POST /auth/callback
```

**説明**: Google 認証後のコールバック処理。ユーザー情報の同期を行います。

**リクエスト:**

```json
{
  "auth0_user_id": "google-oauth2|115744780651980123456",
  "email": "user@gmail.com",
  "email_verified": true,
  "name": "John Doe",
  "picture": "https://lh3.googleusercontent.com/...",
  "provider": "google-oauth2",
  "google_id": "115744780651980123456"
}
```

**レスポンス:**

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "auth0_user_id": "auth0|550e8400e29b41d4a716",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2024-01-01T09:00:00Z"
  },
  "is_new_user": true
}
```

### ユーザー同期

```
POST /auth/sync
```

**説明**: Auth0 のユーザー情報をシステムに同期します。初回ログイン時や情報更新時に使用。

**リクエストヘッダー:**

```
Authorization: Bearer <auth0_access_token>
```

**レスポンス:**

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "auth0_user_id": "auth0|550e8400e29b41d4a716",
    "email": "user@example.com",
    "name": "John Doe",
    "synchronized_at": "2024-01-01T09:00:00Z"
  }
}
```

### ログアウト（クライアント側実装）

ログアウトは Auth0 のログアウトエンドポイントを直接呼び出します：

```
GET https://{auth0_domain}/v2/logout?client_id={client_id}&returnTo={return_url}
```

### Auth0 設定要件

#### 必要な設定

1. **Application 設定**:

   - Application Type: Single Page Application (SPA)
   - Allowed Callback URLs: `https://app.finsight.com/callback`
   - Allowed Logout URLs: `https://app.finsight.com`
   - Allowed Web Origins: `https://app.finsight.com`

2. **API 設定**:

   - API Identifier: `https://api.finsight.com`
   - Signing Algorithm: RS256
   - Enable RBAC: true
   - Add Permissions: true

3. **接続設定（Connections）**:

   - Google OAuth2: 有効化
   - Database Connection: 無効化（パスワード認証は使用しない）
   - その他のソーシャル: 必要に応じて追加

4. **スコープ/パーミッション**:
   - `openid` - OpenID Connect 認証
   - `profile` - プロフィール情報
   - `email` - メールアドレス
   - `read:transactions` - 取引の読み取り
   - `write:transactions` - 取引の作成・更新
   - `delete:transactions` - 取引の削除
   - `read:budgets` - 予算の読み取り
   - `write:budgets` - 予算の作成・更新
   - `read:reports` - レポートの閲覧
   - `admin:all` - 全権限（管理者用）

### トークン検証

すべての API エンドポイントで以下のトークン検証を実施：

1. **署名検証**: Auth0 の JWKS エンドポイントから公開鍵を取得して検証
2. **Issuer 検証**: `iss`クレームが`https://{auth0_domain}/`と一致
3. **Audience 検証**: `aud`クレームが API Identifier と一致
4. **有効期限検証**: `exp`クレームが現在時刻より未来
5. **スコープ検証**: 必要なスコープが含まれているか確認

---

## 1. ユーザー管理 API

### ユーザープロファイル取得

```
GET /users/me
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "created_at": "2024-01-01T09:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

### ユーザープロファイル更新

```
PUT /users/me
```

**リクエスト:**

```json
{
  "email": "newemail@example.com"
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "newemail@example.com",
  "created_at": "2024-01-01T09:00:00Z",
  "updated_at": "2024-01-20T11:00:00Z"
}
```

### アカウント連携管理

```
GET /users/me/connections
```

**説明**: 連携されているソーシャルアカウントの一覧を取得

**レスポンス:**

```json
{
  "connections": [
    {
      "provider": "google-oauth2",
      "email": "user@gmail.com",
      "connected_at": "2024-01-01T09:00:00Z",
      "profile_image": "https://lh3.googleusercontent.com/..."
    }
  ]
}
```

### アカウント削除

```
DELETE /users/me
```

**説明**: ユーザーアカウントの削除。システム内のデータを削除後、Auth0 のユーザーも削除します。

**リクエスト:**

```json
{
  "confirmation": "DELETE MY ACCOUNT",
  "reason": "利用終了のため"
}
```

**レスポンス:**

```
204 No Content
```

**処理内容:**

1. すべての取引データを匿名化
2. 口座情報を削除
3. 予算・カテゴリデータを削除
4. Auth0 Management API を使用して Auth0 ユーザーを削除

---

## 2. 口座管理 API

### 口座一覧取得

```
GET /accounts
```

**レスポンス:**

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "みずほ銀行",
      "type": "checking",
      "initial_balance": 100000,
      "current_balance": 1200000,
      "created_at": "2024-01-01T09:00:00Z",
      "updated_at": "2024-01-20T10:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "name": "楽天証券",
      "type": "investment",
      "initial_balance": 0,
      "current_balance": 250000,
      "created_at": "2024-01-01T09:00:00Z",
      "updated_at": "2024-01-20T10:00:00Z"
    }
  ],
  "pagination": {
    "total": 2,
    "page": 1,
    "limit": 20,
    "total_pages": 1
  }
}
```

### 口座詳細取得

```
GET /accounts/{accountId}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "みずほ銀行",
  "type": "checking",
  "initial_balance": 100000,
  "current_balance": 1200000,
  "created_at": "2024-01-01T09:00:00Z",
  "updated_at": "2024-01-20T10:00:00Z"
}
```

### 口座作成

```
POST /accounts
```

**リクエスト:**

```json
{
  "name": "三井住友銀行",
  "type": "savings",
  "initial_balance": 500000
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440003",
  "name": "三井住友銀行",
  "type": "savings",
  "initial_balance": 500000,
  "current_balance": 500000,
  "created_at": "2024-01-21T09:00:00Z",
  "updated_at": "2024-01-21T09:00:00Z"
}
```

### 口座更新

```
PUT /accounts/{accountId}
```

**リクエスト:**

```json
{
  "name": "三井住友銀行（定期）"
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440003",
  "name": "三井住友銀行（定期）",
  "type": "savings",
  "initial_balance": 500000,
  "current_balance": 500000,
  "created_at": "2024-01-21T09:00:00Z",
  "updated_at": "2024-01-21T10:00:00Z"
}
```

### 口座削除

```
DELETE /accounts/{accountId}
```

**レスポンス:**

```
204 No Content
```

### 口座残高更新

```
POST /accounts/{accountId}/movements
```

**リクエスト:**

```json
{
  "amount": -50000,
  "note": "投資信託購入"
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440012",
  "account_id": "550e8400-e29b-41d4-a716-446655440001",
  "amount": -50000,
  "occurred_at": "2024-01-21T11:00:00Z",
  "note": "投資信託購入"
}
```

---

## 3. カテゴリ管理 API

### カテゴリマスター一覧取得

```
GET /categories/master
```

**レスポンス:**

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "parent_id": null,
      "name": "食費",
      "sort_order": 1,
      "is_active": true,
      "created_at": "2024-01-01T09:00:00Z",
      "updated_at": "2024-01-01T09:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440021",
      "parent_id": "550e8400-e29b-41d4-a716-446655440020",
      "name": "外食",
      "sort_order": 1,
      "is_active": true,
      "created_at": "2024-01-01T09:00:00Z",
      "updated_at": "2024-01-01T09:00:00Z"
    }
  ]
}
```

### カテゴリマスター作成

```
POST /categories/master
```

**必要な権限:** `admin:all`

**リクエスト:**

```json
{
  "name": "娯楽費",
  "parent_id": null,
  "sort_order": 5
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440022",
  "parent_id": null,
  "name": "娯楽費",
  "sort_order": 5,
  "is_active": true,
  "created_at": "2024-01-21T09:00:00Z",
  "updated_at": "2024-01-21T09:00:00Z"
}
```

### カテゴリマスター更新

```
PUT /categories/master/{masterCategoryId}
```

**必要な権限:** `admin:all`

**リクエスト:**

```json
{
  "name": "娯楽・レジャー費",
  "sort_order": 6
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440022",
  "parent_id": null,
  "name": "娯楽・レジャー費",
  "sort_order": 6,
  "is_active": true,
  "created_at": "2024-01-21T09:00:00Z",
  "updated_at": "2024-01-21T10:00:00Z"
}
```

### カテゴリマスター削除（非活性化）

```
DELETE /categories/master/{masterCategoryId}
```

**必要な権限:** `admin:all`

**説明:** カテゴリマスターを物理削除せず、非活性化します。既にユーザーが使用しているカテゴリに影響を与えません。

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440022",
  "is_active": false,
  "updated_at": "2024-01-21T11:00:00Z"
}
```

### ユーザーカテゴリ一覧取得

```
GET /categories
```

**レスポンス:**

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440030",
      "master_id": "550e8400-e29b-41d4-a716-446655440020",
      "parent_id": null,
      "name": "食費",
      "is_custom": false,
      "sort_order": 1,
      "is_active": true
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440031",
      "master_id": null,
      "parent_id": null,
      "name": "趣味",
      "is_custom": true,
      "sort_order": 10,
      "is_active": true
    }
  ]
}
```

### カテゴリ作成

```
POST /categories
```

**リクエスト:**

```json
{
  "name": "医療費",
  "parent_id": null,
  "sort_order": 11
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440032",
  "master_id": null,
  "parent_id": null,
  "name": "医療費",
  "is_custom": true,
  "sort_order": 11,
  "is_active": true,
  "created_at": "2024-01-21T09:00:00Z"
}
```

### カテゴリ更新

```
PUT /categories/{categoryId}
```

**リクエスト:**

```json
{
  "name": "医療・健康",
  "sort_order": 12
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440032",
  "master_id": null,
  "parent_id": null,
  "name": "医療・健康",
  "is_custom": true,
  "sort_order": 12,
  "is_active": true,
  "updated_at": "2024-01-21T10:00:00Z"
}
```

### カテゴリ削除（非活性化）

```
DELETE /categories/{categoryId}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440032",
  "is_active": false,
  "updated_at": "2024-01-21T11:00:00Z"
}
```

---

## 4. 取引管理 API

### 取引一覧取得

```
GET /transactions
```

**クエリパラメータ:**

| パラメータ  | 型      | 必須 | 説明           |
| ----------- | ------- | ---- | -------------- |
| from        | date    | No   | 開始日         |
| to          | date    | No   | 終了日         |
| type        | string  | No   | income/expense |
| category_id | uuid    | No   | カテゴリ ID    |
| account_id  | uuid    | No   | 口座 ID        |
| min_amount  | integer | No   | 最小金額       |
| max_amount  | integer | No   | 最大金額       |

**レスポンス:**

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440040",
      "type": "expense",
      "amount": 1200,
      "transaction_date": "2024-01-20",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440030",
        "name": "食費"
      },
      "description": "ランチ",
      "created_at": "2024-01-20T12:30:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440041",
      "type": "income",
      "amount": 250000,
      "transaction_date": "2024-01-25",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440035",
        "name": "給与"
      },
      "description": "1月分給与",
      "created_at": "2024-01-25T09:00:00Z"
    }
  ],
  "pagination": {
    "total": 150,
    "page": 1,
    "limit": 20,
    "total_pages": 8
  },
  "summary": {
    "total_income": 250000,
    "total_expense": 106672,
    "net": 143328
  }
}
```

### 取引作成

```
POST /transactions
```

**リクエスト:**

```json
{
  "type": "expense",
  "amount": 3500,
  "transaction_date": "2024-01-21",
  "category_id": "550e8400-e29b-41d4-a716-446655440030",
  "description": "スーパーで買い物"
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440042",
  "type": "expense",
  "amount": 3500,
  "transaction_date": "2024-01-21",
  "category": {
    "id": "550e8400-e29b-41d4-a716-446655440030",
    "name": "食費"
  },
  "description": "スーパーで買い物",
  "created_at": "2024-01-21T15:00:00Z"
}
```

### 取引更新

```
PUT /transactions/{transactionId}
```

**リクエスト:**

```json
{
  "amount": 3800,
  "description": "スーパーで買い物（修正）"
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440042",
  "type": "expense",
  "amount": 3800,
  "transaction_date": "2024-01-21",
  "category": {
    "id": "550e8400-e29b-41d4-a716-446655440030",
    "name": "食費"
  },
  "description": "スーパーで買い物（修正）",
  "updated_at": "2024-01-21T16:00:00Z"
}
```

### 取引削除

```
DELETE /transactions/{transactionId}
```

**レスポンス:**

```
204 No Content
```

### 取引一括作成

```
POST /transactions/bulk
```

**リクエスト:**

```json
{
  "transactions": [
    {
      "type": "expense",
      "amount": 1000,
      "transaction_date": "2024-01-21",
      "category_id": "550e8400-e29b-41d4-a716-446655440030",
      "description": "朝食"
    },
    {
      "type": "expense",
      "amount": 1200,
      "transaction_date": "2024-01-21",
      "category_id": "550e8400-e29b-41d4-a716-446655440030",
      "description": "昼食"
    }
  ]
}
```

**レスポンス:**

```json
{
  "created": 2,
  "transactions": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440043",
      "type": "expense",
      "amount": 1000,
      "transaction_date": "2024-01-21",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440030",
        "name": "食費"
      },
      "description": "朝食"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440044",
      "type": "expense",
      "amount": 1200,
      "transaction_date": "2024-01-21",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440030",
        "name": "食費"
      },
      "description": "昼食"
    }
  ]
}
```

---

## 5. 定期取引 API

### 定期取引一覧取得

```
GET /recurring-transactions
```

**レスポンス:**

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440050",
      "name": "家賃",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440033",
        "name": "住居費"
      },
      "amount": 90000,
      "execution_day": 25,
      "last_executed_date": "2024-01-25",
      "next_execution_date": "2024-02-25",
      "is_active": true,
      "description": "毎月25日に家賃支払い",
      "created_at": "2024-01-01T09:00:00Z",
      "updated_at": "2024-01-25T10:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440051",
      "name": "光熱費",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440034",
        "name": "光熱費"
      },
      "amount": 15000,
      "execution_day": 31,
      "last_executed_date": null,
      "next_execution_date": "2024-02-29",
      "is_active": true,
      "description": "毎月末に光熱費支払い",
      "created_at": "2024-01-21T10:00:00Z",
      "updated_at": "2024-01-21T10:00:00Z"
    }
  ],
  "pagination": {
    "total": 5,
    "page": 1,
    "limit": 20,
    "total_pages": 1
  }
}
```

### 定期取引作成

```
POST /recurring-transactions
```

**リクエスト:**

```json
{
  "name": "電話代",
  "category_id": "550e8400-e29b-41d4-a716-446655440034",
  "amount": 8000,
  "execution_day": 15,
  "description": "毎月15日に携帯電話料金支払い"
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440052",
  "name": "電話代",
  "category": {
    "id": "550e8400-e29b-41d4-a716-446655440034",
    "name": "光熱費"
  },
  "amount": 8000,
  "execution_day": 15,
  "last_executed_date": null,
  "next_execution_date": "2024-02-15",
  "is_active": true,
  "description": "毎月15日に携帯電話料金支払い",
  "created_at": "2024-01-21T10:00:00Z",
  "updated_at": "2024-01-21T10:00:00Z"
}
```

### 定期取引更新

```
PUT /recurring-transactions/{recurringTransactionId}
```

**リクエスト:**

```json
{
  "name": "家賃（更新）",
  "amount": 95000,
  "execution_day": 27,
  "description": "毎月27日に家賃支払い（値上げ後）"
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440050",
  "name": "家賃（更新）",
  "category": {
    "id": "550e8400-e29b-41d4-a716-446655440033",
    "name": "住居費"
  },
  "amount": 95000,
  "execution_day": 27,
  "last_executed_date": "2024-01-25",
  "next_execution_date": "2024-02-27",
  "is_active": true,
  "description": "毎月27日に家賃支払い（値上げ後）",
  "updated_at": "2024-01-21T11:00:00Z"
}
```

### 定期取引無効化

```
DELETE /recurring-transactions/{recurringTransactionId}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440051",
  "is_active": false,
  "updated_at": "2024-01-21T11:00:00Z"
}
```


### 月次定期取引の一括実行

```
POST /recurring-transactions/execute-monthly
```

**説明:** 指定した月の全ての定期取引を一括実行します。バッチ処理での使用を想定。

**リクエスト:**

```json
{
  "target_year": 2024,
  "target_month": 2,
  "dry_run": false
}
```

**レスポンス:**

```json
{
  "execution_summary": {
    "target_month": "2024-02",
    "total_recurring_transactions": 5,
    "executed": 3,
    "skipped": 2,
    "failed": 0
  },
  "executed_transactions": [
    {
      "recurring_transaction_id": "550e8400-e29b-41d4-a716-446655440050",
      "transaction_id": "550e8400-e29b-41d4-a716-446655440062",
      "name": "家賃",
      "execution_date": "2024-02-25",
      "amount": 90000
    },
    {
      "recurring_transaction_id": "550e8400-e29b-41d4-a716-446655440052",
      "transaction_id": "550e8400-e29b-41d4-a716-446655440063",
      "name": "電話代",
      "execution_date": "2024-02-15",
      "amount": 8000
    },
    {
      "recurring_transaction_id": "550e8400-e29b-41d4-a716-446655440051",
      "transaction_id": "550e8400-e29b-41d4-a716-446655440064",
      "name": "光熱費",
      "execution_date": "2024-02-29",
      "amount": 15000
    }
  ],
  "skipped_transactions": [
    {
      "recurring_transaction_id": "550e8400-e29b-41d4-a716-446655440053",
      "name": "保険料",
      "reason": "already_executed_this_month"
    },
    {
      "recurring_transaction_id": "550e8400-e29b-41d4-a716-446655440054",
      "name": "ジム会費",
      "reason": "inactive"
    }
  ]
}
```

### 今月の定期取引予定取得

```
GET /recurring-transactions/schedule
```

**クエリパラメータ:**

| パラメータ | 型      | 必須 | 説明 |
| ---------- | ------- | ---- | ---- |
| year       | integer | No   | 年   |
| month      | integer | No   | 月   |

**レスポンス:**

```json
{
  "month": "2024-02",
  "schedule": [
    {
      "date": "2024-02-15",
      "transactions": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440052",
          "name": "電話代",
          "amount": 8000,
          "category": "光熱費",
          "status": "pending"
        }
      ]
    },
    {
      "date": "2024-02-25",
      "transactions": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440050",
          "name": "家賃",
          "amount": 90000,
          "category": "住居費",
          "status": "pending"
        }
      ]
    },
    {
      "date": "2024-02-29",
      "transactions": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440051",
          "name": "光熱費",
          "amount": 15000,
          "category": "光熱費",
          "status": "pending"
        }
      ]
    }
  ],
  "summary": {
    "total_amount": 113000,
    "total_transactions": 3,
    "pending": 3,
    "executed": 0
  }
}
```

---

## 6. 予算管理 API

### 予算一覧取得

```
GET /budgets
```

**クエリパラメータ:**

| パラメータ | 型      | 必須 | 説明 |
| ---------- | ------- | ---- | ---- |
| year       | integer | No   | 年   |
| month      | integer | No   | 月   |

**レスポンス:**

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440070",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440030",
        "name": "食費"
      },
      "amount": 50000,
      "month": "2024-01",
      "spent": 33368,
      "remaining": 16632,
      "percentage": 66.7,
      "created_at": "2024-01-01T09:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440071",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440033",
        "name": "住居費"
      },
      "amount": 90000,
      "month": "2024-01",
      "spent": 90000,
      "remaining": 0,
      "percentage": 100.0,
      "created_at": "2024-01-01T09:00:00Z"
    }
  ],
  "summary": {
    "total_budget": 200000,
    "total_spent": 106672,
    "total_remaining": 93328,
    "percentage": 53.3
  }
}
```

### 予算設定

```
POST /budgets
```

**リクエスト:**

```json
{
  "budgets": [
    {
      "category_id": "550e8400-e29b-41d4-a716-446655440030",
      "amount": 50000,
      "month": "2024-02"
    },
    {
      "category_id": "550e8400-e29b-41d4-a716-446655440033",
      "amount": 90000,
      "month": "2024-02"
    }
  ]
}
```

**レスポンス:**

```json
{
  "created": 2,
  "budgets": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440072",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440030",
        "name": "食費"
      },
      "amount": 50000,
      "month": "2024-02"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440073",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440033",
        "name": "住居費"
      },
      "amount": 90000,
      "month": "2024-02"
    }
  ]
}
```

### 予算更新

```
PUT /budgets/{budgetId}
```

**リクエスト:**

```json
{
  "amount": 45000
}
```

### 予算削除

```
DELETE /budgets/{budgetId}
```

**レスポンス:**

```
204 No Content
```

### 今月の予算取得

```
GET /budgets/current
```

**説明:** 現在の月の予算情報を取得します。月とカテゴリの設定がない場合は空の予算情報を返します。

**レスポンス:**

```json
{
  "month": "2024-02",
  "current_date": "2024-02-15",
  "days_in_month": 29,
  "days_elapsed": 15,
  "days_remaining": 14,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440072",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440030",
        "name": "食費"
      },
      "amount": 50000,
      "spent": 22150,
      "remaining": 27850,
      "percentage": 44.3,
      "daily_average_budget": 1724,
      "daily_average_spent": 1477,
      "projected_month_end": 42630,
      "status": "on_track",
      "created_at": "2024-02-01T09:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440073",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440033",
        "name": "住居費"
      },
      "amount": 90000,
      "spent": 90000,
      "remaining": 0,
      "percentage": 100.0,
      "daily_average_budget": 3103,
      "daily_average_spent": 6000,
      "projected_month_end": 90000,
      "status": "completed",
      "created_at": "2024-02-01T09:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440074",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440031",
        "name": "交通費"
      },
      "amount": 10000,
      "spent": 12500,
      "remaining": -2500,
      "percentage": 125.0,
      "daily_average_budget": 345,
      "daily_average_spent": 833,
      "projected_month_end": 24150,
      "status": "over_budget",
      "created_at": "2024-02-01T09:00:00Z"
    }
  ],
  "summary": {
    "total_budget": 200000,
    "total_spent": 139850,
    "total_remaining": 60150,
    "percentage": 69.9,
    "daily_average_budget": 6897,
    "daily_average_spent": 9323,
    "projected_month_end": 270367,
    "status": "over_projected"
  },
  "alerts": [
    {
      "type": "over_budget",
      "category": "交通費",
      "message": "交通費が予算を25%超過しています"
    },
    {
      "type": "projection_warning",
      "message": "現在のペースでは月末に予算を35%超過する見込みです"
    }
  ],
  "recommendations": [
    {
      "type": "spending_adjustment",
      "category": "食費",
      "message": "食費は順調です。このペースを維持してください"
    },
    {
      "type": "budget_review",
      "category": "交通費",
      "message": "交通費の予算見直しを検討してください"
    }
  ]
}
```

**ステータスの説明:**
- `on_track`: 予算内で順調
- `completed`: 定期支払いなどで予算消化完了
- `over_budget`: 予算超過
- `under_spent`: 支出が少なすぎる（予算の30%以下）
- `at_risk`: 月末超過の可能性あり（現在のペースで110%以上）
- `over_projected`: 全体として月末予算超過の見込み


---

## 7. AI 予算提案 API

### 予算提案取得

```
GET /budget-suggestions
```

**クエリパラメータ:**

| パラメータ | 型     | 必須 | 説明     |
| ---------- | ------ | ---- | -------- |
| month      | string | Yes  | 対象年月 |

**レスポンス:**

```json
{
  "month": "2024-02",
  "suggestions": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440080",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440030",
        "name": "食費"
      },
      "suggested_amount": 45000,
      "current_budget": 50000,
      "last_month_actual": 33368,
      "three_month_average": 38500,
      "reason": "過去3ヶ月の実績に基づき、5,000円の削減が可能です",
      "confidence": 0.85,
      "status": "pending"
    }
  ],
  "summary": {
    "total_suggested": 185000,
    "total_current": 200000,
    "potential_savings": 15000
  }
}
```

### 予算提案生成

```
POST /budget-suggestions/generate
```

**リクエスト:**

```json
{
  "month": "2024-02",
  "optimization_type": "balanced"
}
```

**レスポンス:**

```json
{
  "month": "2024-02",
  "suggestions": [...],
  "created_at": "2024-01-21T10:00:00Z"
}
```

### 予算提案承認

```
PUT /budget-suggestions/{suggestionId}/accept
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440080",
  "status": "accepted",
  "budget_id": "550e8400-e29b-41d4-a716-446655440074",
  "updated_at": "2024-01-21T11:00:00Z"
}
```

### 予算提案却下

```
PUT /budget-suggestions/{suggestionId}/reject
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440080",
  "status": "rejected",
  "updated_at": "2024-01-21T11:00:00Z"
}
```

---

## 8. 資産管理 API

### 資産スナップショット取得

```
GET /assets/snapshots
```

**クエリパラメータ:**

| パラメータ | 型   | 必須 | 説明   |
| ---------- | ---- | ---- | ------ |
| from       | date | No   | 開始日 |
| to         | date | No   | 終了日 |

**レスポンス:**

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440090",
      "snapshot_date": "2024-01-21",
      "total_assets": 1500000,
      "change_from_previous": 50000,
      "change_percentage": 3.4,
      "breakdown": [
        {
          "account_id": "550e8400-e29b-41d4-a716-446655440001",
          "account_name": "みずほ銀行",
          "balance": 1200000
        },
        {
          "account_id": "550e8400-e29b-41d4-a716-446655440002",
          "account_name": "楽天証券",
          "balance": 250000
        },
        {
          "account_id": "550e8400-e29b-41d4-a716-446655440003",
          "account_name": "SBI証券",
          "balance": 50000
        }
      ],
      "created_at": "2024-01-21T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 30,
    "page": 1,
    "limit": 20,
    "total_pages": 2
  }
}
```

### 資産スナップショット作成

```
POST /assets/snapshots
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440091",
  "snapshot_date": "2024-01-22",
  "total_assets": 1550000,
  "created_at": "2024-01-22T00:00:00Z"
}
```

---

## 9. 資産予測 API

### 資産予測取得

```
GET /assets/forecasts/latest
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440100",
  "horizon_months": 12,
  "forecast_date": "2025-01-21",
  "current_assets": 1500000,
  "predicted_assets": 1680000,
  "growth_rate": 12.0,
  "method": "linear_regression",
  "confidence_interval": {
    "lower": 1600000,
    "upper": 1760000
  },
  "assumptions": {
    "monthly_income": 250000,
    "monthly_expense": 200000,
    "savings_rate": 0.2,
    "investment_return": 0.05,
    "inflation_rate": 0.02
  },
  "monthly_forecast": [
    {
      "month": "2024-02",
      "predicted_assets": 1550000
    },
    {
      "month": "2024-03",
      "predicted_assets": 1600000
    }
  ],
  "created_at": "2024-01-21T10:00:00Z"
}
```

### 資産予測作成

```
POST /assets/forecasts
```

**リクエスト:**

```json
{
  "horizon_months": 12,
  "scenario": "realistic",
  "assumptions": {
    "monthly_income": 250000,
    "monthly_expense": 200000,
    "major_expenses": [
      {
        "month": "2024-06",
        "amount": 300000,
        "description": "旅行"
      }
    ]
  }
}
```

**レスポンス:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440101",
  "horizon_months": 12,
  "forecast_date": "2025-01-21",
  "predicted_assets": 1680000,
  "method": "linear_regression",
  "created_at": "2024-01-21T11:00:00Z"
}
```

---

## 10. レポート API

### 月次レポート取得

```
GET /reports/monthly
```

**クエリパラメータ:**

| パラメータ | 型      | 必須 | 説明 |
| ---------- | ------- | ---- | ---- |
| year       | integer | Yes  | 年   |
| month      | integer | Yes  | 月   |

**レスポンス:**

```json
{
  "period": {
    "year": 2024,
    "month": 1
  },
  "summary": {
    "total_income": 250000,
    "total_expense": 106672,
    "net_income": 143328,
    "savings_rate": 57.3
  },
  "top_categories": [
    {
      "category": "食費",
      "amount": 33368,
      "percentage": 31.3
    },
    {
      "category": "住居費",
      "amount": 24964,
      "percentage": 23.4
    }
  ],
  "daily_trend": [...],
  "comparison": {
    "previous_month": {
      "total_expense": 94222,
      "change_amount": 12450,
      "change_percentage": 13.2
    },
    "same_month_last_year": {
      "total_expense": 98500,
      "change_amount": 8172,
      "change_percentage": 8.3
    }
  }
}
```

### 年次レポート取得

```
GET /reports/yearly
```

**クエリパラメータ:**

| パラメータ | 型      | 必須 | 説明 |
| ---------- | ------- | ---- | ---- |
| year       | integer | Yes  | 年   |

**レスポンス:**

```json
{
  "year": 2024,
  "summary": {
    "total_income": 3000000,
    "total_expense": 2400000,
    "net_income": 600000,
    "average_monthly_expense": 200000,
    "savings_rate": 20.0
  },
  "monthly_breakdown": [...],
  "category_ranking": [...],
  "asset_growth": {
    "beginning_balance": 1000000,
    "ending_balance": 1600000,
    "growth_amount": 600000,
    "growth_rate": 60.0
  }
}
```

### レポートエクスポート

```
POST /reports/export
```

**リクエスト:**

```json
{
  "type": "monthly",
  "year": 2024,
  "month": 1,
  "format": "pdf",
  "include_charts": true
}
```

**レスポンス:**

```json
{
  "download_url": "https://api.finsight.com/downloads/reports/550e8400-e29b-41d4-a716-446655440110.pdf",
  "expires_at": "2024-01-22T10:00:00Z"
}
```

### レポートメール送信

```
POST /reports/email
```

**リクエスト:**

```json
{
  "type": "monthly",
  "year": 2024,
  "month": 1,
  "recipient_email": "user@example.com"
}
```

**レスポンス:**

```json
{
  "message": "レポートを送信しました",
  "sent_at": "2024-01-21T10:00:00Z"
}
```

---

## 11. 通知設定 API

### 通知設定取得

```
GET /notifications/settings
```

**レスポンス:**

```json
{
  "monthly_report": {
    "enabled": true,
    "send_day": 1,
    "send_time": "09:00",
    "email": "user@gmail.com"
  },
  "budget_exceeded_alert": {
    "enabled": true,
    "email": "user@gmail.com"
  }
}
```

### 通知設定更新

```
PUT /notifications/settings
```

**リクエスト:**

```json
{
  "monthly_report": {
    "enabled": true,
    "send_day": 5,
    "send_time": "08:00"
  },
  "budget_exceeded_alert": {
    "enabled": false
  }
}
```

**レスポンス:**

```json
{
  "message": "通知設定を更新しました",
  "updated_at": "2024-01-21T10:00:00Z"
}
```

### 月次レポートメール送信（手動）

```
POST /notifications/send-monthly-report
```

**説明:** 指定した月の月次レポートを即座にメール送信します。

**リクエスト:**

```json
{
  "year": 2024,
  "month": 1
}
```

**レスポンス:**

```json
{
  "message": "月次レポートを送信しました",
  "sent_to": "user@gmail.com",
  "report_period": "2024-01",
  "sent_at": "2024-01-21T10:00:00Z"
}
```

### 予算超過通知送信（システム用）

```
POST /notifications/send-budget-alert
```

**説明:** 予算超過が検出された際にシステムが自動的に呼び出すエンドポイント。

**リクエスト:**

```json
{
  "categories": [
    {
      "category_id": "550e8400-e29b-41d4-a716-446655440031",
      "category_name": "交通費",
      "budget_amount": 10000,
      "spent_amount": 12500,
      "exceeded_amount": 2500,
      "exceeded_date": "2024-02-15"
    }
  ],
  "month": "2024-02"
}
```

**レスポンス:**

```json
{
  "message": "予算超過通知を送信しました",
  "sent_to": "user@gmail.com",
  "categories_alerted": 1,
  "sent_at": "2024-02-16T09:00:00Z"
}
```

---

## エラーコード一覧

| コード                  | HTTP ステータス | 説明                                 |
| ----------------------- | --------------- | ------------------------------------ |
| UNAUTHORIZED            | 401             | 認証が必要です                       |
| FORBIDDEN               | 403             | アクセス権限がありません             |
| NOT_FOUND               | 404             | リソースが見つかりません             |
| VALIDATION_ERROR        | 400             | 入力値が不正です                     |
| DUPLICATE_ENTRY         | 409             | 重複するエントリーが存在します       |
| BUDGET_EXCEEDED         | 400             | 予算を超過しています                 |
| INSUFFICIENT_BALANCE    | 400             | 残高が不足しています                 |
| CATEGORY_IN_USE         | 400             | カテゴリが使用中のため削除できません |
| INVALID_DATE_RANGE      | 400             | 日付範囲が不正です                   |
| FILE_TOO_LARGE          | 413             | ファイルサイズが大きすぎます         |
| UNSUPPORTED_FILE_FORMAT | 415             | サポートされていないファイル形式です |
| RATE_LIMIT_EXCEEDED     | 429             | レート制限を超えました               |
| INTERNAL_SERVER_ERROR   | 500             | サーバーエラーが発生しました         |
