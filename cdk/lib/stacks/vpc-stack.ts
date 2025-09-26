import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import {
  Vpc,
  SubnetType,
  SecurityGroup,
  Port,
  IpAddresses,
} from 'aws-cdk-lib/aws-ec2';

export interface VpcStackProps extends StackProps {
  environment: string;
}

export class VpcStack extends Stack {
  public readonly vpc: Vpc;
  public readonly rdsSecurityGroup: SecurityGroup;

  constructor(scope: Construct, id: string, props: VpcStackProps) {
    super(scope, id, props);

    // VPC作成（NAT Gateway削除でコスト削減）
    this.vpc = new Vpc(this, 'FinSightVpc', {
      ipAddresses: IpAddresses.cidr('10.0.0.0/16'),
      maxAzs: 2,
      subnetConfiguration: [
        {
          cidrMask: 24,
          name: 'Public',
          subnetType: SubnetType.PUBLIC,
        },
        {
          cidrMask: 24,
          name: 'Database',
          subnetType: SubnetType.PRIVATE_ISOLATED,
        },
      ],
      natGateways: 0, // NAT Gateway削除（月9,000円削減）
    });

    // RDS用セキュリティグループ
    this.rdsSecurityGroup = new SecurityGroup(this, 'RdsSecurityGroup', {
      vpc: this.vpc,
      description: 'Security group for RDS database',
      allowAllOutbound: false,
    });
  }
}