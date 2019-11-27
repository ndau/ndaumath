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
	"math"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ericlagergren/decimal/v3"
)

func TestAdd(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"simple", args{6, 7}, 13, false},
		{"bigger than int32", args{600000000, 700000000}, 1300000000, false},
		{"max value adding 0", args{math.MaxUint64, 0}, math.MaxUint64, false},
		{"crossing int32 border", args{math.MaxUint32, math.MaxUint32}, 2 * math.MaxUint32, false},
		{"just barely overflowing", args{math.MaxUint64, 1}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSub(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"simple", args{13, 7}, 6, false},
		{"bigger than int32", args{1300000000, 700000000}, 600000000, false},
		{"max value subtracting 0", args{math.MaxUint64, 0}, math.MaxUint64, false},
		{"crossing int32 border", args{2 * math.MaxUint32, math.MaxUint32}, math.MaxUint32, false},
		{"underflowing", args{40, 61}, 0, true},
		{"just barely underflowing", args{0, 1}, 0, true},
		{"zero result", args{10, 10}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sub(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMul(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"simple", args{6, 7}, 42, false},
		{"bigger than int32", args{600000000, 700000000}, 420000000000000000, false},
		{"multiply by zero", args{10, 0}, 0, false},
		{"multiply by zero the other way", args{0, 10}, 0, false},
		{"too big to fit", args{600000000000, 700000000000}, 0, true},
		{"getting close to the limit", args{math.MaxUint32, math.MaxUint32}, math.MaxUint32 * math.MaxUint32, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mul(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"simple", args{42, 7}, 6, false},
		{"divide zero by", args{0, 7}, 0, false},
		{"big numbers", args{420000000000000000, 700000000}, 600000000, false},
		{"divide by zero", args{600000000000, 0}, 0, true},
		{"close to the limit", args{math.MaxUint32 * math.MaxUint32, math.MaxUint32}, math.MaxUint32, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Div(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Div() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMod(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"simple 0 remainder", args{42, 7}, 0, false},
		{"simple with remainder", args{42, 5}, 2, false},
		{"divide zero by", args{0, 7}, 0, false},
		{"big with remainder", args{420000000000000000, 700000001}, 100000001, false},
		{"divide by zero", args{12, 0}, 0, true},
		{"big", args{math.MaxUint32 * math.MaxUint32, math.MaxInt32}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mod(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDivMod(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		want1   uint64
		wantErr bool
	}{
		{"simple no remainder", args{42, 7}, 6, 0, false},
		{"simple with remainder", args{42, 5}, 8, 2, false},
		{"divide zero by", args{0, 7}, 0, 0, false},
		{"big numbers", args{420000000000000000, 700000001}, 599999999, 100000001, false},
		{"divide by zero", args{12, 0}, 0, 0, true},
		{"at the limit", args{math.MaxUint32 * math.MaxUint32, math.MaxInt32}, (math.MaxUint32 * math.MaxUint32) / math.MaxInt32, 1, false},
		{"zero result with remainder", args{42, 55}, 0, 42, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DivMod(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("DivMod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DivMod() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DivMod() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMulDiv(t *testing.T) {
	type args struct {
		v uint64
		n uint64
		d uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"simple with exact result", args{80, 2, 5}, 32, false},
		{"simple with result rounded down", args{82, 2, 5}, 32, false},
		{"simple with result rounded up", args{83, 2, 5}, 33, false},
		{"zero for v", args{0, 2, 5}, 0, false},
		{"zero for n", args{0, 2, 5}, 0, false},
		{"simple with numbers > maxint32", args{80000000000, 2000000000, 5000000000}, 32000000000, false},
		{"big number with small ratio", args{80000000000, 2, 5}, 32000000000, false},
		{"divide by zero", args{80000000000, 2, 0}, 0, true},
		{"approximate with ratio > 1", args{147, 155, 132}, 172, false},
		{"too big with ratio > 1", args{14717364050318377211, 15574702891736741942, 1324724618575407633}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MulDiv(tt.args.v, tt.args.n, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("MulDiv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MulDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func bigmuldiv(a, b, c uint64) uint64 {
	x := big.NewInt(0).SetUint64(a)
	y := big.NewInt(0).SetUint64(b)
	z := big.NewInt(0).SetUint64(c)
	r := big.NewInt(0)
	q := big.NewInt(0)
	r.Mul(x, y)
	q.Div(r, z)
	return q.Uint64()
}

func compareOne(r *rand.Rand, t *testing.T) {
	// make sure they're never negative
	a := r.Uint64() //& 0x7FFFFFFFFFFFFFFF
	b := r.Uint64() //& 0x7FFFFFFFFFFFFFFF
	c := r.Uint64() //& 0x7FFFFFFFFFFFFFFF
	if b > c {
		b, c = c, b
	}
	p, err := MulDiv(a, b, c)
	if err != nil {
		t.Error(err)
	}
	q := bigmuldiv(a, b, c)
	if p != q {
		t.Errorf("muldiv didn't match results from big.Int: (%v %v %v) %v != %v", a, b, c, p, q)
		t.Errorf("a=%x b=%x\n", a, b)
		a /= 10
		b /= 10
		c /= 10
		p, _ = MulDiv(a, b, c)
		q := bigmuldiv(a, b, c)
		t.Errorf("/10: (%v %v %v) %v == %v", a, b, c, p, q)
	}
}

func TestMulDivFuzz(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 10000; i++ {
		compareOne(r, t)
	}
}

func TestConversion(t *testing.T) {
	x := decimal.WithContext(decimal.Context128).SetUint64(math.MaxUint64)
	y, ok := x.Uint64()
	if !ok {
		t.Error("failed to convert back")
	}

	if y == math.MaxUint64 {
		t.Logf("the bug in the decimal library (https://github.com/ericlagergren/decimal/issues/104) has been fixed")
	} else {
		t.Error("bug in decimal library (https://github.com/ericlagergren/decimal/issues/104) remains but makeDecimal has already been nerfed")
	}
}
