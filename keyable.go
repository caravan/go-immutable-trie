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

func CompareKeys[Key Keyable](l, r Key) int {
	lb := ([]byte)(l)
	rb := ([]byte)(r)
	return bytes.Compare(lb, rb)
}

func EqualKeys[Key Keyable](l, r Key) bool {
	return CompareKeys[Key](l, r) == 0
}

func LessThanKeys[Key Keyable](l, r Key) bool {
	return CompareKeys[Key](l, r) == -1
}
