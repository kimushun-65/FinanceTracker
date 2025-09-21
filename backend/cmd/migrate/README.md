# データベースマイグレーションガイド

## 概要
このプロジェクトでは、本番環境のデータベースマイグレーションには **Atlas** を使用し、開発の利便性のために **GORM AutoMigrate** を使用しています。

## マイグレーション戦略

### 本番環境・ステージング環境（Atlas）
Atlasは本番環境における主要なマイグレーションツールです：

```bash
# マイグレーションの適用
make migrate

# または直接実行
go run cmd/migrate/main.go apply

# マイグレーションステータスの確認
go run cmd/migrate/main.go check

# 新しいマイグレーションの生成
go run cmd/migrate/main.go diff マイグレーション名

# マイグレーションの検証
go run cmd/migrate/main.go validate
```

### 開発環境のみ（GORM AutoMigrate）
素早い開発セットアップのために、GORM AutoMigrateが利用可能です：

```bash
# 開発環境の初期セットアップ（開発環境でのみ実行可能）
go run cmd/migrate/main.go dev-setup
```

**⚠️ 警告**: `dev-setup`コマンドについて：
- `APP_ENV=development`の場合のみ動作します
- 本番環境では絶対に使用しないでください
- トリガーは管理しません（Atlasが管理）
- プロトタイピング専用です

## ベストプラクティス

1. **スキーマ変更**: 必ずAtlasマイグレーションを使用
   - `schema.hcl`を編集して理想の状態を定義
   - `make atlas-diff name=マイグレーション名`を実行してマイグレーションを生成
   - 生成されたSQLファイルをレビュー
   - マイグレーションファイルと更新された`atlas.sum`の両方をコミット

2. **開発ワークフロー**:
   - 初期のデータベースセットアップには`dev-setup`を使用
   - スキーマ変更にはAtlasマイグレーションを使用
   - スキーマバージョン管理にGORM AutoMigrateを頼らない

3. **本番環境へのデプロイ**:
   - 必ず`make migrate`またはAtlas applyを使用
   - 本番環境では`dev-setup`を使用しない
   - ステージング環境で必ずマイグレーションをテスト

## マイグレーションファイル

- `schema.hcl`: 理想のデータベーススキーマ（信頼できる情報源）
- `migrations/`: Atlasマイグレーションファイル
- `atlas.sum`: マイグレーション整合性のためのチェックサムファイル

## 環境変数

- `APP_ENV`: 利用可能なコマンドを制御（developmentでdev-setupが有効）
- `DATABASE_URL`: データベース接続文字列
- `atlas.hcl`内のAtlas設定

## トラブルシューティング

### "trigger_set_timestamp() does not exist"エラー
トリガーが参照されているが関数が存在しない場合に発生します。
解決策：初期マイグレーションにトリガー関数定義が含まれていることを確認してください。

### "connected database is not clean"エラー
既存のデータベースにマイグレーションを適用しようとした際に発生します。
解決策：`--baseline`フラグを使用するか、データベースを最初にクリーンアップしてください。

### チェックサムエラー
`atlas migrate hash --dir file://migrations`を実行してチェックサムを再生成してください。