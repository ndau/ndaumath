package pricecurve

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/stretchr/testify/require"
)

func Test_ApproxPriceAtUnit(t *testing.T) {
	tests := []struct {
		name       string
		nunitsSold types.Ndau
		want       float64
	}{
		{"0", 0, 1.00},
		{"1", 1, 1.000000970503574},
		{"1000", 1000, 1.0009709741936168},
		{"714214", 714214, 1.9999994454268752},
		{"714215", 714215, 2.000001386433485},
		{"9,999,000", 9999000, 16384},
		{"15,000,000", 15000000, 121198.72375},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApproxPriceAtUnit(tt.nunitsSold * constants.QuantaPerUnit); got != tt.want {
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
		{"16384", 16384, 9998000},
		{"100000", 100000, 14100000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApproxUnitAtPrice(tt.price); got != tt.want {
				t.Errorf("UnitAtPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTotalPriceFor(t *testing.T) {
	type args struct {
		numNdau     types.Ndau
		alreadySold types.Ndau
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"first ndau", args{100000000, 0}, 1},
		{"first block", args{100000000000, 0}, 1000},
		{"second block", args{100000000000, 100000000000}, 1000.9709741936168},
		{"ten blocks at start", args{1000000000000, 0}, 10043.807166082466},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApproxTotalPriceFor(tt.args.numNdau, tt.args.alreadySold); got != tt.want {
				t.Errorf("TotalPriceFor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_phase1_increases_monotonically(t *testing.T) {
	var prev Nanocent
	var curr Nanocent
	for i := 0; i < 10000; i++ {
		curr = phase1(uint64(i), true)
		if curr <= prev {
			t.Log("block:", i)
			t.Log("curr: ", curr)
			t.Log("prev: ", prev)
		}
		require.True(t, curr > prev)
		prev = curr
	}
}

func TestPhase1(t *testing.T) {
	var dataOut io.Writer

	if false { // probably disable this later sometime
		f, err := os.Create("phase1errors.csv")
		require.NoError(t, err)
		defer f.Close()
		dataOut = bufio.NewWriter(f)
		fmt.Fprintln(dataOut, "block,using floats,using ints,epsilon")
	}

	for block := uint64(0); block < 10000; block++ {
		sold := block * constants.QuantaPerUnit * SaleBlockQty
		apau := ApproxPriceAtUnit(types.Ndau(sold))
		pau := phase1(block, true)
		paud := float64(pau) / float64(Dollar)

		epsilon := (apau - paud) / apau

		if dataOut != nil {
			fmt.Fprintf(dataOut, "%d,%f,%f,%f\n", block, apau, paud, epsilon)
		}

		t.Run(fmt.Sprint(block), func(t *testing.T) {
			require.True(t, math.Abs(epsilon) < 0.000001, "abs epsilon must be < 0.0000001")
		})
	}
}

func Test_phase23IncreasesMonotonically(t *testing.T) {
	var prev Nanocent
	var curr Nanocent
	var err error
	for i := 10000; i < 30000; i++ {
		curr, err = phase23(int64(i))
		require.NoError(t, err)
		if curr <= prev {
			t.Log("block:", i)
			t.Log("curr: ", curr)
			t.Log("prev: ", prev)
		}
		require.True(t, curr > prev)
		prev = curr
	}
}

func TestPhase23(t *testing.T) {
	var dataOut io.Writer

	if false { // probably disable this later sometime
		f, err := os.Create("phase23errors.csv")
		require.NoError(t, err)
		defer f.Close()
		dataOut = bufio.NewWriter(f)
		fmt.Fprintln(dataOut, "block,using floats,using ints,epsilon")
	}

	for block := int64(10000); block < 30000; block++ {
		sold := block * constants.QuantaPerUnit * SaleBlockQty
		apau := ApproxPriceAtUnit(types.Ndau(sold))
		pau, err := phase23(block)
		require.NoError(t, err)
		paud := float64(pau) / float64(Dollar)

		epsilon := (apau - paud) / apau

		if dataOut != nil {
			fmt.Fprintf(dataOut, "%d,%f,%f,%f\n", block, apau, paud, epsilon)
		}

		t.Run(fmt.Sprint(block), func(t *testing.T) {
			t.Log("epsilon:", epsilon)
			require.True(t, math.Abs(epsilon) < 0.0002, "abs epsilon must be < 0.0002")
		})
	}
}
