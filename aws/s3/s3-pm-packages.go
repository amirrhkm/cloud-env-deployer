package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ProvisioningStackProps struct {
	awscdk.StackProps
}

func ProvisioningStack(scope constructs.Construct, id string, props *ProvisioningStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	bucket := awss3.NewBucket(stack, jsii.String("ProvisioningBucket"), &awss3.BucketProps{
		BucketName:        jsii.String("provisioning"),
		Versioned:         jsii.Bool(true),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects: jsii.Bool(true),
	})

	awscdk.NewCfnOutput(stack, jsii.String("BucketName"), &awscdk.CfnOutputProps{
		Value:       bucket.BucketName(),
		Description: jsii.String("Provisioning Bucket"),
		ExportName:  jsii.String("ProvisioningBucket"),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	ProvisioningStack(app, "ProvisioningStack", &ProvisioningStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("715841329405"),
		Region:  jsii.String("ap-southeast-1"),
	}
}
