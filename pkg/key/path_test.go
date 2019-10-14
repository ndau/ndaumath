package key

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

func Test_newPath(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    path
		wantErr bool
	}{
		{"root", "/", path{}, false},
		{"1 level", "/123", path{pathElement{123, false}}, false},
		{"1 level hardened", "/123'", path{pathElement{123, true}}, false},
		{"3 levels", "/123/4/567890", path{
			pathElement{123, false},
			pathElement{4, false},
			pathElement{567890, false},
		}, false},
		{"3 levels w/hardened", "/123'/4'/567890", path{
			pathElement{123, true},
			pathElement{4, true},
			pathElement{567890, false},
		}, false},
		{"bad path 1", "/foo", nil, true},
		{"bad path 2", "/'", nil, true},
		{"bad path 3", "/123/123749327234979", nil, true},
		{"bad path 4", "/foo//bar", nil, true},
		{"bad path 5", "//", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newPath(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("newPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
