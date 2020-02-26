package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type InnerField struct {
	InnerString string
}

type EmptyInnerField struct {
}

type TestField struct {
	TestString string
	TestMap map[string]string
	Inner InnerField
	EmptyInner EmptyInnerField
	Ptr *int
}

func TestFieldsVerifier(t *testing.T) {
	ptr := 1
	errorTestCases :=  []*TestField{
		&TestField{},
		&TestField{TestString: "test"},
		&TestField{TestString: "test", TestMap: map[string]string{"Test": "test"}},
		&TestField{TestString: "test", TestMap: map[string]string{"Test": "test"}, Inner: InnerField{"InnerString"}},
	}

	noErrorTestCase := &TestField{TestString: "test", TestMap: map[string]string{"Test": "test"}, Inner: InnerField{"InnerString"}, Ptr: &ptr}

	for _, testCase := range errorTestCases {
		assert.Error(t, FieldsVerifier(testCase))
	}

	assert.NoError(t, FieldsVerifier(noErrorTestCase))
}