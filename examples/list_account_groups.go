package examples

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/CityOfNewYork/prisma-cloud-remediation/api"
)

func ListAccountGroups() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	svc := secretsmanager.New(sess)

	prismaClient := api.CreatePrismaClient("api3")
	client, err := api.LoginPrismaWithAWSSecret("Prisma", "AlertDismisser", svc, prismaClient)
	if err != nil {
		fmt.Println("Login failed")
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Request for account groups")
	resp, err := client.ListAccountGroups()
	if err != nil {
		fmt.Println("Error: call list Account groups failed")
	}
	fmt.Println(resp)
}
