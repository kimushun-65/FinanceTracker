import { Stack, StackProps, CfnOutput } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { ConfigurationSet } from 'aws-cdk-lib/aws-ses';
import { Topic } from 'aws-cdk-lib/aws-sns';
import { EmailSubscription } from 'aws-cdk-lib/aws-sns-subscriptions';
import { Function } from 'aws-cdk-lib/aws-lambda';
import { PolicyStatement, Effect } from 'aws-cdk-lib/aws-iam';
import { EnvironmentConfig } from '../interfaces/config';

export interface SesStackProps extends StackProps {
  config: EnvironmentConfig;
  lambdaFunctions: { [key: string]: Function };
}

export class SesStack extends Stack {
  public readonly configurationSet: ConfigurationSet;

  constructor(scope: Construct, id: string, props: SesStackProps) {
    super(scope, id, props);

    // メール送信用のドメインIDを作成
    const emailDomain = props.config.customDomain || 'finsight.local';

    // Configuration Set作成（メール送信の設定とイベント追跡）
    this.configurationSet = new ConfigurationSet(this, 'ConfigSet', {
      configurationSetName: `finsight-ses-config-${props.config.environment}`,
    });

    // バウンスと苦情の通知用SNSトピック
    const bounceTopic = new Topic(this, 'BounceTopic', {
      topicName: `finsight-bounce-${props.config.environment}`,
      displayName: 'FinSight Email Bounce Notifications',
    });

    const complaintTopic = new Topic(this, 'ComplaintTopic', {
      topicName: `finsight-complaint-${props.config.environment}`,
      displayName: 'FinSight Email Complaint Notifications',
    });

    // 開発環境では通知メールを送信
    if (props.config.environment !== 'prod') {
      bounceTopic.addSubscription(new EmailSubscription('dev-bounce@finsight.local'));
      complaintTopic.addSubscription(new EmailSubscription('dev-complaint@finsight.local'));
    }

    // イベント宛先の設定（基本的なSNS通知のみ）
    // CloudWatchとEventDestinationは後で設定

    // Lambda関数にSES送信権限を付与
    const sesSendPolicy = new PolicyStatement({
      effect: Effect.ALLOW,
      actions: ['ses:SendEmail', 'ses:SendRawEmail'],
      resources: ['*'],
      conditions: {
        StringEquals: {
          'ses:FromAddress': props.config.sesConfig.fromEmail,
        },
      },
    });

    // reports と notifications Lambda関数にSES権限を付与
    if (props.lambdaFunctions.reports) {
      props.lambdaFunctions.reports.addToRolePolicy(sesSendPolicy);
    }

    if (props.lambdaFunctions.notifications) {
      props.lambdaFunctions.notifications.addToRolePolicy(sesSendPolicy);
    }

    // 出力
    new CfnOutput(this, 'BounceTopicArn', {
      value: bounceTopic.topicArn,
      description: 'SES Bounce Topic ARN',
    });

    new CfnOutput(this, 'ConfigurationSetName', {
      value: this.configurationSet.configurationSetName,
      description: 'SES Configuration Set Name',
    });

    new CfnOutput(this, 'FromEmailAddress', {
      value: props.config.sesConfig.fromEmail,
      description: 'From Email Address for sending',
    });
  }
}