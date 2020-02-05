package eai

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"reflect"
	"testing"

	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestRate_String(t *testing.T) {
	tests := []struct {
		name string
		r    Rate
		want string
	}{
		{"1", RateFromPercent(1), "1%"},
		{"2", RateFromPercent(2), "2%"},
		{"1000", RateFromPercent(1000), "1000%"},
		{"0.5", RateFromPercent(1) / 2, "0.5%"},
		{"0.001", RateFromPercent(1) / 1000, "0.001%"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); got != tt.want {
				t.Errorf("Rate.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRate(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Rate
		wantErr bool
	}{
		{"notpct", args{"1"}, Rate(0), true},
		{"1", args{"1%"}, RateFromPercent(1), false},
		{"2", args{"2%"}, RateFromPercent(2), false},
		{"1000", args{"1000.0%"}, RateFromPercent(1000), false},
		{"1l", args{"1.0000000000%"}, RateFromPercent(1), false},
		{"2l", args{"2.0000000000%"}, RateFromPercent(2), false},
		{"1000l", args{"1000.0000000000%"}, RateFromPercent(1000), false},
		{"0.5l", args{"0.5000000000%"}, RateFromPercent(1) / 2, false},
		{"0.001l", args{"0.0010000000%"}, RateFromPercent(1) / 1000, false},
		{"1t", args{"1.0%"}, RateFromPercent(1), false},
		{"2t", args{"2.0%"}, RateFromPercent(2), false},
		{"1000t", args{"1000.0%"}, RateFromPercent(1000), false},
		{"0.5t", args{"0.5%"}, RateFromPercent(1) / 2, false},
		{"0.001t", args{"0.001%"}, RateFromPercent(1) / 1000, false},
		{"too much precision", args{"1.00000000001%"}, RateFromPercent(0), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRate(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRTRow_MarshalText(t *testing.T) {
	type fields struct {
		From math.Duration
		Rate Rate
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{"zero", fields{From: math.Duration(0), Rate: Rate(0)}, []byte("t0s:0%"), false},
		{"one", fields{From: math.Duration(1 * math.Day), Rate: RateFromPercent(1)}, []byte("1d:1%"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RTRow{
				From: tt.fields.From,
				Rate: tt.fields.Rate,
			}
			got, err := r.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("RTRow.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RTRow.MarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRTRow_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		want    RTRow
		args    args
		wantErr bool
	}{
		{"zero", RTRow{From: math.Duration(0), Rate: Rate(0)}, args{[]byte("t0s:0%")}, false},
		{"one", RTRow{From: math.Duration(1 * math.Day), Rate: RateFromPercent(1)}, args{[]byte("1d:1%")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RTRow{}
			if err := r.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("RTRow.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.Equal(t, r, &tt.want)
			if !reflect.DeepEqual(*r, tt.want) {
				t.Errorf("RTRow.UnmarshalText() = %v, want %v", *r, tt.want)
			}
		})
	}
}
