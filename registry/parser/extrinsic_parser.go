package parser

import (
	"bytes"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/chain/generic"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

// DefaultExtrinsic is the Extrinsic with defaults for the generic types:
//
// Address - types.MultiAddress
// Signature - types.MultiSignature
// PaymentFields - generic.DefaultPaymentFields
type DefaultExtrinsic = Extrinsic[
	types.MultiAddress,
	types.MultiSignature,
	generic.DefaultPaymentFields,
]

// Extrinsic holds all the information of a decoded block extrinsic.
//
// This type is generic over types A, S, P, please check generic.GenericExtrinsicSignature for more
// information about these generic types.
type Extrinsic[A, S, P any] struct {
	Name       string
	CallFields registry.DecodedFields
	CallIndex  types.CallIndex
	Version    byte
	Signature  generic.GenericExtrinsicSignature[A, S, P]
}

//nolint:lll
//go:generate mockery --name ExtrinsicParser --structname ExtrinsicParserMock --filename extrinsic_parser_mock.go --inpackage

// ExtrinsicParser is the interface used for parsing a block's extrinsics into []*Extrinsic.
//
// This interface is generic over types A, S, P, please check generic.GenericExtrinsicSignature for more
// information about these generic types.
//
//nolint:lll
type ExtrinsicParser[A, S, P any] interface {
	ParseExtrinsics(callRegistry registry.CallRegistry, block generic.GenericSignedBlock[A, S, P]) ([]*Extrinsic[A, S, P], error)
}

// ExtrinsicParserFn implements ExtrinsicParser.
//
//nolint:lll
type ExtrinsicParserFn[A, S, P any] func(callRegistry registry.CallRegistry, block generic.GenericSignedBlock[A, S, P]) ([]*Extrinsic[A, S, P], error)

// ParseExtrinsics is the function required for satisfying the ExtrinsicParser interface.
//
//nolint:lll
func (e ExtrinsicParserFn[A, S, P]) ParseExtrinsics(callRegistry registry.CallRegistry, block generic.GenericSignedBlock[A, S, P]) ([]*Extrinsic[A, S, P], error) {
	return e(callRegistry, block)
}

// NewExtrinsicParser creates a new ExtrinsicParser.
func NewExtrinsicParser[A, S, P any]() ExtrinsicParser[A, S, P] {
	// The ExtrinsicParserFn provided here is attempting to decode the types.Args of an extrinsic's method
	// into a map of fields and their respective decoded values.
	//
	//nolint:lll
	return ExtrinsicParserFn[A, S, P](func(callRegistry registry.CallRegistry, block generic.GenericSignedBlock[A, S, P]) ([]*Extrinsic[A, S, P], error) {
		var extrinsics []*Extrinsic[A, S, P]

		for i, extrinsic := range block.GetGenericBlock().GetExtrinsics() {
			callIndex := extrinsic.GetCall().CallIndex

			callDecoder, ok := callRegistry[callIndex]

			if !ok {
				return nil, ErrCallDecoderNotFound.Wrap(fmt.Errorf("extrinsic #%d", i))
			}

			decoder := scale.NewDecoder(bytes.NewReader(extrinsic.GetCall().Args))

			callFields, err := callDecoder.Decode(decoder)

			if err != nil {
				return nil, ErrCallFieldsDecoding.Wrap(fmt.Errorf("extrinsic #%d: %w", i, err))
			}

			call := &Extrinsic[A, S, P]{
				Name:       callDecoder.Name,
				CallFields: callFields,
				CallIndex:  callIndex,
				Version:    extrinsic.GetVersion(),
				Signature:  extrinsic.GetSignature(),
			}

			extrinsics = append(extrinsics, call)
		}

		return extrinsics, nil
	})
}

// DefaultExtrinsicParser is the ExtrinsicParser interface with defaults for the generic types:
//
// Address - types.MultiAddress
// Signature - types.MultiSignature
// PaymentFields - generic.DefaultPaymentFields
type DefaultExtrinsicParser = ExtrinsicParser[
	types.MultiAddress,
	types.MultiSignature,
	generic.DefaultPaymentFields,
]

// NewDefaultExtrinsicParser returns a DefaultExtrinsicParser.
func NewDefaultExtrinsicParser() DefaultExtrinsicParser {
	return NewExtrinsicParser[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]()
}
