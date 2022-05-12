package types_test

import (
	"math/big"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testClassDetails = ClassDetails{
		Owner:             NewAccountID([]byte("acc_id")),
		Issuer:            NewAccountID([]byte("acc_id2")),
		Admin:             NewAccountID([]byte("acc_id3")),
		Freezer:           NewAccountID([]byte("acc_id4")),
		TotalDeposit:      NewU128(*big.NewInt(123)),
		FreeHolding:       true,
		Instances:         4,
		InstanceMetadatas: 5,
		Attributes:        6,
		IsFrozen:          true,
	}
)

func TestClassDetails_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testClassDetails)
}

func TestClassDetails_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{
			testClassDetails,
			MustHexDecodeString("0x6163635f696400000000000000000000000000000000000000000000000000006163635f696432000000000000000000000000000000000000000000000000006163635f696433000000000000000000000000000000000000000000000000006163635f696434000000000000000000000000000000000000000000000000007b0000000000000000000000000000000104000000050000000600000001"),
		},
	})
}

func TestClassDetails_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{
			MustHexDecodeString("0x6163635f696400000000000000000000000000000000000000000000000000006163635f696432000000000000000000000000000000000000000000000000006163635f696433000000000000000000000000000000000000000000000000006163635f696434000000000000000000000000000000000000000000000000007b0000000000000000000000000000000104000000050000000600000001"),
			testClassDetails,
		},
	})
}
