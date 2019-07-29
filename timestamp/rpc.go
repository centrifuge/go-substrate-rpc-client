package timestamp

import (
	"time"

	"github.com/centrifuge/go-substrate-rpc-client"
)

func Now(client substrate.Client) (*time.Time, error) {
	m, err := client.MetaData(true)
	if err != nil {
		return nil, err
	}

	key, err := substrate.NewStorageKey(*m, "Timestamp", "Now", nil)
	if err != nil {
		return nil, err
	}

	s := substrate.NewStateRPC(client)
	res, err := s.Storage(key, nil)
	if err != nil {
		return nil, err
	}

	tempDec := res.Decoder()
	var ts uint64
	err = tempDec.Decode(&ts)
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(ts), 0)
	return &t, nil
}
