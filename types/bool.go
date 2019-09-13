package types

// Bool represents boolean values
type Bool bool

// NewBool creates a new Bool
func NewBool(b bool) Bool {
	return Bool(b)
}
