import { Stack, StackProps, CfnOutput } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { CfnWebACL, CfnWebACLAssociation } from 'aws-cdk-lib/aws-wafv2';
import { RestApi } from 'aws-cdk-lib/aws-apigateway';
import { EnvironmentConfig } from '../interfaces/config';

export interface SecurityStackProps extends StackProps {
  config: EnvironmentConfig;
  api: RestApi;
}

export class SecurityStack extends Stack {
  public readonly webAcl: CfnWebACL;

  constructor(scope: Construct, id: string, props: SecurityStackProps) {
    super(scope, id, props);

    // WAF Web ACL
    this.webAcl = new CfnWebACL(this, 'ApiWebAcl', {
      name: `finsight-api-waf-${props.config.environment}`,
      scope: 'REGIONAL',
      defaultAction: { allow: {} },
      rules: [
        // レート制限ルール
        {
          name: 'RateLimitRule',
          priority: 1,
          statement: {
            rateBasedStatement: {
              limit: 2000, // Production rate limit
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
        // GeoBlockルール（必要に応じて有効化）
        {
          name: 'GeoBlockRule',
          priority: 2,
          statement: {
            notStatement: {
              statement: {
                geoMatchStatement: {
                  countryCodes: ['JP', 'US'], // 許可する国
                },
              },
            },
          },
          action: { block: {} }, // Always block in production
          visibilityConfig: {
            sampledRequestsEnabled: true,
            cloudWatchMetricsEnabled: true,
            metricName: 'GeoBlockRule',
          },
        },
        // SQLインジェクション対策
        {
          name: 'SQLInjectionRule',
          priority: 3,
          statement: {
            orStatement: {
              statements: [
                {
                  sqliMatchStatement: {
                    fieldToMatch: { body: {} },
                    textTransformations: [{
                      priority: 0,
                      type: 'URL_DECODE',
                    }, {
                      priority: 1,
                      type: 'HTML_ENTITY_DECODE',
                    }],
                  },
                },
                {
                  sqliMatchStatement: {
                    fieldToMatch: { queryString: {} },
                    textTransformations: [{
                      priority: 0,
                      type: 'URL_DECODE',
                    }, {
                      priority: 1,
                      type: 'HTML_ENTITY_DECODE',
                    }],
                  },
                },
              ],
            },
          },
          action: { block: {} },
          visibilityConfig: {
            sampledRequestsEnabled: true,
            cloudWatchMetricsEnabled: true,
            metricName: 'SQLInjectionRule',
          },
        },
        // XSS対策
        {
          name: 'XSSRule',
          priority: 4,
          statement: {
            orStatement: {
              statements: [
                {
                  xssMatchStatement: {
                    fieldToMatch: { body: {} },
                    textTransformations: [{
                      priority: 0,
                      type: 'URL_DECODE',
                    }, {
                      priority: 1,
                      type: 'HTML_ENTITY_DECODE',
                    }],
                  },
                },
                {
                  xssMatchStatement: {
                    fieldToMatch: { queryString: {} },
                    textTransformations: [{
                      priority: 0,
                      type: 'URL_DECODE',
                    }, {
                      priority: 1,
                      type: 'HTML_ENTITY_DECODE',
                    }],
                  },
                },
              ],
            },
          },
          action: { block: {} },
          visibilityConfig: {
            sampledRequestsEnabled: true,
            cloudWatchMetricsEnabled: true,
            metricName: 'XSSRule',
          },
        },
        // 大きなリクエストボディのブロック
        {
          name: 'LargeBodyRule',
          priority: 5,
          statement: {
            sizeConstraintStatement: {
              fieldToMatch: { body: {} },
              comparisonOperator: 'GT',
              size: 8192, // 8KB
              textTransformations: [{
                priority: 0,
                type: 'NONE',
              }],
            },
          },
          action: { block: {} },
          visibilityConfig: {
            sampledRequestsEnabled: true,
            cloudWatchMetricsEnabled: true,
            metricName: 'LargeBodyRule',
          },
        },
        // 既知の悪意あるボットをブロック
        {
          name: 'BadBotRule',
          priority: 6,
          statement: {
            orStatement: {
              statements: [
                {
                  byteMatchStatement: {
                    searchString: 'bot',
                    fieldToMatch: {
                      singleHeader: { name: 'user-agent' },
                    },
                    textTransformations: [{
                      priority: 0,
                      type: 'LOWERCASE',
                    }],
                    positionalConstraint: 'CONTAINS',
                  },
                },
                {
                  byteMatchStatement: {
                    searchString: 'crawler',
                    fieldToMatch: {
                      singleHeader: { name: 'user-agent' },
                    },
                    textTransformations: [{
                      priority: 0,
                      type: 'LOWERCASE',
                    }],
                    positionalConstraint: 'CONTAINS',
                  },
                },
              ],
            },
          },
          action: { block: {} }, // Always block in production
          visibilityConfig: {
            sampledRequestsEnabled: true,
            cloudWatchMetricsEnabled: true,
            metricName: 'BadBotRule',
          },
        },
      ],
      visibilityConfig: {
        sampledRequestsEnabled: true,
        cloudWatchMetricsEnabled: true,
        metricName: `finsight-waf-${props.config.environment}`,
      },
    });

    // WAFとAPI Gatewayの関連付け
    new CfnWebACLAssociation(this, 'WebAclAssociation', {
      resourceArn: `arn:aws:apigateway:${this.region}::/restapis/${props.api.restApiId}/stages/${props.api.deploymentStage.stageName}`,
      webAclArn: this.webAcl.attrArn,
    });

    // 出力
    new CfnOutput(this, 'WebAclId', {
      value: this.webAcl.ref,
      description: 'WAF Web ACL ID',
    });

    new CfnOutput(this, 'WebAclArn', {
      value: this.webAcl.attrArn,
      description: 'WAF Web ACL ARN',
    });
  }
}