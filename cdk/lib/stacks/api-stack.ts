import { Stack, StackProps, Duration, CfnOutput } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { ServicePrincipal } from 'aws-cdk-lib/aws-iam';
import {
  RestApi,
  CfnAuthorizer,
  Cors,
  LambdaIntegration,
  MethodOptions,
  AuthorizationType,
} from 'aws-cdk-lib/aws-apigateway';
import { Function, Runtime, Code, Tracing } from 'aws-cdk-lib/aws-lambda';
import { Vpc, SecurityGroup } from 'aws-cdk-lib/aws-ec2';
import { Secret } from 'aws-cdk-lib/aws-secretsmanager';
import { DatabaseInstance } from 'aws-cdk-lib/aws-rds';
import { EnvironmentConfig } from '../interfaces/config';

export interface ApiStackProps extends StackProps {
  vpc: Vpc;
  lambdaSecurityGroup: SecurityGroup;
  database: DatabaseInstance;
  dbSecret: Secret;
  config: EnvironmentConfig;
}

export class ApiStack extends Stack {
  public readonly api: RestApi;
  public readonly authorizer: CfnAuthorizer;
  public readonly lambdaFunctions: { [key: string]: Function } = {};

  constructor(scope: Construct, id: string, props: ApiStackProps) {
    super(scope, id, props);

    // REST API作成
    this.api = new RestApi(this, 'FinSightApi', {
      restApiName: `finsight-api-${props.config.environment}`,
      description: `FinSight REST API for ${props.config.environment}`,
      defaultCorsPreflightOptions: {
        allowOrigins: Cors.ALL_ORIGINS,
        allowMethods: Cors.ALL_METHODS,
        allowHeaders: [
          'Content-Type',
          'Authorization',
          'X-Api-Key',
          'X-Amz-Date',
          'X-Amz-Security-Token',
        ],
        allowCredentials: true,
      },
    });

    // Lambda Authorizer関数
    const authorizerFunction = new Function(this, 'AuthorizerFunction', {
      runtime: Runtime.NODEJS_18_X,
      handler: 'index.handler',
      code: Code.fromAsset('lambda/authorizer'),
      environment: {
        AUTH0_DOMAIN: props.config.auth0Domain,
        AUTH0_AUDIENCE: props.config.auth0Audience,
      },
      timeout: Duration.seconds(10),
      memorySize: 256,
      tracing: Tracing.ACTIVE,
    });

    // API Gateway Authorizer
    this.authorizer = new CfnAuthorizer(this, 'ApiAuthorizer', {
      name: `finsight-authorizer-${props.config.environment}`,
      type: 'TOKEN',
      identitySource: 'method.request.header.Authorization',
      authorizerUri: `arn:aws:apigateway:${this.region}:lambda:path/2015-03-31/functions/${authorizerFunction.functionArn}/invocations`,
      authorizerResultTtlInSeconds: 300,
      restApiId: this.api.restApiId,
    });

    // Lambda invoke permission for authorizer
    authorizerFunction.addPermission('ApiGatewayAuthorizerPermission', {
      principal: new ServicePrincipal('apigateway.amazonaws.com'),
      sourceArn: `arn:aws:execute-api:${this.region}:${this.account}:${this.api.restApiId}/*`,
    });

    // Lambda関数の共通環境変数
    const commonEnv = {
      DB_SECRET_ARN: props.dbSecret.secretArn,
      DB_ENDPOINT: props.database.instanceEndpoint.hostname,
      ENVIRONMENT: props.config.environment,
    };

    // Lambda関数定義
    const lambdaConfigs = [
      { name: 'users', path: 'users' },
      { name: 'accounts', path: 'accounts' },
      { name: 'transactions', path: 'transactions' },
      { name: 'categories', path: 'categories' },
      { name: 'budgets', path: 'budgets' },
      { name: 'reports', path: 'reports' },
      { name: 'auth', path: 'auth' },
      { name: 'notifications', path: 'notifications' },
    ];

    // Lambda関数作成
    lambdaConfigs.forEach(config => {
      const lambdaFunction = new Function(this, `${config.name}Function`, {
        runtime: Runtime.PROVIDED_AL2,
        handler: 'bootstrap',
        code: Code.fromAsset(`lambda/${config.path}`),
        vpc: props.vpc,
        securityGroups: [props.lambdaSecurityGroup],
        environment: commonEnv,
        timeout: Duration.seconds(props.config.lambdaConfig.timeout),
        memorySize: props.config.lambdaConfig.memorySize,
        tracing: Tracing.ACTIVE,
      });

      // データベースシークレットへの読み取り権限付与
      props.dbSecret.grantRead(lambdaFunction);

      // 関数を保存
      this.lambdaFunctions[config.name] = lambdaFunction;
    });

    // APIエンドポイント設定
    this.setupApiEndpoints(props);

    // API出力
    new CfnOutput(this, 'ApiEndpoint', {
      value: this.api.url,
      description: 'API Gateway endpoint URL',
    });

    new CfnOutput(this, 'ApiId', {
      value: this.api.restApiId,
      description: 'API Gateway ID',
    });
  }

  private setupApiEndpoints(props: ApiStackProps) {
    // 認証が必要なメソッドオプション
    const authMethodOptions: MethodOptions = {
      authorizationType: AuthorizationType.CUSTOM,
      authorizer: {
        authorizerId: this.authorizer.ref,
      },
    };

    // APIバージョン
    const v1 = this.api.root.addResource('v1');

    // Users endpoints
    const users = v1.addResource('users');
    const userMe = users.addResource('me');
    userMe.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.users), authMethodOptions);
    userMe.addMethod('PUT', new LambdaIntegration(this.lambdaFunctions.users), authMethodOptions);

    // Accounts endpoints
    const accounts = v1.addResource('accounts');
    accounts.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.accounts), authMethodOptions);
    accounts.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.accounts), authMethodOptions);
    const accountId = accounts.addResource('{id}');
    accountId.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.accounts), authMethodOptions);
    accountId.addMethod('PUT', new LambdaIntegration(this.lambdaFunctions.accounts), authMethodOptions);
    accountId.addMethod('DELETE', new LambdaIntegration(this.lambdaFunctions.accounts), authMethodOptions);

    // Transactions endpoints
    const transactions = v1.addResource('transactions');
    transactions.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.transactions), authMethodOptions);
    transactions.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.transactions), authMethodOptions);
    const transactionId = transactions.addResource('{id}');
    transactionId.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.transactions), authMethodOptions);
    transactionId.addMethod('PUT', new LambdaIntegration(this.lambdaFunctions.transactions), authMethodOptions);
    transactionId.addMethod('DELETE', new LambdaIntegration(this.lambdaFunctions.transactions), authMethodOptions);

    // Categories endpoints
    const categories = v1.addResource('categories');
    categories.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.categories), authMethodOptions);
    categories.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.categories), authMethodOptions);
    const categoryId = categories.addResource('{id}');
    categoryId.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.categories), authMethodOptions);
    categoryId.addMethod('PUT', new LambdaIntegration(this.lambdaFunctions.categories), authMethodOptions);
    categoryId.addMethod('DELETE', new LambdaIntegration(this.lambdaFunctions.categories), authMethodOptions);

    // Budgets endpoints
    const budgets = v1.addResource('budgets');
    budgets.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.budgets), authMethodOptions);
    budgets.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.budgets), authMethodOptions);
    const budgetId = budgets.addResource('{id}');
    budgetId.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.budgets), authMethodOptions);
    budgetId.addMethod('PUT', new LambdaIntegration(this.lambdaFunctions.budgets), authMethodOptions);
    budgetId.addMethod('DELETE', new LambdaIntegration(this.lambdaFunctions.budgets), authMethodOptions);

    // Reports endpoints
    const reports = v1.addResource('reports');
    reports.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.reports), authMethodOptions);
    const reportsSummary = reports.addResource('summary');
    reportsSummary.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.reports), authMethodOptions);
    const reportsBudget = reports.addResource('budget');
    reportsBudget.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.reports), authMethodOptions);

    // Auth endpoints (認証不要)
    const auth = v1.addResource('auth');
    const authCallback = auth.addResource('callback');
    authCallback.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.auth));
    const authRefresh = auth.addResource('refresh');
    authRefresh.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.auth));
    const authLogout = auth.addResource('logout');
    authLogout.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.auth), authMethodOptions);

    // Notifications endpoints
    const notifications = v1.addResource('notifications');
    notifications.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.notifications), authMethodOptions);
    notifications.addMethod('POST', new LambdaIntegration(this.lambdaFunctions.notifications), authMethodOptions);
    const notificationPreferences = notifications.addResource('preferences');
    notificationPreferences.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.notifications), authMethodOptions);
    notificationPreferences.addMethod('PUT', new LambdaIntegration(this.lambdaFunctions.notifications), authMethodOptions);

    // Health check endpoint (認証不要)
    const health = v1.addResource('health');
    health.addMethod('GET', new LambdaIntegration(this.lambdaFunctions.auth));
  }
}