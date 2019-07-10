# go-substrate

The substrate RPC client for go. Some slides https://docs.google.com/presentation/d/1lvCP7wZpsl2ES6fAkHg8-LfLhgtMHcvHeITp77Bgn2w/edit#slide=id.g5d531cd2a8_0_70

## How to run 

- You need to install SubKey command from `https://github.com/centrifuge/substrate`
   - Clone it and `CD` to subkey module and run `cargo install --force --path .`
   
- Now adjust the hardcoded const parameters in `test/main.go` according to your env + chain state.
- Run `test/main.go`

## POC FLow

![Alt text](extrinsic-execution.png?raw=true "Extrinsic Execution")
