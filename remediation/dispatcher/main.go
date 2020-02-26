package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	invokeLambda "github.com/aws/aws-sdk-go/service/lambda"
)

type AlertName string

// Alert struct only partially parse the alert. Most of information is omitted
type Alert struct {
	ResourceID    string    `json:"resourceId"`
	AlertRuleName AlertName `json:"alertRuleName"`
	AccountName   string    `json:"accountName"`
	CloudType     string    `json:"cloudType"`
	Severity      string    `json:"severity"`
	PolicyName    string    `json:"policyName"`
}

const (
	SuspiciousTrafficAlert AlertName = "Suspicious Traffic Alert"
	VPCKiller              AlertName = "VPCKiller"
	ScienceLogic           AlertName = "ScienceLogic"
	RegionViolation        AlertName = "AWS Region Violation"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {

		fmt.Printf("The message %s for event source %s \n", message.MessageId, message.EventSource)

		var alert Alert

		if err := json.Unmarshal([]byte(message.Body), &alert); err != nil {
			return err
		}

		switch alert.AlertRuleName {
		case SuspiciousTrafficAlert:
			{
				fmt.Println(alert.AlertRuleName)
				invokeFunction("PrismaAlertNotification", "Event", []byte(message.Body))
			}
		case VPCKiller:
			{
				fmt.Println(alert.AlertRuleName)
				invokeFunction("PrismaVPCKiller", "Event", []byte(message.Body))
			}
		case ScienceLogic:
			{
				fmt.Println(alert.AlertRuleName)
				invokeFunction("PrismaScienceLogic", "Event", []byte(message.Body))
			}
		case RegionViolation:
			{
				fmt.Println(alert.AlertRuleName)
				invokeFunction("PrismaFalseAlertRemover", "Event", []byte(message.Body))
			}
		default:
			fmt.Printf("Unsupported Alert: %s ", alert.AlertRuleName)
		}

	}

	return nil

}

func invokeFunction(functionName string, invocationType string, payload []byte) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := invokeLambda.New(sess, &aws.Config{Region: aws.String(os.Getenv("REGION"))})

	client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String(functionName), InvocationType: aws.String(invocationType), Payload: payload})
}

func main() {
	lambda.Start(handler)
}
