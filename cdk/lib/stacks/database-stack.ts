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
export interface DatabaseStackProps extends StackProps {
  vpc: Vpc;
  rdsSecurityGroup: SecurityGroup;
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
      instanceType: InstanceType.of(InstanceClass.T3, InstanceSize.SMALL),
      vpc: props.vpc,
      vpcSubnets: {
        subnetType: SubnetType.PRIVATE_ISOLATED,
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

    // データベース初期化はECSアプリケーション起動時に実行
    // Lambda関数は削除し、Goアプリケーションの初期化処理で対応
  }
}