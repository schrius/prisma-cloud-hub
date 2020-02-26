package examples

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/CityOfNewYork/prisma-cloud-remediation/api"
	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma"
)

// coming soon
func DismissAlert() {
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

	size := 50150
	input := &prisma.DismissAlertInput{
		Alerts:        []string{},
		DismissalNote: "Non Virginia Region",
	}

	for i := 50066; i < size; i++ {
		input.Alerts = append(input.Alerts, "P-"+strconv.Itoa(i))
	}
	fmt.Println(input)

	fmt.Println("Request for alerts")
	dismissErr := api.DismissAlert(input, client)
	if dismissErr != nil {
		fmt.Printf("Dismiss alert failed: \n%s\n", dismissErr.Error())
	}

}
