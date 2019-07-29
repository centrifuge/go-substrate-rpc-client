package substrate

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
)

type Client interface {
	Call(result interface{}, method string, args ...interface{}) error

	Subscribe(ctx context.Context, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error)

	MetaData(cache bool) (*MetadataVersioned, error)
}

type client struct {
	*rpc.Client

	// metadataVersioned is the metadata cache to prevent unnecessary requests
	metadataVersioned *MetadataVersioned
}

func (c *client) MetaData(cache bool) (m *MetadataVersioned, err error) {
	if cache && c.metadataVersioned != nil {
		m = c.metadataVersioned
	} else {
		s := NewStateRPC(c)
		m, err = s.MetaData(nil)
		if err != nil {
			return nil, err
		}
		// set cache
		c.metadataVersioned = m
	}
	return
}

// Connect
func Connect(url string) (Client, error) {
	c, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	cc := client{c, nil}
	return &cc, nil
}
