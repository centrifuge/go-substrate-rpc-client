package types_test

import (
	"math/big"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testJunction1 = JunctionV0{
		IsParent: true,
	}
	testJunction2 = JunctionV0{
		IsParachain: true,
		ParachainID: 11,
	}
	testJunction3 = JunctionV0{
		IsAccountId32: true,
		AccountId32NetworkID: NetworkID{
			IsAny: true,
		},
		AccountID: []U8{1, 2, 3},
	}
	testJunction4 = JunctionV0{
		IsAccountIndex64: true,
		AccountIndex64NetworkID: NetworkID{
			IsAny: true,
		},
		AccountIndex: 16,
	}
	testJunction5 = JunctionV0{
		IsAccountKey20: true,
		AccountKey20NetworkID: NetworkID{
			IsKusama: true,
		},
	}
	testJunction6 = JunctionV0{
		IsPalletInstance: true,
		PalletIndex:      4,
	}
	testJunction7 = JunctionV0{
		IsGeneralIndex: true,
		GeneralIndex:   NewU128(*big.NewInt(42)),
	}
	testJunction8 = JunctionV0{
		IsGeneralKey: true,
		GeneralKey:   []U8{6, 8},
	}
	testJunction9 = JunctionV0{
		IsOnlyChild: true,
	}
	testJunction10 = JunctionV0{
		IsPlurality: true,
		PluralityID: BodyID{
			IsUnit: true,
		},
		PluralityPart: BodyPart{
			IsVoice: true,
		},
	}
)

func TestJunctionV0_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testJunction1)
	assertRoundtrip(t, testJunction2)
	assertRoundtrip(t, testJunction3)
	assertRoundtrip(t, testJunction4)
	assertRoundtrip(t, testJunction5)
	assertRoundtrip(t, testJunction6)
	assertRoundtrip(t, testJunction7)
	assertRoundtrip(t, testJunction8)
	assertRoundtrip(t, testJunction9)
	assertRoundtrip(t, testJunction10)
}

func TestJunctionV0_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testJunction1, MustHexDecodeString("0x00")},
		{testJunction2, MustHexDecodeString("0x010b000000")},
		{testJunction3, MustHexDecodeString("0x02000c010203")},
		{testJunction4, MustHexDecodeString("0x03001000000000000000")},
		{testJunction5, MustHexDecodeString("0x040300")},
		{testJunction6, MustHexDecodeString("0x0504")},
		{testJunction7, MustHexDecodeString("0x062a000000000000000000000000000000")},
		{testJunction8, MustHexDecodeString("0x07080608")},
		{testJunction9, MustHexDecodeString("0x08")},
		{testJunction10, MustHexDecodeString("0x090000")},
	})
}

func TestJunctionV0_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x00"), testJunction1},
		{MustHexDecodeString("0x010b000000"), testJunction2},
		{MustHexDecodeString("0x02000c010203"), testJunction3},
		{MustHexDecodeString("0x03001000000000000000"), testJunction4},
		{MustHexDecodeString("0x040300"), testJunction5},
		{MustHexDecodeString("0x0504"), testJunction6},
		{MustHexDecodeString("0x062a000000000000000000000000000000"), testJunction7},
		{MustHexDecodeString("0x07080608"), testJunction8},
		{MustHexDecodeString("0x08"), testJunction9},
		{MustHexDecodeString("0x090000"), testJunction10},
	})
}
