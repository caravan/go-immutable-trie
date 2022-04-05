package trie_test

import (
	"fmt"
	"sort"
	"testing"

	trie "github.com/caravan/go-immutable-trie"
	"github.com/caravan/go-immutable-trie/key"
	"github.com/stretchr/testify/assert"
)

type testEntry struct {
	key   string
	value int
}

func testResults(t *testing.T, q trie.Query[string, int], entries []testEntry) {
	as := assert.New(t)
	i := 0
	q.ForEach(func(k string, v int) {
		as.Equal(entries[i].key, k)
		as.Equal(entries[i].value, v)
		i++
	})
	as.Equal(len(entries), i)
}

func TestWhereQuery(t *testing.T) {
	q := makeTestTrie().Select().All().Where(func(k string, v int) bool {
		return k[0] == 'h'
	})

	testResults(t, q, []testEntry{
		{"hear", 32},
		{"hello", 1},
		{"how", 9},
	})
}

func TestWhileQuery(t *testing.T) {
	q := makeTestTrie().Select().All().While(func(k string, v int) bool {
		return k != "bit"
	})

	testResults(t, q, []testEntry{
		{"a", 16},
		{"are", 5},
	})
}

func TestFromQuery(t *testing.T) {
	q := makeTestTrie().Select().Ascending().From("today")

	testResults(t, q, []testEntry{
		{"today", 4},
		{"you", 37},
	})
}

func TestAscending(t *testing.T) {
	q := makeTestTrie().Select().Ascending().All()
	testResults(t, q, []testEntry{
		{"a", 16},
		{"are", 5},
		{"bit", 1024},
		{"curious", 128},
		{"hear", 32},
		{"hello", 1},
		{"how", 9},
		{"there", 2},
		{"to", 64},
		{"today", 4},
		{"you", 37},
	})

	q = makeTestTrie().Select().From("hear")
	testResults(t, q, []testEntry{
		{"hear", 32},
		{"hello", 1},
		{"how", 9},
		{"there", 2},
		{"to", 64},
		{"today", 4},
		{"you", 37},
	})

	q = makeTestTrie().Select().From("heart")
	testResults(t, q, []testEntry{
		{"hello", 1},
		{"how", 9},
		{"there", 2},
		{"to", 64},
		{"today", 4},
		{"you", 37},
	})
}

func TestLargePrefixRemove(t *testing.T) {
	as := assert.New(t)

	in := make([]testEntry, 100000)
	for i := 0; i < len(in); i++ {
		in[i] = testEntry{
			key:   fmt.Sprintf("%d", i),
			value: i,
		}
	}
	m := map[string]int{}
	for _, e := range in {
		m[e.key] = e.value
	}

	t1 := trie.From[int](m)
	c1 := t1.Count()
	t2, ok := t1.RemovePrefix("2")
	as.True(ok)

	as.Equal(c1-11111, t2.Count())

	c3 := 0
	t2.Select().All().ForEach(func(k string, v int) {
		if k[0] == '2' {
			as.FailNow("key returned starting with a 2")
		}
		c3++
	})
	as.Equal(t2.Count(), c3)
	as.Equal(c1-11111, c3)
}

func TestLargeDescending(t *testing.T) {
	as := assert.New(t)

	in := make([]testEntry, 100000)
	for i := 0; i < len(in); i++ {
		in[i] = testEntry{
			key:   fmt.Sprintf("%d", i),
			value: i,
		}
	}
	m := map[string]int{}
	for _, e := range in {
		m[e.key] = e.value
	}

	sorted := in[:]
	sort.Slice(sorted, func(l, r int) bool {
		le := sorted[l]
		re := sorted[r]
		return key.Compare[string](le.key, re.key) == -1
	})

	tr := trie.From[int](m)
	as.Equal(len(m), tr.Count())

	i := len(in) - 1
	tr.Select().Descending().All().ForEach(func(k string, v int) {
		e := sorted[i]
		as.Equal(e.key, k)
		as.Equal(e.value, v)
		i--
	})
}

func TestDescendingAll(t *testing.T) {
	q := makeTestTrie().Select().Descending().All()
	testResults(t, q, []testEntry{
		{"you", 37},
		{"today", 4},
		{"to", 64},
		{"there", 2},
		{"how", 9},
		{"hello", 1},
		{"hear", 32},
		{"curious", 128},
		{"bit", 1024},
		{"are", 5},
		{"a", 16},
	})
}

func TestDescendingFrom(t *testing.T) {
	q := makeTestTrie().Select().Descending().From("hear")
	testResults(t, q, []testEntry{
		{"hear", 32},
		{"curious", 128},
		{"bit", 1024},
		{"are", 5},
		{"a", 16},
	})
}

func TestDescendingFromBetween(t *testing.T) {
	q := makeTestTrie().Select().Descending().From("heart")
	testResults(t, q, []testEntry{
		{"hear", 32},
		{"curious", 128},
		{"bit", 1024},
		{"are", 5},
		{"a", 16},
	})
}
