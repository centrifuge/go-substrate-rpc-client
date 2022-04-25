//go:build types_test

package test

import (
	"os"
	"path"
	"testing"

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
				if err := types.DecodeFromBytes(b, &metadata); err != nil {
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
