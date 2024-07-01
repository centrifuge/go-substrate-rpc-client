package extrinsic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"math/big"
)

type DynamicExtrinsic struct {
	// Version is the encoded version flag (which encodes the raw transaction version and signing information in one byte)
	Version   byte
	Signature DynamicExtrinsicSignature
	// Method is the call this extrinsic wraps
	Method types.Call
}

// NewExtrinsic creates a new Extrinsic from the provided Call
func NewDynamicExtrinsic(c types.Call) DynamicExtrinsic {
	return DynamicExtrinsic{
		Version: types.ExtrinsicVersion4,
		Method:  c,
	}
}

// MarshalJSON returns a JSON encoded byte array of Extrinsic
func (e DynamicExtrinsic) MarshalJSON() ([]byte, error) {
	s, err := codec.EncodeToHex(e)
	if err != nil {
		return nil, err
	}
	return json.Marshal(s)
}

// IsSigned returns true if the extrinsic is signed
func (e DynamicExtrinsic) IsSigned() bool {
	return e.Version&types.ExtrinsicBitSigned == types.ExtrinsicBitSigned
}

// Type returns the raw transaction version (not flagged with signing information)
func (e DynamicExtrinsic) Type() uint8 {
	return e.Version & types.ExtrinsicUnmaskVersion
}

// Sign adds a signature to the extrinsic
func (e *DynamicExtrinsic) Sign(signer signature.KeyringPair, meta *types.Metadata, opts ...DynamicExtrinsicSigningOption) error {
	if e.Type() != types.ExtrinsicVersion4 {
		return fmt.Errorf("unsupported extrinsic version: %v (isSigned: %v, type: %v)", e.Version, e.IsSigned(), e.Type())
	}

	encodedMethod, err := codec.Encode(e.Method)
	if err != nil {
		return err
	}

	extrinsicPayload := DynamicExtrinsicPayload{}
	extrinsicPayload.Method = encodedMethod
	extrinsicSignature := DynamicExtrinsicSignature{}

	for _, opt := range opts {
		opt(&extrinsicPayload, &extrinsicSignature)
	}

	signerPubKey, err := types.NewMultiAddressFromAccountID(signer.PublicKey)

	if err != nil {
		return err
	}

	sig, err := extrinsicPayload.Sign(signer)
	if err != nil {
		return err
	}

	extrinsicSignature.Signer = signerPubKey
	extrinsicSignature.Signature = types.MultiSignature{IsSr25519: true, AsSr25519: sig}

	if err := checkDynamicExtrinsicData(meta, &extrinsicPayload, &extrinsicSignature); err != nil {
		return fmt.Errorf("dynamic extrinsic signature check failed: %w", err)
	}

	e.Signature = extrinsicSignature

	// mark the extrinsic as signed
	e.Version |= types.ExtrinsicBitSigned

	return nil
}

func (e DynamicExtrinsic) Encode(encoder scale.Encoder) error {
	if e.Type() != types.ExtrinsicVersion4 {
		return fmt.Errorf("unsupported extrinsic version: %v (isSigned: %v, type: %v)", e.Version, e.IsSigned(),
			e.Type())
	}
	var bb = bytes.Buffer{}
	tempEnc := scale.NewEncoder(&bb)

	err := tempEnc.Encode(e.Version)
	if err != nil {
		return err
	}

	if e.IsSigned() {
		err = tempEnc.Encode(e.Signature)
		if err != nil {
			return err
		}
	}

	// encode the method
	err = tempEnc.Encode(e.Method)
	if err != nil {
		return err
	}

	// take the temporary buffer to determine length, write that as prefix
	eb := bb.Bytes()
	err = encoder.EncodeUintCompact(*big.NewInt(0).SetUint64(uint64(len(eb))))
	if err != nil {
		return err
	}

	// write the actual encoded transaction
	err = encoder.Write(eb)
	if err != nil {
		return err
	}

	return nil
}
