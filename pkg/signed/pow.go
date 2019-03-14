package signed

// PowFrac generalizes ExpFrac to arbitrary bases.
//
// It computes `c`, where `c/denominator = (a/denominator)^(b/denominator)`.
//
// It uses the identity `a ^ b = c => e ^ (b ln a) = c`, generalized to operate
// on rational numbers without overflow.
//
// Note that internally, this function uses both ExpFrac and LnFrac, both of
// which approximate real numbers using integer rational pairs. The results
// are deterministic, but are likely to have large error margins with small
// denominators. Large denominators more closely approximate the true result.
func PowFrac(a, b, denominator int64) (int64, error) {
	lnA, err := LnFrac(a, denominator)
	if err != nil {
		return 0, err
	}
	exp, err := MulDiv(b, lnA, denominator)
	if err != nil {
		return 0, err
	}
	return ExpFrac(exp, denominator)
}
