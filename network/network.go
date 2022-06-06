package network

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func AuthenticateAws(roleArn string, mfaSerialArn string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("Failed to load configuration, %v", err)
	}

	stsClient := sts.NewFromConfig(cfg)

	provider := stscreds.NewAssumeRoleProvider(stsClient, roleArn, func(o *stscreds.AssumeRoleOptions) {
		o.SerialNumber = aws.String(mfaSerialArn)
		o.TokenProvider = stscreds.StdinTokenProvider
	})

	cfg.Credentials = aws.NewCredentialsCache(provider)
	creds, err := cfg.Credentials.Retrieve(context.Background())

	if err != nil {
		log.Fatalf("Failed to retrieve role configuration, %v", err)
	}

	os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)

	println("Succesfully configured AWS credentials using role ARN: ", roleArn)
}

func CreateVpc(ctx *pulumi.Context) (*ec2.Vpc, error) {
	vpc, err := ec2.NewVpc(ctx, "ApplicationVpc", &ec2.VpcArgs{
		CidrBlock:        pulumi.String("10.0.0.0/16"),
		EnableDnsSupport: pulumi.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	return vpc, err
}

// func createSubnets(ctx *pulumi.Context, vpc *ec2.Vpc) {
// 	ec2.NewSubnet(ctx, "PrivateSubnet", &ec2.SubnetArgs{
// 		VpcId: pulumi.Any(vpc.),
// 		CidrBlock: pulumi.String("10.0.0.0/16"),
// 	})
// }

func Execute() {

	AuthenticateAws("", "")
	pulumi.Run(func(ctx *pulumi.Context) error {

		vpc, err := CreateVpc(ctx)

		if err != nil {
			return err
		}

		println(vpc)

		return nil
	})
}
