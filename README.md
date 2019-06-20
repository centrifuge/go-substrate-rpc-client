# go-substrate

The substrate RPC client for go.

## How to run 

- You need to install SubKey command from `https://github.com/centrifuge/substrate`
   - Clone it and `CD` to subkey module and run `cargo install --force --path .`
   
- Now adjust the hardcoded const parameters in `test/main.go` according to your env + chain state.
- Run `test/main.go`