package signature

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_checksumWidth(t *testing.T) {
	type args struct {
		inputLen int
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{"0", args{0}, 4},
		{"1", args{1}, 3},
		{"2", args{2}, 7},
		{"3", args{3}, 6},
		{"4", args{4}, 5},
		{"5", args{5}, 4},
		{"6", args{6}, 3},
		{"7", args{7}, 7},
		{"8", args{8}, 6},
		{"9", args{9}, 5},
		{"10", args{10}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checksumWidth(tt.args.inputLen); got != tt.want {
				t.Errorf("checksumWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cksumN(t *testing.T) {
	type args struct {
		input []byte
		n     byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"one", args{[]byte("one"), checksumWidth(3)}, []byte{176, 192, 177, 150, 249, 175}},
		{"two", args{[]byte("two"), checksumWidth(3)}, []byte{14, 35, 190, 49, 58, 226}},
		{"three", args{[]byte("three"), checksumWidth(5)}, []byte{166, 243, 135, 67}},
		{"four", args{[]byte("four"), checksumWidth(4)}, []byte{4, 149, 35, 186, 113}},
		{"five", args{[]byte("five"), checksumWidth(4)}, []byte{88, 10, 161, 100, 95}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cksumN(tt.args.input, tt.args.n)
			if len(got) != int(tt.args.n) {
				t.Errorf("len(cksumN()) = %v, want %v", len(got), int(tt.args.n))
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cksumN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChecksum(t *testing.T) {
	tests := []struct {
		message string
	}{
		{"// TODO: Add test cases."},
		{"fox"},
		{"socks"},
		{"box"},
		{"knox"},
	}
	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			bmsg := []byte(tt.message)
			checksummed := AddChecksum(bmsg)
			gotMessage, gotChecksumOk := CheckChecksum(checksummed)
			if !bytes.Equal(bmsg, gotMessage) {
				t.Errorf("Checksum() gotMessage = %v, want %v", gotMessage, bmsg)
			}
			if !gotChecksumOk {
				t.Errorf("CheckChecksum() gotChecksumOk = %v", gotChecksumOk)
			}
		})
	}
}
