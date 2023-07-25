package parser

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/chain/generic"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestExtrinsicParserFn_ParseExtrinsics(t *testing.T) {
	testExtrinsics := []testExtrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]{
		{
			Name: "extrinsic_1",
			CallIndex: types.CallIndex{
				SectionIndex: 0,
				MethodIndex:  1,
			},
			CallFields: []testField{
				{
					Name:  "bool_field",
					Value: true,
				},
				{
					Name:  "byte_field",
					Value: byte(32),
				},
				{
					Name:  "string_field",
					Value: "test",
				},
			},
			Version: 1,
			Signature: &generic.ExtrinsicSignature[
				types.MultiAddress,
				types.MultiSignature,
				generic.DefaultPaymentFields,
			]{},
		},
		{
			Name: "extrinsic_2",
			CallIndex: types.CallIndex{
				SectionIndex: 1,
				MethodIndex:  0,
			},
			CallFields: []testField{
				{
					Name:  "u8_value",
					Value: types.NewU8(11),
				},
				{
					Name:  "u16_value",
					Value: types.NewU16(121),
				},
				{
					Name:  "u32_value",
					Value: types.NewU32(12134),
				},
				{
					Name:  "u64_value",
					Value: types.NewU64(128678),
				},
				{
					Name:  "u128_value",
					Value: types.NewU128(*big.NewInt(56346)),
				},
				{
					Name:  "u256_value",
					Value: types.NewU256(*big.NewInt(5674)),
				},
			},
			Version: 2,
			Signature: &generic.ExtrinsicSignature[
				types.MultiAddress,
				types.MultiSignature,
				generic.DefaultPaymentFields,
			]{},
		},
		{
			Name: "extrinsic_3",
			CallIndex: types.CallIndex{
				SectionIndex: 1,
				MethodIndex:  1,
			},
			CallFields: []testField{
				{
					Name:  "i8_value",
					Value: types.NewI8(45),
				},
				{
					Name:  "i16_value",
					Value: types.NewI16(445),
				},
				{
					Name:  "i32_value",
					Value: types.NewI32(545),
				},
				{
					Name:  "i64_value",
					Value: types.NewI64(4789),
				},
				{
					Name:  "i128_value",
					Value: types.NewI128(*big.NewInt(56747)),
				},
				{
					Name:  "i256_value",
					Value: types.NewI256(*big.NewInt(45356747)),
				},
			},
			Version: 2,
			Signature: &generic.ExtrinsicSignature[
				types.MultiAddress,
				types.MultiSignature,
				generic.DefaultPaymentFields,
			]{},
		},
	}

	block, reg, err := getExtrinsicParsingTestData[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	](testExtrinsics)
	assert.NoError(t, err)

	extrinsicParser := NewExtrinsicParser[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]()

	res, err := extrinsicParser.ParseExtrinsics(reg, block)
	assert.NoError(t, err)
	assert.Len(t, res, len(testExtrinsics))

	for i, testExtrinsic := range testExtrinsics {
		assert.Equal(t, testExtrinsic.Name, res[i].Name)
		assert.Equal(t, testExtrinsic.Version, res[i].Version)
		assertExtrinsicFieldInformationIsCorrect(t, testExtrinsic.CallFields, res[i])
		assert.Equal(t, testExtrinsic.CallIndex, res[i].CallIndex)
		assert.Equal(t, testExtrinsic.Signature, res[i].Signature)
	}
}

func TestExtrinsicParserFn_ParseExtrinsics_MissingCallDecoder(t *testing.T) {
	testExtrinsics := []testExtrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]{
		{
			Name: "extrinsic_1",
			CallIndex: types.CallIndex{
				SectionIndex: 0,
				MethodIndex:  1,
			},
			CallFields: []testField{
				{
					Name:  "bool_field",
					Value: true,
				},
				{
					Name:  "byte_field",
					Value: byte(32),
				},
				{
					Name:  "string_field",
					Value: "test",
				},
			},
			Version: 1,
			Signature: &generic.ExtrinsicSignature[
				types.MultiAddress,
				types.MultiSignature,
				generic.DefaultPaymentFields,
			]{},
		},
	}

	block, _, err := getExtrinsicParsingTestData[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	](testExtrinsics)
	assert.NoError(t, err)

	extrinsicParser := NewExtrinsicParser[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]()

	// Empty registry, decoding should fail.
	res, err := extrinsicParser.ParseExtrinsics(registry.CallRegistry{}, block)
	assert.ErrorIs(t, err, ErrCallDecoderNotFound)
	assert.Nil(t, res)
}

func TestExtrinsicParserFn_ParseExtrinsics_DecodeError(t *testing.T) {
	testExtrinsics := []testExtrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]{
		{
			Name: "extrinsic_1",
			CallIndex: types.CallIndex{
				SectionIndex: 0,
				MethodIndex:  1,
			},
			CallFields: []testField{
				{
					Name:  "bool_field",
					Value: true,
				},
				{
					Name:  "byte_field",
					Value: byte(32),
				},
				{
					Name:  "string_field",
					Value: "test",
				},
			},
			Version: 1,
			Signature: &generic.ExtrinsicSignature[
				types.MultiAddress,
				types.MultiSignature,
				generic.DefaultPaymentFields,
			]{},
		},
	}

	block, reg, err := getExtrinsicParsingTestData[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	](testExtrinsics)
	assert.NoError(t, err)

	extrinsicParser := NewExtrinsicParser[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]()

	// No args for extrinsics should trigger an error.
	block.Block.Extrinsics[0].Method.Args = []byte{}

	res, err := extrinsicParser.ParseExtrinsics(reg, block)
	assert.ErrorIs(t, err, ErrCallFieldsDecoding)
	assert.Nil(t, res)
}

func assertExtrinsicFieldInformationIsCorrect[A, S, P any](t *testing.T, testFields []testField, extrinsic *Extrinsic[A, S, P]) {
	for testFieldIndex, testField := range testFields {
		assert.Equal(t, testField.Value, extrinsic.CallFields[testFieldIndex].Value)
	}
}

func getExtrinsicParsingTestData[A, S, P any](testExtrinsics []testExtrinsic[A, S, P]) (*generic.SignedBlock[A, S, P], registry.CallRegistry, error) {
	callRegistry, err := getRegistryForTestExtrinsic(testExtrinsics)

	if err != nil {
		return nil, nil, err
	}

	block := &generic.SignedBlock[A, S, P]{
		Block: &generic.Block[A, S, P]{},
	}

	for _, testExtrinsic := range testExtrinsics {
		encodedExtrinsicCall, err := testExtrinsic.Encode()

		if err != nil {
			return nil, nil, err
		}

		block.Block.Extrinsics = append(block.Block.Extrinsics, &generic.Extrinsic[A, S, P]{
			Version:   testExtrinsic.Version,
			Signature: testExtrinsic.Signature,
			Method: types.Call{
				CallIndex: testExtrinsic.CallIndex,
				Args:      encodedExtrinsicCall,
			},
		})
	}

	return block, callRegistry, nil
}

func getRegistryForTestExtrinsic[A, S, P any](testExtrinsics []testExtrinsic[A, S, P]) (registry.CallRegistry, error) {
	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.TypeDecoder{})

	for _, testExtrinsic := range testExtrinsics {
		regFields, err := getTestRegistryFields(testExtrinsic.CallFields)

		if err != nil {
			return nil, err
		}

		callRegistry[testExtrinsic.CallIndex] = &registry.TypeDecoder{
			Name:   testExtrinsic.Name,
			Fields: regFields,
		}
	}

	return callRegistry, nil
}

type testExtrinsic[A, S, P any] struct {
	Name       string
	CallIndex  types.CallIndex
	CallFields []testField
	Version    byte
	Signature  *generic.ExtrinsicSignature[A, S, P]
}

func (t testExtrinsic[A, S, P]) Encode() ([]byte, error) {
	var b []byte

	buf := bytes.NewBuffer(b)

	encoder := scale.NewEncoder(buf)

	for _, field := range t.CallFields {
		if err := encoder.Encode(field.Value); err != nil {
			return nil, fmt.Errorf("couldn't encode field %v: %w", field, err)
		}
	}

	return buf.Bytes(), nil
}
