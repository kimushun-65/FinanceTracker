# バックエンド実装タスク一覧

## 概要
FinanceTrackerバックエンドをDocker環境で開発するためのタスクリストです。
AWS Lambda版の実装計画を基に、Gin framework + PostgreSQLで構築します。

## Phase 0: 開発環境構築
### 0.1 Docker環境セットアップ
- [ ] Docker Composeファイル作成
  - Goアプリケーション用コンテナ (Gin)
  - PostgreSQL用コンテナ
  - pgAdmin用コンテナ（DB管理用）
- [ ] Dockerfileの作成（マルチステージビルド）
- [ ] docker-compose.ymlの作成
- [ ] .env.exampleファイル作成（環境変数テンプレート）

### 0.2 Ginプロジェクト初期化
- [ ] backend/ディレクトリ作成
- [ ] go mod init
- [ ] 基本的な依存関係のインストール
  - gin-gonic/gin
  - gorm.io/gorm
  - gorm.io/driver/postgres
  - joho/godotenv
  - go-playground/validator/v10
  - uber-go/zap（ロギング）

### 0.3 データベースセットアップ
- [ ] PostgreSQL初期化スクリプト作成
- [ ] Atlas（ariga.io/atlas）のセットアップ
  - [ ] atlasfileの作成（プロジェクトルート）
  - [ ] スキーマファイル（schema.hcl）の作成（cmd/migrate/schema.hcl）
  - [ ] マイグレーションディレクトリ構成（cmd/migrate/migrations/）
- [ ] 初期スキーマ作成（テーブル定義）
- [ ] cmd/migrate/main.go作成（マイグレーション実行コマンド）
- [ ] cmd/seed/main.go作成（シードデータ投入コマンド）
- [ ] cmd/seed/data/にシードデータファイル配置

### 0.4 開発ツール設定
- [ ] Makefileの作成（ビルド、テスト、実行コマンド）
  - [ ] make api - APIサーバー起動
  - [ ] make migrate - マイグレーション実行
  - [ ] make seed - シードデータ投入
  - [ ] make atlas-* - Atlas関連コマンド
- [ ] air.tomlの作成（ホットリロード設定）
- [ ] .gitignoreの更新
- [ ] atlasfile（Atlas設定ファイル）の作成

## Phase 1: プロジェクト構造構築（オニオンアーキテクチャ）
### 1.1 ディレクトリ構造作成
```
backend/
├── cmd/
│   ├── api/
│   │   └── main.go              # APIサーバーエントリポイント
│   ├── migrate/
│   │   ├── main.go              # マイグレーション実行
│   │   └── migrations/          # マイグレーションファイル
│   └── seed/
│       ├── main.go              # シードデータ投入
│       └── data/                # シードデータファイル
├── internal/
│   ├── domain/                  # ドメイン層
│   │   ├── user/
│   │   ├── account/
│   │   ├── transaction/
│   │   ├── category/
│   │   ├── budget/
│   │   └── common/
│   ├── application/             # アプリケーション層
│   │   ├── dto/
│   │   ├── usecase/
│   │   └── service/
│   ├── infrastructure/          # インフラストラクチャ層
│   │   ├── postgres/
│   │   ├── auth/
│   │   └── external/
│   └── interface/               # インターフェース層
│       ├── http/
│       │   ├── handler/
│       │   ├── middleware/
│       │   └── router/
│       └── validator/
├── pkg/                         # 共有パッケージ
│   ├── config/
│   ├── errors/
│   └── logger/
├── scripts/                     # ユーティリティスクリプト
├── docker/                      # Docker関連ファイル
├── docs/                        # ドキュメント
└── tests/                       # テストファイル
```

### 1.2 基本設定ファイル作成
- [ ] pkg/config/config.go（設定管理）
- [ ] pkg/logger/logger.go（ログ設定）
- [ ] pkg/errors/errors.go（カスタムエラー）

### 1.3 HTTPサーバー基本実装
- [ ] cmd/api/main.go（APIサーバーエントリポイント）
- [ ] internal/interface/http/router/router.go（ルーティング設定）
- [ ] middleware実装（CORS, ロギング, エラーハンドリング）

### 1.4 マイグレーション・シード実装
- [ ] cmd/migrate/main.go（Atlasマイグレーション実行）
- [ ] cmd/migrate/schema.hcl（スキーマ定義）
- [ ] cmd/seed/main.go（シードデータ投入）
- [ ] cmd/seed/data/*.json（シードデータファイル）

## Phase 2: ドメイン層実装
### 2.1 共通ドメイン要素
- [ ] base_entity.go（UUID, timestamps）
- [ ] 値オブジェクト実装
  - [ ] money.go（金額）
  - [ ] email.go（メールアドレス）
  - [ ] hex_color.go（カラーコード）
- [ ] base_repository.go（リポジトリインターフェース）
- [ ] ドメインエラー定義

### 2.2 各ドメインエンティティ実装
- [ ] User エンティティ
  - [ ] entity/user.go
  - [ ] value/user_id.go
  - [ ] value/auth0_id.go
  - [ ] repository/user_repository.go
- [ ] Account エンティティ
  - [ ] entity/account.go
  - [ ] value/account_type.go
  - [ ] value/balance.go
  - [ ] repository/account_repository.go
- [ ] Transaction エンティティ
  - [ ] entity/transaction.go
  - [ ] value/transaction_type.go
  - [ ] repository/transaction_repository.go
- [ ] Category エンティティ
  - [ ] entity/category.go
  - [ ] entity/category_master.go
  - [ ] repository/category_repository.go
- [ ] Budget エンティティ
  - [ ] entity/budget.go
  - [ ] entity/budget_suggestion.go
  - [ ] repository/budget_repository.go

## Phase 3: アプリケーション層実装
### 3.1 DTOの定義
- [ ] user_dto.go
- [ ] account_dto.go
- [ ] transaction_dto.go
- [ ] category_dto.go
- [ ] budget_dto.go
- [ ] report_dto.go

### 3.2 ユースケース実装
- [ ] ユーザー管理
  - [ ] user_service.go
  - [ ] auth_service.go
- [ ] 口座管理
  - [ ] account_service.go
  - [ ] account_movement_service.go
- [ ] 取引管理
  - [ ] transaction_service.go
- [ ] カテゴリ管理
  - [ ] category_service.go
  - [ ] category_master_service.go
- [ ] 予算管理
  - [ ] budget_service.go
  - [ ] budget_suggestion_service.go
- [ ] レポート生成
  - [ ] report_service.go
  - [ ] summary_service.go
  - [ ] dashboard_service.go

### 3.3 トランザクション管理
- [ ] transaction_manager.go

## Phase 4: インフラストラクチャ層実装
### 4.1 データベース接続
- [ ] connection.go（GORM接続管理）
- [ ] transaction.go（トランザクション管理）
- [ ] atlas_migration.go（Atlasマイグレーション実行）

### 4.2 GORMモデル定義
- [ ] user_model.go
- [ ] account_model.go
- [ ] transaction_model.go
- [ ] category_model.go
- [ ] budget_model.go
- [ ] その他関連モデル

### 4.3 リポジトリ実装
- [ ] base_repository.go（共通処理）
- [ ] user_repository.go
- [ ] account_repository.go
- [ ] transaction_repository.go
- [ ] category_repository.go
- [ ] budget_repository.go

### 4.4 外部サービス連携
- [ ] Auth0認証サービス
- [ ] メール送信サービス
- [ ] キャッシュサービス（Redis）

## Phase 5: HTTPインターフェース層実装
### 5.1 ハンドラー実装
- [ ] user_handler.go（GET,PUT /users/me）
- [ ] account_handler.go（CRUD /accounts）
- [ ] transaction_handler.go（CRUD /transactions）
- [ ] category_handler.go（CRUD /categories）
- [ ] budget_handler.go（CRUD /budgets）
- [ ] report_handler.go（GET /reports/*）
- [ ] auth_handler.go（POST /auth/*）

### 5.2 バリデーション
- [ ] リクエストバリデーション実装
- [ ] カスタムバリデーター作成

### 5.3 レスポンス処理
- [ ] 統一レスポンス形式
- [ ] エラーレスポンス処理

## Phase 6: テスト実装
### 6.1 単体テスト
- [ ] ドメイン層テスト
- [ ] アプリケーション層テスト
- [ ] インフラストラクチャ層テスト（モック使用）

### 6.2 統合テスト
- [ ] API統合テスト
- [ ] データベース統合テスト

### 6.3 E2Eテスト
- [ ] 主要フロー のE2Eテスト

## Phase 7: CI/CD・運用設定
### 7.1 GitHub Actions設定
- [ ] テスト自動実行
- [ ] ビルド自動化
- [ ] Dockerイメージビルド

### 7.2 監視・ログ設定
- [ ] 構造化ログ実装
- [ ] メトリクス収集
- [ ] ヘルスチェックエンドポイント

### 7.3 ドキュメント作成
- [ ] API仕様書（OpenAPI/Swagger）
- [ ] 開発者向けREADME
- [ ] 環境構築手順書

## Phase 8: Lambda統合（オプション）
### 8.1 Lambda対応
- [ ] Lambda用エントリポイント作成
- [ ] CDKスタックとの統合
- [ ] ビルドスクリプト更新

### 8.2 デプロイメント
- [ ] Lambda関数デプロイ設定
- [ ] API Gateway設定

## 技術スタック
- **言語**: Go 1.21+
- **Webフレームワーク**: Gin
- **ORM**: GORM v2
- **データベース**: PostgreSQL 15
- **キャッシュ**: Redis
- **認証**: Auth0
- **コンテナ**: Docker/Docker Compose
- **CI/CD**: GitHub Actions

## 開発規約
1. **コミット規約**: Conventional Commits
2. **ブランチ戦略**: Git Flow
3. **コードレビュー**: PR必須
4. **テストカバレッジ**: 80%以上目標

## 注意事項
- Docker環境での開発を前提とする
- ローカル開発ではホットリロードを活用
- 環境変数は.envファイルで管理（.envはgitignore）
- データベースマイグレーションはAtlasで管理
  - スキーマ定義はHCL形式で記述
  - マイグレーションの自動生成と適用
  - スキーマのバージョン管理