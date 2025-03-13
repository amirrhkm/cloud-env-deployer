package instances

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func CreateRHELInstance(scope constructs.Construct, vpc awsec2.IVpc, sg awsec2.ISecurityGroup, keyPair awsec2.IKeyPair) {
	rhel_userData := awsec2.UserData_ForLinux(
		&awsec2.LinuxUserDataOptions{
			Shebang: jsii.String("#!/bin/bash"),
		},
	)

	rhel_userData.AddCommands(
		jsii.String("sudo dnf update -y"),
		jsii.String("sudo dnf install -y docker"),
		jsii.String("sudo systemctl enable --now docker"),
		jsii.String("sudo dnf install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm -y"),
		jsii.String("sudo dnf install -y neofetch"),
		jsii.String("sudo yum install -y unzip"),
		jsii.String("curl \"https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip\" -o \"awscliv2.zip\""),
		jsii.String("unzip awscliv2.zip"),
		jsii.String("sudo ./aws/install"),
		jsii.String("sudo yum install -y yum-utils"),
		jsii.String("sudo yum-config-manager --add-repo https://download.docker.com/linux/rhel/docker-ce.repo"),
		jsii.String("sudo yum update -y"),
		jsii.String("sudo yum install -y docker-compose-plugin"),
		jsii.String("sudo dnf install -y dpkg"),
	)

	awsec2.NewInstance(scope, jsii.String("amir/MonitoringHubRHELInstance"), &awsec2.InstanceProps{
		InstanceType: awsec2.NewInstanceType(jsii.String("t4g.small")),
		MachineImage: awsec2.MachineImage_Lookup(&awsec2.LookupMachineImageProps{
			Name:   jsii.String("RHEL-9*"),
			Owners: jsii.Strings("309956199498"),
			Filters: &map[string]*[]*string{
				"architecture":        jsii.Strings("arm64"),
				"root-device-type":    jsii.Strings("ebs"),
				"virtualization-type": jsii.Strings("hvm"),
			},
		}),
		Vpc:                      vpc,
		VpcSubnets:               &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC},
		KeyPair:                  keyPair,
		AssociatePublicIpAddress: jsii.Bool(true),
		SecurityGroup:            sg,
		UserData:                 rhel_userData,
	})
}
