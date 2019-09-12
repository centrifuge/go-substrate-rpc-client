package gsrpc

import "github.com/centrifuge/go-substrate-rpc-client/client"

type System struct {
	client client.Client
}

func NewSystemRPC(client client.Client) *System {
	return &System{client: client}
}
