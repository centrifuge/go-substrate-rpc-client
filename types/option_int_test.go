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
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
)

func TestOptionI8_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI8(NewI8(7)))
	assertRoundtrip(t, NewOptionI8(NewI8(0)))
	assertRoundtrip(t, NewOptionI8Empty())
}

func TestOptionI16_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI16(NewI16(14)))
	assertRoundtrip(t, NewOptionI16(NewI16(0)))
	assertRoundtrip(t, NewOptionI16Empty())
}

func TestOptionI32_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI32(NewI32(21)))
	assertRoundtrip(t, NewOptionI32(NewI32(0)))
	assertRoundtrip(t, NewOptionI32Empty())
}

func TestOptionI64_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI64(NewI64(28)))
	assertRoundtrip(t, NewOptionI64(NewI64(0)))
	assertRoundtrip(t, NewOptionI64Empty())
}
