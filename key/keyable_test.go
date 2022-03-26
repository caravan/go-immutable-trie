package key_test

import (
	"testing"

	"github.com/caravan/go-immutable-trie/key"
	"github.com/stretchr/testify/assert"
)

func TestEqualKeys(t *testing.T) {
	as := assert.New(t)
	as.True(key.EqualTo[string]("hello", "hello"))
	as.False(key.EqualTo[string]("hell", "hello"))
}

func TestLessThanKeys(t *testing.T) {
	as := assert.New(t)
	as.True(key.LessThan[string]("ace", "barley"))
	as.False(key.LessThan[string]("barley", "ace"))
}

func TestGreaterThanKeys(t *testing.T) {
	as := assert.New(t)
	as.True(key.GreaterThan[string]("barley", "ace"))
	as.False(key.GreaterThan[string]("ace", "barley"))
}
