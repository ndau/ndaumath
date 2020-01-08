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
	"testing"

	"github.com/ericlagergren/decimal"
	dmath "github.com/ericlagergren/decimal/math"
)

// this is a test helper that uses big math to calculate the results our routine should
// be returning.
func bigexp(a, b int64) int64 {
	af := decimal.WithContext(decimal.Context128)
	af.SetUint64(uint64(a))
	bf := decimal.WithContext(decimal.Context128)
	bf.SetUint64(uint64(b))
	q := af.Quo(af, bf)
	e := dmath.Exp(q, q)
	e.Mul(e, bf)
	e.RoundToInt()
	r, _ := e.Int64()
	return r
}

func TestExpFrac(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"zero", args{0, 1}, false},
		{"zero thousandths", args{0, 1000}, false},
		{"ten hundredths", args{10, 100}, false},
		{"1/10 in napu", args{10000000, 100000000}, false},
		{"1% in napu", args{1000000, 100000000}, false},
		{"2% in napu", args{2000000, 100000000}, false},
		{"3% in napu", args{3000000, 100000000}, false},
		{"4% in napu", args{4000000, 100000000}, false},
		{"5% in napu", args{5000000, 100000000}, false},
		{"6% in napu", args{6000000, 100000000}, false},
		{"7% in napu", args{7000000, 100000000}, false},
		{"8% in napu", args{8000000, 100000000}, false},
		{"9% in napu", args{9000000, 100000000}, false},
		{"10% in napu", args{10000000, 100000000}, false},
		{"11% in napu", args{11000000, 100000000}, false},
		{"12% in napu", args{12000000, 100000000}, false},
		{"13% in napu", args{13000000, 100000000}, false},
		{"14% in napu", args{14000000, 100000000}, false},
		{"15% in napu", args{15000000, 100000000}, false},
		{"bad denom", args{150000000000, 1000000000000}, true},
		{"negative numerator", args{-15000000, 100000000}, true},
		{"negative denominator", args{15000000, -100000000}, true},
		{"a>b", args{150000000, 100000000}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpFrac(tt.args.a, tt.args.b)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ExpFrac() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			want := bigexp(tt.args.a, tt.args.b)
			if got != want {
				t.Errorf("ExpFrac() = %v, want %v", got, want)
			}
		})
	}
}

// this prevents optimization of the return value
var v int64

func BenchmarkExp(b *testing.B) {
	for n := 0; n < b.N; n++ {
		v, _ = ExpFrac(15000000, 100000000)
	}
}

func BenchmarkBigExp(b *testing.B) {
	for n := 0; n < b.N; n++ {
		v = bigexp(15000000, 100000000)
	}
}
