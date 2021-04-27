package mmr

import "github.com/snowfork/go-substrate-rpc-client/v2/client"

// MMR exposes methods for retrieval of MMR data
type MMR struct {
	client client.Client
}

// NewMMR creates a new MMR struct
func NewMMR(c client.Client) *MMR {
	return &MMR{client: c}
}
