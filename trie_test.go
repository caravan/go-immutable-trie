package trie_test

import (
	"testing"

	trie "github.com/caravan/go-immutable-trie"
	"github.com/stretchr/testify/assert"
)

var testMap = map[string]int{
	"hello":   1,
	"there":   2,
	"how":     9,
	"are":     5,
	"you":     37,
	"today":   4,
	"curious": 128,
	"to":      64,
	"hear":    32,
	"a":       16,
	"bit":     1024,
}

func makeTestTrie() trie.Trie[string, int] {
	return trie.From[int](testMap)
}

func TestEmptyTrie(t *testing.T) {
	as := assert.New(t)

	tr := trie.New[string, int]()
	as.NotNil(tr)
	as.Equal(0, tr.Count())
	as.True(tr.IsEmpty())

	v, ok := tr.Get("hello")
	as.False(ok)
	as.Equal(0, v)

	v, r, ok := tr.Remove("blah")
	as.False(ok)
	as.Equal(0, v)
	as.Equal(r, tr)

	f := tr.First()
	as.Nil(f)

	r = tr.Rest()
	as.Equal(r, tr)

	r = tr.Put("hello", 10)
	as.NotEqual(r, tr)
	as.Equal(1, r.Count())
	as.False(r.IsEmpty())

	v, ok = r.Get("hello")
	as.True(ok)
	as.Equal(10, v)
}

func TestRetrieval(t *testing.T) {
	as := assert.New(t)
	tr := makeTestTrie()
	as.Equal(len(testMap), tr.Count())

	for k, v := range testMap {
		res, ok := tr.Get(k)
		as.True(ok)
		as.Equal(v, res)
	}

	var prev string
	for f, r, ok := tr.Split(); ok; f, r, ok = r.Split() {
		res, ok := testMap[f.Key()]
		as.True(ok)
		as.Equal(f.Value(), res)
		as.Less(prev, f.Key())
		prev = f.Key()
	}

	r, ok := tr.Get("b")
	as.False(ok)
	as.Equal(0, r)
}

func TestReplacement(t *testing.T) {
	as := assert.New(t)

	t1 := makeTestTrie()
	res, ok := t1.Get("today")
	as.True(ok)
	as.Equal(4, res)

	t2 := t1.Put("today", 32)
	res, ok = t2.Get("today")
	as.True(ok)
	as.Equal(32, res)
	as.Equal(len(testMap), t2.Count())

	res, ok = t1.Get("today")
	as.True(ok)
	as.Equal(4, res)
	as.Equal(len(testMap), t1.Count())
}

func TestRemoval(t *testing.T) {
	as := assert.New(t)

	r := makeTestTrie()
	cnt := r.Count()
	as.Equal(len(testMap), cnt)

	var f int
	var ok bool
	or := r
	f, r, ok = r.Remove("bogus")
	as.False(ok)
	as.Equal(0, f)
	as.Equal(or, r)

	for k, v := range testMap {
		f, r, ok = r.Remove(k)
		as.True(ok)
		as.Equal(v, f)
		cnt--
		as.Equal(cnt, r.Count())
	}

	f, ok = r.Get("hello")
	as.False(ok)
	as.True(r.IsEmpty())
	as.Equal(0, f)

	or = r
	f, r, ok = r.Remove("hello")
	as.False(ok)
	as.Equal(0, f)
	as.Equal(or, r)

	r = r.Put("hello", 96)
	f, ok = r.Get("hello")
	as.True(ok)
	as.Equal(96, f)
}

func TestSplitting(t *testing.T) {
	as := assert.New(t)

	check := map[string]int{}
	for k, v := range testMap {
		check[k] = v
	}

	tr := makeTestTrie()
	for !tr.IsEmpty() {
		f := tr.First()
		tr = tr.Rest()

		v := f.Value()
		as.Equal(v, check[f.Key()])
		delete(check, f.Key())
	}

	as.Equal(0, len(check))
}

func TestRemovePrefix(t *testing.T) {
	as := assert.New(t)

	t1 := makeTestTrie()
	t2, ok := t1.RemovePrefix("h")

	as.True(ok)
	c := t1.Count()
	as.Equal(c-3, t2.Count())
}
