package parser

import libErr "github.com/centrifuge/go-substrate-rpc-client/v4/error"

const (
	ErrEventsCountDecoding  = libErr.Error("events count decoding")
	ErrEventPhaseDecoding   = libErr.Error("event phase decoding")
	ErrEventIDDecoding      = libErr.Error("event ID decoding")
	ErrEventDecoderNotFound = libErr.Error("event decoder not found")
	ErrEventFieldsDecoding  = libErr.Error("event fields decoding")
	ErrEventTopicsDecoding  = libErr.Error("event topics decoding")
	ErrCallDecoderNotFound  = libErr.Error("call decoder not found")
	ErrCallFieldsDecoding   = libErr.Error("call fields decoding")
)
