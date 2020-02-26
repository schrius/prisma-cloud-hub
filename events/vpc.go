package events

type VPCAlert struct {
	ResourceID           string                `json:"resourceId"`
	AlertRuleName        string                `json:"alertRuleName"`
	AccountName          string                `json:"accountName"`
	HasFinding           bool                  `json:"hasFinding"`
	ResourceRegionID     string                `json:"resourceRegionId"`
	AlertRemediationCli  interface{}           `json:"alertRemediationCli"`
	Source               string                `json:"source"`
	CloudType            string                `json:"cloudType"`
	ComplianceMetadata   interface{}           `json:"complianceMetadata"`
	CallbackURL          string                `json:"callbackUrl"`
	AlertID              string                `json:"alertId"`
	PolicyLabels         []interface{}         `json:"policyLabels"`
	AlertAttribution     interface{}           `json:"alertAttribution"`
	Severity             string                `json:"severity"`
	PolicyName           string                `json:"policyName"`
	Resource             resource              `json:"resource"`
	ResourceName         string                `json:"resourceName"`
	RiskRating           string                `json:"riskRating"`
	ResourceRegion       string                `json:"resourceRegion"`
	PolicyDescription    string                `json:"policyDescription"`
	PolicyRecommendation string                `json:"policyRecommendation"`
	AccountID            string                `json:"accountId"`
	ResourceConfig       resourceConfiguration `json:"resourceConfig"`
	PolicyID             string                `json:"policyId"`
	ResourceCloudService string                `json:"resourceCloudService"`
	AlertTs              int64                 `json:"alertTs"`
	FindingSummary       interface{}           `json:"findingSummary"`
	ResourceType         string                `json:"resourceType"`
}

type resourceConfiguration struct {
	InstanceTenancy         string `json:"instanceTenancy"`
	CidrBlock               string `json:"cidrBlock"`
	CidrBlockAssociationSet []struct {
		CidrBlock      string `json:"cidrBlock"`
		CidrBlockState struct {
			State string `json:"state"`
		} `json:"cidrBlockState"`
		AssociationID string `json:"associationId"`
	} `json:"cidrBlockAssociationSet"`
	OwnerID string `json:"ownerId"`
	Tags    []struct {
		Value string `json:"value"`
		Key   string `json:"key"`
	} `json:"tags"`
	Default                     bool          `json:"default"`
	IsDefault                   bool          `json:"isDefault"`
	DhcpOptionsID               string        `json:"dhcpOptionsId"`
	VpcID                       string        `json:"vpcId"`
	State                       string        `json:"state"`
	SecurityGroupCount          int           `json:"securityGroupCount"`
	SubnetCount                 int           `json:"subnetCount"`
	Ipv6CidrBlockAssociationSet []interface{} `json:"ipv6CidrBlockAssociationSet"`
}

type resource struct {
	Data            data        `json:"data"`
	URL             string      `json:"url"`
	Rrn             string      `json:"rrn"`
	AccountID       string      `json:"accountId"`
	RegionID        string      `json:"regionId"`
	CloudType       string      `json:"cloudType"`
	ResourceAPIName string      `json:"resourceApiName"`
	Name            string      `json:"name"`
	AdditionalInfo  interface{} `json:"additionalInfo"`
	ID              string      `json:"id"`
	Region          string      `json:"region"`
	Account         string      `json:"account"`
	ResourceType    string      `json:"resourceType"`
}

type data struct {
	InstanceTenancy         string `json:"instanceTenancy"`
	CidrBlock               string `json:"cidrBlock"`
	CidrBlockAssociationSet []struct {
		CidrBlock      string `json:"cidrBlock"`
		CidrBlockState struct {
			State string `json:"state"`
		} `json:"cidrBlockState"`
		AssociationID string `json:"associationId"`
	} `json:"cidrBlockAssociationSet"`
	OwnerID string `json:"ownerId"`
	Tags    []struct {
		Value string `json:"value"`
		Key   string `json:"key"`
	} `json:"tags"`
	Default                     bool          `json:"default"`
	IsDefault                   bool          `json:"isDefault"`
	DhcpOptionsID               string        `json:"dhcpOptionsId"`
	VpcID                       string        `json:"vpcId"`
	State                       string        `json:"state"`
	SecurityGroupCount          int           `json:"securityGroupCount"`
	SubnetCount                 int           `json:"subnetCount"`
	Ipv6CidrBlockAssociationSet []interface{} `json:"ipv6CidrBlockAssociationSet"`
}
