package events_test

import (
	"encoding/json"
	"testing"

	"github.com/CityOfNewYork/prisma-cloud-remediation/events"
	"github.com/aws/aws-lambda-go/events/test"
	"github.com/stretchr/testify/assert"
)

func TestVPCKillerEventMarshaling(t *testing.T) {
	inputJSON := test.ReadJSONFromFile(t, "./testdata/vpc-killer-event.json")

	var inputEvent events.VPCAlert
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}

	outputJSON, err := json.Marshal(inputEvent)
	if err != nil {
		t.Errorf("could not marshal event. details: %v", err)
	}

	assert.JSONEq(t, string(inputJSON), string(outputJSON))
}

func TestVPCKillerEventMarshalingMalformedJson(t *testing.T) {
	test.TestMalformedJson(t, events.VPCAlert{})
}
