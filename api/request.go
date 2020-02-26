package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"

	awssecret "github.com/CityOfNewYork/prisma-cloud-remediation/api/aws"
	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma"
	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma/prismaiface"
)

const (
	tenant  = "api3"
	roleArn = "arn:aws:iam::812969027137:role/PrismaCloudReadOnlyRole"
)

// Authenticate prisma Authenticate
type Authenticate struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	CustomerName string `json:"customerName"`
}

// AWSSecretSvc return Secrete svc
func AWSSecretSvc(awsRegion string) *secretsmanager.SecretsManager {
	return secretsmanager.New(session.New(), &aws.Config{Region: aws.String(awsRegion)})
}

// CreatePrismaClient accept tenant and return a new prismaClient
func CreatePrismaClient(tenant string) *prisma.PrismaClient {
	return &prisma.PrismaClient{
		Tenant:          tenant,
		PrismaHTTPiface: &http.Client{},
	}
}

// LoginPrismaWithAWSSecret request secret from AWS Secret manager and login with the id and key
func LoginPrismaWithAWSSecret(secret string, customerName string, svc secretsmanageriface.SecretsManagerAPI, client prismaiface.PrismaAPI) (prismaiface.PrismaAPI, error) {
	auth, err := GetAuth(secret, customerName, svc)
	if err != nil {
		return nil, err
	}
	if err := client.LoginPrisma(&prisma.LoginPrismaInput{
		Auth: auth,
	}); err != nil {
		return nil, err
	}
	return client, nil
}

// GetAuth return authenticate to get token
func GetAuth(secret string, customerName string, svc secretsmanageriface.SecretsManagerAPI) ([]byte, error) {
	secretKey, err := awssecret.GetSecret(secret, svc)
	if err != nil {
		return nil, err
	}
	auth, err := json.Marshal(&Authenticate{Username: secretKey.ID, Password: secretKey.Key, CustomerName: customerName})
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// GetAccountGroupID return arrary of string contain the list of all account groups
func GetAccountGroupID(groupNames []string, accountsGroups *prisma.AccountGroups) []string {
	accountIDs := []string{}

	accountsGroupMap := map[string]string{}

	for _, accountsGroup := range *accountsGroups {
		accountsGroupMap[accountsGroup.Name] = accountsGroup.ID
	}
	for _, groupName := range groupNames {
		if groupID, ok := accountsGroupMap[groupName]; ok {
			accountIDs = append(accountIDs, groupID)
		}
	}
	return accountIDs
}

// DismissAlert call Prisma Client to dismiss the alert pass in the DismissAlertInputer
// it returns an error if error occur
func DismissAlert(input *prisma.DismissAlertInput, prismaClient prismaiface.PrismaAPI) error {
	fmt.Println("Dismissing alert...")
	resp, err := prismaClient.DismissAlerts(input)
	if err != nil {
		return err
	}
	fmt.Println(string(resp))
	return nil
}

// LookUpAccountName return true if the account exist
func LookUpAccountName(accountName string, prismaClient prismaiface.PrismaAPI) (string, error) {
	resp, err := prismaClient.ListAccountNames()
	if err != nil {
		return "", err
	}

	for _, account := range *resp {
		if account.Name == accountName {
			return account.ID, nil
		}
	}
	return "", fmt.Errorf("No account name: %s", accountName)
}

// DefaultHeader accept token and return an http.request.header
func DefaultHeader(token string) map[string]string {
	return map[string]string{
		"Content-Type":   "application/json",
		"x-redlock-auth": token,
	}
}

// RefreshSession refresh the prismaClient token
func RefreshSession(token string, prismaClient prismaiface.PrismaAPI) error {
	resp, err := prismaClient.Request(&prisma.PrismaAPIRequestInput{
		"Get",
		"auth_token/extend",
		nil,
		DefaultHeader(token),
	})
	if err != nil {
		return err
	}
	resp.Close()
	return nil
}
