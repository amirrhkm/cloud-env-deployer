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

	vpc := awsec2.NewVpc(stack, jsii.String("TempInstance/ProactiveMonitoringVPC"), &awsec2.VpcProps{
		MaxAzs:      jsii.Number(1),
		NatGateways: jsii.Number(0),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				CidrMask:   jsii.Number(24),
				Name:       jsii.String("PM-PublicSubnet"),
				SubnetType: awsec2.SubnetType_PUBLIC,
			},
		},
	})

	sg := awsec2.NewSecurityGroup(stack, jsii.String("TempInstance/ProactiveMonitoringSG"), &awsec2.SecurityGroupProps{
		Vpc:              vpc,
		Description:      jsii.String("ProactiveMonitoringSG"),
		AllowAllOutbound: jsii.Bool(true),
	})

	sg.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_AllTraffic(), jsii.String("Allow all inbound traffic"), nil)

	keyPair := awsec2.KeyPair_FromKeyPairName(stack, jsii.String("ExistingKeyPair"), jsii.String("amir"))

	// <--- Linux --->
	linux_userData := awsec2.UserData_ForLinux(
		&awsec2.LinuxUserDataOptions{
			Shebang: jsii.String("#!/bin/bash"),
		},
	)

	// <--- Linux User Data --->
	linux_userData.AddCommands(
		jsii.String("sudo apt update"),
		jsii.String("sudo apt install -y docker.io docker-compose awscli rsyslog"),
	)

	// <--- Debian 12 AMI - EC2 Instance --->
	awsec2.NewInstance(stack, jsii.String("TempInstance/RsyslogSandbox"), &awsec2.InstanceProps{
		InstanceType: awsec2.NewInstanceType(jsii.String("t4g.small")),
		MachineImage: awsec2.MachineImage_Lookup(&awsec2.LookupMachineImageProps{
			Name:   jsii.String("debian-12-arm64-*"),
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
		UserData:                 linux_userData,
	})

	// <--- Ubuntu 22.04 AMI - EC2 Instance --->
	awsec2.NewInstance(stack, jsii.String("RsyslogSandbox"), &awsec2.InstanceProps{
		InstanceType: awsec2.NewInstanceType(jsii.String("t3.medium")),
		MachineImage: awsec2.MachineImage_Lookup(&awsec2.LookupMachineImageProps{
			Name:   jsii.String("ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"),
			Owners: jsii.Strings("099720109477"),
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
		UserData:                 linux_userData,
	})

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

	MonitoringHubStack(app, "TempInstance/AmirStack", &MonitoringHubProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// monitoring: 715841329405
// ronpos-staging: 530830676072
func env() *awscdk.Environment {

	return &awscdk.Environment{
		Account: jsii.String("715841329405"),
		Region:  jsii.String("ap-southeast-1"),
	}
}
