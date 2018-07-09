package unsigned

import (
	"github.com/ericlagergren/decimal"
	"github.com/oneiro-ndev/ndaumath/pkg/types"
)

// Mul multiplies two uint64s and errors if there is an overflow
func Mul(a, b uint64) (uint64, error) {
	x := decimal.WithContext(decimal.Context128).SetUint64(a)
	y := decimal.WithContext(decimal.Context128).SetUint64(b)
	x.Mul(x, y)
	ret, ok := x.Uint64()
	if !ok {
		return 0, types.ErrorOverflow
	}
	return ret, nil
}

// Div divides two uint64s and throws errors if there are problems
func Div(a, b uint64) (uint64, error) {
	if b == 0 {
		return 0, types.ErrorDivideByZero
	}

	x := decimal.WithContext(decimal.Context128).SetUint64(a)
	y := decimal.WithContext(decimal.Context128).SetUint64(b)
	x.QuoInt(x, y)
	ret, ok := x.Uint64()
	if !ok {
		return 0, types.ErrorMath
	}
	return ret, nil
}

// Mod calculates the remainder of dividing a by b and returns errors
// if there are issues.
func Mod(a, b uint64) (uint64, error) {
	if b == 0 {
		return 0, types.ErrorDivideByZero
	}

	x := decimal.WithContext(decimal.Context128).SetUint64(a)
	y := decimal.WithContext(decimal.Context128).SetUint64(b)
	x.Rem(x, y)
	ret, ok := x.Uint64()
	if !ok {
		return 0, types.ErrorMath
	}
	return ret, nil
}

// DivMod calculates the quotient and the remainder of dividing a by b,
// returns both, and and returns errors if there are issues.
func DivMod(a, b uint64) (uint64, uint64, error) {
	if b == 0 {
		return 0, 0, types.ErrorDivideByZero
	}

	x := decimal.WithContext(decimal.Context128).SetUint64(a)
	y := decimal.WithContext(decimal.Context128).SetUint64(b)
	x.QuoRem(x, y, y)
	q, ok := x.Uint64()
	if !ok {
		return 0, 0, types.ErrorMath
	}
	r, ok := y.Uint64()
	if !ok {
		return 0, 0, types.ErrorMath
	}
	return q, r, nil
}

// MulDiv multiplies a uint64 value by the ratio n/d without overflowing the uint64,
// provided that the final result does not overflow. Returns error if the result
// cannot be converted back to uint64.
func MulDiv(v, n, d uint64) (uint64, error) {
	if d == 0 {
		return 0, types.ErrorDivideByZero
	}

	x := decimal.WithContext(decimal.Context128).SetUint64(v)
	y := decimal.WithContext(decimal.Context128).SetUint64(n)
	z := decimal.WithContext(decimal.Context128).SetUint64(d)
	x.Mul(x, y)
	x.QuoInt(x, z)
	ret, ok := x.Uint64()
	if !ok {
		return 0, types.ErrorOverflow
	}
	return ret, nil
}
