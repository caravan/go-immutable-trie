package trie

func New[Key Keyable, Value any]() Trie[Key, Value] {
	return empty[Key, Value]{}
}

func From[Value any](in map[string]Value) Trie[string, Value] {
	res := New[string, Value]()
	for k, v := range in {
		res = res.Put(k, v)
	}
	return res
}
