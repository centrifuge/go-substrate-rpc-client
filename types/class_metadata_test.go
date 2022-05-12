package types_test

import (
	"math/big"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testClassMetadata = ClassMetadata{
		Deposit:  NewU128(*big.NewInt(123)),
		Data:     Bytes("data"),
		IsFrozen: true,
	}
)

func TestClassMetadata_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testClassMetadata)
}

func TestClassMetadata_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testClassMetadata, MustHexDecodeString("0x7b000000000000000000000000000000106461746101")},
	})
}

func TestClassMetadata_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x7b000000000000000000000000000000106461746101"), testClassMetadata},
	})
}
