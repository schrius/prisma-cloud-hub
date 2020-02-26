package alert

import (
	"github.com/CityOfNewYork/prisma-cloud-remediation/events"
)

const (
	DetachInternetGateway = "DetachInternetGateway"
	DeleteInternetGateway = "DeleteInternetGateway"
	DeleteSubnet          = "DeleteSubnet"
	DeleteVpc             = "DeleteVpc"
	Virginia              = "us-east-1"
)

// FalseAlert
// return true if CloudType is not aws,
// non-Virginia region, ResourceID is
// DetachInternetGateway, DeleteInternetGateway,
// DeleteSubnet, DeleteVpc
func FalseAWSRegionViolationAlert(alert *events.Alert) bool {
	if alert.ResourceRegionID == Virginia {
		return false
	}
	if alert.CloudType != "aws" {
		return false
	}
	switch alert.ResourceID {
	case DetachInternetGateway, DeleteInternetGateway, DeleteSubnet, DeleteVpc:
		return true
	}
	return false
}
