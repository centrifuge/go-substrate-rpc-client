package parser

import (
	"bytes"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Extrinsic struct {
	Name       string
	CallFields map[string]any
	CallIndex  types.CallIndex
	Version    byte
	Signature  types.ExtrinsicSignatureV4
}

type ExtrinsicParser interface {
	ParseExtrinsics(callRegistry registry.CallRegistry, block *types.SignedBlock) ([]*Extrinsic, error)
}

type ExtrinsicParserFn func(callRegistry registry.CallRegistry, block *types.SignedBlock) ([]*Extrinsic, error)

func (e ExtrinsicParserFn) ParseExtrinsics(callRegistry registry.CallRegistry, block *types.SignedBlock) ([]*Extrinsic, error) {
	return e(callRegistry, block)
}

func NewExtrinsicParser() ExtrinsicParser {
	return ExtrinsicParserFn(func(callRegistry registry.CallRegistry, block *types.SignedBlock) ([]*Extrinsic, error) {
		var extrinsics []*Extrinsic

		for i, extrinsic := range block.Block.Extrinsics {
			callIndex := extrinsic.Method.CallIndex

			callDecoder, ok := callRegistry[callIndex]

			if !ok {
				return nil, fmt.Errorf("couldn't find call decoder for extrinsic #%d", i)
			}

			decoder := scale.NewDecoder(bytes.NewReader(extrinsic.Method.Args))

			callFields, err := callDecoder.Decode(decoder)

			if err != nil {
				return nil, fmt.Errorf("couldn't decode call fields for extrinsic #%d: %w", i, err)
			}

			call := &Extrinsic{
				Name:       callDecoder.Name,
				CallFields: callFields,
				CallIndex:  callIndex,
				Version:    extrinsic.Version,
				Signature:  extrinsic.Signature,
			}

			extrinsics = append(extrinsics, call)
		}

		return extrinsics, nil
	})
}
