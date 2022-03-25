package trie

type empty[Key Keyable, Value any] struct{}

func (empty[_, _]) trie() {}

func (empty[Key, Value]) Get(Key) (Value, bool) {
	var zero Value
	return zero, false
}

func (empty[Key, Value]) Put(k Key, v Value) Trie[Key, Value] {
	return &trie[Key, Value]{
		pair: pair[Key, Value]{k, v},
	}
}

func (empty[Key, Value]) Remove(_ Key) (Value, Trie[Key, Value], bool) {
	var zero Value
	return zero, empty[Key, Value]{}, false
}

func (empty[_, _]) IsEmpty() bool {
	return true
}

func (empty[_, _]) Count() int {
	return 0
}

func (empty[Key, Value]) Split() (Pair[Key, Value], Trie[Key, Value], bool) {
	return nil, empty[Key, Value]{}, false
}

func (e empty[Key, Value]) First() Pair[Key, Value] {
	f, _, _ := e.Split()
	return f
}

func (e empty[Key, Value]) Rest() Trie[Key, Value] {
	_, r, _ := e.Split()
	return r
}
