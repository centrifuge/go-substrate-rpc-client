package substrate

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

// Hash is 256 bit by default
type Hash []byte

func (h *Hash) String() string {
	b := *h
	return hexutil.Encode(b[:])
}

/**
const PREFIX_1BYTE = 0xef;
const PREFIX_2BYTE = 0xfc;
const PREFIX_4BYTE = 0xfd;
const PREFIX_8BYTE = 0xfe;
 */
type Address struct {
	PubKey [32]byte
}

func NewAddress(b []byte) *Address {
	s := &Address{}
	copy(s.PubKey[:], b)
	return s
}

func (a *Address) Decode(decoder scale.Decoder) error {
	err := decoder.Read(a.PubKey[:])
	if err != nil {
		return err
	}
	return nil
}

func (a Address) Encode(encoder scale.Encoder) error {
	// type of address - public key
	err := encoder.Write([]byte{255})
	if err != nil {
		return err
	}

	err = encoder.Write(a.PubKey[:])
	if err != nil {
		return err
	}
	return nil
}

type Index uint64

type Signature struct {
	Hash [64]byte
}

func NewSignature(b []byte) *Signature {
	s := &Signature{}
	copy(s.Hash[:], b)
	return s
}

func (s *Signature) Decode(decoder scale.Decoder) error {
	err := decoder.Read(s.Hash[:])
	if err != nil {
		return err
	}

	return nil
}

func (s Signature) Encode(encoder scale.Encoder) error {
	err := encoder.Write(s.Hash[:])
	if err != nil {
		return err
	}
	return nil
}
