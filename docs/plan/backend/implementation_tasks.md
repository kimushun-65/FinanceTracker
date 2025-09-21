# バックエンド実装タスク一覧

## 概要
FinanceTrackerバックエンドをDocker環境で開発するためのタスクリストです。
AWS Lambda版の実装計画を基に、Gin framework + PostgreSQLで構築します。

## 進捗サマリー
- **Phase 0: 開発環境構築** - 100%完了 ✅
- **Phase 1: プロジェクト構造構築** - 100%完了 ✅
- **Phase 2: ドメイン層実装** - 未着手
- **Phase 3: アプリケーション層実装** - 未着手
- **Phase 4: インフラストラクチャ層実装** - 40%完了（GORMモデル定義済み、DB接続未実装）
- **Phase 5: HTTPインターフェース層実装** - 30%完了（ルーティング・ミドルウェア実装済み、ハンドラー未実装）
- **Phase 6: テスト実装** - 未着手
- **Phase 7: CI/CD・運用設定** - 60%完了
- **Phase 8: Lambda統合** - 未着手

## Phase 0: 開発環境構築
### 0.1 Docker環境セットアップ
- [x] Docker Composeファイル作成
  - Goアプリケーション用コンテナ (Gin)
  - PostgreSQL用コンテナ
  - pgAdmin用コンテナ（DB管理用）
- [x] Dockerfileの作成（マルチステージビルド）
- [x] docker-compose.ymlの作成
- [x] .env.exampleファイル作成（環境変数テンプレート）

### 0.2 Ginプロジェクト初期化
- [x] backend/ディレクトリ作成
- [x] go mod init
- [x] 基本的な依存関係のインストール
  - gin-gonic/gin
  - gorm.io/gorm
  - gorm.io/driver/postgres
  - joho/godotenv
  - go-playground/validator/v10
  - uber-go/zap（ロギング）

### 0.3 データベースセットアップ
- [x] PostgreSQL初期化スクリプト作成
- [x] Atlas（ariga.io/atlas）のセットアップ
  - [x] atlas.hclの作成（backend/atlas.hcl）
  - [x] スキーマファイル（schema.hcl）の作成（cmd/migrate/schema.hcl）
  - [x] マイグレーションディレクトリ構成（cmd/migrate/migrations/）
- [x] 初期スキーマ作成（テーブル定義）
- [x] cmd/migrate/main.go作成（マイグレーション実行コマンド）
- [x] cmd/seed/main.go作成（シードデータ投入コマンド）
- [x] cmd/seed/にモデル別シードファイル作成

### 0.4 開発ツール設定
- [x] Makefileの作成（ビルド、テスト、実行コマンド）
  - [x] make api - APIサーバー起動
  - [x] make migrate - マイグレーション実行
  - [x] make seed - シードデータ投入
  - [x] make atlas-* - Atlas関連コマンド
- [x] air.tomlの作成（ホットリロード設定）
- [x] .gitignoreの更新
- [x] atlas.hcl（Atlas設定ファイル）の作成

## Phase 1: プロジェクト構造構築（オニオンアーキテクチャ）
### 1.1 ディレクトリ構造作成
```
backend/
├── cmd/
│   ├── api/
│   │   └── main.go              # APIサーバーエントリポイント ✓
│   ├── migrate/
│   │   ├── main.go              # マイグレーション実行 ✓
│   │   └── migrations/          # マイグレーションファイル ✓
│   └── seed/
│       ├── main.go              # シードデータ投入 ✓
│       ├── seed_*.go            # モデル別シードファイル ✓
│       └── helpers.go           # ヘルパー関数 ✓
├── internal/
│   ├── domain/                  # ドメイン層
│   │   ├── user/
│   │   ├── account/
│   │   ├── transaction/
│   │   ├── category/
│   │   ├── budget/
│   │   └── common/              ✓
│   ├── application/             # アプリケーション層 ✓
│   ├── infrastructure/          # インフラストラクチャ層
│   │   └── gorm/
│   │       └── model/           # GORMモデル定義 ✓
│   └── interface/               # インターフェース層
│       ├── middleware/          # HTTPミドルウェア ✓
│       └── router/              # HTTPルーター ✓
├── pkg/                         # 共有パッケージ
│   ├── config/                  ✓
│   ├── errors/                  ✓
│   └── logger/                  ✓
├── scripts/                     # ユーティリティスクリプト ✓
│   └── check-architecture.sh    ✓
├── docs/                        # ドキュメント
└── tests/                       # テストファイル
```

### 1.2 基本設定ファイル作成
- [x] pkg/config/config.go（設定管理）
- [x] pkg/logger/logger.go（ログ設定）
- [x] pkg/errors/errors.go（カスタムエラー）

### 1.3 HTTPサーバー基本実装
- [x] cmd/api/main.go（APIサーバーエントリポイント）
- [x] internal/interface/router/router.go（ルーティング設定）
  - [x] 全APIエンドポイントのルート定義（501 Not Implemented返却）
  - [x] ヘルスチェックエンドポイント実装（/health）
- [x] middleware実装
  - [x] CORS（Cross-Origin Resource Sharing）
  - [x] ロギング（リクエスト/レスポンスログ）
  - [x] エラーハンドリング（構造化エラーレスポンス）
  - [x] 認証（Auth0 JWT検証、JWK取得は未完了）
  - [x] リクエストID（トレーサビリティ）
  - [x] 権限チェック（パーミッションベース認可）

### 1.4 マイグレーション・シード実装
- [x] cmd/migrate/main.go（Atlasマイグレーション実行）
- [x] cmd/migrate/schema.hcl（スキーマ定義）
- [x] cmd/seed/main.go（シードデータ投入）
- [x] cmd/seed/seed_*.go（モデル別シードファイル）

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
- [x] atlas_migration.go（Atlasマイグレーション実行）※cmd/migrate/main.goで実装済み

### 4.2 GORMモデル定義
- [x] base.go（共通フィールド）
- [x] user.go
- [x] account.go
- [x] transaction.go
- [x] category.go
- [x] budget.go
- [x] report.go（AssetSnapshot, AssetForecast）
- [x] notification.go

### 4.3 リポジトリ実装
- [ ] base_repository.go（共通処理）
- [ ] user_repository.go
- [ ] account_repository.go
- [ ] transaction_repository.go
- [ ] category_repository.go
- [ ] budget_repository.go

### 4.4 外部サービス連携
- [ ] Auth0認証サービス
  - [x] JWT検証ミドルウェア実装
  - [ ] JWKセット取得・キャッシュ実装
- [ ] メール送信サービス
- [ ] キャッシュサービス（Redis）

## Phase 5: HTTPインターフェース層実装
### 5.1 ルーティング実装
- [x] 全APIエンドポイントのルート定義完了（501 Not Implemented返却）
  - [x] ユーザー管理: GET /users/me
  - [x] 口座管理: CRUD /accounts, GET /accounts/:id/movements
  - [x] 取引管理: CRUD /transactions, GET /transactions/search
  - [x] カテゴリ管理: CRUD /categories, GET /categories/masters
  - [x] 予算管理: CRUD /budgets, GET /budgets/suggestions
  - [x] レポート: GET /reports/summary, /reports/dashboard
  - [x] 通知設定: CRUD /notification-settings

### 5.2 ハンドラー実装
- [ ] user_handler.go（GET,PUT /users/me）
- [ ] account_handler.go（CRUD /accounts）
- [ ] transaction_handler.go（CRUD /transactions）
- [ ] category_handler.go（CRUD /categories）
- [ ] budget_handler.go（CRUD /budgets）
- [ ] report_handler.go（GET /reports/*）
- [ ] notification_handler.go（CRUD /notification-settings）

### 5.3 バリデーション
- [ ] リクエストバリデーション実装
- [ ] カスタムバリデーター作成

### 5.4 レスポンス処理
- [x] 統一エラーレスポンス形式（middleware/error_handler.go）
- [ ] 成功レスポンス統一形式
- [ ] ページネーションレスポンス

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
- [x] テスト自動実行
- [x] ビルド自動化
- [ ] Dockerイメージビルド

### 7.2 監視・ログ設定
- [x] 構造化ログ実装
- [ ] メトリクス収集
- [x] ヘルスチェックエンドポイント

### 7.3 ドキュメント作成
- [ ] API仕様書（OpenAPI/Swagger）
- [x] 開発者向けREADME
- [x] 環境構築手順書（README.mdに記載）
- [x] マイグレーション手順書（cmd/migrate/README.md）

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
- ローカル開発ではホットリロード（Air）を活用
- 環境変数は.envファイルで管理（.envはgitignore）
- データベースマイグレーションはAtlasで管理
  - スキーマ定義はHCL形式で記述（cmd/migrate/schema.hcl）
  - マイグレーションの自動生成と適用（make atlas-migrate-apply）
  - スキーマのバージョン管理（backend/cmd/migrate/migrations/）
- アーキテクチャ検証スクリプト（scripts/check-architecture.sh）でクリーンアーキテクチャ準拠を確認

## 次のステップ
1. **ドメイン層の実装**（Phase 2）
   - エンティティ、値オブジェクト、リポジトリインターフェースの定義
   - ビジネスロジックの実装
2. **アプリケーション層の実装**（Phase 3）
   - ユースケース（サービス）の実装
   - DTOの定義
3. **リポジトリ実装**（Phase 4の残り）
   - データベース接続管理
   - 各エンティティのリポジトリ実装
4. **APIハンドラー実装**（Phase 5の残り）
   - 各エンドポイントの実装
   - リクエスト/レスポンスの処理