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

package author

import (
	"bytes"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PendingExtrinsics returns all pending extrinsics, potentially grouped by sender
func (a *author) PendingExtrinsics(meta *types.Metadata) ([]*registry.DecodedExtrinsic, error) {
	var hexEncodedExtrinsics []string
	err := a.client.Call(&hexEncodedExtrinsics, "author_pendingExtrinsics")
	if err != nil {
		return nil, err
	}

	extrinsicDecoder, err := registry.NewFactory().CreateExtrinsicDecoder(meta)

	if err != nil {
		return nil, err
	}

	var decodedExtrinsics []*registry.DecodedExtrinsic

	for _, hexEncodedExtrinsic := range hexEncodedExtrinsics {
		extrinsicBytes, err := hexutil.Decode(hexEncodedExtrinsic)

		if err != nil {
			panic(err)
		}

		scaleDecoder := scale.NewDecoder(bytes.NewReader(extrinsicBytes))

		decodedExtrinsic, err := extrinsicDecoder.Decode(scaleDecoder)

		if err != nil {
			return nil, err
		}

		decodedExtrinsics = append(decodedExtrinsics, decodedExtrinsic)
	}

	return decodedExtrinsics, err
}
