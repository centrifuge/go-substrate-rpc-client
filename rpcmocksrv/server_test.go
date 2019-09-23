// Copyright 2018 Jsgenesis
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

package rpcmocksrv

import (
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

// const (
// 	SubKeySign = "sign-blob"

// 	// SubKeyCmd subkey command to create signatures
// 	SubKeyCmd = "/Users/vimukthi/.cargo/bin/subkey"
// )

// type AnchorParams struct {
// 	AnchorIDPreimage [32]byte
// 	DocRoot          [32]byte
// 	Proof            [32]byte
// }

// func NewRandomAnchor() AnchorParams {
// 	ap := AnchorParams{}
// 	copy(ap.AnchorIDPreimage[:], utils.RandomSlice(32))
// 	copy(ap.DocRoot[:], utils.RandomSlice(32))
// 	copy(ap.Proof[:], utils.RandomSlice(32))
// 	return ap
// }

// func (a *AnchorParams) Decode(decoder scale.Decoder) error {
// 	decoder.Read(a.AnchorIDPreimage[:])
// 	decoder.Read(a.DocRoot[:])
// 	decoder.Read(a.Proof[:])
// 	return nil
// }

// func (a AnchorParams) Encode(encoder scale.Encoder) error {
// 	encoder.Write(a.AnchorIDPreimage[:])
// 	encoder.Write(a.DocRoot[:])
// 	encoder.Write(a.Proof[:])
// 	return nil
// }

type TestService struct {
}

func (ts *TestService) Ping(s string) string {
	return s
}

func TestServer(t *testing.T) {
	s := New()

	ts := new(TestService)
	err := s.RegisterName("testserv3", ts)
	assert.NoError(t, err)

	c, err := rpc.Dial(s.URL)
	assert.NoError(t, err)

	var res string
	err = c.Call(&res, "testserv3_ping", "hello")
	assert.NoError(t, err)

	assert.Equal(t, "hello", res)
}
