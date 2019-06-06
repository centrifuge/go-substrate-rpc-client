package main

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/vimukthi-git/go-substrate/withreflect"
)

type MetadataV4 struct {
	Modules []ModuleMetaData
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

func main() {


	// Connect the client.
	client, err := rpc.Dial("ws://127.0.0.1:9944")
	if err != nil {
		panic(err)
	}

	// state_getRuntimeVersion
	var res string
	err = client.Call(&res, "state_getMetadata", "0x602f513ba2fca731f934c93110045de627826413e4cd37304fa989821c820287")
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

}