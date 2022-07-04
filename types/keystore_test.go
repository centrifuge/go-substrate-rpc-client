package types_test

import (
	"math/big"
	"testing"

	fuzz "github.com/google/gofuzz"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	keyPurpose1        = KeyPurposeP2PDiscovery
	keyPurpose2        = KeyPurposeP2PDocumentSigning
	keyPurposeFuzzOpts = []fuzzOpt{
		withFuzzFuncs(func(k *KeyPurpose, c fuzz.Continue) {
			*k = KeyPurpose(c.Intn(2))
		}),
	}
)

func TestKeyPurpose_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[KeyPurpose](t, 1000, keyPurposeFuzzOpts...)
	assertDecodeNilData[KeyPurpose](t)
	assertEncodeEmptyObj[KeyPurpose](t, 1)
}

func TestKeyPurpose_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{keyPurpose1, MustHexDecodeString("0x00")},
		{keyPurpose2, MustHexDecodeString("0x01")},
	})
}

func TestKeyPurpose_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x00"), keyPurpose1},
		{MustHexDecodeString("0x01"), keyPurpose2},
	})
}

var (
	keyType1        = KeyTypeECDSA
	keyType2        = KeyTypeEDDSA
	keyTypeFuzzOpts = []fuzzOpt{
		withFuzzFuncs(func(k *KeyType, c fuzz.Continue) {
			*k = KeyType(c.Intn(2))
		}),
	}
)

func TestKeyType_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[KeyType](t, 1000, keyTypeFuzzOpts...)
	assertDecodeNilData[KeyType](t)
	assertEncodeEmptyObj[KeyType](t, 1)
}

func TestKeyType_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{keyType1, MustHexDecodeString("0x00")},
		{keyType2, MustHexDecodeString("0x01")},
	})
}

func TestKeyType_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x00"), keyType1},
		{MustHexDecodeString("0x01"), keyType2},
	})
}

var (
	testKey = Key{
		KeyPurpose: keyPurpose1,
		KeyType:    keyType2,
		RevokedAt:  NewOptionBlockNumber(BlockNumber(3)),
		Deposit:    NewU128(*big.NewInt(123)),
	}
	keyFuzzOpts = combineFuzzOpts(
		keyPurposeFuzzOpts,
		keyTypeFuzzOpts,
		optionBlockNumberFuzzOpts,
	)
)

func TestKey_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[Key](t, 1000, keyFuzzOpts...)
	assertDecodeNilData[Key](t)
	assertEncodeEmptyObj[Key](t, 19)
}

func TestKey_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{
			testKey,
			MustHexDecodeString("0x0001010c7b000000000000000000000000000000"),
		},
	})
}

func TestKey_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{
			MustHexDecodeString("0x0001010c7b000000000000000000000000000000"),
			testKey,
		},
	})
}

var (
	testKeyID = KeyID{
		Hash:       NewHash([]byte("some_hash")),
		KeyPurpose: keyPurpose2,
	}
	keyIDFuzzOpts = combineFuzzOpts(
		keyPurposeFuzzOpts,
		[]fuzzOpt{
			withFuzzFuncs(func(k *KeyID, c fuzz.Continue) {
				c.Fuzz(&k.Hash)
				c.Fuzz(&k.KeyPurpose)
			}),
		},
	)
)

func TestKeyID_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[KeyID](t, 1000, keyIDFuzzOpts...)
	assertDecodeNilData[KeyID](t)
	assertEncodeEmptyObj[KeyID](t, 33)
}

func TestKeyID_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testKey, MustHexDecodeString("0x0001010c7b000000000000000000000000000000")},
	})
}

func TestKeyID_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x0001010c7b000000000000000000000000000000"), testKey},
	})
}

var (
	testAddKey = AddKey{
		Key:     NewHash([]byte("some_hash")),
		Purpose: keyPurpose1,
		KeyType: keyType1,
	}
	addKeyFuzzOpts = combineFuzzOpts(
		keyPurposeFuzzOpts,
		keyTypeFuzzOpts,
		[]fuzzOpt{
			withFuzzFuncs(func(a *AddKey, c fuzz.Continue) {
				c.Fuzz(&a.Key)
				c.Fuzz(&a.Purpose)
				c.Fuzz(&a.KeyType)
			}),
		},
	)
)

func TestAddKey_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[AddKey](t, 1000, addKeyFuzzOpts...)
	assertDecodeNilData[AddKey](t)
	assertEncodeEmptyObj[AddKey](t, 34)
}

func TestAddKey_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{
			testAddKey,
			MustHexDecodeString("0x736f6d655f6861736800000000000000000000000000000000000000000000000000"),
		},
	})
}

func TestAddKey_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{
			MustHexDecodeString("0x736f6d655f6861736800000000000000000000000000000000000000000000000000"),
			testAddKey,
		},
	})
}