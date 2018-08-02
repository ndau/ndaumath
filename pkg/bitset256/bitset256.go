package bitset256

import (
	"encoding/hex"
	"errors"
	"math/bits"
)

// Bitset256 is an efficient way to store individual bits corresponding to 256 values (i.e.,
// using a byte as an index). The bits are stored in an array of 4 64-bit words, in little-endian
// word order (the 0 bit is the 0 bit of the 0th word).
type Bitset256 [4]uint64

// New creates a new bitset and allows setting some of its bits at the same time.
func New(ixs ...int) *Bitset256 {
	b := &Bitset256{}
	for _, i := range ixs {
		b.Set(i)
	}
	return b
}

// Clone creates a copy of a bitset.
func (b *Bitset256) Clone() *Bitset256 {
	c := *b
	return &c
}

// wmask is a helper function which, given an index,
// returns the word to index into and the mask to use for selecting the given bit
func wmask(ix int) (int, uint64) {
	w := (ix & 0xFF) >> 6 // faster divide by 64
	mask := uint64(1) << byte(ix&0x3F)
	return w, mask
}

// Get retrieves the value of a single bit at the given index.
// The index is taken mod256.
func (b *Bitset256) Get(ix int) bool {
	w, mask := wmask(ix)
	return (b[w] & mask) != 0
}

// Set unconditionally forces a single bit at the index to 1 and returns the pointer to the bitset.
// The index is taken mod256.
func (b *Bitset256) Set(ix int) *Bitset256 {
	w, mask := wmask(ix)
	b[w] |= mask
	return b
}

// Clear unconditionally forces a single bit to 0 and returns the pointer to the bitset.
// The index is taken mod256.
func (b *Bitset256) Clear(ix int) *Bitset256 {
	w, mask := wmask(ix)
	b[w] &= ^mask
	return b
}

// Equals returns true if the two bitsets have identical contents.
func (b *Bitset256) Equals(other *Bitset256) bool {
	for i := 0; i < 4; i++ {
		if b[i] != other[i] {
			return false
		}
	}
	return true
}

// Intersect returns a pointer to a new Bitset256 that is the intersection
// of its two source bitsets (the only bits that are set are the ones where
// both source sets had a 1 bit).
func (b *Bitset256) Intersect(other *Bitset256) *Bitset256 {
	r := b.Clone()
	for i := 0; i < 4; i++ {
		r[i] &= other[i]
	}
	return r
}

// Union returns a pointer to a new Bitset256 that is the union
// of its two source bitsets (the only bits that are set are the ones where
// either source set had a 1 bit).
func (b *Bitset256) Union(other *Bitset256) *Bitset256 {
	r := b.Clone()
	for i := 0; i < 4; i++ {
		r[i] |= other[i]
	}
	return r
}

// IsSubsetOf returns true if all of the bits in a bitset are also in the other bitset.
func (b *Bitset256) IsSubsetOf(other *Bitset256) bool {
	return b.Equals(b.Intersect(other))
}

// Count returns the number of 1 bits that are set.
func (b *Bitset256) Count() int {
	c := 0
	for i := 0; i < 4; i++ {
		c += bits.OnesCount64(b[i])
	}
	return c
}

// AsBytes returns the bitset as a slice of 32 bytes, where the 0 bits in the bitset are in the
// last element of the slice (basically, big-endian format). This is so that rendering the slice
// to a visual format will show the bits in an expected order.
func (b *Bitset256) AsBytes() []byte {
	ba := make([]byte, 32)
	for i := uint(0); i < 4; i++ {
		for j := uint(0); j < 8; j++ {
			ba[j*4+i] = byte((b[3-i] >> ((7 - j) * 8)) & 0xFF)
		}
	}
	return ba
}

// FromBytes takes a slice of 32 bytes and builds a Bitset256 from it, following
// the same rules as AsBytes (the last byte in the slice corresponds to the zeroth bits).
func FromBytes(ba []byte) (*Bitset256, error) {
	if len(ba) != 32 {
		return nil, errors.New("wrong number of bytes")
	}
	b := New()
	for i := uint(0); i < 4; i++ {
		for j := uint(0); j < 8; j++ {
			b[3-i] |= uint64(ba[j*4+i]) << ((7 - j) * 8)
		}
	}
	return b, nil
}

// AsHex returns a string representation of the bitset as a 256-bit number in hex (64 characters).
// It follows the same rules as AsBytes.
func (b *Bitset256) AsHex() string {
	return hex.EncodeToString(b.AsBytes())
}

// FromHex builds a Bitset256 from a hex string like the one AsHex generates.
func FromHex(s string) (*Bitset256, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return FromBytes(b)
}

// String implements Stringer for Bitset256.
func (b *Bitset256) String() string {
	return b.AsHex()
}
