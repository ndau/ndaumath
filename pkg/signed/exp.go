package signed

import (
	"errors"
	"math"
)

// This file contains an implementation of e^x (the exp function) that works for fractions
// between 0 and 1 in a 64-bit fixed point world.
// This frees us from the use of big math and it is also literally 25 times faster than the
// big package and has no memory allocation.

// ExpFrac calculates e^x, where x is a fraction a/b between
// 0 and 1. We use a Taylor Series expansion of e^x that converges well in the target range.
// This expansion is
// x^0/0! + x^1/1! + x^2/2! ...
// We can collapse the first two terms for convenience to 1+x.
// In addition, we make use of the fact that (a/b)^2 = a^2/b^2 so we can use muldiv
// and we require that b <= maxint32, and that a < b.
// Basically, we compute (b + a/b + a^2/2b + a^3/6b) which is b times our desired result
// (so that we have the implied denominator of b).
//
// The return value is the numerator for the fraction; the denominator is still b.
// This fixed point calculation tends to produce values that are slightly off in the last digit
// (as compared to a floating point implementation) because of accumulated rounding errors.
// Therefore, what we do is scale the input fraction by multiplying both numerator and denominator by
// a scaling value and then divide by it again at the end.
// This means that the practical limit for b is maxint32 / 10, which is still larger than our
// napu multiplication factor of 100,000,000.
func ExpFrac(a, b int64) (int64, error) {
	rounder := int64(10)
	a *= rounder
	b *= rounder
	if b > math.MaxInt32 {
		return 0, errors.New("denominator too large")
	}
	if a > b || a < 0 || b < 0 {
		return 0, errors.New("fraction must be between 0 and 1")
	}
	// start the sum at 1 + x, which is b/b + a/b, and we only care about the
	// numerator, so it's just b+a
	sum := b + a

	// we accumulate a product by starting with the original numerator,
	// then multiplying it by the fraction using muldiv; we don't square
	// the denominator because there's an implied division by b in the result.
	// In other words, to square 1200/10000, we muldiv(1200, 1200, 10000) and get 144;
	// the result is 144/10000 which is correct.
	// This converges at a rate of approximately one decimal digit per loop.
	product := a
	fact := int64(1)
	var err error
	for i := int64(2); product != 0; i++ {
		product, err = MulDiv(product, a, b)
		if err != nil {
			return 0, err
		}
		fact *= i
		sum += product / fact
	}
	return (sum + rounder/2) / rounder, nil
}
