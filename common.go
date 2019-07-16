package substrate

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/centrifuge/go-substrate-rpc-client/scalecodec"
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

func (a *Address) ParityDecode(decoder scalecodec.Decoder) {
	decoder.Read(a.PubKey[:])
}

func (a Address) ParityEncode(encoder scalecodec.Encoder) {
	// type of address - public key
	encoder.Write([]byte{255})
	encoder.Write(a.PubKey[:])
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

func (s *Signature) ParityDecode(decoder scalecodec.Decoder) {
	decoder.Read(s.Hash[:])
}

func (s Signature) ParityEncode(encoder scalecodec.Encoder) {
	encoder.Write(s.Hash[:])
}
