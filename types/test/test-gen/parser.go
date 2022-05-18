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
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	WsURLFormat  = "wss://%s.api.onfinality.io/public-ws"
	APIURLFormat = "https://%s.api.subscan.io/api/scan/events"

	BlockchainTagName = "test-gen-blockchain"
	WsURLTagName      = "test-gen-ws"
	APIURLTagName     = "test-gen-api"
	SkipTagName       = "test-gen-skip"
)

type FieldInfo struct {
	ReqData *ReqData

	ClientOpts *ClientOpts
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) CanSkip(field reflect.StructField) bool {
	return field.Tag.Get(SkipTagName) == "true"
}

func (p *Parser) ParseField(field reflect.StructField) (*FieldInfo, error) {
	clientOpts, err := p.parseClientOpts(field)

	if err != nil {
		return nil, fmt.Errorf("couldn't parse client opts: %w", err)
	}

	reqData, err := p.parseReqData(field)

	if err != nil {
		return nil, fmt.Errorf("couldn't parse client opts: %w", err)
	}

	return &FieldInfo{
		reqData,
		clientOpts,
	}, nil
}

var (
	errAPIURLMissing     = errors.New("api URL missing")
	errWsURLMissing      = errors.New("ws URL missing")
	errBlockchainMissing = errors.New("blockchain missing")
)

func (p *Parser) parseClientOpts(field reflect.StructField) (*ClientOpts, error) {
	wsURL := field.Tag.Get(WsURLTagName)
	apiURL := field.Tag.Get(APIURLTagName)

	if wsURL != "" && apiURL != "" {
		return &ClientOpts{WsURL: wsURL, APIURL: apiURL}, nil
	}

	if wsURL != "" {
		return nil, errAPIURLMissing
	}

	if apiURL != "" {
		return nil, errWsURLMissing
	}

	blockchain := field.Tag.Get(BlockchainTagName)

	if blockchain == "" {
		return nil, errBlockchainMissing
	}

	return &ClientOpts{
		Blockchain: blockchain,
		WsURL:      fmt.Sprintf(WsURLFormat, blockchain),
		APIURL:     fmt.Sprintf(APIURLFormat, blockchain),
	}, nil
}

const (
	fieldNameSeparator = "_"
	partsNumber        = 2
)

func (p *Parser) parseReqData(field reflect.StructField) (*ReqData, error) {
	fieldName := field.Name

	parts := strings.Split(fieldName, fieldNameSeparator)

	if plen := len(parts); plen != partsNumber {
		return nil, fmt.Errorf("expected %d parts, got %d", partsNumber, plen)
	}

	return &ReqData{
		Module: strings.ToLower(parts[0]),
		Call:   strings.ToLower(parts[1]),
	}, nil
}
