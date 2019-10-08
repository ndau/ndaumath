package bitset256

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"encoding/hex"
	"errors"
	"math/bits"
)

// The bitset256 package supports an efficient array of 256 boolean values with
// associated operations like get, set, intersection, union, as well as
// conversion to/from strings and arrays of bytes. The bits are stored in an array of 4
// 64-bit words, in little-endian word order (the 0 bit is the 0 bit of the 0th
// word).
//
// It is not implemented with a third party bitset package because the ones I
// looked at allowed for arbitrary sizes; we only needed a 256-bit one and
// performance is pretty important in this use case.
//
// It is intended for use in chaincode to manage the list of valid opcodes.

// Bitset256 is an efficient way to store individual bits corresponding to 256
// values (i.e., using a byte as an index).
type Bitset256 [4]uint64

// New creates a new bitset and allows setting some of its bits at the same time.
func New(ixs ...byte) *Bitset256 {
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
// The index is taken mod256.
func wmask(ix byte) (int, uint64) {
	w := int(ix) >> 6 // faster divide by 64
	mask := uint64(1) << byte(ix&0x3F)
	return w, mask
}

// Get retrieves the value of a single bit at the given index.
func (b *Bitset256) Get(ix byte) bool {
	w, mask := wmask(ix)
	return (b[w] & mask) != 0
}

// Set unconditionally forces a single bit at the index to 1 and returns the pointer to the bitset.
func (b *Bitset256) Set(ix byte) *Bitset256 {
	w, mask := wmask(ix)
	b[w] |= mask
	return b
}

// Clear unconditionally forces a single bit to 0 and returns the pointer to the bitset.
func (b *Bitset256) Clear(ix byte) *Bitset256 {
	w, mask := wmask(ix)
	b[w] &= ^mask
	return b
}

// Toggle reverses the state of a single bit at the index and returns the pointer to the bitset.
func (b *Bitset256) Toggle(ix byte) *Bitset256 {
	w, mask := wmask(ix)
	b[w] ^= mask
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

// Less returns true if, when expressed as a number,
// b would be strictly less than other.
func (b *Bitset256) Less(other *Bitset256) bool {
	for i := 3; i >= 0; i-- {
		// if they're equal, move along
		if b[i] == other[i] {
			continue
		}
		// otherwise return the result of the comparison
		return b[i] < other[i]
	}
	return false
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

// Indices returns an []byte where the values are the indices of all the 1 bits that are set,
// in sorted order from 0. The length of the slice is equal to b.Count().
// This is fairly heavily optimized.
func (b *Bitset256) Indices() []byte {
	n := b.Count()
	result := make([]byte, n)
	c := 0
	for i := 0; i < 4; i++ {
		x := b[i]
		for x != 0 {
			lowest := bits.TrailingZeros64(x)
			if lowest == 64 {
				continue
			}
			result[c] = byte((i * 64) + lowest)
			c++
			if c == n {
				return result
			}
			m := uint64(1) << uint(lowest)
			x &= ^m
		}
	}
	return result
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
