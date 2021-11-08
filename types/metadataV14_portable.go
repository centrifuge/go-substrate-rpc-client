package types

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
)

type PortableTypeV14 struct {
	ID   Si1LookupTypeID
	Type Si1Type

	// TODO(nuno): complete following https://tinyurl.com/zzznj3wd
}

func (d *PortableTypeV14) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&d.ID)
	if err != nil {
		return fmt.Errorf("decode Si1LookupTypeID error: %v", err)
	}

	return decoder.Decode(&d.Type)
}

//----------------v0------------

type Si0LookupTypeID UCompact

type Si0Path []Text

type Si0TypeDefPrimitive struct {
	Value string
}

func (d *Si0TypeDefPrimitive) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch b {
	case 0:
		d.Value = "Bool"
	case 1:
		d.Value = "Char"
	case 2:
		d.Value = "Str"
	case 3:
		d.Value = "U8"
	case 4:
		d.Value = "U16"
	case 5:
		d.Value = "U32"
	case 6:
		d.Value = "U64"
	case 7:
		d.Value = "U128"
	case 8:
		d.Value = "U256"
	case 9:
		d.Value = "I8"
	case 10:
		d.Value = "I16"
	case 11:
		d.Value = "I32"
	case 12:
		d.Value = "I64"
	case 13:
		d.Value = "I128"
	case 14:
		d.Value = "I256"
	default:
		return fmt.Errorf("Si0TypeDefPrimitive do not support this type: %d", b)
	}
	return nil
}

//------------------v1-----------

type Si1LookupTypeID big.Int

func NewSi1LookupTypeID(value *big.Int) Si1LookupTypeID {
	return Si1LookupTypeID(*value)
}

func NewSi1LookupTypeIDFromUInt(value uint64) Si1LookupTypeID {
	return NewSi1LookupTypeID(new(big.Int).SetUint64(value))
}
func (d *Si1LookupTypeID) Int64() int64 {
	i := big.Int(*d)
	return i.Int64()
}

// func (d *Si1LookupTypeID) UnmarshalJSON(bytes []byte) error {
// 	var s int64
// 	if err := json.Unmarshal(bytes, &s); err != nil {
// 		return err
// 	}
// 	dd := NewSi1LookupTypeIDFromUInt(uint64(s))
// 	d = &dd
// 	return nil
// }

func (d Si1LookupTypeID) MarshalJSON() ([]byte, error) {
	s := d.Int64()
	return json.Marshal(s)
}

func (d *Si1LookupTypeID) Decode(decoder scale.Decoder) error {
	ui, err := decoder.DecodeUintCompact()
	if err != nil {
		return err
	}

	*d = Si1LookupTypeID(*ui)
	return nil
}

func (d Si1LookupTypeID) Encode(encoder scale.Encoder) error {
	err := encoder.EncodeUintCompact(big.Int(d))
	if err != nil {
		return err
	}
	return nil
}

type Si1Path Si0Path

type Si1Type struct {
	Path   Si1Path
	Params []Si1TypeParameter
	Def    Si1TypeDef
	Docs   []Text
}

func (d *Si1Type) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&d.Path)
	if err != nil {
		return err
	}
	err = decoder.Decode(&d.Params)
	if err != nil {
		return err
	}
	err = decoder.Decode(&d.Def)
	if err != nil {
		return err
	}
	return decoder.Decode(&d.Docs)
}

type Si1TypeParameter struct {
	Name Text
	Type Si1LookupTypeID
}

func (d *Si1TypeParameter) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&d.Name)
	if err != nil {
		return err
	}
	var hasValue bool
	err = decoder.DecodeOption(&hasValue, &d.Type)
	if err != nil {
		return err
	}
	if !hasValue {
		d.Type = NewSi1LookupTypeID(big.NewInt(0))
	}
	return nil
}

type Si1TypeDef struct {
	IsComposite          bool
	Composite            Si1TypeDefComposite
	IsVariant            bool
	Variant              Si1TypeDefVariant
	IsSequence           bool
	Sequence             Si1TypeDefSequence
	IsArray              bool
	Array                Si1TypeDefArray
	IsTuple              bool
	Tuple                Si1TypeDefTuple
	IsPrimitive          bool
	Primitive            Si1TypeDefPrimitive
	IsCompact            bool
	Compact              Si1TypeDefCompact
	IsBitSequence        bool
	BitSequence          Si1TypeDefBitSequence
	IsHistoricMetaCompat bool
	HistoricMetaCompat   Type
}

func (d *Si1TypeDef) Decode(decoder scale.Decoder) error {
	num, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch num {
	case 0:
		d.IsComposite = true
		return decoder.Decode(&d.Composite)
	case 1:
		d.IsVariant = true
		return decoder.Decode(&d.Variant)
	case 2:
		d.IsSequence = true
		return decoder.Decode(&d.Sequence)
	case 3:
		d.IsArray = true
		return decoder.Decode(&d.Array)
	case 4:
		d.IsTuple = true
		return decoder.Decode(&d.Tuple)
	case 5:
		d.IsPrimitive = true
		return decoder.Decode(&d.Primitive)
	case 6:
		d.IsCompact = true
		return decoder.Decode(&d.Compact)
	case 7:
		d.IsBitSequence = true
		return decoder.Decode(&d.BitSequence)
	case 8:
		d.IsHistoricMetaCompat = true
		return decoder.Decode(&d.HistoricMetaCompat)

	default:
		return fmt.Errorf("Si1TypeDef un know type : %d", num)
	}
}

func (d *Si1TypeDef) GetSi1TypeDefData() {

}

type Si1TypeDefComposite struct {
	Fields []Si1Field
}

func (d *Si1TypeDefComposite) Decode(decoder scale.Decoder) error {
	return decoder.Decode(&d.Fields)
}

type Si1Field struct {
	Name     Text
	Type     Si1LookupTypeID
	TypeName Text
	Docs     []Text
}

func (d *Si1Field) Decode(decoder scale.Decoder) error {
	var hasValue bool
	err := decoder.DecodeOption(&hasValue, &d.Name)
	if err != nil {
		return err
	}
	err = decoder.Decode(&d.Type)
	if err != nil {
		return err
	}
	err = decoder.DecodeOption(&hasValue, &d.TypeName)
	if err != nil {
		return err
	}
	return decoder.Decode(&d.Docs)
}

type Si1TypeDefVariant struct {
	Variants []Si1Variant `json:"variants"`
}

func (d *Si1TypeDefVariant) Decode(decoder scale.Decoder) error {
	return decoder.Decode(&d.Variants)
}

type Si1Variant struct {
	Name   Text       `json:"name"`
	Fields []Si1Field `json:"fields"`
	Index  U8         `json:"index"`
	Docs   []Text     `json:"docs"`
}

func (d *Si1Variant) Decode(decoder scale.Decoder) error {
	var err error
	err = decoder.Decode(&d.Name)
	if err != nil {
		return err
	}
	err = decoder.Decode(&d.Fields)
	if err != nil {
		return err
	}
	err = decoder.Decode(&d.Index)
	if err != nil {
		return err
	}
	err = decoder.Decode(&d.Docs)
	if err != nil {
		return err
	}
	return nil
}

type Si1TypeDefSequence struct {
	Type Si1LookupTypeID
}

func (d *Si1TypeDefSequence) Decode(decoder scale.Decoder) error {
	return decoder.Decode(&d.Type)
}

type Si1TypeDefArray struct {
	Len  U32
	Type Si1LookupTypeID
}

func (d *Si1TypeDefArray) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&d.Len)
	if err != nil {
		return err
	}
	return decoder.Decode(&d.Type)
}

type Si1TypeDefTuple []Si1LookupTypeID

type Si1TypeDefPrimitive struct {
	Si0TypeDefPrimitive
}

type Si1TypeDefCompact struct {
	Type Si1LookupTypeID
}

func (d *Si1TypeDefCompact) Decode(decoder scale.Decoder) error {
	return decoder.Decode(&d.Type)
}

type Si1TypeDefBitSequence struct {
	BitStoreType Si1LookupTypeID
	BitOrderType Si1LookupTypeID
}

func (d *Si1TypeDefBitSequence) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&d.BitStoreType)
	if err != nil {
		return err
	}
	return decoder.Decode(&d.BitOrderType)
}
