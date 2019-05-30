package types

import (
	"fmt"
	gomath "math"
	"regexp"
	"strconv"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/oneiro-ndev/ndaumath/pkg/signed"
	"github.com/pkg/errors"
)

//go:generate msgp -tests=0

// Ndau is a value that holds a single amount
// of ndau. Unlike an int64, it is prevented from overflowing.
type Ndau int64

// Add adds a value to an Ndau
// It may return an overflow error
func (n Ndau) Add(other Ndau) (Ndau, error) {
	t, err := signed.Add(int64(n), int64(other))
	return Ndau(t), err
}

// Sub subtracts, and may overflow
func (n Ndau) Sub(other Ndau) (Ndau, error) {
	t, err := signed.Sub(int64(n), int64(other))
	return Ndau(t), err
}

// Abs returns the absolute value without converting to float
// NOTE THAT THIS FAILS IN THE CASE WHERE n == MinInt64 (this
// value acts as much like -0 as it does MinInt). In particular,
// the value consists of only the negative (high) bit and the
// rest are zeros.
//
// As quantities on the blockchain can't be negative, we are going
// to ignore this case in favor of simplicity.
//
// In particular, this function a) can be inlined, and b) has no
// conditionals.
func (n Ndau) Abs() Ndau {
	y := n >> 63       // sign extended, so this is either -1 (0xFFF...) or 0
	return (n ^ y) - y // twos complement if it was negative
}

// Compare is the sorting operator; it returns -1 if n < rhs, 1 if n > rhs,
// and 0 if they are equal.
func (n Ndau) Compare(rhs Ndau) int {
	if n < rhs {
		return -1
	} else if n > rhs {
		return 1
	}
	return 0
}

// String returns the value of n formatted in a standard format, as if it is a
// decimal value of ndau. The full napu value is displayed, but trailing zeros
// are suppressed.
func (n Ndau) String() string {
	var sign int64 = 1
	if n < 0 {
		sign = -1
	}
	na := n.Abs()
	ndau := na / constants.NapuPerNdau
	napu := na % constants.NapuPerNdau
	if napu == 0 {
		return strconv.FormatInt(int64(sign*int64(ndau)), 10)
	}
	s := fmt.Sprintf("%d.%08d", sign*int64(ndau), napu)
	t := len(s)
	// trim off trailing zeros
	for ; s[t-1] == '0'; t-- {
	}
	return s[:t]
}

var (
	fracdigits int
	ndaure     *regexp.Regexp
)

func init() {
	// fracdigits: how many digits go behind the decimal?
	// computed here so that if constants.NapuPerNdau ever changes,
	// this stays automatically in sync
	fracdigits = int(gomath.Floor(gomath.Log10(constants.NapuPerNdau)))
	// ndaure: parse a string into whole (before the decimal) and frac (after the decimal)
	// strings, which can be used to regenerate the ndau
	ndaure = regexp.MustCompile(fmt.Sprintf(`^\s*(?P<whole>\d+)(\.(?P<frac>\d{1,%d}))?\s*$`, fracdigits))
}

// ParseNdau inverts n.String(): it converts a quantity of ndau expressed as
// a decimal number into a quantity of Ndau, without ever going through an
// intermediate floating-point step in which it may lose precision or behave
// nondeterministically.
func ParseNdau(s string) (Ndau, error) {
	match := ndaure.FindStringSubmatch(s)
	result := make(map[string]string)
	for i, name := range ndaure.SubexpNames() {
		if i != 0 && name != "" && i < len(match) {
			result[name] = match[i]
		}
	}

	wholes, ok := result["whole"]
	if !ok {
		return 0, errors.New("failed to parse ndau")
	}
	whole, err := strconv.ParseUint(wholes, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "parsing ndau")
	}

	out := Ndau(whole) * constants.NapuPerNdau

	fracs, ok := result["frac"]
	if ok {
		if len(fracs) > fracdigits {
			fracs = fracs[:fracdigits]
		} else if len(fracs) < fracdigits {
			iters := fracdigits - len(fracs)
			for i := 0; i < iters; i++ {
				fracs += "0"
			}
		}

		frac, err := strconv.ParseUint(fracs, 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, "parsing frac component")
		}
		out += Ndau(frac)
	}

	return out, nil
}
