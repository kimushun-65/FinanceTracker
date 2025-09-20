# インフラセットアップ後の作業リスト

## 概要
CDKによるインフラストラクチャのデプロイは完了しました。本番環境への全スタックのデプロイが完了し、以下のリソースが利用可能です：

### デプロイ済みリソース
- ✅ VPC (ネットワーク基盤)
- ✅ RDS PostgreSQL (データベース) 
- ✅ API Gateway + Lambda (APIエンドポイント: https://65x5ziikn3.execute-api.ap-northeast-1.amazonaws.com/prod/)
- ✅ Amplify (フロントエンド: https://main.d3ppd99k9cae8.amplifyapp.com)
- ✅ SES (メール送信)
- ✅ CloudWatch (監視: https://ap-northeast-1.console.aws.amazon.com/cloudwatch/home?region=ap-northeast-1#dashboards:name=FinSight-prod)
- ✅ WAF (セキュリティ)

実際にアプリケーションを動作させるためには以下の追加作業が必要です。

## 必須作業

### 1. 環境設定ファイルの更新

#### 対象ファイル
- `cdk/config/prod.json`

#### 更新が必要な項目
```json
{
  "auth0Domain": "your-actual-domain.auth0.com",
  "auth0Audience": "https://api.finsight.yourdomain",
  "auth0ClientId": "your-actual-client-id",
  "githubOwner": "kimushun-65", // ✅ 設定済み
  "customDomain": "", // オプション：独自ドメインを使用する場合
  "sesConfig": {
    "fromEmail": "noreply@finsight.com" // ✅ 設定済み（要検証）
  }
}
```

### 2. ✅ AWS Secrets Manager設定（完了）

#### GitHub Personal Access Token
```bash
# ✅ 設定済み (2025-09-20)
aws secretsmanager create-secret \
  --name github-token \
  --secret-string "ghp_xxxxxxxxxxxxxxxxxxxx"
```

#### Auth0 シークレット（必要に応じて）
```bash
aws secretsmanager create-secret \
  --name auth0-client-secret \
  --secret-string "your-auth0-client-secret"
```

### 3. Amazon SES設定

#### メールアドレスの検証
```bash
# ⏳ 検証待ち（確認メール送信済み）
aws ses verify-email-identity --email-address noreply@finsight.com

# 本番環境（ドメイン全体を検証する場合）
aws ses verify-domain-identity --domain yourdomain.com
```

#### サンドボックスからの移行（本番環境）
- AWSサポートに申請が必要
- 申請フォームから送信制限の解除をリクエスト

### 4. Lambda関数の実装

#### 未実装の関数（現在プレースホルダー）
以下のLambda関数はビジネスロジックの実装が必要：

- **transactions** - 取引記録の CRUD 操作
- **budgets** - 予算管理機能
- **categories** - カテゴリ管理
- **reports** - レポート生成ロジック
- **notifications** - 通知送信ロジック

#### 実装例（transactions）
```go
// lambda/transactions/main.go
func handleGetTransactions(userID string) (events.APIGatewayProxyResponse, error) {
    // データベース接続
    // トランザクション取得
    // レスポンス返却
}
```

### 5. ✅ データベーススキーマの確認と調整（完了）

#### 現在のスキーマ
```bash
# ✅ ER図に準拠したスキーマに更新済み
# lambda/db-init/index.js
```

#### 作成済みテーブル
- ✅ users
- ✅ accounts  
- ✅ category_master
- ✅ categories
- ✅ transactions
- ✅ budgets
- ✅ asset_snapshots
- ✅ account_movements
- ✅ budget_suggestions
- ✅ asset_forecasts
- ✅ notification_settings

## オプション作業

### 6. フロントエンド開発

#### React アプリケーションの実装
```bash
cd frontend
npm install
npm start
```

#### 必要な実装
- Auth0 統合
- API クライアント
- 各画面の実装（ダッシュボード、取引一覧、レポート等）

### 7. 本番環境の追加設定

#### Route53 設定（カスタムドメイン使用時）
```bash
# ホストゾーンの作成
aws route53 create-hosted-zone --name yourdomain.com

# ACM証明書の作成（us-east-1リージョンで作成必要）
aws acm request-certificate \
  --domain-name "*.yourdomain.com" \
  --validation-method DNS
```

#### セキュリティ強化
```bash
# CloudTrail有効化
aws cloudtrail create-trail \
  --name finsight-trail \
  --s3-bucket-name finsight-logs

# GuardDuty有効化
aws guardduty create-detector --enable
```

### 8. GitHub Actions設定

#### リポジトリSecrets設定
GitHubリポジトリの Settings > Secrets で以下を設定：

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `SLACK_WEBHOOK`（オプション）

### 9. モニタリング強化

#### カスタムメトリクス追加
- ビジネスメトリクス（日次アクティブユーザー等）
- エラー率の詳細追跡

#### ログ分析
```bash
# CloudWatch Insights クエリ例
fields @timestamp, @message
| filter @message like /ERROR/
| stats count() by bin(5m)
```

## 実行手順

### 本番環境セットアップ

1. **本番用設定**
   ```bash
   vi cdk/config/prod.json
   ```

2. **Secrets設定**
   ```bash
   aws secretsmanager create-secret --name github-token --secret-string "your-token"
   ```

3. **✅ デプロイ実行（完了）**
   ```bash
   cd cdk
   npx cdk deploy --all --context env=prod
   ```

4. **動作確認**
   ```bash
   ./scripts/test-deployment.sh prod
   ./scripts/test-integration.sh prod
   ```

5. **SESメール検証**
   ```bash
   aws ses verify-email-identity --email-address noreply@yourdomain.com
   ```

6. **GitHub Actions有効化**
   - リポジトリSecrets設定
   - mainブランチへのプッシュで自動デプロイ

7. **セキュリティ設定**
   - CloudTrail有効化
   - GuardDuty有効化
   - セキュリティグループの見直し

## チェックリスト

### 必須項目
- [ ] Auth0 設定完了
- [x] 環境設定ファイル更新（一部完了）
- [x] GitHub Token設定
- [ ] SESメール検証（確認メール送信済み）
- [ ] Lambda関数の基本実装

### 推奨項目  
- [ ] カスタムドメイン設定（現在はAmplifyデフォルトドメイン使用）
- [ ] SSL証明書設定（カスタムドメイン使用時）
- [ ] CloudTrail有効化
- [x] バックアップ設定確認（RDS: 7日間保持設定済み）
- [ ] アラート通知先設定（SNSトピック作成済み、通知先未設定）

### 開発項目
- [ ] フロントエンド実装
- [ ] Lambda関数の完全実装
- [ ] 統合テスト作成
- [ ] ドキュメント整備

## トラブルシューティング

問題が発生した場合は、`docs/TROUBLESHOOTING.md` を参照してください。

## 次のステップ

1. Lambda関数のビジネスロジック実装
2. フロントエンドアプリケーション開発
3. エンドツーエンドテストの実装
4. 本番環境への移行準備