package types

import (
	"fmt"
	"testing"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/stretchr/testify/require"
)

func TestDuration_UpdateWeightedAverageAge(t *testing.T) {
	// we derive the tests from some canonical data
	// computed in excel and validated by hand
	data := []struct {
		day      int
		transfer int
		balance  int
		waa      int
	}{
		{0, 0, 0, 0},       // dummy entry
		{0, 0, 0, 0},       // create an empty account
		{0, 100, 100, 0},   // give it a balance
		{30, 0, 100, 30},   // eai calculations; no transfer
		{30, 50, 150, 20},  // transfer in
		{40, -50, 100, 30}, // withdraw
		{60, 100, 200, 25}, // transfer in
		{80, -200, 0, 45},  // withdraw everything
		{100, 100, 100, 0}, // start again from 0
	}

	for index := range data {
		if index > 0 {
			sinceLastUpdate := Duration((data[index].day - data[index-1].day) * Day)
			transferQty := Ndau(data[index].transfer * constants.QuantaPerUnit)
			previousBalance := Ndau(data[index-1].balance * constants.QuantaPerUnit)
			waa := Duration(data[index-1].waa * Day)
			expectedWAA := Duration(data[index].waa * Day)

			t.Run(fmt.Sprintf("row %d", index), func(t *testing.T) {
				err := (&waa).UpdateWeightedAverageAge(sinceLastUpdate, transferQty, previousBalance)
				if err != nil {
					t.Errorf("Update weighted average age returned err: %s", err.Error())
				}
				if waa != expectedWAA {
					t.Errorf("WAA: %d; expected %d", waa, expectedWAA)
				}
			})
		}
	}
}

func TestParseDuration(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Duration
		wantErr bool
	}{
		{"<blank>", args{""}, Duration(0), false},
		{"t0s", args{"t0s"}, Duration(0), false},
		{"t1s", args{"t1s"}, Duration(1 * Second), false},
		{"1m", args{"1m"}, Duration(1 * Month), false},
		{"t1m", args{"t1m"}, Duration(1 * Minute), false},
		{"p1y2m3dt4h5m6s", args{"p1y2m3dt4h5m6s"}, Duration(36993906000000), false},
		{"P1Y2M3DT4H5M6S", args{"P1Y2M3DT4H5M6S"}, Duration(36993906000000), false},
		{"1y2m3dt4h5m6s7u", args{"1y2m3dt4h5m6s7u"}, Duration(36993906000007), false},
		{"1h", args{"1h"}, Duration(0), true},               // needs t
		{"100y", args{"100y"}, Duration(100 * Year), false}, // 3 digit year
		{"100m", args{"100m"}, Duration(0), true},           // 3 digit anything else
		{"100d", args{"100d"}, Duration(0), true},           // 3 digit anything else
		{"t100h", args{"t100h"}, Duration(0), true},         // 3 digit anything else
		{"t100m", args{"t100m"}, Duration(0), true},         // 3 digit anything else
		{"t100s", args{"t100s"}, Duration(0), true},         // 3 digit anything else
		{"t1u", args{"t1u"}, Duration(1), false},
		{"t1us", args{"t1us"}, Duration(1), false},
		{"t1μ", args{"t1μ"}, Duration(1), false},
		{"t1μs", args{"t1μs"}, Duration(1), false},
		{"t999999μ", args{"t999999μ"}, Duration(999999), false},
		{"t1000000μ", args{"t1000000μ"}, Duration(0), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

// MarshalText not tested because it's trivial
func TestDuration_UnmarshalText(t *testing.T) {
	d0 := Duration(0)
	tests := []struct {
		name    string
		t       *Duration
		text    string
		wantErr bool
	}{
		{"nil", nil, "", true},
		{"1234567", new(Duration), "1y2m3dt4h5m6s7μs", false},
		{"year", &d0, "1y", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.t.UnmarshalText([]byte(tt.text))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tt.t)
				remarshal := tt.t.String()
				require.Equal(t, tt.text, remarshal)
			}
		})
	}
}
