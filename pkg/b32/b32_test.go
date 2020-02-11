package b32

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
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want string
	}{
		{"a", []byte{1, 2, 3, 4, 5}, "aebagbaf"},
		{"b", []byte{}, ""},
		{"c", []byte{0, 0, 0, 0, 0}, "aaaaaaaa"},
		{"d", []byte{99, 100, 21, 0, 0}, "npubkaaa"},
		{"e", []byte{99, 100, 21, 255, 255}, "npubm999"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.b); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    []byte
		wantErr bool
	}{
		{"a", "aebagbaf", []byte{1, 2, 3, 4, 5}, false},
		{"b", "", []byte{}, false},
		{"c", "aaaaaaaa", []byte{0, 0, 0, 0, 0}, false},
		{"d", "npubaaaa", []byte{99, 100, 16, 0, 0}, false},
		{"e", "npvt9999", []byte{99, 103, 31, 255, 255}, false},
		{"f", "tpubaaaa", []byte{139, 100, 16, 0, 0}, false},
		{"g", "tpvt9999", []byte{139, 103, 31, 255, 255}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
