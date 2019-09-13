package state

import (
	"github.com/centrifuge/go-substrate-rpc-client/client"
)

type State struct {
	client *client.Client
}

func NewState(c *client.Client) *State {
	return &State{c}
}
