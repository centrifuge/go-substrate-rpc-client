package timestamp

import (
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client"
)

func TestNow(t *testing.T) {
	c, _ := substrate.Connect("ws://127.0.0.1:9944")
	ts, _ := Now(c)
	fmt.Println(ts)
}
