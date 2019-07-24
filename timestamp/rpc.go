package timestamp

import (
	"bytes"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

func Now(client substrate.Client) (*time.Time, error) {
	m, err := client.MetaData(true)
	if err != nil {
		return nil, err
	}

	key, _ := substrate.NewStorageKey(*m,"Timestamp", "Now", nil)
	s := substrate.NewStateRPC(client)
	res, err := s.Storage(key,  nil)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(res)
	tempDec := scale.NewDecoder(buf)
	var ts uint64
	err = tempDec.Decode(&ts)
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(ts), 0)
	return &t, nil
}