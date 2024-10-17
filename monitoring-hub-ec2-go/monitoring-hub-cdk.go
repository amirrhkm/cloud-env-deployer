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

	userData := awsec2.UserData_ForLinux(
		&awsec2.LinuxUserDataOptions{
			Shebang: jsii.String("#!/bin/bash"),
		},
	)

	awsec2.NewInstance(stack, jsii.String("amir/MonitoringHubInstance"), &awsec2.InstanceProps{
		InstanceType: awsec2.NewInstanceType(jsii.String("t4g.nano")),
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
		UserData:                 userData,
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
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	// 	Account: jsii.String("530830676072"),
	// 	Region:  jsii.String("ap-southeast-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
