package types

import (
	"encoding/json"
	"github.com/ChainSafe/gossamer/codec"
	"strconv"
)

type U8 uint8

func NewU8(u uint8) U8 {
	return U8(u)
}

func (u *U8) UnmarshalJSON(b []byte) error {
	var tmp uint8
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U8(tmp)
	return nil
}

func (u *U8) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint8(*u))
}

func (u U8) EncodedLength() (int, error) {
	return encodedLength(u)
}

func (u U8) Hash() ([32]byte, error) {
	return hash(u)
}

func (u U8) IsEmpty() bool {
	return u == 0
}

func (u U8) Encode() ([]byte, error) {
	return codec.Encode(u)
}

func (u U8) Hex() (string, error) {
	return _hex(u)
}

func (u U8) String() string {
	return strconv.FormatUint(uint64(u), 10)
}

func (u U8) Eq(o Codec) bool {
	if ov, ok := o.(U8); ok {
		return u == ov
	}
	return false
}

type U16 uint16

func NewU16(u uint16) U16 {
	return U16(u)
}

func (u *U16) UnmarshalJSON(b []byte) error {
	var tmp uint16
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U16(tmp)
	return nil
}

func (u *U16) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint16(*u))
}

func (u U16) EncodedLength() (int, error) {
	return encodedLength(u)
}

func (u U16) Hash() ([32]byte, error) {
	return hash(u)
}

func (u U16) IsEmpty() bool {
	return u == 0
}

func (u U16) Encode() ([]byte, error) {
	return codec.Encode(u)
}

func (u U16) Hex() (string, error) {
	return _hex(u)
}

func (u U16) String() string {
	return strconv.FormatUint(uint64(u), 10)
}

func (u U16) Eq(o Codec) bool {
	if ov, ok := o.(U16); ok {
		return u == ov
	}
	return false
}

type U32 uint32

func NewU32(u uint32) U32 {
	return U32(u)
}

func (u *U32) UnmarshalJSON(b []byte) error {
	var tmp uint32
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U32(tmp)
	return nil
}

func (u *U32) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint32(*u))
}

func (u U32) EncodedLength() (int, error) {
	return encodedLength(u)
}

func (u U32) Hash() ([32]byte, error) {
	return hash(u)
}

func (u U32) IsEmpty() bool {
	return u == 0
}

func (u U32) Encode() ([]byte, error) {
	return codec.Encode(u)
}

func (u U32) Hex() (string, error) {
	return _hex(u)
}

func (u U32) String() string {
	return strconv.FormatUint(uint64(u), 10)
}

func (u U32) Eq(o Codec) bool {
	if ov, ok := o.(U32); ok {
		return u == ov
	}
	return false
}

type U64 uint64

func NewU64(u uint64) U64 {
	return U64(u)
}

func (u *U64) UnmarshalJSON(b []byte) error {
	var tmp uint64
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U64(tmp)
	return nil
}

func (u *U64) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint64(*u))
}

func (u U64) EncodedLength() (int, error) {
	return encodedLength(u)
}

func (u U64) Hash() ([32]byte, error) {
	return hash(u)
}

func (u U64) IsEmpty() bool {
	return u == 0
}

func (u U64) Encode() ([]byte, error) {
	return codec.Encode(u)
}

func (u U64) Hex() (string, error) {
	return _hex(u)
}

func (u U64) String() string {
	return strconv.FormatUint(uint64(u), 10)
}

func (u U64) Eq(o Codec) bool {
	if ov, ok := o.(U64); ok {
		return u == ov
	}
	return false
}
