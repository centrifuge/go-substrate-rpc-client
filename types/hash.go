package types

// H160 is a hash containing 160 bits (20 bytes), typically used in blocks, extrinsics and as a sane default
type H160 [20]byte

func NewH160(b [20]byte) H160 {
	return H160(b)
}

// H256 is a hash containing 256 bits (32 bytes), typically used in blocks, extrinsics and as a sane default
type H256 [32]byte

func NewH256(b [32]byte) H256 {
	return H256(b)
}

// H512 is a hash containing 512 bits (64 bytes), typically used for signature
type H512 [64]byte

func NewH512(b [64]byte) H512 {
	return H512(b)
}

// Hash is the default hash that is used across the system. It is just a thin wrapper around H256
type Hash H256

func NewHash(b [32]byte) Hash {
	return Hash(b)
}
