package types_test

import (
	"math/big"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testInstanceMetadata = InstanceMetadata{
		Deposit:  NewU128(*big.NewInt(1234)),
		Data:     Bytes("some_data"),
		IsFrozen: true,
	}
)

func TestInstanceMetadata_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testInstanceMetadata)
}

func TestInstanceMetadata_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testInstanceMetadata, MustHexDecodeString("0xd204000000000000000000000000000024736f6d655f6461746101")},
	})
}

func TestInstanceMetadata_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0xd204000000000000000000000000000024736f6d655f6461746101"), testInstanceMetadata},
	})
}
