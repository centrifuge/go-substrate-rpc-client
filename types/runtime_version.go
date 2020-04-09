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

package types

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

type RuntimeVersion struct {
	APIs             []RuntimeVersionAPI `json:"apis"`
	AuthoringVersion U32                 `json:"authoringVersion"`
	ImplName         string              `json:"implName"`
	ImplVersion      U32                 `json:"implVersion"`
	SpecName         string              `json:"specName"`
	SpecVersion      U32                 `json:"specVersion"`
}

func NewRuntimeVersion() *RuntimeVersion {
	return &RuntimeVersion{APIs: make([]RuntimeVersionAPI, 0)}
}

func (r *RuntimeVersion) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&r.APIs)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.AuthoringVersion)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.ImplName)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.ImplVersion)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.SpecName)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.SpecVersion)
	if err != nil {
		return err
	}

	return nil
}

func (r RuntimeVersion) Encode(ctx context.Context, encoder scale.Encoder) error {
	err := encoder.Encode(ctx, r.APIs)
	if err != nil {
		return err
	}

	err = encoder.Encode(ctx, r.AuthoringVersion)
	if err != nil {
		return err
	}

	err = encoder.Encode(ctx, r.ImplName)
	if err != nil {
		return err
	}

	err = encoder.Encode(ctx, r.ImplVersion)
	if err != nil {
		return err
	}

	err = encoder.Encode(ctx, r.SpecName)
	if err != nil {
		return err
	}

	err = encoder.Encode(ctx, r.SpecVersion)
	if err != nil {
		return err
	}

	return nil
}

type RuntimeVersionAPI struct {
	APIID   string
	Version U32
}

func (r *RuntimeVersionAPI) UnmarshalJSON(b []byte) error {
	tmp := []interface{}{&r.APIID, &r.Version}
	wantLen := len(tmp)
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Notification: %d != %d", g, e)
	}
	return nil
}

func (r RuntimeVersionAPI) MarshalJSON() ([]byte, error) {
	tmp := []interface{}{r.APIID, r.Version}
	return json.Marshal(tmp)
}

func (r *RuntimeVersionAPI) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&r.APIID)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.Version)
	if err != nil {
		return err
	}

	return nil
}

func (r RuntimeVersionAPI) Encode(ctx context.Context, encoder scale.Encoder) error {
	err := encoder.Encode(ctx, r.APIID)
	if err != nil {
		return err
	}

	err = encoder.Encode(ctx, r.Version)
	if err != nil {
		return err
	}

	return nil
}
