package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/stretchr/testify/assert"
)

func TestMetadataV14_TestFindCallIndexWithUnknownFunction(t *testing.T) {
	var metadata Metadata
	err := DecodeFromHexString(MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.NoError(t, err)

	_, err = metadata.FindCallIndex("Module2_14.unknownFunction")
	assert.Error(t, err)
}

// Test that decoding the example metadata v14 works and that
// encoding it produces the original value.
func TestNewMetadataV14_Decode(t *testing.T) {
	// Verify that we can succcessfully decode metadata v14
	var metadata Metadata
	err := DecodeFromHexString(MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.NoError(t, err)

	// Verify that (encoding . decoding) equals the original input
	data, err := EncodeToBytes(metadata)
	assert.NoError(t, err)
	assert.Equal(t, MetadataV14Data, HexEncodeToString(data))
}

// TODO(nuno): make verifications more meaningful
func TestMetadataV14FindCallIndex(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	_, err = meta.FindCallIndex("Balances.transfer")
	assert.NoError(t, err)
}

// TODO(nuno): make verifications more meaningful
func TestMetadataV14FindEventNamesForEventID(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	id := EventID{}
	id[0] = 5
	id[1] = 2
	_, _, err = meta.FindEventNamesForEventID(id)
	assert.NoError(t, err)
}

// TODO(nuno): make verifications more meaningful
func TestMetadataV14FindStorageEntryMetadata(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	assert.NoError(t, err)

	_, err = meta.FindStorageEntryMetadata("System", "Account")
	assert.NoError(t, err)
}

func TestMetadataV14ExistsModuleMetadata(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	res := meta.ExistsModuleMetadata("System")
	assert.True(t, res)
}
