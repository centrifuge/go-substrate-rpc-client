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

package types_test

import (
	"bytes"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestBalanceStatusEncodeDecode(t *testing.T) {
	// encode
	bs := Reserved
	var buf bytes.Buffer
	encoder := scale.NewEncoder(&buf)
	assert.NoError(t, encoder.Encode(bs))
	assert.Equal(t, buf.Len(), 1)
	assert.Equal(t, buf.Bytes(), []byte{1})

	//decode
	decoder := scale.NewDecoder(bytes.NewReader(buf.Bytes()))
	bs0 := BalanceStatus(0)
	err := decoder.Decode(&bs0)
	assert.NoError(t, err)
	assert.Equal(t, bs0, Reserved)

	//decode error
	decoder = scale.NewDecoder(bytes.NewReader([]byte{5}))
	bs0 = BalanceStatus(0)
	err = decoder.Decode(&bs0)
	assert.Error(t, err)
}
