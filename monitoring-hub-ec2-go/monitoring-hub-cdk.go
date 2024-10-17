package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MonitoringHubProps struct {
	awscdk.StackProps
}

func MonitoringHubStack(scope constructs.Construct, id string, props *MonitoringHubProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	silentmode_owner := os.Getenv("SILENTMODE_OWNER")
	silentmode_environment := os.Getenv("SILENTMODE_ENVIRONMENT")
	silentmode_service := os.Getenv("SILENTMODE_SERVICE")

	awscdk.Tags_Of(stack).Add(jsii.String("silentmode:owner"), jsii.String(silentmode_owner), nil)
	awscdk.Tags_Of(stack).Add(jsii.String("silentmode:environment"), jsii.String(silentmode_environment), nil)
	awscdk.Tags_Of(stack).Add(jsii.String("silentmode:service"), jsii.String(silentmode_service), nil)

	vpc := awsec2.NewVpc(stack, jsii.String("amir/MonitoringHubVPC"), &awsec2.VpcProps{
		MaxAzs:      jsii.Number(1),
		NatGateways: jsii.Number(0),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				CidrMask:   jsii.Number(24),
				Name:       jsii.String("amir/MH-PublicSubnet"),
				SubnetType: awsec2.SubnetType_PUBLIC,
			},
		},
	})

	sg := awsec2.NewSecurityGroup(stack, jsii.String("amir/MonitoringHubSG"), &awsec2.SecurityGroupProps{
		Vpc:              vpc,
		Description:      jsii.String("amir/MonitoringHubSG"),
		AllowAllOutbound: jsii.Bool(true),
	})

	sg.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_AllTraffic(), jsii.String("Allow all inbound traffic"), nil)

	keyPair := awsec2.KeyPair_FromKeyPairName(stack, jsii.String("ExistingKeyPair"), jsii.String("amir"))

	// <--- RHEL 9 --->
	rhel_userData := awsec2.UserData_ForLinux(
		&awsec2.LinuxUserDataOptions{
			Shebang: jsii.String("#!/bin/bash"),
		},
	)

	// <--- RHEL 9 User Data --->
	rhel_userData.AddCommands(
		jsii.String("sudo dnf update -y"),
		jsii.String("sudo dnf install -y docker"),
		jsii.String("sudo systemctl enable --now docker"),
		jsii.String("sudo dnf install -y docker-compose"),
		jsii.String("sudo dnf install -y awscli"),
	)

	// <--- RHEL 9 AMI - EC2 Instance --->
	awsec2.NewInstance(stack, jsii.String("amir/MonitoringHubRHELInstance"), &awsec2.InstanceProps{
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

	// <--- Debian 12  --->
	debian_userData := awsec2.UserData_ForLinux(
		&awsec2.LinuxUserDataOptions{
			Shebang: jsii.String("#!/bin/bash"),
		},
	)

	// <--- Debian 12 User Data --->
	debian_userData.AddCommands(
		jsii.String("sudo apt update"),
		jsii.String("sudo apt install -y docker.io docker-compose awscli"),
	)

	// <--- Debian 12 AMI - EC2 Instance --->
	awsec2.NewInstance(stack, jsii.String("amir/MonitoringHubDebianInstance"), &awsec2.InstanceProps{
		InstanceType: awsec2.NewInstanceType(jsii.String("t4g.small")),
		MachineImage: awsec2.MachineImage_Lookup(&awsec2.LookupMachineImageProps{
			Name:   jsii.String("RHEL-9.0.0_HVM-20220513-arm64-0-Hourly2-GP2"),
			Owners: jsii.Strings("136693071363"),
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

	// <--- Windows --->
	windows_userData := awsec2.UserData_ForWindows(&awsec2.WindowsUserDataOptions{})

	// <--- Windows User Data --->
	windows_userData.AddCommands(
		jsii.String("powershell -Command \"[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12\""),
		jsii.String("powershell -Command \"Invoke-WebRequest -Uri https://awscli.amazonaws.com/AWSCLIV2.msi -OutFile C:\\AWSCLIV2.msi\""),
		jsii.String("msiexec.exe /i C:\\AWSCLIV2.msi /qn"),
		jsii.String("powershell -Command \"Invoke-WebRequest -Uri https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe -OutFile C:\\DockerDesktopInstaller.exe\""),
		jsii.String("C:\\DockerDesktopInstaller.exe install --quiet"),
	)

	// <--- Windows Server 2022 AMI - EC2 Instance --->
	awsec2.NewInstance(stack, jsii.String("amir/MonitoringHubWindowsInstance"), &awsec2.InstanceProps{
		InstanceType: awsec2.NewInstanceType(jsii.String("t3.medium")),
		MachineImage: awsec2.MachineImage_Lookup(&awsec2.LookupMachineImageProps{
			Name:   jsii.String("Windows_Server-2022-English-Full-Base-*"),
			Owners: jsii.Strings("801119661308"),
			Filters: &map[string]*[]*string{
				"architecture":        jsii.Strings("x86_64"),
				"root-device-type":    jsii.Strings("ebs"),
				"virtualization-type": jsii.Strings("hvm"),
			},
		}),
		Vpc:                      vpc,
		VpcSubnets:               &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC},
		KeyPair:                  keyPair,
		AssociatePublicIpAddress: jsii.Bool(true),
		SecurityGroup:            sg,
		UserData:                 windows_userData,
	})

	return stack
}
func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	MonitoringHubStack(app, "MonitoringHubCdkStack", &MonitoringHubProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {

	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("AWS_ACCOUNT")),
		Region:  jsii.String(os.Getenv("AWS_REGION")),
	}
}
