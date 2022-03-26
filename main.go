package trie

// New returns a new empty Trie instance
func New[Key Keyable, Value any]() Trie[Key, Value] {
	return empty[Key, Value]{}
}

// From builds a Trie from a map with string keys
func From[Value any](in map[string]Value) Trie[string, Value] {
	res := New[string, Value]()
	for k, v := range in {
		res = res.Put(k, v)
	}
	return res
}
