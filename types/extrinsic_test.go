// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types_test

import (
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/signature"
	. "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

func TestExtrinsic_Unsigned_EncodeDecode(t *testing.T) {
	addr, err := NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	assert.NoError(t, err)

	c, err := NewCall(ExamplaryMetadataV4, "balances.transfer", addr, NewUCompactFromUInt(6969))
	assert.NoError(t, err)

	ext := NewExtrinsic(c)

	extEnc, err := EncodeToHexString(ext)
	assert.NoError(t, err)

	assert.Equal(t, "0x"+
		"98"+ // length prefix, compact
		"04"+ // version
		"0300"+ // call index (section index and method index)
		"ff"+
		"8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48"+ // target address
		"e56c", // amount, compact
		extEnc)

	var extDec Extrinsic
	err = DecodeFromHexString(extEnc, &extDec)
	assert.NoError(t, err)

	assert.Equal(t, ext, extDec)
}

func TestExtrinsic_Signed_EncodeDecode(t *testing.T) {
	extEnc, err := EncodeToHexString(ExamplaryExtrinsic)
	assert.NoError(t, err)

	var extDec Extrinsic
	err = DecodeFromHexString(extEnc, &extDec)
	assert.NoError(t, err)

	assert.Equal(t, ExamplaryExtrinsic, extDec)
}

func TestExtrinsic_Sign(t *testing.T) {
	c, err := NewCall(ExamplaryMetadataV4,
		"balances.transfer", NewAddressFromAccountID(MustHexDecodeString(
			"0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")),
		NewUCompactFromUInt(6969))
	assert.NoError(t, err)

	ext := NewExtrinsic(c)

	o := SignatureOptions{
		BlockHash: NewHash(MustHexDecodeString("0xec7afaf1cca720ce88c1d1b689d81f0583cc15a97d621cf046dd9abf605ef22f")),
		// Era: ExtrinsicEra{IsImmortalEra: true},
		GenesisHash: NewHash(MustHexDecodeString("0xdcd1346701ca8396496e52aa2785b1748deb6db09551b72159dcb3e08991025b")),
		Nonce:       NewUCompactFromUInt(1),
		SpecVersion: 123,
		Tip:         NewUCompactFromUInt(2),
		TransactionVersion: 1,
	}

	assert.False(t, ext.IsSigned())

	err = ext.Sign(signature.TestKeyringPairAlice, o)
	assert.NoError(t, err)

	// fmt.Printf("%#v", ext)

	assert.True(t, ext.IsSigned())

	extEnc, err := EncodeToHexString(ext)
	assert.NoError(t, err)

	// extEnc will have the structure of the following. It can't be tested, since the signature is different on every
	// call to sign. Instead we verify here.
	// "0x"+
	// "2902"+ // length prefix, compact
	// "83"+ // version
	// "ff"+
	// "d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"+ // signer address
	// "6667a2afe5272b327c3886036d2906ceac90fe959377a2d47fa92b6ebe345318379fff37e48a4e8fd552221796dd6329d028f80237"+
	// 		"ebc0abb229ca2235778308"+ // signature
	// "000408"+ // era, nonce, tip
	// "0300" + // call index (section index and method index)
	// "ff"+
	// "8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48"+ // target address
	// "e56c", // amount, compact

	var extDec Extrinsic
	err = DecodeFromHexString(extEnc, &extDec)
	assert.NoError(t, err)

	assert.Equal(t, uint8(ExtrinsicVersion4), extDec.Type())
	assert.Equal(t, signature.TestKeyringPairAlice.PublicKey, extDec.Signature.Signer.AsAccountID[:])

	mb, err := EncodeToBytes(extDec.Method)
	assert.NoError(t, err)

	verifyPayload := ExtrinsicPayloadV4{
		ExtrinsicPayloadV3: ExtrinsicPayloadV3{
			Method:      mb,
			Era:         extDec.Signature.Era,
			Nonce:       extDec.Signature.Nonce,
			Tip:         extDec.Signature.Tip,
			SpecVersion: o.SpecVersion,
			GenesisHash: o.GenesisHash,
			BlockHash:   o.BlockHash,
		},
		TransactionVersion: 1,
	}

	// verify sig
	b, err := EncodeToBytes(verifyPayload)
	assert.NoError(t, err)
	ok, err := signature.Verify(b, extDec.Signature.Signature.AsSr25519[:], HexEncodeToString(signature.TestKeyringPairAlice.PublicKey))
	assert.NoError(t, err)
	assert.True(t, ok)
}

func ExampleExtrinsic() {
	bob, err := NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}

	c, err := NewCall(ExamplaryMetadataV4, "balances.transfer", bob, NewUCompactFromUInt(6969))
	if err != nil {
		panic(err)
	}

	ext := NewExtrinsic(c)

	ext.Method.CallIndex.SectionIndex = 5
	ext.Method.CallIndex.MethodIndex = 0

	era := ExtrinsicEra{IsMortalEra: true, AsMortalEra: MortalEra{0x95, 0x00}}

	o := SignatureOptions{
		BlockHash:   NewHash(MustHexDecodeString("0x223e3eb79416e6258d262b3a76e827aa0886b884a96bf96395cdd1c52d0eeb45")),
		Era:         era,
		GenesisHash: NewHash(MustHexDecodeString("0x81ad0bfe2a0bccd91d2e89852d79b7ff696d4714758e5f7c6f17ec7527e1f550")),
		Nonce:       NewUCompactFromUInt(1),
		SpecVersion: 170,
		Tip:         NewUCompactFromUInt(0),
	}

	err = ext.Sign(signature.TestKeyringPairAlice, o)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", ext)

	extEnc, err := EncodeToHexString(ext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", extEnc)
}

func TestCall(t *testing.T) {
	c := Call{CallIndex{6, 1}, Args{0, 0, 0}}

	enc, err := EncodeToHexString(c)
	assert.NoError(t, err)
	assert.Equal(t, "0x0601000000", enc)
}

func TestNewCallV4(t *testing.T) {
	addr, err := NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	assert.NoError(t, err)

	c, err := NewCall(ExamplaryMetadataV4, "balances.transfer", addr, NewUCompactFromUInt(1000))
	assert.NoError(t, err)

	enc, err := EncodeToHexString(c)
	assert.NoError(t, err)

	assert.Equal(t, "0x0300ff8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48a10f", enc)
}

func TestNewCallV7(t *testing.T) {
	c, err := NewCall(&exampleMetadataV7, "Module2.my function", U8(3))
	assert.NoError(t, err)

	enc, err := EncodeToHexString(c)
	assert.NoError(t, err)

	assert.Equal(t, "0x010003", enc)
}

func TestNewCallV8(t *testing.T) {
	c, err := NewCall(&exampleMetadataV8, "Module2.my function", U8(3))
	assert.NoError(t, err)

	enc, err := EncodeToHexString(c)
	assert.NoError(t, err)

	assert.Equal(t, "0x010003", enc)
}
