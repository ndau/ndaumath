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
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	// make sure an empty set acts like a bunch of zeros
	b := New()
	assert.False(t, b.Get(5))
	assert.Equal(t, 0, b.Count())
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", b.AsHex())
}

func TestSimple(t *testing.T) {
	// set, clear, and toggle a single bit, making sure that we get equivalent pointers back
	b := New()
	assert.False(t, b.Get(1))
	c := b.Set(1)
	assert.True(t, c.Get(1))
	assert.True(t, b.Get(1))
	assert.Equal(t, b, c)
	d := b.Clear(1)
	assert.False(t, d.Get(1))
	assert.False(t, b.Get(1))
	assert.Equal(t, b, d)
	e := d.Toggle(1)
	assert.True(t, e.Get(1))
	assert.True(t, b.Get(1))
	assert.Equal(t, d, e)
	f := d.Toggle(1)
	assert.False(t, f.Get(1))
	assert.False(t, c.Get(1))
	assert.Equal(t, d, f)
	// and just make sure that we're not setting all the bits at once
	assert.False(t, b.Get(2))
}

func TestClone(t *testing.T) {
	// make sure clones are distinct sets
	b := New().Set(1).Set(35)
	assert.Equal(t, 2, b.Count())
	c := b.Clone()
	assert.Equal(t, 2, c.Count())
	assert.Equal(t, b, c)
	c.Set(28)
	assert.Equal(t, 3, c.Count())
	assert.Equal(t, 2, b.Count())
}

func TestAllBits(t *testing.T) {
	// run through and check all indices
	b := New()
	for i := 0; i < 256; i++ {
		assert.False(t, b.Get(byte(i)))
		b.Set(byte(i))
		assert.True(t, b.Get(byte(i)))
		b.Clear(byte(i))
		assert.False(t, b.Get(byte(i)))
		b.Toggle(byte(i))
		assert.Equal(t, int(i+1), b.Count())
	}
}

// setMultiples is a helper function that iterates the set and sets every Nth bit
func setMultiples(n int) *Bitset256 {
	b := New()
	for i := 0; i < 256; i += n {
		b.Set(byte(i))
	}
	return b
}

func TestAsBytes(t *testing.T) {
	// check the AsBytes function to make sure it round trips properly
	// and errors if given bad data
	b := setMultiples(7)
	assert.Equal(t, 37, b.Count())
	ba := b.AsBytes()
	c, err := FromBytes(ba)
	assert.Nil(t, err)
	assert.Equal(t, b, c)
	d, err := FromBytes(ba[:len(ba)-1])
	assert.NotNil(t, err)
	assert.Nil(t, d)
}

func TestAsHex(t *testing.T) {
	// check AsHex for roundtrip and errors
	b := setMultiples(7)
	assert.Equal(t, 37, b.Count())
	s := b.AsHex()
	c, err := FromHex(s)
	assert.Nil(t, err)
	assert.Equal(t, b, c)
	d, err := FromHex(s[:len(s)-1])
	assert.NotNil(t, err)
	assert.Nil(t, d)
}

func TestSubset(t *testing.T) {
	// Tests the subset function
	b := New().Set(1).Set(2).Set(3)
	c := New().Set(1).Set(2)
	assert.Equal(t, 3, b.Count())
	assert.Equal(t, 2, c.Count())
	assert.True(t, c.IsSubsetOf(b))
	assert.False(t, b.IsSubsetOf(c))
	assert.True(t, c.IsSubsetOf(c))
	assert.True(t, b.IsSubsetOf(b))
}

func TestNewMulti(t *testing.T) {
	// Checks that New() with multiple arguments does the right thing
	b := New(1, 2, 3)
	c := New().Set(1).Set(2)
	assert.Equal(t, 3, b.Count())
	assert.Equal(t, 2, c.Count())
	assert.True(t, c.IsSubsetOf(b))
	assert.False(t, b.IsSubsetOf(c))
	assert.True(t, c.IsSubsetOf(c))
	assert.True(t, b.IsSubsetOf(b))
}

func TestIntersect(t *testing.T) {
	fizz := setMultiples(3)
	buzz := setMultiples(5)
	fizzbuzz := fizz.Intersect(buzz)
	assert.Equal(t, 18, fizzbuzz.Count())
	all := fizz.Union(buzz)
	assert.Equal(t, 120, all.Count())
	assert.True(t, fizzbuzz.IsSubsetOf(fizz))
	assert.True(t, fizzbuzz.IsSubsetOf(buzz))
	assert.True(t, fizzbuzz.IsSubsetOf(fizzbuzz))
	assert.True(t, fizzbuzz.IsSubsetOf(all))
}

func TestUnion(t *testing.T) {
	fizz := setMultiples(3)
	assert.Equal(t, 256/3+1, fizz.Count())
	buzz := setMultiples(5)
	assert.Equal(t, 256/5+1, buzz.Count())
	all := fizz.Union(buzz)
	assert.Equal(t, 256/3+1, fizz.Count())
	assert.Equal(t, 256/5+1, buzz.Count())
	assert.Equal(t, 120, all.Count())
	assert.True(t, fizz.IsSubsetOf(all))
	assert.True(t, buzz.IsSubsetOf(all))
	assert.True(t, all.IsSubsetOf(all))
}

func TestAsHexString(t *testing.T) {
	// Checks that the hex function gets endianness correct
	b := setMultiples(4)
	assert.Equal(t, "1111111111111111111111111111111111111111111111111111111111111111", b.AsHex())
	c := New(0, 17, 34, 51)
	assert.Equal(t, "0000000000000008000000000000000400000000000000020000000000000001", c.AsHex())
	for i := 0; i < 256; i++ {
		c.Toggle(byte(i))
	}
	assert.Equal(t, "fffffffffffffff7fffffffffffffffbfffffffffffffffdfffffffffffffffe", c.AsHex())
	assert.Equal(t, c.String(), c.AsHex())
}

func TestCompare(t *testing.T) {
	b1 := New(1)
	b2 := New(2)
	assert.True(t, b1.Less(b2))
	assert.False(t, b2.Less(b1))
	assert.False(t, b2.Equals(b1))
	assert.False(t, b1.Equals(b2))
	b1 = New(2)
	b2 = New(2)
	assert.False(t, b1.Less(b2))
	assert.False(t, b2.Less(b1))
	assert.True(t, b2.Equals(b1))
	assert.True(t, b1.Equals(b2))
	b1 = New(200, 5)
	b2 = New(100, 9)
	assert.False(t, b1.Less(b2))
	assert.True(t, b2.Less(b1))
	assert.False(t, b2.Equals(b1))
	assert.False(t, b1.Equals(b2))
}

func TestIndices(t *testing.T) {
	for i := 0; i < 10; i++ {
		b1 := New()
		for j := 0; j < i; j++ {
			x := rand.Intn(256)
			b1.Set(byte(x))
		}
		ind := b1.Indices()
		assert.Equal(t, b1.Count(), len(ind))
		b2 := New(ind...)
		assert.True(t, b1.Equals(b2))
	}
}
