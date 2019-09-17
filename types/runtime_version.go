// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Centrifuge GmbH
//
// This file is part of Go Substrate RPC Client (GSRPC).
//
// GSRPC is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// GSRPC is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package types

import (
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

func (r RuntimeVersion) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(r.APIs)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.AuthoringVersion)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.ImplName)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.ImplVersion)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.SpecName)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.SpecVersion)
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

func (r *RuntimeVersionAPI) MarshalJSON() ([]byte, error) {
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

func (r RuntimeVersionAPI) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(r.APIID)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.Version)
	if err != nil {
		return err
	}

	return nil
}
