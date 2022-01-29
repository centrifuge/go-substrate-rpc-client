// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
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

package gsrpc_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func Example_simpleConnect() {
	// The following example shows how to instantiate a Substrate API and use it to connect to a node

	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	chain, err := api.RPC.System.Chain()
	if err != nil {
		panic(err)
	}
	nodeName, err := api.RPC.System.Name()
	if err != nil {
		panic(err)
	}
	nodeVersion, err := api.RPC.System.Version()
	if err != nil {
		panic(err)
	}

	fmt.Printf("You are connected to chain %v using %v v%v\n", chain, nodeName, nodeVersion)
	// Output: You are connected to chain Development using Substrate Node v3.0.0-dev-1b646b2-x86_64-linux-gnu
}

func Example_listenToNewBlocks() {
	// This example shows how to subscribe to new blocks.
	//
	// It displays the block number every time a new block is seen by the node you are connected to.
	//
	// NOTE: The example runs until 10 blocks are received or until you stop it with CTRL+C

	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	sub, err := api.RPC.Chain.SubscribeNewHeads()
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	count := 0

	for {
		head := <-sub.Chan()
		fmt.Printf("Chain is at block: #%v\n", head.Number)
		count++

		if count == 10 {
			sub.Unsubscribe()
			break
		}
	}
}

func Example_listenToBalanceChange() {
	// This example shows how to instantiate a Substrate API and use it to connect to a node and retrieve balance
	// updates
	//
	// NOTE: The example runs until you stop it with CTRL+C

	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	alice := signature.TestKeyringPairAlice.PublicKey
	key, err := types.CreateStorageKey(meta, "System", "Account", alice)
	if err != nil {
		panic(err)
	}

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

	previous := accountInfo.Data.Free
	fmt.Printf("%#x has a balance of %v\n", alice, previous)
	fmt.Printf("You may leave this example running and transfer any value to %#x\n", alice)

	// Here we subscribe to any balance changes
	sub, err := api.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	// outer for loop for subscription notifications
	for {
		// inner loop for the changes within one of those notifications
		for _, chng := range (<-sub.Chan()).Changes {
			if !chng.HasStorageData {
				continue
			}

			var acc types.AccountInfo
			if err = types.DecodeFromBytes(chng.StorageData, &acc); err != nil {
				panic(err)
			}

			// Calculate the delta
			current := acc.Data.Free
			var change = types.U128{Int: big.NewInt(0).Sub(current.Int, previous.Int)}

			// Only display positive value changes (Since we are pulling `previous` above already,
			// the initial balance change will also be zero)
			if change.Cmp(big.NewInt(0)) != 0 {
				fmt.Printf("New balance change of: %v %v %v\n", change, previous, current)
				previous = current
				return
			}
		}
	}
}

func Example_unsubscribeFromListeningToUpdates() {
	// This example shows how to subscribe to and later unsubscribe from listening to block updates.
	//
	// In this example we're calling the built-in unsubscribe() function after a timeOut of 20s to cleanup and
	// unsubscribe from listening to updates.

	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	sub, err := api.RPC.Chain.SubscribeNewHeads()
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	timeout := time.After(20 * time.Second)

	for {
		select {
		case head := <-sub.Chan():
			fmt.Printf("Chain is at block: #%v\n", head.Number)
		case <-timeout:
			sub.Unsubscribe()
			fmt.Println("Unsubscribed")
			return
		}
	}
}

func Example_makeASimpleTransfer() {
	// This sample shows how to create a transaction to make a transfer from one an account to another.

	// Instantiate the API
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to Bob
	bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}

	// 1 unit of transfer
	bal, ok := new(big.Int).SetString("100000000000000", 10)
	if !ok {
		panic(fmt.Errorf("failed to convert balance"))
	}

	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompact(bal))
	if err != nil {
		panic(err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", signature.TestKeyringPairAlice.PublicKey)
	if err != nil {
		panic(err)
	}

	var accountInfo types.AccountInfo
	ok, err = api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

	nonce := uint32(accountInfo.Nonce)
	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(100),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction using Alice's default account
	err = ext.Sign(signature.TestKeyringPairAlice, o)
	if err != nil {
		panic(err)
	}

	// Send the extrinsic
	_, err = api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Balance transferred from Alice to Bob: %v\n", bal.String())
	// Output: Balance transferred from Alice to Bob: 100000000000000
}

func Example_makeASimpleTransferEd25519AccountTwo() {
	// This sample shows how to create a transaction to make a transfer from one an account to another.

	time.Sleep(10 * time.Second)

	// Instantiate the API
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to Ed25519 Account Two
	accountTwo, err := types.NewMultiAddressFromHexAccountID("0x05aeedbab13e7c50ee8e46421104e213395921ea067aa5fbb6dcc3617dace6a7")
	if err != nil {
		panic(err)
	}

	// 1 unit of transfer
	bal, ok := new(big.Int).SetString("100000000000000", 10)
	if !ok {
		panic(fmt.Errorf("failed to convert balance"))
	}

	c, err := types.NewCall(meta, "Balances.transfer", accountTwo, types.NewUCompact(bal))
	if err != nil {
		panic(err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", signature.TestKeyringPairAlice.PublicKey, nil)
	if err != nil {
		panic(err)
	}

	var accountInfo types.AccountInfo
	ok, err = api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

	nonce := uint32(accountInfo.Nonce)
	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	ext.SetSignScheme("Sr25519")

	// Sign the transaction using Alice's default account
	err = ext.Sign(signature.TestKeyringPairAlice, o)
	if err != nil {
		panic(err)
	}

	// Send the extrinsic
	_, err = api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Balance transferred from Alice to AccountTwo: %v\n", bal.String())
	// Output: Balance transferred from Alice to AccountTwo: 100000000000000
}

// TODO: ECDSA Account
// func Example_makeASimpleTransferEcdsaAccountThree() {
// 	// This sample shows how to create a transaction to make a transfer from one an account to another.

// 	time.Sleep(8 * time.Second)

// 	// Instantiate the API
// 	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	meta, err := api.RPC.State.GetMetadataLatest()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Create a call, transferring 12345 units to AccountThree Ecdsa
// 	// 	subkey generate --scheme Ecdsa --network polkadot
// 	// 	Secret phrase `clerk april useful wish balcony narrow member planet later explain cake harvest` is account:
// 	//   Secret seed:       0x954e40cec3ffbfe7f12e5552b9807a21100d5dc0e0e8b70320d588402159090c
// 	//   Public key (hex):  0x02ef86282d5a0c79ecfabaae31612d1d07c69dc5fce5e02761931e7febdcd1f732
// 	//   Public key (SS58): 1HzNpEUKiBZCAVAB3FZv12jJRk4zNvzBh5UoSSq7mbAE8yVu
// 	//   Account ID:        0x51ecb60d234ce572341bdfb6720b453e67e7b7de155a8843c7078c52dd1cb311
// 	//   SS58 Address:      12rRCQDcDP7YzxjirKEV8hbJu77PtgX16rBoBTYMbDcyoFxw
// 	accountTwo, err := types.NewMultiAddressFromHexAccountID("0x51ecb60d234ce572341bdfb6720b453e67e7b7de155a8843c7078c52dd1cb311")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 1 unit of transfer
// 	bal, ok := new(big.Int).SetString("100000000000000", 10)
// 	if !ok {
// 		panic(fmt.Errorf("failed to convert balance"))
// 	}

// 	c, err := types.NewCall(meta, "Balances.transfer", accountTwo, types.NewUCompact(bal))
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Create the extrinsic
// 	ext := types.NewExtrinsic(c)

// 	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
// 	if err != nil {
// 		panic(err)
// 	}

// 	rv, err := api.RPC.State.GetRuntimeVersionLatest()
// 	if err != nil {
// 		panic(err)
// 	}

// 	key, err := types.CreateStorageKey(meta, "System", "Account", signature.TestKeyringPairAlice.PublicKey, nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var accountInfo types.AccountInfo
// 	ok, err = api.RPC.State.GetStorageLatest(key, &accountInfo)
// 	if err != nil || !ok {
// 		panic(err)
// 	}

// 	nonce := uint32(accountInfo.Nonce)
// 	o := types.SignatureOptions{
// 		BlockHash:          genesisHash,
// 		Era:                types.ExtrinsicEra{IsMortalEra: false},
// 		GenesisHash:        genesisHash,
// 		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
// 		SpecVersion:        rv.SpecVersion,
// 		Tip:                types.NewUCompactFromUInt(0),
// 		TransactionVersion: rv.TransactionVersion,
// 	}

// 	ext.SetSignScheme("Ecdsa")

// 	// Sign the transaction using Alice's default account
// 	err = ext.Sign(signature.TestKeyringPairAlice, o)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Send the extrinsic
// 	_, err = api.RPC.Author.SubmitExtrinsic(ext)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("Balance transferred from Alice to AccountTwo: %v\n", bal.String())
// 	// Output: Balance transferred from Alice to AccountTwo: 100000000000000
// }

func Example_makeASimpleTransferEd25519() {

	// This sample shows how to create a transaction to make a transfer from one an account to another.
	time.Sleep(10 * time.Second)
	// Instantiate the API
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to Bob Account
	bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}
	// fmt.Println("BXL Example_makeASimpleTransferEd25519 bob : ", bob)

	// 1 unit of transfer
	bal, ok := new(big.Int).SetString("50000000000000", 10)
	if !ok {
		panic(fmt.Errorf("failed to convert balance"))
	}
	// fmt.Println("BXL Example_makeASimpleTransferEd25519 bal : ", bal)

	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompact(bal))
	if err != nil {
		panic(err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)
	// fmt.Println("BXL Example_makeASimpleTransferEd25519 ext : ", ext)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	// Ed25519 Account Two
	// 	subkey generate --scheme ed25519 --network polkadot
	// Secret phrase `panda claim miracle furnace dynamic battle kick cigar grain book fox badge` is account:
	//   Secret seed:       0x922b17a2230ca612424c9dc3990eaffa42e26cf74d795966819b6d04c980e81c
	//   Public key (hex):  0x05aeedbab13e7c50ee8e46421104e213395921ea067aa5fbb6dcc3617dace6a7
	//   Public key (SS58): 18TCqXdPyzrXWBDbd36foHbSQ7YitnqdNa8qdd1zUP6nUNF
	//   Account ID:        0x05aeedbab13e7c50ee8e46421104e213395921ea067aa5fbb6dcc3617dace6a7
	//   SS58 Address:      18TCqXdPyzrXWBDbd36foHbSQ7YitnqdNa8qdd1zUP6nUNF
	//   SS58 Address: 		5CCA4WGZYCjP5yAhdyz6XeTSan7u2bEhYsqegLdfSPMac53h

	pk, err := hex.DecodeString("05aeedbab13e7c50ee8e46421104e213395921ea067aa5fbb6dcc3617dace6a7")
	if err != nil {
		panic(err)
	}
	// fmt.Println("BXL Example_makeASimpleTransferEd25519 pk : ", pk)

	var testKeypair = signature.KeyringPair{
		URI:       "0x922b17a2230ca612424c9dc3990eaffa42e26cf74d795966819b6d04c980e81c",
		Address:   "5CCA4WGZYCjP5yAhdyz6XeTSan7u2bEhYsqegLdfSPMac53h",
		PublicKey: pk,
	}

	// fmt.Println("BXL:  testKeypair.PublicKey : ", testKeypair.PublicKey)
	// fmt.Println("BXL:  testKeypair.URI : ", testKeypair.URI)
	// fmt.Println("BXL:  testKeypair.Address : ", testKeypair.Address)

	key, err := types.CreateStorageKey(meta, "System", "Account", testKeypair.PublicKey, nil)
	if err != nil {
		panic(err)
	}
	// fmt.Println("BXL Example_makeASimpleTransferEd25519 key : ", key)

	var accountInfo types.AccountInfo
	ok, err = api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		fmt.Println("BXL Example_makeASimpleTransferEd25519 GetStorageLatest err : ", err)
		panic(err)
	}
	// fmt.Println("BXL Example_makeASimpleTransferEd25519 accountInfo : ", accountInfo)

	nonce := uint32(accountInfo.Nonce)
	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	ext.SetSignScheme("Ed25519")
	// fmt.Println("BXL Example_makeASimpleTransferEd25519 SetSignScheme :  IsEd25519	")

	// Sign the transaction using Alice's default account
	err = ext.Sign(testKeypair, o)
	if err != nil {
		panic(err)
	}
	// fmt.Println("BXL Example_makeASimpleTransferEd25519 ext : ", ext)

	// Debug encoding
	// enc, err := types.EncodeToHexString(ext)
	// fmt.Println("BXL Example_makeASimpleTransferEd25519 enc : ", enc)
	// if err != nil {
	// fmt.Println("BXL Example_makeASimpleTransferEd25519 err : ", err)

	// return types.Hash{}, err
	// }

	// Debug encoding
	// enc, err := types.EncodeToHexString(ext)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("BXL Example_makeASimpleTransferEd25519 enc : ", enc)

	// Send the extrinsic
	_, err = api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		fmt.Println("BXL Example_makeASimpleTransferEd25519 err : ", err)
		panic(err)
	}

	fmt.Printf("Balance transferred from Account2 to Bob: %v\n", bal.String())
	// Output: Balance transferred from Account2 to Bob: 50000000000000
}

// func Example_makeASimpleTransferEcdsa() {
// 	// This sample shows how to create a transaction to make a transfer from one an account to another.
// 	time.Sleep(8 * time.Second)

// 	// Instantiate the API
// 	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	meta, err := api.RPC.State.GetMetadataLatest()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Create a call, transferring 12345 units to Bob
// 	bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
// 	if err != nil {
// 		panic(err)
// 	}
// 	// fmt.Println("BXL Example_makeASimpleTransferEcdsa bob : ", bob)

// 	// 1 unit of transfer
// 	bal, ok := new(big.Int).SetString("50000000000000", 10)
// 	if !ok {
// 		panic(fmt.Errorf("failed to convert balance"))
// 	}
// 	// fmt.Println("BXL Example_makeASimpleTransferEcdsa bal : ", bal)

// 	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompact(bal))
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Create the extrinsic
// 	ext := types.NewExtrinsic(c)
// 	// fmt.Println("BXL Example_makeASimpleTransferEcdsa ext : ", ext)

// 	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
// 	if err != nil {
// 		panic(err)
// 	}

// 	rv, err := api.RPC.State.GetRuntimeVersionLatest()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Ecdsa Account
// 	// 	subkey generate --scheme Ecdsa --network polkadot
// 	// 	Secret phrase `clerk april useful wish balcony narrow member planet later explain cake harvest` is account:
// 	//   Secret seed:       0x954e40cec3ffbfe7f12e5552b9807a21100d5dc0e0e8b70320d588402159090c
// 	//   Public key (hex):  0x02ef86282d5a0c79ecfabaae31612d1d07c69dc5fce5e02761931e7febdcd1f732
// 	//   Public key (SS58): 1HzNpEUKiBZCAVAB3FZv12jJRk4zNvzBh5UoSSq7mbAE8yVu
// 	//   Account ID:        0x51ecb60d234ce572341bdfb6720b453e67e7b7de155a8843c7078c52dd1cb311
// 	//   SS58 Address:      12rRCQDcDP7YzxjirKEV8hbJu77PtgX16rBoBTYMbDcyoFxw

// 	pk, err := hex.DecodeString("51ecb60d234ce572341bdfb6720b453e67e7b7de155a8843c7078c52dd1cb311")
// 	if err != nil {
// 		panic(err)
// 	}
// 	// fmt.Println("BXL Example_makeASimpleTransferEcdsa pk : ", pk)

// 	var testKeypair = signature.KeyringPair{
// 		URI:       "0x954e40cec3ffbfe7f12e5552b9807a21100d5dc0e0e8b70320d588402159090c",
// 		Address:   "12rRCQDcDP7YzxjirKEV8hbJu77PtgX16rBoBTYMbDcyoFxw",
// 		PublicKey: pk,
// 	}

// 	fmt.Println("BXL:  testKeypair.PublicKey : ", testKeypair.PublicKey)
// 	fmt.Println("BXL:  testKeypair.URI : ", testKeypair.URI)
// 	fmt.Println("BXL:  testKeypair.Address : ", testKeypair.Address)

// 	key, err := types.CreateStorageKey(meta, "System", "Account", testKeypair.PublicKey, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("BXL Example_makeASimpleTransferEcdsa key : ", key.Hex())

// 	var accountInfo types.AccountInfo
// 	ok, err = api.RPC.State.GetStorageLatest(key, &accountInfo)
// 	if err != nil || !ok {
// 		fmt.Println("BXL Example_makeASimpleTransferEcdsa GetStorageLatest err : ", err)
// 		panic(err)
// 	}

// 	fmt.Println("BXL Example_makeASimpleTransferEcdsa accountInfo: ", accountInfo)

// 	fmt.Println("BXL Example_makeASimpleTransferEcdsa accountInfo.Data : ", accountInfo.Data)
// 	fmt.Printf("BXL %#x has a balance of %v\n", key, accountInfo.Data.Free)

// 	nonce := uint32(accountInfo.Nonce)
// 	o := types.SignatureOptions{
// 		BlockHash:          genesisHash,
// 		Era:                types.ExtrinsicEra{IsMortalEra: false},
// 		GenesisHash:        genesisHash,
// 		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
// 		SpecVersion:        rv.SpecVersion,
// 		Tip:                types.NewUCompactFromUInt(0),
// 		TransactionVersion: rv.TransactionVersion,
// 	}

// 	ext.SetSignScheme("Ecdsa")
// 	// fmt.Println("BXL Example_makeASimpleTransferEcdsa SetSignScheme :  IsEd25519	")

// 	fmt.Println("BXL Example_makeASimpleTransferEcdsa ext : ", ext)

// 	// Sign the transaction using Alice's default account
// 	err = ext.Sign(testKeypair, o)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// fmt.Println("BXL Example_makeASimpleTransferEcdsa ext : ", ext)

// 	// Debug encoding
// 	enc, err := types.EncodeToHexString(ext)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("BXL Example_makeASimpleTransferEcdsa enc : ", enc)

// 	// Send the extrinsic
// 	// _, err = api.RPC.Author.SubmitExtrinsic(ext)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// fmt.Println("BXL Example_makeASimpleTransferEd25519 result : ", result)

// 	fmt.Printf("Balance transferred from Account2 to Bob: %v\n", bal.String())
// 	// Output: Balance transferred from Account2 to Bob: 50000000000000
// }

func Example_displaySystemEvents() {
	// Query the system events and extract information from them. This example runs until exited via Ctrl-C

	// Create our API with a default connection to the local node
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Subscribe to system events via storage
	key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		panic(err)
	}

	sub, err := api.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	// outer for loop for subscription notifications
	for {
		set := <-sub.Chan()
		// inner loop for the changes within one of those notifications
		for _, chng := range set.Changes {
			if !types.Eq(chng.StorageKey, key) || !chng.HasStorageData {
				// skip, we are only interested in events with content
				continue
			}

			// Decode the event records
			events := types.EventRecords{}
			err = types.EventRecordsRaw(chng.StorageData).DecodeEventRecords(meta, &events)
			if err != nil {
				panic(err)
			}

			// Show what we are busy with
			for _, e := range events.Balances_Endowed {
				fmt.Printf("\tBalances:Endowed:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x, %v\n", e.Who, e.Balance)
			}
			for _, e := range events.Balances_DustLost {
				fmt.Printf("\tBalances:DustLost:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x, %v\n", e.Who, e.Balance)
			}
			for _, e := range events.Balances_Transfer {
				fmt.Printf("\tBalances:Transfer:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v, %v, %v\n", e.From, e.To, e.Value)
			}
			for _, e := range events.Balances_BalanceSet {
				fmt.Printf("\tBalances:BalanceSet:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v, %v, %v\n", e.Who, e.Free, e.Reserved)
			}
			for _, e := range events.Balances_Deposit {
				fmt.Printf("\tBalances:Deposit:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v, %v\n", e.Who, e.Balance)
			}
			for _, e := range events.Grandpa_NewAuthorities {
				fmt.Printf("\tGrandpa:NewAuthorities:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.NewAuthorities)
			}
			for _, e := range events.Grandpa_Paused {
				fmt.Printf("\tGrandpa:Paused:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.Grandpa_Resumed {
				fmt.Printf("\tGrandpa:Resumed:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.ImOnline_HeartbeatReceived {
				fmt.Printf("\tImOnline:HeartbeatReceived:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x\n", e.AuthorityID)
			}
			for _, e := range events.ImOnline_AllGood {
				fmt.Printf("\tImOnline:AllGood:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.ImOnline_SomeOffline {
				fmt.Printf("\tImOnline:SomeOffline:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.IdentificationTuples)
			}
			for _, e := range events.Indices_IndexAssigned {
				fmt.Printf("\tIndices:IndexAssigned:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x%v\n", e.AccountID, e.AccountIndex)
			}
			for _, e := range events.Indices_IndexFreed {
				fmt.Printf("\tIndices:IndexFreed:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.AccountIndex)
			}
			for _, e := range events.Offences_Offence {
				fmt.Printf("\tOffences:Offence:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v%v\n", e.Kind, e.OpaqueTimeSlot)
			}
			for _, e := range events.Session_NewSession {
				fmt.Printf("\tSession:NewSession:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.SessionIndex)
			}
			for _, e := range events.Staking_Reward {
				fmt.Printf("\tStaking:Reward:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.Amount)
			}
			for _, e := range events.Staking_Slash {
				fmt.Printf("\tStaking:Slash:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x%v\n", e.AccountID, e.Balance)
			}
			for _, e := range events.Staking_OldSlashingReportDiscarded {
				fmt.Printf("\tStaking:OldSlashingReportDiscarded:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.SessionIndex)
			}
			for _, e := range events.System_ExtrinsicSuccess {
				fmt.Printf("\tSystem:ExtrinsicSuccess:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.System_ExtrinsicFailed {
				fmt.Printf("\tSystem:ExtrinsicFailed:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.DispatchError)
			}
			for _, e := range events.System_CodeUpdated {
				fmt.Printf("\tSystem:CodeUpdated:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.System_NewAccount {
				fmt.Printf("\tSystem:NewAccount:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x\n", e.Who)
			}
			for _, e := range events.System_KilledAccount {
				fmt.Printf("\tSystem:KilledAccount:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#X\n", e.Who)
			}
		}
	}
}

func Example_transactionWithEvents() {
	// Display the events that occur during a transfer by sending a value to bob

	// Instantiate the API
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to Bob
	bob, err := types.NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}

	amount := types.NewUCompactFromUInt(12345)

	c, err := types.NewCall(meta, "Balances.transfer", bob, amount)
	if err != nil {
		panic(err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	// Get the nonce for Alice
	key, err := types.CreateStorageKey(meta, "System", "Account", signature.TestKeyringPairAlice.PublicKey)
	if err != nil {
		panic(err)
	}

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

	nonce := uint32(accountInfo.Nonce)

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	fmt.Printf("Sending %v from %#x to %#x with nonce %v", amount, signature.TestKeyringPairAlice.PublicKey, bob.AsAccountID, nonce)

	// Sign the transaction using Alice's default account
	err = ext.Sign(signature.TestKeyringPairAlice, o)
	if err != nil {
		panic(err)
	}

	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		status := <-sub.Chan()
		fmt.Printf("Transaction status: %#v\n", status)

		if status.IsInBlock {
			fmt.Printf("Completed at block hash: %#x\n", status.AsInBlock)
			return
		}
	}
}
