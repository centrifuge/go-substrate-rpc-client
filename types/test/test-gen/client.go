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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
)

type ClientOpts struct {
	Blockchain string
	APIURL     string
	WsURL      string
}

type ReqData struct {
	Module string
	Call   string
}

type Client struct {
	blockchain string

	apiURL string
	http   *http.Client

	wsURL string
	sapi  *gsrpc.SubstrateAPI
}

func NewClient(opts ClientOpts) (*Client, error) {
	sapi, err := gsrpc.NewSubstrateAPI(opts.WsURL)

	if err != nil {
		return nil, fmt.Errorf("couldn't create substrate client: %w", err)
	}

	return &Client{
		blockchain: opts.Blockchain,
		apiURL:     opts.APIURL,
		http:       http.DefaultClient,
		wsURL:      opts.WsURL,
		sapi:       sapi,
	}, nil
}

type EventListResponseEvent struct {
	EventIndex     string `json:"event_index"`
	BlockNum       int    `json:"block_num"`
	ExtrinsicIdx   int    `json:"extrinsic_idx"`
	ModuleID       string `json:"module_id"`
	EventID        string `json:"event_id"`
	Params         string `json:"params"`
	Phase          int    `json:"phase"`
	EventIdx       int    `json:"event_idx"`
	ExtrinsicHash  string `json:"extrinsic_hash"`
	Finalized      bool   `json:"finalized"`
	BlockTimestamp int    `json:"block_timestamp"`
}

type EventListResponseData struct {
	Count  int                      `json:"count"`
	Events []EventListResponseEvent `json:"events"`
}

type EventListResponseBody struct {
	Code        int                   `json:"code"`
	Message     string                `json:"message"`
	GeneratedAt int                   `json:"generated_at"`
	Data        EventListResponseData `json:"data"`
}

func (c *Client) GetTestData(ctx context.Context, reqData *ReqData) (*TestData, error) {
	blockNumber, err := c.getBlockNumber(ctx, reqData)

	if err != nil {
		return nil, fmt.Errorf("couldn't get block number: %w", err)
	}

	return c.getTestData(blockNumber)
}

type TestData struct {
	Blockchain  string
	BlockNumber int

	APIURL string
	WsURL  string

	Meta        []byte
	StorageData []byte
}

func (c *Client) getTestData(blockNumber int) (*TestData, error) {
	meta, err := c.sapi.RPC.State.GetMetadataLatest()

	if err != nil {
		return nil, fmt.Errorf("couldn't get latest metadata: %w", err)
	}

	encodedMetadata, err := codec.Encode(meta)

	if err != nil {
		return nil, fmt.Errorf("couldn't encode metadata: %w", err)
	}

	key, err := types.CreateStorageKey(meta, "System", "Events", nil)

	if err != nil {
		return nil, fmt.Errorf("couldn't create storage key: %w", err)
	}

	bh, err := c.sapi.RPC.Chain.GetBlockHash(uint64(blockNumber))

	if err != nil {
		return nil, fmt.Errorf("couldn't get block hash '%d': %w", blockNumber, err)
	}

	storageData, err := c.sapi.RPC.State.GetStorageRaw(key, bh)

	if err != nil {
		return nil, fmt.Errorf("couldn't get raw storage data with key '%s': %w", key, err)
	}

	return &TestData{
		Blockchain:  c.blockchain,
		BlockNumber: blockNumber,
		APIURL:      c.apiURL,
		WsURL:       c.wsURL,
		Meta:        encodedMetadata,
		StorageData: *storageData,
	}, nil
}

const (
	maxBlockAge = 180 * 24 * time.Hour // 180 days
)

var (
	errNoEventsFound = errors.New("no events found")
	errBlockIsTooOld = errors.New("block is too old")
)

func (c *Client) getBlockNumber(ctx context.Context, reqData *ReqData) (int, error) {
	req, err := c.createRequest(ctx, reqData)

	if err != nil {
		return 0, fmt.Errorf("couldn't create request: %w", err)
	}

	res, err := c.http.Do(req)

	if err != nil {
		return 0, fmt.Errorf("couldn't perform request: %w", err)
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return 0, fmt.Errorf("couldn't read response body: %w", err)
	}

	var r EventListResponseBody

	if err := json.Unmarshal(b, &r); err != nil {
		return 0, fmt.Errorf("couldn't unmarshal response body: %w", err)
	}

	if len(r.Data.Events) == 0 {
		return 0, errNoEventsFound
	}

	firstEvent := r.Data.Events[0]

	blockTimestamp := time.Unix(int64(firstEvent.BlockTimestamp), 0)

	if time.Since(blockTimestamp) > maxBlockAge {
		return 0, errBlockIsTooOld
	}

	return firstEvent.BlockNum, nil
}

type EventListReqBody struct {
	Row    int    `json:"row"`
	Page   int    `json:"page"`
	Module string `json:"module"`
	Call   string `json:"call"`
}

var (
	errMissingModuleName = errors.New("missing module name")
	errMissingCallName   = errors.New("missing call name")
)

func (c *Client) createRequest(ctx context.Context, reqData *ReqData) (*http.Request, error) {
	if reqData.Module == "" {
		return nil, errMissingModuleName
	}

	if reqData.Call == "" {
		return nil, errMissingCallName
	}

	b, err := json.Marshal(EventListReqBody{
		Row:    1,
		Page:   0,
		Module: reqData.Module,
		Call:   reqData.Call,
	})

	if err != nil {
		return nil, fmt.Errorf("couldn't marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(b))

	if err != nil {
		return nil, fmt.Errorf("couldn't create request: %w", err)
	}

	return req, nil
}
