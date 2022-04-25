package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_ParseField(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		field             reflect.StructField
		expectedFieldInfo *FieldInfo
	}{
		{
			field: reflect.StructField{
				Name: "Test_field",
				Tag:  reflect.StructTag(fmt.Sprintf(`%s:"ws://test" %s:"http://test"`, WsURLTagName, ApiURLTagName)),
			},
			expectedFieldInfo: &FieldInfo{
				ReqData: &ReqData{
					Module: "test",
					Call:   "field",
				},
				ClientOpts: &ClientOpts{
					ApiURL: "http://test",
					WsURL:  "ws://test",
				},
			},
		},
		{
			field: reflect.StructField{
				Name: "Test_field",
				Tag:  reflect.StructTag(fmt.Sprintf(`%s:"blockchainname"`, BlockchainTagName)),
			},
			expectedFieldInfo: &FieldInfo{
				ReqData: &ReqData{
					Module: "test",
					Call:   "field",
				},
				ClientOpts: &ClientOpts{
					Blockchain: "blockchainname",
					ApiURL:     fmt.Sprintf(ApiURLFormat, "blockchainname"),
					WsURL:      fmt.Sprintf(WsURLFormat, "blockchainname"),
				},
			},
		},
	}

	for _, test := range tests {
		r, err := parser.ParseField(test.field)
		assert.Nil(t, err)

		assert.EqualValues(t, r, test.expectedFieldInfo)
	}
}

func TestParser_CanSkip(t *testing.T) {
	parser := NewParser()

	r := parser.CanSkip(reflect.StructField{
		Tag: reflect.StructTag(fmt.Sprintf(`%s:"true"`, SkipTagName)),
	})

	assert.True(t, r)

	r = parser.CanSkip(reflect.StructField{})

	assert.False(t, r)
}

func TestParser_parseClientOpts(t *testing.T) {
	parser := NewParser()

	field := reflect.StructField{
		Tag: reflect.StructTag(fmt.Sprintf(`%s:"ws://test" %s:"http://test"`, WsURLTagName, ApiURLTagName)),
	}

	_, err := parser.parseClientOpts(field)

	assert.Nil(t, err)

	field = reflect.StructField{
		Tag: reflect.StructTag(fmt.Sprintf(`%s:"blockchainname""`, BlockchainTagName)),
	}

	_, err = parser.parseClientOpts(field)

	assert.Nil(t, err)

	field = reflect.StructField{
		Tag: reflect.StructTag(fmt.Sprintf(`%s:"ws://test"`, WsURLTagName)),
	}

	_, err = parser.parseClientOpts(field)

	assert.True(t, errors.Is(err, errApiURLMissing))

	field = reflect.StructField{
		Tag: reflect.StructTag(fmt.Sprintf(`%s:"http://test"`, ApiURLTagName)),
	}

	_, err = parser.parseClientOpts(field)

	assert.True(t, errors.Is(err, errWsURLMissing))

	field = reflect.StructField{
		Tag: "",
	}

	_, err = parser.parseClientOpts(field)

	assert.True(t, errors.Is(err, errBlockchainMissing))
}

func TestParser_parseReqData(t *testing.T) {
	parser := NewParser()

	field := reflect.StructField{
		Name: "Test_Field",
	}

	_, err := parser.parseReqData(field)

	assert.Nil(t, err)

	field = reflect.StructField{
		Name: "TestField",
	}

	_, err = parser.parseReqData(field)

	assert.NotNil(t, err)

	field = reflect.StructField{
		Name: "Test-Field",
	}

	_, err = parser.parseReqData(field)

	assert.NotNil(t, err)

	field = reflect.StructField{
		Name: "Test_Field_1",
	}

	_, err = parser.parseReqData(field)

	assert.NotNil(t, err)

}
