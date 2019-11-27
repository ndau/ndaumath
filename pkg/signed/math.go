package signed

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"math/big"

	"github.com/oneiro-ndev/ndaumath/pkg/ndauerr"
)

func op(a, b int64, operand func(*big.Int, *big.Int) *big.Int) (int64, error) {
	x := big.NewInt(a)
	y := big.NewInt(b)
	x = operand(x, y)
	if !x.IsInt64() {
		return 0, ndauerr.ErrOverflow
	}
	return x.Int64(), nil
}

// Add adds two int64s and errors if there is an overflow
func Add(a, b int64) (int64, error) {
	return op(a, b, func(x, y *big.Int) *big.Int { return x.Add(x, y) })
}

// Sub subtracts two int64s and errors if there is an overflow
func Sub(a, b int64) (int64, error) {
	return op(a, b, func(x, y *big.Int) *big.Int { return x.Sub(x, y) })
}

// Mul multiplies two int64s and errors if there is an overflow
func Mul(a, b int64) (int64, error) {
	return op(a, b, func(x, y *big.Int) *big.Int { return x.Mul(x, y) })
}

// Div divides two int64s and throws errors if there are problems
func Div(a, b int64) (int64, error) {
	if b == 0 {
		return 0, ndauerr.ErrDivideByZero
	}

	return op(a, b, func(x, y *big.Int) *big.Int { return x.Div(x, y) })
}

// Mod calculates the remainder of dividing a by b and returns errors
// if there are issues.
func Mod(a, b int64) (int64, error) {
	if b == 0 {
		return 0, ndauerr.ErrDivideByZero
	}

	return op(a, b, func(x, y *big.Int) *big.Int { return x.Mod(x, y) })
}

// DivMod calculates the quotient and the remainder of dividing a by b,
// returns both, and and returns errors if there are issues.
func DivMod(a, b int64) (int64, int64, error) {
	if b == 0 {
		return 0, 0, ndauerr.ErrDivideByZero
	}

	x := big.NewInt(a)
	y := big.NewInt(b)
	q, r := x.QuoRem(x, y, big.NewInt(0))
	if !q.IsInt64() || !r.IsInt64() {
		return 0, 0, ndauerr.ErrOverflow
	}
	return q.Int64(), r.Int64(), nil
}

// MulDiv multiplies a int64 value by the ratio n/d without overflowing the int64,
// provided that the final result does not overflow. Returns error if the result
// cannot be converted back to int64.
func MulDiv(v, n, d int64) (int64, error) {
	if d == 0 {
		return 0, ndauerr.ErrDivideByZero
	}

	x := big.NewInt(v)
	y := big.NewInt(n)
	z := big.NewInt(d)
	x = x.Mul(x, y)
	x = x.Quo(x, z)
	if !x.IsInt64() {
		return 0, ndauerr.ErrOverflow
	}
	return x.Int64(), nil
}
