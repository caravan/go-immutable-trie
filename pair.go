package trie

import "github.com/caravan/go-immutable-trie/key"

type (
	// Pair stores Key/Value pairs
	Pair[Key key.Keyable, Value any] interface {
		pair() // marker
		Key() Key
		Value() Value
	}

	pair[Key key.Keyable, Value any] struct {
		key   Key
		value Value
	}
)

func (*pair[_, _]) pair() {}

func (p *pair[Key, _]) Key() Key {
	return p.key
}

func (p *pair[_, Value]) Value() Value {
	return p.value
}
