# FinanceTracker バックエンド 最終ディレクトリ構成

## 概要
FinanceTrackerバックエンドは、Clean Architecture（オニオンアーキテクチャ）に基づいて設計された、DDD（ドメイン駆動設計）実装のGo言語プロジェクトです。認証システム、主要なドメインモデル、データベース管理機能が実装済みで、拡張性と保守性に優れた構成となっています。

## 技術スタック
- **言語**: Go 1.24.0
- **Webフレームワーク**: Gin v1.11.0
- **ORM**: GORM v1.25.12
- **データベース**: PostgreSQL 15
- **認証**: Auth0
- **ログ**: Zap v1.27.0
- **マイグレーション**: Atlas CLI
- **ドキュメント**: Swagger/OpenAPI
- **コンテナ**: Docker

## 完全ディレクトリ構成

```
backend/
├── README.md                                     # プロジェクト概要・セットアップ手順
├── Dockerfile                                    # マルチステージビルド対応Docker設定
├── go.mod                                        # Go Module設定
├── go.sum                                        # Goライブラリ依存関係ロックファイル
├── atlas.hcl                                     # Atlas CLI設定（スキーマ管理）
├── .air.toml                                     # ホットリロード設定（推測）
├── .gitignore                                    # Git除外設定
│
├── scripts/                                      # 開発・運用スクリプト
│   └── check-architecture.sh                    # Clean Architecture依存関係チェックツール
│
├── tmp/                                          # 一時ファイル（開発時生成）
│   ├── build-errors.log                         # ビルドエラーログ
│   └── main                                      # ビルド成果物
│
├── docs/                                         # API ドキュメント（Swagger生成）
│   ├── docs.go                                   # Swagger設定
│   ├── swagger.json                              # OpenAPI JSON
│   └── swagger.yaml                              # OpenAPI YAML
│
├── db/                                           # データベース関連ファイル
│   └── init/
│       └── 01_init.sql                           # PostgreSQL初期化SQL
│
├── cmd/                                          # コマンドラインツール群
│   ├── api/
│   │   └── main.go                               # 🟢 APIサーバーエントリーポイント (実装済み)
│   ├── migrate/
│   │   ├── README.md                             # マイグレーション手順書
│   │   ├── main.go                               # 🟢 Atlasマイグレーション実行ツール (実装済み)
│   │   ├── schema.hcl                            # 🟢 データベーススキーマ定義 (実装済み)
│   │   └── migrations/                           # マイグレーションファイル保存先
│   │       ├── 20240101000000_initial.sql        # 🟢 初期スキーマ (233行, 実装済み)
│   │       ├── 20241201000000_add_account_movements.sql # 🟢 口座移動テーブル追加 (23行, 実装済み)
│   │       └── atlas.sum                         # Atlasマイグレーションチェックサム
│   └── seed/
│       ├── main.go                               # 🟢 シードデータ実行コマンド (実装済み)
│       ├── helpers.go                            # 🟢 シード用ヘルパー関数 (実装済み)
│       ├── seed_users.go                         # 🟢 ユーザーシードデータ (実装済み)
│       ├── seed_accounts.go                      # 🟢 口座シードデータ (実装済み)
│       ├── seed_categories.go                    # 🟢 カテゴリシードデータ (実装済み)
│       ├── seed_category_masters.go              # 🟢 カテゴリマスターシードデータ (実装済み)
│       ├── seed_transactions.go                  # 🟢 取引シードデータ (実装済み)
│       └── seed_budgets.go                       # 🟢 予算シードデータ (実装済み)
│
├── pkg/                                          # 共通パッケージ群（アプリケーション横断）
│   ├── config/
│   │   └── config.go                             # 🟢 環境変数ベース設定管理 (99行, 実装済み)
│   ├── database/
│   │   └── database.go                           # 🟢 GORM データベース接続管理 (実装済み)
│   ├── errors/
│   │   └── errors.go                             # 🟢 アプリケーション共通エラー (258行, 実装済み)
│   └── logger/
│       └── logger.go                             # 🟢 Zap ベース構造化ログ機能 (実装済み)
│
└── internal/                                     # Clean Architecture実装（プライベートパッケージ）
    │
    ├── domain/                                   # 🎯 ドメイン層（最内層・ビジネスロジック）
    │   ├── common/                               # 共通ドメインコンポーネント
    │   │   ├── base_entity.go                    # 🟢 ベースエンティティ（UUID・タイムスタンプ）
    │   │   ├── errors.go                         # 🟢 ドメインエラー定義 (191行, 実装済み)
    │   │   ├── repository/
    │   │   │   └── base_repository.go            # 🟢 ベースリポジトリインターフェース
    │   │   └── value/                            # 共通値オブジェクト
    │   │       ├── email.go                      # 🟢 メールアドレス値オブジェクト
    │   │       ├── hex_color.go                  # 🟢 16進カラー値オブジェクト
    │   │       ├── money.go                      # 🟢 金額値オブジェクト (178行, 実装済み)
    │   │       └── time.go                       # 🟢 時間関連値オブジェクト
    │   │
    │   ├── user/                                 # 👤 ユーザードメインコンテキスト
    │   │   ├── entity/
    │   │   │   └── user.go                       # 🟢 ユーザーエンティティ (112行, 実装済み)
    │   │   ├── repository/
    │   │   │   └── user_repository.go            # 🟢 ユーザーリポジトリインターフェース (21行, 実装済み)
    │   │   └── value/
    │   │       ├── auth0_id.go                   # 🟢 Auth0ID値オブジェクト
    │   │       └── user_id.go                    # 🟢 ユーザーID値オブジェクト
    │   │
    │   ├── account/                              # 🏦 口座ドメインコンテキスト
    │   │   ├── entity/
    │   │   │   ├── account.go                    # 🟢 口座エンティティ (187行, 実装済み)
    │   │   │   └── account_movement.go           # 🟢 口座残高変動エンティティ
    │   │   ├── repository/
    │   │   │   └── account_repository.go         # 🟢 口座リポジトリインターフェース
    │   │   └── value/
    │   │       ├── account_name.go               # 🟢 口座名値オブジェクト
    │   │       ├── account_type.go               # 🟢 口座種別値オブジェクト
    │   │       └── balance.go                    # 🟢 残高値オブジェクト (178行, 実装済み)
    │   │
    │   ├── transaction/                          # 💰 取引ドメインコンテキスト
    │   │   ├── entity/
    │   │   │   └── transaction.go                # 🟢 取引エンティティ (177行, 実装済み)
    │   │   ├── repository/
    │   │   │   └── transaction_repository.go     # 🟢 取引リポジトリインターフェース
    │   │   └── value/
    │   │       ├── description.go                # 🟢 取引説明値オブジェクト
    │   │       └── transaction_type.go           # 🟢 取引種別値オブジェクト
    │   │
    │   ├── category/                             # 📂 カテゴリドメインコンテキスト
    │   │   ├── entity/
    │   │   │   ├── category.go                   # 🟢 ユーザーカテゴリエンティティ
    │   │   │   └── category_master.go            # 🟢 システムカテゴリマスターエンティティ
    │   │   ├── repository/
    │   │   │   ├── category_repository.go        # 🟢 カテゴリリポジトリインターフェース
    │   │   │   └── category_master_repository.go # 🟢 カテゴリマスターリポジトリインターフェース
    │   │   └── value/
    │   │       ├── category_name.go              # 🟢 カテゴリ名値オブジェクト
    │   │       └── category_type.go              # 🟢 カテゴリタイプ値オブジェクト
    │   │
    │   ├── budget/                               # 📊 予算ドメインコンテキスト
    │   │   ├── entity/
    │   │   │   ├── budget.go                     # 🟢 予算エンティティ (187行, 実装済み)
    │   │   │   └── budget_suggestion.go          # 🟢 AI予算提案エンティティ (215行, 実装済み)
    │   │   ├── repository/
    │   │   │   ├── budget_repository.go          # 🟢 予算リポジトリインターフェース
    │   │   │   └── budget_suggestion_repository.go # 🟢 予算提案リポジトリインターフェース
    │   │   └── value/
    │   │       ├── period_type.go                # 🟢 期間タイプ値オブジェクト
    │   │       └── suggestion_status.go          # 🟢 提案ステータス値オブジェクト
    │   │
    │   ├── asset/                                # 📈 資産管理ドメインコンテキスト
    │   │   ├── entity/
    │   │   │   ├── asset_forecast.go             # 🟢 資産予測エンティティ
    │   │   │   └── asset_snapshot.go             # 🟢 資産スナップショットエンティティ
    │   │   ├── repository/
    │   │   │   ├── asset_forecast_repository.go  # 🟢 資産予測リポジトリインターフェース
    │   │   │   └── asset_snapshot_repository.go  # 🟢 資産スナップショットリポジトリインターフェース
    │   │   └── value/
    │   │       ├── account_breakdown.go          # 🟢 口座別内訳値オブジェクト
    │   │       └── assumptions.go                # 🟢 予測前提パラメータ値オブジェクト (187行, 実装済み)
    │   │
    │   └── notification/                         # 🔔 通知ドメインコンテキスト
    │       ├── entity/
    │       │   └── notification_settings.go      # 🟢 通知設定エンティティ (189行, 実装済み)
    │       └── repository/
    │           └── notification_settings_repository.go # 🟢 通知設定リポジトリインターフェース
    │
    ├── application/                              # 🎯 アプリケーション層（ユースケース）
    │   ├── dto/                                  # データ転送オブジェクト
    │   │   ├── auth_dto.go                       # 🟢 認証用DTO (実装済み)
    │   │   ├── user_dto.go                       # 🔴 ユーザー用DTO (未実装)
    │   │   ├── account_dto.go                    # 🔴 口座用DTO (未実装)
    │   │   ├── transaction_dto.go                # 🔴 取引用DTO (未実装)
    │   │   ├── category_dto.go                   # 🔴 カテゴリ用DTO (未実装)
    │   │   ├── budget_dto.go                     # 🔴 予算用DTO (未実装)
    │   │   ├── asset_dto.go                      # 🔴 資産管理用DTO (未実装)
    │   │   └── notification_dto.go               # 🔴 通知用DTO (未実装)
    │   ├── service/                              # アプリケーションサービス（ユースケース）
    │   │   ├── auth_service.go                   # 🟢 認証サービス (89行, 実装済み)
    │   │   ├── user_service.go                   # 🔴 ユーザー管理サービス (未実装)
    │   │   ├── account_service.go                # 🔴 口座管理サービス (未実装)
    │   │   ├── transaction_service.go            # 🔴 取引管理サービス (未実装)
    │   │   ├── category_service.go               # 🔴 カテゴリ管理サービス (未実装)
    │   │   ├── budget_service.go                 # 🔴 予算管理サービス (未実装)
    │   │   ├── asset_service.go                  # 🔴 資産管理サービス (未実装)
    │   │   ├── notification_service.go           # 🔴 通知管理サービス (未実装)
    │   │   └── report_service.go                 # 🔴 レポート生成サービス (未実装)
    │   └── transaction/
    │       └── transaction_manager.go            # 🔴 DBトランザクション管理 (未実装)
    │
    ├── infrastructure/                           # 🎯 インフラストラクチャ層（外部システム統合）
    │   ├── auth0/                                # Auth0外部認証プロバイダー統合
    │   │   ├── client.go                         # 🟢 Auth0クライアント (235行, 実装済み)
    │   │   └── middleware.go                     # 🟢 Auth0ミドルウェア (198行, 実装済み)
    │   ├── email/                                # 🔴 メール送信サービス (未実装)
    │   ├── cache/                                # 🔴 Redis キャッシュサービス (未実装)
    │   └── gorm/                                 # GORM ORM実装
    │       ├── model/                            # データベースモデル（テーブルマッピング）
    │       │   ├── base.go                       # 🟢 ベースモデル（共通フィールド）
    │       │   ├── user.go                       # 🟢 ユーザーテーブルモデル
    │       │   ├── account.go                    # 🟢 口座テーブルモデル
    │       │   ├── transaction.go                # 🟢 取引テーブルモデル
    │       │   ├── category.go                   # 🟢 カテゴリテーブルモデル
    │       │   ├── budget.go                     # 🟢 予算テーブルモデル
    │       │   ├── report.go                     # 🟢 レポート関連テーブルモデル
    │       │   └── notification.go               # 🟢 通知設定テーブルモデル
    │       └── repository/                       # リポジトリ実装（ドメインリポジトリIF実装）
    │           ├── base_repository.go            # 🔴 ベースリポジトリ実装 (未実装)
    │           ├── user_repository.go            # 🟢 ユーザーリポジトリ実装 (実装済み)
    │           ├── account_repository.go         # 🔴 口座リポジトリ実装 (未実装)
    │           ├── transaction_repository.go     # 🔴 取引リポジトリ実装 (未実装)
    │           ├── category_repository.go        # 🔴 カテゴリリポジトリ実装 (未実装)
    │           ├── budget_repository.go          # 🔴 予算リポジトリ実装 (未実装)
    │           ├── asset_repository.go           # 🔴 資産管理リポジトリ実装 (未実装)
    │           └── notification_repository.go    # 🔴 通知設定リポジトリ実装 (未実装)
    │
    └── interface/                                # 🎯 インターフェース層（HTTP API・外部I/F）
        ├── handler/                              # HTTPハンドラー（コントローラー）
        │   ├── auth_handler.go                   # 🟢 認証API (296行, 実装済み)
        │   ├── user_handler.go                   # 🔴 ユーザー管理API (未実装)
        │   ├── account_handler.go                # 🔴 口座管理API (未実装)
        │   ├── transaction_handler.go            # 🔴 取引管理API (未実装)
        │   ├── category_handler.go               # 🔴 カテゴリ管理API (未実装)
        │   ├── budget_handler.go                 # 🔴 予算管理API (未実装)
        │   ├── asset_handler.go                  # 🔴 資産管理API (未実装)
        │   ├── report_handler.go                 # 🔴 レポート生成API (未実装)
        │   └── notification_handler.go           # 🔴 通知設定API (未実装)
        ├── middleware/                           # HTTPミドルウェア（横断的関心事）
        │   ├── auth.go                           # 🟢 認証ミドルウェア (350行, 実装済み)
        │   ├── cors.go                           # 🟢 CORSミドルウェア (実装済み)
        │   ├── error_handler.go                  # 🟢 エラーハンドリングミドルウェア (実装済み)
        │   ├── logger.go                         # 🟢 リクエストログミドルウェア (実装済み)
        │   └── request_id.go                     # 🟢 リクエストIDミドルウェア (実装済み)
        ├── router/                               # HTTPルーティング
        │   ├── router.go                         # 🟢 メインルーター設定 (実装済み)
        │   ├── auth_routes.go                    # 🟢 認証関連ルート (実装済み)
        │   └── interfaces.go                     # 🟢 ルーターインターフェース (実装済み)
        └── validator/                            # 🔴 リクエストバリデーション (未実装)
```

## Clean Architecture 層分離

### 🎯 ドメイン層（Domain Layer）- 最内層
**場所**: `internal/domain/`
**責務**: ビジネスロジック・ドメインルール
**依存関係**: 他の層に依存しない純粋なビジネスロジック

**構成要素:**
- **エンティティ**: ビジネスオブジェクトと不変ルール
- **値オブジェクト**: 不変の値を表現するオブジェクト
- **リポジトリインターフェース**: データ永続化の抽象化
- **ドメインサービス**: 複数エンティティに跨るビジネスロジック
- **ドメインエラー**: ビジネスルール違反エラー

### 🎯 アプリケーション層（Application Layer）
**場所**: `internal/application/`
**責務**: ユースケース実装・アプリケーション固有ロジック
**依存関係**: ドメイン層のみに依存

**構成要素:**
- **アプリケーションサービス**: ユースケース実装
- **DTO (Data Transfer Object)**: レイヤー間データ転送
- **トランザクション管理**: 複数操作の整合性担保

### 🎯 インフラストラクチャ層（Infrastructure Layer）
**場所**: `internal/infrastructure/`
**責務**: 外部システム・技術的詳細との統合
**依存関係**: ドメイン層・アプリケーション層に依存

**構成要素:**
- **リポジトリ実装**: データベースアクセス実装
- **外部API連携**: Auth0、メール送信サービス等
- **ORMモデル**: データベーステーブルマッピング
- **キャッシュサービス**: Redis等

### 🎯 インターフェース層（Interface Layer）- 最外層
**場所**: `internal/interface/`
**責務**: 外部とのインターフェース（HTTP API、CLI等）
**依存関係**: 全ての内側の層に依存可能

**構成要素:**
- **HTTPハンドラー**: REST APIエンドポイント実装
- **ミドルウェア**: 認証、ログ、エラーハンドリング
- **ルーター**: エンドポイント定義・ルーティング
- **バリデーター**: 入力検証

## データベーススキーマ

### 実装済みテーブル（12テーブル）
1. **users** - ユーザー管理
2. **accounts** - 口座管理
3. **account_movements** - 口座残高変動履歴
4. **transactions** - 取引記録
5. **transfers** - 振替取引
6. **categories** - ユーザー固有カテゴリ
7. **category_masters** - システム標準カテゴリ（22種類）
8. **budgets** - 予算管理
9. **budget_suggestions** - AI予算提案
10. **asset_snapshots** - 資産スナップショット
11. **asset_forecasts** - 資産予測
12. **notification_settings** - 通知設定

### データベース機能
- **主キー**: UUID使用
- **タイムスタンプ**: 自動更新（トリガー）
- **外部キー制約**: 参照整合性保証
- **インデックス**: クエリパフォーマンス最適化
- **カスケード削除**: データ整合性保証

## API エンドポイント構成

### 🟢 実装済みエンドポイント
```
認証関連:
GET    /auth/login           # Auth0ログインリダイレクト
GET    /auth/callback        # Auth0認証コールバック
POST   /auth/logout          # ログアウト
GET    /auth/user            # 現在のユーザー情報取得（認証必須）
GET    /auth/check           # 認証状態チェック
POST   /auth/token           # HttpOnlyクッキートークン設定
DELETE /auth/token           # トークン削除

システム関連:
GET    /health               # ヘルスチェック
GET    /api/v1/              # API情報
```

### 🔴 未実装エンドポイント（実装予定）
```
ユーザー関連:
GET    /api/v1/users/me              # 現在のユーザー情報取得
PUT    /api/v1/users/me              # ユーザー情報更新

口座管理:
GET    /api/v1/accounts              # 口座一覧取得
POST   /api/v1/accounts              # 口座作成
GET    /api/v1/accounts/:id          # 口座詳細取得
PUT    /api/v1/accounts/:id          # 口座更新
DELETE /api/v1/accounts/:id          # 口座削除
POST   /api/v1/accounts/:id/movements # 口座残高変動記録

取引管理:
GET    /api/v1/transactions          # 取引一覧取得
POST   /api/v1/transactions          # 取引作成
GET    /api/v1/transactions/:id      # 取引詳細取得
PUT    /api/v1/transactions/:id      # 取引更新
DELETE /api/v1/transactions/:id      # 取引削除
GET    /api/v1/transactions/summary/monthly # 月次取引サマリー

カテゴリ管理:
GET    /api/v1/categories            # ユーザーカテゴリ一覧
PUT    /api/v1/categories/:id        # カテゴリ更新
DELETE /api/v1/categories/:id        # カテゴリ削除
GET    /api/v1/categories/master     # システムカテゴリマスター一覧

予算管理:
GET    /api/v1/budgets               # 予算一覧取得
POST   /api/v1/budgets               # 予算作成
GET    /api/v1/budgets/:id           # 予算詳細取得
PUT    /api/v1/budgets/:id           # 予算更新
DELETE /api/v1/budgets/:id           # 予算削除
GET    /api/v1/budgets/current       # 現在の予算取得
GET    /api/v1/budget-suggestions    # 予算提案一覧
POST   /api/v1/budget-suggestions/generate # AI予算提案生成

レポート・資産管理:
GET    /api/v1/reports/assets/snapshots # 資産スナップショット一覧
GET    /api/v1/reports/assets/forecasts/latest # 最新資産予測
GET    /api/v1/summary/monthly        # 月次サマリー

通知管理:
GET    /api/v1/notifications/settings # 通知設定取得
PUT    /api/v1/notifications/settings # 通知設定更新
```

## 実装状況サマリー

### ✅ 完全実装済み（～100行以上）
1. **認証システム** - Auth0完全統合
   - JWTトークン検証・キャッシュ
   - ログイン・ログアウト・コールバック
   - HttpOnlyクッキートークン管理
   - 認証ミドルウェア

2. **ドメインモデル** - DDD実装
   - 全7ドメインのエンティティ・値オブジェクト
   - リポジトリインターフェース定義
   - ドメインエラー定義

3. **データベース基盤**
   - 12テーブルスキーマ実装
   - マイグレーション管理（Atlas）
   - シードデータ（カテゴリマスター含む）

4. **HTTP基盤**
   - ミドルウェア（認証・CORS・ログ・エラー）
   - ルーティング設定
   - Swagger API文書化

5. **設定・ログ**
   - 環境変数ベース設定管理
   - 構造化ログ（Zap）
   - エラーハンドリング

### 🔴 未実装（優先実装）
1. **ビジネスAPIハンドラー** - 認証以外の全API
2. **アプリケーションサービス** - 認証以外の全サービス
3. **リポジトリ実装** - User以外の全リポジトリ
4. **DTO定義** - 認証以外の全DTO
5. **バリデーション** - リクエスト検証
6. **テストコード** - 単体・統合・E2Eテスト

## 開発規約・品質管理

### アーキテクチャ品質
- **依存関係チェック**: `scripts/check-architecture.sh`
- **Clean Architecture準拠**: 層間依存関係の方向性確認
- **循環依存検出**: Import文解析による検証

### コード品質
- **型安全性**: Go の型システム活用
- **エラーハンドリング**: 構造化エラー処理
- **ログ**: 構造化ログ（JSON形式）
- **設定管理**: 環境変数ベース

### データベース品質
- **マイグレーション**: Atlas CLIによるスキーマ管理
- **トランザクション**: ACID特性保証
- **インデックス**: パフォーマンス最適化
- **制約**: 参照整合性・データ整合性

## 次のステップ

### 🔥 最優先（Week 4）
1. **トランザクション管理API** - 家計簿の中核機能
2. **アカウント管理API** - 口座・残高管理
3. **ユーザー管理API完成** - プロファイル管理

### 🎯 高優先（Week 5）
4. **カテゴリー管理** - 取引分類機能
5. **予算管理** - 予算設定・進捗管理

### 📊 中優先（Week 6）
6. **レポート・資産管理** - データ分析・可視化
7. **通知機能** - アラート・リマインダー

### 🔧 基盤強化
8. **テスト実装** - 品質保証
9. **パフォーマンス最適化** - レスポンス改善
10. **セキュリティ強化** - 脆弱性対策

## 結論

FinanceTrackerバックエンドは、Clean Architectureに基づいた堅実な設計が実装されており、認証システムと主要なドメインモデルは完成度が高く、品質管理の仕組みも整備されています。

現在の構造は拡張性と保守性に優れており、新機能追加や変更に対して柔軟に対応できる設計となっています。次のステップとして、未実装のビジネスAPIとサービスの実装、テストコードの追加、そしてフロントエンドとの統合テストが重要になります。