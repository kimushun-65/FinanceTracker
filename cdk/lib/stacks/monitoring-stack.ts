import { Stack, StackProps, CfnOutput, Duration } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { Dashboard, GraphWidget, Metric, SingleValueWidget, TextWidget } from 'aws-cdk-lib/aws-cloudwatch';
import { DatabaseInstance } from 'aws-cdk-lib/aws-rds';
import { RestApi } from 'aws-cdk-lib/aws-apigateway';
import { FargateService } from 'aws-cdk-lib/aws-ecs';
import { Alarm, ComparisonOperator, TreatMissingData } from 'aws-cdk-lib/aws-cloudwatch';
import { SnsAction } from 'aws-cdk-lib/aws-cloudwatch-actions';
import { Topic } from 'aws-cdk-lib/aws-sns';
import { EmailSubscription } from 'aws-cdk-lib/aws-sns-subscriptions';
import { EnvironmentConfig } from '../interfaces/config';

export interface MonitoringStackProps extends StackProps {
  config: EnvironmentConfig;
  ecsService: FargateService;
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

    // Production email notifications
    // NOTE: Update this email address to your actual alert email
    this.alertTopic.addSubscription(new EmailSubscription('alerts@finsight.com'));

    // ECSサービスのメトリクス
    const ecsServiceMetrics = [
      new Metric({
        namespace: 'AWS/ECS',
        metricName: 'CPUUtilization',
        dimensionsMap: {
          ServiceName: props.ecsService.serviceName,
          ClusterName: props.ecsService.cluster.clusterName,
        },
        period: Duration.minutes(5),
      }),
      new Metric({
        namespace: 'AWS/ECS',
        metricName: 'MemoryUtilization',
        dimensionsMap: {
          ServiceName: props.ecsService.serviceName,
          ClusterName: props.ecsService.cluster.clusterName,
        },
        period: Duration.minutes(5),
      }),
    ];

    // ECSサービスのアラーム
    const cpuAlarm = new Alarm(this, 'EcsCpuAlarm', {
      metric: ecsServiceMetrics[0],
      threshold: 80,
      evaluationPeriods: 2,
      comparisonOperator: ComparisonOperator.GREATER_THAN_THRESHOLD,
      treatMissingData: TreatMissingData.NOT_BREACHING,
      alarmDescription: 'High CPU utilization for ECS service',
    });
    cpuAlarm.addAlarmAction(new SnsAction(this.alertTopic));

    const memoryAlarm = new Alarm(this, 'EcsMemoryAlarm', {
      metric: ecsServiceMetrics[1],
      threshold: 80,
      evaluationPeriods: 2,
      comparisonOperator: ComparisonOperator.GREATER_THAN_THRESHOLD,
      treatMissingData: TreatMissingData.NOT_BREACHING,
      alarmDescription: 'High memory utilization for ECS service',
    });
    memoryAlarm.addAlarmAction(new SnsAction(this.alertTopic));

    const ecsWidget = new GraphWidget({
      title: 'ECS Service Metrics',
      width: 12,
      height: 6,
      left: [ecsServiceMetrics[0]],
      right: [ecsServiceMetrics[1]],
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
- ECS CPU: >80% for 10 minutes
- ECS Memory: >80% for 10 minutes
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
        [ecsWidget],
        [rdsWidget],
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