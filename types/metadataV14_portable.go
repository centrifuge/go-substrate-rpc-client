package types

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
)

type PortableTypeV14 struct {
	ID   Si1LookupTypeID
	Type Si1Type
}

func (d *PortableTypeV14) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&d.ID)
	if err != nil {
		return fmt.Errorf("decode Si1LookupTypeID error: %v", err)
	}

	return decoder.Decode(&d.Type)
}

func (x PortableTypeV14) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(x.ID)
	if err != nil {
		return err
	}

	return encoder.Encode(x.Type)
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

func (d Si0TypeDefPrimitive) Encode(encoder scale.Encoder) error {
	switch d.Value {
	case "Bool":
		return encoder.PushByte(0)
	case "Char":
		return encoder.PushByte(1)
	case "Str":
		return encoder.PushByte(2)
	case "U8":
		return encoder.PushByte(3)
	case "U16":
		return encoder.PushByte(4)
	case "U32":
		return encoder.PushByte(5)
	case "U64":
		return encoder.PushByte(6)
	case "U128":
		return encoder.PushByte(7)
	case "U256":
		return encoder.PushByte(8)
	case "I8":
		return encoder.PushByte(9)
	case "I16":
		return encoder.PushByte(10)
	case "I32":
		return encoder.PushByte(11)
	case "I64":
		return encoder.PushByte(12)
	case "I128":
		return encoder.PushByte(13)
	case "I256":
		return encoder.PushByte(14)
	default:
		//TODO(nuno): Not sure what to do
		return nil
	}
}

//------------------v1-----------

type Si1LookupTypeID struct {
	UCompact
}

func NewSi1LookupTypeID(value *big.Int) Si1LookupTypeID {
	return Si1LookupTypeID{NewUCompact(value)}
}

func NewSi1LookupTypeIDFromUInt(value uint64) Si1LookupTypeID {
	return NewSi1LookupTypeID(new(big.Int).SetUint64(value))
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
	Name    Text
	HasType bool
	Type    Si1LookupTypeID
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
	d.HasType = hasValue
	if !d.HasType {
		d.Type = NewSi1LookupTypeID(big.NewInt(0))
	}
	return nil
}

func (d Si1TypeParameter) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(d.Name)
	if err != nil {
		return err
	}

	return encoder.EncodeOption(d.HasType, &d.Type)

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
		return fmt.Errorf("Si1TypeDef unknow type : %d", num)
	}
}

func (d Si1TypeDef) Encode(encoder scale.Encoder) error {
	switch {
	case d.IsComposite:
		err := encoder.PushByte(0)
		if err != nil {
			return err
		}
		return encoder.Encode(&d.Composite)
	case d.IsVariant:
		err := encoder.PushByte(1)
		if err != nil {
			return err
		}
		return encoder.Encode(&d.Variant)
	case d.IsSequence:
		err := encoder.PushByte(2)
		if err != nil {
			return err
		}
		return encoder.Encode(&d.Sequence)
	case d.IsArray:
		err := encoder.PushByte(3)
		if err != nil {
			return err
		}
		return encoder.Encode(&d.Array)
	case d.IsTuple:
		err := encoder.PushByte(4)
		if err != nil {
			return err
		}
		return encoder.Encode(&d.Tuple)
	case d.IsPrimitive:
		err := encoder.PushByte(5)
		if err != nil {
			return err
		}
		return encoder.Encode(&d.Primitive)
	case d.IsCompact:
		err := encoder.PushByte(6)
		if err != nil {
			return err
		}
		return encoder.Encode(&d.Compact)
	case d.IsBitSequence:
		err := encoder.PushByte(7)
		if err != nil {
			return err
		}
		return encoder.Encode(&d.BitSequence)
	case d.IsHistoricMetaCompat:
		err := encoder.PushByte(8)
		if err != nil {
			return err
		}
		d.IsHistoricMetaCompat = true
		return encoder.Encode(&d.HistoricMetaCompat)

	default:
		return errors.New("expected Si1TypeDef instance to be one of the valid variants")
	}
}

type Si1TypeDefComposite struct {
	Fields []Si1Field
}

func (d *Si1TypeDefComposite) Decode(decoder scale.Decoder) error {
	return decoder.Decode(&d.Fields)
}

type Si1Field struct {
	HasName     bool
	Name        Text
	Type        Si1LookupTypeID
	HasTypeName bool
	TypeName    Text
	Docs        []Text
}

func (d *Si1Field) Decode(decoder scale.Decoder) error {
	var hasValue bool
	err := decoder.DecodeOption(&hasValue, &d.Name)
	if err != nil {
		return err
	}
	d.HasName = hasValue

	err = decoder.Decode(&d.Type)
	if err != nil {
		return err
	}

	err = decoder.DecodeOption(&hasValue, &d.TypeName)
	if err != nil {
		return err
	}

	d.HasTypeName = hasValue

	return decoder.Decode(&d.Docs)
}

func (d Si1Field) Encode(encoder scale.Encoder) error {
	// TODO(nuno): may need to handle optional Name and TypeName
	err := encoder.EncodeOption(d.HasName, d.Name)
	if err != nil {
		return err
	}
	err = encoder.Encode(d.Type)
	if err != nil {
		return err
	}
	err = encoder.EncodeOption(d.HasTypeName, d.TypeName)
	if err != nil {
		return err
	}
	return encoder.Encode(&d.Docs)
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
