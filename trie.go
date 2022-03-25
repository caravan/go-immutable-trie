package trie

type (
	// Trie maps a set of Keys to another set of Values
	Trie[Key Keyable, Value any] interface {
		trie() // marker
		Put(Key, Value) Trie[Key, Value]
		Get(Key) (Value, bool)
		Remove(Key) (Value, Trie[Key, Value], bool)
		First() Pair[Key, Value]
		Rest() Trie[Key, Value]
		Split() (Pair[Key, Value], Trie[Key, Value], bool)
		Count() int
		IsEmpty() bool
	}

	Pair[Key Keyable, Value any] interface {
		pair() // marker
		Key() Key
		Value() Value
	}

	trie[Key Keyable, Value any] struct {
		pair    pair[Key, Value]
		buckets [nibbleSize]*trie[Key, Value]
	}

	pair[Key Keyable, Value any] struct {
		key   Key
		value Value
	}

	empty[Key Keyable, Value any] struct{}
)

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

func (*trie[_, _]) trie() {}

func (t *trie[Key, Value]) Get(k Key) (Value, bool) {
	h := Nibble(k)
	return t.get(k, h)
}

func (t *trie[Key, Value]) get(k Key, n Nibbles[Key]) (Value, bool) {
	if EqualKeys[Key](t.pair.key, k) {
		return t.pair.value, true
	}
	if idx, rest, ok := n.Consume(); ok {
		bucket := t.buckets[idx]
		if bucket != nil {
			return bucket.get(k, rest)
		}
	}
	var zero Value
	return zero, false
}

func (t *trie[Key, Value]) Put(k Key, v Value) Trie[Key, Value] {
	p := &pair[Key, Value]{k, v}
	h := Nibble[Key](p.key)
	return t.put(p, h)
}

func (t *trie[Key, Value]) put(
	p *pair[Key, Value], n Nibbles[Key],
) *trie[Key, Value] {
	switch CompareKeys[Key](p.key, t.pair.key) {
	case EqualTo:
		return t.replacePair(p)
	case LessThan:
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
	p *pair[Key, Value], n Nibbles[Key],
) *trie[Key, Value] {
	res := t.demoted(n)
	res.pair = *p
	return res
}

func (t *trie[Key, Value]) demoted(n Nibbles[Key]) *trie[Key, Value] {
	res := *t
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
	p *pair[Key, Value], n Nibbles[Key],
) *trie[Key, Value] {
	res := *t
	if idx, rest, ok := n.Consume(); ok {
		bucket := t.buckets[idx]
		if bucket == nil {
			res.buckets[idx] = &trie[Key, Value]{pair: *p}
		} else {
			res.buckets[idx] = bucket.put(p, rest)
		}
	} else {
		panic("programmer error: appended a non-consumable key")
	}
	return &res
}

func (t *trie[Key, Value]) Remove(k Key) (Value, Trie[Key, Value], bool) {
	h := Nibble(k)
	if v, r, ok := t.remove(k, h); ok {
		if r != nil {
			return v, r, true
		}
		return v, empty[Key, Value]{}, true
	}
	var zero Value
	return zero, t, false
}

func (t *trie[Key, Value]) remove(
	k Key, n Nibbles[Key],
) (Value, *trie[Key, Value], bool) {
	if EqualKeys[Key](t.pair.key, k) {
		return t.pair.value, t.promote(), true
	}
	if idx, rest, ok := n.Consume(); ok {
		if bucket := t.buckets[idx]; bucket != nil {
			if v, r, ok := bucket.remove(k, rest); ok {
				res := *t
				res.buckets[idx] = r
				return v, &res, true
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
		if k := bucket.pair.Key(); first || LessThanKeys[Key](k, low) {
			idx = i
			low = k
			res = bucket
			first = false
		}
	}
	return res, idx
}

func (t *trie[Key, Value]) First() Pair[Key, Value] {
	f, _, _ := t.Split()
	return f
}

func (t *trie[Key, Value]) Rest() Trie[Key, Value] {
	_, r, _ := t.Split()
	return r
}

func (t *trie[Key, Value]) Split() (Pair[Key, Value], Trie[Key, Value], bool) {
	if r := t.promote(); r != nil {
		return t.pair, r, true
	}
	return t.pair, empty[Key, Value]{}, true
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

func (pair[_, _]) pair() {}

func (p pair[Key, _]) Key() Key {
	return p.key
}

func (p pair[_, Value]) Value() Value {
	return p.value
}

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
