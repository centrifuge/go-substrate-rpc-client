# GSRPC Registry
The GSRPC Registry can parse target metadata information into an in-memory registry of complex structures. 

By leveraging the on-chain metadata, GSRPC is more robust to changes on types, allowing clients to only keep updated the types that are relevant to their business operation.

This registry can be used afterwards to decode data read from live chains (events & extrinsics).

## Usage
Since docs get outdated fairly quick, here are links to tests that will always be up-to-date.
### Populate Call, Error & Events Registries
[Browse me](registry_test.go)

### Event retriever
[TestLive_EventRetriever_GetEvents](retriever/event_retriever_live_test.go)
### Extrinsic retriever
Since chain runtimes can be customized, modifying core types such as Accounts, Signature payloads or Payment payloads, the code supports a customizable way of passing those custom types to the extrinsic retriever.

On the other hand, since a great majority of chains do not need to change these types, the tool provides a default for the most common used ones.
#### Using Chain Defaults
[TestExtrinsicRetriever_NewDefault](retriever/extrinsic_retriever_test.go#L179)
#### Using Custom core types
[TestLive_ExtrinsicRetriever_GetExtrinsics](retriever/extrinsic_retriever_live_test.go)