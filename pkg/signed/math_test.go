package signed

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
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
)

func TestAdd(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"simple", args{6, 7}, 13, false},
		{"bigger than int32", args{600000000, 700000000}, 1300000000, false},
		{"max value adding 0", args{math.MaxInt64, 0}, math.MaxInt64, false},
		{"crossing int32 border", args{math.MaxUint32, math.MaxUint32}, 2 * math.MaxUint32, false},
		{"just barely overflowing", args{math.MaxInt64, 1}, 0, true},
		{"simple neg", args{1, -1}, 0, false},
		{"max-1", args{int64(math.MaxInt64), -1}, int64(math.MaxInt64 - 1), false},
		{"max+1", args{int64(math.MaxInt64), 1}, 0, true},
		{"max possible sum", args{int64(math.MaxInt64 / 2), int64(math.MaxInt64 / 2)}, int64(math.MaxInt64) - 1, false},
		{"max negative + 1", args{int64(math.MinInt64), 1}, int64(math.MinInt64 + 1), false},
		{"sum of max and min", args{int64(math.MaxInt64), int64(math.MinInt64)}, -1, false},
		{"half of min", args{int64(math.MinInt64 / 2), int64(math.MinInt64 / 2)}, int64(math.MinInt64), false},
		{"negative overflow", args{int64(math.MinInt64), -1}, 0, true},
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
		a int64
		b int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"subtract a negative", args{1, -1}, 2, false},
		{"simple sum", args{1, 1}, 0, false},
		{"result less than zero", args{1, 100}, -99, false},
		{"simple subtraction", args{654321, 123456}, 530865, false},
		{"close to max", args{int64(math.MaxInt64), 1}, int64(math.MaxInt64 - 1), false},
		{"overflow", args{int64(math.MaxInt64), -1}, 0, true},
		{"close to max result", args{int64(math.MaxInt64 / 2), -int64(math.MaxInt64 / 2)}, int64(math.MaxInt64 - 1), false},
		{"close to neg max", args{int64(math.MinInt64), -1}, int64(math.MinInt64 + 1), false},
		{"pos and neg maxes", args{int64(math.MaxInt64), int64(math.MaxInt64)}, 0, false},
		{"neg max", args{int64(math.MinInt64 / 2), -int64(math.MinInt64 / 2)}, int64(math.MinInt64), false},
		{"negative overflow", args{int64(math.MinInt64), 1}, 0, true}}
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
		a int64
		b int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"simple", args{6, 7}, 42, false},
		{"bigger than int32", args{600000000, 700000000}, 420000000000000000, false},
		{"multiply by zero", args{10, 0}, 0, false},
		{"multiply by zero the other way", args{0, 10}, 0, false},
		{"too big to fit", args{600000000000, 700000000000}, 0, true},
		{"at the limit", args{math.MaxInt32, math.MaxInt32}, math.MaxInt32 * math.MaxInt32, false},
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
		a int64
		b int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"simple", args{42, 7}, 6, false},
		{"divide zero by", args{0, 7}, 0, false},
		{"big numbers", args{420000000000000000, 700000000}, 600000000, false},
		{"divide by zero", args{600000000000, 0}, 0, true},
		{"at the limit", args{math.MaxInt32 * math.MaxInt32, math.MaxInt32}, math.MaxInt32, false},
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
		a int64
		b int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"simple 0 remainder", args{42, 7}, 0, false},
		{"simple with remainder", args{42, 5}, 2, false},
		{"divide zero by", args{0, 7}, 0, false},
		{"big with remainder", args{420000000000000000, 700000001}, 100000001, false},
		{"divide by zero", args{12, 0}, 0, true},
		{"big", args{math.MaxInt32 * math.MaxInt32, math.MaxInt32}, 0, false},
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
		a int64
		b int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		want1   int64
		wantErr bool
	}{
		{"simple no remainder", args{42, 7}, 6, 0, false},
		{"simple with remainder", args{42, 5}, 8, 2, false},
		{"divide zero by", args{0, 7}, 0, 0, false},
		{"big numbers", args{420000000000000000, 700000001}, 599999999, 100000001, false},
		{"divide by zero", args{12, 0}, 0, 0, true},
		{"zero result with remainder", args{42, 55}, 0, 42, false},
		{"at the limit", args{math.MaxInt32 * math.MaxInt32, math.MaxInt32}, math.MaxInt32, 0, false},
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
		v int64
		n int64
		d int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"simple with exact result", args{80, 2, 5}, 32, false},
		{"truncation 1", args{82, 2, 5}, 32, false},
		{"truncation 2", args{83, 2, 5}, 33, false},
		{"positive truncated towards zero", args{100, 1, 3}, 33, false},
		{"positive truncated toward zero 2", args{100, 2, 3}, 66, false},
		{"negative truncated towards zero", args{-100, 1, 3}, -33, false},
		{"negative truncated toward zero 2", args{-100, 2, 3}, -66, false},
		{"zero for v", args{0, 2, 5}, 0, false},
		{"zero for n", args{0, 2, 5}, 0, false},
		{"simple with numbers > maxint32", args{80000000000, 2000000000, 5000000000}, 32000000000, false},
		{"big number with small ratio", args{80000000000, 2, 5}, 32000000000, false},
		{"divide by zero", args{80000000000, 2, 0}, 0, true},
		{"approximate with ratio > 1", args{147, 155, 132}, 172, false},
		{"too big with ratio > 1", args{math.MaxInt64, 1557470289173674194, 132472461857540763}, 0, true},
		{"too big with ratio < 1", args{math.MaxInt64, 132472461857540763, 1557470289173674194}, 784504724644480276, false},
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

func bigmuldiv(a, b, c int64) int64 {
	x := big.NewInt(0).SetInt64(a)
	y := big.NewInt(0).SetInt64(b)
	z := big.NewInt(0).SetInt64(c)
	r := big.NewInt(0)
	q := big.NewInt(0)
	r.Mul(x, y)
	q.Div(r, z)
	return q.Int64()
}

func compareOne(r *rand.Rand, t *testing.T) {
	a := r.Int63()
	b := r.Int63()
	c := r.Int63()
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
