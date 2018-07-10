package types

import (
	"fmt"
	"strconv"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/oneiro-ndev/ndaumath/pkg/signed"
)

//go:generate msgp -tests=0

// Ndau is a value that holds a single amount
// of ndau. Unlike an int64, it is prevented from overflowing.
type Ndau int64

// Add adds a value to an Ndau
// It may return an overflow error
func (n Ndau) Add(other Ndau) (Ndau, error) {
	t, err := signed.Add(int64(n), int64(other))
	return Ndau(t), err
}

// Sub subtracts, and may overflow
func (n Ndau) Sub(other Ndau) (Ndau, error) {
	t, err := signed.Sub(int64(n), int64(other))
	return Ndau(t), err
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
func (n Ndau) Abs() Ndau {
	y := n >> 63       // sign extended, so this is either -1 (0xFFF...) or 0
	return (n ^ y) - y // twos complement if it was negative
}

// Compare is the sorting operator; it returns -1 if n < rhs, 1 if n > rhs,
// and 0 if they are equal.
func (n Ndau) Compare(rhs Ndau) int {
	if n < rhs {
		return -1
	} else if n > rhs {
		return 1
	}
	return 0
}

// String returns the value of n formatted in a standard format, as if it is a
// decimal value of ndau. The full napu value is displayed, but trailing zeros
// are suppressed.
func (n Ndau) String() string {
	var sign int64 = 1
	if n < 0 {
		sign = -1
	}
	na := n.Abs()
	ndau := na / constants.QuantaPerUnit
	napu := na % constants.QuantaPerUnit
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
