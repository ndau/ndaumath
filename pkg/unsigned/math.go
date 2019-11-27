package unsigned

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

func bigu(x uint64) *big.Int {
	b := new(big.Int)
	b.SetUint64(x)
	return b
}

func op(a, b uint64, operand func(*big.Int, *big.Int) *big.Int) (uint64, error) {
	x := bigu(a)
	y := bigu(b)
	x = operand(x, y)
	if !x.IsUint64() {
		return 0, ndauerr.ErrOverflow
	}
	return x.Uint64(), nil
}

// Add adds two int64s and errors if there is an overflow
func Add(a, b uint64) (uint64, error) {
	return op(a, b, func(x, y *big.Int) *big.Int { return x.Add(x, y) })
}

// Sub subtracts two int64s and errors if there is an overflow
func Sub(a, b uint64) (uint64, error) {
	return op(a, b, func(x, y *big.Int) *big.Int { return x.Sub(x, y) })
}

// Mul multiplies two int64s and errors if there is an overflow
func Mul(a, b uint64) (uint64, error) {
	return op(a, b, func(x, y *big.Int) *big.Int { return x.Mul(x, y) })
}

// Div divides two int64s and throws errors if there are problems
func Div(a, b uint64) (uint64, error) {
	if b == 0 {
		return 0, ndauerr.ErrDivideByZero
	}

	return op(a, b, func(x, y *big.Int) *big.Int { return x.Div(x, y) })
}

// Mod calculates the remainder of dividing a by b and returns errors
// if there are issues.
func Mod(a, b uint64) (uint64, error) {
	if b == 0 {
		return 0, ndauerr.ErrDivideByZero
	}

	return op(a, b, func(x, y *big.Int) *big.Int { return x.Mod(x, y) })
}

// DivMod calculates the quotient and the remainder of dividing a by b,
// returns both, and and returns errors if there are issues.
func DivMod(a, b uint64) (uint64, uint64, error) {
	if b == 0 {
		return 0, 0, ndauerr.ErrDivideByZero
	}

	x := bigu(a)
	y := bigu(b)
	q, r := x.QuoRem(x, y, big.NewInt(0))
	if !q.IsUint64() || !r.IsUint64() {
		return 0, 0, ndauerr.ErrOverflow
	}
	return q.Uint64(), r.Uint64(), nil
}

// MulDiv multiplies a uint64 value by the ratio n/d without overflowing the uint64,
// provided that the final result does not overflow. Returns error if the result
// cannot be converted back to uint64.
func MulDiv(v, n, d uint64) (uint64, error) {
	if d == 0 {
		return 0, ndauerr.ErrDivideByZero
	}

	x := bigu(v)
	y := bigu(n)
	z := bigu(d)
	x = x.Mul(x, y)
	x = x.Quo(x, z)
	if !x.IsUint64() {
		return 0, ndauerr.ErrOverflow
	}
	return x.Uint64(), nil
}
