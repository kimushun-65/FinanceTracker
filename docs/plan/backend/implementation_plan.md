# バックエンド実装計画書（AWS Lambda版）

## 概要
本ドキュメントは、FinanceTrackerプロジェクトのバックエンドをAWS Lambdaを使用して実装するための詳細な計画書です。
オニオンアーキテクチャを採用し、保守性・テスタビリティ・拡張性を確保します。

## 現在のインフラストラクチャ状況

### 既存のCDK基盤
- **CDKスタック**: VPC、Database、API、Amplify、SES、Monitoring、Security
- **データベース**: PostgreSQL 15 (RDS)、マルチAZ構成
- **Lambda**: Go言語 (PROVIDED_AL2ランタイム)
- **API Gateway**: REST API、Auth0カスタムオーソライザー
- **認証**: Auth0統合

### 既存のLambda関数
- users、accounts、transactions、categories、budgets、reports、auth、notifications
- 各関数は基本的なプレースホルダー実装のみ
- 共通ユーティリティ (DB接続、レスポンス処理) は実装済み

## 実装フェーズ

### フェーズ1: プロジェクト構造の再編成（Day 1-2）

#### 1.1 既存構造のリファクタリング
現在の構造:
```
cdk/lambda/
├── [function-name]/
│   ├── main.go
│   ├── go.mod
│   └── bootstrap
└── common/
    ├── db/
    ├── models/
    └── utils/
```

新しい構造（オニオンアーキテクチャ）:
```
backend/
├── cmd/
│   └── lambda/
│       ├── users/
│       ├── accounts/
│       ├── transactions/
│       ├── categories/
│       ├── budgets/
│       ├── reports/
│       ├── auth/
│       └── notifications/
├── internal/
│   ├── domain/                # ドメイン層
│   │   ├── user/
│   │   ├── account/
│   │   ├── transaction/
│   │   ├── category/
│   │   ├── budget/
│   │   └── common/
│   ├── application/           # アプリケーション層
│   │   ├── dto/
│   │   ├── usecase/
│   │   └── transaction/
│   ├── infrastructure/        # インフラストラクチャ層
│   │   ├── postgres/
│   │   ├── auth0/
│   │   └── aws/
│   └── interface/             # インターフェース層
│       └── lambda/
├── pkg/                       # 共有パッケージ
│   ├── config/
│   └── errors/
├── scripts/                   # ビルド・デプロイスクリプト
├── go.mod
├── go.sum
└── Makefile
```

#### 1.2 依存関係の整理
既存のパッケージに加えて:
- GORM v2 (gorm.io/gorm)
- GORM PostgreSQLドライバー (gorm.io/driver/postgres)
- go-playground/validator/v10
- uber-go/zap (ロギング)

#### 1.3 CDKとの統合
- build-go-functions.sh スクリプトの更新
- Lambda関数のビルドパスをbackend/に変更

### フェーズ2: ドメイン層実装（Day 3-5）

#### 2.1 共通ドメイン要素
```go
// internal/domain/common/
├── errors.go              # 共通エラー定義
├── base_entity.go         # 基本エンティティ（UUID, timestamps）
├── value/
│   ├── money.go           # 金額値オブジェクト
│   ├── email.go           # メールアドレス値オブジェクト
│   └── hex_color.go       # HEX形式の色値オブジェクト
└── base_repository.go     # 基本リポジトリインターフェース
```

#### 2.2 ユーザーコンテキスト（user）
```go
// internal/domain/user/
├── entity/
│   └── user.go           # ユーザーエンティティ
├── value/
│   ├── user_id.go        # ユーザーID値オブジェクト
│   └── auth0_id.go       # Auth0 ID値オブジェクト
├── repository/
│   └── user_repository.go
└── errors.go             # ユーザードメインエラー
```

#### 2.3 アカウントコンテキスト（account）
```go
// internal/domain/account/
├── entity/
│   └── account.go        # アカウントエンティティ
├── value/
│   ├── account_type.go   # アカウントタイプ値オブジェクト
│   └── balance.go        # 残高値オブジェクト
├── repository/
│   └── account_repository.go
└── errors.go
```

#### 2.4 取引コンテキスト（transaction）
```go
// internal/domain/transaction/
├── entity/
│   └── transaction.go    # 取引エンティティ
├── value/
│   ├── transaction_type.go   # 取引タイプ値オブジェクト
│   └── amount.go             # 金額値オブジェクト
├── repository/
│   └── transaction_repository.go
└── errors.go
```

#### 2.5 カテゴリコンテキスト（category）
```go
// internal/domain/category/
├── entity/
│   ├── category.go           # カテゴリエンティティ
│   └── category_master.go    # カテゴリマスタエンティティ
├── value/
│   └── category_type.go      # カテゴリタイプ値オブジェクト
├── repository/
│   └── category_repository.go
└── errors.go
```

#### 2.6 予算コンテキスト（budget）
```go
// internal/domain/budget/
├── entity/
│   ├── budget.go                # 予算エンティティ
│   └── budget_suggestion.go     # AI予算提案エンティティ
├── value/
│   └── suggestion_status.go     # 提案ステータス値オブジェクト
├── repository/
│   ├── budget_repository.go
│   └── budget_suggestion_repository.go
└── errors.go
```

#### 2.7 資産管理コンテキスト（asset）
```go
// internal/domain/asset/
├── entity/
│   ├── asset_snapshot.go        # 資産スナップショットエンティティ
│   └── asset_forecast.go        # 資産予測エンティティ
├── value/
│   ├── forecast_method.go       # 予測手法値オブジェクト
│   └── assumptions.go           # 予測前提条件値オブジェクト
├── repository/
│   ├── asset_snapshot_repository.go
│   └── asset_forecast_repository.go
└── errors.go
```

#### 2.8 通知管理コンテキスト（notification）
```go
// internal/domain/notification/
├── entity/
│   └── notification_settings.go # 通知設定エンティティ
├── repository/
│   └── notification_settings_repository.go
└── errors.go
```

### フェーズ3: アプリケーション層実装（Day 6-8）

#### 3.1 DTOの定義
```go
// internal/application/dto/
├── user/
│   └── user_dto.go
├── account/
│   └── account_dto.go
├── transaction/
│   └── transaction_dto.go
├── category/
│   └── category_dto.go
├── budget/
│   └── budget_dto.go
└── report/
    └── report_dto.go
```

#### 3.2 ユースケース実装
```go
// internal/application/usecase/
├── user/
│   ├── user_service.go       # ユーザー管理
│   └── auth_service.go       # Auth0認証サービス
├── account/
│   ├── account_service.go    # アカウント管理
│   └── movement_service.go   # 残高調整サービス
├── transaction/
│   └── transaction_service.go # 取引管理
├── category/
│   ├── category_service.go   # カテゴリ管理
│   └── category_master_service.go # マスターカテゴリ管理
├── budget/
│   ├── budget_service.go     # 予算管理
│   └── budget_suggestion_service.go # AI予算提案
├── asset/
│   ├── asset_snapshot_service.go # 資産スナップショット
│   └── asset_forecast_service.go # 資産予測
├── notification/
│   └── notification_service.go # 通知設定管理
└── report/
    ├── report_service.go     # レポート生成
    ├── summary_service.go    # サマリー生成
    └── dashboard_service.go  # ダッシュボード用データ
```

#### 3.3 トランザクション管理
```go
// internal/application/transaction/
└── manager.go               # トランザクションマネージャー
```

### フェーズ4: インフラストラクチャ層実装（Day 9-11）

#### 4.1 データベース接続
```go
// internal/infrastructure/database/
├── connection.go           # GORM DB接続管理
├── transaction.go          # トランザクション管理
└── migration/              # GORMマイグレーション
```

#### 4.2 GORMモデル定義
```go
// internal/infrastructure/gorm/model/
├── user.go
├── account.go
├── transaction.go
├── category.go
├── budget.go
├── asset_snapshot.go
└── notification_setting.go
```

#### 4.3 リポジトリ実装
```go
// internal/infrastructure/gorm/repository/
├── base_repository.go      # GORM基本リポジトリ実装
├── user_repository.go
├── account_repository.go
├── transaction_repository.go
├── category_repository.go
└── budget_repository.go
```

#### 4.4 外部サービス実装
```go
// internal/infrastructure/
├── auth0/
│   └── auth0_service.go    # Auth0認証実装
├── aws/
│   ├── s3_service.go       # S3連携
│   ├── ses_service.go      # メール送信
│   └── secrets_manager.go  # シークレット管理
└── cache/
    └── elasticache_service.go # ElastiCache実装
```

### フェーズ5: Lambda インターフェース層実装（Day 12-14）

#### 5.1 Lambda関数ハンドラー
```go
// internal/interface/lambda/
├── users/
│   └── handler.go               # GET,PUT /users/me
├── accounts/
│   └── handler.go               # CRUD /accounts, POST /accounts/{id}/movements
├── transactions/
│   └── handler.go               # CRUD /transactions, GET /transactions/summary/*
├── categories/
│   └── handler.go               # CRUD /categories, GET /categories/master
├── budgets/
│   └── handler.go               # CRUD /budgets, GET /budgets/current, /budget-suggestions/*
├── reports/
│   └── handler.go               # GET /reports/*, /summary/*, /assets/*
├── auth/
│   └── handler.go               # POST /auth/callback, /auth/sync
└── notifications/
    └── handler.go               # GET,PUT /notifications/settings
```

#### 5.1.1 APIエンドポイントマッピング
| Lambda関数 | HTTPメソッド | エンドポイント | 説明 |
|-----------|------------|------------|------|
| users | GET | /users/me | ユーザー情報取得 |
| users | PUT | /users/me | ユーザー情報更新 |
| accounts | GET | /accounts | 口座一覧取得 |
| accounts | POST | /accounts | 口座作成 |
| accounts | PUT | /accounts/{id} | 口座更新 |
| accounts | DELETE | /accounts/{id} | 口座削除 |
| accounts | POST | /accounts/{id}/movements | 残高調整 |
| transactions | GET | /transactions | 取引一覧取得 |
| transactions | POST | /transactions | 取引作成 |
| transactions | PUT | /transactions/{id} | 取引更新 |
| transactions | DELETE | /transactions/{id} | 取引削除 |
| transactions | GET | /transactions/summary/monthly | 月次サマリー |
| categories | GET | /categories | カテゴリ一覧取得 |
| categories | PUT | /categories/{id} | カテゴリ更新 |
| categories | DELETE | /categories/{id} | カテゴリ削除 |
| categories | GET | /categories/master | マスターカテゴリ一覧 |
| budgets | GET | /budgets | 予算一覧取得 |
| budgets | POST | /budgets | 予算設定 |
| budgets | PUT | /budgets/{id} | 予算更新 |
| budgets | GET | /budgets/current | 今月の予算取得 |
| budgets | GET | /budget-suggestions | AI予算提案取得 |
| budgets | POST | /budget-suggestions/generate | AI予算提案生成 |
| reports | GET | /assets/snapshots | 資産スナップショット |
| reports | GET | /assets/forecasts/latest | 最新資産予測 |
| reports | GET | /summary/monthly | 月別収支サマリー |
| auth | POST | /auth/callback | Auth0コールバック |
| auth | POST | /auth/sync | ユーザー同期 |
| notifications | GET | /notifications/settings | 通知設定取得 |
| notifications | PUT | /notifications/settings | 通知設定更新 |

#### 5.2 共通処理
```go
// pkg/lambda/
├── response.go             # レスポンスヘルパー
├── error.go                # エラーハンドリング
├── context.go              # コンテキスト管理
└── middleware.go           # ミドルウェア
```

#### 5.3 各Lambda関数のエントリポイント
```go
// cmd/lambda/users/main.go
package main

import (
    "github.com/aws/aws-lambda-go/lambda"
    handler "finance-tracker/internal/interface/lambda/users"
    "finance-tracker/internal/infrastructure/postgres"
    // DI設定
)

func main() {
    // 依存関係の初期化
    db := postgres.NewConnection()
    userRepo := postgres.NewUserRepository(db)
    userService := usecase.NewUserService(userRepo)
    h := handler.NewHandler(userService)
    
    lambda.Start(h.Handle)
}
```

### フェーズ6: デプロイ設定の更新（Day 15-16）

#### 6.1 CDKの更新
```typescript
// cdk/lib/stacks/api-stack.tsの更新
// Lambda関数のコードパスをbackend/に変更
code: Code.fromAsset(`backend/build/${config.path}`),
```

#### 6.2 ビルドスクリプトの更新
```bash
# backend/scripts/build.sh
#!/bin/bash

FUNCTIONS=(users accounts transactions categories budgets reports auth notifications)

for func in "${FUNCTIONS[@]}"; do
    echo "Building $func function..."
    cd cmd/lambda/$func
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
    cp bootstrap ../../../build/$func/
    cd -
done

# CDKビルドスクリプトの更新
# cdk/lambda/build-go-functions.shの修正
```

### フェーズ7: テスト実装（Day 17-18）

#### 7.1 単体テスト
```go
// ドメイン層テスト
// internal/domain/auth/entity/user_test.go

// アプリケーション層テスト
// internal/application/usecase/auth/login_service_test.go

// インフラストラクチャ層テスト
// internal/infrastructure/gorm/repository/user_repository_test.go
```

#### 7.2 統合テスト
```go
// Lambda関数の統合テスト
// tests/integration/lambda/auth_test.go
```

#### 7.3 E2Eテスト
```go
// エンドツーエンドテスト
// tests/e2e/auth_flow_test.go
```

### フェーズ8: デプロイと運用準備（Day 19-20）

#### 8.1 CI/CDパイプライン設定
- GitHub Actionsワークフロー作成
- 自動テスト
- 自動ビルド
- 自動デプロイ

#### 8.2 監視・ログ設定
- CloudWatchログ設定
- X-Ray トレーシング
- アラート設定

#### 8.3 ドキュメント作成
- API仕様書
- 開発者向けガイド
- 運用手順書

## 技術スタック

### コア技術
- **言語**: Go 1.21+
- **Lambdaランタイム**: PROVIDED_AL2 (Amazon Linux 2)
- **フレームワーク**: AWS Lambda for Go

### 主要ライブラリ
- **ORM**: GORM v2
- **認証**: Auth0 SDK
- **バリデーション**: go-playground/validator/v10
- **ログ**: uber-go/zap
- **AWS SDK**: AWS SDK for Go v2

### インフラ（既存CDKスタック）
- **コンピューティング**: AWS Lambda
- **API Gateway**: AWS API Gateway (REST API)
- **データベース**: Amazon RDS (PostgreSQL 15)
- **キャッシュ**: ElastiCache for Redis
- **ストレージ**: Amazon S3
- **メール**: Amazon SES
- **認証**: Auth0

## 考慮事項

### Lambda特有の制約
1. **コールドスタート対策**
   - 軽量な関数設計
   - プロビジョニング済み同時実行の検討
   - 初期化処理の最適化

2. **タイムアウト設定**
   - API Gateway: 29秒
   - Lambda: 15分（適切に設定）

3. **メモリ割り当て**
   - 各関数に適切なメモリサイズ設定
   - メモリ使用量の監視

4. **同時実行数**
   - アカウントレベルの制限
   - 関数レベルの予約

### セキュリティ
1. **IAMロール**
   - 最小権限の原則
   - 関数ごとの適切なロール設定

2. **環境変数の暗号化**
   - KMSによる暗号化
   - SecretsManagerの活用

3. **VPC設定**
   - 必要に応じてVPC内で実行
   - セキュリティグループの適切な設定

### パフォーマンス最適化
1. **データベース接続**
   - コネクションプーリング
   - RDS Proxyの活用

2. **キャッシュ戦略**
   - Redisによる頻繁アクセスデータのキャッシュ
   - Lambda内メモリキャッシュ

3. **非同期処理**
   - SQSを使用した非同期ジョブ
   - Step Functionsによるワークフロー

## ドメインモデル詳細

### エンティティ関係図
```
User (1) --- (*) Account
User (1) --- (*) Category  
User (1) --- (*) Budget
User (1) --- (*) Transaction
User (1) --- (*) AssetSnapshot
User (1) --- (*) AssetForecast
User (1) --- (*) BudgetSuggestion
User (1) --- (1) NotificationSettings
Account (1) --- (*) AccountMovement
Category (1) --- (*) Transaction
Category (1) --- (*) Budget
Category (1) --- (*) BudgetSuggestion
CategoryMaster (1) --- (*) Category
CategoryMaster (0..1) --- (*) CategoryMaster (self-reference)
```

### ビジネスルール
1. **アカウント残高**: AccountMovementを通じてのみ変更
2. **取引記録**: 過去1年以上前の取引は編集不可
3. **カテゴリ階層**: 2階層まで（親・子）
4. **予算**: ユーザー、カテゴリ、年月の組み合わせは一意
5. **資産スナップショット**: 日次で自動作成、過去のスナップショットは変更不可
6. **通知設定**: ユーザーごとに1つのみ
7. **メールアドレス**: RFC 5322準拠、最大255文字
8. **金額**: 整数管理（最小単位：円）

### ドメインイベント
1. **UserCreated**: 各コンテキストでユーザー参照を作成、通知設定初期化
2. **TransactionCreated**: 予算消化計算、口座残高更新
3. **AccountMovementCreated**: 口座残高更新
4. **BudgetUpdated**: 予算達成率の再計算
5. **BudgetExceeded**: 予算超過通知トリガー
6. **AssetSnapshotCreated**: 資産予測の再計算トリガー
7. **AccountBalanceChanged**: 資産スナップショットの更新

## 成功指標

### 技術的指標
- API応答時間: 95%が300ms以内
- エラー率: 0.1%以下
- 可用性: 99.9%以上
- テストカバレッジ: 80%以上

### ビジネス指標
- 開発期間: 18営業日以内
- 既存インフラの活用
- スケーラビリティ: Lambdaの自動スケーリング

## リスクと対策

### 技術的リスク
1. **コールドスタート問題**
   - 対策: 軽量な関数設計、共通初期化処理の最適化

2. **データベース接続プール**
   - 対策: RDS Proxyの活用、接続プール管理

3. **デバッグの困難さ**
   - 対策: 構造化ログ、X-Rayトレーシング

### 運用リスク
1. **コスト管理**
   - 対策: CloudWatchモニタリング、予算アラート

2. **CDKとの統合**
   - 対策: ビルドプロセスの自動化

## まとめ
この実装計画に従って、既存のCDKインフラを活用しながら、FinanceTrackerのバックエンドをオニオンアーキテクチャで再構築します。現在のプレースホルダー実装から、保守性・テスタビリティ・拡張性に優れたプロダクショングレードのシステムへと成長させます。

## 実装優先順位

1. **フェーズ1-2**: 構造の再編成とドメイン層の実装
2. **フェーズ3-4**: アプリケーション層とインフラ層
3. **フェーズ5**: Lambdaハンドラーの実装
4. **フェーズ6-7**: CDK統合とテスト
5. **フェーズ8**: デプロイと運用準備