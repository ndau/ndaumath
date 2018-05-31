package types

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
