package key

import "bytes"

type (
	Keyable interface {
		~[]byte | ~string
	}

	Comparison int
)

const (
	Less Comparison = iota - 1
	Equal
	Greater
)

// Compare performs a byte-level comparison of two Keys
func Compare[Key Keyable](l, r Key) Comparison {
	lb := ([]byte)(l)
	rb := ([]byte)(r)
	return Comparison(bytes.Compare(lb, rb))
}

// EqualTo returns whether the provided Keys are equal
func EqualTo[Key Keyable](l, r Key) bool {
	return Compare[Key](l, r) == Equal
}

// LessThan returns whether the left Key is less than the right Key
func LessThan[Key Keyable](l, r Key) bool {
	return Compare[Key](l, r) == Less
}

// GreaterThan returns whether the left Key is greater than the right Key
func GreaterThan[Key Keyable](l, r Key) bool {
	return Compare[Key](l, r) == Greater
}
