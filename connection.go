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

	// m is the metadata cache to prevent unnecessary requests
	m *MetadataVersioned
}

func (c *client) MetaData(cache bool) (m *MetadataVersioned, err error) {
	if cache && c.m != nil {
		m = c.m
	} else {
		s := NewStateRPC(c)
		m, err = s.MetaData(nil)
		if err != nil {
			return nil, err
		}
		// set cache
		c.m = m
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
