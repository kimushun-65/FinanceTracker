# FinSight インフラストラクチャ (CDK)

このディレクトリには、FinSight個人財務管理アプリケーションのAWS CDKインフラストラクチャコードが含まれています。

## アーキテクチャ概要

インフラストラクチャは以下のスタックで構成されています：

1. **VPC Stack** - パブリック/プライベートサブネットを持つネットワークインフラ
2. **Database Stack** - Secrets Manager統合を持つRDS PostgreSQL
3. **API Stack** - API Gateway、Lambda関数、Auth0統合
4. **Amplify Stack** - GitHub統合によるフロントエンドホスティング
5. **SES Stack** - メール送信設定
6. **Monitoring Stack** - CloudWatchダッシュボードとアラーム
7. **Security Stack** - WAFルールとセキュリティ設定

## 前提条件

- Node.js 18.x以上
- 適切な認証情報で設定されたAWS CLI
- CDK CLI (`npm install -g aws-cdk`)
- Go 1.21以上（Lambda関数用）

## プロジェクト構造

```
cdk/
├── bin/            # CDKアプリのエントリポイント
├── lib/            # スタック定義
│   ├── stacks/     # 個別スタック実装
│   └── interfaces/ # TypeScriptインターフェース
├── lambda/         # Lambda関数のソースコード
├── config/         # 環境設定
└── scripts/        # デプロイとテストのスクリプト
```

## 環境設定

環境固有の設定は `config/` ディレクトリに保存されています：
- `prod.json` - 本番環境

## デプロイ

### 初期セットアップ

1. 依存関係のインストール:
   ```bash
   npm install
   ```

2. CDKのブートストラップ（初回のみ）:
   ```bash
   cdk bootstrap
   ```

### 本番環境へのデプロイ

```bash
./scripts/deploy-prod.sh
```

## テスト

### 統合テスト
```bash
./scripts/test-integration.sh prod
```

### パフォーマンステスト
```bash
./scripts/test-performance.sh prod
```

### セキュリティチェック
```bash
./scripts/security-check.sh prod
```

## Lambda関数

アプリケーションには以下のLambda関数が含まれています：

- **auth** - 認証と認可
- **users** - ユーザープロファイル管理
- **accounts** - 金融アカウント管理
- **transactions** - 取引記録
- **categories** - 取引カテゴリ分け
- **budgets** - 予算管理
- **reports** - 財務レポート
- **notifications** - メール通知

## コスト最適化

以下によりコスト最適化されています：
- マルチAZ構成でのNATゲートウェイ（高可用性）
- 適切なRDSインスタンスサイズ（T3.small）
- Lambda予約同時実行数の制限
- CloudWatchログ保持ポリシー

## セキュリティ機能

- レート制限付きWAF保護
- LambdaとRDSのVPC分離
- データベース認証情報のSecrets Manager管理
- 最小権限のIAMロール
- 全Lambda関数のX-Rayトレーシング

## モニタリング

CloudWatchダッシュボードへのアクセス:
- https://ap-northeast-1.console.aws.amazon.com/cloudwatch/home?region=ap-northeast-1#dashboards:name=FinSight-prod

## 便利なコマンド

* `npm run build`   TypeScriptをJavaScriptにコンパイル
* `npm run watch`   変更を監視してコンパイル
* `npm run test`    Jestユニットテストの実行
* `npx cdk deploy`  デフォルトのAWSアカウント/リージョンにデプロイ
* `npx cdk diff`    デプロイ済みスタックと現在の状態を比較
* `npx cdk synth`   合成されたCloudFormationテンプレートを出力
