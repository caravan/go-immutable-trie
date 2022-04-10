package trie

import (
	"github.com/caravan/go-immutable-trie/key"
	"github.com/caravan/go-immutable-trie/nibble"
)

type (
	// Trie maps a set of Keys to another set of Values
	Trie[Key key.Keyable, Value any] interface {
		trie() // marker
		Read[Key, Value]
		Split[Key, Value]
		Write[Key, Value]
	}

	Read[Key key.Keyable, Value any] interface {
		Get(Key) (Value, bool)
		Count() int
		IsEmpty() bool
		Select() Direction[Key, Value]
	}

	Split[Key key.Keyable, Value any] interface {
		First() Pair[Key, Value]
		Rest() Trie[Key, Value]
		Split() (Pair[Key, Value], Trie[Key, Value], bool)
	}

	Write[Key key.Keyable, Value any] interface {
		Put(Key, Value) Trie[Key, Value]
		Remove(Key) (Value, Trie[Key, Value], bool)
		RemovePrefix(Key) (Trie[Key, Value], bool)
	}

	trie[Key key.Keyable, Value any] struct {
		pair[Key, Value]
		*buckets[Key, Value]
	}

	buckets[Key key.Keyable, Value any] [nibble.Size]*trie[Key, Value]

	bucketsMutator[Key key.Keyable, Value any] func(*buckets[Key, Value])
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
	if idx, n, ok := n.Consume(); ok && t.buckets != nil {
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
	if idx, next, ok := n.Branch(t.pair.key).Consume(); ok {
		res := t.mutateBuckets(func(buckets *buckets[Key, Value]) {
			bucket := buckets[idx]
			if bucket == nil {
				buckets[idx] = &trie[Key, Value]{pair: t.pair}
			} else {
				buckets[idx] = bucket.put(&t.pair, next)
			}
		})
		res.pair = *p
		return res
	}
	panic("programmer error: demoted a non-consumable key")

}

func (t *trie[Key, Value]) appendPair(
	p *pair[Key, Value], n nibble.Nibbles[Key],
) *trie[Key, Value] {
	if idx, n, ok := n.Consume(); ok {
		return t.mutateBuckets(func(buckets *buckets[Key, Value]) {
			bucket := buckets[idx]
			if bucket == nil {
				buckets[idx] = &trie[Key, Value]{pair: *p}
			} else {
				buckets[idx] = bucket.put(p, n)
			}
		})
	}
	panic("programmer error: appended a non-consumable key")
}

func (t *trie[Key, Value]) RemovePrefix(k Key) (Trie[Key, Value], bool) {
	n := nibble.Make(k)
	if res, ok := t.removePrefix(k, n); ok {
		return res, ok
	}
	return t, false
}

func (t *trie[Key, Value]) removePrefix(
	k Key, n nibble.Nibbles[Key],
) (*trie[Key, Value], bool) {
	if key.StartsWith(t.key, k) {
		if res := t.promote(); res != nil {
			res, _ = res.removePrefix(k, n)
			return res, true
		}
		return nil, true
	}
	if idx, n, ok := n.Consume(); ok && t.buckets != nil {
		if bucket := t.buckets[idx]; bucket != nil {
			if bucket, ok := bucket.removePrefix(k, n); ok {
				return t.mutateBuckets(func(buckets *buckets[Key, Value]) {
					buckets[idx] = bucket
				}), true
			}
		}
	}
	return t, false
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
	if idx, n, ok := n.Consume(); ok && t.buckets != nil {
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

func (t *trie[Key, Value]) mutateBuckets(
	mutate bucketsMutator[Key, Value],
) *trie[Key, Value] {
	res := *t
	var storage buckets[Key, Value]
	if res.buckets != nil {
		storage = *res.buckets
	}
	res.buckets = &storage
	mutate(res.buckets)
	return &res
}

func (t *trie[Key, Value]) promote() *trie[Key, Value] {
	if bucket, idx := t.leastBucket(); bucket != nil {
		res := t.mutateBuckets(func(buckets *buckets[Key, Value]) {
			buckets[idx] = bucket.promote()
		})
		res.pair = bucket.pair
		return res
	}
	return nil
}

func (t *trie[Key, Value]) leastBucket() (*trie[Key, Value], int) {
	var res *trie[Key, Value]
	var low Key
	idx := -1
	first := true
	if t.buckets != nil {
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
	if t.buckets != nil {
		for _, bucket := range t.buckets {
			if bucket != nil {
				res += bucket.Count()
			}
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
