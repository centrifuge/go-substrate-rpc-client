// Copyright 2018 Jsgenesis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpcmocksrv

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

// type authorService struct {
// }

// func (s *authorService) SubmitExtrinsic(hex string) string {
// 	// TODO make this more dynamic and have proper tests when subkey is included in the build
// 	return hex
// }

// type stateService struct {
// 	metadata string

// 	storage map[string]string

// 	storageForBlock map[string]map[string]string
// }

// func newStateService(metadata string) *stateService {
// 	return &stateService{metadata: metadata, storageForBlock: make(map[string]map[string]string),
//	  storage: make(map[string]string)}
// }

// func (s *stateService) GetMetadata(blocknum *string) string {
// 	return s.metadata
// }

// func (s *stateService) GetStorage(key *string, blocknum *string) string {
// 	if key != nil && blocknum != nil {
// 		return s.storageForBlock[*key][*blocknum]
// 	} else if key != nil {
// 		return s.storage[*key]
// 	}
// 	return ""
// }

type Server struct {
	// author *authorService
	// state  *stateService

	*rpc.Server
	// Host consists of hostname and port
	Host string
	// URL consists of protocol, hostname and port
	URL string
}

// Following 4 methods are not go routine safe

// func (s *Server) AddStorageKey(key, value string) {
// 	s.state.storage[key] = value
// }

// func (s *Server) AddStorageKeyForBlock(key, blocknum, value string) {
// 	s.state.storageForBlock[key][blocknum] = value
// }

// func (s *Server) RemoveStorageKey(key string) {
// 	delete(s.state.storage, key)
// }

// func (s *Server) RemoveStorageKeyForBlock(key, blocknum string) {
// 	delete(s.state.storageForBlock[key], blocknum)
// }

// New creates a new RPC mock server with a random port that allows registration of services
func New() *Server {
	port := randomPort()
	host := "localhost:" + strconv.Itoa(port)

	_, rpcServ, err := rpc.StartWSEndpoint(host, []rpc.API{}, []string{}, []string{"*"}, true)
	if err != nil {
		panic(err)
	}
	s := Server{
		Server: rpcServ,
		Host:   host,
		URL:    "ws://" + host,
	}
	return &s
}

// Init inits the testrpc server. rpcURL is the rpc url, eg: localhost:8080
// func (ts *Server) Init(metadata string, rpcURL *string) (string, error) {
// ts.author = new(authorService)
// ts.state = newStateService(metadata)
// server := rpc.NewServer()
// err := server.RegisterName("author", ts.author)
// if err != nil {
// 	return "", err
// }

// err = server.RegisterName("state", ts.state)
// if err != nil {
// 	return "", err
// }

// http.Handle("/", server.WebsocketHandler([]string{"*"}))
// port := randomPort()
// url := ""
// if rpcURL == nil {
// 	url = "localhost:" + strconv.Itoa(port)
// } else {
// 	url = *rpcURL
// }
// go http.ListenAndServe(url, nil)

// allow sometime for the server to start
// time.Sleep(10 * time.Millisecond)

// return "ws://" + url, nil
// }

func randomPort() int {
	rand.Seed(time.Now().UnixNano())
	min := 10000
	max := 30000
	return rand.Intn(max-min+1) + min
}
