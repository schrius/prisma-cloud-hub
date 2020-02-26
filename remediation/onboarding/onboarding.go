package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CityOfNewYork/prisma-cloud-remediation/api"
	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma/prismaiface"
	"github.com/CityOfNewYork/prisma-cloud-remediation/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const (
	tenant     = "api3"
	roleArn    = "arn:aws:iam::1234567890:role/PrismaCloudReadOnlyRole"
	externalID = "FrenchEllaReturns"
)

func login(svc *secretsmanager.SecretsManager) (prismaiface.PrismaAPI, error) {
	prismaClient := api.CreatePrismaClient("api3")
	client, err := api.LoginPrismaWithAWSSecret("Prisma", "OnBoarding", svc, prismaClient)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func handler(ctx context.Context, event events.OnBoardEvent) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	prismaClient, err := login(secretsmanager.New(sess))
	if err != nil {
		fmt.Println("Login failed")
		fmt.Println(err.Error())
		return err
	}

	accountGroups, accoutGroupsErr := prismaClient.ListAccountGroups()
	if accoutGroupsErr != nil {
		fmt.Println(accoutGroupsErr.Error())
		return accoutGroupsErr
	}
	accountGroupIDs := api.GetAccountGroupID(event.GroupNames, accountGroups)

	switch event.CloudType {
	case "aws":
		if accountID := looUpAccount(event.AWS.Name, prismaClient); accountID != "" {
			return nil
		}
		event.AWS.GroupIds = accountGroupIDs
		payload, err := json.Marshal(event.AWS)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if err := prismaClient.RegisterAccount(payload); err != nil {
			fmt.Println(err.Error())
			return err
		}

	case "azure":
		if accountID := looUpAccount(event.Azure.CloudAccount.Name, prismaClient); accountID != "" {
			return nil
		}
		event.Azure.CloudAccount.GroupIds = accountGroupIDs
		payload, err := json.Marshal(event.Azure)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if err := prismaClient.RegisterAccount(payload); err != nil {
			fmt.Println(err.Error())
			return err
		}

	case "gcp":
		if accountID := looUpAccount(event.GCP.CloudAccount.Name, prismaClient); accountID != "" {
			return nil
		}
		event.GCP.CloudAccount.GroupIds = accountGroupIDs
		payload, err := json.Marshal(event.GCP)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if err := prismaClient.RegisterAccount(payload); err != nil {
			fmt.Println(err.Error())
			return err
		}

	default:
		fmt.Println("Not support cloud type: %s\n", event.CloudType)
	}
	return nil
}

func looUpAccount(name string, prismaClient prismaiface.PrismaAPI) string {
	accountID, accountIDErr := api.LookUpAccountName(name, prismaClient)
	if accountIDErr != nil {
		fmt.Println(accountIDErr.Error())
		return ""
	}

	if accountID != "" {
		fmt.Println("Account exist")
		return accountID
	}
	return ""
}

func main() {
	lambda.Start(handler)
}
