import { Stack, StackProps, CfnOutput, Duration } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Dashboard, GraphWidget, Metric, SingleValueWidget, TextWidget } from 'aws-cdk-lib/aws-cloudwatch';
import { Function } from 'aws-cdk-lib/aws-lambda';
import { DatabaseInstance } from 'aws-cdk-lib/aws-rds';
import { RestApi } from 'aws-cdk-lib/aws-apigateway';
import { Alarm, ComparisonOperator, TreatMissingData } from 'aws-cdk-lib/aws-cloudwatch';
import { SnsAction } from 'aws-cdk-lib/aws-cloudwatch-actions';
import { Topic } from 'aws-cdk-lib/aws-sns';
import { EmailSubscription } from 'aws-cdk-lib/aws-sns-subscriptions';
import { EnvironmentConfig } from '../interfaces/config';

export interface MonitoringStackProps extends StackProps {
  config: EnvironmentConfig;
  lambdaFunctions: { [key: string]: Function };
  database: DatabaseInstance;
  api: RestApi;
}

export class MonitoringStack extends Stack {
  public readonly dashboard: Dashboard;
  public readonly alertTopic: Topic;

  constructor(scope: Construct, id: string, props: MonitoringStackProps) {
    super(scope, id, props);

    // アラート通知用SNSトピック
    this.alertTopic = new Topic(this, 'AlertTopic', {
      topicName: `finsight-alerts-${props.config.environment}`,
      displayName: 'FinSight System Alerts',
    });

    // 開発環境ではメール通知を設定
    if (props.config.environment !== 'prod') {
      this.alertTopic.addSubscription(new EmailSubscription('dev-alerts@finsight.local'));
    }

    // Lambda関数のメトリクス
    const lambdaWidgets = Object.entries(props.lambdaFunctions).map(([name, func]) => {
      // Lambda関数のアラーム設定
      const errorAlarm = new Alarm(this, `${name}ErrorAlarm`, {
        metric: func.metricErrors({
          period: Duration.minutes(5),
        }),
        threshold: 5,
        evaluationPeriods: 2,
        comparisonOperator: ComparisonOperator.GREATER_THAN_OR_EQUAL_TO_THRESHOLD,
        treatMissingData: TreatMissingData.NOT_BREACHING,
        alarmDescription: `High error rate for ${name} function`,
      });
      errorAlarm.addAlarmAction(new SnsAction(this.alertTopic));

      const durationAlarm = new Alarm(this, `${name}DurationAlarm`, {
        metric: func.metricDuration({
          period: Duration.minutes(5),
        }),
        threshold: 10000, // 10秒
        evaluationPeriods: 3,
        comparisonOperator: ComparisonOperator.GREATER_THAN_THRESHOLD,
        treatMissingData: TreatMissingData.NOT_BREACHING,
        alarmDescription: `High duration for ${name} function`,
      });
      durationAlarm.addAlarmAction(new SnsAction(this.alertTopic));

      return new GraphWidget({
        title: `Lambda: ${name}`,
        width: 12,
        height: 6,
        left: [
          func.metricInvocations({ period: Duration.minutes(5) }),
          func.metricErrors({ period: Duration.minutes(5) }),
        ],
        right: [
          func.metricDuration({ period: Duration.minutes(5) }),
        ],
      });
    });

    // API Gatewayメトリクス
    const apiWidget = new GraphWidget({
      title: 'API Gateway Metrics',
      width: 12,
      height: 6,
      left: [
        props.api.metricCount({ period: Duration.minutes(5) }),
        props.api.metricClientError({ period: Duration.minutes(5) }),
        props.api.metricServerError({ period: Duration.minutes(5) }),
      ],
      right: [
        props.api.metricLatency({ period: Duration.minutes(5) }),
      ],
    });

    // API Gateway アラーム
    const apiErrorAlarm = new Alarm(this, 'ApiServerErrorAlarm', {
      metric: props.api.metricServerError({
        period: Duration.minutes(5),
      }),
      threshold: 10,
      evaluationPeriods: 2,
      comparisonOperator: ComparisonOperator.GREATER_THAN_THRESHOLD,
      treatMissingData: TreatMissingData.NOT_BREACHING,
      alarmDescription: 'High API Gateway server error rate',
    });
    apiErrorAlarm.addAlarmAction(new SnsAction(this.alertTopic));

    const apiLatencyAlarm = new Alarm(this, 'ApiLatencyAlarm', {
      metric: props.api.metricLatency({
        period: Duration.minutes(5),
      }),
      threshold: 5000, // 5秒
      evaluationPeriods: 3,
      comparisonOperator: ComparisonOperator.GREATER_THAN_THRESHOLD,
      treatMissingData: TreatMissingData.NOT_BREACHING,
      alarmDescription: 'High API Gateway latency',
    });
    apiLatencyAlarm.addAlarmAction(new SnsAction(this.alertTopic));

    // RDSメトリクス
    const rdsWidget = new GraphWidget({
      title: 'RDS Database Metrics',
      width: 12,
      height: 6,
      left: [
        props.database.metricCPUUtilization({ period: Duration.minutes(5) }),
        props.database.metricDatabaseConnections({ period: Duration.minutes(5) }),
      ],
      right: [
        props.database.metricFreeStorageSpace({ period: Duration.minutes(5) }),
      ],
    });

    // RDS アラーム
    const rdsConnectAlarm = new Alarm(this, 'RdsConnectionAlarm', {
      metric: props.database.metricDatabaseConnections({
        period: Duration.minutes(5),
      }),
      threshold: 80,
      evaluationPeriods: 2,
      comparisonOperator: ComparisonOperator.GREATER_THAN_THRESHOLD,
      treatMissingData: TreatMissingData.NOT_BREACHING,
      alarmDescription: 'High RDS connection count',
    });
    rdsConnectAlarm.addAlarmAction(new SnsAction(this.alertTopic));

    const rdsStorageAlarm = new Alarm(this, 'RdsStorageAlarm', {
      metric: props.database.metricFreeStorageSpace({
        period: Duration.minutes(5),
      }),
      threshold: 2000000000, // 2GB
      evaluationPeriods: 1,
      comparisonOperator: ComparisonOperator.LESS_THAN_THRESHOLD,
      treatMissingData: TreatMissingData.BREACHING,
      alarmDescription: 'Low RDS free storage space',
    });
    rdsStorageAlarm.addAlarmAction(new SnsAction(this.alertTopic));

    // システム概要ウィジェット
    const systemOverviewWidget = new TextWidget({
      markdown: `# FinSight System Overview - ${props.config.environment.toUpperCase()}

## Environment: ${props.config.environment}
## Region: ${props.config.region}

### Key Metrics to Monitor:
- **Lambda Functions**: Error rates, duration, invocation count
- **API Gateway**: Request count, latency, error rates  
- **RDS Database**: CPU, connections, storage, latency
- **System Health**: Overall availability and performance

### Alert Thresholds:
- Lambda errors: ≥5 errors in 10 minutes
- Lambda duration: >10 seconds for 15 minutes
- API errors: >10 server errors in 10 minutes  
- API latency: >5 seconds for 15 minutes
- RDS connections: >80 connections
- RDS storage: <2GB free space`,
      width: 24,
      height: 6,
    });

    // ダッシュボード作成
    this.dashboard = new Dashboard(this, 'FinSightDashboard', {
      dashboardName: `FinSight-${props.config.environment}`,
      widgets: [
        [systemOverviewWidget],
        [apiWidget],
        [rdsWidget],
        ...lambdaWidgets.map(widget => [widget]),
      ],
    });

    // ログ保持期間設定（将来の実装予定）
    // Lambda関数のログ保持期間を環境に応じて設定する場合はカスタムリソースが必要

    // 出力
    new CfnOutput(this, 'DashboardUrl', {
      value: `https://${props.config.region}.console.aws.amazon.com/cloudwatch/home?region=${props.config.region}#dashboards:name=${this.dashboard.dashboardName}`,
      description: 'CloudWatch Dashboard URL',
    });

    new CfnOutput(this, 'AlertTopicArn', {
      value: this.alertTopic.topicArn,
      description: 'SNS Topic ARN for alerts',
    });
  }
}