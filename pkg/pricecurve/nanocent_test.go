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
	"reflect"
	"testing"
)

func Test_ParseDollars(t *testing.T) {
	tests := []struct {
		in      string
		want    Nanocent
		wantErr bool
	}{
		{"// TODO: Add test cases.", 0, true},
		{"1", 100000000000, false},
		{"-1", -100000000000, false},
		{"$1", 100000000000, false},
		{"-$1", -100000000000, false},
		{"1.5", 0, true},
		{"-1.5", 0, true},
		{"$1.5", 0, true},
		{"-$1.5", 0, true},
		{"1.50", 150000000000, false},
		{"-1.50", -150000000000, false},
		{"$1.50", 150000000000, false},
		{"-$1.50", -150000000000, false},
		{"0.00000000001", 1, false},
		{"-0.00000000001", -1, false},
		{"$0.00000000001", 1, false},
		{"-$0.00000000001", -1, false},
		{"0.000000000001", 0, true},
		{"-0.000000000001", 0, true},
		{"$0.000000000001", 0, true},
		{"-$0.000000000001", 0, true},
		{"0.00_000_000_001", 1, false},
		{"-0.00_000_000_001", -1, false},
		{"$0.00_000_000_001", 1, false},
		{"-$0.00_000_000_001", -1, false},
		{"0.00,000,000,001", 1, false},
		{"-0.00,000,000,001", -1, false},
		{"$0.00,000,000,001", 1, false},
		{"-$0.00,000,000,001", -1, false},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseDollars(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDollars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDollars() = %v, want %v", got, tt.want)
			}
		})
	}
}
