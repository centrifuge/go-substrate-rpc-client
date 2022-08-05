package types

import (
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
)

type EventKeystoreKeyAdded struct {
	Phase      Phase
	Owner      AccountID
	Key        Hash
	KeyPurpose KeyPurpose
	KeyType    KeyType
	Topics     []Hash
}

type EventKeystoreKeyRevoked struct {
	Phase       Phase
	Owner       AccountID
	Key         Hash
	BlockNumber BlockNumber
	Topics      []Hash
}

type EventKeystoreDepositSet struct {
	Phase      Phase
	NewDeposit U128
	Topics     []Hash
}

type KeyPurpose uint

const (
	KeyPurposeP2PDiscovery KeyPurpose = iota
	KeyPurposeP2PDocumentSigning
)

var (
	keyPurposeMap = map[KeyPurpose]struct{}{
		KeyPurposeP2PDiscovery:       {},
		KeyPurposeP2PDocumentSigning: {},
	}
)

func (k *KeyPurpose) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	kp := KeyPurpose(b)

	if _, ok := keyPurposeMap[kp]; !ok {
		return errors.New("unknown key purpose")
	}

	*k = kp

	return nil
}

func (k KeyPurpose) Encode(encoder scale.Encoder) error {
	return encoder.PushByte(byte(k))
}

type KeyType uint

const (
	KeyTypeECDSA KeyType = iota
	KeyTypeEDDSA
)

var (
	keyTypeMap = map[KeyType]struct{}{
		KeyTypeECDSA: {},
		KeyTypeEDDSA: {},
	}
)

func (k *KeyType) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	kt := KeyType(b)

	if _, ok := keyTypeMap[kt]; !ok {
		return errors.New("unknown key type")
	}

	*k = kt

	return nil
}

func (k KeyType) Encode(encoder scale.Encoder) error {
	return encoder.PushByte(byte(k))
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
