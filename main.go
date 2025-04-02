package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Cambiar la región a "us-east-2" para la capa gratuita
		provider, err := aws.NewProvider(ctx, "aws-provider", &aws.ProviderArgs{
			Region: pulumi.String("us-east-2"),
		})
		if err != nil {
			return fmt.Errorf("error creating AWS provider: %v", err)
		}

		// Crear un VPC por defecto
		defaultVpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{
			Default: pulumi.BoolRef(true),
		}, pulumi.Provider(provider))
		if err != nil {
			return fmt.Errorf("error finding default VPC: %v", err)
		}

		// Crear una subred pública dentro del VPC por defecto
		subnet, err := ec2.GetSubnets(ctx, &ec2.GetSubnetsArgs{
			Filters: []ec2.GetSubnetsFilter{
				{
					Name:   "vpc-id",
					Values: []string{defaultVpc.Id},
				},
				{
					Name:   "default-for-az",
					Values: []string{"true"},
				},
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return fmt.Errorf("error finding default subnet: %v", err)
		}

		// Crear un grupo de seguridad para permitir el tráfico SSH y HTTP
		sg, err := ec2.NewSecurityGroup(ctx, "web-sg", &ec2.SecurityGroupArgs{
			Description: pulumi.String("Allow SSH and HTTP"),
			VpcId:       pulumi.String(defaultVpc.Id),
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					Protocol: pulumi.String("tcp"),
					FromPort: pulumi.Int(22),
					ToPort:   pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				&ec2.SecurityGroupIngressArgs{
					Protocol: pulumi.String("tcp"),
					FromPort: pulumi.Int(80),
					ToPort:   pulumi.Int(80),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				&ec2.SecurityGroupEgressArgs{
					Protocol: pulumi.String("-1"), // Permitir todo el tráfico de salida
					FromPort: pulumi.Int(0),
					ToPort:   pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return fmt.Errorf("error creating security group: %v", err)
		}

		// Crear una instancia EC2 dentro de la capa gratuita (t2.micro) con Ubuntu 22.04
		instance, err := ec2.NewInstance(ctx, "web-instance", &ec2.InstanceArgs{
			Ami:             pulumi.String("ami-08c18c49bdb7f38f7"), // Ubuntu 22.04 LTS en la región us-east-2
			InstanceType:    pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			SubnetId:        pulumi.String(subnet.Ids[0]),
			KeyName:         pulumi.String("ec2-django-key"),
			Tags: 	pulumi.StringMap{"Name": pulumi.String("WebInstance")},
		}, pulumi.Provider(provider))
		if err != nil {
			return fmt.Errorf("error creating EC2 instance: %v", err)
		}

		// Exportar la dirección pública de la instancia
		ctx.Export("instancePublicIP", instance.PublicIp)

		return nil

	})
}