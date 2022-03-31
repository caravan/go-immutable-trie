package trie

import "github.com/caravan/go-immutable-trie/key"

type empty[Key key.Keyable, Value any] struct{}

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

func (e empty[Key, Value]) Select() Query[Key, Value] {
	return e
}

func (e empty[Key, Value]) From(Key) Query[Key, Value] {
	return e
}

func (e empty[Key, Value]) Ascending() Query[Key, Value] {
	return e
}

func (e empty[Key, Value]) Descending() Query[Key, Value] {
	return e
}

func (e empty[Key, Value]) Next() (Pair[Key, Value], Query[Key, Value], bool) {
	return nil, nil, false
}

func (e empty[Key, Value]) ForEach(ForEach[Key, Value]) {}

func (e empty[Key, Value]) Where(Filter[Key, Value]) Query[Key, Value] {
	return e
}

func (e empty[Key, Value]) While(Filter[Key, Value]) Query[Key, Value] {
	return e
}
