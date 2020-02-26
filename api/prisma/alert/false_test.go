package alert_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma/alert"
	"github.com/CityOfNewYork/prisma-cloud-remediation/events"
)

func TestFalseAWSRegionViolationAlert(t *testing.T) {
	testCases := []struct {
		name     string
		alert    *events.Alert
		expected bool
	}{
		{
			name: "Azure CloudType",
			alert: &events.Alert{
				ResourceRegionID: "us-east-1",
				ResourceID:       "DetachInternetGateway",
				CloudType:        "azure",
			},
			expected: false,
		},
		{
			name: "Virginia Region",
			alert: &events.Alert{
				ResourceRegionID: "us-east-1",
				ResourceID:       "DetachInternetGateway",
				CloudType:        "aws",
			},
			expected: false,
		},
		{
			name: "Not region violation ResourceID",
			alert: &events.Alert{
				ResourceRegionID: "us-east-2",
				ResourceID:       "CreateEC2",
				CloudType:        "aws",
			},
			expected: false,
		},
		{
			name: "False alert DetachInternetGateway alert in Non-Virginia region",
			alert: &events.Alert{
				ResourceRegionID: "us-east-2",
				ResourceID:       "DetachInternetGateway",
				CloudType:        "aws",
			},
			expected: true,
		},
		{
			name: "False alert DeleteInternetGateway alert in Non-Virginia region",
			alert: &events.Alert{
				ResourceRegionID: "us-west-1",
				ResourceID:       "DeleteInternetGateway",
				CloudType:        "aws",
			},
			expected: true,
		},
		{
			name: "False alert DeleteSubnet alert in Non-Virginia region",
			alert: &events.Alert{
				ResourceRegionID: "us-east-2",
				ResourceID:       "DeleteSubnet",
				CloudType:        "aws",
			},
			expected: true,
		},
		{
			name: "False alert DeleteVpc alert in Non-Virginia region",
			alert: &events.Alert{
				ResourceRegionID: "us-west-2",
				ResourceID:       "DeleteVpc",
				CloudType:        "aws",
			},
			expected: true,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s", i, testCase.name), func(t *testing.T) {
			resp := alert.FalseAWSRegionViolationAlert(testCase.alert)
			assert.Equal(t, testCase.expected, resp)
		})
	}
}
