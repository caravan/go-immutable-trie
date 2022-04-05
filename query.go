package trie

import (
	"github.com/caravan/go-immutable-trie/key"
	"github.com/caravan/go-immutable-trie/nibble"
)

type (
	Direction[Key key.Keyable, Value any] interface {
		Select[Key, Value]
		Ascending() Select[Key, Value]
		Descending() Select[Key, Value]
	}

	Select[Key key.Keyable, Value any] interface {
		All() Query[Key, Value]
		From(Key) Query[Key, Value]
	}

	Iterator[Key key.Keyable, Value any] interface {
		Next() (Pair[Key, Value], Query[Key, Value], bool)
	}

	Query[Key key.Keyable, Value any] interface {
		Iterator[Key, Value]
		ForEach(ForEach[Key, Value])
		Where(Filter[Key, Value]) Query[Key, Value]
		While(Filter[Key, Value]) Query[Key, Value]
	}

	ForEach[Key key.Keyable, Value any] func(Key, Value)
	Filter[Key key.Keyable, Value any]  func(Key, Value) bool

	iterator[Key key.Keyable, Value any] struct {
		parent     *iterator[Key, Value]
		descending bool
		*trie[Key, Value]
		idx int
	}

	fetcher[Key key.Keyable, Value any] func() (
		Pair[Key, Value], Iterator[Key, Value],
	)

	decorated[Key key.Keyable, Value any] struct {
		Iterator[Key, Value]
	}

	where[Key key.Keyable, Value any] struct {
		Query[Key, Value]
		Filter[Key, Value]
	}

	while[Key key.Keyable, Value any] struct {
		Query[Key, Value]
		Filter[Key, Value]
	}
)

func makeQuery[Key key.Keyable, Value any](
	t *trie[Key, Value],
) Direction[Key, Value] {
	return &iterator[Key, Value]{
		trie: t,
	}
}

func (i *iterator[Key, Value]) mutate(
	mutate func(*iterator[Key, Value]),
) *iterator[Key, Value] {
	res := &(*i)
	mutate(res)
	return res
}

func (i *iterator[Key, Value]) Ascending() Select[Key, Value] {
	return i
}

func (i *iterator[Key, Value]) Descending() Select[Key, Value] {
	return i.mutate(func(i *iterator[Key, Value]) {
		i.descending = true
	})
}

func (i *iterator[Key, Value]) All() Query[Key, Value] {
	if i.descending {
		return i.last().decorate()
	}
	return i.decorate()
}

func (i *iterator[Key, Value]) From(k Key) Query[Key, Value] {
	n := nibble.Make(k)
	pos, ok := i.seek(k, n)
	if !ok && i.descending {
		return decorate(pos.prevParent())
	}
	return pos.decorate()
}

func (i *iterator[Key, Value]) last() *iterator[Key, Value] {
	for idx := len(i.buckets) - 1; idx >= 0; idx-- {
		if bucket := i.buckets[idx]; bucket != nil {
			return i.setIndex(idx).child(bucket).last()
		}
	}
	return i.setIndex(-1)
}

func (i *iterator[Key, Value]) seek(
	k Key, n nibble.Nibbles[Key],
) (*iterator[Key, Value], bool) {
	t := i.trie
	if key.EqualTo[Key](t.pair.key, k) {
		return i, true
	}
	if idx, n, ok := n.Consume(); ok {
		bucket := t.buckets[idx]
		if bucket != nil {
			parent := i.setIndex(int(idx))
			child := parent.child(bucket)
			return child.seek(k, n)
		}
	}
	return i, false
}

func (i *iterator[Key, Value]) Next() (
	Pair[Key, Value], Query[Key, Value], bool,
) {
	fetch := i.getFetcher()
	p, n := fetch()
	return p, decorate(n), true
}

func (i *iterator[Key, Value]) getFetcher() fetcher[Key, Value] {
	if i.descending {
		return i.fetchPrev
	}
	return i.fetchNext
}

func (i *iterator[Key, Value]) fetchNext() (
	Pair[Key, Value], Iterator[Key, Value],
) {
	p := i.pair
	if res, ok := i.nextBucket(); ok {
		return &p, res
	}
	return &p, empty[Key, Value]{}
}

func (i *iterator[Key, Value]) nextBucket() (*iterator[Key, Value], bool) {
	for idx := i.idx; idx < len(i.buckets); idx++ {
		if bucket := i.buckets[idx]; bucket != nil {
			parent := i.setIndex(idx)
			return parent.child(bucket), true
		}
	}
	if parent := i.parent; parent != nil {
		return parent.advanceIndex().nextBucket()
	}
	return nil, false
}

func (i *iterator[Key, Value]) fetchPrev() (
	Pair[Key, Value], Iterator[Key, Value],
) {
	p := i.pair
	return &p, i.prevParent()
}

func (i *iterator[Key, Value]) prevParent() Iterator[Key, Value] {
	if parent := i.parent; parent != nil {
		return parent.prevBucket()
	}
	return empty[Key, Value]{}
}

func (i *iterator[Key, Value]) prevBucket() *iterator[Key, Value] {
	for idx := i.idx - 1; idx >= 0; idx-- {
		if bucket := i.buckets[idx]; bucket != nil {
			return i.setIndex(idx).child(bucket).last()
		}
	}
	return i.setIndex(-1)
}

func (i *iterator[Key, Value]) advanceIndex() *iterator[Key, Value] {
	return i.setIndex(i.idx + 1)
}

func (i *iterator[Key, Value]) setIndex(idx int) *iterator[Key, Value] {
	return i.mutate(func(i *iterator[Key, Value]) {
		i.idx = idx
	})
}

func (i *iterator[Key, Value]) child(
	t *trie[Key, Value],
) *iterator[Key, Value] {
	return &iterator[Key, Value]{
		parent:     i,
		descending: i.descending,
		trie:       t,
	}
}

func (i *iterator[Key, Value]) decorate() Query[Key, Value] {
	return decorate[Key, Value](i)
}

func (w *where[Key, Value]) Next() (
	Pair[Key, Value], Query[Key, Value], bool,
) {
	for p, c, ok := w.Query.Next(); ok; p, c, ok = c.Next() {
		if w.Filter(p.Key(), p.Value()) {
			q := (&where[Key, Value]{c, w.Filter}).decorate()
			return p, q, true
		}
	}
	return nil, decoratedEmpty[Key, Value](), false
}

func (w *where[Key, Value]) decorate() Query[Key, Value] {
	return decorate[Key, Value](w)
}

func (w *while[Key, Value]) Next() (
	Pair[Key, Value], Query[Key, Value], bool,
) {
	if p, c, ok := w.Query.Next(); ok && w.Filter(p.Key(), p.Value()) {
		q := (&while[Key, Value]{c, w.Filter}).decorate()
		return p, q, true
	}
	return nil, decoratedEmpty[Key, Value](), false
}

func (w *while[Key, Value]) decorate() Query[Key, Value] {
	return decorate[Key, Value](w)
}

func decorate[Key key.Keyable, Value any](
	i Iterator[Key, Value],
) Query[Key, Value] {
	return &decorated[Key, Value]{i}
}

func (d *decorated[Key, Value]) ForEach(f ForEach[Key, Value]) {
	for p, q, ok := d.Next(); ok; p, q, ok = q.Next() {
		f(p.Key(), p.Value())
	}
}

func (d *decorated[Key, Value]) Where(f Filter[Key, Value]) Query[Key, Value] {
	return (&where[Key, Value]{d, f}).decorate()
}

func (d *decorated[Key, Value]) While(f Filter[Key, Value]) Query[Key, Value] {
	return (&while[Key, Value]{d, f}).decorate()
}
