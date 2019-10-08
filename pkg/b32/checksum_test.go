package b32

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

func TestCheck(t *testing.T) {
	type args struct {
		b   []byte
		ckb []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"a", args{[]byte("this is a test"), []byte{111, 238}}, true},
		{"b", args{[]byte(""), []byte{29, 15}}, true},
		{"c", args{[]byte("this was a test"), []byte{165, 254}}, true},
		{"d", args{[]byte("this Is a test"), []byte{200, 18}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Check(tt.args.b, tt.args.ckb); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChecksum(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want []byte
	}{
		// TODO: Add test cases.
		{"a", []byte("this is a test"), []byte{111, 238}},
		{"b", []byte(""), []byte{29, 15}},
		{"c", []byte("this was a test"), []byte{165, 254}},
		{"d", []byte("this Is a test"), []byte{200, 18}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Checksum16(tt.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}
