import { Stack, StackProps, CfnOutput } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { CfnApp, CfnBranch } from 'aws-cdk-lib/aws-amplify';
import { EnvironmentConfig } from '../interfaces/config';

export interface AmplifyStackProps extends StackProps {
  config: EnvironmentConfig;
  apiEndpoint: string;
}

export class AmplifyStackOAuth extends Stack {
  public readonly amplifyApp: CfnApp;
  public readonly branch: CfnBranch;

  constructor(scope: Construct, id: string, props: AmplifyStackProps) {
    super(scope, id, props);

    // L1 Constructを使用（より細かい制御が可能）
    this.amplifyApp = new CfnApp(this, 'FinSightApp', {
      name: `finsight-frontend-${props.config.environment}`,
      description: `FinSight Frontend for ${props.config.environment}`,
      repository: `https://github.com/${props.config.githubOwner}/${props.config.repositoryName}`,
      // OAuthトークンを直接指定
      oauthToken: '{{resolve:secretsmanager:github-token:SecretString}}',
      environmentVariables: [
        {
          name: 'REACT_APP_API_ENDPOINT',
          value: props.apiEndpoint,
        },
        {
          name: 'REACT_APP_AUTH0_DOMAIN',
          value: props.config.auth0Domain,
        },
        {
          name: 'REACT_APP_AUTH0_CLIENT_ID',
          value: props.config.auth0ClientId,
        },
        {
          name: 'REACT_APP_AUTH0_AUDIENCE',
          value: props.config.auth0Audience,
        },
        {
          name: 'REACT_APP_ENVIRONMENT',
          value: props.config.environment,
        },
      ],
      buildSpec: `version: 1
frontend:
  phases:
    preBuild:
      commands:
        - cd frontend
        - npm ci
    build:
      commands:
        - npm run build
  artifacts:
    baseDirectory: frontend/build
    files:
      - '**/*'
  cache:
    paths:
      - frontend/node_modules/**/*`,
      customRules: [
        {
          source: '</^[^.]+$/>',
          target: '/index.html',
          status: '200',
        },
      ],
    });

    // ブランチ設定
    this.branch = new CfnBranch(this, 'MainBranch', {
      appId: this.amplifyApp.attrAppId,
      branchName: 'main',
      stage: 'PRODUCTION',
      enableAutoBuild: true,
      enablePerformanceMode: true,
    });

    // 出力
    new CfnOutput(this, 'AmplifyAppId', {
      value: this.amplifyApp.attrAppId,
      description: 'Amplify App ID',
    });

    new CfnOutput(this, 'AmplifyAppUrl', {
      value: `https://main.${this.amplifyApp.attrDefaultDomain}`,
      description: 'Amplify App URL',
    });
  }
}