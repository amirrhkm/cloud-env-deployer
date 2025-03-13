package instances

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func CreateAL2Instance(scope constructs.Construct, vpc awsec2.IVpc, sg awsec2.ISecurityGroup, keyPair awsec2.IKeyPair, userData awsec2.UserData) {
	awsec2.NewInstance(scope, jsii.String("RsyslogSandbox"), &awsec2.InstanceProps{
		InstanceType: awsec2.NewInstanceType(jsii.String("t3.small")),
		MachineImage: awsec2.MachineImage_Lookup(&awsec2.LookupMachineImageProps{
			Name:   jsii.String("amzn2-ami-hvm-*-x86_64-gp2"),
			Owners: jsii.Strings("amazon"),
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
		UserData:                 userData,
	})
}
