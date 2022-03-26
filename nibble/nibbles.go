package nibble

import "github.com/caravan/go-immutable-trie/key"

type (
	// Nibbles is an interface used to consume 4-bit values from a Key
	Nibbles[Key key.Keyable] interface {
		Consume() (uint8, Nibbles[Key], bool)
		ByteOffset() int
		Branch(Key) Nibbles[Key]
	}

	nibbles[Key key.Keyable] struct {
		data Key
		off  int
	}

	highNibbles[Key key.Keyable]  struct{ nibbles[Key] }
	lowNibbles[Key key.Keyable]   struct{ nibbles[Key] }
	emptyNibbles[Key key.Keyable] struct{ nibbles[Key] }
)

const Size = 16

// Make constructs a new set of Nibbles from the provided Key
func Make[Key key.Keyable](k Key) Nibbles[Key] {
	if len(k) > 0 {
		n := makeNibbles[Key](k, 0)
		return &highNibbles[Key]{n}
	}
	return &emptyNibbles[Key]{}
}

func makeNibbles[Key key.Keyable](data Key, off int) nibbles[Key] {
	return nibbles[Key]{
		data: data,
		off:  off,
	}
}

func (n *nibbles[Key]) ByteOffset() int {
	return n.off
}

func (n *nibbles[Key]) Branch(k Key) Nibbles[Key] {
	b := makeNibbles[Key](k, n.off)
	return &highNibbles[Key]{b}
}

func (n *highNibbles[Key]) Consume() (uint8, Nibbles[Key], bool) {
	next := (lowNibbles[Key])(*n)
	return n.data[n.off] >> 4 & 0x0F, &next, true
}

func (n *lowNibbles[Key]) Consume() (uint8, Nibbles[Key], bool) {
	res := n.data[n.off] & 0x0F
	next := makeNibbles[Key](n.data, n.off+1)
	if next.off < len(n.data) {
		return res, &highNibbles[Key]{next}, true
	}
	return res, &emptyNibbles[Key]{next}, true
}

func (n *lowNibbles[Key]) Branch(k Key) Nibbles[Key] {
	b := makeNibbles[Key](k, n.off)
	return &lowNibbles[Key]{b}
}

func (n *emptyNibbles[Key]) Consume() (uint8, Nibbles[Key], bool) {
	return 0, n, false
}
