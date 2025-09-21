# データベースマイグレーション管理

このディレクトリはFinanceTrackerのデータベースマイグレーションを管理します。

## マイグレーション戦略

### 開発環境・本番環境共通
- **Atlas Migration**: 正式なマイグレーション管理
- コマンド: `make migrate-apply`
- 利点: バージョン管理、ロールバック可能、環境間の一貫性

### 緊急時のオプション
- **GORM AutoMigration**: デバッグや緊急時のみ
- コマンド: `make migrate-gorm`
- 注意: 本番環境では使用しない

## ディレクトリ構造

```
cmd/migrate/
├── main.go          # マイグレーション実行エントリポイント
├── schema.hcl       # Atlasスキーマ定義
├── migrations/      # Atlasマイグレーションファイル
│   └── atlas.sum    # マイグレーション整合性チェックファイル
└── README.md        # このファイル
```

## 使用方法

### 初回セットアップ

1. **データベースとマイグレーション**
```bash
make setup  # DB作成、マイグレーション適用、シードデータ投入
```

2. **マイグレーションファイルの確認**
```bash
make migrate-status  # 適用状態を確認
```

### スキーマ変更時のワークフロー

1. **schema.hclを編集**
```bash
vim cmd/migrate/schema.hcl  # スキーマ定義を変更
```

2. **差分確認**
```bash
make migrate-check  # 現在のDBとの差分を確認
```

3. **新しいマイグレーション作成**
```bash
make atlas-migrate-new  # 対話的に名前を指定
# または
make migrate-diff  # タイムスタンプで自動生成
```

4. **マイグレーション適用**
```bash
make migrate-apply  # 新しいマイグレーションを適用
```

### データベースリセット

```bash
make db-reset      # DB削除→作成→マイグレーション
make seed          # シードデータ投入
```

## スキーマ定義

スキーマは以下の2箇所で管理されています：

1. **GORMモデル** (`internal/infrastructure/gorm/model/`)
   - 開発時のプライマリ定義
   - Go構造体でスキーマを表現

2. **Atlas HCL** (`cmd/migrate/schema.hcl`)
   - 本番用のスキーマ定義
   - HCL形式でスキーマを記述

## トラブルシューティング

### マイグレーションエラー
- `make migrate-validate` でマイグレーションの妥当性を確認
- `make migrate-history` で適用履歴を確認

### スキーマ不整合
- `make migrate-check` で現在のDBとスキーマの差分を確認
- 必要に応じて `make db-reset` でクリーンな状態に戻す

## 注意事項

- 開発環境ではGORM AutoMigrationを使用するため、破壊的変更（カラム削除等）は自動適用されません
- 本番環境への移行前に必ずAtlasマイグレーションファイルを生成してください
- マイグレーションファイルは一度適用されたら編集しないでください