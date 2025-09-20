# FinSight 実装計画書

## プロジェクト概要

家計簿アプリ「FinSight」の包括的実装計画です。フロントエンド（Next.js）、バックエンド（Go + Gin）、データベース（PostgreSQL）、インフラ（AWS）を段階的に構築します。

## 技術スタック

### フロントエンド
- **フレームワーク**: Next.js 14 (App Router)
- **言語**: TypeScript
- **スタイリング**: Tailwind CSS + shadcn/ui
- **状態管理**: Zustand
- **認証**: Auth0
- **アーキテクチャ**: Feature-Sliced Design (FSD)

### バックエンド
- **言語**: Go 1.21+
- **フレームワーク**: Gin
- **アーキテクチャ**: Clean Architecture (Onion Architecture)
- **ORM**: GORM
- **認証**: Auth0 JWT検証
- **API**: RESTful

### データベース
- **RDBMS**: PostgreSQL 15
- **マイグレーション**: GORM AutoMigrate + 手動SQL
- **接続**: RDS Proxy

### インフラ
- **クラウド**: AWS
- **IaC**: AWS CDK (TypeScript)
- **コンピュート**: Lambda (Go runtime)
- **API**: API Gateway
- **ホスティング**: Amplify
- **メール**: SES
- **監視**: CloudWatch + X-Ray

## 実装フェーズ

### Phase 1: インフラ基盤構築 (週1-2)

#### 1.1 AWS CDK セットアップ
**期間**: 2日間

**タスク**:
- [ ] CDKプロジェクト初期化
- [ ] VPCスタック作成
- [ ] セキュリティグループ設定
- [ ] RDS PostgreSQL構築
- [ ] Secrets Manager設定

**成果物**:
```
infrastructure/
├── bin/finsight.ts
├── lib/stacks/
│   ├── vpc-stack.ts
│   ├── database-stack.ts
│   └── certificate-stack.ts
├── config/
│   ├── dev.json
│   └── prod.json
└── package.json
```

**検証**:
- [ ] VPC作成確認
- [ ] RDS接続テスト
- [ ] Secrets Manager動作確認

#### 1.2 API Gateway + Lambda基盤
**期間**: 3日間

**タスク**:
- [ ] API Gatewayスタック作成
- [ ] Lambda関数デプロイ設定
- [ ] Auth0 Lambda Authorizer実装
- [ ] CORS設定
- [ ] CloudWatch Logs設定

**成果物**:
```
infrastructure/lib/stacks/
├── api-stack.ts
└── monitoring-stack.ts

lambda/
├── auth/
├── users/
└── shared/
```

**検証**:
- [ ] API Gateway動作確認
- [ ] Lambda関数実行確認
- [ ] JWT認証テスト
- [ ] ログ出力確認

#### 1.3 フロントエンドホスティング
**期間**: 2日間

**タスク**:
- [ ] Amplifyスタック作成
- [ ] GitHub連携設定
- [ ] 環境変数設定
- [ ] カスタムドメイン設定
- [ ] SSL証明書設定

**成果物**:
```
infrastructure/lib/stacks/
├── amplify-stack.ts
└── certificate-stack.ts
```

**検証**:
- [ ] Amplifyデプロイ確認
- [ ] GitHub連携動作確認
- [ ] カスタムドメインアクセス
- [ ] SSL証明書確認

### Phase 2: バックエンド開発 (週3-6)

#### 2.1 プロジェクト構造とドメイン層
**期間**: 4日間

**タスク**:
- [ ] Go プロジェクト初期化
- [ ] Clean Architecture構造作成
- [ ] ドメインエンティティ実装
- [ ] 値オブジェクト実装
- [ ] リポジトリインターフェース定義

**成果物**:
```
backend/
├── cmd/api/main.go
├── internal/
│   ├── domain/
│   │   ├── user/
│   │   ├── account/
│   │   ├── transaction/
│   │   ├── category/
│   │   ├── budget/
│   │   └── common/
│   ├── application/
│   ├── infrastructure/
│   └── interface/
└── go.mod
```

**検証**:
- [ ] ドメインロジック単体テスト
- [ ] 値オブジェクト検証テスト
- [ ] エンティティ関係テスト

#### 2.2 インフラストラクチャ層
**期間**: 3日間

**タスク**:
- [ ] GORM設定とモデル定義
- [ ] リポジトリ実装
- [ ] データベース接続設定
- [ ] マイグレーション実装
- [ ] JWT認証ミドルウェア

**成果物**:
```
internal/infrastructure/
├── database/connection.go
├── gorm/
│   ├── model/
│   └── repository/
└── jwt/auth.go
```

**検証**:
- [ ] データベース接続確認
- [ ] CRUD操作テスト
- [ ] マイグレーション実行
- [ ] 認証ミドルウェアテスト

#### 2.3 アプリケーション層とユースケース
**期間**: 4日間

**タスク**:
- [ ] ユースケース実装
- [ ] DTO定義
- [ ] バリデーション実装
- [ ] エラーハンドリング
- [ ] トランザクション管理

**成果物**:
```
internal/application/
├── usecase/
├── dto/
└── validation/
```

**検証**:
- [ ] ユースケース単体テスト
- [ ] トランザクション動作確認
- [ ] エラーハンドリングテスト

#### 2.4 API コントローラー
**期間**: 3日間

**タスク**:
- [ ] Ginルーター設定
- [ ] コントローラー実装
- [ ] ミドルウェア適用
- [ ] APIレスポンス標準化
- [ ] ヘルスチェックAPI

**成果物**:
```
internal/interface/
├── controller/
├── middleware/
└── router/
```

**検証**:
- [ ] 全APIエンドポイントテスト
- [ ] 認証付きリクエストテスト
- [ ] エラーレスポンステスト

### Phase 3: データベース設計・実装 (週5-6)

#### 3.1 テーブル設計
**期間**: 2日間

**タスク**:
- [ ] ER図に基づくテーブル定義
- [ ] インデックス設計
- [ ] 制約設定
- [ ] シーケンス設定

**成果物**:
```
database/
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_accounts.sql
│   ├── 003_create_categories.sql
│   ├── 004_create_transactions.sql
│   └── 005_create_budgets.sql
└── seeds/
    └── master_categories.sql
```

#### 3.2 マスターデータ投入
**期間**: 1日間

**タスク**:
- [ ] カテゴリマスター作成
- [ ] シードデータ作成
- [ ] データ投入スクリプト

### Phase 4: フロントエンド開発 (週7-12)

#### 4.1 プロジェクト構造とベース実装
**期間**: 3日間

**タスク**:
- [ ] Next.js プロジェクト初期化
- [ ] FSD構造作成
- [ ] Tailwind CSS + shadcn/ui設定
- [ ] Auth0設定
- [ ] 共通レイアウト作成

**成果物**:
```
frontend/src/
├── app/
├── page-components/
├── widgets/
├── features/
├── entities/
└── shared/
    ├── ui/
    ├── lib/
    └── config/
```

#### 4.2 認証・ユーザー管理機能
**期間**: 3日間

**タスク**:
- [ ] Auth0 Provider設定
- [ ] ログイン/ログアウト機能
- [ ] プロファイル管理
- [ ] 認証ガード実装

**成果物**:
```
src/features/auth/
├── auth-login/
├── auth-logout/
└── auth-profile/

src/entities/user/
├── api/
├── model/
└── ui/
```

#### 4.3 ダッシュボード画面
**期間**: 4日間

**タスク**:
- [ ] レスポンシブレイアウト
- [ ] 収支サマリーカード
- [ ] 月次グラフ
- [ ] 最近の取引一覧
- [ ] 予算達成状況

**成果物**:
```
src/page-components/dashboard/
src/widgets/dashboard/
├── monthly-summary/
├── transaction-chart/
├── recent-transactions/
└── budget-progress/
```

#### 4.4 取引管理機能
**期間**: 5日間

**タスク**:
- [ ] 取引一覧表示
- [ ] 取引作成フォーム
- [ ] 取引編集機能
- [ ] フィルタリング
- [ ] ページネーション
- [ ] CSV/PDF エクスポート

**成果物**:
```
src/page-components/transactions/
src/widgets/transactions/
├── transaction-list/
├── transaction-form/
├── transaction-filters/
└── transaction-export/

src/features/transaction/
├── transaction-get/
├── transaction-create/
├── transaction-update/
└── transaction-delete/
```

#### 4.5 予算管理機能
**期間**: 4日間

**タスク**:
- [ ] 予算設定画面
- [ ] 予算対実績表示
- [ ] 進捗バー
- [ ] アラート表示
- [ ] AI予算提案機能

**成果物**:
```
src/page-components/budget/
src/widgets/budget/
├── budget-setup/
├── budget-progress/
├── budget-alerts/
└── budget-suggestions/
```

#### 4.6 口座・資産管理機能
**期間**: 4日間

**タスク**:
- [ ] 口座一覧表示
- [ ] 口座作成・編集
- [ ] 資産推移グラフ
- [ ] 資産予測表示

**成果物**:
```
src/page-components/assets/
src/widgets/assets/
├── account-list/
├── account-form/
├── asset-chart/
└── asset-forecast/
```

#### 4.7 設定・レポート機能
**期間**: 3日間

**タスク**:
- [ ] ユーザー設定画面
- [ ] 通知設定
- [ ] カテゴリ管理
- [ ] レポート生成
- [ ] メール送信設定

**成果物**:
```
src/page-components/settings/
src/widgets/settings/
├── user-settings/
├── notification-settings/
├── category-settings/
└── report-settings/
```

### Phase 5: 統合・テスト (週13-14)

#### 5.1 E2E テスト
**期間**: 3日間

**タスク**:
- [ ] Playwright設定
- [ ] 主要フローE2Eテスト作成
- [ ] 認証フローテスト
- [ ] データ操作テスト

#### 5.2 パフォーマンス最適化
**期間**: 2日間

**タスク**:
- [ ] フロントエンドバンドル最適化
- [ ] API レスポンス最適化
- [ ] データベースクエリ最適化

#### 5.3 セキュリティ検証
**期間**: 2日間

**タスク**:
- [ ] 認証・認可テスト
- [ ] XSS/CSRF 対策確認
- [ ] SQLインジェクション対策確認
- [ ] 機密情報漏洩チェック

### Phase 6: デプロイ・本番化 (週15-16)

#### 6.1 CI/CDパイプライン
**期間**: 3日間

**タスク**:
- [ ] GitHub Actions設定
- [ ] 自動テスト実行
- [ ] 自動デプロイ設定
- [ ] 環境別デプロイ

#### 6.2 本番環境構築
**期間**: 2日間

**タスク**:
- [ ] 本番環境設定
- [ ] ドメイン設定
- [ ] SSL証明書設定
- [ ] 監視設定

#### 6.3 運用準備
**期間**: 2日間

**タスク**:
- [ ] ログ監視設定
- [ ] アラート設定
- [ ] バックアップ設定
- [ ] 運用ドキュメント作成

## 各週の詳細タスク

### 第1週: インフラストラクチャ基盤

#### Day 1-2: AWS CDK基盤構築
```bash
# 1. CDKプロジェクト初期化
mkdir infrastructure && cd infrastructure
npm init -y
npm install -g aws-cdk
cdk init app --language typescript

# 2. 依存関係インストール
npm install @aws-cdk/aws-ec2 @aws-cdk/aws-rds @aws-cdk/aws-secretsmanager

# 3. VPCスタック作成
touch lib/stacks/vpc-stack.ts
touch lib/stacks/database-stack.ts
```

#### Day 3-5: API Gateway + Lambda設定
```bash
# 1. APIスタック作成
touch lib/stacks/api-stack.ts
mkdir -p lambda/{auth,users,transactions,budgets}

# 2. Lambda関数用Go環境セットアップ
cd lambda/auth && go mod init auth
cd ../users && go mod init users
```

#### Day 6-7: フロントエンドホスティング
```bash
# 1. Next.jsプロジェクト作成
npx create-next-app@latest frontend --typescript --tailwind --app
cd frontend && npm install @auth0/nextjs-auth0

# 2. Amplifyスタック作成
touch infrastructure/lib/stacks/amplify-stack.ts
```

### 第2週: インフラ完成とバックエンド開始

#### Day 8-10: SES設定と監視
```bash
# 1. SESスタック作成
touch infrastructure/lib/stacks/ses-stack.ts

# 2. 監視スタック作成
touch infrastructure/lib/stacks/monitoring-stack.ts

# 3. 全スタックデプロイ
cdk deploy --all --context env=dev
```

#### Day 11-14: バックエンドプロジェクト構築
```bash
# 1. Go プロジェクト初期化
mkdir backend && cd backend
go mod init github.com/your-org/finsight-backend

# 2. ディレクトリ構造作成
mkdir -p {cmd/api,internal/{domain,application,infrastructure,interface}}
mkdir -p internal/domain/{user,account,transaction,category,budget}/entity
mkdir -p internal/domain/{user,account,transaction,category,budget}/repository
```

### 第3-6週: バックエンド開発詳細

#### ドメイン層実装例
```go
// internal/domain/user/entity/user.go
package entity

import (
    "time"
    "github.com/google/uuid"
    "github.com/your-org/finsight-backend/internal/domain/user/value"
)

type User struct {
    ID           uuid.UUID
    Auth0UserID  string
    Email        value.Email
    Name         string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func NewUser(auth0UserID, email, name string) (*User, error) {
    emailVO, err := value.NewEmail(email)
    if err != nil {
        return nil, err
    }
    
    return &User{
        ID:          uuid.New(),
        Auth0UserID: auth0UserID,
        Email:       emailVO,
        Name:        name,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }, nil
}
```

### 第7-12週: フロントエンド開発詳細

#### FSD構造実装例
```typescript
// src/entities/user/api/userApi.ts
export const userApi = {
  getCurrentUser: async (): Promise<User> => {
    const response = await fetch('/api/users/me', {
      headers: { Authorization: `Bearer ${getToken()}` }
    });
    return response.json();
  }
};

// src/features/auth/auth-login/ui/LoginButton.tsx
export const LoginButton = () => {
  const { loginWithRedirect } = useAuth0();
  
  return (
    <Button onClick={() => loginWithRedirect()}>
      ログイン
    </Button>
  );
};
```

## 成果物チェックリスト

### インフラ
- [ ] VPC、サブネット、セキュリティグループ
- [ ] RDS PostgreSQL（Multi-AZ本番）
- [ ] API Gateway + Lambda関数群
- [ ] Amplify ホスティング
- [ ] SES メール送信
- [ ] CloudWatch 監視・アラート
- [ ] X-Ray 分散トレーシング

### バックエンド
- [ ] Clean Architecture構造
- [ ] 7つのドメインコンテキスト実装
- [ ] 全APIエンドポイント（50+個）
- [ ] JWT認証・認可
- [ ] エラーハンドリング
- [ ] 単体・統合テスト

### フロントエンド
- [ ] FSD アーキテクチャ
- [ ] 8つの主要画面
- [ ] レスポンシブデザイン
- [ ] Auth0 統合
- [ ] 状態管理（Zustand）
- [ ] E2E テスト

### データベース
- [ ] 12テーブル設計・実装
- [ ] インデックス・制約設定
- [ ] マスターデータ投入
- [ ] マイグレーション管理

## リスク管理

### 技術的リスク
1. **Auth0統合の複雑性**
   - 軽減策: 早期プロトタイプ作成、ドキュメント詳細調査

2. **AWS Lambda Cold Start**
   - 軽減策: 適切なメモリ設定、Provisioned Concurrency検討

3. **フロントエンドバンドルサイズ**
   - 軽減策: Code Splitting、Dynamic Import活用

### プロジェクト管理リスク
1. **スコープクリープ**
   - 軽減策: 各フェーズの成果物明確化、変更管理

2. **技術的負債の蓄積**
   - 軽減策: コードレビュー実施、リファクタリング時間確保

## 品質保証

### コード品質
- [ ] TypeScript/Go の型安全性
- [ ] ESLint/golangci-lint
- [ ] Prettier/gofmt
- [ ] コードレビュー

### テスト戦略
- [ ] 単体テスト（70%以上カバレッジ）
- [ ] 統合テスト
- [ ] E2Eテスト（主要フロー）
- [ ] パフォーマンステスト

### セキュリティ
- [ ] OWASP Top 10対策
- [ ] 依存関係脆弱性スキャン
- [ ] セキュリティヘッダー設定
- [ ] 機密情報管理

## 運用計画

### 監視・ログ
- [ ] アプリケーションメトリクス
- [ ] インフラメトリクス
- [ ] エラーアラート
- [ ] パフォーマンス監視

### バックアップ・復旧
- [ ] RDS自動バックアップ
- [ ] ポイントインタイム復旧
- [ ] 災害復旧計画

### スケーリング
- [ ] Lambda自動スケーリング
- [ ] RDS垂直スケーリング
- [ ] CDN活用

このロードマップに従って実装を進めることで、16週間（約4ヶ月）でFinSightアプリケーションを完成させることができます。