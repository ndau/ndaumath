package types

import (
	"fmt"
	"testing"
	"time"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/stretchr/testify/require"
)

func TestTimestampFrom(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    Timestamp
		wantErr bool
	}{
		{"a", args{constants.Epoch}, 0, false},
		{"b", args{time.Date(2018, time.January, 18, 14, 21, 0, 0, time.UTC)},
			1000000 * (24*60*60*17 + 14*60*60 + 21*60), false},
		{"c", args{time.Date(2010, time.January, 18, 14, 21, 0, 0, time.UTC)}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimestampFrom(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimestampFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TimestampFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Timestamp
		wantErr bool
	}{
		{"a", args{"2018-01-01T00:00:00Z"}, 0, false},
		{"b", args{"2018-01-18T14:21:00Z"}, 1000000 * (24*60*60*17 + 14*60*60 + 21*60), false},
		{"c", args{"2010-01-01T00:00:00Z"}, 0, true},
		{"d", args{"BLAH"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimestamp(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_Compare(t *testing.T) {
	type args struct {
		o Timestamp
	}
	tests := []struct {
		name string
		t    Timestamp
		args args
		want int
	}{
		{"a", 10000000, args{0}, 1},
		{"b", 10000000, args{20000000}, -1},
		{"c", 10000000, args{10000000}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Compare(tt.args.o); got != tt.want {
				t.Errorf("Timestamp.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_Since(t *testing.T) {
	type args struct {
		o Timestamp
	}
	tests := []struct {
		name string
		t    Timestamp
		args args
		want Duration
	}{
		{"a", 10000000, args{0}, 10000000},
		{"b", 10000000, args{20000000}, -10000000},
		{"c", 10000000, args{10000000}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Since(tt.args.o); got != tt.want {
				t.Errorf("Timestamp.Since() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_Add(t *testing.T) {
	type args struct {
		d Duration
	}
	tests := []struct {
		name string
		t    Timestamp
		args args
		want Timestamp
	}{
		{"a", 10000000, args{0}, 10000000},
		{"b", 10000000, args{20000000}, 30000000},
		{"c", 0, args{10000000}, 10000000},
		{"d", 30000000, args{-20000000}, 10000000},
		{"e", 0, args{-10000000}, 0},
		{"f", constants.MaxTimestamp / 2, args{constants.MaxDuration / 2}, constants.MaxTimestamp - 1},
		{"g", constants.MaxTimestamp / 2, args{constants.MaxDuration}, constants.MaxTimestamp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Add(tt.args.d); got != tt.want {
				t.Errorf("Timestamp.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_Sub(t *testing.T) {
	type args struct {
		d Duration
	}
	tests := []struct {
		name string
		t    Timestamp
		args args
		want Timestamp
	}{
		{"a", 10000000, args{0}, 10000000},
		{"b", 20000000, args{10000000}, 10000000},
		{"c", 0, args{10000000}, 0},
		{"d", 30000000, args{-20000000}, 50000000},
		{"e", 0, args{-10000000}, 10000000},
		{"f", constants.MaxTimestamp / 2, args{-constants.MaxDuration / 2}, constants.MaxTimestamp - 1},
		{"g", constants.MaxTimestamp, args{-10}, constants.MaxTimestamp},
		{"h", constants.MaxTimestamp / 2, args{constants.MaxDuration / 2}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Sub(tt.args.d); got != tt.want {
				t.Errorf("Timestamp.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_String(t *testing.T) {
	tests := []struct {
		name string
		t    Timestamp
		want string
	}{
		{"a", 0, constants.EpochStart},
		{"b", 1000000 * (24*60*60*17 + 14*60*60 + 21*60), "2018-01-18T14:21:00Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("Timestamp.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_UpdateWeightedAverageAge(t *testing.T) {
	// we derive the tests from some canonical data
	// computed in excel and validated by hand
	data := []struct {
		day      int
		transfer int
		balance  int
		waa      int
	}{
		{0, 100, 100, 0},   // create an account
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
	good := []string{
		"",
		"t1s",
		"1m",
		"t1m",
		"p1y2m3dt4h5m6s",
		"P1Y2M3DT4H5M6S",
	}

	for _, good_example := range good {
		t.Log(good_example)
		_, err := ParseDuration(good_example)
		require.NoError(t, err)
	}
}
