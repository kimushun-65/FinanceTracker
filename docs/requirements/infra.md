# インフラストラクチャ要件定義書

## 概要

FinSight のインフラストラクチャは、AWS のサーバーレスアーキテクチャを採用し、スケーラビリティ、可用性、コスト効率を最大化します。EC2 インスタンスは使用せず、AWS Lambda を中心としたフルマネージドサービスで構成します。インフラストラクチャの管理には AWS CDK（Cloud Development Kit）を使用し、Infrastructure as Code を実現します。

## アーキテクチャ概要

### 基本方針

- **サーバーレスファースト**: EC2 不使用、Lambda 中心のアーキテクチャ
- **マネージドサービス活用**: 運用負荷を最小化
- **高可用性**: マルチ AZ 構成で 99.9%の可用性を確保
- **セキュリティ重視**: 最小権限の原則、暗号化の徹底
- **コスト最適化**: 使用量に応じた従量課金モデル

## AWS サービス構成

### 1. API 層

#### Amazon API Gateway (REST API)

- **用途**: RESTful API のエンドポイント管理
- **設定**:
  - リージョナルエンドポイント
  - カスタムドメイン設定（api.finsight.com）
  - SSL/TLS 証明書（AWS Certificate Manager）
  - API キー管理とレート制限
  - CORS 設定

#### AWS Lambda

- **用途**: API リクエストの処理
- **ランタイム**: Go 1.x (provided.al2)
- **設定**:
  - メモリ: 512MB〜1024MB（エンドポイントごとに最適化）
  - タイムアウト: 30 秒
  - 同時実行数: 1000（予約済み同時実行数: 100）
  - 環境変数による設定管理

### 2. 認証・認可

#### Auth0 Integration

- **用途**: ユーザー認証（Google OAuth）
- **連携方法**:
  - Lambda Authorizer で JWT 検証
  - Auth0 の JWKS エンドポイントから公開鍵取得

### 3. データ層

#### Amazon RDS for PostgreSQL

- **用途**: メインデータベース（全てのアプリケーションデータを管理）
- **構成**:
  - エンジンバージョン: PostgreSQL 15
  - インスタンスクラス: db.t3.micro（初期フェーズ）
  - Single-AZ 配置（初期フェーズ）
  - 自動バックアップ（7 日間保持）
  - 暗号化: AWS KMS
- **管理対象データ**:
  - ユーザー情報
  - 取引データ
  - カテゴリ
  - 予算
  - 口座情報
  - 資産スナップショット
  - 通知設定

### 4. フロントエンドホスティング

#### AWS Amplify

- **用途**: React アプリケーションのホスティングとデプロイ
- **設定**:
  - GitHub リポジトリとの自動連携
  - ブランチベースのデプロイ（main, develop）
  - 環境変数管理
  - カスタムドメイン設定（www.finsight.com）
  - SSL/TLS 証明書自動管理
- **ビルド設定**:
  - Node.js バージョン: 18.x
  - ビルドコマンド: `npm run build`
  - 出力ディレクトリ: `build/`

### 5. CDN・配信

#### Amazon CloudFront

- **用途**: API Gateway のキャッシュとセキュリティ強化
- **設定**:
  - オリジン: API Gateway（API）
  - キャッシュ動作設定
  - WAF 統合
  - カスタムエラーページ

### 6. メール送信

#### Amazon SES

- **用途**: トランザクションメール送信
- **設定**:
  - 送信レート: 1 メール/秒（初期フェーズ）
  - バウンス・苦情処理
  - DKIM 設定

### 7. 監視・ログ

#### Amazon CloudWatch

- **用途**:
  - メトリクス収集
  - ログ集約
  - アラーム設定
- **ログストリーム**:
  - `/aws/lambda/finsight-api`
  - `/aws/rds/instance/finsight-db`

#### AWS X-Ray

- **用途**: 分散トレーシング
- **設定**:
  - Lambda 関数のトレース
  - サービスマップ生成

### 8. セキュリティ

#### AWS WAF

- **用途**: Web アプリケーションファイアウォール
- **ルール**:
  - SQL インジェクション防御
  - XSS 防御
  - レート制限

#### AWS Secrets Manager

- **用途**: 機密情報管理
- **管理対象**:
  - データベース認証情報
  - API キー
  - Auth0 設定

#### AWS KMS

- **用途**: 暗号化キー管理
- **対象**:
  - RDS データベース
  - Secrets Manager

### 9. ネットワーク

#### Amazon VPC

- **構成**:
  - CIDR: 10.0.0.0/16
  - パブリックサブネット: 2 AZ
  - プライベートサブネット: 2 AZ
  - NAT ゲートウェイ: 1 つ（初期フェーズ）

#### セキュリティグループ

- **Lambda 用**: アウトバウンドのみ許可
- **RDS 用**: Lambda SG からのみ 5432 番ポート許可

## CDK構成

### プロジェクト構造
```
infrastructure/
├── bin/
│   └── finsight.ts              # CDKアプリケーションエントリーポイント
├── lib/
│   ├── stacks/
│   │   ├── api-stack.ts         # API Gateway + Lambda
│   │   ├── database-stack.ts    # RDS
│   │   ├── frontend-stack.ts    # Amplify設定
│   │   ├── monitoring-stack.ts  # CloudWatch + X-Ray
│   │   └── security-stack.ts    # WAF + Secrets Manager
│   └── constructs/
│       ├── lambda-api.ts        # Lambda関数構成
│       └── vpc-resources.ts     # VPCリソース
├── config/
│   ├── dev.json                 # 開発環境設定
│   ├── staging.json             # ステージング環境設定
│   └── prod.json                # 本番環境設定
├── cdk.json                     # CDK設定
└── package.json
```

### スタック構成
1. **VPCスタック**: ネットワーク基盤
2. **データベーススタック**: RDS
3. **APIスタック**: API Gateway、Lambda、Authorizer、CloudFront
4. **フロントエンドスタック**: Amplify
5. **監視スタック**: CloudWatch、X-Ray、アラーム
6. **セキュリティスタック**: WAF、Secrets Manager、KMS

### デプロイメント戦略
```bash
# 開発環境
cdk deploy --all --context env=dev

# ステージング環境
cdk deploy --all --context env=staging

# 本番環境
cdk deploy --all --context env=prod --require-approval broadening
```

## Lambda 関数設計

### 関数分割戦略

- **認証系**: `/auth/*` エンドポイント用
- **ユーザー系**: `/users/*` エンドポイント用
- **取引系**: `/transactions/*` エンドポイント用
- **予算系**: `/budgets/*` エンドポイント用
- **資産系**: `/accounts/*`, `/assets/*` エンドポイント用
- **カテゴリ系**: `/categories/*` エンドポイント用
- **レポート系**: `/reports/*` エンドポイント用
- **通知系**: `/notifications/*` エンドポイント用

### コールドスタート対策

- **Provisioned Concurrency**: 重要なエンドポイントに設定
- **関数の最適化**:
  - 初期化処理の最小化
  - 依存関係の削減
  - コンテナイメージではなく ZIP デプロイ

## データベース接続管理

### RDS Proxy

- **用途**: コネクションプーリング
- **設定**:
  - 最大接続数: 20（初期フェーズ）
  - アイドルタイムアウト: 30 分
  - 認証: IAM データベース認証

## CI/CD パイプライン

### フロントエンド（Amplify）
- **自動デプロイ**: GitHub リポジトリの push で自動ビルド・デプロイ
- **プレビュー環境**: PR ごとに独立したプレビュー環境を自動生成
- **ロールバック**: ワンクリックで前のバージョンに戻せる

### バックエンド（GitHub Actions）

```yaml
# .github/workflows/deploy-backend.yml
name: Deploy Backend to AWS
on:
  push:
    branches:
      - main
      - develop
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - name: Install CDK
        run: npm install -g aws-cdk
      - name: Deploy
        run: cdk deploy --all --require-approval never
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
```

## コスト最適化

### 見積もり（月額） - 10人/日利用想定

- **Lambda**: $0-5（9,000リクエスト/月想定）
- **API Gateway**: $1（10,000 APIコール想定）  
- **RDS**: $15（db.t3.micro、Single-AZ）
- **Amplify**: $5（ビルド時間・ホスティング）
- **CloudFront**: $1（1GB転送想定）
- **その他**: $5（CloudWatch、Secrets Manager等）
- **合計**: 約$27-32/月

### スケールアップ時の見積もり（1,000人/日）

- **Lambda**: $50-100
- **API Gateway**: $35
- **RDS**: $100（db.t3.medium、Multi-AZへアップグレード）
- **Amplify**: $20（増加するトラフィック対応）
- **CloudFront**: $20
- **その他**: $40
- **合計**: 約$265-315/月

### コスト削減施策

- **開発環境**: 夜間・週末の自動停止
- **キャッシュ活用**: API Gateway、CloudFront キャッシュ
- **Lambda**: コールドスタート対策でProvisioned Concurrencyを避ける
- **RDS**: db.t3.microからスタートし、必要に応じてスケールアップ

## 災害復旧計画

### バックアップ戦略

- **RDS**: 自動バックアップ（7 日間）+ 手動スナップショット（月次）
- **Lambdaコード**: Gitリポジトリでバージョン管理
- **設定情報**: Secrets Managerで暗号化保存

### RTO/RPO 目標

- **RTO（目標復旧時間）**: 4 時間
- **RPO（目標復旧時点）**: 1 時間

## セキュリティチェックリスト

- [ ] 全てのデータ転送で TLS 使用
- [ ] 保存データの暗号化
- [ ] IAM ロールの最小権限設定
- [ ] セキュリティグループの最小開放
- [ ] CloudTrail による監査ログ
- [ ] AWS Config によるコンプライアンスチェック Automate
- [ ] 定期的なセキュリティ評価

## 運用手順書

### 日常運用

1. **監視**: CloudWatch ダッシュボード確認
2. **ログ確認**: CloudWatch Logs Insights
3. **パフォーマンス**: X-Ray トレース分析

### 障害対応

1. **アラート受信**: CloudWatch アラーム → SNS → Slack
2. **初動対応**: 影響範囲の特定
3. **復旧作業**: Lambda 再デプロイ、RDS フェイルオーバー等
4. **事後対応**: 障害レポート作成

### スケーリング

- **Lambda**: 自動スケーリング（同時実行数上限に注意）
- **RDS**: 垂直スケーリング（インスタンスタイプ変更）
