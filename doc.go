// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Philip Stanislaus, Philip Stehlik, Vimukthi Wickramasinghe
//
// This file is part of Go Substrate RPC Client (GSRPC).
//
// GSRPC is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// GSRPC is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

/*
Package gsrpc (Go Substrate RPC Client) provides APIs and types around Polkadot and any Substrate-based chain RPC calls.
This client is modelled after [polkadot-js/api](https://github.com/polkadot-js/api).

Calling RPC methods

Simply instantiate the gsrpc with a URL of your choice, e. g.

	api, err := gsrpc.NewSubstrateAPI("wss://substrate-rpc.parity.io/")

and run any of the provided RPC methods from the api:

	hash, err := api.RPC.Chain.GetBlockHashLatest()

Further examples can be found below.

Types

The package [types](http://localhost:6060/pkg/github.com/centrifuge/go-substrate-rpc-client/types/) exports a number
of useful basic types including functions for encoding and decoding them.

To use your own custom types, you can simply create structs and arrays composing those basic types. Here are some
examples using composition of a mix of these basic and builtin Go types:

1. Vectors, lists, series, sets, arrays, slices: http://localhost:6060/pkg/github.com/centrifuge/go-substrate-rpc-client/types/#example_Vec_simple

2. Structs: http://localhost:6060/pkg/github.com/centrifuge/go-substrate-rpc-client/types/#example_Struct_simple

There are some caveats though that you should be aware of:

1. The order of the values in your structs is of relevance to the encoding. The scale codec Substrate/Polkadot
uses does not encode labels/keys.

2. Some types do not have corresponding types in Go. Working with them requires a custom struct with Encoding/Decoding
methods that implement the Encodeable/Decodeable interfaces. Examples for that are enums, tuples and vectors with any
types, you can find reference implementations of those here: types/enum_test.go , types/tuple_test.go and
types/vec_any_test.go

For more information about the types sub-package, see http://localhost:6060/pkg/github.com/centrifuge/go-substrate-rpc-client/types
*/
package gsrpc
