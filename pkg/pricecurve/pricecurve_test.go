package pricecurve

import (
	"testing"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/oneiro-ndev/ndaumath/pkg/types"
)

func Test_PriceAtUnit(t *testing.T) {
	tests := []struct {
		name       string
		nunitsSold types.Ndau
		want       float64
	}{
		{"0", 0, 1.00},
		{"1", 1, 1.0000009704065236},
		{"1000", 1000, 1.0009708770490777},
		{"714000", 714000, 1.9994455591209304},
		{"10,000,000", 10000000, 16384},
		{"15,000,000", 15000000, 121198.72375},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PriceAtUnit(tt.nunitsSold * constants.QuantaPerUnit); got != tt.want {
				t.Errorf("PriceAtUnit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UnitAtPrice(t *testing.T) {
	tests := []struct {
		name  string
		price float64
		want  int
	}{
		{"1", 1.0, 0},
		{"2", 2.0, 714000},
		{"16.90", 16.90, 2913000},
		{"16384", 16384, 9999000},
		{"100000", 100000, 14100000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnitAtPrice(tt.price); got != tt.want {
				t.Errorf("UnitAtPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
