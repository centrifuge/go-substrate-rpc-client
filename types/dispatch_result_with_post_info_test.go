package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testDispatchResultWithPostInfo1 = DispatchResultWithPostInfo{
		IsOk: true,
		Ok: PostDispatchInfo{
			ActualWeight: NewOptionWeight(123),
			PaysFee: Pays{
				IsYes: true,
			},
		},
	}
	testDispatchResultWithPostInfo2 = DispatchResultWithPostInfo{
		IsError: true,
		Error: DispatchErrorWithPostInfo{
			PostInfo: PostDispatchInfo{
				ActualWeight: NewOptionWeight(456),
				PaysFee: Pays{
					IsNo: true,
				},
			},
			Error: DispatchError{
				IsOther: true,
			},
		},
	}
)

func TestDispatchResultWithPostInfo_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testDispatchResultWithPostInfo1)
	assertRoundtrip(t, testDispatchResultWithPostInfo2)
}

func TestDispatchResultWithPostInfo_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testDispatchResultWithPostInfo1, MustHexDecodeString("0x00017b0000000000000000")},
		{testDispatchResultWithPostInfo2, MustHexDecodeString("0x0101c8010000000000000100")},
	})
}

func TestDispatchResultWithPostInfo_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x00017b0000000000000000"), testDispatchResultWithPostInfo1},
		{MustHexDecodeString("0x0101c8010000000000000100"), testDispatchResultWithPostInfo2},
	})
}
