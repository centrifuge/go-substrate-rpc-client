package substrate

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/minio/blake2b-simd"

	bbb "golang.org/x/crypto/blake2b"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestState_GetMetaData(t *testing.T) {
	s := State{nonetwork: true}
	res, err := s.MetaData([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, "system", res.Metadata.Modules[0].Name)
	// fmt.Println(res)
}

func TestBlake(t *testing.T) {
	bb, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000000901")
	b := blake2b.Sum256(bb)
 	b2 := bbb.Sum256(bb)
 	fmt.Println(hexutil.Encode(b[:]))
	fmt.Println(hexutil.Encode(b2[:]))
}
