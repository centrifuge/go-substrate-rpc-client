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

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
)

type ClientOpts struct {
	Blockchain string
	ApiURL     string
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
		apiURL:     opts.ApiURL,
		http:       http.DefaultClient,
		wsURL:      opts.WsURL,
		sapi:       sapi,
	}, nil
}

type EventListResponseBody struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	GeneratedAt int    `json:"generated_at"`
	Data        struct {
		Count  int `json:"count"`
		Events []struct {
			EventIndex     string `json:"event_index"`
			BlockNum       int    `json:"block_num"`
			ExtrinsicIdx   int    `json:"extrinsic_idx"`
			ModuleId       string `json:"module_id"`
			EventId        string `json:"event_id"`
			Params         string `json:"params"`
			Phase          int    `json:"phase"`
			EventIdx       int    `json:"event_idx"`
			ExtrinsicHash  string `json:"extrinsic_hash"`
			Finalized      bool   `json:"finalized"`
			BlockTimestamp int    `json:"block_timestamp"`
		} `json:"events"`
	} `json:"data"`
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

	ApiURL string
	WsURL  string

	Meta        []byte
	StorageData []byte
}

func (c *Client) getTestData(blockNumber int) (*TestData, error) {
	meta, err := c.sapi.RPC.State.GetMetadataLatest()

	if err != nil {
		return nil, fmt.Errorf("couldn't get latest metadata: %w", err)
	}

	encodedMetadata, err := types.EncodeToBytes(meta)

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
		ApiURL:      c.apiURL,
		WsURL:       c.wsURL,
		Meta:        encodedMetadata,
		StorageData: *storageData,
	}, nil
}

const (
	maxBlockAge = 180 * 24 * time.Hour // 180 days
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
		return 0, errors.New("no events found")
	}

	firstEvent := r.Data.Events[0]

	blockTimestamp := time.Unix(int64(firstEvent.BlockTimestamp), 0)

	if time.Since(blockTimestamp) > maxBlockAge {
		return 0, fmt.Errorf("block with timestamp %s is too old", blockTimestamp.Format(time.RFC3339))
	}

	return firstEvent.BlockNum, nil
}

type EventListReqBody struct {
	Row    int    `json:"row"`
	Page   int    `json:"page"`
	Module string `json:"module"`
	Call   string `json:"call"`
}

func (c *Client) createRequest(ctx context.Context, reqData *ReqData) (*http.Request, error) {
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
