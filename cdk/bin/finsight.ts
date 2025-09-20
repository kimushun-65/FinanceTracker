#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { VpcStack } from '../lib/stacks/vpc-stack';
import { EnvironmentConfig } from '../lib/interfaces/config';

const app = new cdk.App();

// 環境を取得（デフォルトは'dev'）
const environment = app.node.tryGetContext('env') || 'dev';

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

// タグを追加
cdk.Tags.of(app).add('Environment', environment);
cdk.Tags.of(app).add('Project', 'FinSight');