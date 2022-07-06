package types_test

import (
	"math/big"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	proxyDefinition1 = ProxyDefinition{
		Delegate:  NewAccountID([]byte("some-account")),
		ProxyType: ProxyTypePrice,
		Delay:     3,
	}
	proxyDefinition2 = ProxyDefinition{
		Delegate:  NewAccountID([]byte("some-other-account")),
		ProxyType: KeystoreManagement,
		Delay:     0,
	}
	proxyDefinitionFuzzOpts = proxyTypeFuzzOpts
)

func TestProxyDefinition_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[ProxyDefinition](t, 1000, proxyDefinitionFuzzOpts...)
	assertDecodeNilData[ProxyDefinition](t)
	assertEncodeEmptyObj[ProxyDefinition](t, 34)
}

func TestProxyDefinition_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{
			proxyDefinition1,
			MustHexDecodeString("0x736f6d652d6163636f756e740000000000000000000000000000000000000000060c"),
		},
		{
			proxyDefinition2,
			MustHexDecodeString("0x736f6d652d6f746865722d6163636f756e7400000000000000000000000000000900"),
		},
	})
}

func TestProxyDefinition_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{
			MustHexDecodeString("0x736f6d652d6163636f756e740000000000000000000000000000000000000000060c"),
			proxyDefinition1,
		},
		{
			MustHexDecodeString("0x736f6d652d6f746865722d6163636f756e7400000000000000000000000000000900"),
			proxyDefinition2,
		},
	})
}

var (
	proxyStorageEntry1 = ProxyStorageEntry{
		ProxyDefinitions: []ProxyDefinition{
			proxyDefinition1,
			proxyDefinition2,
		},
		Balance: NewU128(*big.NewInt(1234)),
	}
	proxyStorageEntryFuzzOpts = proxyDefinitionFuzzOpts
)

func TestProxyStorageEntry_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[ProxyStorageEntry](t, 1000, proxyStorageEntryFuzzOpts...)
	assertDecodeNilData[ProxyStorageEntry](t)
	assertEncodeEmptyObj[ProxyStorageEntry](t, 17)
}

func TestProxyStorageEntry_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{
			proxyStorageEntry1,
			MustHexDecodeString("0x08736f6d652d6163636f756e740000000000000000000000000000000000000000060c736f6d652d6f746865722d6163636f756e7400000000000000000000000000000900d2040000000000000000000000000000"),
		},
	})
}

func TestProxyStorageEntry_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{
			MustHexDecodeString("0x08736f6d652d6163636f756e740000000000000000000000000000000000000000060c736f6d652d6f746865722d6163636f756e7400000000000000000000000000000900d2040000000000000000000000000000"),
			proxyStorageEntry1,
		},
	})
}
