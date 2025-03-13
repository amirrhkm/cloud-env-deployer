package main

import (
	"provisioning-cdk/instances"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ProvisioningProps struct {
	awscdk.StackProps
}

func ProvisioningStack(scope constructs.Construct, id string, props *ProvisioningProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	silentmode_owner := "silentmode"
	silentmode_environment := "staging"
	silentmode_service := "proactive-monitoring"

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
		jsii.String("sudo apt install -y docker.io docker-compose awscli rsyslog ufw neofetch make"),
	)

	// Create instances
	instances.CreateDebianInstance(stack, vpc, sg, keyPair, linux_userData)
	instances.CreateUbuntuInstance(stack, vpc, sg, keyPair, linux_userData)
	instances.CreateAL2Instance(stack, vpc, sg, keyPair, linux_userData)
	instances.CreateRHELInstance(stack, vpc, sg, keyPair)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	ProvisioningStack(app, "TempInstance/AmirStack", &ProvisioningProps{
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
