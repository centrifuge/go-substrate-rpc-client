package types

// Bytes represents byte slices
type Bytes []byte

// NewBytes creates a new Bytes type
func NewBytes(b []byte) Bytes {
	return Bytes(b)
}
