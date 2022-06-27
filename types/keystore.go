package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

type KeyPurpose struct {
	IsP2PDiscovery bool

	IsP2PDocumentSigning bool
}

func (k *KeyPurpose) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	switch b {
	case 0:
		k.IsP2PDiscovery = true
	case 1:
		k.IsP2PDocumentSigning = true
	}

	return nil
}

func (k KeyPurpose) Encode(encoder scale.Encoder) error {
	switch {
	case k.IsP2PDiscovery:
		return encoder.PushByte(0)
	case k.IsP2PDocumentSigning:
		return encoder.PushByte(1)
	}

	return nil
}

type KeyType struct {
	IsECDSA bool

	IsEDDSA bool
}

func (k *KeyType) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	switch b {
	case 0:
		k.IsECDSA = true
	case 1:
		k.IsEDDSA = true
	}

	return nil
}

func (k KeyType) Encode(encoder scale.Encoder) error {
	switch {
	case k.IsECDSA:
		return encoder.PushByte(0)
	case k.IsEDDSA:
		return encoder.PushByte(1)
	}

	return nil
}

type Key struct {
	KeyPurpose KeyPurpose
	KeyType    KeyType
	RevokedAt  OptionBlockNumber
	Deposit    U128
}

type KeyID struct {
	Hash       Hash
	KeyPurpose KeyPurpose
}

type AddKey struct {
	Key     Hash
	Purpose KeyPurpose
	KeyType KeyType
}
