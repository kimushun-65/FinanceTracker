#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { VpcStack } from '../lib/stacks/vpc-stack';
import { DatabaseStack } from '../lib/stacks/database-stack';
import { ApiStack } from '../lib/stacks/api-stack';
// import { AmplifyStack } from '../lib/stacks/amplify-stack'; // 手動管理に移行
import { SesStack } from '../lib/stacks/ses-stack';
import { MonitoringStack } from '../lib/stacks/monitoring-stack';
import { SecurityStack } from '../lib/stacks/security-stack';
import { EnvironmentConfig } from '../lib/interfaces/config';

const app = new cdk.App();

// 環境を取得（デフォルトは'prod'）
const environment = app.node.tryGetContext('env') || 'prod';

// 環境設定を読み込み
const config: EnvironmentConfig = require(`../config/${environment}.json`);

// 共通の環境設定
const env = {
  account: process.env.CDK_DEFAULT_ACCOUNT,
  region: config.region,
};

// VPCスタック
const vpcStack = new VpcStack(app, `VpcStack-${environment}`, {
  environment,
  env,
});

// データベーススタック
const databaseStack = new DatabaseStack(app, `DatabaseStack-${environment}`, {
  vpc: vpcStack.vpc,
  rdsSecurityGroup: vpcStack.rdsSecurityGroup,
  lambdaSecurityGroup: vpcStack.lambdaSecurityGroup,
  environment,
  env,
});
databaseStack.addDependency(vpcStack);

// APIスタック
const apiStack = new ApiStack(app, `ApiStack-${environment}`, {
  vpc: vpcStack.vpc,
  lambdaSecurityGroup: vpcStack.lambdaSecurityGroup,
  database: databaseStack.database,
  dbSecret: databaseStack.dbSecret,
  config,
  env,
});
apiStack.addDependency(databaseStack);

// Amplifyスタック（フロントエンド）- 手動管理に移行
// const amplifyStack = new AmplifyStack(app, `AmplifyStack-${environment}`, {
//   config,
//   apiEndpoint: apiStack.api.url,
//   env,
// });
// amplifyStack.addDependency(apiStack);

// SESスタック（メール送信）
const sesStack = new SesStack(app, `SesStack-${environment}`, {
  config,
  lambdaFunctions: apiStack.lambdaFunctions,
  env,
});
sesStack.addDependency(apiStack);

// MonitoringStack（監視・ダッシュボード）
const monitoringStack = new MonitoringStack(app, `MonitoringStack-${environment}`, {
  config,
  lambdaFunctions: apiStack.lambdaFunctions,
  database: databaseStack.database,
  api: apiStack.api,
  env,
});
monitoringStack.addDependency(apiStack);

// SecurityStack（WAF設定）
const securityStack = new SecurityStack(app, `SecurityStack-${environment}`, {
  config,
  api: apiStack.api,
  env,
});
securityStack.addDependency(apiStack);

// タグを追加
cdk.Tags.of(app).add('Environment', environment);
cdk.Tags.of(app).add('Project', 'FinSight');