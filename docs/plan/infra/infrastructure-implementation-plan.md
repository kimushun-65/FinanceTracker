# FinSight インフラストラクチャ実装詳細計画

## 概要

FinSightアプリケーションのAWSインフラストラクチャを段階的に構築するための詳細な実装計画です。AWS CDKを使用したInfrastructure as Codeアプローチで、開発・ステージング・本番環境を一貫してデプロイします。

## 実装期間: 14日間（Phase 1: 週1-2）

### 📅 **Day 1-2: プロジェクト初期化とVPC基盤**

#### Day 1: 環境セットアップ（作業時間: 6時間）

**目標**: 開発環境の準備とCDKプロジェクトの初期化

**午前（3時間）**:
1. **AWS環境準備**
   ```bash
   # AWS CLI設定
   aws configure
   aws sts get-caller-identity  # 認証確認
   
   # AWS CDK CLI インストール
   npm install -g aws-cdk@latest
   cdk --version
   ```

2. **プロジェクト構造作成**
   ```bash
   mkdir -p FinanceTracker/infrastructure
   cd FinanceTracker/infrastructure
   
   # CDK プロジェクト初期化
   cdk init app --language typescript
   ```

3. **依存関係インストール**
   ```bash
   npm install @aws-cdk/aws-ec2 \
               @aws-cdk/aws-rds \
               @aws-cdk/aws-secretsmanager \
               @aws-cdk/aws-lambda \
               @aws-cdk/aws-apigateway \
               @aws-cdk/aws-amplify \
               @aws-cdk/aws-ses \
               @aws-cdk/aws-cloudwatch \
               @aws-cdk/aws-certificatemanager \
               @aws-cdk/aws-route53 \
               @aws-cdk/aws-iam
   ```

**午後（3時間）**:
4. **設定ファイル作成**
   ```bash
   mkdir -p config lib/stacks lib/constructs lib/interfaces scripts
   
   # 環境設定ファイル作成
   touch config/{dev,staging,prod}.json
   touch lib/interfaces/config.ts
   ```

5. **環境設定の実装**
   ```typescript
   // lib/interfaces/config.ts
   export interface EnvironmentConfig {
     environment: string;
     region: string;
     customDomain?: string;
     auth0Domain: string;
     auth0Audience: string;
     auth0ClientId: string;
     githubOwner: string;
     repositoryName: string;
     databaseConfig: DatabaseConfig;
     lambdaConfig: LambdaConfig;
     sesConfig: SesConfig;
   }
   ```

**成果物**:
- [ ] CDKプロジェクト初期化完了
- [ ] 依存関係インストール完了
- [ ] 基本ディレクトリ構造作成
- [ ] 環境設定インターフェース定義

**検証**:
```bash
# CDK動作確認
cdk list
cdk synth
```

#### Day 2: VPCスタック実装（作業時間: 8時間）

**目標**: ネットワーク基盤の構築

**午前（4時間）**:
1. **VPCスタック実装**
   ```typescript
   // lib/stacks/vpc-stack.ts
   import { Stack, StackProps, Construct } from '@aws-cdk/core';
   import { Vpc, SubnetType, SecurityGroup, Port } from '@aws-cdk/aws-ec2';
   
   export interface VpcStackProps extends StackProps {
     environment: string;
   }
   
   export class VpcStack extends Stack {
     public readonly vpc: Vpc;
     public readonly lambdaSecurityGroup: SecurityGroup;
     public readonly rdsSecurityGroup: SecurityGroup;
   
     constructor(scope: Construct, id: string, props: VpcStackProps) {
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
         natGateways: props.environment === 'prod' ? 2 : 1,
       });
   
       // セキュリティグループ作成
       this.lambdaSecurityGroup = new SecurityGroup(this, 'LambdaSecurityGroup', {
         vpc: this.vpc,
         description: 'Security group for Lambda functions',
         allowAllOutbound: true,
       });
   
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

**午後（4時間）**:
2. **メインアプリファイル更新**
   ```typescript
   // bin/finsight.ts
   import * as cdk from '@aws-cdk/core';
   import { VpcStack } from '../lib/stacks/vpc-stack';
   import { EnvironmentConfig } from '../lib/interfaces/config';
   
   const app = new cdk.App();
   const environment = app.node.tryGetContext('env') || 'dev';
   const config: EnvironmentConfig = require(`../config/${environment}.json`);
   
   // VPCスタック
   const vpcStack = new VpcStack(app, `VpcStack-${environment}`, {
     environment,
     env: {
       account: process.env.CDK_DEFAULT_ACCOUNT,
       region: config.region,
     },
   });
   ```

3. **設定ファイル実装**
   ```json
   // config/dev.json
   {
     "environment": "dev",
     "region": "ap-northeast-1",
     "customDomain": null,
     "auth0Domain": "finsight-dev.auth0.com",
     "auth0Audience": "https://api-dev.finsight.local",
     "auth0ClientId": "${DEV_AUTH0_CLIENT_ID}",
     "githubOwner": "your-github-username",
     "repositoryName": "finsight",
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
       "fromEmail": "noreply@dev.finsight.local",
       "sendingQuota": 200,
       "sendingRate": 1
     }
   }
   ```

**成果物**:
- [ ] VPCスタック実装完了
- [ ] セキュリティグループ設定完了
- [ ] 環境設定ファイル作成完了
- [ ] メインアプリファイル更新完了

**検証**:
```bash
# VPCスタック検証
cdk synth VpcStack-dev
cdk deploy VpcStack-dev --require-approval never

# リソース確認
aws ec2 describe-vpcs --filters "Name=tag:Name,Values=VpcStack-dev/FinSightVpc"
```

### 📅 **Day 3-4: データベース基盤**

#### Day 3: RDS設定とSecrets Manager（作業時間: 8時間）

**目標**: PostgreSQLデータベースとシークレット管理の構築

**午前（4時間）**:
1. **Secrets Manager設定**
   ```typescript
   // lib/stacks/database-stack.ts
   import { Stack, StackProps, Construct, Duration } from '@aws-cdk/core';
   import { 
     DatabaseInstance, 
     DatabaseInstanceEngine, 
     PostgresEngineVersion,
     Credentials,
     InstanceType,
     InstanceClass,
     InstanceSize,
     SubnetGroup
   } from '@aws-cdk/aws-rds';
   import { Secret } from '@aws-cdk/aws-secretsmanager';
   import { Vpc, SecurityGroup, SubnetType } from '@aws-cdk/aws-ec2';
   
   export interface DatabaseStackProps extends StackProps {
     vpc: Vpc;
     rdsSecurityGroup: SecurityGroup;
     environment: string;
   }
   
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
     }
   }
   ```

**午後（4時間）**:
2. **RDS PostgreSQL実装**
   ```typescript
   // lib/stacks/database-stack.ts (続き)
   
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
     removalPolicy: props.environment === 'prod' 
       ? RemovalPolicy.RETAIN 
       : RemovalPolicy.DESTROY,
   });
   ```

**成果物**:
- [ ] DatabaseStack実装完了
- [ ] Secrets Manager設定完了
- [ ] RDS PostgreSQL設定完了
- [ ] セキュリティ設定適用完了

**検証**:
```bash
# データベーススタック検証
cdk deploy DatabaseStack-dev --require-approval never

# RDS確認
aws rds describe-db-instances --db-instance-identifier $(aws rds describe-db-instances --query 'DBInstances[0].DBInstanceIdentifier' --output text)
```

#### Day 4: データベース初期化Lambda（作業時間: 6時間）

**目標**: データベース初期化とマイグレーション機能

**午前（3時間）**:
1. **Lambda初期化関数作成**
   ```bash
   mkdir -p lambda/db-init
   cd lambda/db-init
   npm init -y
   npm install pg aws-sdk
   ```

   ```javascript
   // lambda/db-init/index.js
   const { Client } = require('pg');
   const AWS = require('aws-sdk');
   
   const secretsManager = new AWS.SecretsManager();
   
   exports.handler = async (event) => {
     try {
       // Secrets Managerから認証情報取得
       const secretValue = await secretsManager.getSecretValue({
         SecretId: process.env.DB_SECRET_ARN
       }).promise();
       
       const secret = JSON.parse(secretValue.SecretString);
       
       // PostgreSQL接続
       const client = new Client({
         host: process.env.DB_ENDPOINT,
         port: 5432,
         user: secret.username,
         password: secret.password,
         database: 'finsight',
       });
       
       await client.connect();
       
       // 基本テーブル作成
       await client.query(`
         CREATE TABLE IF NOT EXISTS users (
           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
           auth0_user_id VARCHAR(255) UNIQUE NOT NULL,
           email VARCHAR(255) UNIQUE NOT NULL,
           name VARCHAR(255) NOT NULL,
           created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
           updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
         );
       `);
       
       await client.end();
       
       return {
         statusCode: 200,
         body: JSON.stringify({ message: 'Database initialized successfully' })
       };
     } catch (error) {
       console.error('Database initialization failed:', error);
       throw error;
     }
   };
   ```

**午後（3時間）**:
2. **Database Stack に初期化Lambda追加**
   ```typescript
   // lib/stacks/database-stack.ts (追加)
   import { Function, Runtime, Code } from '@aws-cdk/aws-lambda';
   
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
   ```

**成果物**:
- [ ] データベース初期化Lambda作成
- [ ] 基本テーブルDDL実装
- [ ] Secrets Manager統合
- [ ] Lambda権限設定

**検証**:
```bash
# Lambda実行テスト
aws lambda invoke --function-name $(aws lambda list-functions --query 'Functions[?contains(FunctionName, `DbInit`)].FunctionName' --output text) /tmp/response.json
```

### 📅 **Day 5-7: API Gateway + Lambda基盤**

#### Day 5: API Gateway設定（作業時間: 8時間）

**目標**: REST API基盤とCORS設定

**午前（4時間）**:
1. **API Gateway Stack作成**
   ```typescript
   // lib/stacks/api-stack.ts
   import { Stack, StackProps, Construct, Duration } from '@aws-cdk/core';
   import { 
     RestApi, 
     Cors, 
     TokenAuthorizer,
     LambdaIntegration,
     MethodOptions 
   } from '@aws-cdk/aws-apigateway';
   import { Function, Runtime, Code } from '@aws-cdk/aws-lambda';
   import { Certificate } from '@aws-cdk/aws-certificatemanager';
   import { Secret } from '@aws-cdk/aws-secretsmanager';
   import { Vpc, SecurityGroup } from '@aws-cdk/aws-ec2';
   
   export interface ApiStackProps extends StackProps {
     vpc: Vpc;
     lambdaSecurityGroup: SecurityGroup;
     dbSecret: Secret;
     environment: string;
     customDomain?: string;
     certificate?: Certificate;
     auth0Domain: string;
     auth0Audience: string;
     sesFromEmail: string;
   }
   
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
     }
   }
   ```

**午後（4時間）**:
2. **Lambda Authorizer実装**
   ```bash
   mkdir -p lambda/authorizer
   cd lambda/authorizer
   go mod init authorizer
   ```

   ```go
   // lambda/authorizer/main.go
   package main
   
   import (
       "context"
       "encoding/json"
       "fmt"
       "github.com/aws/aws-lambda-go/events"
       "github.com/aws/aws-lambda-go/lambda"
       "github.com/golang-jwt/jwt/v4"
       "os"
   )
   
   type Response struct {
       PrincipalID    string                 `json:"principalId"`
       PolicyDocument PolicyDocument         `json:"policyDocument"`
       Context        map[string]interface{} `json:"context,omitempty"`
   }
   
   type PolicyDocument struct {
       Version   string      `json:"Version"`
       Statement []Statement `json:"Statement"`
   }
   
   type Statement struct {
       Action   string `json:"Action"`
       Effect   string `json:"Effect"`
       Resource string `json:"Resource"`
   }
   
   func handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (Response, error) {
       token := event.AuthorizationToken
       
       // "Bearer " プレフィックスを削除
       if len(token) > 7 && token[:7] == "Bearer " {
           token = token[7:]
       }
       
       // JWT検証（簡易版 - 実際はAuth0の公開鍵で検証）
       claims := jwt.MapClaims{}
       _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
           // 実際の実装ではAuth0のJWKSから公開鍵を取得
           return []byte(os.Getenv("JWT_SECRET")), nil
       })
       
       if err != nil {
           return Response{}, fmt.Errorf("unauthorized")
       }
       
       return Response{
           PrincipalID: claims["sub"].(string),
           PolicyDocument: PolicyDocument{
               Version: "2012-10-17",
               Statement: []Statement{
                   {
                       Action:   "execute-api:Invoke",
                       Effect:   "Allow",
                       Resource: event.MethodArn,
                   },
               },
           },
           Context: map[string]interface{}{
               "userId": claims["sub"].(string),
           },
       }, nil
   }
   
   func main() {
       lambda.Start(handler)
   }
   ```

**成果物**:
- [ ] API Gateway Stack実装
- [ ] CORS設定完了
- [ ] Lambda Authorizer実装
- [ ] JWT検証ロジック実装

#### Day 6: Lambda関数群実装（作業時間: 8時間）

**目標**: 各APIエンドポイント用Lambda関数作成

**全日（8時間）**:
1. **API用Lambda関数作成**
   ```bash
   # 各API用ディレクトリ作成
   mkdir -p lambda/{auth,users,transactions,budgets,accounts,categories,reports,notifications}
   
   # 共通ライブラリ作成
   mkdir -p lambda/shared
   ```

2. **共通ライブラリ実装**
   ```go
   // lambda/shared/go.mod
   module shared
   
   go 1.21
   
   require (
       github.com/aws/aws-lambda-go v1.41.0
       github.com/gin-gonic/gin v1.9.1
       gorm.io/gorm v1.25.5
       gorm.io/driver/postgres v1.5.4
   )
   ```

   ```go
   // lambda/shared/database.go
   package shared
   
   import (
       "fmt"
       "gorm.io/driver/postgres"
       "gorm.io/gorm"
   )
   
   func NewDB(host, user, password, dbname string, port int) (*gorm.DB, error) {
       dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require TimeZone=Asia/Tokyo",
           host, user, password, dbname, port)
       
       return gorm.Open(postgres.Open(dsn), &gorm.Config{})
   }
   ```

3. **ユーザーAPI Lambda実装**
   ```go
   // lambda/users/main.go
   package main
   
   import (
       "context"
       "encoding/json"
       "github.com/aws/aws-lambda-go/events"
       "github.com/aws/aws-lambda-go/lambda"
       "shared"
   )
   
   func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
       // データベース接続
       db, err := shared.NewDB(
           os.Getenv("DB_HOST"),
           os.Getenv("DB_USER"),
           os.Getenv("DB_PASSWORD"),
           os.Getenv("DB_NAME"),
           5432,
       )
       if err != nil {
           return events.APIGatewayProxyResponse{
               StatusCode: 500,
               Body:       `{"error": "Database connection failed"}`,
           }, nil
       }
   
       switch request.HTTPMethod {
       case "GET":
           return getUserProfile(ctx, request, db)
       case "PUT":
           return updateUserProfile(ctx, request, db)
       default:
           return events.APIGatewayProxyResponse{
               StatusCode: 405,
               Body:       `{"error": "Method not allowed"}`,
           }, nil
       }
   }
   
   func main() {
       lambda.Start(handler)
   }
   ```

**成果物**:
- [ ] 8つのLambda関数ディレクトリ作成
- [ ] 共通ライブラリ実装
- [ ] ユーザーAPI Lambda実装
- [ ] データベース接続ロジック

#### Day 7: API エンドポイント設定（作業時間: 8時間）

**目標**: API Gateway リソースとメソッド設定

**午前（4時間）**:
1. **API Stack Lambda関数作成部分実装**
   ```typescript
   // lib/stacks/api-stack.ts (続き)
   
   private createApiLambdas(props: ApiStackProps): { [key: string]: Function } {
     const functions: { [key: string]: Function } = {};
   
     const lambdaConfigs = [
       { name: 'auth', path: 'auth', needsAuth: false },
       { name: 'users', path: 'users', needsAuth: true },
       { name: 'transactions', path: 'transactions', needsAuth: true },
       { name: 'budgets', path: 'budgets', needsAuth: true },
       { name: 'accounts', path: 'accounts', needsAuth: true },
       { name: 'categories', path: 'categories', needsAuth: true },
       { name: 'reports', path: 'reports', needsAuth: true },
       { name: 'notifications', path: 'notifications', needsAuth: true },
     ];
   
     lambdaConfigs.forEach(config => {
       functions[config.name] = new Function(this, `${config.name}Function`, {
         runtime: Runtime.PROVIDED_AL2, // Go runtime
         handler: 'bootstrap',
         code: Code.fromAsset(`lambda/${config.path}`),
         vpc: props.vpc,
         securityGroups: [props.lambdaSecurityGroup],
         environment: {
           DB_SECRET_ARN: props.dbSecret.secretArn,
           ENVIRONMENT: props.environment,
           AUTH0_DOMAIN: props.auth0Domain,
           AUTH0_AUDIENCE: props.auth0Audience,
         },
         timeout: Duration.seconds(30),
         memorySize: 512,
       });
   
       // Secrets Manager読み取り権限
       props.dbSecret.grantRead(functions[config.name]);
     });
   
     return functions;
   }
   ```

**午後（4時間）**:
2. **API エンドポイント設定**
   ```typescript
   // lib/stacks/api-stack.ts (続き)
   
   private setupApiEndpoints(authorizer: TokenAuthorizer): void {
     // ヘルスチェック（認証不要）
     const healthResource = this.api.root.addResource('health');
     healthResource.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.auth));
   
     // /auth - 認証不要
     const authResource = this.api.root.addResource('auth');
     authResource.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.auth));
     
     const callbackResource = authResource.addResource('callback');
     callbackResource.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.auth));
     
     const syncResource = authResource.addResource('sync');
     syncResource.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.auth), {
       authorizer,
     });
   
     // /users - 認証必要
     const usersResource = this.api.root.addResource('users');
     const meResource = usersResource.addResource('me');
     meResource.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.users), {
       authorizer,
     });
     meResource.addMethod('PUT', new LambdaIntegration(this.lambdaFunctions.users), {
       authorizer,
     });
   
     // /transactions - 認証必要
     const transactionsResource = this.api.root.addResource('transactions');
     transactionsResource.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.transactions), {
       authorizer,
     });
     transactionsResource.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.transactions), {
       authorizer,
     });
     
     const transactionResource = transactionsResource.addResource('{transactionId}');
     transactionResource.addMethod('PUT', new LambdaIntegration(this.lambdaFunctions.transactions), {
       authorizer,
     });
     transactionResource.addMethod('DELETE', new LambdaIntegration(this.lambdaFunctions.transactions), {
       authorizer,
     });
   
     // 残りのエンドポイントも同様に設定...
   }
   ```

**成果物**:
- [ ] Lambda関数作成ロジック実装
- [ ] 認証・認可設定
- [ ] 全APIエンドポイント設定
- [ ] CRUD操作パス設定

**検証**:
```bash
# API Stack デプロイ
cdk deploy ApiStack-dev --require-approval never

# API Gateway動作確認
curl https://$(aws apigateway get-rest-apis --query 'items[0].id' --output text).execute-api.ap-northeast-1.amazonaws.com/prod/health
```

### 📅 **Day 8-10: フロントエンドホスティングとSES**

#### Day 8: SSL証明書設定（作業時間: 6時間）

**目標**: SSL/TLS証明書管理とドメイン設定

**午前（3時間）**:
1. **Certificate Stack実装**
   ```typescript
   // lib/stacks/certificate-stack.ts
   import { Stack, StackProps, Construct } from '@aws-cdk/core';
   import { Certificate, CertificateValidation } from '@aws-cdk/aws-certificatemanager';
   import { HostedZone } from '@aws-cdk/aws-route53';
   
   export interface CertificateStackProps extends StackProps {
     domainName: string;
     environment: string;
   }
   
   export class CertificateStack extends Stack {
     public readonly certificate: Certificate;
     public readonly hostedZone: HostedZone;
   
     constructor(scope: Construct, id: string, props: CertificateStackProps) {
       super(scope, id, props);
   
       // Route53 ホストゾーン（既存のものを参照）
       this.hostedZone = HostedZone.fromLookup(this, 'HostedZone', {
         domainName: props.domainName,
       });
   
       // SSL証明書作成
       this.certificate = new Certificate(this, 'FinSightCertificate', {
         domainName: props.domainName,
         subjectAlternativeNames: [
           `*.${props.domainName}`,
           `api.${props.domainName}`,
           `www.${props.domainName}`,
         ],
         validation: CertificateValidation.fromDns(this.hostedZone),
       });
     }
   }
   ```

**午後（3時間）**:
2. **メインアプリに証明書スタック追加**
   ```typescript
   // bin/finsight.ts (更新)
   import { CertificateStack } from '../lib/stacks/certificate-stack';
   
   // 証明書スタック（カスタムドメインがある場合のみ）
   let certificateStack: CertificateStack | undefined;
   if (config.customDomain) {
     certificateStack = new CertificateStack(app, `CertificateStack-${environment}`, {
       domainName: config.customDomain,
       environment,
       env: {
         account: process.env.CDK_DEFAULT_ACCOUNT,
         region: 'us-east-1', // CloudFrontで使用するため
       },
     });
   }
   ```

**成果物**:
- [ ] Certificate Stack実装
- [ ] Route53統合
- [ ] ワイルドカード証明書設定
- [ ] DNS検証設定

#### Day 9: Amplifyホスティング設定（作業時間: 8時間）

**目標**: フロントエンドホスティング環境構築

**午前（4時間）**:
1. **Amplify Stack実装**
   ```typescript
   // lib/stacks/amplify-stack.ts
   import { Stack, StackProps, Construct } from '@aws-cdk/core';
   import { App, GitHubSourceCodeProvider } from '@aws-cdk/aws-amplify';
   import { SecretValue } from '@aws-cdk/core';
   import { Certificate } from '@aws-cdk/aws-certificatemanager';
   
   export interface AmplifyStackProps extends StackProps {
     environment: string;
     apiUrl: string;
     auth0Domain: string;
     auth0ClientId: string;
     githubOwner: string;
     repositoryName: string;
     customDomain?: string;
     certificate?: Certificate;
   }
   
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
           NEXT_PUBLIC_API_URL: props.apiUrl,
           NEXT_PUBLIC_AUTH0_DOMAIN: props.auth0Domain,
           NEXT_PUBLIC_AUTH0_CLIENT_ID: props.auth0ClientId,
           NEXT_PUBLIC_ENVIRONMENT: props.environment,
         },
         buildSpec: {
           version: '1.0',
           frontend: {
             phases: {
               preBuild: {
                 commands: [
                   'npm ci',
                 ],
               },
               build: {
                 commands: [
                   'npm run build',
                 ],
               },
             },
             artifacts: {
               baseDirectory: '.next',
               files: ['**/*'],
             },
           },
         },
       });
   
       // ブランチ設定
       const mainBranch = this.app.addBranch('main', {
         autoBuild: true,
         stage: props.environment === 'prod' ? 'PRODUCTION' : 'DEVELOPMENT',
         environmentVariables: this.app.environmentVariables,
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

**午後（4時間）**:
2. **GitHub連携とビルド設定**
   ```bash
   # GitHub Personal Access Token をSecrets Managerに保存
   aws secretsmanager create-secret \
     --name github-token \
     --description "GitHub Personal Access Token for Amplify" \
     --secret-string '{"token":"ghp_your_token_here"}'
   ```

3. **フロントエンドプロジェクト準備**
   ```bash
   # Next.js プロジェクト作成（別作業）
   mkdir -p ../frontend
   cd ../frontend
   npx create-next-app@latest . --typescript --tailwind --app --src-dir --import-alias "@/*"
   
   # Auth0依存関係追加
   npm install @auth0/nextjs-auth0
   ```

**成果物**:
- [ ] Amplify Stack実装
- [ ] GitHub連携設定
- [ ] 環境変数設定
- [ ] ビルド設定実装
- [ ] カスタムドメイン設定

#### Day 10: SES設定（作業時間: 8時間）

**目標**: メール送信機能の構築

**午前（4時間）**:
1. **SES Stack実装**
   ```typescript
   // lib/stacks/ses-stack.ts
   import { Stack, StackProps, Construct } from '@aws-cdk/core';
   import { 
     EmailIdentity, 
     Identity, 
     ConfigurationSet,
     EventDestination,
     EventType 
   } from '@aws-cdk/aws-ses';
   import { Topic } from '@aws-cdk/aws-sns';
   import { PolicyDocument, PolicyStatement } from '@aws-cdk/aws-iam';
   
   export interface SesStackProps extends StackProps {
     environment: string;
     customDomain?: string;
     fromEmail: string;
   }
   
   export class SesStack extends Stack {
     public readonly verifiedDomain: string;
     public readonly sesIdentity: EmailIdentity;
     public readonly configurationSet: ConfigurationSet;
   
     constructor(scope: Construct, id: string, props: SesStackProps) {
       super(scope, id, props);
   
       // 送信設定セット
       this.configurationSet = new ConfigurationSet(this, 'SesConfigurationSet', {
         configurationSetName: `finsight-${props.environment}`,
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
   
       // ドメイン認証（カスタムドメインがある場合）
       if (props.customDomain) {
         this.sesIdentity = new EmailIdentity(this, 'SesEmailIdentity', {
           identity: Identity.domain(props.customDomain),
           dkimSigning: true,
           configurationSet: this.configurationSet,
         });
         this.verifiedDomain = props.customDomain;
       } else {
         // 開発環境用：検証済みメールアドレス
         this.sesIdentity = new EmailIdentity(this, 'SesEmailIdentity', {
           identity: Identity.email(props.fromEmail),
           configurationSet: this.configurationSet,
         });
         this.verifiedDomain = props.fromEmail;
       }
     }
   }
   ```

**午後（4時間）**:
2. **SES監視とアラート設定**
   ```typescript
   // lib/stacks/ses-stack.ts (続き)
   
   // バウンス・苦情処理用SNSトピック
   const bouncesTopic = new Topic(this, 'SesBouncestopic', {
     topicName: `finsight-ses-bounces-${props.environment}`,
   });
   
   const complaintsTopic = new Topic(this, 'SesComplaintsTopic', {
     topicName: `finsight-ses-complaints-${props.environment}`,
   });
   
   // SES送信統計をCloudWatchに送信
   this.configurationSet.addEventDestination('CloudWatchDestination', {
     destination: EventDestination.cloudWatchDimensions({
       source: 'messageTag',
       defaultValue: 'default',
     }),
     events: [EventType.SEND, EventType.REJECT, EventType.BOUNCE, EventType.COMPLAINT, EventType.DELIVERY],
   });
   
   // バウンス・苦情をSNSに送信
   this.configurationSet.addEventDestination('SnsDestination', {
     destination: EventDestination.snsTopic(bouncesTopic),
     events: [EventType.BOUNCE, EventType.COMPLAINT],
   });
   ```

3. **Lambda関数にSES権限追加**
   ```typescript
   // lib/stacks/api-stack.ts (SES権限追加)
   
   // SES送信権限（reportsとnotifications関数のみ）
   if (config.name === 'reports' || config.name === 'notifications') {
     functions[config.name].addToRolePolicy(new PolicyStatement({
       actions: [
         'ses:SendEmail',
         'ses:SendRawEmail',
         'ses:GetSendQuota',
         'ses:GetSendStatistics',
       ],
       resources: [
         `arn:aws:ses:${this.region}:${this.account}:identity/${props.sesFromEmail}`,
         `arn:aws:ses:${this.region}:${this.account}:configuration-set/finsight-${props.environment}`,
       ],
     }));
   }
   ```

**成果物**:
- [ ] SES Stack実装
- [ ] 設定セット作成
- [ ] ドメイン/メール認証設定
- [ ] 監視・アラート設定
- [ ] Lambda SES権限設定

**検証**:
```bash
# SES設定確認
aws ses get-configuration-set --configuration-set-name finsight-dev
aws ses list-identities
```

### 📅 **Day 11-12: 監視設定とデプロイ自動化**

#### Day 11: CloudWatch監視設定（作業時間: 8時間）

**目標**: 包括的な監視とアラート構築

**午前（4時間）**:
1. **Monitoring Stack実装**
   ```typescript
   // lib/stacks/monitoring-stack.ts
   import { Stack, StackProps, Construct, Duration } from '@aws-cdk/core';
   import { 
     Dashboard, 
     GraphWidget, 
     Metric, 
     Alarm, 
     ComparisonOperator,
     TreatMissingData 
   } from '@aws-cdk/aws-cloudwatch';
   import { Topic } from '@aws-cdk/aws-sns';
   import { EmailSubscription } from '@aws-cdk/aws-sns-subscriptions';
   import { RestApi } from '@aws-cdk/aws-apigateway';
   import { Function } from '@aws-cdk/aws-lambda';
   import { DatabaseInstance } from '@aws-cdk/aws-rds';
   
   export interface MonitoringStackProps extends StackProps {
     environment: string;
     api: RestApi;
     lambdaFunctions: { [key: string]: Function };
     database: DatabaseInstance;
     alertEmail: string;
   }
   
   export class MonitoringStack extends Stack {
     constructor(scope: Construct, id: string, props: MonitoringStackProps) {
       super(scope, id, props);
   
       // アラート用SNSトピック
       const alertTopic = new Topic(this, 'AlertTopic', {
         topicName: `finsight-alerts-${props.environment}`,
       });
       
       alertTopic.addSubscription(new EmailSubscription(props.alertEmail));
   
       // CloudWatchダッシュボード
       const dashboard = new Dashboard(this, 'FinSightDashboard', {
         dashboardName: `finsight-${props.environment}`,
       });
   
       // API Gatewayメトリクス
       dashboard.addWidgets(
         new GraphWidget({
           title: 'API Gateway Requests',
           left: [props.api.metricRequestCount()],
           right: [props.api.metricLatency()],
           width: 12,
         }),
         new GraphWidget({
           title: 'API Gateway Errors',
           left: [props.api.metricClientError(), props.api.metricServerError()],
           width: 12,
         })
       );
   
       // Lambda関数メトリクス
       const lambdaMetrics = Object.entries(props.lambdaFunctions).map(([name, func]) => 
         func.metricInvocations({ label: name })
       );
       
       dashboard.addWidgets(
         new GraphWidget({
           title: 'Lambda Invocations',
           left: lambdaMetrics,
           width: 24,
         })
       );
     }
   }
   ```

**午後（4時間）**:
2. **アラーム設定実装**
   ```typescript
   // lib/stacks/monitoring-stack.ts (続き)
   
   private createAlarms(props: MonitoringStackProps, alertTopic: Topic): void {
     // API Gatewayエラー率アラーム
     new Alarm(this, 'ApiErrorAlarm', {
       metric: props.api.metricClientError({
         statistic: 'Sum',
         period: Duration.minutes(5),
       }),
       threshold: 10,
       evaluationPeriods: 2,
       comparisonOperator: ComparisonOperator.GREATER_THAN_THRESHOLD,
       alarmDescription: 'API Gateway error rate is too high',
       treatMissingData: TreatMissingData.NOT_BREACHING,
     }).addAlarmAction(new SnsAction(alertTopic));
   
     // Lambda関数エラーアラーム
     Object.entries(props.lambdaFunctions).forEach(([name, func]) => {
       new Alarm(this, `${name}LambdaErrorAlarm`, {
         metric: func.metricErrors({
           period: Duration.minutes(5),
         }),
         threshold: 5,
         evaluationPeriods: 2,
         alarmDescription: `Lambda function ${name} error rate is high`,
       }).addAlarmAction(new SnsAction(alertTopic));
   
       // Lambda実行時間アラーム
       new Alarm(this, `${name}LambdaDurationAlarm`, {
         metric: func.metricDuration({
           period: Duration.minutes(5),
         }),
         threshold: 25000, // 25秒
         evaluationPeriods: 3,
         alarmDescription: `Lambda function ${name} execution time is too long`,
       }).addAlarmAction(new SnsAction(alertTopic));
     });
   
     // データベースアラーム
     new Alarm(this, 'DatabaseConnectionAlarm', {
       metric: props.database.metricDatabaseConnections({
         period: Duration.minutes(5),
       }),
       threshold: 80,
       evaluationPeriods: 2,
       alarmDescription: 'Database connection count is high',
     }).addAlarmAction(new SnsAction(alertTopic));
   
     new Alarm(this, 'DatabaseCpuAlarm', {
       metric: props.database.metricCPUUtilization({
         period: Duration.minutes(5),
       }),
       threshold: 80,
       evaluationPeriods: 3,
       alarmDescription: 'Database CPU utilization is high',
     }).addAlarmAction(new SnsAction(alertTopic));
   }
   ```

**成果物**:
- [ ] CloudWatchダッシュボード
- [ ] 包括的なアラーム設定
- [ ] SNS通知設定
- [ ] メトリクス可視化

#### Day 12: X-Ray設定とログ管理（作業時間: 6時間）

**目標**: 分散トレーシングとログ管理

**午前（3時間）**:
1. **X-Ray設定実装**
   ```typescript
   // lib/stacks/monitoring-stack.ts (X-Ray追加)
   
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
       
       // 環境変数でX-Ray有効化
       func.addEnvironment('_X_AMZN_TRACE_ID', '');
     });
   
     // API GatewayでX-Ray有効化
     props.api.deploymentStage.node.addDependency(
       new CfnStage(this, 'ApiStageXRay', {
         restApiId: props.api.restApiId,
         deploymentId: props.api.latestDeployment?.deploymentId || '',
         stageName: 'prod',
         tracingConfig: {
           tracingEnabled: true,
         },
       })
     );
   }
   ```

**午後（3時間）**:
2. **ログ管理設定**
   ```typescript
   // lib/stacks/monitoring-stack.ts (ログ管理追加)
   import { LogGroup, RetentionDays } from '@aws-cdk/aws-logs';
   
   private setupLogManagement(props: MonitoringStackProps): void {
     // Lambda関数のログ保持期間設定
     Object.entries(props.lambdaFunctions).forEach(([name, func]) => {
       new LogGroup(this, `${name}LogGroup`, {
         logGroupName: `/aws/lambda/${func.functionName}`,
         retention: props.environment === 'prod' 
           ? RetentionDays.ONE_MONTH 
           : RetentionDays.TWO_WEEKS,
       });
     });
   
     // API Gatewayアクセスログ設定
     const apiLogGroup = new LogGroup(this, 'ApiGatewayLogGroup', {
       logGroupName: `/aws/apigateway/${props.api.restApiName}`,
       retention: RetentionDays.TWO_WEEKS,
     });
   
     // アクセスログ形式設定
     props.api.deploymentStage.addPropertyOverride('AccessLogSetting', {
       DestinationArn: apiLogGroup.logGroupArn,
       Format: JSON.stringify({
         requestId: '$context.requestId',
         ip: '$context.identity.sourceIp',
         user: '$context.identity.user',
         requestTime: '$context.requestTime',
         httpMethod: '$context.httpMethod',
         resourcePath: '$context.resourcePath',
         status: '$context.status',
         protocol: '$context.protocol',
         responseLength: '$context.responseLength',
       }),
     });
   }
   ```

**成果物**:
- [ ] X-Ray分散トレーシング設定
- [ ] ログ保持期間設定
- [ ] API Gatewayアクセスログ
- [ ] 構造化ログ形式設定

### 📅 **Day 13-14: CI/CDとセキュリティ**

#### Day 13: GitHub Actions設定（作業時間: 8時間）

**目標**: 自動デプロイパイプラインの構築

**午前（4時間）**:
1. **GitHub Actions ワークフロー作成**
   ```bash
   mkdir -p .github/workflows
   ```

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
         - uses: actions/checkout@v4
         
         - name: Setup Node.js
           uses: actions/setup-node@v4
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
         - uses: actions/checkout@v4
         
         - name: Setup Node.js
           uses: actions/setup-node@v4
           with:
             node-version: '18'
             cache: 'npm'
             cache-dependency-path: infrastructure/package-lock.json
         
         - name: Install dependencies
           run: |
             cd infrastructure
             npm ci
         
         - name: Configure AWS credentials
           uses: aws-actions/configure-aws-credentials@v4
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
         - uses: actions/checkout@v4
         
         - name: Setup Node.js
           uses: actions/setup-node@v4
           with:
             node-version: '18'
             cache: 'npm'
             cache-dependency-path: infrastructure/package-lock.json
         
         - name: Install dependencies
           run: |
             cd infrastructure
             npm ci
         
         - name: Configure AWS credentials
           uses: aws-actions/configure-aws-credentials@v4
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

**午後（4時間）**:
2. **デプロイスクリプト作成**
   ```bash
   # scripts/deploy-dev.sh
   #!/bin/bash
   set -e
   
   echo "🚀 Deploying to Development Environment"
   
   # 環境変数チェック
   if [ -z "$DEV_AUTH0_CLIENT_ID" ]; then
     echo "❌ DEV_AUTH0_CLIENT_ID is not set"
     exit 1
   fi
   
   # CDK Bootstrap（初回のみ）
   echo "📦 Bootstrapping CDK..."
   npx cdk bootstrap aws://$(aws sts get-caller-identity --query Account --output text)/ap-northeast-1
   
   # スタック依存関係に基づく順序でデプロイ
   echo "🏗️  Deploying VPC Stack..."
   npx cdk deploy VpcStack-dev --context env=dev --require-approval never
   
   echo "🏗️  Deploying Database Stack..."
   npx cdk deploy DatabaseStack-dev --context env=dev --require-approval never
   
   echo "🏗️  Deploying API Stack..."
   npx cdk deploy ApiStack-dev --context env=dev --require-approval never
   
   echo "🏗️  Deploying SES Stack..."
   npx cdk deploy SesStack-dev --context env=dev --require-approval never
   
   echo "🏗️  Deploying Monitoring Stack..."
   npx cdk deploy MonitoringStack-dev --context env=dev --require-approval never
   
   if [ ! -z "$CUSTOM_DOMAIN" ]; then
     echo "🏗️  Deploying Certificate Stack..."
     npx cdk deploy CertificateStack-dev --context env=dev --require-approval never
     
     echo "🏗️  Deploying Amplify Stack..."
     npx cdk deploy AmplifyStack-dev --context env=dev --require-approval never
   fi
   
   echo "✅ Development deployment completed!"
   
   # エンドポイント情報出力
   echo "📋 Deployment Information:"
   aws cloudformation describe-stacks --stack-name ApiStack-dev --query 'Stacks[0].Outputs'
   ```

**成果物**:
- [ ] GitHub Actions ワークフロー
- [ ] 環境別デプロイスクリプト
- [ ] 依存関係管理
- [ ] シークレット管理設定

#### Day 14: セキュリティ設定とテスト（作業時間: 8時間）

**目標**: セキュリティ強化と動作検証

**午前（4時間）**:
1. **WAF設定実装**
   ```typescript
   // lib/stacks/security-stack.ts
   import { Stack, StackProps, Construct } from '@aws-cdk/core';
   import { CfnWebACL, CfnWebACLAssociation } from '@aws-cdk/aws-wafv2';
   import { RestApi } from '@aws-cdk/aws-apigateway';
   
   export interface SecurityStackProps extends StackProps {
     api: RestApi;
     environment: string;
   }
   
   export class SecurityStack extends Stack {
     constructor(scope: Construct, id: string, props: SecurityStackProps) {
       super(scope, id, props);
   
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
           {
             name: 'SQLInjectionRule',
             priority: 2,
             statement: {
               managedRuleGroupStatement: {
                 vendorName: 'AWS',
                 name: 'AWSManagedRulesSQLiRuleSet',
               },
             },
             action: { block: {} },
             visibilityConfig: {
               sampledRequestsEnabled: true,
               cloudWatchMetricsEnabled: true,
               metricName: 'SQLInjectionRule',
             },
           },
           {
             name: 'XSSRule',
             priority: 3,
             statement: {
               managedRuleGroupStatement: {
                 vendorName: 'AWS',
                 name: 'AWSManagedRulesCommonRuleSet',
               },
             },
             action: { block: {} },
             visibilityConfig: {
               sampledRequestsEnabled: true,
               cloudWatchMetricsEnabled: true,
               metricName: 'XSSRule',
             },
           },
         ],
       });
   
       // WAFをAPI Gatewayに関連付け
       new CfnWebACLAssociation(this, 'WebAclAssociation', {
         resourceArn: props.api.deploymentStage.stageArn,
         webAclArn: webAcl.attrArn,
       });
     }
   }
   ```

**午後（4時間）**:
2. **統合テスト実装**
   ```bash
   # 統合テスト用スクリプト作成
   touch scripts/test-deployment.sh
   chmod +x scripts/test-deployment.sh
   ```

   ```bash
   # scripts/test-deployment.sh
   #!/bin/bash
   set -e
   
   ENVIRONMENT=${1:-dev}
   echo "🧪 Testing $ENVIRONMENT environment deployment"
   
   # API Gateway エンドポイント取得
   API_URL=$(aws cloudformation describe-stacks \
     --stack-name ApiStack-$ENVIRONMENT \
     --query 'Stacks[0].Outputs[?OutputKey==`ApiGatewayUrl`].OutputValue' \
     --output text)
   
   if [ -z "$API_URL" ]; then
     echo "❌ API Gateway URL not found"
     exit 1
   fi
   
   echo "🔗 API URL: $API_URL"
   
   # ヘルスチェック
   echo "🏥 Testing health check..."
   HEALTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $API_URL/health)
   if [ "$HEALTH_RESPONSE" = "200" ]; then
     echo "✅ Health check passed"
   else
     echo "❌ Health check failed (HTTP $HEALTH_RESPONSE)"
     exit 1
   fi
   
   # データベース接続テスト
   echo "🗄️  Testing database connection..."
   DB_ENDPOINT=$(aws rds describe-db-instances \
     --query 'DBInstances[?DBName==`finsight`].Endpoint.Address' \
     --output text)
   
   if [ -z "$DB_ENDPOINT" ]; then
     echo "❌ Database endpoint not found"
     exit 1
   fi
   
   echo "✅ Database endpoint found: $DB_ENDPOINT"
   
   # Lambda関数テスト
   echo "⚡ Testing Lambda functions..."
   LAMBDA_FUNCTIONS=$(aws lambda list-functions \
     --query 'Functions[?contains(FunctionName, `finsight`)].FunctionName' \
     --output text)
   
   for FUNCTION in $LAMBDA_FUNCTIONS; do
     echo "  Testing $FUNCTION..."
     RESPONSE=$(aws lambda invoke \
       --function-name $FUNCTION \
       --payload '{"httpMethod":"GET","path":"/health"}' \
       /tmp/lambda-response.json 2>/dev/null)
     
     if [ $? -eq 0 ]; then
       echo "  ✅ $FUNCTION is working"
     else
       echo "  ❌ $FUNCTION failed"
     fi
   done
   
   # SES設定確認
   echo "📧 Testing SES configuration..."
   SES_IDENTITY=$(aws ses list-identities --output text)
   if [ -z "$SES_IDENTITY" ]; then
     echo "❌ No SES identities found"
   else
     echo "✅ SES identities found: $SES_IDENTITY"
   fi
   
   echo "🎉 All tests completed!"
   ```

3. **セキュリティ検証スクリプト**
   ```bash
   # scripts/security-check.sh
   #!/bin/bash
   set -e
   
   echo "🔒 Running security checks..."
   
   # IAMロール権限確認
   echo "👤 Checking IAM roles..."
   aws iam list-roles --query 'Roles[?contains(RoleName, `finsight`)].{RoleName:RoleName,AssumeRolePolicyDocument:AssumeRolePolicyDocument}'
   
   # セキュリティグループ確認
   echo "🛡️  Checking Security Groups..."
   aws ec2 describe-security-groups --filters "Name=group-name,Values=*finsight*" --query 'SecurityGroups[].{GroupName:GroupName,GroupId:GroupId,IpPermissions:IpPermissions}'
   
   # Secrets Manager確認
   echo "🔐 Checking Secrets Manager..."
   aws secretsmanager list-secrets --query 'SecretList[?contains(Name, `finsight`)].{Name:Name,Description:Description}'
   
   # RDS暗号化確認
   echo "🗄️  Checking RDS encryption..."
   aws rds describe-db-instances --query 'DBInstances[?DBName==`finsight`].{StorageEncrypted:StorageEncrypted,KmsKeyId:KmsKeyId}'
   
   echo "✅ Security check completed!"
   ```

**成果物**:
- [ ] WAF設定実装
- [ ] 統合テストスクリプト
- [ ] セキュリティ検証スクリプト
- [ ] 自動化されたテスト実行

**最終検証**:
```bash
# 全スタックデプロイテスト
cd infrastructure
npm test
./scripts/deploy-dev.sh

# セキュリティチェック
./scripts/security-check.sh

# 統合テスト実行
./scripts/test-deployment.sh dev
```

## 📋 実装チェックリスト

### Phase 1完了時点での成果物

#### インフラストラクチャ
- [ ] **VPC**: 2AZ構成、パブリック/プライベートサブネット
- [ ] **セキュリティグループ**: Lambda、RDS用の適切な設定
- [ ] **RDS PostgreSQL**: 暗号化、バックアップ設定済み
- [ ] **Secrets Manager**: データベース認証情報管理
- [ ] **API Gateway**: CORS、認証設定済み
- [ ] **Lambda関数**: 8つのAPI機能別関数
- [ ] **Lambda Authorizer**: Auth0 JWT検証
- [ ] **Amplify**: フロントエンドホスティング設定
- [ ] **SES**: メール送信、ドメイン認証
- [ ] **CloudWatch**: メトリクス、ログ、アラーム
- [ ] **X-Ray**: 分散トレーシング
- [ ] **WAF**: セキュリティルール設定

#### 自動化・運用
- [ ] **CDK**: Infrastructure as Code実装
- [ ] **GitHub Actions**: CI/CDパイプライン
- [ ] **デプロイスクリプト**: 環境別自動化
- [ ] **テストスクリプト**: 統合テスト自動化
- [ ] **監視設定**: 包括的なアラート設定

#### セキュリティ
- [ ] **暗号化**: 保存時・転送時暗号化
- [ ] **IAM**: 最小権限の原則適用
- [ ] **ネットワーク**: プライベートサブネット配置
- [ ] **WAF**: 基本的な攻撃対策
- [ ] **シークレット管理**: 認証情報の安全な管理

## 🚀 次フェーズへの準備

Phase 1完了後、以下のフェーズに進む準備が整います：

1. **Phase 2: バックエンド開発**
   - Goアプリケーション開発
   - クリーンアーキテクチャ実装
   - データベースモデル定義

2. **Phase 3: フロントエンド開発**
   - Next.js アプリケーション開発
   - FSDアーキテクチャ実装
   - Auth0統合

3. **Phase 4: 統合・テスト**
   - E2Eテスト実装
   - パフォーマンス最適化
   - セキュリティ検証強化

このインフラ実装計画により、スケーラブルで安全なFinSightアプリケーションの基盤が確立されます。