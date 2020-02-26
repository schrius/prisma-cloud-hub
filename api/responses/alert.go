package response

import (
	"time"
)

// Alert Prisma simple alert response
type Alert []struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Policy struct {
		PolicyID   string `json:"policyId"`
		PolicyType string `json:"policyType"`
	} `json:"policy"`
	Resource struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Account      string `json:"account"`
		AccountID    string `json:"accountId"`
		Region       string `json:"region"`
		RegionID     string `json:"regionId"`
		ResourceType string `json:"resourceType"`
		CloudType    string `json:"cloudType"`
	} `json:"resource"`
}

// AWSAlert Prisma AWS alert response
type AWSAlert []struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	FirstSeen     int64  `json:"firstSeen"`
	LastSeen      int64  `json:"lastSeen"`
	AlertTime     int64  `json:"alertTime"`
	EventOccurred int64  `json:"eventOccurred"`
	Policy        struct {
		PolicyID      string `json:"policyId"`
		PolicyType    string `json:"policyType"`
		SystemDefault bool   `json:"systemDefault"`
		Remediable    bool   `json:"remediable"`
	} `json:"policy"`
	RiskDetail struct {
		RiskScore struct {
			Score    int `json:"score"`
			MaxScore int `json:"maxScore"`
		} `json:"riskScore"`
		Rating string `json:"rating"`
		Score  string `json:"score"`
	} `json:"riskDetail"`
	Resource struct {
		ID                 string        `json:"id"`
		Name               string        `json:"name"`
		Account            string        `json:"account"`
		AccountID          string        `json:"accountId"`
		CloudAccountGroups []interface{} `json:"cloudAccountGroups"`
		Region             string        `json:"region"`
		RegionID           string        `json:"regionId"`
		ResourceType       string        `json:"resourceType"`
		Data               struct {
			EventVersion string `json:"eventVersion"`
			UserIdentity struct {
				Type        string `json:"type"`
				PrincipalID string `json:"principalId"`
				Arn         string `json:"arn"`
				AccountID   string `json:"accountId"`
			} `json:"userIdentity"`
			EventTime         time.Time   `json:"eventTime"`
			EventSource       string      `json:"eventSource"`
			EventName         string      `json:"eventName"`
			AwsRegion         string      `json:"awsRegion"`
			SourceIPAddress   string      `json:"sourceIPAddress"`
			UserAgent         string      `json:"userAgent"`
			RequestParameters interface{} `json:"requestParameters"`
			ResponseElements  struct {
				RenewRole string `json:"RenewRole"`
			} `json:"responseElements"`
			AdditionalEventData struct {
				RenewedBy  string `json:"RenewedBy"`
				RedirectTo string `json:"RedirectTo"`
			} `json:"additionalEventData"`
			EventID            string `json:"eventID"`
			EventType          string `json:"eventType"`
			RecipientAccountID string `json:"recipientAccountId"`
		} `json:"data"`
		CloudType string `json:"cloudType"`
	} `json:"resource"`
	TriggeredBy        string `json:"triggeredBy"`
	InvestigateOptions struct {
		SearchID string `json:"searchId"`
		StartTs  int64  `json:"startTs"`
		EndTs    int64  `json:"endTs"`
	} `json:"investigateOptions"`
}
