package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	WsURLFormat  = "wss://%s.api.onfinality.io/public-ws"
	ApiURLFormat = "https://%s.api.subscan.io/api/scan/events"

	BlockchainTagName = "test-gen-blockchain"
	WsURLTagName      = "test-gen-ws"
	ApiURLTagName     = "test-gen-api"
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

func (p *Parser) parseClientOpts(field reflect.StructField) (*ClientOpts, error) {
	wsURL := field.Tag.Get(WsURLTagName)
	apiURL := field.Tag.Get(ApiURLTagName)

	if wsURL != "" && apiURL != "" {
		return &ClientOpts{WsURL: wsURL, ApiURL: apiURL}, nil
	}

	if wsURL != "" {
		return nil, errors.New("api URL missing")
	}

	if apiURL != "" {
		return nil, errors.New("ws URL missing")
	}

	blockchain := field.Tag.Get(BlockchainTagName)

	if blockchain == "" {
		return nil, errors.New("no blockchain provided")
	}

	return &ClientOpts{
		Blockchain: blockchain,
		WsURL:      fmt.Sprintf(WsURLFormat, blockchain),
		ApiURL:     fmt.Sprintf(ApiURLFormat, blockchain),
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
