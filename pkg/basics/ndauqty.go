package basics

import (
	"fmt"
	"strconv"
)

// NdauQty is a value that holds a single amount
// of ndau. Unlike an int64, it is prevented from overflowing.
type NdauQty int64

// Add adds a value to an NdauQty
// It may return an overflow error
func (n NdauQty) Add(other NdauQty) (NdauQty, error) {
	t := n + other
	// if the signs are opposite there's no way it can overflow
	if (n > 0) == (other < 0) {
		return t, nil
	}
	// otherwise, if the sum doesn't have the same sign
	// we overflowed
	if (n > 0) == (t < 0) {
		return t, OverflowError{}
	}
	return t, nil
}

// Sub subtracts, and may overflow
func (n NdauQty) Sub(other NdauQty) (NdauQty, error) {
	return n.Add(-other)
}

// Abs returns the absolute value without converting to float
// NOTE THAT THIS FAILS IN THE CASE WHERE n == MinInt64 (this
// value acts as much like -0 as it does MinInt). In particular,
// the value consists of only the negative (high) bit and the
// rest are zeros.
//
// As quantities on the blockchain can't be negative, we are going
// to ignore this case in favor of simplicity.
//
// In particular, this function a) can be inlined, and b) has no
// conditionals.
func (n NdauQty) Abs() NdauQty {
	y := n >> 63       // sign extended, so this is either -1 (0xFFF...) or 0
	return (n ^ y) - y // twos complement if it was negative
}

// Compare is the sorting operator; it returns -1 if n < rhs, 1 if n > rhs,
// and 0 if they are equal.
func (n NdauQty) Compare(rhs NdauQty) int {
	if n < rhs {
		return -1
	} else if n > rhs {
		return 1
	}
	return 0
}

func (n NdauQty) String() string {
	var sign int64 = 1
	if n < 0 {
		sign = -1
	}
	na := n.Abs()
	ndau := na / QuantaPerUnit
	napu := na % QuantaPerUnit
	if napu == 0 {
		return strconv.FormatInt(int64(sign*int64(ndau)), 10)
	}
	s := fmt.Sprintf("%d.%08d", sign*int64(ndau), napu)
	t := len(s)
	// trim off trailing zeros
	for ; s[t-1] == '0'; t-- {
	}
	return s[:t]
}
