//go:build types_test

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

package test

import (
	"os"
	"path"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

//go:generate go run ./test-gen test_data meta_bytes storage_bytes

func TestTypesDecode(t *testing.T) {
	dirEntries, err := os.ReadDir("test_data")

	if err != nil {
		t.Errorf("Couldn't read test data dir - %s", err)
	}

	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}

		subdirPath := path.Join("test_data", dirEntry.Name())

		subdirEntries, err := os.ReadDir(subdirPath)

		if err != nil {
			t.Errorf("Couldn't read sub dir - %s", err)
		}

		var metadata types.Metadata
		var storageData []byte

		for _, subdirEntry := range subdirEntries {
			if path.Ext(subdirEntry.Name()) == ".txt" {
				// Skip info file
				continue
			}

			testDataPath := path.Join(subdirPath, subdirEntry.Name())

			b, err := os.ReadFile(testDataPath)

			if err != nil {
				t.Errorf("Couldn't read meta file - %s", err)
			}

			switch subdirEntry.Name() {
			case "meta_bytes":
				if err := codec.Decode(b, &metadata); err != nil {
					t.Errorf("Couldn't decode Metadata - %s", err)
				}
			case "storage_bytes":
				storageData = b
			}
		}

		events := types.EventRecords{}

		if err := types.EventRecordsRaw(storageData).DecodeEventRecords(&metadata, &events); err != nil {
			t.Errorf("Couldn't decode events for %s - %s", subdirPath, err)
		}
	}
}
