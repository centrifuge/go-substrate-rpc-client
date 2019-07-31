package substrate

import (
	"bytes"
	"errors"
	"fmt"
	"hash"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/minio/blake2b-simd"
	"github.com/pierrec/xxHash/xxHash64"
)

// MethodIDX [sectionIndex, methodIndex] 16bits
type MethodIDX struct {
	SectionIndex uint8
	MethodIndex  uint8
}

func (e *MethodIDX) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&e.SectionIndex)
	if err != nil {
		return err
	}

	err = decoder.Decode(&e.MethodIndex)
	if err != nil {
		return err
	}

	return nil
}

func (m MethodIDX) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.SectionIndex)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.MethodIndex)
	if err != nil {
		return err
	}

	return nil
}

type MetadataV4 struct {
	Modules []ModuleMetaData
}

func (m *MetadataV4) MethodIndex(method string) MethodIDX {
	s := strings.Split(method, ".")
	var sIDX, mIDX uint8 = 0, 0
	// section index
	var sCounter = 0

	for _, n := range m.Modules {
		if n.CallsOptional == 1 {
			if n.Name == s[0] {
				sIDX = uint8(sCounter)
				for j, f := range n.Calls {
					if f.Name == s[1] {
						mIDX = uint8(j)
					}
				}
			}
			sCounter++
		}
	}

	return MethodIDX{sIDX, mIDX}
}

func (m *MetadataV4) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Modules)
	if err != nil {
		return err
	}
	return nil
}

type FunctionArgumentMetadata struct {
	Name string
	Type string
}

func (m *FunctionArgumentMetadata) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Type)
	if err != nil {
		return err
	}

	return nil
}

type FunctionMetaData struct {
	Name          string
	Args          []FunctionArgumentMetadata
	Documentation []string
}

func (m *FunctionMetaData) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Args)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Documentation)
	if err != nil {
		return err
	}

	return nil
}

type EventMetadata struct {
	Name          string
	Args          []string
	Documentation []string
}

func (m *EventMetadata) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Args)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Documentation)
	if err != nil {
		return err
	}

	return nil
}

/**
[{"name":"AccountNonce","modifier":"Default","type":{"MapType":{"hasher":"Blake2_256","key":"AccountId","value":"Index","isLinked":false}},"fallback":"0x0000000000000000","documentation":[" Extrinsics nonce for accounts."]},{"name":"ExtrinsicCount","modifier":"Optional","type":{"PlainType":"u32"},"fallback":"0x00","documentation":[" Total extrinsics count for the current block."]},{"name":"AllExtrinsicsLen","modifier":"Optional","type":{"PlainType":"u32"},"fallback":"0x00","documentation":[" Total length in bytes for all extrinsics put together, for the current block."]},{"name":"BlockHash","modifier":"Default","type":{"MapType":{"hasher":"Blake2_256","key":"BlockNumber","value":"Hash","isLinked":false}},"fallback":"0x0000000000000000000000000000000000000000000000000000000000000000","documentation":[" Map of block numbers to block hashes."]},{"name":"ExtrinsicData","modifier":"Default","type":{"MapType":{"hasher":"Blake2_256","key":"u32","value":"Bytes","isLinked":false}},"fallback":"0x00","documentation":[" Extrinsics data for the current block (maps extrinsic's index to its data)."]},{"name":"RandomSeed","modifier":"Default","type":{"PlainType":"Hash"},"fallback":"0x0000000000000000000000000000000000000000000000000000000000000000","documentation":[" Random seed of the current block."]},{"name":"Number","modifier":"Default","type":{"PlainType":"BlockNumber"},"fallback":"0x0000000000000000","documentation":[" The current block number being processed. Set by `execute_block`."]},{"name":"ParentHash","modifier":"Default","type":{"PlainType":"Hash"},"fallback":"0x0000000000000000000000000000000000000000000000000000000000000000","documentation":[" Hash of the previous block."]},{"name":"ExtrinsicsRoot","modifier":"Default","type":{"PlainType":"Hash"},"fallback":"0x0000000000000000000000000000000000000000000000000000000000000000","documentation":[" Extrinsics root of the current block, also part of the block header."]},{"name":"Digest","modifier":"Default","type":{"PlainType":"Digest"},"fallback":"0x00","documentation":[" Digest of the current block, also part of the block header."]},{"name":"Events","modifier":"Default","type":{"PlainType":"Vec<EventRecord>"},"fallback":"0x00","documentation":[" Events deposited for the current block."]}]
*/

type TypMap struct {
	Hasher   uint8
	Key      string
	Value    string
	IsLinked bool
}

func (t TypMap) HashFunc() (hash.Hash, error) {
	if t.Hasher == 1 {
		return blake2b.New(&blake2b.Config{Size: 32})
	}
	return nil, errors.New("hash function type not supported")
}

func (m *TypMap) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Hasher)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Key)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Value)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.IsLinked)
	if err != nil {
		return err
	}

	return nil
}

type TypDoubleMap struct {
	Hasher     uint8
	Key        string
	Key2       string
	Value      string
	Key2Hasher string
}

func (m *TypDoubleMap) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Hasher)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Key)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Key2)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Value)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Key2Hasher)
	if err != nil {
		return err
	}

	return nil
}

type StorageFunctionMetadata struct {
	Name          string
	Modifier      uint8
	Type          uint8
	Plane         string
	Map           TypMap
	DMap          TypDoubleMap
	Fallback      []byte
	Documentation []string
}

func (s StorageFunctionMetadata) isMap() bool {
	return s.Type == 1
}

func (s StorageFunctionMetadata) isDMap() bool {
	return s.Type == 2
}

func (m *StorageFunctionMetadata) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Modifier)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Type)
	if err != nil {
		return err
	}

	switch m.Type {
	case 0:
		err = decoder.Decode(&m.Plane)
		if err != nil {
			return err
		}
	case 1:
		err = decoder.Decode(&m.Map)
		if err != nil {
			return err
		}
	default:
		err = decoder.Decode(&m.DMap)
		if err != nil {
			return err
		}
	}
	err = decoder.Decode(&m.Fallback)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Documentation)
	if err != nil {
		return err
	}
	// fmt.Println(metadataVersioned.Documentation)
	return nil
}

type ModuleMetaData struct {
	Name            string
	Prefix          string
	StorageOptional uint8
	Storage         []StorageFunctionMetadata
	CallsOptional   uint8
	Calls           []FunctionMetaData
	EventsOptional  uint8
	Events          []EventMetadata
}

func (m *ModuleMetaData) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Prefix)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.StorageOptional)
	if err != nil {
		return err
	}

	if m.StorageOptional == 1 {
		err = decoder.Decode(&m.Storage)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.CallsOptional)
	if err != nil {
		return err
	}

	if m.CallsOptional == 1 {
		err = decoder.Decode(&m.Calls)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.EventsOptional)
	if err != nil {
		return err
	}

	if m.EventsOptional == 1 {
		err = decoder.Decode(&m.Events)
		if err != nil {
			return err
		}
		// fmt.Println(metadataVersioned.Events)
	}
	return nil
}

// MetadataVersioned only supports v4
type MetadataVersioned struct {
	// 1635018093
	MagicNumber uint32
	Version     uint8
	Metadata    MetadataV4
}

func NewMetadataVersioned() *MetadataVersioned {
	return &MetadataVersioned{Metadata: MetadataV4{make([]ModuleMetaData, 0)}}
}

func (m *MetadataVersioned) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.MagicNumber)
	if err != nil {
		return err
	}
	// we need to decide which struct to use based on the following number(enum), for now its hardcoded
	err = decoder.Decode(&m.Version)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Metadata)
	if err != nil {
		return err
	}

	return nil
}

type State struct {
	client    Client
}

func NewStateRPC(client Client) *State {
	return &State{client: client}
}

func (s *State) MetaData(blockHash Hash) (*MetadataVersioned, error) {
	var res string
	// block hash can give error - Error(Client(UnknownBlock("State already discarded for Hash(0xxxx)")), State { next_error: None, backtrace: InternalBacktrace { backtrace: None } })
	var err error
	if blockHash == nil {
		err = s.client.Call(&res, "state_getMetadata")
	} else {
		err = s.client.Call(&res, "state_getMetadata", blockHash.String())
	}
	if err != nil {
		return nil, err
	}

	b, err := hexutil.Decode(res)
	if err != nil {
		return nil, err
	}

	dec := scale.NewDecoder(bytes.NewReader(b))
	n := NewMetadataVersioned()
	err = dec.Decode(n)
	if err != nil {
		return nil, err
	}

	return n, nil
}

type StorageKey []byte

func NewStorageKey(meta MetadataVersioned, module string, fn string, key []byte) (StorageKey, error) {
	var fnMeta *StorageFunctionMetadata
	for _, m := range meta.Metadata.Modules {
		if m.Prefix == module {
			for _, s := range m.Storage {
				if s.Name == fn {
					fnMeta = &s
					break
				}
			}
		}
	}
	if fnMeta == nil {
		return nil, fmt.Errorf("no meta data found for module %s function %s", module, fn)
	}

	var hasher hash.Hash
	var err error
	if fnMeta.isMap() {
		hasher, err = fnMeta.Map.HashFunc()
		if err != nil {
			return nil, err
		}
	} else if fnMeta.isDMap() {
		// TODO define hashing for 2 keys
	}

	afn := []byte(module + " " + fn)
	// TODO why is add length prefix step in JS client doesn't add anything to the hashed key?
	if hasher != nil {
		hasher.Write(append(afn, key...))
		return hasher.Sum(nil), nil
	} else {
		if key != nil {
			return createMultiXxhash(append(afn, key...), 2), nil
		}
		return createMultiXxhash(append(afn), 2), nil
	}
}

func (s StorageKey) Encode(encoder scale.Encoder) error {
	return encoder.Encode(s)
}

type StorageData []byte

func (s StorageData) Decoder() *scale.Decoder {
	buf := bytes.NewBuffer(s[:])
	return scale.NewDecoder(buf)
}

func (s *State) Storage(key StorageKey, block []byte) (StorageData, error) {
	var res string
	var err error
	if block != nil {
		err = s.client.Call(&res, "state_getStorage", hexutil.Encode(key), hexutil.Encode(block))
	} else {
		err = s.client.Call(&res, "state_getStorage", hexutil.Encode(key))
	}

	if err != nil {
		return nil, err
	}

	if res == "" {
		return nil, errors.New("empty result")
	}

	return hexutil.Decode(res)
}

func createMultiXxhash(data []byte, rounds int) []byte {
	res := make([]byte, 0)
	for i := 0; i < rounds; i++ {
		h := xxHash64.New(uint64(i))
		h.Write(data)
		res = append(res, h.Sum(nil)...)
	}
	return res
}
