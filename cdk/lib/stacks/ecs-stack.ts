import { Stack, StackProps, Duration, CfnOutput } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import {
  Vpc,
  SecurityGroup,
  Port,
  SubnetType,
  Peer,
} from 'aws-cdk-lib/aws-ec2';
import {
  Cluster,
  FargateService,
  FargateTaskDefinition,
  ContainerImage,
  Protocol,
  LogDrivers,
} from 'aws-cdk-lib/aws-ecs';
import {
  ApplicationLoadBalancer,
  ApplicationTargetGroup,
  TargetType,
  ApplicationProtocol,
  ListenerAction,
} from 'aws-cdk-lib/aws-elasticloadbalancingv2';
import { LogGroup, RetentionDays } from 'aws-cdk-lib/aws-logs';
import { Secret } from 'aws-cdk-lib/aws-secretsmanager';
import { DatabaseInstance } from 'aws-cdk-lib/aws-rds';
import { Repository } from 'aws-cdk-lib/aws-ecr';
import { EnvironmentConfig } from '../interfaces/config';

export interface EcsStackProps extends StackProps {
  vpc: Vpc;
  database: DatabaseInstance;
  dbSecret: Secret;
  config: EnvironmentConfig;
}

export class EcsStack extends Stack {
  public readonly cluster: Cluster;
  public readonly service: FargateService;
  public readonly loadBalancer: ApplicationLoadBalancer;
  public readonly ecrRepository: Repository;

  constructor(scope: Construct, id: string, props: EcsStackProps) {
    super(scope, id, props);

    // ECRリポジトリ作成
    this.ecrRepository = new Repository(this, 'FinanceTrackerRepository', {
      repositoryName: `finance-tracker-${props.config.environment}`,
      imageScanOnPush: true,
    });

    // ECSクラスター作成
    this.cluster = new Cluster(this, 'FinanceTrackerCluster', {
      vpc: props.vpc,
      clusterName: `finance-tracker-${props.config.environment}`,
      containerInsights: true,
    });

    // セキュリティグループ作成
    const ecsSecurityGroup = new SecurityGroup(this, 'EcsSecurityGroup', {
      vpc: props.vpc,
      description: 'Security group for ECS tasks',
      allowAllOutbound: true,
    });

    const albSecurityGroup = new SecurityGroup(this, 'AlbSecurityGroup', {
      vpc: props.vpc,
      description: 'Security group for Application Load Balancer',
      allowAllOutbound: false,
    });

    // ALB → ECS通信許可
    ecsSecurityGroup.addIngressRule(
      albSecurityGroup,
      Port.tcp(8080),
      'Allow ALB to access ECS'
    );

    // インターネット → ALB通信許可
    albSecurityGroup.addIngressRule(
      Peer.anyIpv4(),
      Port.tcp(80),
      'Allow HTTP traffic'
    );

    albSecurityGroup.addIngressRule(
      Peer.anyIpv4(),
      Port.tcp(443),
      'Allow HTTPS traffic'
    );

    // RDS接続許可（既存のRDSセキュリティグループに追加）
    const rdsSecurityGroup = SecurityGroup.fromSecurityGroupId(
      this,
      'ImportedRdsSecurityGroup',
      props.database.connections.securityGroups[0].securityGroupId
    );

    rdsSecurityGroup.addIngressRule(
      ecsSecurityGroup,
      Port.tcp(5432),
      'Allow ECS to access RDS'
    );

    // Application Load Balancer
    this.loadBalancer = new ApplicationLoadBalancer(this, 'FinanceTrackerALB', {
      vpc: props.vpc,
      internetFacing: true,
      loadBalancerName: `finance-tracker-alb-${props.config.environment}`,
      securityGroup: albSecurityGroup,
      vpcSubnets: {
        subnetType: SubnetType.PUBLIC,
      },
    });

    // ログ設定
    const logGroup = new LogGroup(this, 'FinanceTrackerLogGroup', {
      logGroupName: `/ecs/finance-tracker-${props.config.environment}`,
      retention: RetentionDays.ONE_WEEK,
    });

    // Fargateタスク定義
    const taskDefinition = new FargateTaskDefinition(this, 'FinanceTrackerTaskDef', {
      memoryLimitMiB: props.config.ecsConfig.memoryMB,
      cpu: props.config.ecsConfig.cpu,
    });

    // データベースシークレットへの読み取り権限付与
    props.dbSecret.grantRead(taskDefinition.taskRole);

    // コンテナ定義
    const container = taskDefinition.addContainer('FinanceTrackerContainer', {
      image: ContainerImage.fromRegistry('public.ecr.aws/docker/library/golang:1.21-alpine'),
      logging: LogDrivers.awsLogs({
        streamPrefix: 'finance-tracker',
        logGroup,
      }),
      environment: {
        ENVIRONMENT: props.config.environment,
        DB_ENDPOINT: props.database.instanceEndpoint.hostname,
        DB_SECRET_ARN: props.dbSecret.secretArn,
        PORT: '8080',
      },
      // 一時的なコマンド（後でGoアプリケーションに変更）
      command: [
        'sh', '-c',
        'echo "Finance Tracker API starting..." && go run -c "package main; import(\\"fmt\\"; \\"net/http\\"); func main() { http.HandleFunc(\\"/health\\", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, \\"{\\\\\\"status\\\\\\":\\\\\\"healthy\\\\\\",\\\\\\"service\\\\\\":\\\\\\"Finance Tracker Go API\\\\\\"}\\") }); fmt.Println(\\"Server starting on :8080\\"); http.ListenAndServe(\\":8080\\", nil) }" /tmp/main.go'
      ],
    });

    container.addPortMappings({
      containerPort: 8080,
      protocol: Protocol.TCP,
    });

    // ターゲットグループ
    const targetGroup = new ApplicationTargetGroup(this, 'FinanceTrackerTargetGroup', {
      vpc: props.vpc,
      port: 8080,
      protocol: ApplicationProtocol.HTTP,
      targetType: TargetType.IP,
      healthCheck: {
        path: '/health',
        interval: Duration.seconds(30),
        timeout: Duration.seconds(5),
        healthyThresholdCount: 2,
        unhealthyThresholdCount: 3,
      },
    });

    // ALBリスナー
    const listener = this.loadBalancer.addListener('FinanceTrackerListener', {
      port: 80,
      defaultAction: ListenerAction.forward([targetGroup]),
    });

    // Fargateサービス
    this.service = new FargateService(this, 'FinanceTrackerService', {
      cluster: this.cluster,
      taskDefinition,
      serviceName: `finance-tracker-${props.config.environment}`,
      desiredCount: props.config.ecsConfig.desiredCount,
      securityGroups: [ecsSecurityGroup],
      vpcSubnets: {
        subnetType: SubnetType.PUBLIC,
      },
      enableExecuteCommand: true,
    });

    // ターゲットグループにサービスを追加
    this.service.attachToApplicationTargetGroup(targetGroup);

    // 出力
    new CfnOutput(this, 'LoadBalancerDNS', {
      value: this.loadBalancer.loadBalancerDnsName,
      description: 'DNS name of the load balancer',
    });

    new CfnOutput(this, 'ClusterName', {
      value: this.cluster.clusterName,
      description: 'Name of the ECS cluster',
    });

    new CfnOutput(this, 'ServiceName', {
      value: this.service.serviceName,
      description: 'Name of the ECS service',
    });
  }
}