package trie

import "bytes"

type Keyable interface {
	~[]byte | ~string
}

const (
	LessThan    = -1
	EqualTo     = 0
	GreaterThan = 1
)

// CompareKeys performs a byte-level comparison of two Keys
func CompareKeys[Key Keyable](l, r Key) int {
	lb := ([]byte)(l)
	rb := ([]byte)(r)
	return bytes.Compare(lb, rb)
}

// EqualKeys returns whether the provided Keys are equal
func EqualKeys[Key Keyable](l, r Key) bool {
	return CompareKeys[Key](l, r) == 0
}

// LessThanKeys returns whether the left Key is less than the right Key
func LessThanKeys[Key Keyable](l, r Key) bool {
	return CompareKeys[Key](l, r) == -1
}

// GreaterThanKeys returns whether the left Key is greater than the right Key
func GreaterThanKeys[Key Keyable](l, r Key) bool {
	return CompareKeys[Key](l, r) == 1
}
