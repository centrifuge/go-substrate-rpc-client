package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testDisputeLocation1 = DisputeLocation{
		IsLocal: true,
	}
	testDisputeLocation2 = DisputeLocation{
		IsRemote: true,
	}

	testDisputeResult1 = DisputeResult{
		IsValid: true,
	}

	testDisputeResult2 = DisputeResult{
		IsInvalid: true,
	}
)

func TestDisputeLocation_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testDisputeLocation1)
	assertRoundtrip(t, testDisputeLocation2)
}

func TestDisputeLocation_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testDisputeLocation1, MustHexDecodeString("0x00")},
		{testDisputeLocation2, MustHexDecodeString("0x01")},
	})
}

func TestDisputeLocation_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x00"), testDisputeLocation1},
		{MustHexDecodeString("0x01"), testDisputeLocation2},
	})
}

func TestDisputeResult_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testDisputeResult1)
	assertRoundtrip(t, testDisputeResult2)
}

func TestDisputeResult_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testDisputeResult1, MustHexDecodeString("0x00")},
		{testDisputeResult2, MustHexDecodeString("0x01")},
	})
}

func TestDisputeResult_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x00"), testDisputeResult1},
		{MustHexDecodeString("0x01"), testDisputeResult2},
	})
}
