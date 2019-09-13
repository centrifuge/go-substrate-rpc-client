//nolint:dupl
package types

import (
	"encoding/json"
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
