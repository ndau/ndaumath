package types

import (
	"math"
	"testing"
)

func TestNdauQty_Add(t *testing.T) {
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
				t.Errorf("NdauQty.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got != tt.want {
				t.Errorf("NdauQty.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNdauQty_Sub(t *testing.T) {
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
				t.Errorf("NdauQty.Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got != tt.want {
				t.Errorf("NdauQty.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNdauQty_Abs(t *testing.T) {
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
				t.Errorf("NdauQty.Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}
