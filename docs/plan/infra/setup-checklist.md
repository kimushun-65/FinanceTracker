# インフラセットアップ チェックリスト & 検証手順

## 概要

FinSightインフラストラクチャのセットアップから運用開始までの包括的なチェックリストです。各ステップでの検証方法とトラブルシューティングガイドを含みます。

## 🔧 事前準備チェックリスト

### 開発環境準備
- [ ] **Node.js 18+** インストール済み
  ```bash
  node --version  # v18.0.0以上
  npm --version   # v8.0.0以上
  ```

- [ ] **AWS CLI** 設定済み
  ```bash
  aws --version          # v2.0.0以上
  aws sts get-caller-identity  # 認証確認
  ```

- [ ] **Git** 設定済み
  ```bash
  git --version
  git config --global user.name "Your Name"
  git config --global user.email "your.email@example.com"
  ```

- [ ] **Go 1.21+** インストール済み（Lambda開発用）
  ```bash
  go version  # go1.21以上
  ```

### AWS アカウント準備
- [ ] **IAM ユーザー権限** 確認
  - [ ] AdministratorAccess（開発環境）
  - [ ] CloudFormationFullAccess
  - [ ] Route53FullAccess（カスタムドメイン使用時）

- [ ] **サービス制限** 確認
  - [ ] VPC制限（デフォルト5個）
  - [ ] RDS制限（デフォルト40個）
  - [ ] Lambda同時実行制限（デフォルト1000）

- [ ] **請求アラート** 設定済み
  ```bash
  aws budgets describe-budgets --account-id $(aws sts get-caller-identity --query Account --output text)
  ```

### ドメイン準備（本番環境）
- [ ] **ドメイン取得** 済み
- [ ] **Route53 ホストゾーン** 作成済み
- [ ] **DNS委任** 設定済み

---

## 📋 実装フェーズ別チェックリスト

### Phase 1: 基盤構築（Day 1-2）

#### Day 1: プロジェクト初期化
- [ ] **CDKプロジェクト作成**
  ```bash
  cd infrastructure
  cdk init app --language typescript
  npm list  # 依存関係確認
  ```

- [ ] **CDKブートストラップ**
  ```bash
  cdk bootstrap aws://$(aws sts get-caller-identity --query Account --output text)/ap-northeast-1
  ```

- [ ] **プロジェクト構造確認**
  ```
  infrastructure/
  ├── bin/finsight.ts ✓
  ├── lib/
  │   ├── stacks/ ✓
  │   ├── constructs/ ✓
  │   └── interfaces/ ✓
  ├── config/ ✓
  └── scripts/ ✓
  ```

**検証コマンド**:
```bash
cdk list                    # スタック一覧表示
cdk synth                   # CloudFormation テンプレート生成
npm test                    # テスト実行（設定済みの場合）
```

#### Day 2: VPCスタック
- [ ] **VPCスタック実装**
  ```typescript
  // lib/stacks/vpc-stack.ts 実装確認
  export class VpcStack extends Stack {
    public readonly vpc: Vpc;
    public readonly lambdaSecurityGroup: SecurityGroup;
    public readonly rdsSecurityGroup: SecurityGroup;
  }
  ```

- [ ] **設定ファイル準備**
  ```json
  // config/dev.json 必須項目確認
  {
    "environment": "dev",
    "region": "ap-northeast-1",
    "customDomain": null,
    "auth0Domain": "finsight-dev.auth0.com"
  }
  ```

**検証手順**:
```bash
# 1. 構文チェック
cdk synth VpcStack-dev

# 2. デプロイ実行
cdk deploy VpcStack-dev --require-approval never

# 3. リソース確認
aws ec2 describe-vpcs --filters "Name=tag:Name,Values=VpcStack-dev/FinSightVpc"
aws ec2 describe-security-groups --filters "Name=group-name,Values=*FinSight*"
```

**成功判定**:
- [ ] VPC作成成功（CIDR: 10.0.0.0/16）
- [ ] パブリックサブネット2個作成
- [ ] プライベートサブネット2個作成
- [ ] セキュリティグループ2個作成
- [ ] NATゲートウェイ作成（dev: 1個, prod: 2個）

### Phase 2: データベース基盤（Day 3-4）

#### Day 3: RDS + Secrets Manager
- [ ] **Secrets Manager設定**
  ```bash
  # シークレット確認
  aws secretsmanager list-secrets --query 'SecretList[?contains(Name, `finsight-db-credentials`)].Name'
  ```

- [ ] **RDS PostgreSQL設定**
  ```bash
  # RDS インスタンス確認
  aws rds describe-db-instances --query 'DBInstances[?DBName==`finsight`]'
  ```

**検証手順**:
```bash
# 1. DatabaseStack デプロイ
cdk deploy DatabaseStack-dev --require-approval never

# 2. 接続テスト（Lambda経由）
aws lambda invoke \
  --function-name $(aws lambda list-functions --query 'Functions[?contains(FunctionName, `DbInit`)].FunctionName' --output text) \
  /tmp/db-test.json

# 3. RDS接続確認
aws rds describe-db-instances \
  --db-instance-identifier $(aws rds describe-db-instances --query 'DBInstances[0].DBInstanceIdentifier' --output text)
```

**成功判定**:
- [ ] RDS PostgreSQL 15 インスタンス起動
- [ ] 暗号化設定有効
- [ ] バックアップ設定有効（7日間）
- [ ] VPCプライベートサブネット配置
- [ ] Secrets Manager認証情報作成

#### Day 4: データベース初期化
- [ ] **初期化Lambda実装**
  ```javascript
  // lambda/db-init/index.js
  // PostgreSQL接続 + 基本テーブル作成
  ```

**検証手順**:
```bash
# Lambda実行テスト
aws lambda invoke \
  --function-name finsight-db-init-dev \
  --log-type Tail \
  /tmp/init-response.json

# ログ確認
aws logs tail /aws/lambda/finsight-db-init-dev --follow
```

**成功判定**:
- [ ] Lambda関数正常実行
- [ ] データベース接続成功
- [ ] 基本テーブル作成成功
- [ ] エラーログなし

### Phase 3: API基盤（Day 5-7）

#### Day 5: API Gateway基盤
- [ ] **Lambda Authorizer実装**
  ```go
  // lambda/authorizer/main.go
  // JWT検証ロジック実装確認
  ```

- [ ] **API Gateway設定**
  ```typescript
  // lib/stacks/api-stack.ts
  // REST API + CORS設定確認
  ```

**検証手順**:
```bash
# 1. API Stack デプロイ
cdk deploy ApiStack-dev --require-approval never

# 2. API Gateway確認
API_ID=$(aws apigateway get-rest-apis --query 'items[0].id' --output text)
echo "API Gateway ID: $API_ID"

# 3. CORS設定確認
aws apigateway get-method \
  --rest-api-id $API_ID \
  --resource-id $(aws apigateway get-resources --rest-api-id $API_ID --query 'items[0].id' --output text) \
  --http-method OPTIONS
```

**成功判定**:
- [ ] API Gateway作成成功
- [ ] CORS設定有効
- [ ] Lambda Authorizer動作
- [ ] CloudWatch Logs設定

#### Day 6-7: Lambda関数 + エンドポイント
- [ ] **Lambda関数実装**
  ```bash
  # 各Lambda関数確認
  ls lambda/{auth,users,transactions,budgets,accounts,categories,reports,notifications}/
  ```

**検証手順**:
```bash
# 1. 全Lambda関数確認
aws lambda list-functions --query 'Functions[?contains(FunctionName, `finsight`)].FunctionName'

# 2. ヘルスチェック
API_URL="https://${API_ID}.execute-api.ap-northeast-1.amazonaws.com/prod"
curl -X GET "$API_URL/health"

# 3. 認証エンドポイント確認
curl -X POST "$API_URL/auth/callback" \
  -H "Content-Type: application/json" \
  -d '{"test": true}'
```

**成功判定**:
- [ ] 8つのLambda関数デプロイ成功
- [ ] 全エンドポイント作成成功
- [ ] ヘルスチェック200応答
- [ ] エラーハンドリング動作

### Phase 4: フロントエンド基盤（Day 8-9）

#### Day 8: SSL証明書 + Amplify
- [ ] **Certificate Stack**（カスタムドメイン使用時）
  ```bash
  # 証明書確認
  aws acm list-certificates --region us-east-1
  ```

- [ ] **Amplify設定**
  ```bash
  # Amplifyアプリ確認
  aws amplify list-apps
  ```

**検証手順**:
```bash
# 1. 証明書デプロイ（カスタムドメインの場合）
cdk deploy CertificateStack-dev --require-approval never

# 2. Amplify デプロイ
cdk deploy AmplifyStack-dev --require-approval never

# 3. GitHub連携確認
aws amplify get-app --app-id $(aws amplify list-apps --query 'apps[0].appId' --output text)
```

**成功判定**:
- [ ] SSL証明書発行成功（該当する場合）
- [ ] Amplifyアプリ作成成功
- [ ] GitHub連携設定成功
- [ ] 環境変数設定成功

### Phase 5: メール送信（Day 9-10）

#### Day 9: SES設定
- [ ] **SES Identity確認**
  ```bash
  # SES設定確認
  aws ses list-identities
  aws ses get-configuration-set --configuration-set-name finsight-dev
  ```

**検証手順**:
```bash
# 1. SES Stack デプロイ
cdk deploy SesStack-dev --require-approval never

# 2. 送信テスト
aws ses send-email \
  --source "noreply@dev.finsight.local" \
  --destination "ToAddresses=your-email@example.com" \
  --message "Subject={Data=Test},Body={Text={Data=Test Email}}"

# 3. 設定セット確認
aws ses describe-configuration-set --configuration-set-name finsight-dev
```

**成功判定**:
- [ ] SES Identity確認済み
- [ ] 設定セット作成成功
- [ ] 送信テスト成功
- [ ] バウンス処理設定成功

### Phase 6: 監視設定（Day 10-11）

#### Day 10: CloudWatch + X-Ray
- [ ] **監視設定確認**
  ```bash
  # ダッシュボード確認
  aws cloudwatch list-dashboards
  
  # アラーム確認
  aws cloudwatch describe-alarms --query 'MetricAlarms[?contains(AlarmName, `finsight`)].AlarmName'
  ```

**検証手順**:
```bash
# 1. Monitoring Stack デプロイ
cdk deploy MonitoringStack-dev --require-approval never

# 2. X-Ray確認
aws xray get-service-graph

# 3. ログ確認
aws logs describe-log-groups --log-group-name-prefix "/aws/lambda/finsight"
```

**成功判定**:
- [ ] CloudWatchダッシュボード作成
- [ ] 主要アラーム設定成功
- [ ] X-Rayトレーシング有効
- [ ] ログ保持期間設定

### Phase 7: CI/CD + セキュリティ（Day 11-12）

#### Day 11: GitHub Actions
- [ ] **ワークフロー設定**
  ```yaml
  # .github/workflows/deploy-infrastructure.yml
  # 必須ジョブ確認: validate, deploy-dev, deploy-prod
  ```

**検証手順**:
```bash
# 1. GitHub Secrets設定確認
# AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY等

# 2. ワークフロー実行テスト
git push origin develop  # 開発環境デプロイトリガー

# 3. デプロイ結果確認
# GitHub Actions UI で確認
```

**成功判定**:
- [ ] GitHub Actions ワークフロー作成
- [ ] 環境別デプロイ設定成功
- [ ] シークレット管理設定成功
- [ ] 自動デプロイ成功

#### Day 12: セキュリティ設定
- [ ] **WAF設定**
  ```bash
  # WAF確認
  aws wafv2 list-web-acls --scope REGIONAL
  ```

**検証手順**:
```bash
# 1. Security Stack デプロイ
cdk deploy SecurityStack-dev --require-approval never

# 2. WAF設定確認
aws wafv2 get-web-acl \
  --scope REGIONAL \
  --id $(aws wafv2 list-web-acls --scope REGIONAL --query 'WebACLs[0].Id' --output text)

# 3. レート制限テスト
for i in {1..100}; do curl -s "$API_URL/health" > /dev/null; done
```

**成功判定**:
- [ ] WAF WebACL作成成功
- [ ] レート制限ルール動作
- [ ] SQLインジェクション対策有効
- [ ] XSS対策有効

---

## 🧪 総合テストチェックリスト

### 機能テスト
- [ ] **API エンドポイント全体テスト**
  ```bash
  # 全エンドポイントテスト実行
  ./scripts/test-deployment.sh dev
  ```

- [ ] **認証フローテスト**
  ```bash
  # Auth0統合テスト（手動）
  # 1. Auth0でテストユーザー作成
  # 2. JWTトークン取得
  # 3. API認証テスト
  ```

- [ ] **データベース操作テスト**
  ```bash
  # CRUD操作確認
  # Lambda経由でデータベース操作テスト
  ```

### パフォーマンステスト
- [ ] **レスポンス時間確認**
  ```bash
  # API レスポンス時間測定
  curl -w "Time: %{time_total}s\n" -o /dev/null -s "$API_URL/health"
  ```

- [ ] **同時接続テスト**
  ```bash
  # 負荷テスト（簡易版）
  ab -n 100 -c 10 "$API_URL/health"
  ```

### セキュリティテスト
- [ ] **脆弱性スキャン**
  ```bash
  # セキュリティ設定確認
  ./scripts/security-check.sh
  ```

- [ ] **権限確認**
  ```bash
  # IAM権限最小化確認
  aws iam list-roles --query 'Roles[?contains(RoleName, `finsight`)]'
  ```

---

## 🚨 トラブルシューティング

### よくある問題と解決法

#### 1. CDKデプロイエラー
**問題**: `cdk deploy`でエラーが発生
```
Error: This stack uses assets, so the toolkit stack must be deployed to the environment
```

**解決法**:
```bash
# CDK再ブートストラップ
cdk bootstrap aws://$(aws sts get-caller-identity --query Account --output text)/ap-northeast-1
```

#### 2. VPC制限エラー
**問題**: VPC作成で制限エラー
```
VpcLimitExceeded: The maximum number of VPCs has been reached
```

**解決法**:
```bash
# 既存VPC確認・削除
aws ec2 describe-vpcs
aws ec2 delete-vpc --vpc-id vpc-xxxxxxxx
```

#### 3. RDS接続エラー
**問題**: LambdaからRDS接続失敗

**確認項目**:
```bash
# 1. セキュリティグループ確認
aws ec2 describe-security-groups --group-ids sg-xxxxxxxx

# 2. VPC設定確認
aws lambda get-function-configuration --function-name function-name

# 3. Secrets Manager権限確認
aws lambda get-policy --function-name function-name
```

#### 4. API Gateway 502エラー
**問題**: API Gateway で502 Bad Gateway

**確認手順**:
```bash
# 1. Lambda関数ログ確認
aws logs tail /aws/lambda/function-name --follow

# 2. Lambda実行権限確認
aws lambda get-function --function-name function-name

# 3. VPC設定確認（タイムアウトの可能性）
aws lambda get-function-configuration --function-name function-name
```

#### 5. SES送信エラー
**問題**: メール送信が失敗

**確認項目**:
```bash
# 1. SES Identity確認
aws ses get-identity-verification-attributes

# 2. 送信制限確認
aws ses get-send-quota

# 3. 設定セット確認
aws ses describe-configuration-set --configuration-set-name finsight-dev
```

### エスカレーション手順

#### レベル1: 自己解決
1. **エラーメッセージ確認**
2. **AWS CloudWatch Logs確認**
3. **このドキュメントの該当セクション確認**

#### レベル2: コミュニティ相談
1. **Stack Overflow検索**
2. **AWS re:Post検索**
3. **GitHub Issues検索**

#### レベル3: AWSサポート
1. **AWS サポートケース作成**
2. **問題の詳細情報収集**
3. **再現手順の文書化**

---

## 📊 完了確認チェックリスト

### インフラストラクチャ完成度
- [ ] **100%** VPC基盤構築
- [ ] **100%** データベース設定
- [ ] **100%** API Gateway + Lambda
- [ ] **100%** フロントエンドホスティング
- [ ] **100%** メール送信機能
- [ ] **100%** 監視・ログ設定
- [ ] **100%** CI/CD パイプライン
- [ ] **100%** セキュリティ設定

### 運用準備完了度
- [ ] **100%** 自動デプロイ機能
- [ ] **100%** 監視・アラート設定
- [ ] **100%** バックアップ設定
- [ ] **100%** セキュリティ設定
- [ ] **100%** 文書化完了

### 次フェーズ準備完了度
- [ ] **100%** バックエンド開発環境準備
- [ ] **100%** フロントエンド開発環境準備
- [ ] **100%** 開発チーム向け情報整理

## 🎯 成果物一覧

### コードベース
```
infrastructure/
├── bin/finsight.ts                    # CDK アプリエントリーポイント
├── lib/stacks/                        # CDK スタック実装
│   ├── vpc-stack.ts                   # VPC基盤
│   ├── database-stack.ts              # RDS PostgreSQL
│   ├── api-stack.ts                   # API Gateway + Lambda
│   ├── certificate-stack.ts           # SSL証明書
│   ├── amplify-stack.ts               # フロントエンドホスティング
│   ├── ses-stack.ts                   # メール送信
│   ├── monitoring-stack.ts            # 監視・ログ
│   └── security-stack.ts              # WAF + セキュリティ
├── lib/constructs/                    # 再利用可能コンポーネント
├── lib/interfaces/config.ts           # 設定インターフェース
├── config/                            # 環境別設定
│   ├── dev.json
│   ├── staging.json
│   └── prod.json
├── scripts/                           # 運用スクリプト
│   ├── deploy-dev.sh
│   ├── deploy-staging.sh
│   ├── deploy-prod.sh
│   ├── test-deployment.sh
│   └── security-check.sh
├── lambda/                            # Lambda関数ソースコード
│   ├── auth/
│   ├── users/
│   ├── transactions/
│   ├── budgets/
│   ├── accounts/
│   ├── categories/
│   ├── reports/
│   ├── notifications/
│   ├── authorizer/
│   ├── db-init/
│   └── shared/
└── .github/workflows/                 # CI/CD設定
    └── deploy-infrastructure.yml
```

### AWSリソース
- **14個** のCloudFormationスタック
- **50+** のAWSリソース
- **3つ** の環境（dev/staging/prod）
- **完全自動化** されたデプロイメント

このチェックリストにより、インフラストラクチャの構築から運用開始までを確実に完了できます。