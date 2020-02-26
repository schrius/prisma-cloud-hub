package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type AlertName string

// Alert struct only partially parse the alert. Most of information is omitted
type Alert struct {
	ResourceID    string      `json:"resourceId"`
	AlertRuleName AlertName   `json:"alertRuleName"`
	AccountName   string      `json:"accountName"`
	CloudType     string      `json:"cloudType"`
	Severity      string      `json:"severity"`
	PolicyName    string      `json:"policyName"`
	Resource      interface{} `json:"resource"`
}

func handler(ctx context.Context, event Alert) error {

	fmt.Println("Event received.")

	message := fmt.Sprintf(`
	Alert Rule: %s
	Resource: %s
	CloudType: %s
	Account: %s
	Policy: %s
	ResourceID: %s`, event.AlertRuleName, event.ResourceID, event.CloudType, event.AccountName, event.PolicyName, event.ResourceID)

	fmt.Printf("Alert: %s\n", message)
	err := sendSNS(message)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Notification sent.")
	return nil
}

func sendSNS(message string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	snsClient := sns.New(sess, &aws.Config{Region: aws.String(os.Getenv("REGION"))})

	fmt.Println("Sending notification.")
	_, err := snsClient.Publish(&sns.PublishInput{
		TopicArn: aws.String(os.Getenv("SNS")),
		Message:  aws.String(message),
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
