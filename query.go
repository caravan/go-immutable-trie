package trie

import "github.com/caravan/go-immutable-trie/key"

type (
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
		parent *iterator[Key, Value]
		*trie[Key, Value]
		idx int
	}

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
) Query[Key, Value] {
	return (&iterator[Key, Value]{
		trie: t,
	}).decorate()
}

func (i *iterator[Key, Value]) mutate(
	mutate func(*iterator[Key, Value]),
) *iterator[Key, Value] {
	res := &(*i)
	mutate(res)
	return res
}

func (i *iterator[Key, Value]) Next() (
	Pair[Key, Value], Query[Key, Value], bool,
) {
	p := i.pair
	n, ok := i.fetchNext()
	return &p, decorate(n), ok
}

func (i *iterator[Key, Value]) fetchNext() (Iterator[Key, Value], bool) {
	for idx, bucket := range i.buckets[i.idx:] {
		if bucket == nil {
			continue
		}
		return &iterator[Key, Value]{
			parent: i.advanceIndex(idx + 1),
			trie:   bucket,
		}, true
	}
	if i.parent != nil {
		return i.parent.fetchNext()
	}
	return nil, false
}

func (i *iterator[Key, Value]) advanceIndex(idx int) *iterator[Key, Value] {
	return i.mutate(func(i *iterator[Key, Value]) {
		i.idx += idx
	})
}

func (i *iterator[Key, Value]) nextDescending() (
	Pair[Key, Value], Query[Key, Value], bool,
) {
	panic("not implemented")
}

func (i *iterator[Key, Value]) decorate() Query[Key, Value] {
	return decorate[Key, Value](i)
}

func (w *where[Key, Value]) Next() (
	Pair[Key, Value], Query[Key, Value], bool,
) {
	for p, c, ok := w.Query.Next(); ok; p, c, ok = c.Next() {
		if w.Filter(p.Key(), p.Value()) {
			return p, (&where[Key, Value]{c, w.Filter}).decorate(), true
		}
	}
	return nil, nil, false
}

func (w *where[Key, Value]) decorate() Query[Key, Value] {
	return decorate[Key, Value](w)
}

func (w *while[Key, Value]) Next() (
	Pair[Key, Value], Query[Key, Value], bool,
) {
	if p, c, ok := w.Query.Next(); ok && w.Filter(p.Key(), p.Value()) {
		return p, (&while[Key, Value]{c, w.Filter}).decorate(), true
	}
	return nil, nil, false
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
	for p, n, ok := d.Next(); ok; p, n, ok = n.Next() {
		f(p.Key(), p.Value())
	}
}

func (d *decorated[Key, Value]) Where(f Filter[Key, Value]) Query[Key, Value] {
	return (&where[Key, Value]{d, f}).decorate()
}

func (d *decorated[Key, Value]) While(f Filter[Key, Value]) Query[Key, Value] {
	return (&while[Key, Value]{d, f}).decorate()
}
