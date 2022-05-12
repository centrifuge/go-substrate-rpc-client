package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testHRMPChannelID = HRMPChannelID{
		Sender:    11,
		Recipient: 45,
	}
)

func TestHRMPChannelID_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testHRMPChannelID)
}

func TestHRMPChannelID_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testHRMPChannelID, MustHexDecodeString("0x0b0000002d000000")},
	})
}

func TestHRMPChannelID_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x0b0000002d000000"), testHRMPChannelID},
	})
}
