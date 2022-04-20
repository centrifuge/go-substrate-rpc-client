package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var testMultiLocation = MultiLocationV1{
	Parents: 1,
	Interior: JunctionsV1{
		IsHere: true,
	},
}

func TestOptionMultiLocation_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionMultiLocationV1(testMultiLocation))
	assertRoundtrip(t, NewOptionMultiLocationV1Empty())
}

func TestOptionMultiLocation_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewOptionMultiLocationV1(testMultiLocation), MustHexDecodeString("0x010100")},
		{NewOptionMultiLocationV1Empty(), MustHexDecodeString("0x00")},
	})
}

func TestOptionMultiLocation_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x010100"), NewOptionMultiLocationV1(testMultiLocation)},
		{MustHexDecodeString("0x00"), NewOptionMultiLocationV1Empty()},
	})
}
