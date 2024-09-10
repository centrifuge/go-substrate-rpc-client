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

package extrinsic

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic/extensions"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestExtrinsic_Unsigned_Encode(t *testing.T) {
	var meta types.Metadata

	err := codec.DecodeFromHex(types.MetadataV14Data, &meta)
	assert.NoError(t, err)

	c, err := types.NewCall(&meta, "System.remark", []byte("test"))
	assert.NoError(t, err)

	ext := NewExtrinsic(c)

	extEnc, err := codec.EncodeToHex(ext)
	assert.NoError(t, err)

	assert.Equal(t, "0x"+
		"20"+ // length prefix, compact
		"04"+ // version
		"0000"+ // call index
		"10"+"74657374", // remark
		extEnc,
	)
}

func TestExtrinsic_Signed_Encode(t *testing.T) {
	var meta types.Metadata

	err := codec.DecodeFromHex(types.MetadataV14Data, &meta)
	assert.NoError(t, err)

	c, err := types.NewCall(&meta, "System.remark", []byte("test"))
	assert.NoError(t, err)

	ext := NewExtrinsic(c)

	var genesisHash types.Hash

	err = codec.Decode([]byte("0xc787b4dfaa5c0b163fa24eeeb8bf2d06188f81c1beb7ebea76e581549f55254d"), &genesisHash)
	assert.NoError(t, err)

	specVersion := types.U32(1402)
	txVersion := types.U32(2)

	err = ext.Sign(signature.TestKeyringPairAlice, &meta,
		WithEra(types.ExtrinsicEra{IsImmortalEra: true}, genesisHash),
		WithNonce(types.NewUCompactFromUInt(uint64(1))),
		WithTip(types.NewUCompactFromUInt(0)),
		WithSpecVersion(specVersion),
		WithTransactionVersion(txVersion),
		WithGenesisHash(genesisHash),
		WithMetadataMode(extensions.CheckMetadataModeDisabled, extensions.CheckMetadataHash{Hash: types.NewEmptyOption[types.H256]()}),
	)
	assert.NoError(t, err)

	encodedSignature, err := codec.EncodeToHex(ext.Signature)
	assert.NoError(t, err)

	assert.True(t, strings.HasPrefix(
		encodedSignature,
		"0x"+
			"00"+"d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d", // signer
	),
	)

	assert.True(t, strings.HasSuffix(
		encodedSignature,
		"00"+ // era
			"04"+ // nonce compact
			"00"+ // tip
			"00", // mode
	),
	)

	extEnc, err := codec.EncodeToHex(ext)
	assert.NoError(t, err)

	assert.Equal(
		t,
		extEnc,
		"0x"+
			"b901"+"84"+ // prefix
			strings.TrimPrefix(encodedSignature, "0x")+ // signature
			"0000"+ // call index
			"10"+"74657374", // remark
	)
}

func TestExtrinsic_Sign(t *testing.T) {
	call := types.Call{}
	extrinsic := NewExtrinsic(call)

	var meta types.Metadata

	err := codec.DecodeFromHex(types.MetadataV14Data, &meta)
	assert.NoError(t, err)

	err = extrinsic.Sign(signature.TestKeyringPairAlice, &meta,
		WithEra(types.ExtrinsicEra{IsImmortalEra: true}, types.Hash{}),
		WithNonce(types.NewUCompactFromUInt(uint64(0))),
		WithTip(types.NewUCompactFromUInt(0)),
		WithSpecVersion(123),
		WithTransactionVersion(456),
		WithGenesisHash(types.Hash{}),
		WithMetadataMode(extensions.CheckMetadataModeDisabled, extensions.CheckMetadataHash{Hash: types.NewEmptyOption[types.H256]()}),
	)
	assert.NoError(t, err)
	assert.True(t, extrinsic.Version&BitSigned > 1)
	assert.NotNil(t, extrinsic.Signature)
}

func TestExtrinsic_Sign_InvalidVersionError(t *testing.T) {
	extrinsic := &Extrinsic{}

	var meta types.Metadata

	err := codec.DecodeFromHex(types.MetadataV14Data, &meta)
	assert.NoError(t, err)

	err = extrinsic.Sign(signature.TestKeyringPairAlice, &meta,
		WithEra(types.ExtrinsicEra{IsImmortalEra: true}, types.Hash{}),
		WithNonce(types.NewUCompactFromUInt(uint64(0))),
		WithTip(types.NewUCompactFromUInt(0)),
		WithSpecVersion(123),
		WithTransactionVersion(456),
		WithGenesisHash(types.Hash{}),
		WithMetadataMode(extensions.CheckMetadataModeDisabled, extensions.CheckMetadataHash{Hash: types.NewEmptyOption[types.H256]()}),
	)
	assert.ErrorIs(t, err, ErrInvalidVersion)
}

func TestExtrinsic_Sign_PayloadCreationError(t *testing.T) {
	call := types.Call{}
	extrinsic := NewExtrinsic(call)

	var meta types.Metadata

	err := codec.DecodeFromHex(types.MetadataV14Data, &meta)
	assert.NoError(t, err)

	meta.AsMetadataV14.Extrinsic.SignedExtensions = append(
		meta.AsMetadataV14.Extrinsic.SignedExtensions,
		types.SignedExtensionMetadataV14{
			Identifier:       "unsupported_extension",
			Type:             types.Si1LookupTypeID{},
			AdditionalSigned: types.Si1LookupTypeID{},
		},
	)

	err = extrinsic.Sign(signature.TestKeyringPairAlice, &meta,
		WithEra(types.ExtrinsicEra{IsImmortalEra: true}, types.Hash{}),
		WithNonce(types.NewUCompactFromUInt(uint64(0))),
		WithTip(types.NewUCompactFromUInt(0)),
		WithSpecVersion(123),
		WithTransactionVersion(456),
		WithGenesisHash(types.Hash{}),
		WithMetadataMode(extensions.CheckMetadataModeDisabled, extensions.CheckMetadataHash{Hash: types.NewEmptyOption[types.H256]()}),
	)
	assert.ErrorIs(t, err, ErrPayloadCreation)
}

func TestExtrinsic_Sign_MultiAddressCreationError(t *testing.T) {
	call := types.Call{}
	extrinsic := NewExtrinsic(call)

	var meta types.Metadata

	err := codec.DecodeFromHex(types.MetadataV14Data, &meta)
	assert.NoError(t, err)

	err = extrinsic.Sign(signature.KeyringPair{}, &meta,
		WithEra(types.ExtrinsicEra{IsImmortalEra: true}, types.Hash{}),
		WithNonce(types.NewUCompactFromUInt(uint64(0))),
		WithTip(types.NewUCompactFromUInt(0)),
		WithSpecVersion(123),
		WithTransactionVersion(456),
		WithGenesisHash(types.Hash{}),
		WithMetadataMode(extensions.CheckMetadataModeDisabled, extensions.CheckMetadataHash{Hash: types.NewEmptyOption[types.H256]()}),
	)
	assert.ErrorIs(t, err, ErrMultiAddressCreation)
}
