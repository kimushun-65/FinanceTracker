import { Stack, StackProps, RemovalPolicy, Duration } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import {
  DatabaseInstance,
  DatabaseInstanceEngine,
  PostgresEngineVersion,
  Credentials,
} from 'aws-cdk-lib/aws-rds';
import { Secret } from 'aws-cdk-lib/aws-secretsmanager';
import { Vpc, SecurityGroup, SubnetType, InstanceType, InstanceClass, InstanceSize } from 'aws-cdk-lib/aws-ec2';
import { Function, Runtime, Code } from 'aws-cdk-lib/aws-lambda';

export interface DatabaseStackProps extends StackProps {
  vpc: Vpc;
  rdsSecurityGroup: SecurityGroup;
  lambdaSecurityGroup: SecurityGroup;
  environment: string;
}

export class DatabaseStack extends Stack {
  public readonly database: DatabaseInstance;
  public readonly dbSecret: Secret;

  constructor(scope: Construct, id: string, props: DatabaseStackProps) {
    super(scope, id, props);

    // データベース認証情報をSecrets Managerで管理
    this.dbSecret = new Secret(this, 'DatabaseSecret', {
      secretName: `finsight-db-credentials-${props.environment}`,
      generateSecretString: {
        secretStringTemplate: JSON.stringify({ username: 'postgres' }),
        generateStringKey: 'password',
        excludeCharacters: '"@/\\',
      },
    });

    // RDS インスタンス作成
    this.database = new DatabaseInstance(this, 'FinSightDatabase', {
      engine: DatabaseInstanceEngine.postgres({
        version: PostgresEngineVersion.VER_15,
      }),
      instanceType: props.environment === 'prod'
        ? InstanceType.of(InstanceClass.T3, InstanceSize.SMALL)
        : InstanceType.of(InstanceClass.T3, InstanceSize.MICRO),
      vpc: props.vpc,
      vpcSubnets: {
        subnetType: SubnetType.PRIVATE_WITH_EGRESS,
      },
      securityGroups: [props.rdsSecurityGroup],
      credentials: Credentials.fromSecret(this.dbSecret),
      multiAz: props.environment === 'prod',
      storageEncrypted: true,
      backupRetention: Duration.days(7),
      deletionProtection: props.environment === 'prod',
      databaseName: 'finsight',
      removalPolicy: props.environment === 'prod'
        ? RemovalPolicy.RETAIN
        : RemovalPolicy.DESTROY,
    });

    // データベース初期化Lambda
    const dbInitFunction = new Function(this, 'DbInitFunction', {
      runtime: Runtime.NODEJS_18_X,
      handler: 'index.handler',
      code: Code.fromAsset('lambda/db-init'),
      vpc: props.vpc,
      securityGroups: [props.lambdaSecurityGroup],
      environment: {
        DB_SECRET_ARN: this.dbSecret.secretArn,
        DB_ENDPOINT: this.database.instanceEndpoint.hostname,
      },
      timeout: Duration.minutes(5),
    });

    // Secret読み取り権限付与
    this.dbSecret.grantRead(dbInitFunction);
  }
}