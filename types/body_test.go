package types_test

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	bodyID1 = types.BodyID{
		IsUnit: true,
	}
	bodyID2 = types.BodyID{
		IsNamed: true,
		Body:    []types.U8{4},
	}
	bodyID3 = types.BodyID{
		IsIndex: true,
		Index:   6,
	}
	bodyID4 = types.BodyID{
		IsExecutive: true,
	}
	bodyID5 = types.BodyID{
		IsTechnical: true,
	}
	bodyID6 = types.BodyID{
		IsLegislative: true,
	}
	bodyID7 = types.BodyID{
		IsJudicial: true,
	}
)

func TestBodyID_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, bodyID1)
	assertRoundtrip(t, bodyID2)
	assertRoundtrip(t, bodyID3)
	assertRoundtrip(t, bodyID4)
	assertRoundtrip(t, bodyID5)
	assertRoundtrip(t, bodyID6)
	assertRoundtrip(t, bodyID7)
}

func TestBodyID_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{bodyID1, types.MustHexDecodeString("0x00")},
		{bodyID2, types.MustHexDecodeString("0x010404")},
		{bodyID3, types.MustHexDecodeString("0x0206000000")},
		{bodyID4, types.MustHexDecodeString("0x03")},
		{bodyID5, types.MustHexDecodeString("0x04")},
		{bodyID6, types.MustHexDecodeString("0x05")},
		{bodyID7, types.MustHexDecodeString("0x06")},
	})
}

func TestBodyID_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{types.MustHexDecodeString("0x00"), bodyID1},
		{types.MustHexDecodeString("0x010404"), bodyID2},
		{types.MustHexDecodeString("0x0206000000"), bodyID3},
		{types.MustHexDecodeString("0x03"), bodyID4},
		{types.MustHexDecodeString("0x04"), bodyID5},
		{types.MustHexDecodeString("0x05"), bodyID6},
		{types.MustHexDecodeString("0x06"), bodyID7},
	})
}

var (
	bodyPart1 = types.BodyPart{
		IsVoice: true,
	}
	bodyPart2 = types.BodyPart{
		IsMembers:    true,
		MembersCount: 3,
	}
	bodyPart3 = types.BodyPart{
		IsFraction:    true,
		FractionNom:   1,
		FractionDenom: 2,
	}
	bodyPart4 = types.BodyPart{
		IsAtLeastProportion:    true,
		AtLeastProportionNom:   3,
		AtLeastProportionDenom: 4,
	}
	bodyPart5 = types.BodyPart{
		IsMoreThanProportion:    true,
		MoreThanProportionNom:   6,
		MoreThanProportionDenom: 7,
	}
)

func TestBodyPart_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, bodyPart1)
	assertRoundtrip(t, bodyPart2)
	assertRoundtrip(t, bodyPart3)
	assertRoundtrip(t, bodyPart4)
	assertRoundtrip(t, bodyPart5)
}

func TestBodyPart_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{bodyPart1, types.MustHexDecodeString("0x00")},
		{bodyPart2, types.MustHexDecodeString("0x0103000000")},
		{bodyPart3, types.MustHexDecodeString("0x020100000002000000")},
		{bodyPart4, types.MustHexDecodeString("0x030300000004000000")},
		{bodyPart5, types.MustHexDecodeString("0x040600000007000000")},
	})
}

func TestBodyPart_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{types.MustHexDecodeString("0x00"), bodyPart1},
		{types.MustHexDecodeString("0x0103000000"), bodyPart2},
		{types.MustHexDecodeString("0x020100000002000000"), bodyPart3},
		{types.MustHexDecodeString("0x030300000004000000"), bodyPart4},
		{types.MustHexDecodeString("0x040600000007000000"), bodyPart5},
	})
}
