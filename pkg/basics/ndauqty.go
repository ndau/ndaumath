package basics

// Ndau is a value that holds a single amount
// of ndau. It is prevented from overflow
type Ndau int64

// Add adds a value to an NdauQty
// It may return an overflow error
func (n Ndau) Add(other Ndau) (Ndau, error) {
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
func (n Ndau) Sub(other Ndau) (Ndau, error) {
	return n.Add(-other)
}

// Abs returns the absolute value without converting to float
// NOTE THAT THIS FAILS IN THE CASE WHERE n == MinInt64 -- this
// value acts as much like -0 as MinInt. In particular, the value
// consists of only the negative (high) bit and the rest are
// zeros.
// As this is the only error case, and quantities on the
// blockchain can't be negative, we are going to ignore this
// case.
func (n Ndau) Abs() Ndau {
	y := n >> 63       // y ← x ⟫ 63
	return (n ^ y) - y // (x ⨁ y) - y
}
