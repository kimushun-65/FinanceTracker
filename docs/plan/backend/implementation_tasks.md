# バックエンド実装タスク一覧

## 概要
FinanceTrackerバックエンドをDocker環境で開発するためのタスクリストです。
AWS Lambda版の実装計画を基に、Gin framework + PostgreSQLで構築します。

## 進捗サマリー (2025-09-23更新)
- **Phase 0: 開発環境構築** - 100%完了 ✅
- **Phase 1: プロジェクト構造構築** - 100%完了 ✅
- **Phase 2: ドメイン層実装** - 100%完了 ✅
- **Phase 2.5: Auth0認証実装** - 100%完了 ✅ (JWT、クッキー、認証API全て実装完了)
- **Phase 2.7: DI層実装** - 100%完了 ✅
- **Phase 3: アプリケーション層実装** - 15%完了（認証サービスのみ実装済み）
- **Phase 4: インフラストラクチャ層実装** - 70%完了（DB接続・ユーザーリポジトリ実装済み、その他リポジトリ未実装）
- **Phase 5: HTTPインターフェース層実装** - 60%完了（認証API完成、その他API未実装）
- **Phase 6: テスト実装** - 未着手
- **Phase 7: CI/CD・運用設定** - 60%完了
- **Phase 8: Lambda統合** - 未着手

### 🎯 更新された重要マイルストーン
1. **Week 4 Day 1-2**: ✅ DI層実装完了
2. **Week 4 Day 3-5**: アプリケーション層実装（DTO・サービス）
3. **Week 5目標**: 全主要ビジネスAPI完成
4. **Week 6目標**: テスト実装・統合完了

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
│   ├── di/                      # 🆕 DI（依存関係注入）層
│   │   ├── container.go         # メインDIコンテナ
│   │   ├── config.go            # 統一設定管理
│   │   ├── interfaces.go        # DIインターフェース
│   │   └── providers/           # プロバイダー群
│   │       ├── database.go      # データベース・外部サービス
│   │       ├── repository.go    # リポジトリ初期化
│   │       ├── service.go       # アプリケーションサービス
│   │       ├── handler.go       # HTTPハンドラー
│   │       └── middleware.go    # ミドルウェア
│   ├── domain/                  # ドメイン層
│   │   ├── user/
│   │   ├── account/
│   │   ├── transaction/
│   │   ├── category/
│   │   ├── budget/
│   │   ├── asset/               # 追加予定
│   │   ├── notification/        # 追加予定
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
  - [x] `internal/di/config.go` - 環境変数の統一管理（117行）
  - [x] 設定検証とデフォルト値設定
- [x] DIコンテナインターフェース定義 ✅
  - [x] `internal/di/interfaces.go` - コンテナインターフェース（10行）
- [x] メインDIコンテナ実装 ✅
  - [x] `internal/di/container.go` - 依存関係管理の中核（209行）

**成果物**:
```
internal/di/
├── container.go              # メインDIコンテナ
├── config.go                 # 統一設定管理
├── interfaces.go             # DIインターフェース
├── test_container.go         # テスト用DIコンテナ
└── providers/                # プロバイダー群
    ├── database.go           # データベース・外部サービス
    ├── repository.go         # リポジトリ初期化
    ├── service.go            # アプリケーションサービス
    ├── handler.go            # HTTPハンドラー
    └── middleware.go         # ミドルウェア
```

### 2.7.2 実装方針変更
**注記**: プロバイダーを個別ファイルに分割せず、container.go内にメソッドとして統合実装
- [x] initInfrastructure() - データベース・Auth0接続
- [x] initRepositories() - リポジトリ初期化
- [x] initServices() - サービス初期化  
- [x] initHandlers() - ハンドラー初期化

### 2.7.3 既存コード統合 ✅
**期間**: 完了
**目的**: 既存の実装をDIコンテナに統合し、main.goを簡略化

**タスク**:
- [x] main.go の大幅簡略化 ✅
  - [x] DIコンテナ初期化によるワンライン依存関係解決
- [x] 環境変数管理の統一化 ✅
  - [x] 既存の分散した設定読み込みをDI層に集約

### 2.7.4 テスト環境構築
**注記**: 必要になった時点で実装予定

**検証**:
- [ ] DIコンテナの正常初期化確認
- [ ] 全既存機能の動作確認（認証API）
- [ ] テスト用コンテナでのモック注入確認
- [ ] メモリリーク・接続プールの動作確認

**DI層導入の効果測定**:

**Before（DIなし）**:
```go
// main.go - 100+ lines の手動依存関係構築
func main() {
    config := loadConfig()
    db := setupDatabase(config)
    auth0Client := setupAuth0(config)
    userRepo := repository.NewUserRepository(db)
    authService := service.NewAuthService(userRepo, auth0Client)
    authHandler := handler.NewAuthHandler(authService)
    // ... 多数の初期化コード
    router := setupRouter(authHandler, /* many handlers */)
    router.Run()
}
```

**After（DIあり）**:
```go
// main.go - 20 lines で完結
func main() {
    container, err := di.NewContainer()
    if err != nil {
        log.Fatal(err)
    }
    defer container.Shutdown(context.Background())
    
    router := router.SetupRouter(container)
    server := setupServer(router, container.GetConfig())
    gracefulShutdown(server, container)
}
```

**期待される効果**:
1. **開発効率50-70%向上**: 新機能追加時の依存関係構築が自動化
2. **テスタビリティ大幅改善**: モックオブジェクトの簡単な注入
3. **Clean Architecture完全実現**: 依存関係逆転の適切な実装
4. **保守性向上**: 設定管理の統一化・エラーハンドリングの改善

## Phase 3: アプリケーション層実装 🔄 (15%完了)

### 3.1 DTOの定義 (DI層導入後に効率化)
- [x] auth_dto.go (UserInfo, TokenClaims)
- [ ] user_dto.go 🔥 **DI層導入後の優先実装**
- [ ] account_dto.go 🔥 **DI層導入後の優先実装**
- [ ] transaction_dto.go 🔥 **DI層導入後の優先実装**
- [ ] category_dto.go
- [ ] budget_dto.go
- [ ] asset_dto.go（AssetSnapshot, AssetForecast）
- [ ] notification_dto.go（NotificationSettings）

### 3.2 ユースケース実装 (DI層により依存関係注入が自動化)
- [x] ユーザー管理
  - [ ] user_service.go ⚠️ **DI層導入後に完成**
  - [x] auth_service.go ✅ (SyncUser, GetUserByAuth0ID実装済み)
- [ ] 口座管理 🔥 **DI層導入後の優先実装**
  - [ ] account_service.go (DIコンテナで自動注入)
  - [ ] account_movement_service.go (DIコンテナで自動注入)
- [ ] 取引管理 🔥 **DI層導入後の優先実装**
  - [ ] transaction_service.go (DIコンテナで自動注入)
- [ ] カテゴリ管理 (DI層により効率化)
  - [ ] category_service.go (DIコンテナで自動注入)
  - [ ] category_master_service.go (DIコンテナで自動注入)
- [ ] 予算管理 (DI層により効率化)
  - [ ] budget_service.go (DIコンテナで自動注入)
  - [ ] budget_suggestion_service.go (DIコンテナで自動注入)
- [ ] 資産管理 (DI層により効率化)
  - [ ] asset_snapshot_service.go (DIコンテナで自動注入)
  - [ ] asset_forecast_service.go (DIコンテナで自動注入)
  - [ ] asset_analysis_service.go (DIコンテナで自動注入)
- [ ] 通知管理 (DI層により効率化)
  - [ ] notification_settings_service.go (DIコンテナで自動注入)
  - [ ] notification_sender_service.go (DIコンテナで自動注入)
- [ ] レポート生成 (DI層により効率化)
  - [ ] summary_service.go (DIコンテナで自動注入)
  - [ ] dashboard_service.go (DIコンテナで自動注入)

### 3.3 トランザクション管理 (DI層で一元管理)
- [ ] transaction_manager.go 🔥 **DI層実装と同時実装**

## Phase 4: インフラストラクチャ層実装 🔄 (70%完了)

### 4.1 データベース接続
- [x] connection.go（GORM接続管理）✅ `pkg/database/database.go`実装済み
- [ ] transaction.go（トランザクション管理）🔥 **優先実装**
- [x] atlas_migration.go（Atlasマイグレーション実行）✅ `cmd/migrate/main.go`で実装済み

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
- [ ] base_repository.go（共通処理）🔥 **優先実装**
- [x] user_repository.go ✅ (Save, FindByID, FindByAuth0UserID, Exists, Delete実装済み)
- [ ] account_repository.go 🔥 **優先実装**
- [ ] account_movement_repository.go 🔥 **優先実装**
- [ ] transaction_repository.go 🔥 **優先実装**
- [ ] category_repository.go
- [ ] category_master_repository.go
- [ ] budget_repository.go
- [ ] budget_suggestion_repository.go
- [ ] asset_snapshot_repository.go
- [ ] asset_forecast_repository.go
- [ ] notification_settings_repository.go

### 4.4 外部サービス連携
- [x] Auth0認証サービス ✅ **完全実装済み**
  - [x] JWT検証ミドルウェア実装 ✅ `auth0/middleware.go`
  - [x] JWKセット取得・キャッシュ実装 ✅ `auth0/client.go`
  - [x] ユーザー情報取得 ✅
  - [x] ログイン/ログアウトURL生成 ✅
- [ ] メール送信サービス
- [ ] キャッシュサービス（Redis）

## Phase 5: HTTPインターフェース層実装 🔄 (60%完了)

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
  - [x] ログイン・ログアウト・コールバック・ユーザー情報取得
  - [x] トークン管理・認証状態チェック
- [ ] user_handler.go（GET,PUT /users/me）🔥 **優先実装**
- [ ] account_handler.go（CRUD /accounts）🔥 **優先実装**
- [ ] transaction_handler.go（CRUD /transactions）🔥 **優先実装**
- [ ] category_handler.go（CRUD /categories）
- [ ] budget_handler.go（CRUD /budgets）
- [ ] asset_handler.go（GET /assets/snapshots, /assets/forecasts）
- [ ] report_handler.go（GET /reports/*）
- [ ] notification_handler.go（CRUD /notification-settings）

### 5.3 バリデーション
- [ ] リクエストバリデーション実装 🔥 **優先実装**
- [ ] カスタムバリデーター作成

### 5.4 レスポンス処理
- [x] 統一エラーレスポンス形式 ✅ `middleware/error_handler.go`
- [ ] 成功レスポンス統一形式 🔥 **優先実装**
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
- Reportテーブル/エンティティは実装しない（レポートは既存データの集計表示のみ）

## 🚀 更新された次のステップ（優先順位順）

### 🔥 **最優先（Week 4 Day 3-5）**: インフラストラクチャ層実装
1. **コアリポジトリ実装**
   - `account_repository.go` - 口座管理の基盤
   - `transaction_repository.go` - 取引記録の基盤
   - `category_repository.go` - カテゴリ管理の基盤
   - `base_repository.go` - 共通CRUD処理

### 🎯 **高優先（Week 5 Day 1-3）**: アプリケーション層実装
2. **DTO・サービス実装**
   - ユーザー管理（user_dto.go, user_service.go）
   - 口座管理（account_dto.go, account_service.go）
   - 取引管理（transaction_dto.go, transaction_service.go）

### 📋 **中優先（Week 5 Day 4-5）**: インターフェース層実装
3. **HTTPハンドラー実装**
   - user_handler.go（GET/PUT /users/me）
   - account_handler.go（CRUD /accounts）
   - transaction_handler.go（CRUD /transactions）

### 🎯 **高優先（Week 5）**: 残りビジネス機能
5. **カテゴリー管理** (DI層効率活用)
   - `category_repository.go`, `category_master_repository.go`
   - `category_service.go`, `category_master_service.go`
   - `category_handler.go`

6. **予算管理** (DI層効率活用)
   - `budget_repository.go`, `budget_suggestion_repository.go`
   - `budget_service.go`
   - `budget_handler.go`

### 📊 **中優先（Week 6）**: 分析・レポート機能
7. **レポート・資産管理** (DI層効率活用)
   - `summary_service.go`, `dashboard_service.go`
   - `asset_snapshot_repository.go`, `asset_forecast_repository.go`
   - `report_handler.go`, `asset_handler.go`

### 🔧 **基盤強化** (DI層により簡素化)
8. **共通機能実装**
   - `base_repository.go` (共通CRUD処理)
   - `transaction_manager.go` (DB トランザクション管理)
   - 成功レスポンス統一形式
   - リクエストバリデーション

### ✅ **完了済み**
- ✅ DI層実装（依存性注入、コンテナ管理）
- ✅ Auth0認証システム（フル実装完了）
- ✅ プロジェクト基盤（アーキテクチャ、ドメイン層）
- ✅ インフラ基盤（DB接続、ユーザーリポジトリ）
- ✅ HTTP基盤（ルーティング、ミドルウェア、認証API）

### 📋 更新された実装推奨順序
```
Week 4 Day 1-2: ✅ DI層実装（基盤強化）- 完了
Week 4 Day 3-5: インフラ層実装（リポジトリ群）
Week 5 Day 1-3: アプリケーション層実装（DTO・サービス）
Week 5 Day 4-5: インターフェース層実装（HTTPハンドラー）
Week 6: テスト・統合・最適化（テスト用DIコンテナ活用）
```

### 🎯 DI層導入による開発効率改善
**Before（現在）**: 各機能実装時に手動で依存関係構築
**After（DI層導入後）**: 自動的な依存関係解決、50-70%工数削減

**例: 新機能追加時の工数比較**
```go
// Before（現在）: 15-20行の依存関係構築が必要
func setupNewFeature() {
    db := getDB()
    userRepo := repository.NewUserRepository(db)
    newRepo := repository.NewNewRepository(db)
    newService := service.NewNewService(newRepo, userRepo, /* 他の依存関係 */)
    newHandler := handler.NewNewHandler(newService)
    // ルーターに手動追加...
}

// After（DI導入後）: 1行でコンテナに追加
// DIコンテナが自動的に依存関係を解決
container.GetNewHandler() // 自動的に全依存関係が注入済み
```