package types

type option struct {
	hasValue bool
}

// IsNone returns true if the value is missing
func (o option) IsNone() bool {
	return !o.hasValue
}

// IsNone returns true if a value is present
func (o option) IsSome() bool {
	return o.hasValue
}
