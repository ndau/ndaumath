package types

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
	"testing"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
)

func TestNdau_Add(t *testing.T) {
	tests := []struct {
		name    string
		n       Ndau
		other   Ndau
		want    Ndau
		wantErr bool
	}{
		{"a", 1, 1, 2, false},
		{"b", 1, -1, 0, false},
		{"c", 1, 100, 101, false},
		{"d", 123456, 654321, 777777, false},
		{"e", Ndau(int64(math.MaxInt64)), -1, Ndau(int64(math.MaxInt64 - 1)), false},
		{"f", Ndau(int64(math.MaxInt64)), 1, 0, true},
		{"g", Ndau(int64(math.MaxInt64 / 2)), Ndau(int64(math.MaxInt64 / 2)), Ndau(int64(math.MaxInt64) - 1), false},
		{"h", Ndau(int64(math.MinInt64)), 1, Ndau(int64(math.MinInt64 + 1)), false},
		{"i", Ndau(int64(math.MaxInt64)), Ndau(int64(math.MinInt64)), -1, false},
		{"j", Ndau(int64(math.MinInt64 / 2)), Ndau(int64(math.MinInt64 / 2)), Ndau(int64(math.MinInt64)), false},
		{"k", Ndau(int64(math.MinInt64)), -1, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Add(tt.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ndau.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				s := err.Error()
				if s != "overflow error" {
					t.Errorf("Error type was wrong, got %s, wanted overflow error", s)
				}
				return
			}
			if got != tt.want {
				t.Errorf("Ndau.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNdau_Sub(t *testing.T) {
	tests := []struct {
		name    string
		n       Ndau
		other   Ndau
		want    Ndau
		wantErr bool
	}{
		{"a", 1, -1, 2, false},
		{"b", 1, 1, 0, false},
		{"c", 1, 100, -99, false},
		{"d", 654321, 123456, 530865, false},
		{"e", Ndau(int64(math.MaxInt64)), 1, Ndau(int64(math.MaxInt64 - 1)), false},
		{"f", Ndau(int64(math.MaxInt64)), -1, 0, true},
		{"g", Ndau(int64(math.MaxInt64 / 2)), -Ndau(int64(math.MaxInt64 / 2)), Ndau(int64(math.MaxInt64) - 1), false},
		{"h", Ndau(int64(math.MinInt64)), -1, Ndau(int64(math.MinInt64 + 1)), false},
		{"i", Ndau(int64(math.MaxInt64)), Ndau(int64(math.MaxInt64)), 0, false},
		{"j", Ndau(int64(math.MinInt64 / 2)), -Ndau(int64(math.MinInt64 / 2)), Ndau(int64(math.MinInt64)), false},
		{"k", Ndau(int64(math.MinInt64)), 1, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Sub(tt.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ndau.Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got != tt.want {
				t.Errorf("Ndau.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNdau_Abs(t *testing.T) {
	tests := []struct {
		name string
		n    Ndau
		want Ndau
	}{
		{"a", 1, 1},
		{"b", 100, 100},
		{"c", -101, 101},
		{"d", Ndau(int64(math.MaxInt64)), Ndau(int64(math.MaxInt64))},
		// explicitly test for the abs(MinInt) case which returns MinInt again
		{"e", Ndau(int64(math.MinInt64)), Ndau(int64(math.MinInt64))},
		{"f", -1, 1},
		{"g", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Abs(); got != tt.want {
				t.Errorf("Ndau.Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNdau_Compare(t *testing.T) {
	tests := []struct {
		name string
		n    Ndau
		rhs  Ndau
		want int
	}{
		{"a", 1, 2, -1},
		{"b", 2, 1, 1},
		{"c", 2, 2, 0},
		{"d", -1, 0, -1},
		{"e", 0, -1, 1},
		{"f", 0, 0, 0},
		{"g", 1, math.MaxInt64, -1},
		{"h", math.MaxInt64, 1, 1},
		{"i", math.MaxInt64, math.MaxInt64, 0},
		{"j", 1, math.MinInt64, 1},
		{"k", math.MinInt64, 1, -1},
		{"l", math.MinInt64, math.MinInt64, 0},
		{"m", math.MaxInt64, math.MinInt64, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Compare(tt.rhs); got != tt.want {
				t.Errorf("Ndau.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNdau_String(t *testing.T) {
	tests := []struct {
		name string
		n    Ndau
		want string
	}{
		{"a", constants.QuantaPerUnit, "1"},
		{"b", constants.QuantaPerUnit * 1.5, "1.5"},
		{"c", constants.QuantaPerUnit / 5, "0.2"},
		{"d", 1, "0.00000001"},
		{"e", 17*constants.QuantaPerUnit + 1234, "17.00001234"},
		{"f", -17 * constants.QuantaPerUnit, "-17"},
		{"g", -17*constants.QuantaPerUnit - 1234, "-17.00001234"},
		{"h", 100, "0.000001"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("Ndau.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ndauize(n int) Ndau {
	return Ndau(n * constants.NapuPerNdau)
}

func TestParseNdau(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    Ndau
		wantErr bool
	}{
		{"pct", "1%", Ndau(0), true},
		{"1", "1", ndauize(1), false},
		{"2", "2", ndauize(2), false},
		{"1000", "1000", ndauize(1000), false},
		{"1l", "1.00000000", ndauize(1), false},
		{"2l", "2.00000000", ndauize(2), false},
		{"1000l", "1000.00000000", ndauize(1000), false},
		{"0.5l", "0.50000000", ndauize(1) / 2, false},
		{"0.001l", "0.00100000", ndauize(1) / 1000, false},
		{"1t", "1.0", ndauize(1), false},
		{"2t", "2.0", ndauize(2), false},
		{"1000t", "1000.0", ndauize(1000), false},
		{"0.5t", "0.5", ndauize(1) / 2, false},
		{"0.001t", "0.001", ndauize(1) / 1000, false},
		{"too much precision", "1.000000001", ndauize(0), true},
		{"bare leading decimal", ".1", ndauize(1) / 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNdau(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNdau() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseNdau() = %v, want %v", got, tt.want)
			}
		})
	}
}
