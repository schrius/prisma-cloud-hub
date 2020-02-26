package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/CityOfNewYork/prisma-cloud-remediation/api"
	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma"
	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma/alert"
	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma/prismaiface"
	"github.com/CityOfNewYork/prisma-cloud-remediation/events"
)

func login(svc *secretsmanager.SecretsManager) (prismaiface.PrismaAPI, error) {
	prismaClient := api.CreatePrismaClient("api3")
	client, err := api.LoginPrismaWithAWSSecret("Prisma", "FalseAlertDismisser", svc, prismaClient)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func handler(contxt context.Context, event events.Alert) error {
	if !alert.FalseAWSRegionViolationAlert(&event) {
		fmt.Println("This is not a false alert.")
		return nil
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	prismaClient, err := login(secretsmanager.New(sess))
	if err != nil {
		fmt.Println("Login failed")
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("Dimiss AlertID: %s\n", event.AlertID)
	dismissErr := api.DismissAlert(&prisma.DismissAlertInput{
		Alerts:        []string{event.AlertID},
		DismissalNote: "Non Virginia Region",
	}, prismaClient)
	if dismissErr != nil {
		fmt.Printf("Dismiss alert failed: \n%s\n", dismissErr.Error())
		return dismissErr
	}
	fmt.Printf("Aleret %s is dismissed\n", event.AlertID)
	return nil
}

func main() {
	lambda.Start(handler)
}
