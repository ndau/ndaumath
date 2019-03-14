package signed

// LnFrac computes the natural logarithm of a rational number.
//
// It returns `x` where `ln (numerator / denominator) = x / denominator`
func LnFrac(numerator, denominator int64) (int64, error) {
	// TODO: implement this; this is just a stub right now
	//
	// There are two reasonable ways to calculate `x` such that `ln(n/d) = (x/d)`:
	//
	// - `x = LnInt(n^d) - LnInt(d^d)`
	// - `x = d*(LnInt(n) - LnInt(d))`
	//
	// These are mathematically equivalent, and they're both accurate. However,
	// there are problems actually implementing them: the first way overflows
	// spectacularly, the the second way is going to encounter way more rounding
	// error than I think we're really comfortable with.
	return 0, nil
}
