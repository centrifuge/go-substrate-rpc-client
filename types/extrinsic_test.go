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

	c, err := NewCall(ExamplaryMetadataV4, "balances.transfer", addr, UCompact(6969))
	assert.NoError(t, err)

	ext := NewExtrinsic(c)

	extEnc, err := EncodeToHexString(ext)
	assert.NoError(t, err)

	assert.Equal(t, "0x"+
		"98"+ // length prefix, compact
		"03"+ // version
		"0200"+ // call index (section index and method index)
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
	ext := Extrinsic{Version: 0x83, Signature: ExtrinsicSignatureV3{Signer: Address{IsAccountID: true, AsAccountID: AccountID{0xd4, 0x35, 0x93, 0xc7, 0x15, 0xfd, 0xd3, 0x1c, 0x61, 0x14, 0x1a, 0xbd, 0x4, 0xa9, 0x9f, 0xd6, 0x82, 0x2c, 0x85, 0x58, 0x85, 0x4c, 0xcd, 0xe3, 0x9a, 0x56, 0x84, 0xe7, 0xa5, 0x6d, 0xa2, 0x7d}, IsAccountIndex: false, AsAccountIndex: 0x0}, Signature: Signature{0x5c, 0x77, 0x1d, 0xd5, 0x6a, 0xe0, 0xce, 0xed, 0x68, 0xd, 0xb3, 0xbb, 0x4c, 0x40, 0x7a, 0x38, 0x96, 0x99, 0x97, 0xae, 0xb6, 0xa, 0x2c, 0x62, 0x39, 0x1, 0x6, 0x2f, 0x7f, 0x8e, 0xbf, 0x2f, 0xe7, 0x73, 0x3a, 0x61, 0x3c, 0xf1, 0x6b, 0x78, 0xf6, 0x10, 0xc6, 0x52, 0x32, 0xa2, 0x3c, 0xc5, 0xce, 0x25, 0xda, 0x29, 0xa3, 0xd5, 0x84, 0x85, 0xd8, 0x7b, 0xd8, 0x3d, 0xb8, 0x18, 0x3f, 0x8}, Era: ExtrinsicEra{IsImmortalEra: true, IsMortalEra: false, AsMortalEra: MortalEra{First: 0x0, Second: 0x0}}, Nonce: 0x1, Tip: 0x2}, Method: Call{CallIndex: CallIndex{SectionIndex: 0x3, MethodIndex: 0x0}, Args: Args{0xff, 0x8e, 0xaf, 0x4, 0x15, 0x16, 0x87, 0x73, 0x63, 0x26, 0xc9, 0xfe, 0xa1, 0x7e, 0x25, 0xfc, 0x52, 0x87, 0x61, 0x36, 0x93, 0xc9, 0x12, 0x90, 0x9c, 0xb2, 0x26, 0xaa, 0x47, 0x94, 0xf2, 0x6a, 0x48, 0xe5, 0x6c}}} //nolint:lll

	extEnc, err := EncodeToHexString(ext)
	assert.NoError(t, err)

	var extDec Extrinsic
	err = DecodeFromHexString(extEnc, &extDec)
	assert.NoError(t, err)

	assert.Equal(t, ext, extDec)
}

func TestExtrinsic_Sign(t *testing.T) {
	c, err := NewCall(ExamplaryMetadataV4,
		"balances.transfer", NewAddressFromAccountID(MustHexDecodeString(
			"0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")),
		UCompact(6969))
	assert.NoError(t, err)

	ext := NewExtrinsic(c)

	o := SignatureOptions{
		BlockHash: NewHash(MustHexDecodeString("0xec7afaf1cca720ce88c1d1b689d81f0583cc15a97d621cf046dd9abf605ef22f")),
		// Era: ExtrinsicEra{IsImmortalEra: true},
		GenesisHash: NewHash(MustHexDecodeString("0xdcd1346701ca8396496e52aa2785b1748deb6db09551b72159dcb3e08991025b")),
		Nonce:       1,
		SpecVersion: 123,
		Tip:         2,
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

	assert.Equal(t, uint8(ExtrinsicVersion3), extDec.Type())
	assert.Equal(t, signature.TestKeyringPairAlice.PublicKey, extDec.Signature.Signer.AsAccountID[:])

	mb, err := EncodeToBytes(extDec.Method)
	assert.NoError(t, err)

	verifyPayload := ExtrinsicPayloadV3{
		Method:      mb,
		Era:         extDec.Signature.Era,
		Nonce:       extDec.Signature.Nonce,
		Tip:         extDec.Signature.Tip,
		SpecVersion: o.SpecVersion,
		GenesisHash: o.GenesisHash,
		BlockHash:   o.BlockHash,
	}

	// verify sig
	b, err := EncodeToBytes(verifyPayload)
	assert.NoError(t, err)
	ok, err := signature.Verify(b, extDec.Signature.Signature[:], signature.TestKeyringPairAlice.URI)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func ExampleExtrinsic() {
	bob, err := NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}

	c, err := NewCall(ExamplaryMetadataV4, "balances.transfer", bob, UCompact(6969))
	if err != nil {
		panic(err)
	}

	ext := NewExtrinsic(c)
	if err != nil {
		panic(err)
	}

	ext.Method.CallIndex.SectionIndex = 5
	ext.Method.CallIndex.MethodIndex = 0

	era := ExtrinsicEra{IsMortalEra: true, AsMortalEra: MortalEra{0x95, 0x00}}

	o := SignatureOptions{
		BlockHash:   NewHash(MustHexDecodeString("0x223e3eb79416e6258d262b3a76e827aa0886b884a96bf96395cdd1c52d0eeb45")),
		Era:         era,
		GenesisHash: NewHash(MustHexDecodeString("0x81ad0bfe2a0bccd91d2e89852d79b7ff696d4714758e5f7c6f17ec7527e1f550")),
		Nonce:       1,
		SpecVersion: 170,
		Tip:         0,
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

	c, err := NewCall(ExamplaryMetadataV4, "balances.transfer", addr, UCompact(1000))
	assert.NoError(t, err)

	enc, err := EncodeToHexString(c)
	assert.NoError(t, err)

	assert.Equal(t, "0x0200ff8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48a10f", enc)
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
