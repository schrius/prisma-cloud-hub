package events

type OnBoardEvent struct {
	CloudType  string       `json:"cloudType"`
	GroupNames []string     `json:"groupNames"`
	AWS        AWSAccount   `json:"awsAccount"`
	Azure      AzureAccount `json:"azureAccount"`
	GCP        GCP          `json:"gcpAccount"`
}

type AWSAccount struct {
	AccountID  string   `json:"accountId"`
	Enabled    bool     `json:"enabled"`
	ExternalID string   `json:"externalId"`
	GroupIds   []string `json:"groupIds"`
	Name       string   `json:"name"`
	RoleArn    string   `json:"roleArn"`
}
type AzureAccount struct {
	CloudAccount       CloudAccount `json:"cloudAccount"`
	ClientID           string       `json:"clientId"`
	Key                string       `json:"key"`
	MonitorFlowLogs    bool         `json:"monitorFlowLogs"`
	TenantID           string       `json:"tenantId"`
	ServicePrincipalID string       `json:"servicePrincipalId"`
}

type CloudAccount struct {
	AccountID string   `json:"accountId"`
	Enabled   bool     `json:"enabled"`
	GroupIds  []string `json:"groupIds"`
	Name      string   `json:"name"`
}

type GCP struct {
	CloudAccount           CloudAccount `json:"cloudAccount"`
	CompressionEnabled     bool         `json:"compressionEnabled"`
	DataflowEnabledProject string       `json:"dataflowEnabledProject"`
	FlowLogStorageBucket   string       `json:"flowLogStorageBucket"`
	Credentials            struct {
		Type                    string `json:"type"`
		ProjectID               string `json:"project_id"`
		PrivateKeyID            string `json:"private_key_id"`
		PrivateKey              string `json:"private_key"`
		ClientEmail             string `json:"client_email"`
		ClientID                string `json:"client_id"`
		AuthURI                 string `json:"auth_uri"`
		TokenURI                string `json:"token_uri"`
		AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
		ClientX509CertURL       string `json:"client_x509_cert_url"`
	} `json:"credentials"`
}
