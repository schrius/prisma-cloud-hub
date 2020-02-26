package events

// PrismaEvent SQS  for json marshal
type PrismaEvent struct {
	ResourceID           string         `json:"resourceId"`
	AlertRuleName        string         `json:"alertRuleName"`
	AccountName          string         `json:"accountName"`
	HasFinding           bool           `json:"hasFinding"`
	ResourceRegionID     string         `json:"resourceRegionId"`
	AlertRemediationCli  interface{}    `json:"alertRemediationCli"`
	Source               string         `json:"source"`
	CloudType            string         `json:"cloudType"`
	ComplianceMetadata   interface{}    `json:"complianceMetadata"`
	CallbackURL          string         `json:"callbackUrl"`
	AlertID              string         `json:"alertId"`
	PolicyLabels         []string       `json:"policyLabels"`
	AlertAttribution     interface{}    `json:"alertAttribution"`
	Severity             string         `json:"severity"`
	PolicyName           string         `json:"policyName"`
	Resource             PrismaResource `json:"resource"`
	ResourceName         string         `json:"resourceName"`
	RiskRating           string         `json:"riskRating"`
	ResourceRegion       string         `json:"resourceRegion"`
	PolicyDescription    string         `json:"policyDescription"`
	PolicyRecommendation string         `json:"policyRecommendation"`
	AccountID            string         `json:"accountId"`
	ResourceConfig       interface{}    `json:"resourceConfig"`
	PolicyID             string         `json:"policyId"`
	ResourceCloudService string         `json:"resourceCloudService"`
	AlertTs              int64          `json:"alertTs"`
	FindingSummary       interface{}    `json:"findingSummary"`
	ResourceType         string         `json:"resourceType"`
}

// PrismaResource struct
type PrismaResource struct {
	Data            interface{} `json:"data"`
	URL             interface{} `json:"url"`
	Rrn             string      `json:"rrn"`
	AccountID       string      `json:"accountId"`
	RegionID        string      `json:"regionId"`
	CloudType       string      `json:"cloudType"`
	ResourceAPIName interface{} `json:"resourceApiName"`
	Name            string      `json:"name"`
	AdditionalInfo  interface{} `json:"additionalInfo"`
	ID              string      `json:"id"`
	Region          string      `json:"region"`
	Account         string      `json:"account"`
	ResourceType    string      `json:"resourceType"`
}
