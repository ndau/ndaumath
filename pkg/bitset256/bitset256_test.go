package bitset256

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	b := New()
	assert.False(t, b.Get(5))
	assert.Equal(t, 0, b.Count())
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", b.AsHex())
}

func TestSimple(t *testing.T) {
	b := New()
	assert.False(t, b.Get(1))
	c := b.Set(1)
	assert.Equal(t, b, c)
	assert.True(t, b.Get(1))
	fmt.Println(b)
	b.Clear(1)
	fmt.Println(b)
	assert.False(t, b.Get(1))
	assert.False(t, b.Get(2))
}

func TestClone(t *testing.T) {
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
	b := New()
	for i := 0; i < 256; i++ {
		assert.False(t, b.Get(i))
		b.Set(i)
		assert.True(t, b.Get(i))
		b.Clear(i)
		assert.False(t, b.Get(i))
		b.Set(i)
		assert.Equal(t, int(i+1), b.Count())
	}
}

func setMultiples(n int) *Bitset256 {
	b := New()
	for i := 0; i < 256; i += n {
		b.Set(i)
	}
	return b
}

func TestAsBytes(t *testing.T) {
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
	b := New(1, 2, 3)
	c := New(1, 2)
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

func TestAsHex2(t *testing.T) {
	b := setMultiples(4)
	assert.Equal(t, "1111111111111111111111111111111111111111111111111111111111111111", b.AsHex())
	c := New(0, 17, 34, 51)
	assert.Equal(t, "0000000000000008000000000000000400000000000000020000000000000001", c.AsHex())
	for i := 0; i < 256; i++ {
		c.Toggle(i)
	}
	assert.Equal(t, "fffffffffffffff7fffffffffffffffbfffffffffffffffdfffffffffffffffe", c.AsHex())
}
