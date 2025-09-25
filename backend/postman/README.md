# Postmanを使用したAPIデバッグガイド

## セットアップ手順

### 1. Postmanのインストール
[Postman公式サイト](https://www.postman.com/downloads/)からダウンロード・インストール

### 2. 環境とコレクションのインポート

1. Postmanを開く
2. 左サイドバーの「Import」をクリック
3. 以下のファイルをドラッグ&ドロップ：
   - `FinanceTracker.postman_environment.json`
   - `FinanceTracker.postman_collection.json`

### 3. 環境変数の設定

1. 右上の環境選択から「FinanceTracker Local」を選択
2. 環境変数の確認（デフォルトで設定済み）：
   - `auth0_domain`: `dev-kimushun3765.jp.auth0.com`
   - `auth0_client_id`: `N19XPSt4tzz3Fa7reXz26yRgLfoeuoJ9`
   - `auth0_audience`: `https://api.financetracker.local`
   - `use_dev_auth`: `true`（開発用認証を有効化）
   - `dev_user_id`: `auth0|dev-test-user-123`

### 開発用認証の切り替え

- **開発用認証を使用する場合**（デフォルト）:
  - `use_dev_auth`: `true`に設定
  - 自動的に`X-Dev-User-ID`ヘッダーが追加されます
  
- **実際のAuth0認証を使用する場合**:
  - `use_dev_auth`: `false`に設定
  - `access_token`に有効なJWTトークンを設定

### 4. バックエンドサーバーの起動

```bash
cd backend
make api
```

## デバッグフロー

### 1. ヘルスチェック
まず「Health Check」リクエストを送信してAPIが起動していることを確認

### 2. 認証方法

#### 方法1: 開発用ヘッダー（推奨 - 開発環境）
デフォルトで有効になっています。追加設定は不要です。
- 自動的に`X-Dev-User-ID: auth0|dev-test-user-123`ヘッダーが全リクエストに追加
- 環境変数`use_dev_auth`が`true`に設定されている間有効

#### 方法2: Machine-to-Machine認証（本番環境向け）
1. Auth0ダッシュボードでM2Mアプリケーションを作成
2. 「Get Auth0 Token (Client Credentials)」のbodyを編集：
   - `YOUR_M2M_CLIENT_ID`を実際のClient IDに置換
   - `YOUR_M2M_CLIENT_SECRET`を実際のClient Secretに置換
3. リクエストを実行してトークンを取得

#### 方法3: SPAからの認証フロー
1. ブラウザで `http://localhost:8080/auth/login` を開く
2. Auth0でログイン
3. コールバック後、クッキーにトークンが保存される

#### なぜPassword Grant Typeが使えないのか？
- Auth0のSPAアプリケーションでは、セキュリティ上の理由からPassword Grant Typeは無効化されています
- 代わりにAuthorization Code Flow with PKCE（方法3）を使用してください
- 開発環境では方法1の開発用ヘッダーが最も簡単です

#### 方法4: 開発用ヘッダーの詳細（既に有効）
開発環境では`X-Dev-User-ID`ヘッダーで認証をバイパスできます：

```bash
# コマンドラインでのテスト
curl -H "X-Dev-User-ID: auth0|dev-test-123" http://localhost:8080/api/v1/users/me

# または開発用テストスクリプトを実行
./scripts/test-api-dev.sh
```

Postmanでは：
1. Headers タブで `X-Dev-User-ID: auth0|dev-test-123` を追加
2. Authorization を「No Auth」に設定

### 3. APIのテスト順序

推奨される実行順序：

1. **Users**
   - Get Current User（ユーザーが存在しない場合は404）
   - Update Current User（ユーザーが自動作成される）

2. **Categories**  
   - List Master Categories（マスターカテゴリ一覧を取得）
   - Create Category（マスターIDを使用してユーザーカテゴリ作成）
   - List User Categories

3. **Accounts**
   - Create Account（口座作成）
   - List Accounts
   - Get/Update/Delete Account

4. **Transactions**
   - Create Transaction（口座IDとカテゴリIDが必要）
   - List Transactions
   - Get Monthly Summary

5. **Budgets**
   - Create Budget（カテゴリIDが必要）
   - Get Current Budgets
   - List/Update/Delete Budget

## デバッグのコツ

### 1. リクエストIDの確認
すべてのレスポンスヘッダーに`X-Request-ID`が含まれます。エラー時はこのIDでログを検索できます。

### 2. エラーレスポンスの確認
- 400: バリデーションエラー（リクエストボディを確認）
- 401: 認証エラー（トークンを確認）
- 403: 権限エラー（他ユーザーのリソースへのアクセス）
- 404: リソースが見つからない
- 500: サーバーエラー（ログを確認）

### 3. UUIDの管理
レスポンスから取得したUUIDは、Postman変数に保存すると便利：

```javascript
// Testsタブに追加
if (pm.response.code === 201) {
    const response = pm.response.json();
    pm.environment.set('last_account_id', response.id);
}
```

### 4. データの依存関係
- Transaction作成には有効なAccount IDとCategory IDが必要
- Budget作成には有効なCategory IDが必要
- Category作成には有効なCategory Master IDが必要

### 5. 日付フォーマット
- ISO 8601形式を使用: `2024-01-15T12:00:00Z`
- 日付のみの場合: `2024-01-15`

## トラブルシューティング

### 認証エラーが発生する場合
1. トークンの有効期限を確認
2. Auth0の設定（Audience、Scope）を確認
3. バックエンドの環境変数を確認

### CORSエラーが発生する場合
フロントエンドからの呼び出し時は、`config/config.go`でCORS設定を確認

### データベースエラーが発生する場合
1. PostgreSQLが起動していることを確認
2. マイグレーションが実行されていることを確認
3. シードデータが投入されていることを確認

```bash
make migrate
make seed
```

## ログの確認方法

```bash
# Dockerコンテナのログ
docker-compose logs -f api

# 構造化ログの検索（リクエストID）
docker-compose logs api | grep "request-id-here"
```