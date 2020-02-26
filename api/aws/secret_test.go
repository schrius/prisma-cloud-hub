package aws_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/stretchr/testify/assert"

	awssecret "github.com/CityOfNewYork/prisma-cloud-remediation/api/aws"
)

type mockSecretsManager struct {
	secretsmanageriface.SecretsManagerAPI
	resp secretsmanager.GetSecretValueOutput
}

func (m *mockSecretsManager) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	if input == nil {
		return nil, errors.New("Missing required parameter")
	}
	return &m.resp, nil
}

func TestGetSecret(t *testing.T) {
	secretString := "{\"Key\":\"Key\", \"ID\":\"ID\", \"ExternalID\":\"ExternalID\"}"
	resp := &secretsmanager.GetSecretValueOutput{SecretString: &secretString}
	mockSvc := &mockSecretsManager{resp: *resp}
	expected := awssecret.Secret{
		Key:        "Key",
		ID:         "ID",
		ExternalID: "ExternalID",
	}

	_, err := awssecret.GetSecret("", mockSvc)
	assert.Error(t, err)

	result, err := awssecret.GetSecret("secret", mockSvc)
	assert.Equal(t, expected, *result)
	assert.NoError(t, err)

}
