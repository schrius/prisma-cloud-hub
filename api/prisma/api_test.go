package prisma_test

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

	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma"
)

type mockHttpClient struct {
	mock.Mock
	prisma.PrismaHTTPiface
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

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *mockHttpClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *mockHttpClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestInvalidRequest(t *testing.T) {
	httpClient := &http.Client{}
	testCases := []struct {
		name     string
		client   *prisma.PrismaClient
		input    *prisma.PrismaAPIRequestInput
		expected error
	}{
		{
			name:     "empty prisma client token",
			client:   &prisma.PrismaClient{},
			input:    nil,
			expected: errors.New("required field Token type of string is empty"),
		},
		{
			name:     "empty prisma client Tenant",
			client:   &prisma.PrismaClient{Token: "token"},
			input:    nil,
			expected: errors.New("required field Tenant type of string is empty"),
		},
		{
			name:     "empty PrismaHTTPiface",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api"},
			input:    nil,
			expected: errors.New("required field PrismaHTTPiface type of prisma.PrismaHTTPiface is empty"),
		},
		{
			name:     "empty PrismaAPIRequestInput",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api", PrismaHTTPiface: httpClient},
			input:    nil,
			expected: errors.New("PrismaAPIRequestInput is nil"),
		},
		{
			name:     "empty PrismaAPIRequestInput.Endpoint",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api", PrismaHTTPiface: httpClient},
			input:    &prisma.PrismaAPIRequestInput{Action: http.MethodPost},
			expected: errors.New("required field Endpoint type of string is empty"),
		},
		{
			name:     "empty PrismaAPIRequestInput.Payload ",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api", PrismaHTTPiface: httpClient},
			input:    &prisma.PrismaAPIRequestInput{Action: http.MethodPost, Endpoint: "endpoint"},
			expected: errors.New("required field Payload type of []uint8 is empty"),
		},
		{
			name:     "empty PrismaAPIRequestInput.Header",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api", PrismaHTTPiface: httpClient},
			input:    &prisma.PrismaAPIRequestInput{Action: http.MethodPost, Endpoint: "endpoint", Payload: []byte("OK")},
			expected: errors.New("required field Header type of map[string]string is empty"),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			_, err := testCase.client.Request(testCase.input)
			assert.Equal(t, testCase.expected, err)
		})
	}
}

func TestRequestDo(t *testing.T) {
	mockClient := new(mockHttpClient)
	client := createMockHttpClient(mockClient)
	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Header.Get("x-redlock-auth") == "token" && req.URL.String() == "https://api.prismacloud.io/endpoint"
	})).Return(createHttpResponse(200, []byte("Body")), nil)
	testCases := []struct {
		name     string
		input    *prisma.PrismaAPIRequestInput
		expected io.ReadCloser
	}{
		{
			name:     "get request",
			input:    createPrismaAPIRequestInput(http.MethodGet),
			expected: ioutil.NopCloser(bytes.NewBufferString("Body")),
		},
		{
			name:     "post request",
			input:    createPrismaAPIRequestInput(http.MethodPost, []byte("payload")),
			expected: ioutil.NopCloser(bytes.NewBufferString("Body")),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			resp, err := client.Request(testCase.input)
			call := mockClient.Calls[i]
			assert.Equal(t, "Do", call.Method)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, resp)
		})
	}
	mockClient.AssertNumberOfCalls(t, "Do", len(testCases))
	mockClient.AssertExpectations(t)
}

func TestInvalidListAlerts(t *testing.T) {
	testCases := []struct {
		name     string
		client   *prisma.PrismaClient
		input    *prisma.ListAlertsInput
		expected error
	}{
		{
			name:     "empty prisma client token",
			client:   &prisma.PrismaClient{},
			input:    nil,
			expected: errors.New("required field Token type of string is empty"),
		},
		{
			name:     "empty prisma client Tenant",
			client:   &prisma.PrismaClient{Token: "token"},
			input:    nil,
			expected: errors.New("required field Tenant type of string is empty"),
		},
		{
			name:     "empty PrismaHTTPiface",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api"},
			input:    nil,
			expected: errors.New("required field PrismaHTTPiface type of prisma.PrismaHTTPiface is empty"),
		},
		{
			name:     "nil ListAlertsInput",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api", PrismaHTTPiface: &http.Client{}},
			input:    nil,
			expected: errors.New("ListAlertsInput is nil"),
		},
		{
			name:     "empty ListAlertsInput.Params",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api", PrismaHTTPiface: &http.Client{}},
			input:    &prisma.ListAlertsInput{},
			expected: errors.New("required field Params type of map[string]string is empty"),
		},
		{
			name:     "empty ListAlertsInput.ListAlertsPayload",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api", PrismaHTTPiface: &http.Client{}},
			input:    &prisma.ListAlertsInput{Params: map[string]string{"test": "test"}},
			expected: errors.New("required field ListAlertsPayload type of prisma.ListAlertsPayload is empty"),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			_, err := testCase.client.ListAlerts(testCase.input)
			assert.Equal(t, testCase.expected, err)
		})
	}
}

func TestListAlerts(t *testing.T) {
	mockClient := new(mockHttpClient)
	client := createMockHttpClient(mockClient)
	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Header.Get("x-redlock-auth") == "token" && req.URL.String() == "https://api.prismacloud.io/alert"
	})).Return(createHttpResponse(200, []byte(`[{"id": "123", "status":"Open"}]`)), nil)
	input := &prisma.ListAlertsInput{
		map[string]string{
			"detailed": "false",
		},
		prisma.ListAlertsPayload{
			prisma.Filters{
				{
					Name:     "policy.name",
					Value:    "AWS Region Violation",
					Operator: "=",
				},
			},
			[]string{
				"alert.id",
			},
		},
	}
	expectedResponse := &prisma.Alerts{{
		ID:     "123",
		Status: "Open",
	},
	}
	resp, err := client.ListAlerts(input)
	call := mockClient.Calls[0]
	assert.NoError(t, err)
	assert.Equal(t, "Do", call.Method)
	assert.Equal(t, expectedResponse, resp)
	mockClient.AssertExpectations(t)
}

func TestInvalidDismissAlerts(t *testing.T) {
	testCases := []struct {
		name     string
		client   *prisma.PrismaClient
		input    *prisma.DismissAlertInput
		expected error
	}{
		{
			name:     "nil DismissAlertInput",
			client:   &prisma.PrismaClient{},
			input:    nil,
			expected: errors.New("DismissAlertInput is nil"),
		},
		{
			name:     "nil Alerts",
			input:    &prisma.DismissAlertInput{},
			expected: errors.New("required field Alerts type of []string is empty"),
		},
		{
			name:     "empty dismissalNote",
			client:   &prisma.PrismaClient{},
			input:    &prisma.DismissAlertInput{Alerts: []string{"P-123"}},
			expected: errors.New("required field DismissalNote type of string is empty"),
		},
		{
			name:     "empty prisma client token",
			client:   &prisma.PrismaClient{},
			input:    &prisma.DismissAlertInput{Alerts: []string{"P-123"}, DismissalNote: "Test"},
			expected: errors.New("required field Token type of string is empty"),
		},
		{
			name:     "empty prisma client Tenant",
			client:   &prisma.PrismaClient{Token: "token"},
			input:    &prisma.DismissAlertInput{Alerts: []string{"P-123"}, DismissalNote: "Test"},
			expected: errors.New("required field Tenant type of string is empty"),
		},
		{
			name:     "empty PrismaHTTPiface",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api"},
			input:    &prisma.DismissAlertInput{Alerts: []string{"P-123"}, DismissalNote: "Test"},
			expected: errors.New("required field PrismaHTTPiface type of prisma.PrismaHTTPiface is empty"),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			_, err := testCase.client.DismissAlerts(testCase.input)
			assert.Equal(t, testCase.expected, err)
		})
	}
}

func TestDismissAlerts(t *testing.T) {
	mockClient := new(mockHttpClient)
	client := createMockHttpClient(mockClient)
	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Header.Get("x-redlock-auth") == "token" && req.URL.String() == "https://api.prismacloud.io/alert/dismiss"
	})).Return(createHttpResponse(200, []byte("Done")), nil)

	expectedResponse := []byte("Done")

	resp, err := client.DismissAlerts(&prisma.DismissAlertInput{
		Alerts:        []string{"I-123"},
		DismissalNote: "Test",
	})
	call := mockClient.Calls[0]
	assert.NoError(t, err)
	assert.Equal(t, "Do", call.Method)
	assert.Equal(t, expectedResponse, resp)

	mockClient.AssertExpectations(t)
}

func TestLoginPrismaError(t *testing.T) {
	// Login return token so we do no check token here
	testCases := []struct {
		name     string
		client   *prisma.PrismaClient
		input    *prisma.LoginPrismaInput
		expected error
	}{
		{
			name:     "empty prisma client Tenant",
			client:   &prisma.PrismaClient{Token: "token"},
			input:    nil,
			expected: errors.New("required field Tenant type of string is empty"),
		},
		{
			name:     "empty prisma client PrismaHTTPiface",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api"},
			input:    nil,
			expected: errors.New("required field PrismaHTTPiface type of prisma.PrismaHTTPiface is empty"),
		},
		{
			name:     "nil TokenRequest",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api", PrismaHTTPiface: &http.Client{}},
			input:    nil,
			expected: errors.New("LoginPrismaInput is nil"),
		},
		{
			name:     "empty Auth Input",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api", PrismaHTTPiface: &http.Client{}},
			input:    &prisma.LoginPrismaInput{},
			expected: errors.New("required field Auth type of []uint8 is empty"),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			err := testCase.client.LoginPrisma(testCase.input)
			assert.Equal(t, testCase.expected, err)
		})
	}
}

func TestLoginPrisma(t *testing.T) {
	mockClient := new(mockHttpClient)
	client := &prisma.PrismaClient{
		Tenant:          "api",
		PrismaHTTPiface: mockClient,
	}

	expected := ioutil.NopCloser(bytes.NewBuffer([]byte(`{"token":"token", "message":"message"}`)))
	authInput := []byte(`{"Test":"Test"}`)
	mockClient.On("Post", "https://api.prismacloud.io/login", "application/json", bytes.NewBuffer(authInput)).Return(&http.Response{
		StatusCode: 200,
		Body:       expected,
		Header:     make(http.Header),
	}, nil)

	err := client.LoginPrisma(&prisma.LoginPrismaInput{
		authInput,
	})
	assert.NoError(t, err)
	assert.Equal(t, "token", client.Token)
	mockClient.AssertExpectations(t)
}

func TestListAccountGroupsWithInvalidClient(t *testing.T) {
	testCases := []struct {
		name     string
		client   *prisma.PrismaClient
		expected error
	}{
		{
			name:     "nil prismaClient",
			client:   &prisma.PrismaClient{},
			expected: errors.New("required field Token type of string is empty"),
		},
		{
			name:     "empty prisma client Tenant",
			client:   &prisma.PrismaClient{Token: "token"},
			expected: errors.New("required field Tenant type of string is empty"),
		},
		{
			name:     "empty prisma client PrismaHTTPiface",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api"},
			expected: errors.New("required field PrismaHTTPiface type of prisma.PrismaHTTPiface is empty"),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			_, err := testCase.client.ListAccountGroups()
			assert.Equal(t, testCase.expected, err)
		})
	}
}

func TestListAccountGroups(t *testing.T) {
	mockClient := new(mockHttpClient)
	prismaClient := createMockHttpClient(mockClient)
	expectedAccountGroup := &prisma.AccountGroups{
		{
			ID:   "Test",
			Name: "Test",
		},
	}
	expectedResponse, _ := json.Marshal(expectedAccountGroup)
	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Header.Get("x-redlock-auth") == "token" && req.URL.String() == "https://api.prismacloud.io/cloud/group/name"
	})).Return(createHttpResponse(200, expectedResponse), nil)

	actualData, err := prismaClient.ListAccountGroups()
	call := mockClient.Calls[0]
	assert.NoError(t, err)
	assert.Equal(t, "Do", call.Method)
	assert.Equal(t, expectedAccountGroup, actualData)
}

func TestRegisterAccountErrorHandling(t *testing.T) {
	testCases := []struct {
		name     string
		client   *prisma.PrismaClient
		expected error
	}{
		{
			name:     "nil prismaClient",
			client:   &prisma.PrismaClient{},
			expected: errors.New("required field Token type of string is empty"),
		},
		{
			name:     "empty prisma client Tenant",
			client:   &prisma.PrismaClient{Token: "token"},
			expected: errors.New("required field Tenant type of string is empty"),
		},
		{
			name:     "empty prisma client PrismaHTTPiface",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api"},
			expected: errors.New("required field PrismaHTTPiface type of prisma.PrismaHTTPiface is empty"),
		},
		{
			name:     "empty prisma client PrismaHTTPiface",
			client:   &prisma.PrismaClient{Token: "token", Tenant: "api", PrismaHTTPiface: &http.Client{}},
			expected: errors.New("required parameter payload is empty"),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			err := testCase.client.RegisterAccount(nil)
			assert.Equal(t, testCase.expected, err)
		})
	}
}

func TestRegisterAccount(t *testing.T) {
	mockHTTP := new(mockHttpClient)
	prismaClient := createMockHttpClient(mockHTTP)
	expected := []byte(`done`)
	mockHTTP.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Header.Get("x-redlock-auth") == "token" && req.URL.String() == "https://api.prismacloud.io/cloud/cloud_type"
	})).Return(createHttpResponse(200, expected), nil)

	err := prismaClient.RegisterAccount([]byte("test"))
	call := mockHTTP.Calls[0]
	assert.NoError(t, err)
	assert.Equal(t, "Do", call.Method)
}

func TestDefaultHeader(t *testing.T) {
	expect := map[string]string{
		"Content-Type":   "application/json",
		"x-redlock-auth": "token",
	}
	result := prisma.DefaultHeader("token")
	assert.Equal(t, expect, result)
}

func createPrismaAPIRequestInput(action string, payload ...[]byte) *prisma.PrismaAPIRequestInput {
	if len(payload) == 0 {
		return &prisma.PrismaAPIRequestInput{
			Action:   action,
			Endpoint: "endpoint",
			Header:   prisma.DefaultHeader("token"),
			Payload:  nil,
		}
	}
	return &prisma.PrismaAPIRequestInput{
		Action:   action,
		Endpoint: "endpoint",
		Header:   prisma.DefaultHeader("token"),
		Payload:  payload[0],
	}
}

func createHttpResponse(statusCode int, payload ...[]byte) *http.Response {
	body := []byte{}
	if len(payload) != 0 {
		body = payload[0]
	}
	return &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
		Header: http.Header{
			"Content-Type": {
				"application/json",
			},
		},
	}
}

// creat mock http client match header and url
func createMockHttpClient(mockClient *mockHttpClient) *prisma.PrismaClient {
	return &prisma.PrismaClient{
		"token",
		"api",
		mockClient,
	}
}
