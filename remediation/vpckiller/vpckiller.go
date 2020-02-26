package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/CityOfNewYork/prisma-cloud-remediation/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Client contain session and configure of the EC2 Client to cal ec2 API
type Client struct {
	session *session.Session
	config  *aws.Config
	ec2svc  *ec2.EC2
}

// Role is the assume role that used to access other accounts with the ExternalID
const (
	Role       = "Prisma_VPC_Term_Role"
	ExternalID = "PrismaVPCKiller"
	Virginia   = "us-east-1"
	DevAccount = "1234567890"
)

func (c *Client) initClient(region string, accountID string, role string, externalID string) error {
	if c.session == nil {
		c.session = session.Must(session.NewSession())
	}

	arn := fmt.Sprintf("arn:aws:iam::%v:role/%v", accountID, role)

	creds := stscreds.NewCredentials(
		c.session,
		arn,
		func(provider *stscreds.AssumeRoleProvider) {
			provider.ExternalID = aws.String(externalID)
		})

	c.config = &aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	}

	c.ec2svc = ec2.New(c.session, c.config)

	return nil
}

func handler(ctx context.Context, event events.VPCAlert) error {
	client := &Client{}
	var vpcIds []*string

	fmt.Printf("%+v\n", event)
	if event.ResourceRegionID == Virginia && event.AccountID != DevAccount {
		fmt.Printf("Region: %s is allowed in %s", event.ResourceRegionID, event.AccountID)
	}

	client.initClient(event.ResourceRegionID, event.AccountID, Role, ExternalID)
	vpcs, verifyError := client.verifyVpcs([]*string{
		aws.String(event.ResourceID),
	})
	if verifyError != nil {
		fmt.Println(verifyError.Error())
		return verifyError
	}

	for _, vpc := range vpcs {
		vpcIds = append(vpcIds, vpc.VpcId)
	}

	vpcFilters := []*ec2.Filter{
		{
			Name:   aws.String("vpc-id"),
			Values: vpcIds,
		},
	}

	if err := client.cleanInternetGateways(vpcIds); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err := client.cleanEC2Instances(vpcFilters); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err := client.cleanSecurityGroups(vpcFilters); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err := client.cleanSubnets(vpcFilters); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err := client.deleteVpc(aws.String(event.ResourceID)); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("VPC Killer tasks: done")
	return nil
}

func (c *Client) deleteVpc(vpcID *string) error {
	fmt.Printf("Deleting Vpc %s\n", *vpcID)
	_, err := c.ec2svc.DeleteVpc(&ec2.DeleteVpcInput{
		VpcId: vpcID,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("VPC %s has been deleted\n", *vpcID)
	return nil
}

func (c *Client) cleanEC2Instances(filters []*ec2.Filter) error {
	if filters == nil {
		return errors.New("deleteEC2Instances: missing required variable vpdId")
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	instances, err := c.ec2svc.DescribeInstances(input)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var instanceIds []*string

	if len(instances.Reservations) < 1 {
		fmt.Println("No instances was found")
		return nil
	}
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			instanceIds = append(instanceIds, instance.InstanceId)
		}
	}

	fmt.Printf("%d Instances found", len(instanceIds))
	termInput := &ec2.TerminateInstancesInput{
		InstanceIds: instanceIds,
	}

	fmt.Println("Terminating Instances...")
	_, termError := c.ec2svc.TerminateInstances(termInput)

	if termError != nil {
		fmt.Println(termError.Error())
		return termError
	}

	fmt.Println("Waiting for instances to be terminated")
	if err := c.waitUntilInstanceTerminated(
		&ec2.DescribeInstancesInput{
			Filters: filters,
		}); err != nil {
		return err
	}
	fmt.Println("All instances has been terminiated")

	return nil
}

func (c *Client) waitUntilInstanceTerminated(input *ec2.DescribeInstancesInput) error {
	err := c.ec2svc.WaitUntilInstanceTerminated(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (c *Client) cleanSecurityGroups(filters []*ec2.Filter) error {
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: filters,
	}

	fmt.Println("Searching custom security groups")
	sgs, err := c.ec2svc.DescribeSecurityGroups(input)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if len(sgs.SecurityGroups) == 1 && *sgs.SecurityGroups[0].GroupName == "default" {
		fmt.Println("No custom security group was found")
		return nil
	}
	for _, sg := range sgs.SecurityGroups {
		if *sg.GroupName != "default" {
			fmt.Printf("Deleting custom security group: %s\n", *sg.GroupId)
			_, err := c.ec2svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
				GroupId: sg.GroupId,
			})
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			fmt.Printf("%s has deleted", *sg.GroupId)
		}
	}
	fmt.Println("All custom security groups have been deleted")
	return nil
}

func (c *Client) cleanSubnets(filters []*ec2.Filter) error {
	input := &ec2.DescribeSubnetsInput{
		Filters: filters,
	}
	output, err := c.ec2svc.DescribeSubnets(input)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if len(output.Subnets) == 0 {
		fmt.Println("No subnets was found")
		return nil
	}

	for _, subnet := range output.Subnets {
		fmt.Printf("Deleting subnet: %s\n", *subnet.SubnetId)
		_, err := c.ec2svc.DeleteSubnet(&ec2.DeleteSubnetInput{
			SubnetId: subnet.SubnetId,
		})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		fmt.Printf("%s has deleted", *subnet.SubnetId)
	}
	fmt.Println("All subnets has been deleted")

	return nil
}

func (c *Client) cleanInternetGateways(vpcID []*string) error {

	input := &ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: vpcID,
			},
		},
	}

	output, err := c.ec2svc.DescribeInternetGateways(input)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if len(output.InternetGateways) == 0 {
		fmt.Println("No internet gateways was found")
		return nil
	}

	for _, internetgateway := range output.InternetGateways {
		fmt.Printf("Deleting InternetGateway: %s\n", *internetgateway.InternetGatewayId)
		for _, attachment := range internetgateway.Attachments {
			_, err := c.ec2svc.DetachInternetGateway(&ec2.DetachInternetGatewayInput{
				InternetGatewayId: internetgateway.InternetGatewayId,
				VpcId:             attachment.VpcId,
			})
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			fmt.Printf("%s has detached from %s", *internetgateway.InternetGatewayId, *attachment.VpcId)
		}

		_, derr := c.ec2svc.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
			InternetGatewayId: internetgateway.InternetGatewayId,
		})
		if derr != nil {
			fmt.Println(derr.Error())
			return derr
		}
		fmt.Printf("%s has been deleted", *internetgateway.InternetGatewayId)
	}
	fmt.Println("All InternetGateways has been detached")

	return nil
}

func (c *Client) verifyVpcs(vpcIds []*string) ([]*ec2.Vpc, error) {
	fmt.Println("Searching Vpcs")
	output, err := c.ec2svc.DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: vpcIds,
	},
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	if len(output.Vpcs) == 0 {
		return nil, fmt.Errorf("Vpcs not found")
	}

	fmt.Printf("verify %d Vpcs found %d\n", len(vpcIds), len(output.Vpcs))
	return output.Vpcs, nil
}

func awsErrorHandler(err error) error {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return err
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
