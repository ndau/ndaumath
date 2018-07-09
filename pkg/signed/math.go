package signed

import (
	"github.com/ericlagergren/decimal"
	"github.com/oneiro-ndev/ndaumath/pkg/types"
)

// Mul multiplies two int64s and errors if there is an overflow
func Mul(a, b int64) (int64, error) {
	x := decimal.WithContext(decimal.Context128).SetMantScale(a, 0)
	y := decimal.WithContext(decimal.Context128).SetMantScale(b, 0)
	x.Mul(x, y)
	ret, ok := x.Int64()
	if !ok {
		return 0, types.ErrorOverflow
	}
	return ret, nil
}

// Div divides two int64s and throws errors if there are problems
func Div(a, b int64) (int64, error) {
	if b == 0 {
		return 0, types.ErrorDivideByZero
	}

	x := decimal.WithContext(decimal.Context128).SetMantScale(a, 0)
	y := decimal.WithContext(decimal.Context128).SetMantScale(b, 0)
	x.QuoInt(x, y)
	ret, ok := x.Int64()
	if !ok {
		return 0, types.ErrorMath
	}
	return ret, nil
}

// Mod calculates the remainder of dividing a by b and returns errors
// if there are issues.
func Mod(a, b int64) (int64, error) {
	if b == 0 {
		return 0, types.ErrorDivideByZero
	}

	x := decimal.WithContext(decimal.Context128).SetMantScale(a, 0)
	y := decimal.WithContext(decimal.Context128).SetMantScale(b, 0)
	x.Rem(x, y)
	ret, ok := x.Int64()
	if !ok {
		return 0, types.ErrorMath
	}
	return ret, nil
}

// DivMod calculates the quotient and the remainder of dividing a by b,
// returns both, and and returns errors if there are issues.
func DivMod(a, b int64) (int64, int64, error) {
	if b == 0 {
		return 0, 0, types.ErrorDivideByZero
	}

	x := decimal.WithContext(decimal.Context128).SetMantScale(a, 0)
	y := decimal.WithContext(decimal.Context128).SetMantScale(b, 0)
	x.QuoRem(x, y, y)
	q, ok := x.Int64()
	if !ok {
		return 0, 0, types.ErrorMath
	}
	r, ok := y.Int64()
	if !ok {
		return 0, 0, types.ErrorMath
	}
	return q, r, nil
}

// MulDiv multiplies a int64 value by the ratio n/d without overflowing the int64,
// provided that the final result does not overflow. Returns error if the result
// cannot be converted back to int64.
func MulDiv(v, n, d int64) (int64, error) {
	if d == 0 {
		return 0, types.ErrorDivideByZero
	}

	x := decimal.WithContext(decimal.Context128).SetMantScale(v, 0)
	y := decimal.WithContext(decimal.Context128).SetMantScale(n, 0)
	z := decimal.WithContext(decimal.Context128).SetMantScale(d, 0)
	x.Mul(x, y)
	x.QuoInt(x, z)
	ret, ok := x.Int64()
	if !ok {
		return 0, types.ErrorOverflow
	}
	return ret, nil
}
