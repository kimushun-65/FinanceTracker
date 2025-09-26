# Postman クイックスタートガイド

## 🚀 3ステップで開始

### 1. ファイルをインポート
Postmanを開いて以下のファイルをドラッグ&ドロップ：
- `FinanceTracker.postman_collection.json`
- `FinanceTracker.postman_environment.json`

### 2. 環境を選択
右上から「FinanceTracker Local」を選択

### 3. APIをテスト
任意のリクエストを送信！

## ✨ 特徴

### 自動的に開発用ヘッダーが追加されます
デフォルトで`X-Dev-User-ID: auth0|dev-test-user-123`が全リクエストに追加されるため、手動設定は不要です。

### 環境変数で制御可能
- `use_dev_auth`: `true`（デフォルト）- 開発用認証を使用
- `use_dev_auth`: `false` - 実際のAuth0トークンを使用

## 📝 テスト順序

1. **Health Check** - APIの稼働確認
2. **Users > Update Current User** - ユーザー作成（初回のみ）
3. **Categories > List Master Categories** - カテゴリマスター取得
4. **Accounts > Create Account** - 口座作成
5. その他のエンドポイント

## 🔧 トラブルシューティング

### 500エラーが返される場合
```bash
# バックエンドを再起動
docker-compose restart backend

# ログを確認
docker-compose logs -f backend
```

### 認証エラーが返される場合
環境変数の`use_dev_auth`が`true`になっているか確認

## 📚 詳細情報
詳しい使い方は[README.md](./README.md)を参照してください。