package types

import (
	"testing"
	"time"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
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
