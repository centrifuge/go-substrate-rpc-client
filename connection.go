package substrate

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
)

type Client interface {
	Call(result interface{}, method string, args ...interface{}) error

	Subscribe(ctx context.Context, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error)
}

// Connect
func Connect(url string) (Client, error) {
	client, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}

	return client, nil
}
