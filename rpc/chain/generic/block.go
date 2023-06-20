package generic

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
)

// DefaultGenericSignedBlock is the SignedBlock with defaults for the generic types:
//
// Address - types.MultiAddress
// Signature - types.MultiSignature
// PaymentFields - DefaultPaymentFields
type DefaultGenericSignedBlock = SignedBlock[
	types.MultiAddress,
	types.MultiSignature,
	DefaultPaymentFields,
]

// GenericSignedBlock is the interface that represents the block of a particular chain.
//
// This interface is generic over types A, S, P, please check GenericExtrinsicSignature for more
// information about these generic types.
//
//nolint:revive
type GenericSignedBlock[A, S, P any] interface {
	GetGenericBlock() GenericBlock[A, S, P]
	GetJustification() []byte
}

// GenericBlock is the interface that holds information about the header and extrinsics of a block.
//
// This interface is generic over types A, S, P, please check GenericExtrinsicSignature for more
// information about these generic types.
//
//nolint:revive
type GenericBlock[A, S, P any] interface {
	GetHeader() types.Header
	GetExtrinsics() []GenericExtrinsic[A, S, P]
}

// GenericExtrinsic is the interface that holds the extrinsic information.
//
// This interface is generic over types A, S, P, please check GenericExtrinsicSignature for more
// information about these generic types.
//
//nolint:revive
type GenericExtrinsic[A, S, P any] interface {
	GetVersion() byte
	GetSignature() GenericExtrinsicSignature[A, S, P]
	GetCall() types.Call
}

// GenericExtrinsicSignature is the interface that holds the extrinsic signature information.
//
// This interface is generic over the following types:
//
// A - Signer, the default implementation for this is the types.MultiAddress type
// which can support a variable number of addresses.
//
// S - Signature, the default implementation for this is the types.MultiSignature type which can support
// multiple signature curves.
//
// P - PaymentFields (ChargeAssetTx in substrate), the default implementation for this is the DefaultPaymentFields which
// holds information about the tip sent for the extrinsic.
//
//nolint:revive
type GenericExtrinsicSignature[A, S, P any] interface {
	GetSigner() A
	GetSignature() S
	GetEra() types.ExtrinsicEra
	GetNonce() types.UCompact
	GetPaymentFields() P
}

// SignedBlock implements the GenericSignedBlock interface.
type SignedBlock[A, S, P any] struct {
	Block         *Block[A, S, P] `json:"block"`
	Justification []byte          `json:"justification"`
}

func (s *SignedBlock[A, S, P]) GetGenericBlock() GenericBlock[A, S, P] {
	return s.Block
}

func (s *SignedBlock[A, S, P]) GetJustification() []byte {
	return s.Justification
}

// Block implements the GenericBlock interface.
type Block[A, S, P any] struct {
	Header     types.Header          `json:"header"`
	Extrinsics []*Extrinsic[A, S, P] `json:"extrinsics"`
}

//nolint:revive
func (b *Block[A, S, P]) GetHeader() types.Header {
	return b.Header
}

//nolint:revive
func (b *Block[A, S, P]) GetExtrinsics() []GenericExtrinsic[A, S, P] {
	var res []GenericExtrinsic[A, S, P]

	for _, ext := range b.Extrinsics {
		res = append(res, ext)
	}
	return res
}

// Extrinsic implements the GenericExtrinsic interface.
type Extrinsic[A, S, P any] struct {
	// Version is the encoded version flag (which encodes the raw transaction version and signing information in one byte)
	Version byte `json:"version"`
	// Signature is the ExtrinsicSignature, its presence depends on the Version flag
	Signature *ExtrinsicSignature[A, S, P] `json:"signature"`
	// Method is the call this extrinsic wraps
	Method types.Call `json:"method"`
}

//nolint:revive
func (e *Extrinsic[A, S, P]) GetVersion() byte {
	return e.Version
}

//nolint:revive
func (e *Extrinsic[A, S, P]) GetSignature() GenericExtrinsicSignature[A, S, P] {
	return e.Signature
}

//nolint:revive
func (e *Extrinsic[A, S, P]) GetCall() types.Call {
	return e.Method
}

// UnmarshalJSON fills Extrinsic with the JSON encoded byte array given by bz
//nolint:revive
func (e *Extrinsic[A, S, P]) UnmarshalJSON(bz []byte) error {
	var tmp string
	if err := json.Unmarshal(bz, &tmp); err != nil {
		return err
	}

	var l types.UCompact

	if err := codec.DecodeFromHex(tmp, &l); err != nil {
		return err
	}

	prefix, err := codec.EncodeToHex(l)
	if err != nil {
		return err
	}

	// determine whether length prefix is there
	if strings.HasPrefix(tmp, prefix) {
		return codec.DecodeFromHex(tmp, e)
	}

	// not there, prepend with compact encoded length prefix
	dec, err := codec.HexDecodeString(tmp)
	if err != nil {
		return err
	}
	length := types.NewUCompactFromUInt(uint64(len(dec)))
	bprefix, err := codec.Encode(length)
	if err != nil {
		return err
	}
	bprefix = append(bprefix, dec...)
	return codec.Decode(bprefix, e)
}

// IsSigned returns true if the extrinsic is signed.
//nolint:revive
func (e *Extrinsic[A, S, P]) IsSigned() bool {
	return e.Version&types.ExtrinsicBitSigned == types.ExtrinsicBitSigned
}

// Type returns the raw transaction version.
//nolint:revive
func (e *Extrinsic[A, S, P]) Type() uint8 {
	return e.Version & types.ExtrinsicUnmaskVersion
}

// Decode decodes the extrinsic based on the data present in the decoder.
//nolint:revive
func (e *Extrinsic[A, S, P]) Decode(decoder scale.Decoder) error {
	// compact length encoding (1, 2, or 4 bytes) (may not be there for Extrinsics older than Jan 11 2019)
	if _, err := decoder.DecodeUintCompact(); err != nil {
		return err
	}

	if err := decoder.Decode(&e.Version); err != nil {
		return err
	}

	if e.IsSigned() {
		if e.Type() != types.ExtrinsicVersion4 {
			return fmt.Errorf("unsupported extrinsic version: %v (isSigned: %v, type: %v)", e.Version, e.IsSigned(),
				e.Type())
		}

		e.Signature = new(ExtrinsicSignature[A, S, P])

		if err := decoder.Decode(&e.Signature); err != nil {
			return err
		}
	}

	if err := decoder.Decode(&e.Method); err != nil {
		return err
	}

	return nil
}

// ExtrinsicSignature implements the GenericExtrinsicSignature interface.
type ExtrinsicSignature[A, S, P any] struct {
	Signer        A
	Signature     S
	Era           types.ExtrinsicEra
	Nonce         types.UCompact
	PaymentFields P
}

//nolint:revive
func (e *ExtrinsicSignature[A, S, P]) GetSigner() A {
	return e.Signer
}

//nolint:revive
func (e *ExtrinsicSignature[A, S, P]) GetSignature() S {
	return e.Signature
}

//nolint:revive
func (e *ExtrinsicSignature[A, S, P]) GetEra() types.ExtrinsicEra {
	return e.Era
}

//nolint:revive
func (e *ExtrinsicSignature[A, S, P]) GetNonce() types.UCompact {
	return e.Nonce
}

//nolint:revive
func (e *ExtrinsicSignature[A, S, P]) GetPaymentFields() P {
	return e.PaymentFields
}

// DefaultPaymentFields represents the default payment fields found in the extrinsics of most substrate chains.
type DefaultPaymentFields struct {
	Tip types.UCompact
}

// PaymentFieldsWithAssetID represents the payment fields found on chains that require an asset ID, such as statemint.
type PaymentFieldsWithAssetID struct {
	Tip     types.UCompact
	AssetID types.Option[types.UCompact]
}
