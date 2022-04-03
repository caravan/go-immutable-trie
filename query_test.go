package trie_test

import (
	"testing"

	trie "github.com/caravan/go-immutable-trie"
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
		as.Less(i, len(entries))
		as.Equal(entries[i].key, k)
		as.Equal(entries[i].value, v)
		i++
	})
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
}

func TestDescending(t *testing.T) {
	q := makeTestTrie().Select().Descending().All()
	testResults(t, q, []testEntry{
		{"you", 37},
		{"today", 4},
		{"to", 64},
		{"there", 2},
		{"hello", 1},
		{"how", 9},
		{"hear", 32},
		{"curious", 128},
		{"bit", 1024},
		{"are", 5},
		{"a", 16},
	})
}
