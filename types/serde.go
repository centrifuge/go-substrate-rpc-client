package types

import "sync"

// SerDeOptions are serialise and deserialize options for types
type SerDeOptions struct {
	// NoPalletIndices enable this to work with substrate chains that do not have indices pallet in runtime
	NoPalletIndices bool
}

var defaultOptions = SerDeOptions{}
var mu sync.RWMutex

// SetSerDeOptions overrides default serialise and deserialize options
func SetSerDeOptions(so SerDeOptions) {
	defer mu.Unlock()
	mu.Lock()
	defaultOptions = so
}

// BuildOptsFromMetadata sets certain EncoderOptions based on chain metadata
func BuildOptsFromMetadata(meta *Metadata) SerDeOptions {
	var opts SerDeOptions
	if !meta.ExistsModuleMetadata("Indices") {
		opts.NoPalletIndices = true
	}
	return opts
}
