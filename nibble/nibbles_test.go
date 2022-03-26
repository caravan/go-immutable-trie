package nibble_test

import (
	"testing"

	"github.com/caravan/go-immutable-trie/nibble"
	"github.com/stretchr/testify/assert"
)

func TestLargeNibbles(t *testing.T) {
	as := assert.New(t)

	n := nibble.Make("hello")
	as.NotNil(n)
	as.Equal(0, n.ByteOffset())

	f, r, ok := n.Consume()
	as.Equal(uint8(0x6), f)
	as.NotNil(r)
	as.True(ok)
	as.Equal(0, r.ByteOffset())

	f, r, ok = r.Consume()
	as.Equal(uint8(0x8), f)
	as.NotNil(r)
	as.True(ok)
	as.Equal(1, r.ByteOffset())
}

func TestSmallNibbles(t *testing.T) {
	as := assert.New(t)

	n := nibble.Make([]byte{'h'})
	as.NotNil(n)
	as.Equal(0, n.ByteOffset())

	f, r, ok := n.Consume()
	as.Equal(uint8(0x6), f)
	as.NotNil(r)
	as.True(ok)
	as.Equal(0, r.ByteOffset())

	f, r, ok = r.Consume()
	as.Equal(uint8(0x8), f)
	as.NotNil(r)
	as.True(ok)
	as.Equal(1, r.ByteOffset())

	f, r, ok = r.Consume()
	as.Equal(uint8(0x0), f)
	as.False(ok)
	as.Equal(1, r.ByteOffset())
}

func TestEmptyNibbles(t *testing.T) {
	as := assert.New(t)

	n := nibble.Make([]byte{})
	as.NotNil(n)
	as.Equal(0, n.ByteOffset())

	f, r, ok := n.Consume()
	as.Equal(uint8(0x0), f)
	as.Equal(n, r)
	as.NotNil(r)
	as.False(ok)
	as.Equal(0, r.ByteOffset())
}
