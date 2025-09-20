# AWS CDK デプロイプラン

## 概要

このドキュメントでは、FinSightアプリケーションのAWSインフラストラクチャをAWS CDKを使用してデプロイするための詳細なプランを示します。サーバーレスアーキテクチャを採用し、スケーラビリティ、セキュリティ、コスト効率を最大化します。

## CDKプロジェクト構造

```
infrastructure/
├── bin/
│   └── finsight.ts              # CDKアプリケーションエントリーポイント
├── lib/
│   ├── stacks/
│   │   ├── vpc-stack.ts         # VPCリソース
│   │   ├── database-stack.ts    # RDS PostgreSQL
│   │   ├── api-stack.ts         # API Gateway + Lambda
│   │   ├── certificate-stack.ts # SSL証明書
│   │   ├── amplify-stack.ts     # フロントエンドホスティング
│   │   ├── ses-stack.ts         # SESメール送信
│   │   └── monitoring-stack.ts  # CloudWatch + X-Ray
│   ├── constructs/
│   │   ├── lambda-api.ts        # Lambda関数構成
│   │   ├── rds-postgres.ts      # RDS構成
│   │   └── api-gateway.ts       # API Gateway構成
│   └── interfaces/
│       └── config.ts            # 環境設定インターフェース
├── config/
│   ├── dev.json                 # 開発環境設定
│   ├── staging.json             # ステージング環境設定
│   └── prod.json                # 本番環境設定
├── scripts/
│   ├── deploy-dev.sh            # 開発環境デプロイスクリプト
│   ├── deploy-staging.sh        # ステージング環境デプロイスクリプト
│   └── deploy-prod.sh           # 本番環境デプロイスクリプト
├── cdk.json                     # CDK設定
├── package.json                 # 依存関係
├── tsconfig.json                # TypeScript設定
└── README.md                    # デプロイ手順
```

## スタック設計

### 1. VPCスタック (VpcStack)

**責務**: ネットワーク基盤の構築

```typescript
// lib/stacks/vpc-stack.ts
export class VpcStack extends Stack {
  public readonly vpc: Vpc;
  public readonly lambdaSecurityGroup: SecurityGroup;
  public readonly rdsSecurityGroup: SecurityGroup;

  constructor(scope: Construct, id: string, props: StackProps) {
    super(scope, id, props);

    // VPC作成
    this.vpc = new Vpc(this, 'FinSightVpc', {
      cidr: '10.0.0.0/16',
      maxAzs: 2,
      subnetConfiguration: [
        {
          cidrMask: 24,
          name: 'Public',
          subnetType: SubnetType.PUBLIC,
        },
        {
          cidrMask: 24,
          name: 'Private',
          subnetType: SubnetType.PRIVATE_WITH_EGRESS,
        },
      ],
      natGateways: 1, // コスト削減のため1つのみ
    });

    // Lambda用セキュリティグループ
    this.lambdaSecurityGroup = new SecurityGroup(this, 'LambdaSecurityGroup', {
      vpc: this.vpc,
      description: 'Security group for Lambda functions',
      allowAllOutbound: true,
    });

    // RDS用セキュリティグループ
    this.rdsSecurityGroup = new SecurityGroup(this, 'RdsSecurityGroup', {
      vpc: this.vpc,
      description: 'Security group for RDS database',
      allowAllOutbound: false,
    });

    // Lambda → RDS接続許可
    this.rdsSecurityGroup.addIngressRule(
      this.lambdaSecurityGroup,
      Port.tcp(5432),
      'Allow Lambda to access RDS'
    );
  }
}
```

### 2. データベーススタック (DatabaseStack)

**責務**: RDS PostgreSQLの構築

```typescript
// lib/stacks/database-stack.ts
export class DatabaseStack extends Stack {
  public readonly database: DatabaseInstance;
  public readonly dbSecret: Secret;

  constructor(scope: Construct, id: string, props: DatabaseStackProps) {
    super(scope, id, props);

    // データベース認証情報をSecrets Managerで管理
    this.dbSecret = new Secret(this, 'DatabaseSecret', {
      secretName: `finsight-db-credentials-${props.environment}`,
      generateSecretString: {
        secretStringTemplate: JSON.stringify({ username: 'postgres' }),
        generateStringKey: 'password',
        excludeCharacters: '"@/\\',
      },
    });

    // RDS インスタンス作成
    this.database = new DatabaseInstance(this, 'FinSightDatabase', {
      engine: DatabaseInstanceEngine.postgres({
        version: PostgresEngineVersion.VER_15,
      }),
      instanceType: props.environment === 'prod' 
        ? InstanceType.of(InstanceClass.T3, InstanceSize.SMALL)
        : InstanceType.of(InstanceClass.T3, InstanceSize.MICRO),
      vpc: props.vpc,
      vpcSubnets: {
        subnetType: SubnetType.PRIVATE_WITH_EGRESS,
      },
      securityGroups: [props.rdsSecurityGroup],
      credentials: Credentials.fromSecret(this.dbSecret),
      multiAz: props.environment === 'prod',
      storageEncrypted: true,
      backupRetention: Duration.days(7),
      deletionProtection: props.environment === 'prod',
      databaseName: 'finsight',
    });

    // データベース初期化Lambda
    const dbInitFunction = new Function(this, 'DbInitFunction', {
      runtime: Runtime.NODEJS_18_X,
      handler: 'index.handler',
      code: Code.fromAsset('lambda/db-init'),
      vpc: props.vpc,
      securityGroups: [props.lambdaSecurityGroup],
      environment: {
        DB_SECRET_ARN: this.dbSecret.secretArn,
        DB_ENDPOINT: this.database.instanceEndpoint.hostname,
      },
      timeout: Duration.minutes(5),
    });

    // Secret読み取り権限付与
    this.dbSecret.grantRead(dbInitFunction);
  }
}
```

### 3. APIスタック (ApiStack)

**責務**: API Gateway、Lambda関数、認証の構築

```typescript
// lib/stacks/api-stack.ts
export class ApiStack extends Stack {
  public readonly api: RestApi;
  public readonly lambdaFunctions: { [key: string]: Function };

  constructor(scope: Construct, id: string, props: ApiStackProps) {
    super(scope, id, props);

    // API Gateway作成
    this.api = new RestApi(this, 'FinSightApi', {
      restApiName: `finsight-api-${props.environment}`,
      description: 'FinSight REST API',
      domainName: props.customDomain ? {
        domainName: `api.${props.customDomain}`,
        certificate: props.certificate,
      } : undefined,
      defaultCorsPreflightOptions: {
        allowOrigins: Cors.ALL_ORIGINS,
        allowMethods: Cors.ALL_METHODS,
        allowHeaders: ['Content-Type', 'Authorization'],
      },
    });

    // Lambda Authorizer（Auth0 JWT検証）
    const authorizer = new TokenAuthorizer(this, 'Auth0Authorizer', {
      handler: this.createAuthorizerFunction(props),
      identitySource: 'method.request.header.Authorization',
      resultsCacheTtl: Duration.minutes(5),
    });

    // Lambda関数群の作成
    this.lambdaFunctions = this.createApiLambdas(props);

    // APIエンドポイントの設定
    this.setupApiEndpoints(authorizer);
  }

  private createApiLambdas(props: ApiStackProps): { [key: string]: Function } {
    const functions: { [key: string]: Function } = {};

    const lambdaConfigs = [
      { name: 'auth', path: 'auth' },
      { name: 'users', path: 'users' },
      { name: 'transactions', path: 'transactions' },
      { name: 'budgets', path: 'budgets' },
      { name: 'accounts', path: 'accounts' },
      { name: 'categories', path: 'categories' },
      { name: 'reports', path: 'reports' },
      { name: 'notifications', path: 'notifications' },
    ];

    lambdaConfigs.forEach(config => {
      functions[config.name] = new Function(this, `${config.name}Function`, {
        runtime: Runtime.PROVIDED_AL2,
        handler: 'bootstrap',
        code: Code.fromAsset(`lambda/${config.path}`),
        vpc: props.vpc,
        securityGroups: [props.lambdaSecurityGroup],
        environment: {
          DB_SECRET_ARN: props.dbSecret.secretArn,
          ENVIRONMENT: props.environment,
          AUTH0_DOMAIN: props.auth0Domain,
          AUTH0_AUDIENCE: props.auth0Audience,
          SES_FROM_EMAIL: props.sesFromEmail,
          SES_CONFIGURATION_SET: `finsight-${props.environment}`,
        },
        timeout: Duration.seconds(30),
        memorySize: 512,
      });

      // Secrets Manager読み取り権限
      props.dbSecret.grantRead(functions[config.name]);

      // SES送信権限（reportsとnotifications関数のみ）
      if (config.name === 'reports' || config.name === 'notifications') {
        functions[config.name].addToRolePolicy(new PolicyStatement({
          actions: [
            'ses:SendEmail',
            'ses:SendRawEmail',
          ],
          resources: [
            `arn:aws:ses:${this.region}:${this.account}:identity/${props.sesFromEmail}`,
            `arn:aws:ses:${this.region}:${this.account}:configuration-set/finsight-${props.environment}`,
          ],
        }));
      }
    });

    return functions;
  }

  private setupApiEndpoints(authorizer: TokenAuthorizer): void {
    // /auth - 認証不要
    const authResource = this.api.root.addResource('auth');
    authResource.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.auth));

    // /users - 認証必要
    const usersResource = this.api.root.addResource('users');
    usersResource.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.users), {
      authorizer,
    });
    usersResource.addMethod('PUT', new LambdaIntegration(this.lambdaFunctions.users), {
      authorizer,
    });

    // その他のエンドポイントも同様に設定...
  }
}
```

### 4. 証明書スタック (CertificateStack)

**責務**: SSL/TLS証明書の管理

```typescript
// lib/stacks/certificate-stack.ts
export class CertificateStack extends Stack {
  public readonly certificate: Certificate;

  constructor(scope: Construct, id: string, props: CertificateStackProps) {
    super(scope, id, props);

    this.certificate = new Certificate(this, 'FinSightCertificate', {
      domainName: props.domainName,
      subjectAlternativeNames: [`*.${props.domainName}`],
      validation: CertificateValidation.fromDns(),
    });
  }
}
```

### 5. Amplifyスタック (AmplifyStack)

**責務**: フロントエンドホスティングの構築

```typescript
// lib/stacks/amplify-stack.ts
export class AmplifyStack extends Stack {
  public readonly app: App;

  constructor(scope: Construct, id: string, props: AmplifyStackProps) {
    super(scope, id, props);

    // Amplifyアプリ作成
    this.app = new App(this, 'FinSightApp', {
      sourceCodeProvider: new GitHubSourceCodeProvider({
        owner: props.githubOwner,
        repository: props.repositoryName,
        oauthToken: SecretValue.secretsManager('github-token'),
      }),
      environmentVariables: {
        REACT_APP_API_URL: props.apiUrl,
        REACT_APP_AUTH0_DOMAIN: props.auth0Domain,
        REACT_APP_AUTH0_CLIENT_ID: props.auth0ClientId,
      },
    });

    // ブランチ設定
    const mainBranch = this.app.addBranch('main', {
      autoBuild: true,
      stage: props.environment === 'prod' ? 'PRODUCTION' : 'DEVELOPMENT',
    });

    // カスタムドメイン設定
    if (props.customDomain) {
      const domain = this.app.addDomain('FinSightDomain', {
        domainName: props.customDomain,
        subDomains: [
          {
            branch: mainBranch,
            prefix: props.environment === 'prod' ? '' : props.environment,
          },
        ],
      });
    }
  }
}
```

### 6. SESスタック (SesStack)

**責務**: Amazon SESメール送信サービスの構築

```typescript
// lib/stacks/ses-stack.ts
export class SesStack extends Stack {
  public readonly verifiedDomain: string;
  public readonly sesIdentity: EmailIdentity;

  constructor(scope: Construct, id: string, props: SesStackProps) {
    super(scope, id, props);

    // ドメイン認証（カスタムドメインがある場合）
    if (props.customDomain) {
      this.sesIdentity = new EmailIdentity(this, 'SesEmailIdentity', {
        identity: Identity.domain(props.customDomain),
        dkimSigning: true,
      });
      this.verifiedDomain = props.customDomain;
    } else {
      // 開発環境用：検証済みメールアドレス
      this.sesIdentity = new EmailIdentity(this, 'SesEmailIdentity', {
        identity: Identity.email(props.fromEmail),
      });
      this.verifiedDomain = props.fromEmail;
    }

    // 送信レート制限設定
    new CfnConfigurationSet(this, 'SesConfigurationSet', {
      name: `finsight-${props.environment}`,
      deliveryOptions: {
        tlsPolicy: 'Require',
      },
      reputationOptions: {
        reputationMetricsEnabled: true,
      },
      sendingOptions: {
        sendingEnabled: true,
      },
      suppressionOptions: {
        suppressedReasons: ['BOUNCE', 'COMPLAINT'],
      },
    });

    // バウンス・苦情処理用SNSトピック
    const bouncesTopic = new Topic(this, 'SesBouncestopic', {
      topicName: `finsight-ses-bounces-${props.environment}`,
    });

    const complaintsTopic = new Topic(this, 'SesComplaintsTopic', {
      topicName: `finsight-ses-complaints-${props.environment}`,
    });

    // SES送信統計をCloudWatchに送信
    new CfnConfigurationSetEventDestination(this, 'SesCloudWatchDestination', {
      configurationSetName: `finsight-${props.environment}`,
      eventDestination: {
        name: 'cloudwatch-destination',
        enabled: true,
        matchingEventTypes: ['send', 'reject', 'bounce', 'complaint', 'delivery'],
        cloudWatchDestination: {
          dimensionConfigurations: [
            {
              dimensionName: 'MessageTag',
              dimensionValueSource: 'messageTag',
              defaultDimensionValue: 'default',
            },
          ],
        },
      },
    });

    // Lambda関数にSES送信権限を付与するためのIAMポリシー
    const sesPolicy = new PolicyDocument({
      statements: [
        new PolicyStatement({
          actions: [
            'ses:SendEmail',
            'ses:SendRawEmail',
            'ses:GetSendQuota',
            'ses:GetSendStatistics',
          ],
          resources: [
            `arn:aws:ses:${this.region}:${this.account}:identity/${this.verifiedDomain}`,
            `arn:aws:ses:${this.region}:${this.account}:configuration-set/finsight-${props.environment}`,
          ],
        }),
      ],
    });

    // 出力値
    new CfnOutput(this, 'SesIdentityArn', {
      value: this.sesIdentity.emailIdentityArn,
      description: 'SES Email Identity ARN',
    });

    new CfnOutput(this, 'SesFromEmail', {
      value: this.verifiedDomain,
      description: 'SES From Email Address',
    });
  }
}
```

### 7. 監視スタック (MonitoringStack)

**責務**: CloudWatch、X-Ray、アラームの構築

```typescript
// lib/stacks/monitoring-stack.ts
export class MonitoringStack extends Stack {
  constructor(scope: Construct, id: string, props: MonitoringStackProps) {
    super(scope, id, props);

    // CloudWatchダッシュボード
    const dashboard = new Dashboard(this, 'FinSightDashboard', {
      dashboardName: `finsight-${props.environment}`,
    });

    // APIメトリクス
    dashboard.addWidgets(
      new GraphWidget({
        title: 'API Gateway Requests',
        left: [props.api.metricRequestCount()],
        right: [props.api.metricLatency()],
      }),
      new GraphWidget({
        title: 'Lambda Invocations',
        left: [new Metric({
          namespace: 'AWS/Lambda',
          metricName: 'Invocations',
          dimensionsMap: {
            FunctionName: 'finsight-*',
          },
        })],
      }),
    );

    // アラーム設定
    this.createAlarms(props);

    // X-Ray有効化
    this.enableXRayTracing(props);
  }

  private createAlarms(props: MonitoringStackProps): void {
    // API Gatewayエラー率アラーム
    new Alarm(this, 'ApiErrorAlarm', {
      metric: props.api.metricClientError(),
      threshold: 10,
      evaluationPeriods: 2,
      alarmDescription: 'API Gateway error rate is too high',
    });

    // データベース接続アラーム
    new Alarm(this, 'DatabaseConnectionAlarm', {
      metric: props.database.metricDatabaseConnections(),
      threshold: 80,
      evaluationPeriods: 2,
      alarmDescription: 'Database connection count is high',
    });
  }

  private enableXRayTracing(props: MonitoringStackProps): void {
    // Lambda関数にX-Rayトレーシング有効化
    Object.values(props.lambdaFunctions).forEach(func => {
      func.addToRolePolicy(new PolicyStatement({
        actions: [
          'xray:PutTraceSegments',
          'xray:PutTelemetryRecords',
        ],
        resources: ['*'],
      }));
    });
  }
}
```

## 環境設定

### 開発環境 (dev.json)

```json
{
  "environment": "dev",
  "region": "ap-northeast-1",
  "customDomain": null,
  "auth0Domain": "finsight-dev.auth0.com",
  "auth0Audience": "https://api-dev.finsight.com",
  "auth0ClientId": "${DEV_AUTH0_CLIENT_ID}",
  "githubOwner": "your-org",
  "repositoryName": "finsight-frontend",
  "databaseConfig": {
    "instanceType": "t3.micro",
    "multiAz": false,
    "deletionProtection": false
  },
  "lambdaConfig": {
    "memorySize": 512,
    "timeout": 30
  },
  "sesConfig": {
    "fromEmail": "noreply@dev.finsight.com",
    "sendingQuota": 200,
    "sendingRate": 1
  }
}
```

### ステージング環境 (staging.json)

```json
{
  "environment": "staging",
  "region": "ap-northeast-1",
  "customDomain": "staging.finsight.com",
  "auth0Domain": "finsight-staging.auth0.com",
  "auth0Audience": "https://api-staging.finsight.com",
  "auth0ClientId": "${STAGING_AUTH0_CLIENT_ID}",
  "githubOwner": "your-org",
  "repositoryName": "finsight-frontend",
  "databaseConfig": {
    "instanceType": "t3.micro",
    "multiAz": false,
    "deletionProtection": true
  },
  "lambdaConfig": {
    "memorySize": 512,
    "timeout": 30
  },
  "sesConfig": {
    "fromEmail": "noreply@staging.finsight.com",
    "sendingQuota": 1000,
    "sendingRate": 5
  }
}
```

### 本番環境 (prod.json)

```json
{
  "environment": "prod",
  "region": "ap-northeast-1",
  "customDomain": "finsight.com",
  "auth0Domain": "finsight.auth0.com",
  "auth0Audience": "https://api.finsight.com",
  "auth0ClientId": "${PROD_AUTH0_CLIENT_ID}",
  "githubOwner": "your-org",
  "repositoryName": "finsight-frontend",
  "databaseConfig": {
    "instanceType": "t3.small",
    "multiAz": true,
    "deletionProtection": true
  },
  "lambdaConfig": {
    "memorySize": 1024,
    "timeout": 30
  },
  "sesConfig": {
    "fromEmail": "noreply@finsight.com",
    "sendingQuota": 10000,
    "sendingRate": 14
  }
}
```

## デプロイメント手順

### 1. 初期セットアップ

```bash
# CDKプロジェクト初期化
cd infrastructure
npm install

# CDK CLIインストール（未インストールの場合）
npm install -g aws-cdk

# AWS認証情報設定
aws configure

# CDKブートストラップ（初回のみ）
cdk bootstrap aws://ACCOUNT-NUMBER/ap-northeast-1
```

### 2. 環境別デプロイ

#### 開発環境

```bash
# 環境変数設定
export DEV_AUTH0_CLIENT_ID="your-dev-client-id"

# デプロイ実行
npm run deploy:dev

# または個別実行
cdk deploy --context env=dev --all
```

#### ステージング環境

```bash
# 環境変数設定
export STAGING_AUTH0_CLIENT_ID="your-staging-client-id"

# デプロイ実行
npm run deploy:staging

# または個別実行
cdk deploy --context env=staging --all
```

#### 本番環境

```bash
# 環境変数設定
export PROD_AUTH0_CLIENT_ID="your-prod-client-id"

# デプロイ実行（承認が必要）
npm run deploy:prod

# または個別実行
cdk deploy --context env=prod --all --require-approval broadening
```

### 3. スタック別デプロイ順序

```bash
# 1. VPCスタック
cdk deploy VpcStack --context env=dev

# 2. 証明書スタック（カスタムドメイン使用時）
cdk deploy CertificateStack --context env=dev

# 3. データベーススタック
cdk deploy DatabaseStack --context env=dev

# 4. APIスタック
cdk deploy ApiStack --context env=dev

# 5. Amplifyスタック
cdk deploy AmplifyStack --context env=dev

# 6. SESスタック
cdk deploy SesStack --context env=dev

# 7. 監視スタック
cdk deploy MonitoringStack --context env=dev
```

## CI/CDパイプライン

### GitHub Actions設定

```yaml
# .github/workflows/deploy-infrastructure.yml
name: Deploy Infrastructure

on:
  push:
    branches: [main, develop]
    paths: ['infrastructure/**']
  pull_request:
    branches: [main]
    paths: ['infrastructure/**']

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: infrastructure/package-lock.json
      
      - name: Install dependencies
        run: |
          cd infrastructure
          npm ci
      
      - name: Run tests
        run: |
          cd infrastructure
          npm test
      
      - name: CDK Synth
        run: |
          cd infrastructure
          npx cdk synth --context env=dev

  deploy-dev:
    needs: validate
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'
    environment: development
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: infrastructure/package-lock.json
      
      - name: Install dependencies
        run: |
          cd infrastructure
          npm ci
      
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ap-northeast-1
      
      - name: Deploy to Dev
        run: |
          cd infrastructure
          npx cdk deploy --context env=dev --all --require-approval never
        env:
          DEV_AUTH0_CLIENT_ID: ${{ secrets.DEV_AUTH0_CLIENT_ID }}

  deploy-prod:
    needs: validate
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: production
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: infrastructure/package-lock.json
      
      - name: Install dependencies
        run: |
          cd infrastructure
          npm ci
      
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID_PROD }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY_PROD }}
          aws-region: ap-northeast-1
      
      - name: Deploy to Production
        run: |
          cd infrastructure
          npx cdk deploy --context env=prod --all --require-approval never
        env:
          PROD_AUTH0_CLIENT_ID: ${{ secrets.PROD_AUTH0_CLIENT_ID }}
```

## セキュリティ設定

### IAMロール設計

```typescript
// Lambda実行ロール
const lambdaExecutionRole = new Role(this, 'LambdaExecutionRole', {
  assumedBy: new ServicePrincipal('lambda.amazonaws.com'),
  managedPolicies: [
    ManagedPolicy.fromAwsManagedPolicyName('service-role/AWSLambdaVPCAccessExecutionRole'),
  ],
  inlinePolicies: {
    SecretsManagerAccess: new PolicyDocument({
      statements: [
        new PolicyStatement({
          actions: [
            'secretsmanager:GetSecretValue',
          ],
          resources: [
            `arn:aws:secretsmanager:${this.region}:${this.account}:secret:finsight-db-credentials-*`,
          ],
        }),
      ],
    }),
  },
});
```

### Secrets Manager設定

```typescript
// Auth0設定をSecrets Managerで管理
const auth0Secret = new Secret(this, 'Auth0Secret', {
  secretName: `finsight-auth0-config-${props.environment}`,
  secretObjectValue: {
    domain: SecretValue.unsafePlainText(props.auth0Domain),
    audience: SecretValue.unsafePlainText(props.auth0Audience),
    clientId: SecretValue.unsafePlainText(props.auth0ClientId),
  },
});
```

### WAF設定

```typescript
// Web Application Firewall
const webAcl = new CfnWebACL(this, 'FinSightWebAcl', {
  scope: 'REGIONAL',
  defaultAction: { allow: {} },
  rules: [
    {
      name: 'RateLimitRule',
      priority: 1,
      statement: {
        rateBasedStatement: {
          limit: 2000,
          aggregateKeyType: 'IP',
        },
      },
      action: { block: {} },
      visibilityConfig: {
        sampledRequestsEnabled: true,
        cloudWatchMetricsEnabled: true,
        metricName: 'RateLimitRule',
      },
    },
  ],
});
```

## 運用・監視

### CloudWatchアラーム

```typescript
// Lambda関数エラーアラーム
Object.entries(this.lambdaFunctions).forEach(([name, func]) => {
  new Alarm(this, `${name}LambdaErrorAlarm`, {
    metric: func.metricErrors(),
    threshold: 5,
    evaluationPeriods: 2,
    alarmDescription: `Lambda function ${name} error rate is high`,
  });
});

// データベースCPU使用率アラーム
new Alarm(this, 'DatabaseCpuAlarm', {
  metric: props.database.metricCPUUtilization(),
  threshold: 80,
  evaluationPeriods: 2,
  alarmDescription: 'Database CPU utilization is high',
});
```

### ログ設定

```typescript
// Lambda関数のログ保持期間設定
Object.values(this.lambdaFunctions).forEach(func => {
  new LogGroup(this, `${func.functionName}LogGroup`, {
    logGroupName: `/aws/lambda/${func.functionName}`,
    retention: RetentionDays.TWO_WEEKS,
  });
});
```

## コスト最適化

### Lambda設定

```typescript
// 環境別のLambda設定
const lambdaConfig = props.environment === 'prod' ? {
  memorySize: 1024,
  reservedConcurrentExecutions: 100,
} : {
  memorySize: 512,
  reservedConcurrentExecutions: 10,
};
```

### RDS設定

```typescript
// 環境別のRDS設定
const rdsConfig = props.environment === 'prod' ? {
  instanceType: InstanceType.of(InstanceClass.T3, InstanceSize.SMALL),
  multiAz: true,
  backupRetention: Duration.days(7),
} : {
  instanceType: InstanceType.of(InstanceClass.T3, InstanceSize.MICRO),
  multiAz: false,
  backupRetention: Duration.days(3),
};
```

## 災害復旧

### バックアップ設定

```typescript
// RDS自動バックアップ
const database = new DatabaseInstance(this, 'Database', {
  backupRetention: Duration.days(7),
  preferredBackupWindow: '03:00-04:00',
  preferredMaintenanceWindow: 'Sun:04:00-Sun:05:00',
});

// Lambda関数のバージョニング
const lambdaVersion = func.currentVersion;
const lambdaAlias = new Alias(this, 'LambdaAlias', {
  aliasName: 'live',
  version: lambdaVersion,
});
```

## デプロイ後の確認手順

### 1. インフラストラクチャ確認

```bash
# スタック一覧確認
cdk list --context env=dev

# リソース確認
aws cloudformation describe-stacks --stack-name VpcStack-dev
aws cloudformation describe-stacks --stack-name DatabaseStack-dev
aws cloudformation describe-stacks --stack-name ApiStack-dev
```

### 2. 接続テスト

```bash
# API Gateway動作確認
curl -H "Authorization: Bearer YOUR_TOKEN" \
  https://api-dev.finsight.com/v1/users/me

# データベース接続確認
aws rds describe-db-instances --db-instance-identifier finsight-dev
```

### 3. 監視設定確認

```bash
# CloudWatchアラーム確認
aws cloudwatch describe-alarms

# X-Rayトレース確認
aws xray get-service-graph
```

## トラブルシューティング

### よくある問題と解決方法

1. **デプロイタイムアウト**
   - Lambda関数のタイムアウト設定を確認
   - VPC設定によるコールドスタート遅延

2. **データベース接続エラー**
   - セキュリティグループ設定を確認
   - Secrets Manager権限を確認

3. **Amplifyビルドエラー**
   - 環境変数設定を確認
   - Node.jsバージョン互換性を確認

### ログ確認コマンド

```bash
# Lambda関数ログ
aws logs tail /aws/lambda/finsight-users-dev --follow

# API Gatewayアクセスログ
aws logs tail /aws/apigateway/finsight-api-dev --follow

# RDSエラーログ
aws rds describe-db-log-files --db-instance-identifier finsight-dev
```

## SES機能詳細

### メール送信機能の要件対応
- **月次レポート送信**: settings画面で設定した日時に自動送信
- **予算超過アラート**: 閾値超過時のリアルタイム通知  
- **ドメイン認証**: 本番環境ではカスタムドメインでDKIM設定
- **バウンス処理**: SNSトピックで配信失敗を監視

### 環境別設定
- **開発**: 検証済みメールアドレス使用、200通/月
- **ステージング**: ドメイン認証、1,000通/月  
- **本番**: 完全なドメイン認証、10,000通/月

## まとめ

このCDKプランにより、SESを含むFinSightアプリケーションの完全なインフラストラクチャを自動化してデプロイできます。環境別の設定により、開発からプロダクションまで一貫したデプロイメントプロセスを提供し、セキュリティとコスト効率を両立させた設計となっています。