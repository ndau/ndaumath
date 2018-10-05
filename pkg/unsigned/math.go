package unsigned

import (
	"github.com/ericlagergren/decimal"
	"github.com/oneiro-ndev/ndaumath/pkg/ndauerr"
)

// makeDecimal constructs a decimal object from a uint64
func makeDecimal(n uint64) *decimal.Big {
	return decimal.WithContext(decimal.Context128).SetUint64(n)
}

// Add adds two uint64s and errors if there is an overflow
func Add(a, b uint64) (uint64, error) {
	x := makeDecimal(a)
	y := makeDecimal(b)
	x.Add(x, y)
	ret, ok := x.Uint64()
	if !ok {
		return 0, ndauerr.ErrOverflow
	}
	return ret, nil
}

// Sub adds two uint64s and errors if there is an overflow
func Sub(a, b uint64) (uint64, error) {
	x := makeDecimal(a)
	y := makeDecimal(b)
	x.Sub(x, y)
	ret, ok := x.Uint64()
	if !ok {
		return 0, ndauerr.ErrOverflow
	}
	return ret, nil
}

// Mul multiplies two uint64s and errors if there is an overflow
func Mul(a, b uint64) (uint64, error) {
	x := makeDecimal(a)
	y := makeDecimal(b)
	x.Mul(x, y)
	ret, ok := x.Uint64()
	if !ok {
		return 0, ndauerr.ErrOverflow
	}
	return ret, nil
}

// Div divides two uint64s and throws errors if there are problems
func Div(a, b uint64) (uint64, error) {
	if b == 0 {
		return 0, ndauerr.ErrDivideByZero
	}

	x := makeDecimal(a)
	y := makeDecimal(b)
	x.QuoInt(x, y)
	ret, ok := x.Uint64()
	if !ok {
		return 0, ndauerr.ErrMath
	}
	return ret, nil
}

// Mod calculates the remainder of dividing a by b and returns errors
// if there are issues.
func Mod(a, b uint64) (uint64, error) {
	if b == 0 {
		return 0, ndauerr.ErrDivideByZero
	}

	x := makeDecimal(a)
	y := makeDecimal(b)
	x.Rem(x, y)
	ret, ok := x.Uint64()
	if !ok {
		return 0, ndauerr.ErrMath
	}
	return ret, nil
}

// DivMod calculates the quotient and the remainder of dividing a by b,
// returns both, and and returns errors if there are issues.
func DivMod(a, b uint64) (uint64, uint64, error) {
	if b == 0 {
		return 0, 0, ndauerr.ErrDivideByZero
	}

	x := makeDecimal(a)
	y := makeDecimal(b)
	x.QuoRem(x, y, y)
	q, ok := x.Uint64()
	if !ok {
		return 0, 0, ndauerr.ErrMath
	}
	r, ok := y.Uint64()
	if !ok {
		return 0, 0, ndauerr.ErrMath
	}
	return q, r, nil
}

// MulDiv multiplies a uint64 value by the ratio n/d without overflowing the uint64,
// provided that the final result does not overflow. Returns error if the result
// cannot be converted back to uint64.
func MulDiv(v, n, d uint64) (uint64, error) {
	if d == 0 {
		return 0, ndauerr.ErrDivideByZero
	}

	x := makeDecimal(v)
	y := makeDecimal(n)
	z := makeDecimal(d)
	x.Mul(x, y)
	x.QuoInt(x, z)
	ret, ok := x.Uint64()
	if !ok {
		return 0, ndauerr.ErrOverflow
	}
	return ret, nil
}
