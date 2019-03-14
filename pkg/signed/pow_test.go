package signed

import (
	"testing"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
)

func TestPowFrac(t *testing.T) {
	const n = constants.NapuPerNdau

	type args struct {
		a           int64
		b           int64
		denominator int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"2^3", args{2 * n, 3 * n, n}, 8 * n, false},
		{"3^2", args{3 * n, 2 * n, n}, 9 * n, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PowFrac(tt.args.a, tt.args.b, tt.args.denominator)
			if (err != nil) != tt.wantErr {
				t.Errorf("PowFrac() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PowFrac() = %v, want %v", got, tt.want)
			}
		})
	}
}
