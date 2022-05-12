package types_test

import (
	"math/big"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testInstanceDetails = InstanceDetails{
		Owner:    NewAccountID([]byte("acc_id")),
		Approved: NewOptionAccountID(NewAccountID([]byte("acc_id2"))),
		IsFrozen: true,
		Deposit:  NewU128(*big.NewInt(123)),
	}
)

func TestInstanceDetails_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testInstanceDetails)
}

func TestInstanceDetails_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{
			testInstanceDetails,
			MustHexDecodeString("0x6163635f69640000000000000000000000000000000000000000000000000000016163635f69643200000000000000000000000000000000000000000000000000017b000000000000000000000000000000"),
		},
	})
}

func TestInstanceDetails_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{
			MustHexDecodeString("0x6163635f69640000000000000000000000000000000000000000000000000000016163635f69643200000000000000000000000000000000000000000000000000017b000000000000000000000000000000"),
			testInstanceDetails,
		},
	})
}
