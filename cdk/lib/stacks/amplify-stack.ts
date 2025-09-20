import { Stack, StackProps, CfnOutput, SecretValue } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { App, GitHubSourceCodeProvider, Branch, RedirectStatus } from '@aws-cdk/aws-amplify-alpha';
import { EnvironmentConfig } from '../interfaces/config';

export interface AmplifyStackProps extends StackProps {
  config: EnvironmentConfig;
  apiEndpoint: string;
}

export class AmplifyStack extends Stack {
  public readonly amplifyApp: App;
  public readonly branch: Branch;

  constructor(scope: Construct, id: string, props: AmplifyStackProps) {
    super(scope, id, props);

    // Amplifyアプリケーション作成
    this.amplifyApp = new App(this, 'FinSightApp', {
      appName: `finsight-${props.config.environment}`,
      description: `FinSight Frontend for ${props.config.environment}`,
      sourceCodeProvider: new GitHubSourceCodeProvider({
        owner: props.config.githubOwner,
        repository: props.config.repositoryName,
        oauthToken: SecretValue.secretsManager('github-token'),
      }),
      environmentVariables: {
        REACT_APP_API_ENDPOINT: props.apiEndpoint,
        REACT_APP_AUTH0_DOMAIN: props.config.auth0Domain,
        REACT_APP_AUTH0_CLIENT_ID: props.config.auth0ClientId,
        REACT_APP_AUTH0_AUDIENCE: props.config.auth0Audience,
        REACT_APP_ENVIRONMENT: props.config.environment,
      },
      autoBranchCreation: {
        patterns: ['feature/*', 'bugfix/*'],
        environmentVariables: {
          REACT_APP_ENVIRONMENT: 'preview',
        },
      },
      autoBranchDeletion: true,
      customRules: [
        // SPAのルーティング対応
        {
          source: '</^[^.]+$/>',
          target: '/index.html',
          status: RedirectStatus.REWRITE,
        },
      ],
    });

    // ブランチ設定
    const branchName = props.config.environment === 'prod' ? 'main' : 
                     props.config.environment === 'staging' ? 'staging' : 'develop';
    
    this.branch = this.amplifyApp.addBranch(branchName, {
      branchName: branchName,
      stage: props.config.environment === 'prod' ? 'PRODUCTION' : 'DEVELOPMENT',
      autoBuild: true,
      performanceMode: props.config.environment === 'prod',
    });

    // カスタムドメイン設定（ドメインが設定されている場合）
    if (props.config.customDomain) {
      const domain = this.amplifyApp.addDomain(props.config.customDomain, {
        enableAutoSubdomain: true,
        autoSubdomainCreationPatterns: ['*'],
      });
      
      domain.mapRoot(this.branch);
      domain.mapSubDomain(this.branch, props.config.environment === 'prod' ? 'www' : props.config.environment);
    }

    // 出力
    new CfnOutput(this, 'AmplifyAppId', {
      value: this.amplifyApp.appId,
      description: 'Amplify App ID',
    });

    new CfnOutput(this, 'AmplifyAppUrl', {
      value: `https://${branchName}.${this.amplifyApp.defaultDomain}`,
      description: 'Amplify App URL',
    });

    if (props.config.customDomain) {
      new CfnOutput(this, 'CustomDomainUrl', {
        value: props.config.environment === 'prod' 
          ? `https://${props.config.customDomain}` 
          : `https://${props.config.environment}.${props.config.customDomain}`,
        description: 'Custom Domain URL',
      });
    }
  }
}