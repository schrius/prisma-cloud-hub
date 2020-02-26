package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CityOfNewYork/prisma-cloud-remediation/events"
)

func TestHandler(t *testing.T) {
	testCases := []struct {
		name     string
		event    events.Alert
		expected error
	}{
		{
			name: "false CloudType event",
			event: events.Alert{
				ResourceID:       "DetachInternetGateway",
				CloudType:        "azure",
				ResourceRegionID: "us-east-2",
			},
			expected: nil,
		},
		{
			name: "false Virginia event",
			event: events.Alert{
				ResourceID:       "DetachInternetGateway",
				CloudType:        "aws",
				ResourceRegionID: "us-east-1",
			},
			expected: nil,
		},
		{
			name: "false event",
			event: events.Alert{
				ResourceID:       "CreateEC2",
				CloudType:        "aws",
				ResourceRegionID: "us-east-2",
			},
			expected: nil,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase[%d] %s ", i, testCase.name), func(t *testing.T) {
			err := handler(nil, testCase.event)
			assert.NoError(t, err)
		})

	}
}
