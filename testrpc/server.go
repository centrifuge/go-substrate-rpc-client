package testrpc

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

type authorService struct {
}

func (s *authorService) SubmitExtrinsic(hex string) string {
	// TODO make this more dynamic and have proper tests when subkey is included in the build
	return hex
}

type stateService struct {
	metadata string

	storage map[string]string

	storageForBlock map[string]map[string]string
}

func newStateService(metadata string) *stateService {
	return &stateService{metadata: metadata, storageForBlock: make(map[string]map[string]string), storage: make(map[string]string)}
}

func (s *stateService) GetMetadata(blocknum *string) string {
	return s.metadata
}

func (s *stateService) GetStorage(key *string, blocknum *string) string {
	if key != nil && blocknum != nil {
		return s.storageForBlock[*key][*blocknum]
	} else if key != nil {
		return s.storage[*key]
	}
	return ""
}

type Server struct {
	author *authorService
	state  *stateService

	server *rpc.Server
}

// Following 4 methods are not go routine safe

func (s *Server) AddStorageKey(key, value string) {
	s.state.storage[key] = value
}

func (s *Server) AddStorageKeyForBlock(key, blocknum, value string) {
	s.state.storageForBlock[key][blocknum] = value
}

func (s *Server) RemoveStorageKey(key string) {
	delete(s.state.storage, key)
}

func (s *Server) RemoveStorageKeyForBlock(key, blocknum string) {
	delete(s.state.storageForBlock[key], blocknum)
}

// Init inits the testrpc server. port is the port that should be used.
func (ts *Server) Init(metadata string) (int, error) {
	ts.author = new(authorService)
	ts.state = newStateService(metadata)
	server := rpc.NewServer()
	err := server.RegisterName("author", ts.author)
	if err != nil {
		return 0, err
	}

	err = server.RegisterName("state", ts.state)
	if err != nil {
		return 0, err
	}

	http.Handle("/", server.WebsocketHandler([]string{"*"}))
	port := randomPort()
	go http.ListenAndServe(":"+strconv.Itoa(port), nil)
	// allow sometime for the server to start
	time.Sleep(10 * time.Millisecond)

	return port, nil
}

func randomPort() int {
	rand.Seed(time.Now().UnixNano())
	min := 10000
	max := 30000
	return rand.Intn(max-min+1) + min
}
