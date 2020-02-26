package prisma

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/CityOfNewYork/prisma-cloud-remediation/errors"
)

type PrismaHTTPiface interface {
	CloseIdleConnections()
	Get(string) (*http.Response, error)
	Head(string) (*http.Response, error)
	Do(*http.Request) (*http.Response, error)
	Post(string, string, io.Reader) (*http.Response, error)
	PostForm(string, url.Values) (*http.Response, error)
}
type PrismaClient struct {
	Token  string
	Tenant string
	PrismaHTTPiface
}

type AccountGroups []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CloudResponse []struct {
	Name           string   `json:"name"`
	CloudType      string   `json:"cloudType"`
	AccountType    string   `json:"accountType"`
	Enabled        bool     `json:"enabled"`
	LastModifiedTs int64    `json:"lastModifiedTs"`
	LastModifiedBy string   `json:"lastModifiedBy"`
	IngestionMode  int      `json:"ingestionMode"`
	GroupIds       []string `json:"groupIds"`
	Groups         []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"groups"`
	Status                string `json:"status"`
	NumberOfChildAccounts int    `json:"numberOfChildAccounts"`
	AccountID             string `json:"accountId"`
	AddedOn               int64  `json:"addedOn"`
}

type AccountNames []struct {
	CloudType         string `json:"cloudType"`
	Name              string `json:"name"`
	ID                string `json:"id"`
	ParentAccountName string `json:"parentAccountName,omitempty"`
}

// PrismaAPIRequestInput input for API request
type PrismaAPIRequestInput struct {
	Action   string
	Endpoint string
	Payload  []byte
	Header   map[string]string
}

type PrismaAPIRequestOutput struct {
	Body io.ReadCloser
}

type ListAlertsPayload struct {
	Filters Filters  `json:"filters"`
	Fields  []string `json:"fields"`
}

type ListAlertsInput struct {
	Params map[string]string
	ListAlertsPayload
}

type Filters []struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}

type ListAlertsResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Policy struct {
		Name       string `json:"name"`
		PolicyType string `json:"policyType"`
		Remediable bool   `json:"remediable"`
	} `json:"policy"`
	Resource struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Account string `json:"account"`
		Region  string `json:"region"`
	} `json:"resource"`
}

// DismissAlertInput DissmissAlerts parameter
type DismissAlertInput struct {
	Alerts             []string           `json:"alerts"`
	DismissalNote      string             `json:"dismissalNote"`
	DismissAlertFilter DismissAlertFilter `json:"filter"`
}

type DismissAlertFilter struct {
	TimeRange FilterTimeRange `json:"timeRange"`
	Filters   Filters         `json:"filters"`
}

type FilterTimeRange struct {
	Type  string         `json:"type"`
	Value TimeRangeValue `json:"value"`
}

type TimeRangeValue struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

type Alerts []struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Policy struct {
		PolicyID   string `json:"policyId"`
		PolicyType string `json:"policyType"`
	} `json:"policy,omitempty"`
	Resource struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Account      string `json:"account"`
		AccountID    string `json:"accountId"`
		Region       string `json:"region"`
		RegionID     string `json:"regionId"`
		ResourceType string `json:"resourceType"`
		CloudType    string `json:"cloudType"`
	} `json:"resource,omitempty"`
}

// Authenticate prisma Authenticate
type Authenticate struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	CustomerName string `json:"customerName"`
}

// LoginPrismaResponse struct
type LoginPrismaResponse struct {
	Token         string `json:"token"`
	Message       string `json:"message"`
	CustomerNames []struct {
		CustomerName string `json:"customerName"`
		TosAccepted  bool   `json:"tosAccepted"`
	} `json:"customerNames,omitempty"`
}

// LoginPrismaInput request parameter for Login
type LoginPrismaInput struct {
	Auth []byte
}

type RegisterAccountInput struct {
	Payload []byte
}

// Request accept an input as HTTP API request to call Prisma API
func (pc *PrismaClient) Request(request *PrismaAPIRequestInput) (io.ReadCloser, error) {
	if err := errors.FieldsVerifier(pc); err != nil {
		return nil, err
	}

	if request == nil {
		return nil, fmt.Errorf("PrismaAPIRequestInput is nil")
	}

	if request.Action == http.MethodGet {
		if err := errors.FieldsVerifier(request, "Payload"); err != nil {
			return nil, err
		}
	} else {
		if err := errors.FieldsVerifier(request); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("https://%s.prismacloud.io/%s", pc.Tenant, request.Endpoint)

	req, err := http.NewRequest(request.Action, url, bytes.NewBuffer(request.Payload))
	if err != nil {
		return nil, err
	}

	for key, value := range request.Header {
		req.Header.Add(key, value)
	}

	resp, resperr := pc.Do(req)
	if resperr != nil {
		return nil, resperr
	}

	if resp.StatusCode == http.StatusOK {
		return resp.Body, nil
	}

	return nil, fmt.Errorf("Response error: %d", resp.StatusCode)
}

// ListAlerts return a filered list of alerts
func (pc *PrismaClient) ListAlerts(listAlertInput *ListAlertsInput) (*Alerts, error) {
	if err := errors.FieldsVerifier(pc); err != nil {
		return nil, err
	}

	if listAlertInput == nil {
		return nil, errors.New("ListAlertsInput is nil")
	}

	if err := errors.FieldsVerifier(listAlertInput); err != nil {
		return nil, err
	}

	payload, err := json.Marshal(&listAlertInput.ListAlertsPayload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s.prismacloud.io/alert", pc.Tenant)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-redlock-auth", pc.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, resperr := pc.Do(req)
	if resperr != nil {
		return nil, resperr
	}
	alerts := &Alerts{}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(alerts); err != nil {
			return nil, err
		}
		return alerts, nil
	}
	return nil, fmt.Errorf("Unexpected error %d", resp.StatusCode)
}

// DismissAlerts accpet a DismissAlertInput which contain Alert ID
func (pc *PrismaClient) DismissAlerts(dismissAlertInput *DismissAlertInput) ([]byte, error) {
	if dismissAlertInput == nil {
		return nil, errors.New("DismissAlertInput is nil")
	}

	if err := errors.FieldsVerifier(dismissAlertInput, "DismissAlertFilter"); err != nil {
		return nil, err
	}
	if err := errors.FieldsVerifier(pc); err != nil {
		return nil, err
	}

	if dismissAlertInput.DismissAlertFilter.TimeRange == (FilterTimeRange{}) {
		dismissAlertInput.DismissAlertFilter.TimeRange = FilterTimeRange{
			Type: "relative",
			Value: TimeRangeValue{
				Amount: 1,
				Unit:   "week",
			},
		}
	}

	if len(dismissAlertInput.DismissAlertFilter.Filters) == 0 {
		dismissAlertInput.DismissAlertFilter.Filters = Filters{
			{
				Name:     "alert.status",
				Value:    "Open",
				Operator: "=",
			},
			{
				Name:     "policy.type",
				Value:    "audit_event",
				Operator: "=",
			},
		}
	}

	payload, err := json.Marshal(&dismissAlertInput)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(payload))

	url := fmt.Sprintf("https://%s.prismacloud.io/alert/dismiss", pc.Tenant)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-redlock-auth", pc.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, resperr := pc.Do(req)
	if resperr != nil {
		return nil, resperr
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, fmt.Errorf("Unexpected error %d", resp.StatusCode)
}

// ListAccountGroups return AccountGroups that contain group id and name
func (pc *PrismaClient) ListAccountGroups() (*AccountGroups, error) {
	if err := errors.FieldsVerifier(pc); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%s.prismacloud.io/cloud/group/name", pc.Tenant)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-redlock-auth", pc.Token)
	resp, err := pc.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	accountGroups := &AccountGroups{}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(accountGroups); err != nil {
			return nil, err
		}
		return accountGroups, nil
	}
	return nil, fmt.Errorf("Unexpected error %d", resp.StatusCode)
}

// LoginPrisma get Token from prisma and return the Token
func (pc *PrismaClient) LoginPrisma(reuqest *LoginPrismaInput) error {
	if err := errors.FieldsVerifier(pc, "Token"); err != nil {
		return err
	}

	if reuqest == nil {
		return errors.New("LoginPrismaInput is nil")
	}

	if err := errors.FieldsVerifier(reuqest); err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s.prismacloud.io/login", pc.Tenant)

	resp, err := pc.Post(url, "application/json", bytes.NewBuffer(reuqest.Auth))

	if err != nil {
		return err
	}
	loginResponse := &LoginPrismaResponse{}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(loginResponse); err != nil {
			return err
		}
		pc.Token = loginResponse.Token
		return nil
	}
	return fmt.Errorf("Unexpected error %d", resp.StatusCode)
}

// ListAccountNames returns a list of cloud account IDs and names.
func (pc *PrismaClient) ListAccountNames() (*AccountNames, error) {
	if err := errors.FieldsVerifier(pc, "Token"); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s.prismacloud.io/cloud/name", pc.Tenant)

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-redlock-auth", pc.Token)

	resp, respErr := pc.Do(req)

	if respErr != nil {
		fmt.Println(respErr.Error())
		return nil, respErr
	}

	accountNames := &AccountNames{}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(accountNames); err != nil {
			return nil, err
		}
		return accountNames, nil
	}

	return nil, fmt.Errorf("Unexpected error %d", resp.StatusCode)
}

// RegisterAccount register new account
func (pc *PrismaClient) RegisterAccount(payload []byte) error {
	if err := errors.FieldsVerifier(pc); err != nil {
		return err
	}

	if len(payload) == 0 {
		return fmt.Errorf("required parameter payload is empty")
	}

	url := fmt.Sprintf("https://%s.prismacloud.io/cloud/cloud_type", pc.Tenant)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-redlock-auth", pc.Token)
	resp, err := pc.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	return fmt.Errorf("Unexpected error %d", resp.StatusCode)
}

// DefaultHeader accept Token and return an http.request.header
func DefaultHeader(token string) map[string]string {
	return map[string]string{
		"Content-Type":   "application/json",
		"x-redlock-auth": token,
	}
}
