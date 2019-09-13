# go-substrate-rpc-client (GSRPC)

Substrate RPC client in Go. It provides APIs and Types around Polkadot and any Substrate-based chain RPC calls. This client is modelled after [polkadot-js/api](https://github.com/polkadot-js/api).

## Documentation & Usage Examples

Please refer to https://godoc.org/github.com/centrifuge/go-substrate-rpc-client

## Contributing

1. Install dependencies by running `make` followed by `make install`
1. Run tests `make test`
1. Lint `make lint` (you can use `make lint-fix` to automatically fix issues)

## Run with centrifuge-chain

1. Install subkey command from [https://github.com/centrifuge/substrate](https://github.com/centrifuge/substrate): Clone it, `cd subkey` and run `cargo install --force --path .`
2. Run a centrifuge-chain locally:`docker run -p 9944:9944 -p 30333:30333 centrifugeio/centrifuge-chain:20190814150805-ddd3818 centrifuge-chain --ws-external --dev`
3. Now adjust the hardcoded const parameters in `test/main.go` according to your env + chain state.
4. Run `go run test/main.go`

- You need to install SubKey command from `https://github.com/centrifuge/substrate`
   - Clone it and `CD` to subkey module and run `cargo install --force --path .`
   
- Now adjust the hardcoded const parameters in `test/main.go` according to your env + chain state.
- Run `test/main.go`

## POC FLow

![Alt text](extrinsic-execution.png?raw=true "Extrinsic Execution")
