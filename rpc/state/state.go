package state

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/client"
	"io/ioutil"
	"net/http"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

type State struct {
	client *client.Client
}

func NewState(c *client.Client) *State {
	return &State{c}
}

func (s *State) GetMetadataLatest() (types.Metadata, error) {
	requestBody := bytes.NewBuffer([]byte("{\"id\":1, \"jsonrpc\":\"2.0\", \"method\": \"state_getMetadata\", \"params\":[]}"))

	resp, err := http.Post("http://localhost:9933", "application/json", requestBody)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	resp.Body.Close()

	type RpcResp struct {
		JSONRPC string `json:"jsonrpc"`
		Result  string `json:"result"`
		ID      int    `json:"id"`
	}

	var rpcResp RpcResp

	err = json.Unmarshal(body, &rpcResp)
	if err != nil {
		panic(err)
	}

	bz, err := hex.DecodeString(rpcResp.Result[2:])

	if err != nil {
		panic(err)
	}

	decoder := scale.NewDecoder(bytes.NewReader(bz))

	metadata := types.NewMetadata()

	err = decoder.Decode(metadata)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(metadata)

	// with Gossamer codec
	// res, err := codec.Decode(body, types.Metadata{})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res)

	// var res string

	// with Ethereum Go Client
	// err := (*s.client).Call(&res, "state_getMetadata")
	// if err != nil {
	// 	panic(err)
	// 	return types.Metadata{}, err
	// }

	// fmt.Println("res", res)

	// bz, err := hex.DecodeString(res[2:])
	// if err != nil {
	// 	return types.Hash{}, err
	// }

	// fmt.Println("bz", bz)

	// if len(bz) != 32 {
	// 	return types.Hash{}, fmt.Errorf("Required result to be 32 bytes, but got %v", len(bz))
	// }

	// var bz32 [32]byte
	// copy(bz32[:], bz)

	// fmt.Println("bz32", bz32)

	// hash := types.NewHash(bz32)

	// return hash, nil
	return types.Metadata{}, err
}
