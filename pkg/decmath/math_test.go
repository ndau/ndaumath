package decmath

import (
	"math"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestMul(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"a", args{6, 7}, 42, false},
		{"b", args{600000000, 700000000}, 420000000000000000, false},
		{"c", args{600000000000, 700000000000}, 0, true},
		{"d", args{math.MaxUint32, math.MaxUint32}, math.MaxUint32 * math.MaxUint32, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mul(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"a", args{42, 7}, 6, false},
		{"b", args{420000000000000000, 700000000}, 600000000, false},
		{"c", args{600000000000, 0}, 0, true},
		{"d", args{math.MaxUint32 * math.MaxUint32, math.MaxUint32}, math.MaxUint32, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Div(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Div() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMod(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"a", args{42, 7}, 0, false},
		{"b", args{42, 5}, 2, false},
		{"c", args{420000000000000000, 700000001}, 100000001, false},
		{"d", args{12, 0}, 0, true},
		{"e", args{math.MaxUint32 * math.MaxUint32, math.MaxInt32}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mod(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDivMod(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		want1   uint64
		wantErr bool
	}{
		{"a", args{42, 7}, 6, 0, false},
		{"b", args{42, 5}, 8, 2, false},
		{"c", args{420000000000000000, 700000001}, 599999999, 100000001, false},
		{"d", args{12, 0}, 0, 0, true},
		{"e", args{math.MaxUint32 * math.MaxUint32, math.MaxInt32}, (math.MaxUint32 * math.MaxUint32) / math.MaxInt32, 1, false},
		{"f", args{42, 55}, 0, 42, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DivMod(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("DivMod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DivMod() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DivMod() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMulDiv(t *testing.T) {
	type args struct {
		v uint64
		n uint64
		d uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"a", args{80, 2, 5}, 32, false},
		{"b", args{82, 2, 5}, 32, false},
		{"c", args{83, 2, 5}, 33, false},
		{"d", args{80000000000, 2000000000, 5000000000}, 32000000000, false},
		{"e", args{80000000000, 2, 5}, 32000000000, false},
		{"f", args{80000000000, 2, 0}, 0, true},
		{"g", args{147, 155, 132}, 172, false},
		{"h", args{14717364050318377211, 15574702891736741942, 1324724618575407633}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MulDiv(tt.args.v, tt.args.n, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("MulDiv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MulDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func bigmuldiv(a, b, c uint64) uint64 {
	x := big.NewInt(0).SetUint64(a)
	y := big.NewInt(0).SetUint64(b)
	z := big.NewInt(0).SetUint64(c)
	r := big.NewInt(0)
	q := big.NewInt(0)
	r.Mul(x, y)
	q.Div(r, z)
	return q.Uint64()
}

func compareOne(r *rand.Rand, t *testing.T) {
	// make sure they're never negative
	a := r.Uint64() //& 0x7FFFFFFFFFFFFFFF
	b := r.Uint64() //& 0x7FFFFFFFFFFFFFFF
	c := r.Uint64() //& 0x7FFFFFFFFFFFFFFF
	if b > c {
		b, c = c, b
	}
	p, err := MulDiv(a, b, c)
	if err != nil {
		t.Error(err)
	}
	q := bigmuldiv(a, b, c)
	if p != q {
		t.Errorf("muldiv didn't match results from big.Int: (%v %v %v) %v != %v", a, b, c, p, q)
		t.Errorf("a=%x b=%x\n", a, b)
		a /= 10
		b /= 10
		c /= 10
		p, _ = MulDiv(a, b, c)
		q := bigmuldiv(a, b, c)
		t.Errorf("/10: (%v %v %v) %v == %v", a, b, c, p, q)
	}
}

func TestMulDivFuzz(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 10000; i++ {
		compareOne(r, t)
	}
}
