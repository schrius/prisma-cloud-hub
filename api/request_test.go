package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/CityOfNewYork/prisma-cloud-remediation/api"
	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma"
	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma/prismaiface"
)

type mockHttpClient struct {
	mock.Mock
	prisma.PrismaHTTPiface
}

type mockPrismaClient struct {
	Token string
	mock.Mock
	prismaiface.PrismaAPI
}

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

func (m *mockPrismaClient) LoginPrisma(input *prisma.LoginPrismaInput) error {
	args := m.Called(input)
	m.Token = "token"
	return args.Error(0)
}

func (m *mockPrismaClient) Request(input *prisma.PrismaAPIRequestInput) (io.ReadCloser, error) {
	args := m.Called(input)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *mockHttpClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *mockPrismaClient) ListAccountNames() (*prisma.AccountNames, error) {
	args := m.Called()
	return args.Get(0).(*prisma.AccountNames), args.Error(1)
}

func TestCreatePrismaClient(t *testing.T) {
	client := api.CreatePrismaClient("api")
	assert.Equal(t, "api", client.Tenant)
	assert.IsType(t, &http.Client{}, client.PrismaHTTPiface)
}

func TestGetAuth(t *testing.T) {
	secretString := "{\"Key\":\"TestKey\", \"ID\":\"TestID\", \"ExternalID\":\"ExternalID\"}"
	response := &secretsmanager.GetSecretValueOutput{SecretString: &secretString}
	mockSvc := &mockSecretsManager{resp: *response}
	resp, err := api.GetAuth("Test", "TestName", mockSvc)
	auth, _ := json.Marshal(&api.Authenticate{Username: "TestID", Password: "TestKey", CustomerName: "TestName"})
	assert.NoError(t, err)
	assert.Equal(t, auth, resp)
}

func TestGetAccountGroupsID(t *testing.T) {
	accountGroups := &prisma.AccountGroups{
		{
			ID:   "123",
			Name: "Test1",
		},
		{

			ID:   "456",
			Name: "Test2",
		},
		{

			ID:   "789",
			Name: "Test3",
		},
	}
	testCases := []struct {
		name          string
		groupNames    []string
		accountGroups *prisma.AccountGroups
		expected      []string
	}{
		{
			name:       "1 Groups Name match",
			groupNames: []string{"Test1"},
			expected:   []string{"123"},
		},
		{
			name:       "1 Groups Name match",
			groupNames: []string{"Test1", "Test2"},
			expected:   []string{"123", "456"},
		},
		{
			name:       "0 Group Name match",
			groupNames: []string{"Test4"},
			expected:   []string{},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			groupIDs := api.GetAccountGroupID(testCase.groupNames, accountGroups)
			assert.Equal(t, testCase.expected, groupIDs)
		})
	}
}

func TestLoginPrismaWithAWSSecret(t *testing.T) {

	secretString := "{\"Key\":\"Key\", \"ID\":\"ID\", \"ExternalID\":\"ExternalID\"}"
	resp := &secretsmanager.GetSecretValueOutput{SecretString: &secretString}
	mockSvc := &mockSecretsManager{resp: *resp}

	mockClient := new(mockPrismaClient)
	mockClient.On("LoginPrisma", mock.Anything).Return(nil)

	_, err := api.LoginPrismaWithAWSSecret("test", "Test", mockSvc, mockClient)

	call := mockClient.Calls[0]
	assert.NoError(t, err)
	assert.Equal(t, "LoginPrisma", call.Method)
	assert.Equal(t, "token", mockClient.Token)
}

func TestLooUpAccountName(t *testing.T) {

	response := &prisma.AccountNames{
		{
			CloudType: "aws",
			Name:      "Test1",
			ID:        "123",
		},
		{
			CloudType: "aws",
			Name:      "Test2",
			ID:        "456",
		},
		{
			CloudType: "aws",
			Name:      "Test3",
			ID:        "789",
		},
	}
	testCases := []struct {
		name          string
		accountName   string
		expected      string
		expectedError error
	}{
		{
			name:          "look up existing account",
			accountName:   "Test1",
			expected:      "123",
			expectedError: nil,
		},
		{
			name:          "look up existing account",
			accountName:   "Test2",
			expected:      "456",
			expectedError: nil,
		},
		{
			name:          "look up not existing account",
			accountName:   "Test4",
			expected:      "",
			expectedError: errors.New("No account name: Test4"),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			mockPrisma := new(mockPrismaClient)
			mockPrisma.On("ListAccountNames", mock.Anything).Return(response, nil)
			actualResponse, err := api.LookUpAccountName(testCase.accountName, mockPrisma)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expected, actualResponse)
		})
	}
}

func TestRefreshSession(t *testing.T) {
	mockClient := new(mockPrismaClient)
	mockClient.On("Request", mock.Anything).Return(ioutil.NopCloser(bytes.NewBufferString("Test")), nil)

	err := api.RefreshSession("token", mockClient)
	assert.NoError(t, err)
}
