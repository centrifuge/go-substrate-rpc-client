package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/vimukthi-git/go-substrate/withreflect"
	"log"
	"os/exec"
	"strings"
)

type Address []byte

type Signature []byte

type MortalEra struct {
	Period uint64
	Phase uint64
}

type ExtrinsicSignature struct {
	Version uint8
	Signer Address
	Signature Signature
	Nonce uint32
	Era uint8 // era enum
	ImmortalEra []byte
	MortalEra MortalEra
}

func (e *ExtrinsicSignature) ParityEncode(encoder withreflection.Encoder) {
	panic("implement me")
}

type Args interface {
	withreflection.Encodeable
	//withreflection.Decodeable
}

type Method struct {

	CallIndex MethodIDX
	//  dynamic struct with the list of arguments defined as fields
	Args Args
}

func (m Method) ParityEncode(encoder withreflection.Encoder) {
	encoder.Encode(&m.CallIndex)
	encoder.Encode(m.Args)
}

type Extrinsic struct {
	Signature ExtrinsicSignature
	Method Method
}

func (e Extrinsic) Encode(encoder withreflection.Encoder) []byte {
	b := make([]byte, 0, 1000)
	bb := bytes.NewBuffer(b)
	tempEnc := withreflection.NewEncoder(bb)
	tempEnc.Encode(&e.Method)
	bbb := bb.Bytes()
	encoded := hex.EncodeToString(bbb)


	// use "subKey" command for signature
	//fmt.Println("method", hexutil.Encode(bb.Bytes()))
	/**
	subkey sign-transaction --call "070080de8d09fec089eec30b24e7c28859b627454641a057aeb9d4f583b36307c191e280e6fe3357f989cee5691e8393cfd7c3f870307bfbe3cd8fc8bcc561b8783dbcc18033d822a14e34e54d871208a5cef385b336a42da3b2ba729d71336d51df25f253"
	--nonce 1 --suri "//Alice" --password "" --prior-block-hash "bad7010fc5b729a383599112f97a08d7921a00f2550cb78c71dffb3080162d4d"
	 */
	out, err := exec.Command("/Users/vimukthi/.cargo/bin/subkey", "sign-transaction", "--call", encoded, "--nonce", "1", "--suri", "//Alice", "--password", "", "--prior-block-hash", "bad7010fc5b729a383599112f97a08d7921a00f2550cb78c71dffb3080162d4d").Output()
	if err != nil {
		log.Fatal(err.Error())
	}

	v := strings.TrimSpace(string(out))
	fmt.Println("out", v)

	dec, err := hexutil.Decode(v)
	if err != nil {
		log.Fatal(err.Error())
	}

	dec = append(dec, bbb...)
	//fmt.Println(dec)

	return dec
}

// MethodIDX [sectionIndex, methodIndex] 16bits
type MethodIDX struct {
	SectionIndex uint8
	MethodIndex uint8
}

func (m MethodIDX) ParityEncode(encoder withreflection.Encoder) {
	encoder.Encode(m.SectionIndex)
	encoder.Encode(m.MethodIndex)
}

type MetadataV4 struct {
	Modules []ModuleMetaData
}

func (m *MetadataV4) MethodIndex(method string) MethodIDX {
	s := strings.Split(method, ".")
	var sIDX, mIDX uint8 = 0, 0
	for i, n := range m.Modules {
		if n.Name == s[0] {
			sIDX = uint8(i)
			for j, f := range n.Calls {
				if f.Name == s[1] {
					mIDX = uint8(j)
				}
			}
		}
	}

	return MethodIDX{sIDX, mIDX}
}

func (m *MetadataV4) ParityDecode(decoder withreflection.Decoder) {
	decoder.Decode(&m.Modules)
}

type FunctionArgumentMetadata struct {
	Name string
	Type string
}

func (m *FunctionArgumentMetadata) ParityDecode(decoder withreflection.Decoder) {
	decoder.Decode(&m.Name)
	decoder.Decode(&m.Type)
}

type FunctionMetaData struct {
	Name string
	Args []FunctionArgumentMetadata
	Documentation []string
}

func (m *FunctionMetaData) ParityDecode(decoder withreflection.Decoder) {
	decoder.Decode(&m.Name)
	decoder.Decode(&m.Args)
	decoder.Decode(&m.Documentation)
}

type EventMetadata struct {
	Name string
	Args []string
	Documentation []string
}

func (m *EventMetadata) ParityDecode(decoder withreflection.Decoder) {
	decoder.Decode(&m.Name)
	decoder.Decode(&m.Args)
	decoder.Decode(&m.Documentation)
}

/**
[{"name":"AccountNonce","modifier":"Default","type":{"MapType":{"hasher":"Blake2_256","key":"AccountId","value":"Index","isLinked":false}},"fallback":"0x0000000000000000","documentation":[" Extrinsics nonce for accounts."]},{"name":"ExtrinsicCount","modifier":"Optional","type":{"PlainType":"u32"},"fallback":"0x00","documentation":[" Total extrinsics count for the current block."]},{"name":"AllExtrinsicsLen","modifier":"Optional","type":{"PlainType":"u32"},"fallback":"0x00","documentation":[" Total length in bytes for all extrinsics put together, for the current block."]},{"name":"BlockHash","modifier":"Default","type":{"MapType":{"hasher":"Blake2_256","key":"BlockNumber","value":"Hash","isLinked":false}},"fallback":"0x0000000000000000000000000000000000000000000000000000000000000000","documentation":[" Map of block numbers to block hashes."]},{"name":"ExtrinsicData","modifier":"Default","type":{"MapType":{"hasher":"Blake2_256","key":"u32","value":"Bytes","isLinked":false}},"fallback":"0x00","documentation":[" Extrinsics data for the current block (maps extrinsic's index to its data)."]},{"name":"RandomSeed","modifier":"Default","type":{"PlainType":"Hash"},"fallback":"0x0000000000000000000000000000000000000000000000000000000000000000","documentation":[" Random seed of the current block."]},{"name":"Number","modifier":"Default","type":{"PlainType":"BlockNumber"},"fallback":"0x0000000000000000","documentation":[" The current block number being processed. Set by `execute_block`."]},{"name":"ParentHash","modifier":"Default","type":{"PlainType":"Hash"},"fallback":"0x0000000000000000000000000000000000000000000000000000000000000000","documentation":[" Hash of the previous block."]},{"name":"ExtrinsicsRoot","modifier":"Default","type":{"PlainType":"Hash"},"fallback":"0x0000000000000000000000000000000000000000000000000000000000000000","documentation":[" Extrinsics root of the current block, also part of the block header."]},{"name":"Digest","modifier":"Default","type":{"PlainType":"Digest"},"fallback":"0x00","documentation":[" Digest of the current block, also part of the block header."]},{"name":"Events","modifier":"Default","type":{"PlainType":"Vec<EventRecord>"},"fallback":"0x00","documentation":[" Events deposited for the current block."]}]
 */

type TypMap struct {
	Hasher uint8
	Key string
	Value string
	IsLinked bool
}

func (m *TypMap) ParityDecode(decoder withreflection.Decoder) {
	decoder.Decode(&m.Hasher)
	decoder.Decode(&m.Key)
	decoder.Decode(&m.Value)
	decoder.Decode(&m.IsLinked)
}

type TypDoubleMap struct {
	Hasher uint8
	Key string
	Key2 string
	Value string
	Key2Hasher string
}

func (m *TypDoubleMap) ParityDecode(decoder withreflection.Decoder) {
	decoder.Decode(&m.Hasher)
	decoder.Decode(&m.Key)
	decoder.Decode(&m.Key2)
	decoder.Decode(&m.Value)
	decoder.Decode(&m.Key2Hasher)
}

type StorageFunctionMetadata struct {
	Name string
	Modifier uint8
	Type uint8
	Plane string
	Map TypMap
	DMap TypDoubleMap
	Fallback []byte
	Documentation []string
}

func (m *StorageFunctionMetadata) ParityDecode(decoder withreflection.Decoder) {
	decoder.Decode(&m.Name)
	decoder.Decode(&m.Modifier)
	decoder.Decode(&m.Type)
	switch m.Type {
	case 0:
		decoder.Decode(&m.Plane)
	case 1:
		decoder.Decode(&m.Map)
	default:
		decoder.Decode(&m.DMap)
	}
	decoder.Decode(&m.Fallback)
	decoder.Decode(&m.Documentation)
	fmt.Println(m.Documentation)
}

type ModuleMetaData struct {
	Name string
	Prefix string
	StorageOptional uint8
	Storage []StorageFunctionMetadata
	CallsOptional uint8
	Calls []FunctionMetaData
	EventsOptional uint8
	Events []EventMetadata
}
/**
_class [
  ModuleMetadata [Map] {
    'name' => [String (Text): 'system'],
    'prefix' => [String (Text): 'System'],
    'storage' => Option { raw: [_class] },
    'calls' => Option { raw: Null {} },
    'events' => Option { raw: [_class] },
    _jsonMap: Map {},
    _Types: {
      name: 'Text',
      prefix: 'Text',
      storage: 'Option',
      calls: 'Option',
      events: 'Option'
    }
  },
  ModuleMetadata [Map] {
    'name' => [String (Text): 'timestamp'],
    'prefix' => [String (Text): 'Timestamp'],
    'storage' => Option { raw: [_class] },
    'calls' => Option { raw: [_class] },
    'events' => Option { raw: Null {} },
    _jsonMap: Map {},
    _Types: {
      name: 'Text',
      prefix: 'Text',
      storage: 'Option',
      calls: 'Option',
      events: 'Option'
    }
  },
  ModuleMetadata [Map] {
    'name' => [String (Text): 'consensus'],
    'prefix' => [String (Text): 'Consensus'],
    'storage' => Option { raw: [_class] },
    'calls' => Option { raw: [_class] },
    'events' => Option { raw: Null {} },
    _jsonMap: Map {},
    _Types: {
      name: 'Text',
      prefix: 'Text',
      storage: 'Option',
      calls: 'Option',
      events: 'Option'
    }
  },
  ModuleMetadata [Map] {
    'name' => [String (Text): 'aura'],
    'prefix' => [String (Text): ''],
    'storage' => Option { raw: Null {} },
    'calls' => Option { raw: Null {} },
    'events' => Option { raw: Null {} },
    _jsonMap: Map {},
    _Types: {
      name: 'Text',
      prefix: 'Text',
      storage: 'Option',
      calls: 'Option',
      events: 'Option'
    }
  },
  ModuleMetadata [Map] {
    'name' => [String (Text): 'indices'],
    'prefix' => [String (Text): 'Indices'],
    'storage' => Option { raw: [_class] },
    'calls' => Option { raw: [_class] },
    'events' => Option { raw: [_class] },
    _jsonMap: Map {},
    _Types: {
      name: 'Text',
      prefix: 'Text',
      storage: 'Option',
      calls: 'Option',
      events: 'Option'
    }
  },
  ModuleMetadata [Map] {
    'name' => [String (Text): 'balances'],
    'prefix' => [String (Text): 'Balances'],
    'storage' => Option { raw: [_class] },
    'calls' => Option { raw: [_class] },
    'events' => Option { raw: [_class] },
    _jsonMap: Map {},
    _Types: {
      name: 'Text',
      prefix: 'Text',
      storage: 'Option',
      calls: 'Option',
      events: 'Option'
    }
  },
  ModuleMetadata [Map] {
    'name' => [String (Text): 'sudo'],
    'prefix' => [String (Text): 'Sudo'],
    'storage' => Option { raw: [_class] },
    'calls' => Option { raw: [_class] },
    'events' => Option { raw: [_class] },
    _jsonMap: Map {},
    _Types: {
      name: 'Text',
      prefix: 'Text',
      storage: 'Option',
      calls: 'Option',
      events: 'Option'
    }
  },
  ModuleMetadata [Map] {
    'name' => [String (Text): 'kerplunk'],
    'prefix' => [String (Text): 'Kerplunk'],
    'storage' => Option { raw: [_class] },
    'calls' => Option { raw: [_class] },
    'events' => Option { raw: [_class] },
    _jsonMap: Map {},
    _Types: {
      name: 'Text',
      prefix: 'Text',
      storage: 'Option',
      calls: 'Option',
      events: 'Option'
    }
  },
  _Type: [Function: ModuleMetadata]
]
 */
func (m *ModuleMetaData) ParityDecode(decoder withreflection.Decoder) {
	decoder.Decode(&m.Name)
	decoder.Decode(&m.Prefix)

	decoder.Decode(&m.StorageOptional)
	if m.StorageOptional == 1 {
		decoder.Decode(&m.Storage)
	}

	decoder.Decode(&m.CallsOptional)
	if m.CallsOptional == 1 {
		decoder.Decode(&m.Calls)
	}

	decoder.Decode(&m.EventsOptional)
	if m.EventsOptional == 1 {
		decoder.Decode(&m.Events)
		fmt.Println(m.Events)
	}
}

type MetadataVersioned struct {
	// 1635018093
	MagicNumber uint32
	Version       uint8
	Metadata MetadataV4
}

func (m *MetadataVersioned) ParityDecode(decoder withreflection.Decoder) {
	decoder.Decode(&m.MagicNumber)
	// we need to decide which struct to use based on the following number(enum), for now its hardcoded
	decoder.Decode(&m.Version)
	decoder.Decode(&m.Metadata)
}

type AnchorParams struct {
	AnchorIDPreimage []byte
	DocRoot []byte
	Proof []byte
}

func (a AnchorParams) ParityEncode(encoder withreflection.Encoder) {
	encoder.Encode(a.AnchorIDPreimage)
	encoder.Encode(a.DocRoot)
	encoder.Encode(a.Proof)
}

func main() {


	// Connect the client.
	client, err := rpc.Dial("ws://127.0.0.1:9944")
	if err != nil {
		panic(err)
	}

	// state_getRuntimeVersion
	var res string
	err = client.Call(&res, "state_getMetadata", "0xd133045f0efad58582772cbdb6f5f0cd6af7bb4bf1f30d039a4b18b4bdaf4901")
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)

	b, err := hexutil.Decode(res)
	if err != nil {
		panic(err)
	}

	dec := withreflection.NewDecoder(bytes.NewReader(b))

	//fmt.Println(res, b)
	v := uint32(0)
	vv := uint8(0)
	n := &MetadataVersioned{v, vv, MetadataV4{make([]ModuleMetaData, 10)}}
	dec.Decode(n)
	fmt.Println(res, n)

	//fmt.Println(n.Metadata.MethodIndex("kerplunk.commit"))
	//err = client.Call(&res, "state_queryStorage", "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY")
	//if err != nil {
	//	panic(err)
	//}


	e := Extrinsic{Method:Method{CallIndex:n.Metadata.MethodIndex("kerplunk.commit"), Args:&AnchorParams{utils.RandomSlice(32), utils.RandomSlice(32), utils.RandomSlice(32)}}}

	var bi []byte
	tempEnc := withreflection.NewEncoder(bytes.NewBuffer(bi))
	enc := e.Encode(*tempEnc)
	vs := hexutil.Encode(enc)
	fmt.Println(vs)
	err = client.Call(&res, "author_submitExtrinsic", vs)
	if err != nil {
		panic(err)
	}
}