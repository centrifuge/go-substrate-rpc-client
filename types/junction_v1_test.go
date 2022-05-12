package types_test

import (
	"math/big"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testJunctionV1n1 = JunctionV1{
		IsParachain: true,
		ParachainID: NewUCompactFromUInt(11),
	}
	testJunctionV1n2 = JunctionV1{
		IsAccountId32: true,
		AccountId32NetworkID: NetworkID{
			IsAny: true,
		},
		AccountID: []U8{1, 2, 3},
	}
	testJunctionV1n3 = JunctionV1{
		IsAccountIndex64: true,
		AccountIndex64NetworkID: NetworkID{
			IsAny: true,
		},
		AccountIndex: 16,
	}
	testJunctionV1n4 = JunctionV1{
		IsAccountKey20: true,
		AccountKey20NetworkID: NetworkID{
			IsKusama: true,
		},
	}
	testJunctionV1n5 = JunctionV1{
		IsPalletInstance: true,
		PalletIndex:      4,
	}
	testJunctionV1n6 = JunctionV1{
		IsGeneralIndex: true,
		GeneralIndex:   NewU128(*big.NewInt(42)),
	}
	testJunctionV1n7 = JunctionV1{
		IsGeneralKey: true,
		GeneralKey:   []U8{6, 8},
	}
	testJunctionV1n8 = JunctionV1{
		IsOnlyChild: true,
	}
	testJunctionV1n9 = JunctionV1{
		IsPlurality: true,
		BodyID: BodyID{
			IsUnit: true,
		},
		BodyPart: BodyPart{
			IsVoice: true,
		},
	}
)

func TestJunctionV1_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testJunctionV1n1)
	assertRoundtrip(t, testJunctionV1n2)
	assertRoundtrip(t, testJunctionV1n3)
	assertRoundtrip(t, testJunctionV1n4)
	assertRoundtrip(t, testJunctionV1n5)
	assertRoundtrip(t, testJunctionV1n6)
	assertRoundtrip(t, testJunctionV1n7)
	assertRoundtrip(t, testJunctionV1n8)
	assertRoundtrip(t, testJunctionV1n9)
}

func TestJunctionV1_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testJunctionV1n1, MustHexDecodeString("0x002c")},
		{testJunctionV1n2, MustHexDecodeString("0x01000c010203")},
		{testJunctionV1n3, MustHexDecodeString("0x02001000000000000000")},
		{testJunctionV1n4, MustHexDecodeString("0x030300")},
		{testJunctionV1n5, MustHexDecodeString("0x0404")},
		{testJunctionV1n6, MustHexDecodeString("0x052a000000000000000000000000000000")},
		{testJunctionV1n7, MustHexDecodeString("0x06080608")},
		{testJunctionV1n8, MustHexDecodeString("0x07")},
		{testJunctionV1n9, MustHexDecodeString("0x080000")},
	})
}

func TestJunctionV1_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x002c"), testJunctionV1n1},
		{MustHexDecodeString("0x01000c010203"), testJunctionV1n2},
		{MustHexDecodeString("0x02001000000000000000"), testJunctionV1n3},
		{MustHexDecodeString("0x030300"), testJunctionV1n4},
		{MustHexDecodeString("0x0404"), testJunctionV1n5},
		{MustHexDecodeString("0x052a000000000000000000000000000000"), testJunctionV1n6},
		{MustHexDecodeString("0x06080608"), testJunctionV1n7},
		{MustHexDecodeString("0x07"), testJunctionV1n8},
		{MustHexDecodeString("0x080000"), testJunctionV1n9},
	})
}
