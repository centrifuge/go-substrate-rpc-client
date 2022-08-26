package types_test

import (
	"math/big"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	proxyDefinition1 = ProxyDefinition{
		Delegate:  newTestAccountID(),
		ProxyType: ProxyTypePrice,
		Delay:     3,
	}
	proxyDefinition2 = ProxyDefinition{
		Delegate:  newTestAccountID(),
		ProxyType: KeystoreManagement,
		Delay:     0,
	}
	proxyDefinitionFuzzOpts = proxyTypeFuzzOpts
)

func TestProxyDefinition_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[ProxyDefinition](t, 1000, proxyDefinitionFuzzOpts...)
	assertDecodeNilData[ProxyDefinition](t)
	assertEncodeEmptyObj[ProxyDefinition](t, 37)
}

func TestProxyDefinition_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{
			proxyDefinition1,
			MustHexDecodeString("0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f200603000000"),
		},
		{
			proxyDefinition2,
			MustHexDecodeString("0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f200900000000"),
		},
	})
}

func TestProxyDefinition_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{
			MustHexDecodeString("0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f200603000000"),
			proxyDefinition1,
		},
		{
			MustHexDecodeString("0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f200900000000"),
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
			MustHexDecodeString("0x080102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f2006030000000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f200900000000d2040000000000000000000000000000"),
		},
	})
}

func TestProxyStorageEntry_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{
			MustHexDecodeString("0x080102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f2006030000000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f200900000000d2040000000000000000000000000000"),
			proxyStorageEntry1,
		},
	})
}
