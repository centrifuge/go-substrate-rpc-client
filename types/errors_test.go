package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testDispatchError1 = DispatchError{
		IsOther: true,
	}
	testDispatchError2 = DispatchError{
		IsCannotLookup: true,
	}
	testDispatchError3 = DispatchError{
		IsBadOrigin: true,
	}
	testDispatchError4 = DispatchError{
		IsModule: true,
		ModuleError: ModuleError{
			Index: 4,
			Error: 5,
		},
	}
	testDispatchError5 = DispatchError{
		IsConsumerRemaining: true,
	}
	testDispatchError6 = DispatchError{
		IsNoProviders: true,
	}
	testDispatchError7 = DispatchError{
		IsTooManyConsumers: true,
	}
	testDispatchError8 = DispatchError{
		IsToken: true,
		TokenError: TokenError{
			IsUnsupported: true,
		},
	}
	testDispatchError9 = DispatchError{
		IsArithmetic: true,
		ArithmeticError: ArithmeticError{
			IsDivisionByZero: true,
		},
	}
	testDispatchError10 = DispatchError{
		IsTransactional: true,
		TransactionalError: TransactionalError{
			IsLimitReached: true,
		},
	}
)

func TestDispatchError_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testDispatchError1)
	assertRoundtrip(t, testDispatchError2)
	assertRoundtrip(t, testDispatchError3)
	assertRoundtrip(t, testDispatchError4)
	assertRoundtrip(t, testDispatchError5)
	assertRoundtrip(t, testDispatchError6)
	assertRoundtrip(t, testDispatchError7)
	assertRoundtrip(t, testDispatchError8)
	assertRoundtrip(t, testDispatchError9)
	assertRoundtrip(t, testDispatchError10)
}

func TestDispatchError_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testDispatchError1, MustHexDecodeString("0x00")},
		{testDispatchError2, MustHexDecodeString("0x01")},
		{testDispatchError3, MustHexDecodeString("0x02")},
		{testDispatchError4, MustHexDecodeString("0x030405")},
		{testDispatchError5, MustHexDecodeString("0x04")},
		{testDispatchError6, MustHexDecodeString("0x05")},
		{testDispatchError7, MustHexDecodeString("0x06")},
		{testDispatchError8, MustHexDecodeString("0x0706")},
		{testDispatchError9, MustHexDecodeString("0x0802")},
		{testDispatchError10, MustHexDecodeString("0x0900")},
	})
}

func TestDispatchError_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x00"), testDispatchError1},
		{MustHexDecodeString("0x01"), testDispatchError2},
		{MustHexDecodeString("0x02"), testDispatchError3},
		{MustHexDecodeString("0x030405"), testDispatchError4},
		{MustHexDecodeString("0x04"), testDispatchError5},
		{MustHexDecodeString("0x05"), testDispatchError6},
		{MustHexDecodeString("0x06"), testDispatchError7},
		{MustHexDecodeString("0x0706"), testDispatchError8},
		{MustHexDecodeString("0x0802"), testDispatchError9},
		{MustHexDecodeString("0x0900"), testDispatchError10},
	})
}
