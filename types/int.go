package types

import (
	"encoding/json"
)

type I8 int8

func NewI8(i int8) I8 {
	return I8(i)
}

func (i *I8) UnmarshalJSON(b []byte) error {
	var tmp int8
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*i = I8(tmp)
	return nil
}

func (i *I8) MarshalJSON() ([]byte, error) {
	return json.Marshal(int8(*i))
}

type I16 int16

func NewI16(i int16) I16 {
	return I16(i)
}

func (i *I16) UnmarshalJSON(b []byte) error {
	var tmp int16
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*i = I16(tmp)
	return nil
}

func (i *I16) MarshalJSON() ([]byte, error) {
	return json.Marshal(int16(*i))
}

type I32 int32

func NewI32(i int32) I32 {
	return I32(i)
}

func (i *I32) UnmarshalJSON(b []byte) error {
	var tmp int32
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*i = I32(tmp)
	return nil
}

func (i *I32) MarshalJSON() ([]byte, error) {
	return json.Marshal(int32(*i))
}

type I64 int64

func NewI64(i int64) I64 {
	return I64(i)
}

func (i *I64) UnmarshalJSON(b []byte) error {
	var tmp int64
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*i = I64(tmp)
	return nil
}

func (i *I64) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(*i))
}
