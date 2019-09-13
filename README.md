# go-substrate-rpc-client (GSRPC)

Substrate RPC client in Go. It provides APIs and Types around Polkadot and any Substrate-based chain RPC calls. 
This client is modelled after [polkadot-js/api](https://github.com/polkadot-js/api).

## State

This package is actively developed. The first stable release is expected in November 2019.

## Documentation & Usage Examples

Please refer to https://godoc.org/github.com/centrifuge/go-substrate-rpc-client

## Contributing

1. Install dependencies by running `make` followed by `make install`
1. Run tests `make test`
1. Lint `make lint` (you can use `make lint-fix` to automatically fix issues)

## Run tests in a Docker container against the Substrate Default Docker image

1. Run the docker container `make test-dockerized` 

## Run tests locally against the Substrate Default Docker image

1. Start the Substrate Default Docker image: `make run-substrate-docker`
1. In another terminal, run the tests against that image: `make test`
1. Visit https://polkadot.js.org/apps for inspection

## Run tests locally against any substrate endpoint

1. Set the endpoint: `export RPC_URL="http://example.com:9933"`
1. Run the tests `make test`
