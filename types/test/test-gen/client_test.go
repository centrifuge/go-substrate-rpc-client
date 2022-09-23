// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
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

package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	mockClient "github.com/centrifuge/go-substrate-rpc-client/v4/client/mocks"
	mockChain "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/chain/mocks"
	mockState "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state/mocks"

	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpcmocksrv"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

type mockSrv struct{}

func (m *mockSrv) GetMetadata(hash *string) string {
	return types.ExamplaryMetadataV4String
}

func TestClient_NewClient(t *testing.T) {
	opts := ClientOpts{}

	c, err := NewClient(opts)
	assert.Error(t, err)
	assert.Nil(t, c)

	s := rpcmocksrv.New()
	defer s.Stop()

	err = s.RegisterName("state", &mockSrv{})
	assert.NoError(t, err)

	opts = ClientOpts{
		WsURL: s.URL,
	}

	c, err = NewClient(opts)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestClient_GetTestData(t *testing.T) {
	chainMock := mockChain.NewChain(t)
	stateMock := mockState.NewState(t)
	clientMock := mockClient.NewClient(t)

	moduleName := "test-module"
	callName := "test-call"

	blockNumber := 123

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)

		var rb EventListReqBody

		err = json.Unmarshal(body, &rb)
		assert.Nil(t, err)

		assert.Equal(t, rb.Module, moduleName)
		assert.Equal(t, rb.Call, callName)
		assert.Equal(t, rb.Page, 0)
		assert.Equal(t, rb.Row, 1)

		res := EventListResponseBody{
			Data: EventListResponseData{
				Count: 1,
				Events: []EventListResponseEvent{
					{
						BlockNum:       blockNumber,
						BlockTimestamp: int(time.Now().Unix()),
					},
				},
			},
		}

		resBytes, err := json.Marshal(res)

		assert.Nil(t, err)

		_, err = w.Write(resBytes)

		assert.Nil(t, err)
	}))

	defer ts.Close()

	var metadata types.Metadata
	err := codec.DecodeFromHex(types.MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.Nil(t, err)

	encodedMeta, err := codec.Encode(metadata)

	assert.Nil(t, err)

	stateMock.On("GetMetadataLatest").Return(&metadata, nil)

	hash := types.NewHash([]byte{0, 1, 2})

	chainMock.On("GetBlockHash", uint64(blockNumber)).Return(hash, nil)

	key, err := types.CreateStorageKey(&metadata, "System", "Events", nil)

	assert.Nil(t, err)

	storageData := types.StorageDataRaw([]byte{2, 3, 4})

	stateMock.On("GetStorageRaw", key, hash).Return(&storageData, nil)

	sapi := &gsrpc.SubstrateAPI{
		RPC: &rpc.RPC{
			Chain: chainMock,
			State: stateMock,
		},
		Client: clientMock,
	}

	client := &Client{
		apiURL: ts.URL,
		http:   http.DefaultClient,
		sapi:   sapi,
	}

	testData, err := client.GetTestData(context.Background(), &ReqData{
		Module: moduleName,
		Call:   callName,
	})

	assert.Nil(t, err)

	assert.Equal(t, testData.StorageData, []byte(storageData))
	assert.Equal(t, testData.Meta, encodedMeta)
}

func TestClient_GetTestData_InvalidRequestData(t *testing.T) {
	client := &Client{}

	_, err := client.GetTestData(context.Background(), &ReqData{
		Call: "test",
	})

	assert.True(t, errors.Is(err, errMissingModuleName))

	_, err = client.GetTestData(context.Background(), &ReqData{
		Module: "test",
	})

	assert.True(t, errors.Is(err, errMissingCallName))
}

func TestClient_GetTestData_HttpError(t *testing.T) {
	client := &Client{
		apiURL: "http://localhost/non-there",
		http:   http.DefaultClient,
	}

	_, err := client.GetTestData(context.Background(), &ReqData{
		Module: "moduleName",
		Call:   "callName",
	})

	var netErr *net.OpError
	assert.True(t, errors.As(err, &netErr))
}

func TestClient_GetTestData_NoResBody(t *testing.T) {
	moduleName := "test-module"
	callName := "test-call"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)

		var rb EventListReqBody

		err = json.Unmarshal(body, &rb)
		assert.Nil(t, err)

		assert.Equal(t, rb.Module, moduleName)
		assert.Equal(t, rb.Call, callName)
		assert.Equal(t, rb.Page, 0)
		assert.Equal(t, rb.Row, 1)
	}))

	defer ts.Close()

	client := &Client{
		apiURL: ts.URL,
		http:   http.DefaultClient,
	}

	_, err := client.GetTestData(context.Background(), &ReqData{
		Module: moduleName,
		Call:   callName,
	})

	var syntaxErr *json.SyntaxError
	assert.True(t, errors.As(err, &syntaxErr))
}

func TestClient_GetTestData_NoEvents(t *testing.T) {
	moduleName := "test-module"
	callName := "test-call"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)

		var rb EventListReqBody

		err = json.Unmarshal(body, &rb)
		assert.Nil(t, err)

		assert.Equal(t, rb.Module, moduleName)
		assert.Equal(t, rb.Call, callName)
		assert.Equal(t, rb.Page, 0)
		assert.Equal(t, rb.Row, 1)

		res := EventListResponseBody{
			Data: EventListResponseData{
				Count:  1,
				Events: []EventListResponseEvent{},
			},
		}

		resBytes, err := json.Marshal(res)

		assert.Nil(t, err)

		_, err = w.Write(resBytes)

		assert.Nil(t, err)
	}))

	defer ts.Close()

	client := &Client{
		apiURL: ts.URL,
		http:   http.DefaultClient,
	}

	_, err := client.GetTestData(context.Background(), &ReqData{
		Module: moduleName,
		Call:   callName,
	})

	assert.True(t, errors.Is(err, errNoEventsFound))
}

func TestClient_GetTestData_BlockTooOld(t *testing.T) {
	moduleName := "test-module"
	callName := "test-call"

	blockNumber := 123

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)

		var rb EventListReqBody

		err = json.Unmarshal(body, &rb)
		assert.Nil(t, err)

		assert.Equal(t, rb.Module, moduleName)
		assert.Equal(t, rb.Call, callName)
		assert.Equal(t, rb.Page, 0)
		assert.Equal(t, rb.Row, 1)

		res := EventListResponseBody{
			Data: EventListResponseData{
				Count: 1,
				Events: []EventListResponseEvent{
					{
						BlockNum:       blockNumber,
						BlockTimestamp: int(time.Now().Unix() - maxBlockAge.Milliseconds()),
					},
				},
			},
		}

		resBytes, err := json.Marshal(res)

		assert.Nil(t, err)

		_, err = w.Write(resBytes)

		assert.Nil(t, err)
	}))

	defer ts.Close()

	client := &Client{
		apiURL: ts.URL,
		http:   http.DefaultClient,
	}

	_, err := client.GetTestData(context.Background(), &ReqData{
		Module: moduleName,
		Call:   callName,
	})

	assert.True(t, errors.Is(err, errBlockIsTooOld))
}

func TestClient_GetTestData_MetadataError(t *testing.T) {
	chainMock := mockChain.NewChain(t)
	stateMock := mockState.NewState(t)
	clientMock := mockClient.NewClient(t)

	moduleName := "test-module"
	callName := "test-call"

	blockNumber := 123

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)

		var rb EventListReqBody

		err = json.Unmarshal(body, &rb)
		assert.Nil(t, err)

		assert.Equal(t, rb.Module, moduleName)
		assert.Equal(t, rb.Call, callName)
		assert.Equal(t, rb.Page, 0)
		assert.Equal(t, rb.Row, 1)

		res := EventListResponseBody{
			Data: EventListResponseData{
				Count: 1,
				Events: []EventListResponseEvent{
					{
						BlockNum:       blockNumber,
						BlockTimestamp: int(time.Now().Unix()),
					},
				},
			},
		}

		resBytes, err := json.Marshal(res)

		assert.Nil(t, err)

		_, err = w.Write(resBytes)

		assert.Nil(t, err)
	}))

	defer ts.Close()

	metaErr := errors.New("metadata error")

	stateMock.On("GetMetadataLatest").Return(nil, metaErr)

	sapi := &gsrpc.SubstrateAPI{
		RPC: &rpc.RPC{
			Chain: chainMock,
			State: stateMock,
		},
		Client: clientMock,
	}

	client := &Client{
		apiURL: ts.URL,
		http:   http.DefaultClient,
		sapi:   sapi,
	}

	_, err := client.GetTestData(context.Background(), &ReqData{
		Module: moduleName,
		Call:   callName,
	})

	assert.True(t, errors.Is(err, metaErr))
}

func TestClient_GetTestData_InvalidMetadata(t *testing.T) {
	chainMock := mockChain.NewChain(t)
	stateMock := mockState.NewState(t)
	clientMock := mockClient.NewClient(t)

	moduleName := "test-module"
	callName := "test-call"

	blockNumber := 123

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)

		var rb EventListReqBody

		err = json.Unmarshal(body, &rb)
		assert.Nil(t, err)

		assert.Equal(t, rb.Module, moduleName)
		assert.Equal(t, rb.Call, callName)
		assert.Equal(t, rb.Page, 0)
		assert.Equal(t, rb.Row, 1)

		res := EventListResponseBody{
			Data: EventListResponseData{
				Count: 1,
				Events: []EventListResponseEvent{
					{
						BlockNum:       blockNumber,
						BlockTimestamp: int(time.Now().Unix()),
					},
				},
			},
		}

		resBytes, err := json.Marshal(res)

		assert.Nil(t, err)

		_, err = w.Write(resBytes)

		assert.Nil(t, err)
	}))

	defer ts.Close()

	var metadata types.Metadata

	stateMock.On("GetMetadataLatest").Return(&metadata, nil)

	sapi := &gsrpc.SubstrateAPI{
		RPC: &rpc.RPC{
			Chain: chainMock,
			State: stateMock,
		},
		Client: clientMock,
	}

	client := &Client{
		apiURL: ts.URL,
		http:   http.DefaultClient,
		sapi:   sapi,
	}

	testData, err := client.GetTestData(context.Background(), &ReqData{
		Module: moduleName,
		Call:   callName,
	})

	assert.NotNil(t, err)
	assert.Nil(t, testData)
}

func TestClient_GetTestData_StorageKeyCreationError(t *testing.T) {
	chainMock := mockChain.NewChain(t)
	stateMock := mockState.NewState(t)
	clientMock := mockClient.NewClient(t)

	moduleName := "test-module"
	callName := "test-call"

	blockNumber := 123

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)

		var rb EventListReqBody

		err = json.Unmarshal(body, &rb)
		assert.Nil(t, err)

		assert.Equal(t, rb.Module, moduleName)
		assert.Equal(t, rb.Call, callName)
		assert.Equal(t, rb.Page, 0)
		assert.Equal(t, rb.Row, 1)

		res := EventListResponseBody{
			Data: EventListResponseData{
				Count: 1,
				Events: []EventListResponseEvent{
					{
						BlockNum:       blockNumber,
						BlockTimestamp: int(time.Now().Unix()),
					},
				},
			},
		}

		resBytes, err := json.Marshal(res)

		assert.Nil(t, err)

		_, err = w.Write(resBytes)

		assert.Nil(t, err)
	}))

	defer ts.Close()

	var metadata types.Metadata
	err := codec.DecodeFromHex(types.MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.Nil(t, err)

	metadata.AsMetadataV14.Pallets = nil

	stateMock.On("GetMetadataLatest").Return(&metadata, nil)

	sapi := &gsrpc.SubstrateAPI{
		RPC: &rpc.RPC{
			Chain: chainMock,
			State: stateMock,
		},
		Client: clientMock,
	}

	client := &Client{
		apiURL: ts.URL,
		http:   http.DefaultClient,
		sapi:   sapi,
	}

	testData, err := client.GetTestData(context.Background(), &ReqData{
		Module: moduleName,
		Call:   callName,
	})
	assert.NotNil(t, err)
	assert.Nil(t, testData)
}

func TestClient_GetTestData_BlockHashError(t *testing.T) {
	chainMock := mockChain.NewChain(t)
	stateMock := mockState.NewState(t)
	clientMock := mockClient.NewClient(t)

	moduleName := "test-module"
	callName := "test-call"

	blockNumber := 123

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)

		var rb EventListReqBody

		err = json.Unmarshal(body, &rb)
		assert.Nil(t, err)

		assert.Equal(t, rb.Module, moduleName)
		assert.Equal(t, rb.Call, callName)
		assert.Equal(t, rb.Page, 0)
		assert.Equal(t, rb.Row, 1)

		res := EventListResponseBody{
			Data: EventListResponseData{
				Count: 1,
				Events: []EventListResponseEvent{
					{
						BlockNum:       blockNumber,
						BlockTimestamp: int(time.Now().Unix()),
					},
				},
			},
		}

		resBytes, err := json.Marshal(res)

		assert.Nil(t, err)

		_, err = w.Write(resBytes)

		assert.Nil(t, err)
	}))

	defer ts.Close()

	var metadata types.Metadata
	err := codec.DecodeFromHex(types.MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.Nil(t, err)

	stateMock.On("GetMetadataLatest").Return(&metadata, nil)

	blockHashErr := errors.New("block hash err")

	chainMock.On("GetBlockHash", uint64(blockNumber)).Return(nil, blockHashErr)

	sapi := &gsrpc.SubstrateAPI{
		RPC: &rpc.RPC{
			Chain: chainMock,
			State: stateMock,
		},
		Client: clientMock,
	}

	client := &Client{
		apiURL: ts.URL,
		http:   http.DefaultClient,
		sapi:   sapi,
	}

	_, err = client.GetTestData(context.Background(), &ReqData{
		Module: moduleName,
		Call:   callName,
	})

	assert.True(t, errors.Is(err, blockHashErr))
}

func TestClient_GetTestData_StorageError(t *testing.T) {
	chainMock := mockChain.NewChain(t)
	stateMock := mockState.NewState(t)
	clientMock := mockClient.NewClient(t)

	moduleName := "test-module"
	callName := "test-call"

	blockNumber := 123

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)

		var rb EventListReqBody

		err = json.Unmarshal(body, &rb)
		assert.Nil(t, err)

		assert.Equal(t, rb.Module, moduleName)
		assert.Equal(t, rb.Call, callName)
		assert.Equal(t, rb.Page, 0)
		assert.Equal(t, rb.Row, 1)

		res := EventListResponseBody{
			Data: EventListResponseData{
				Count: 1,
				Events: []EventListResponseEvent{
					{
						BlockNum:       blockNumber,
						BlockTimestamp: int(time.Now().Unix()),
					},
				},
			},
		}

		resBytes, err := json.Marshal(res)

		assert.Nil(t, err)

		_, err = w.Write(resBytes)

		assert.Nil(t, err)
	}))

	defer ts.Close()

	var metadata types.Metadata
	err := codec.DecodeFromHex(types.MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.Nil(t, err)

	stateMock.On("GetMetadataLatest").Return(&metadata, nil)

	hash := types.NewHash([]byte{0, 1, 2})

	chainMock.On("GetBlockHash", uint64(blockNumber)).Return(hash, nil)

	key, err := types.CreateStorageKey(&metadata, "System", "Events", nil)

	assert.Nil(t, err)

	storageErr := errors.New("storage error")

	stateMock.On("GetStorageRaw", key, hash).Return(nil, storageErr)

	sapi := &gsrpc.SubstrateAPI{
		RPC: &rpc.RPC{
			Chain: chainMock,
			State: stateMock,
		},
		Client: clientMock,
	}

	client := &Client{
		apiURL: ts.URL,
		http:   http.DefaultClient,
		sapi:   sapi,
	}

	_, err = client.GetTestData(context.Background(), &ReqData{
		Module: moduleName,
		Call:   callName,
	})

	assert.True(t, errors.Is(err, storageErr))
}
