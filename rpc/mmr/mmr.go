package mmr

import "github.com/centrifuge/go-substrate-rpc-client/v4/client"

// MMR exposes methods for retrieval of MMR data
type MMR struct {
	client client.Client
}

// NewMMR creates a new MMR struct
func NewMMR(c client.Client) *MMR {
	return &MMR{client: c}
}
