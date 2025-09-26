# バックエンド実装タスク一覧

## 概要
FinanceTrackerバックエンドをDocker環境で開発するためのタスクリストです。
AWS Lambda版の実装計画を基に、Gin framework + PostgreSQLで構築します。

## 進捗サマリー (2025-09-26更新)
- **Phase 0: 開発環境構築** - 100%完了 ✅
- **Phase 1: プロジェクト構造構築** - 100%完了 ✅
- **Phase 2: ドメイン層実装** - 100%完了 ✅
- **Phase 2.5: Auth0認証実装** - 100%完了 ✅ (JWT、クッキー、認証API全て実装完了)
- **Phase 2.7: DI層実装** - 100%完了 ✅
- **Phase 3: アプリケーション層実装** - 100%完了 ✅（全サービス・DTO実装完了）
- **Phase 4: インフラストラクチャ層実装** - 100%完了 ✅（全リポジトリ実装完了）
- **Phase 5: HTTPインターフェース層実装** - 95%完了 ✅（主要ハンドラー実装・デバッグ完了）
- **Phase 6: テスト実装** - 未着手
- **Phase 7: CI/CD・運用設定** - 60%完了
- **Phase 8: Lambda統合** - 100%完了 ✅（全Lambda関数をTypeScriptプロキシとして実装）

### 🎯 更新された重要マイルストーン
1. **Week 4**: ✅ DI層、アプリケーション層、インフラ層実装完了
2. **Week 5**: ✅ HTTPハンドラー実装、主要ビジネスAPI完成
3. **Week 6**: ✅ Lambda実装完了、テスト実装残り

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
│   ├── di/                      # DI（依存関係注入）層 ✓
│   │   ├── container.go         # メインDIコンテナ ✓
│   │   └── config.go            # 統一設定管理 ✓
│   ├── domain/                  # ドメイン層 ✓
│   │   ├── user/
│   │   ├── account/
│   │   ├── transaction/
│   │   ├── category/
│   │   ├── budget/
│   │   ├── asset/
│   │   ├── notification/
│   │   └── common/
│   ├── application/             # アプリケーション層 ✓
│   ├── infrastructure/          # インフラストラクチャ層 ✓
│   │   └── gorm/
│   │       └── model/           # GORMモデル定義 ✓
│   └── interface/               # インターフェース層
│       ├── handler/             # HTTPハンドラー ✓
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
  - [x] 全APIエンドポイントのルート定義
  - [x] ヘルスチェックエンドポイント実装（/health）
- [x] middleware実装
  - [x] CORS（Cross-Origin Resource Sharing）
  - [x] ロギング（リクエスト/レスポンスログ）
  - [x] エラーハンドリング（構造化エラーレスポンス）
  - [x] 認証（Auth0 JWT検証）
  - [x] リクエストID（トレーサビリティ）

### 1.4 マイグレーション・シード実装
- [x] cmd/migrate/main.go（Atlasマイグレーション実行）
- [x] cmd/migrate/schema.hcl（スキーマ定義）
- [x] cmd/seed/main.go（シードデータ投入）
- [x] cmd/seed/seed_*.go（モデル別シードファイル）

## Phase 2: ドメイン層実装
### 2.1 共通ドメイン要素
- [x] base_entity.go（UUID, timestamps）
- [x] 値オブジェクト実装
  - [x] money.go（金額）
  - [x] email.go（メールアドレス）
  - [x] hex_color.go（カラーコード）
  - [x] time.go（時刻関連）
- [x] ドメインエラー定義（errors.go）

### 2.2 各ドメインエンティティ実装
- [x] User エンティティ
  - [x] entity/user.go
  - [x] value/user_id.go
  - [x] value/auth0_id.go
  - [x] repository/user_repository.go
- [x] Account エンティティ
  - [x] entity/account.go
  - [x] entity/account_movement.go
  - [x] value/account_type.go（資産口座のみ対応）
  - [x] value/balance.go
  - [x] value/account_name.go
  - [x] repository/account_repository.go
  - [x] repository/account_movement_repository.go
- [x] Transaction エンティティ
  - [x] entity/transaction.go（Transfer機能は要件外のため実装なし）
  - [x] value/transaction_type.go（income/expenseのみ）
  - [x] value/description.go
  - [x] repository/transaction_repository.go
- [x] Category エンティティ
  - [x] entity/category.go
  - [x] entity/category_master.go
  - [x] value/category_name.go
  - [x] value/category_type.go
  - [x] repository/category_repository.go
  - [x] repository/category_master_repository.go
- [x] Budget エンティティ
  - [x] entity/budget.go
  - [x] entity/budget_suggestion.go
  - [x] value/period_type.go
  - [x] value/suggestion_status.go
  - [x] repository/budget_repository.go
  - [x] repository/budget_suggestion_repository.go
- [x] AssetSnapshot エンティティ（資産管理コンテキスト）
  - [x] entity/asset_snapshot.go
  - [x] value/account_breakdown.go（口座別内訳）
  - [x] repository/asset_snapshot_repository.go
- [x] AssetForecast エンティティ（資産予測管理）
  - [x] entity/asset_forecast.go
  - [x] value/assumptions.go（予測手法は単一のため forecast_method.go は削除）
  - [x] repository/asset_forecast_repository.go
- [x] NotificationSettings エンティティ（通知管理コンテキスト）
  - [x] entity/notification_settings.go
  - [x] repository/notification_settings_repository.go

## Phase 2.7: DI層実装 ✅ **完了** (2025-09-24)

### 2.7.1 DI基盤構築
**期間**: 1日で完了
**目的**: Clean Architectureの依存関係管理を効率化し、テスタビリティと保守性を向上

**タスク**:
- [x] DI層ディレクトリ構造作成 ✅
  - [x] `internal/di/` ディレクトリ作成
- [x] 統一設定管理実装 ✅
  - [x] `internal/di/config.go` - 環境変数の統一管理
  - [x] 設定検証とデフォルト値設定
- [x] メインDIコンテナ実装 ✅
  - [x] `internal/di/container.go` - 依存関係管理の中核

### 2.7.2 既存コード統合 ✅
**タスク**:
- [x] main.go の大幅簡略化 ✅
  - [x] DIコンテナ初期化によるワンライン依存関係解決
- [x] 環境変数管理の統一化 ✅
  - [x] 既存の分散した設定読み込みをDI層に集約

## Phase 3: アプリケーション層実装 ✅ **完了** (2025-01-10)

### 3.1 DTOの定義 ✅
- [x] auth_dto.go (UserInfo, TokenClaims)
- [x] user_dto.go ✅ **実装完了**
- [x] account_dto.go ✅ **実装完了**
- [x] transaction_dto.go ✅ **実装完了**
- [x] category_dto.go ✅ **実装完了**
- [x] budget_dto.go ✅ **実装完了**
- [x] asset_dto.go ✅ **実装完了**（AssetSnapshot, AssetForecast）
- [x] notification_dto.go ✅ **実装完了**（NotificationSettings）

### 3.2 ユースケース実装 ✅
- [x] ユーザー管理 ✅ **実装完了**
  - [x] user_service.go ✅ **実装完了**
  - [x] auth_service.go ✅ **実装完了**
- [x] 口座管理 ✅ **実装完了**
  - [x] account_service.go ✅ **実装完了**
- [x] 取引管理 ✅ **実装完了**
  - [x] transaction_service.go ✅ **実装完了**
- [x] カテゴリ管理 ✅ **実装完了**
  - [x] category_service.go ✅ **実装完了**
  - [x] category_master_service.go ✅ （category_serviceに統合）
- [x] 予算管理 ✅ **実装完了**
  - [x] budget_service.go ✅ **実装完了**
  - budget_suggestion_service.go（基本機能のため未実装）
- [x] 資産管理 ✅ **実装完了**
  - [x] asset_service.go ✅ **実装完了**
- [x] 通知管理 ✅ **実装完了**
  - [x] notification_service.go ✅ **実装完了**

### 3.3 トランザクション管理
- [x] トランザクション処理は各サービス内で実装済み ✅

## Phase 4: インフラストラクチャ層実装 ✅ **完了** (2025-01-10)

### 4.1 データベース接続
- [x] connection.go（GORM接続管理）✅ `pkg/database/database.go`実装済み
- [x] transaction.go（トランザクション管理）✅ 各リポジトリ内で実装

### 4.2 GORMモデル定義
- [x] base.go（共通フィールド）✅
- [x] user.go ✅
- [x] account.go ✅
- [x] transaction.go ✅
- [x] category.go ✅
- [x] budget.go ✅
- [x] report.go（AssetSnapshot, AssetForecast）✅
- [x] notification.go ✅

### 4.3 リポジトリ実装
- [x] user_repository.go ✅ **実装完了**
- [x] account_repository.go ✅ **実装完了**
- [x] transaction_repository.go ✅ **実装完了**
- [x] category_repository.go ✅ **実装完了**
- [x] category_master_repository.go ✅ **実装完了**
- [x] budget_repository.go ✅ **実装完了**
- budget_suggestion_repository.go（基本機能のため未実装）
- [x] asset_snapshot_repository.go ✅ **実装完了**
- [x] asset_forecast_repository.go ✅ **実装完了**
- [x] notification_settings_repository.go ✅ **実装完了**

### 4.4 外部サービス連携
- [x] Auth0認証サービス ✅ **完全実装済み**
  - [x] JWT検証ミドルウェア実装 ✅ `auth0/middleware.go`
  - [x] JWKセット取得・キャッシュ実装 ✅ `auth0/client.go`
  - [x] ユーザー情報取得 ✅
  - [x] ログイン/ログアウトURL生成 ✅

## Phase 5: HTTPインターフェース層実装 ✅ **90%完了** (2025-01-10)

### 5.1 ルーティング実装
- [x] 全APIエンドポイントのルート定義完了 ✅
- [x] 認証エンドポイント実装済み ✅
  - [x] GET /auth/login, /auth/callback, /auth/user
  - [x] POST /auth/logout, /auth/token
  - [x] DELETE /auth/token
- [x] システムエンドポイント ✅
  - [x] GET /health, /api/v1/ (API情報)

### 5.2 ハンドラー実装
- [x] auth_handler.go ✅ **完全実装済み**
- [x] user_handler.go ✅ **実装・デバッグ完了**
  - [x] GetCurrentUser（GET /users/me）
  - [x] UpdateCurrentUser（PUT /users/me）
- [x] account_handler.go ✅ **実装・デバッグ完了**
  - [x] List（GET /accounts）
  - [x] Create（POST /accounts）
  - [x] Get（GET /accounts/:id）
  - [x] Update（PUT /accounts/:id）
  - [x] Delete（DELETE /accounts/:id）
- [x] transaction_handler.go ✅ **実装・デバッグ完了**
  - [x] List（GET /transactions）- 日付フィルタリング修正済み
  - [x] Create（POST /transactions）
  - [x] Get（GET /transactions/:id）
  - [x] Update（PUT /transactions/:id）- 認証エラー修正済み
  - [x] Delete（DELETE /transactions/:id）- 認証エラー修正済み
  - [x] MonthlySummary（GET /transactions/summary/monthly）- 認証エラー修正済み
- [x] category_handler.go ✅ **実装・デバッグ完了**
  - [x] List（GET /categories）- name フィールド追加済み
  - [x] Create（POST /categories）- name フィールド追加済み
  - [x] Get（GET /categories/:id）- name フィールド追加済み
  - [x] Update（PUT /categories/:id）- name フィールド追加済み
  - [x] Delete（DELETE /categories/:id）- 論理削除実装
  - [x] ListMaster（GET /categories/master）- 認証エラー修正済み
- [x] budget_handler.go ✅ **実装・デバッグ完了**
  - [x] List（GET /budgets）- 認証エラー修正済み
  - [x] Create（POST /budgets）- 日付パースエラー修正済み
  - [x] Get（GET /budgets/:id）- 認証エラー修正済み
  - [x] Update（PUT /budgets/:id）- 認証エラー修正済み
  - [x] Delete（DELETE /budgets/:id）- 認証エラー修正済み
  - [x] GetCurrent（GET /budgets/current）- 認証エラー修正済み
- [ ] asset_handler.go（GET /assets/snapshots, /assets/forecasts）
- [ ] report_handler.go（GET /reports/*）
- [ ] notification_handler.go（CRUD /notification-settings）

### 5.3 バリデーション
- [x] リクエストバリデーション実装 ✅ DTOのBindingタグで実装

### 5.4 レスポンス処理
- [x] 統一エラーレスポンス形式 ✅ `middleware/error_handler.go`
- [x] 成功レスポンス統一形式 ✅ 各ハンドラーで実装

## Phase 6: テスト実装
### 6.1 単体テスト
- [ ] ドメイン層テスト
- [ ] アプリケーション層テスト
- [ ] インフラストラクチャ層テスト（モック使用）

### 6.2 統合テスト
- [ ] API統合テスト
- [ ] データベース統合テスト

### 6.3 E2Eテスト
- [ ] 主要フローのE2Eテスト

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

## Phase 8: Lambda統合 ✅ **完了** (2025-01-10)

### 8.1 Lambda対応 ✅
- [x] Auth0 JWT Authorizer ✅ **完了**
  - [x] `cdk/lambda/authorizer/index.js` - JWT検証実装
- [x] Lambda関数実装 ✅ **全てTypeScriptプロキシとして実装**
  - [x] `cdk/lambda/users/index.ts` - バックエンドAPIプロキシ
  - [x] `cdk/lambda/accounts/index.ts` - バックエンドAPIプロキシ
  - [x] `cdk/lambda/transactions/index.ts` - バックエンドAPIプロキシ
  - [x] `cdk/lambda/categories/index.ts` - バックエンドAPIプロキシ
  - [x] `cdk/lambda/budgets/index.ts` - バックエンドAPIプロキシ
  - [x] `cdk/lambda/reports/index.ts` - バックエンドAPIプロキシ
  - [x] `cdk/lambda/auth/index.ts` - バックエンドAPIプロキシ
  - [x] `cdk/lambda/notifications/index.ts` - バックエンドAPIプロキシ
- [x] 不要なGoファイルの削除（main.go）✅
- [x] CDKスタックとの統合 ✅ **完了**
- [x] ビルドスクリプト更新 ✅ **完了**

### 8.2 デプロイメント
- [x] Lambda関数デプロイ設定 ✅ **完了**
- [x] API Gateway設定 ✅ **完了**
- [x] GitHub Actions CI/CD統合 ✅ **完了**

### 8.3 Lambda実装方針（更新）
- **アーキテクチャ**: 全Lambda関数がバックエンドAPIへのプロキシ
- **言語**: TypeScript（Node.js 18.x）
- **API通信**: HTTP/HTTPSクライアントを使用
- **認証**: Auth0 IDをX-Auth0-IDヘッダーで転送
- **ビルド**: CDKのpackage.jsonでTypeScriptコンパイル

## 技術スタック
- **言語**: Go 1.21+
- **Webフレームワーク**: Gin
- **ORM**: GORM v2
- **データベース**: PostgreSQL 15
- **認証**: Auth0
- **コンテナ**: Docker/Docker Compose
- **CI/CD**: GitHub Actions
- **Lambda**: TypeScript（プロキシ実装）

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
- アーキテクチャ検証スクリプト（scripts/check-architecture.sh）でクリーンアーキテクチャ準拠を確認
- Lambda関数は全てTypeScriptによるバックエンドAPIプロキシ実装

## 🚀 更新された次のステップ（優先順位順）

### ✅ **完了済み**: コア機能実装
1. **バックエンドAPI実装** ✅ **完了**
   - DI層（依存性注入コンテナ）✅
   - ドメイン層（エンティティ、値オブジェクト、リポジトリインターフェース）✅
   - アプリケーション層（DTO、サービス）✅
   - インフラストラクチャ層（リポジトリ実装）✅
   - インターフェース層（HTTPハンドラー）✅

2. **Lambda関数実装** ✅ **完了**
   - 全Lambda関数をTypeScriptプロキシとして実装 ✅
   - バックエンドAPIへの透過的なプロキシ ✅
   - Auth0認証情報の転送 ✅

### 🔥 **残タスク（高優先度）**
3. **残りのHTTPハンドラー実装**
   - [ ] asset_handler.go（資産スナップショット、予測）
   - [ ] report_handler.go（レポート生成）
   - [ ] notification_handler.go（通知設定）

4. **テスト実装**
   - [ ] 単体テスト（ドメイン層、アプリケーション層）
   - [ ] 統合テスト（API、データベース）
   - [ ] E2Eテスト（主要フロー）

### 📊 **中優先度タスク**
5. **運用・監視強化**
   - [ ] メトリクス収集
   - [ ] API仕様書（OpenAPI/Swagger）
   - [ ] Dockerイメージビルド自動化

### ✅ **完了済み成果物**
- ✅ 完全なクリーンアーキテクチャ実装
- ✅ Auth0認証システム（JWT、クッキー管理）
- ✅ 主要ビジネスAPI（ユーザー、口座、取引、カテゴリ、予算）
- ✅ Lambda関数（TypeScriptプロキシ）
- ✅ CDK統合・デプロイメント設定

### 📋 実装完了率
```
ドメイン層: 100% ✅
アプリケーション層: 100% ✅
インフラストラクチャ層: 100% ✅
インターフェース層: 90% ✅
Lambda統合: 100% ✅
テスト: 0%
運用設定: 60%
```