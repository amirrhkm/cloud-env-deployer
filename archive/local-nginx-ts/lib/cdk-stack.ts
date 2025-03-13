import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';

export class CdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const vpc = new ec2.Vpc(this, 'MyVpc', {
      maxAzs: 1,
      subnetConfiguration: [
        {
          cidrMask: 24,
          name: 'PublicSubnet',
          subnetType: ec2.SubnetType.PUBLIC,
        },
        {
          cidrMask: 24,
          name: 'PrivateSubnet',
          subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
        },
      ],
    });

    const ami = ec2.MachineImage.latestAmazonLinux({
      generation: ec2.AmazonLinuxGeneration.AMAZON_LINUX_2,
    });

    const publicSecurityGroup = new ec2.SecurityGroup(this, 'PublicSecurityGroup', {
      vpc,
      allowAllOutbound: true,
    });
    publicSecurityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(22), 'Allow SSH access');
    publicSecurityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(80), 'Allow HTTP access');

    const publicInstance = new ec2.Instance(this, 'PublicEC2', {
      vpc,
      instanceType: ec2.InstanceType.of(ec2.InstanceClass.T2, ec2.InstanceSize.MICRO),
      machineImage: ami,
      vpcSubnets: { subnetType: ec2.SubnetType.PUBLIC },
      securityGroup: publicSecurityGroup,
      keyName: 'stitch-kp1',
    });

    publicInstance.addUserData(
      `#!/bin/bash`,
      `sudo amazon-linux-extras install nginx1.12 -y`,
      `sudo systemctl start nginx`,
      `sudo systemctl enable nginx`
    );

    const privateSecurityGroup = new ec2.SecurityGroup(this, 'PrivateSecurityGroup', {
      vpc,
      allowAllOutbound: true,
    });
    privateSecurityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(80), 'Allow HTTP access');

    const privateInstance = new ec2.Instance(this, 'PrivateEC2', {
      vpc,
      instanceType: ec2.InstanceType.of(ec2.InstanceClass.T2, ec2.InstanceSize.MICRO),
      machineImage: ami,
      vpcSubnets: { subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS },
      securityGroup: privateSecurityGroup,
      keyName: 'stitch-kp1',
    });

    privateInstance.addUserData(
      `#!/bin/bash`,
      `sudo amazon-linux-extras install nginx1.12 -y`,
      `sudo systemctl start nginx`,
      `sudo systemctl enable nginx`
    );
  }
}
