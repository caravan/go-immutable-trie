package key

import "bytes"

type Keyable interface {
	~[]byte | ~string
}

const (
	Less    = -1
	Equal   = 0
	Greater = 1
)

// Compare performs a byte-level comparison of two Keys
func Compare[Key Keyable](l, r Key) int {
	lb := ([]byte)(l)
	rb := ([]byte)(r)
	return bytes.Compare(lb, rb)
}

// EqualTo returns whether the provided Keys are equal
func EqualTo[Key Keyable](l, r Key) bool {
	return Compare[Key](l, r) == 0
}

// LessThan returns whether the left Key is less than the right Key
func LessThan[Key Keyable](l, r Key) bool {
	return Compare[Key](l, r) == -1
}

// GreaterThan returns whether the left Key is greater than the right Key
func GreaterThan[Key Keyable](l, r Key) bool {
	return Compare[Key](l, r) == 1
}
