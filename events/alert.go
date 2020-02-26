package events

// Alert simple alert struct
type Alert struct {
	ResourceID       string `json:"resourceId"`
	AlertRuleName    string `json:"alertRuleName"`
	AccountName      string `json:"accountName"`
	ResourceRegionID string `json:"resourceRegionId"`
	CloudType        string `json:"cloudType"`
	AlertID          string `json:"alertId"`
	Severity         string `json:"severity"`
	PolicyName       string `json:"policyName"`
	ResourceName     string `json:"resourceName"`
	ResourceRegion   string `json:"resourceRegion"`
	AccountID        string `json:"accountId"`
	PolicyID         string `json:"policyId"`
}
