package trie

import (
	"github.com/caravan/go-immutable-trie/key"
	"github.com/caravan/go-immutable-trie/nibble"
)

type (
	// Trie maps a set of Keys to another set of Values
	Trie[Key key.Keyable, Value any] interface {
		trie() // marker
		Put(Key, Value) Trie[Key, Value]
		Get(Key) (Value, bool)
		Remove(Key) (Value, Trie[Key, Value], bool)
		First() Pair[Key, Value]
		Rest() Trie[Key, Value]
		Split() (Pair[Key, Value], Trie[Key, Value], bool)
		Count() int
		IsEmpty() bool
		Select() Direction[Key, Value]
	}

	trie[Key key.Keyable, Value any] struct {
		pair    pair[Key, Value]
		buckets [nibble.Size]*trie[Key, Value]
	}
)

func (*trie[_, _]) trie() {}

func (t *trie[Key, Value]) Get(k Key) (Value, bool) {
	n := nibble.Make(k)
	return t.get(k, n)
}

func (t *trie[Key, Value]) get(k Key, n nibble.Nibbles[Key]) (Value, bool) {
	if key.EqualTo[Key](t.pair.key, k) {
		return t.pair.value, true
	}
	if idx, n, ok := n.Consume(); ok {
		bucket := t.buckets[idx]
		if bucket != nil {
			return bucket.get(k, n)
		}
	}
	var zero Value
	return zero, false
}

func (t *trie[Key, Value]) Put(k Key, v Value) Trie[Key, Value] {
	p := &pair[Key, Value]{k, v}
	n := nibble.Make[Key](p.key)
	return t.put(p, n)
}

func (t *trie[Key, Value]) put(
	p *pair[Key, Value], n nibble.Nibbles[Key],
) *trie[Key, Value] {
	switch key.Compare[Key](p.key, t.pair.key) {
	case key.Equal:
		return t.replacePair(p)
	case key.Less:
		return t.insertPair(p, n)
	default:
		return t.appendPair(p, n)
	}
}

func (t *trie[Key, Value]) replacePair(p *pair[Key, Value]) *trie[Key, Value] {
	res := *t
	res.pair = *p
	return &res
}

func (t *trie[Key, Value]) insertPair(
	p *pair[Key, Value], n nibble.Nibbles[Key],
) *trie[Key, Value] {
	res := *t
	res.pair = *p
	if idx, next, ok := n.Branch(t.pair.key).Consume(); ok {
		bucket := res.buckets[idx]
		if bucket == nil {
			res.buckets[idx] = &trie[Key, Value]{pair: t.pair}
		} else {
			res.buckets[idx] = bucket.put(&t.pair, next)
		}
	} else {
		panic("programmer error: demoted a non-consumable key")
	}
	return &res
}

func (t *trie[Key, Value]) appendPair(
	p *pair[Key, Value], n nibble.Nibbles[Key],
) *trie[Key, Value] {
	res := *t
	if idx, n, ok := n.Consume(); ok {
		bucket := t.buckets[idx]
		if bucket == nil {
			res.buckets[idx] = &trie[Key, Value]{pair: *p}
		} else {
			res.buckets[idx] = bucket.put(p, n)
		}
	} else {
		panic("programmer error: appended a non-consumable key")
	}
	return &res
}

func (t *trie[Key, Value]) Remove(k Key) (Value, Trie[Key, Value], bool) {
	n := nibble.Make(k)
	if val, rest, ok := t.remove(k, n); ok {
		if rest != nil {
			return val, rest, true
		}
		return val, empty[Key, Value]{}, true
	}
	var zero Value
	return zero, t, false
}

func (t *trie[Key, Value]) remove(
	k Key, n nibble.Nibbles[Key],
) (Value, *trie[Key, Value], bool) {
	if key.EqualTo[Key](t.pair.key, k) {
		return t.pair.value, t.promote(), true
	}
	if idx, n, ok := n.Consume(); ok {
		if bucket := t.buckets[idx]; bucket != nil {
			if val, rest, ok := bucket.remove(k, n); ok {
				res := *t
				res.buckets[idx] = rest
				return val, &res, true
			}
		}
	}
	var zero Value
	return zero, nil, false
}

func (t *trie[Key, Value]) promote() *trie[Key, Value] {
	if bucket, idx := t.leastBucket(); bucket != nil {
		res := *t
		res.pair = bucket.pair
		res.buckets[idx] = bucket.promote()
		return &res
	}
	return nil
}

func (t *trie[Key, Value]) leastBucket() (*trie[Key, Value], int) {
	var res *trie[Key, Value]
	var low Key
	idx := -1
	first := true
	for i, bucket := range t.buckets {
		if bucket == nil {
			continue
		}
		if k := bucket.pair.Key(); first || key.LessThan[Key](k, low) {
			idx = i
			low = k
			res = bucket
			first = false
		}
	}
	return res, idx
}

func (t *trie[Key, Value]) First() Pair[Key, Value] {
	f := t.pair
	return &f
}

func (t *trie[Key, Value]) Rest() Trie[Key, Value] {
	if rest := t.promote(); rest != nil {
		return rest
	}
	return empty[Key, Value]{}
}

func (t *trie[Key, Value]) Split() (Pair[Key, Value], Trie[Key, Value], bool) {
	first := t.pair
	if r := t.promote(); r != nil {
		return &first, r, true
	}
	return &first, empty[Key, Value]{}, true
}

func (t *trie[_, _]) Count() int {
	res := 1
	for _, bucket := range t.buckets {
		if bucket != nil {
			res += bucket.Count()
		}
	}
	return res
}

func (t *trie[_, _]) IsEmpty() bool {
	return false
}

func (t *trie[Key, Value]) Select() Direction[Key, Value] {
	return makeQuery[Key, Value](t)
}
