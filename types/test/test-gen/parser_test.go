// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
				Tag:  reflect.StructTag(fmt.Sprintf(`%s:"ws://test" %s:"http://test"`, WsURLTagName, APIURLTagName)),
			},
			expectedFieldInfo: &FieldInfo{
				ReqData: &ReqData{
					Module: "test",
					Call:   "field",
				},
				ClientOpts: &ClientOpts{
					APIURL: "http://test",
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
					APIURL:     fmt.Sprintf(APIURLFormat, "blockchainname"),
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

func TestParser_ParseField_ParseClientOptsErr(t *testing.T) {
	parser := NewParser()

	f := reflect.StructField{
		Tag: reflect.StructTag(WsURLTagName),
	}

	r, err := parser.ParseField(f)
	assert.NotNil(t, err)
	assert.Nil(t, r)

}

func TestParser_ParseField_ParseReqDataErr(t *testing.T) {
	parser := NewParser()

	f := reflect.StructField{
		Name: "incorrect_field_name",
		Tag:  reflect.StructTag(fmt.Sprintf(`%s:"blockchainname"`, BlockchainTagName)),
	}

	r, err := parser.ParseField(f)
	assert.NotNil(t, err)
	assert.Nil(t, r)

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
		Tag: reflect.StructTag(fmt.Sprintf(`%s:"ws://test" %s:"http://test"`, WsURLTagName, APIURLTagName)),
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

	assert.True(t, errors.Is(err, errAPIURLMissing))

	field = reflect.StructField{
		Tag: reflect.StructTag(fmt.Sprintf(`%s:"http://test"`, APIURLTagName)),
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
