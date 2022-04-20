package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var testExecutionResult = ExecutionResult{
	Outcome: 1,
	Error: XCMError{
		IsOverflow: true,
	},
}

func TestOptionExecutionResult_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionExecutionResult(testExecutionResult))
	assertRoundtrip(t, NewOptionExecutionResultEmpty())
}

func TestOptionExecutionResult_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewOptionExecutionResult(testExecutionResult), MustHexDecodeString("0x010100000000")},
		{NewOptionExecutionResultEmpty(), MustHexDecodeString("0x00")},
	})
}

func TestOptionExecutionResult_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x010100000000"), NewOptionExecutionResult(testExecutionResult)},
		{MustHexDecodeString("0x00"), NewOptionBytesEmpty()},
	})
}
