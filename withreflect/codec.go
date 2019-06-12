// Copyright 2018 Jsgenesis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package withreflection

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

// Implementation for Parity codec in Go.
// Derived from https://github.com/paritytech/parity-codec/
// While Rust implementation uses Rust type system and is highly optimized, this one
// has to rely on Go's reflection and thus is notably slower.
// Feature parity is almost full, apart from the lack of support for u128 (which are missing in Go).

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Encoder is a wrapper around a Writer that allows encoding data items to a stream.
type Encoder struct {
	writer io.Writer
}

func NewEncoder(writer io.Writer) *Encoder {
	return &Encoder{writer: writer}
}

// Write several bytes to the encoder.
func (pe Encoder) Write(bytes []byte) {
	c, err := pe.writer.Write(bytes)
	check(err)
	if c < len(bytes) {
		panic(fmt.Sprintf("Could not write %d bytes to writer", len(bytes)))
	}
}

// PushByte writes a single byte to an encoder.
func (pe Encoder) PushByte(b byte) {
	pe.Write([]byte{b})
}

// EncodeUintCompact writes an unsigned integer to the stream using the compact encoding.
// A typical usage is storing the length of a collection.
// Definition of compact encoding:
// 0b00 00 00 00 / 00 00 00 00 / 00 00 00 00 / 00 00 00 00
//   xx xx xx 00															(0 ... 2**6 - 1)		(u8)
//   yL yL yL 01 / yH yH yH yL												(2**6 ... 2**14 - 1)	(u8, u16)  low LH high
//   zL zL zL 10 / zM zM zM zL / zM zM zM zM / zH zH zH zM					(2**14 ... 2**30 - 1)	(u16, u32)  low LMMH high
//   nn nn nn 11 [ / zz zz zz zz ]{4 + n}									(2**30 ... 2**536 - 1)	(u32, u64, u128, U256, U512, U520) straight LE-encoded
// Rust implementation: see impl<'a> Encode for CompactRef<'a, u64>
func (pe Encoder) EncodeUintCompact(v uint64) {

	// TODO: handle numbers wide than 64 bits (byte slices?)
	// Currently, Rust implementation only seems to support u128

	if v < 1<<30 {
		if v < 1<<6 {
			pe.PushByte(byte(v) << 2)
		} else if v < 1<<14 {
			err := binary.Write(pe.writer, binary.LittleEndian, uint16(v<<2)+1)
			if err != nil {
				panic(err)
			}
		} else {
			err := binary.Write(pe.writer, binary.LittleEndian, uint32(v<<2)+2)
			if err != nil {
				panic(err)
			}
		}
		return
	}

	n := byte(0)
	limit := uint64(1 << 32)
	for v >= limit && limit > 256 { // when overflows, limit will be < 256
		n++
		limit <<= 8
	}
	if n > 4 {
		panic("Assertion error: n>4 needed to compact-encode uint64")
	}
	pe.PushByte((n << 2) + 3)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, v)
	pe.Write(buf[:4+n])
}

// Encode a value to the stream.
func (pe Encoder) Encode(value interface{}) {
	t := reflect.TypeOf(value)
	tk := t.Kind()
	switch tk {

	// Boolean and numbers are trivially encoded via binary.Write
	// It will use reflection again and take a performance hit
	// TODO: consider handling every case directly
	case reflect.Bool:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uintptr:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		err := binary.Write(pe.writer, binary.LittleEndian, value)
		if err != nil {
			panic(err)
		}
	case reflect.Ptr:
		rv := reflect.ValueOf(value)
		if rv.IsNil() {
			panic("Encoding null pointers not supported; consider using Option type")
		} else {
			dereferenced := rv.Elem()
			pe.Encode(dereferenced.Interface())
		}

	// Arrays and slices: first compact-encode length, then each item individually
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		rv := reflect.ValueOf(value)
		len := rv.Len()
		len64 := uint64(len)
		if len64 > math.MaxUint32 {
			panic("Attempted to serialize a collection with too many elements.")
		}
		pe.EncodeUintCompact(len64)
		for i := 0; i < len; i++ {
			pe.Encode(rv.Index(i).Interface())
		}

	// Strings are encoded as UTF-8 byte slices, just as in Rust
	case reflect.String:
		pe.Encode([]byte(value.(string)))

	case reflect.Struct:
		encodeable := reflect.TypeOf((*Encodeable)(nil)).Elem()
		if t.Implements(encodeable) {
			value.(Encodeable).ParityEncode(pe)
		} else {
			panic(fmt.Sprintf("Type %s does not support Encodeable interface", t))
		}

	// Currently unsupported types
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.UnsafePointer:
		fallthrough
	case reflect.Invalid:
		panic(fmt.Sprintf("Type %s cannot be encoded", t.Kind()))
	default:
		fmt.Println("not captured")
	}
}

// EncodeOption stores optionally present value to the stream.
func (pe Encoder) EncodeOption(hasValue bool, value interface{}) {
	if !hasValue {
		pe.PushByte(0)
	} else {
		pe.PushByte(1)
		pe.Encode(value)
	}
}

// Decoder - a wraper around a Reader that allows decoding data items from a stream.
// Unlike Rust implementations, decoder methods do not return success state, but just
// panic on error. Since decoding failue is an "unexpected" error, this approach should
// be justified.
type Decoder struct {
	reader io.Reader
}

func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{reader:reader}
}

// Read reads bytes from a stream into a buffer and panics if cannot read the required
// number of bytes.
func (pd Decoder) Read(bytes []byte) {
	c, err := pd.reader.Read(bytes)
	check(err)
	if c < len(bytes) {
		panic(fmt.Sprintf("Cannot read the required number of bytes %d, only %d available", len(bytes), c))
	}
}

// ReadOneByte reads a next byte from the stream.
// Named so to avoid a linter warning about a clash with io.ByteReader.ReadByte
func (pd Decoder) ReadOneByte() byte {
	buf := []byte{0}
	pd.Read(buf)
	return buf[0]
}

// Decode takes a pointer to a decodable value and populates it from the stream.
func (pd Decoder) Decode(target interface{}) {
	t0 := reflect.TypeOf(target)
	if t0.Kind() != reflect.Ptr {
		panic("Target must be a pointer, but was " + fmt.Sprint(t0))
	}
	val := reflect.ValueOf(target)
	if val.IsNil() {
		panic("Target is a nil pointer")
	}
	pd.DecodeIntoReflectValue(val.Elem())
}

// DecodeIntoReflectValue populates a writable reflect.Value from the stream
func (pd Decoder) DecodeIntoReflectValue(target reflect.Value) {
	t := target.Type()
	if !target.CanSet() {
		panic("Unsettable value " + fmt.Sprint(t))
	}

	switch t.Kind() {

	// Boolean and numbers are trivially decoded via binary.Read
	// It will use reflection again and take a performance hit
	// TODO: consider handling every case directly
	case reflect.Bool:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uintptr:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		intHolder := reflect.New(t)
		intPointer := intHolder.Interface()
		binary.Read(pd.reader, binary.LittleEndian, intPointer)
		target.Set(intHolder.Elem())

	// Pointers are encoded just as the referenced types, with panicking on nil.
	// If you want to replicate Option<T> behavior in Rust, see OptionBool and an
	// example type OptionInt8 in tests.
	case reflect.Ptr:
		pd.DecodeIntoReflectValue(target.Elem())

	// Arrays and slices: first compact-encode length, then each item individually
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		codedLen64 := pd.DecodeUintCompact()
		if codedLen64 > math.MaxUint32 {
			panic("Encoded array length is higher than allowed by the protocol (32-bit unsigned integer)")
		}
		if codedLen64 > uint64(maxInt) {
			panic("Encoded array length is higher than allowed by the platform")
		}
		codedLen := int(codedLen64)
		targetLen := target.Len()
		if codedLen != targetLen {
			if t.Kind() == reflect.Array {
				panic(fmt.Sprintf(
					"We want to decode an array of length %d, but the encoded length is %d",
					target.Len(), codedLen))
			}
			if t.Kind() == reflect.Slice {
				if int(codedLen) > target.Cap() {
					newSlice := reflect.MakeSlice(t, int(codedLen), int(codedLen))
					target.Set(newSlice)
				} else {
					target.SetLen(int(codedLen))
				}
			}
		}
		for i := 0; i < codedLen; i++ {
			pd.DecodeIntoReflectValue(target.Index(i))
		}

	// Strings are encoded as UTF-8 byte slices, just as in Rust
	case reflect.String:
		var bytes []byte
		pd.Decode(&bytes)
		target.SetString(string(bytes))

	case reflect.Struct:
		encodeable := reflect.TypeOf((*Decodeable)(nil)).Elem()
		ptrType := reflect.PtrTo(t)
		if ptrType.Implements(encodeable) {
			ptrVal := reflect.New(t)
			ptrVal.Interface().(Decodeable).ParityDecode(pd)
			target.Set(ptrVal.Elem())
		} else {
			panic(fmt.Sprintf("Type %s does not support Decodeable interface", ptrType))
		}

	// Currently unsupported types
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.UnsafePointer:
		fallthrough
	case reflect.Invalid:
		panic(fmt.Sprintf("Type %s cannot be decoded", t.Kind()))
	}
}

// DecodeUintCompact decodes a compact-encoded integer. See EncodeUintCompact method.
func (pd Decoder) DecodeUintCompact() uint64 {
	b := pd.ReadOneByte()
	mode := b & 3
	switch mode {
	case 0:
		// right shift to remove mode bits
		return uint64(b >> 2)
	case 1:
		r := uint64(pd.ReadOneByte())
		// * 2^6
		r <<= 6
		// right shift to remove mode bits and add to prev
		r += uint64(b >> 2)
		return r
	case 2:
		// value = 32 bits + mode
		buf := make([]byte, 4)
		buf[0] = b
		pd.Read(buf[1:4])
		// set the buffer in little endian order
		r := binary.LittleEndian.Uint32(buf)
		// remove the last 2 mode bits
		r >>= 2
		return uint64(r)
	case 3:
		// remove mode bits
		l := b >> 2
		if l > 4 {
			panic("Not supported: l>4 encountered when decoding a compact-encoded uint")
		}
		buf := make([]byte, 8)
		pd.Read(buf[:l+4])
		return binary.LittleEndian.Uint64(buf)
	default:
		panic("Code should be unreachable")
	}
}

// DecodeOption decodes a optionally available value into a boolean presence field and a value.
func (pd Decoder) DecodeOption(hasValue *bool, valuePointer interface{}) {
	b := pd.ReadOneByte()
	switch b {
	case 0:
		*hasValue = false
	case 1:
		*hasValue = true
		pd.Decode(valuePointer)
	default:
		panic(fmt.Sprintf("Unknown byte prefix for encoded OptionBool: %d", b))
	}
}

// Encodeable is an interface that defines a custom encoding rules for a data type.
// Should be defined for structs (not pointers to them).
// See OptionBool for an example implementation.
type Encodeable interface {
	// ParityEncode encodes and write this structure into a stream
	ParityEncode(encoder Encoder)
}

// Decodeable is an interface that defines a custom encoding rules for a data type.
// Should be defined for pointers to structs.
// See OptionBool for an example implementation.
type Decodeable interface {
	// ParityDecode populates this structure from a stream (overwriting the current contents), return false on failure
	ParityDecode(decoder Decoder)
}

// OptionBool is a structure that can store a boolean or a missing value.
// Note that encoding rules are slightly different from other "Option" fields.
type OptionBool struct {
	hasValue bool
	value    bool
}

// NewOptionBoolEmpty creates an OptionBool without a value.
func NewOptionBoolEmpty() OptionBool {
	return OptionBool{false, false}
}

// NewOptionBool creates an OptionBool with a value.
func NewOptionBool(value bool) OptionBool {
	return OptionBool{true, value}
}

// ParityEncode implements encoding for OptionBool as per Rust implementation.
func (o OptionBool) ParityEncode(encoder Encoder) {
	if !o.hasValue {
		encoder.PushByte(0)
	} else {
		if o.value {
			encoder.PushByte(1)
		} else {
			encoder.PushByte(2)
		}
	}
}

// ParityDecode implements decoding for OptionBool as per Rust implementation.
func (o *OptionBool) ParityDecode(decoder Decoder) {
	b := decoder.ReadOneByte()
	switch b {
	case 0:
		o.hasValue = false
		o.value = false
	case 1:
		o.hasValue = true
		o.value = true
	case 2:
		o.hasValue = true
		o.value = false
	default:
		panic(fmt.Sprintf("Unknown byte prefix for encoded OptionBool: %d", b))
	}
}

// ToKeyedVec replicates the behaviour of Rust's to_keyed_vec helper.
func ToKeyedVec(value interface{}, prependKey []byte) []byte {
	var buffer = bytes.NewBuffer(prependKey)
	Encoder{buffer}.Encode(value)
	return buffer.Bytes()
}
