package trie_test

import (
	"testing"

	trie "github.com/caravan/go-immutable-trie"
	"github.com/stretchr/testify/assert"
)

func TestEqualKeys(t *testing.T) {
	as := assert.New(t)

	as.True(trie.EqualKeys[string]("hello", "hello"))
	as.False(trie.EqualKeys[string]("hell", "hello"))
}
