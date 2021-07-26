// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
// Copyright 2021 Snowfork
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
	"encoding/json"
	"fmt"
)

// StorageChangeSet contains changes from storage subscriptions
type StorageChangeSet struct {
	Block   Hash             `json:"block"`
	Changes []KeyValueOption `json:"changes"`
}

type KeyValueOption struct {
	StorageKey  StorageKey
	StorageData OptionStorageData
}

type OptionStorageData struct {
	option
	value StorageDataRaw
}

func NewOptionStorageData(value StorageDataRaw) OptionStorageData {
	return OptionStorageData{option{true}, value}
}

func NewOptionStorageDataEmpty() OptionStorageData {
	return OptionStorageData{option{false}, nil}
}

func (o OptionStorageData) Unwrap() (ok bool, value StorageDataRaw) {
	return o.hasValue, o.value
}

func (r *KeyValueOption) UnmarshalJSON(b []byte) error {
	var tmp []string
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	switch len(tmp) {
	case 2:
		key, err := HexDecodeString(tmp[0])
		if err != nil {
			return err
		}
		r.StorageKey = key

		if tmp[1] != "" {
			r.StorageData.hasValue = true
			data, err := HexDecodeString(tmp[1])
			if err != nil {
				return err
			}
			r.StorageData.value = data
		} else {
			r.StorageData.hasValue = false
		}
	default:
		return fmt.Errorf("expected 2 entries for StorageChange, got %v", len(tmp))
	}
	return nil
}

func (r KeyValueOption) MarshalJSON() ([]byte, error) {
	var tmp []interface{}
	tmp = append(tmp, r.StorageKey.Hex())
	if r.StorageData.hasValue {
		tmp = append(tmp, r.StorageData.value.Hex())
	} else {
		tmp = append(tmp, nil)
	}
	return json.Marshal(tmp)
}
