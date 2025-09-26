# Finance Tracker デプロイ戦略

## 概要
Lambda関数を削除し、最小構成のECS Fargateでコスト最適化したアーキテクチャへの移行戦略

## アーキテクチャ変更点

### 変更前（Lambda + NAT Gateway）
- API Gateway → Lambda関数 → RDS
- NAT Gateway（月額9,000円）
- Lambda関数（10個以上）

### 変更後（ECS Fargate）
- API Gateway → ALB → ECS Fargate → RDS
- NAT Gateway削除
- 最小構成のECS（0.25 vCPU, 512MB）

## 段階的デプロイ戦略

### Phase 1: 基盤インフラ構築
1. **VPCスタック**
   - CIDR: 10.1.0.0/16（既存と競合回避）
   - パブリックサブネット x2
   - プライベートサブネット x2（データベース用）
   - NAT Gateway: 0（削除）

2. **Databaseスタック**
   - RDS PostgreSQL（t3.small）
   - マルチAZ構成
   - 自動バックアップ7日間

```bash
# Phase 1 デプロイ
cdk deploy VpcStack-prod DatabaseStack-prod --context env=prod
```

### Phase 2: コンテナインフラ構築
1. **ECRリポジトリ作成**
   - リポジトリ名: finance-tracker-prod
   - イメージスキャン有効化

2. **Dockerイメージビルド&プッシュ**
```bash
# ECRログイン
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin [ACCOUNT_ID].dkr.ecr.ap-northeast-1.amazonaws.com

# イメージビルド
cd backend
docker build -t finance-tracker-backend .

# タグ付け&プッシュ
docker tag finance-tracker-backend:latest [ACCOUNT_ID].dkr.ecr.ap-northeast-1.amazonaws.com/finance-tracker-prod:latest
docker push [ACCOUNT_ID].dkr.ecr.ap-northeast-1.amazonaws.com/finance-tracker-prod:latest
```

### Phase 3: ECSサービスデプロイ
1. **ECSスタック**
   - Fargate: 0.25 vCPU, 512MB メモリ
   - Application Load Balancer
   - ヘルスチェック: /health

```bash
cdk deploy EcsStack-prod --context env=prod
```

### Phase 4: API統合&周辺サービス
1. **APIスタック**
   - API Gateway → VPC Link → ALB統合
   - プロキシ統合でシンプル化

2. **その他スタック**
   - SESスタック（メール送信）
   - Monitoringスタック（CloudWatch）
   - Securityスタック（WAF）

```bash
cdk deploy ApiStack-prod SesStack-prod MonitoringStack-prod SecurityStack-prod --context env=prod
```

## コスト見積もり

### 月額費用（東京リージョン）

| サービス | 構成 | 月額費用（円） |
|---------|------|---------------|
| **ECS Fargate** | 0.25 vCPU, 512MB x 1 | 約1,200円 |
| **ALB** | 基本料金 + 処理量 | 約2,500円 |
| **RDS** | t3.micro, シングルAZ | 約2,000円 |
| **API Gateway** | 100万リクエスト/月想定 | 約500円 |
| **VPC** | NAT Gateway削除 | 0円 |
| **CloudWatch** | ログ・メトリクス | 約500円 |
| **その他** | S3, Route53等 | 約300円 |
| **合計** | | **約7,000円** |

### コスト削減効果
- NAT Gateway削除: **-9,000円/月**
- RDS t3.small→t3.micro: **-6,000円/月**
- Lambda関数削除: **-約1,000円/月**
- **削減額合計: 約16,000円/月（-70%）**

## トラブルシューティング

### ECSタスク起動失敗時の対処
1. CloudWatch Logsでコンテナログ確認
2. タスク定義の環境変数確認
3. セキュリティグループ設定確認
4. ヘルスチェックパス確認

### スタック削除時の依存関係エラー
1. 依存スタックから順に削除
   - SecurityStack → MonitoringStack → SesStack → ApiStack → EcsStack → DatabaseStack → VpcStack
2. RDSインスタンスの削除保護確認
3. ネットワークインターフェースの手動削除

## セキュリティ考慮事項

1. **ネットワーク分離**
   - RDSはプライベートサブネット配置
   - セキュリティグループで最小権限

2. **認証・認可**
   - Auth0統合維持
   - API Gatewayでの認証

3. **データ保護**
   - RDS暗号化有効
   - バックアップ自動化

## 今後の最適化案

1. **スポットインスタンス活用**
   - Fargate Spotで最大70%削減可能

2. **オートスケーリング**
   - 負荷に応じた自動スケール設定

3. **コンテナイメージ最適化**
   - マルチステージビルドで軽量化
   - distrolessイメージ使用

## 実装スケジュール

- Phase 1: 30分（VPC + RDS）
- Phase 2: 30分（ECR + イメージビルド）
- Phase 3: 30分（ECS起動確認含む）
- Phase 4: 30分（API統合）
- **合計: 約2時間**

## ロールバック手順

問題発生時は以下の手順でロールバック：

1. 新規作成スタックの削除
```bash
aws cloudformation delete-stack --stack-name [STACK_NAME]
```

2. 古い環境の復旧（必要に応じて）
3. 問題分析と修正
4. 再デプロイ