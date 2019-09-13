package types

import (
	"encoding/json"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

type RuntimeVersion struct {
	APIs             []RuntimeVersionAPI `json:"apis"`
	AuthoringVersion U32                 `json:"authoringVersion"`
	ImplName         string              `json:"implName"`
	ImplVersion      U32                 `json:"implVersion"`
	SpecName         string              `json:"specName"`
	SpecVersion      U32                 `json:"specVersion"`
}

func NewRuntimeVersion() *RuntimeVersion {
	return &RuntimeVersion{APIs: make([]RuntimeVersionAPI, 0)}
}

func (r *RuntimeVersion) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&r.APIs)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.AuthoringVersion)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.ImplName)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.ImplVersion)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.SpecName)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.SpecVersion)
	if err != nil {
		return err
	}

	return nil
}

func (r RuntimeVersion) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(r.APIs)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.AuthoringVersion)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.ImplName)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.ImplVersion)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.SpecName)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.SpecVersion)
	if err != nil {
		return err
	}

	return nil
}

type RuntimeVersionAPI struct {
	APIID   string
	Version U32
}

func (r *RuntimeVersionAPI) UnmarshalJSON(b []byte) error {
	tmp := []interface{}{&r.APIID, &r.Version}
	wantLen := len(tmp)
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Notification: %d != %d", g, e)
	}
	return nil
}

func (r *RuntimeVersionAPI) MarshalJSON() ([]byte, error) {
	tmp := []interface{}{r.APIID, r.Version}
	return json.Marshal(tmp)
}

func (r *RuntimeVersionAPI) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&r.APIID)
	if err != nil {
		return err
	}

	err = decoder.Decode(&r.Version)
	if err != nil {
		return err
	}

	return nil
}

func (r RuntimeVersionAPI) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(r.APIID)
	if err != nil {
		return err
	}

	err = encoder.Encode(r.Version)
	if err != nil {
		return err
	}

	return nil
}
