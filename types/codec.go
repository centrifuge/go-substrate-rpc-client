package types

// Codec is the base interface that all types implement. The Codec Base is required for operating as an encoding/decoding layer.
type Codec interface {
	// The length of the value when encoded as a byte array
	EncodedLength() (int, error)
	// Returns a hash of the value
	Hash() ([32]byte, error)
	// Checks if the value is an empty value
	IsEmpty() bool
	// Compares the value of the input to see if there is a match
	Eq(o Codec) bool
	// Returns a hex string representation of the value. isLe returns a LE (number-only) representation
	Hex() (string, error)
	// Returns the string representation of the value
	String() string
	// Encodes the value as a byte array as per the SCALE specifications
	Encode() ([]byte, error)
}
