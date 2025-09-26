import { Stack, StackProps, Duration, CfnOutput } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import {
  RestApi,
  Cors,
  Integration,
  IntegrationType,
  MethodOptions,
} from 'aws-cdk-lib/aws-apigateway';
import { Vpc } from 'aws-cdk-lib/aws-ec2';
import { ApplicationLoadBalancer } from 'aws-cdk-lib/aws-elasticloadbalancingv2';
import { EnvironmentConfig } from '../interfaces/config';

export interface ApiStackProps extends StackProps {
  vpc: Vpc;
  loadBalancer: ApplicationLoadBalancer;
  config: EnvironmentConfig;
}

export class ApiStack extends Stack {
  public readonly api: RestApi;

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

    // 認証はGoアプリケーション内で処理（Lambda-less アーキテクチャ）

    // バックエンド統合設定（ALBへの直接接続）
    const backendIntegration = new Integration({
      type: IntegrationType.HTTP_PROXY,
      integrationHttpMethod: 'ANY',
      uri: `http://${props.loadBalancer.loadBalancerDnsName}/{proxy}`,
      options: {
        requestParameters: {
          'integration.request.path.proxy': 'method.request.path.proxy',
          'integration.request.header.X-Forwarded-For': 'context.identity.sourceIp',
        },
      },
    });

    // すべてのエンドポイントはGoアプリケーションで認証処理
    const methodOptions: MethodOptions = {
      requestParameters: {
        'method.request.path.proxy': true,
      },
    };

    // APIエンドポイント設定
    this.setupApiEndpoints(backendIntegration, methodOptions);

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

  private setupApiEndpoints(
    backendIntegration: Integration,
    methodOptions: MethodOptions
  ) {
    // APIバージョン
    const v1 = this.api.root.addResource('v1');

    // プロキシ統合 - すべてのリクエストをGoアプリケーションに転送
    const proxy = v1.addResource('{proxy+}');
    proxy.addMethod('ANY', backendIntegration, methodOptions);
    
    // ルートパスも対応
    v1.addMethod('ANY', backendIntegration, methodOptions);
  }
}