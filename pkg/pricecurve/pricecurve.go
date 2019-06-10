package pricecurve

import (
	"math"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/oneiro-ndev/ndaumath/pkg/signed"
	"github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/pkg/errors"
)

//go:generate msgp

// A Nanocent is one billionth of one hundredth of one USD.
//
// It is fundamentally an integer and is computed using only integer math, for
// perfect determinism.
type Nanocent int64

const (
	dollar       = 100000000000 // 10^11 nanocents
	phaseBlocks  = 10000        // number of blocks in a phase
	SaleBlockQty = 1000         // number of ndau in a block
)

// ApproxPriceAtUnit returns the price of the next ndau in USD given the number
// already sold
func ApproxPriceAtUnit(nunitsSold types.Ndau) float64 {
	ndauSold := float64(nunitsSold / constants.QuantaPerUnit)
	saleBlock := ndauSold / SaleBlockQty

	if saleBlock <= phaseBlocks*1 {
		// price in phase 1 has 14 doublings, from a starting point of $1 to a
		// finishing price of $16384 at the 10-millionth unit
		var price1 = math.Pow(2.0, saleBlock*14/10000)
		return price1
	}

	// NOTE: this function replaces the elaborate spreadsheet model for phase 2
	// with a cubic approximation function that was developed from a curve fit
	// of a few of the key points on the phase 2 and phase 3 data. It is off by
	// a little bit from the initially proposed curve but it's vastly easier to
	// calculate. The difference is a little bit high early in phase 2 (at
	// worst, 13% high) and drifts to about 5% low late in phase 2. It's
	// generally slightly high in phase 3, peaking at 8%, but that's probably a
	// good thing as it makes the curve more s-like.
	//
	// Note that phase 1 is exactly as originally proposed and the slope at
	// entry of phase 2 is deliberately smooth.
	if saleBlock < phaseBlocks*3 {
		// determined by a cubic curvefit for phase 2 and 3
		// y = -41633 - 8.286618*x + 0.00167424*x^2 - 2.654015e-8*x^3
		const d = -2.654015e-8
		const c = 0.00167424
		const b = -8.286618
		const a = -41633
		x := saleBlock

		price2 := d*math.Pow(x, 3) + c*math.Pow(x, 2) + b*x + a
		return price2
	}

	// after the end of phase 3 we don't sell any more ndau so just return the
	// final price
	return 500450.83
}

// ApproxUnitAtPrice does a binary search for the lowest multiple of 1000 units
// that exceeds the price
func ApproxUnitAtPrice(price float64) int {
	high := 30000
	low := 0
	guess := 15000
	for high-low > 1 {
		p := ApproxPriceAtUnit(types.Ndau(guess * 1000 * constants.QuantaPerUnit))
		if p >= price {
			high = guess
		}
		if p < price {
			low = guess
		}
		guess = int((high + low) / 2)
	}
	return guess * 1000
}

// ApproxTotalPriceFor returns the total price for a group of ndau given the
// amount to be purchased and the number already sold The numbers passed in are
// integer number of napu NOT ndau
func ApproxTotalPriceFor(numNdau, alreadySold types.Ndau) float64 {
	const numPerBlock = 1000 * constants.QuantaPerUnit
	var totalPrice float64
	for {
		var price = ApproxPriceAtUnit(alreadySold)
		var availableInThisBlock = alreadySold % numPerBlock
		if availableInThisBlock == 0 {
			availableInThisBlock = numPerBlock
		}

		// if what we're buying fits in the current block, just calculate the
		// total price and we're done
		if numNdau <= availableInThisBlock {
			totalPrice += price * float64(numNdau/constants.QuantaPerUnit)
			return totalPrice
		}

		// otherwise, buy the remainder of this block and loop
		numNdau -= availableInThisBlock
		alreadySold += availableInThisBlock
		totalPrice += price * float64(availableInThisBlock/constants.QuantaPerUnit)
	}
}

func pow2(n int) uint64 {
	if n == 0 {
		return 0
	}
	return uint64(1) << uint(n)
}

// price in phase 1 has 14 doublings, increasing every 1,000 ndau from a starting point
// of $1 to a finishing price of $16384 at the 9,999,001st unit.
//
// The ratio between successive blocks is constant: 1.00097097419
//
// To prevent excessive error, we pre-compute a table of doublings, and
// work from there. The 14 entries in this table are the prices of ndau when
// 2 ^ (2 ^ (N - 1) * 14 / 9999) have been sold, where N = 1 to 14.

func phase1(block uint64) (out Nanocent) {
	doublings := []Nanocent{
		100000000000, 100097097419, 100291575187, 100681665003, 101466402368,
		103054274072, 106304953285, 113117158227, 128079155778, 164201982670,
		269884708015, 729084792015, 5320807694887, 283384837710462,
	}

	if block <= 1 {
		return doublings[int(block)]
	}

	// find the appropriate doubling for this block to get the base price.
	// linearly search the list; it's faster than binary for lists of this size.
	var dblock int
	for dblock, out = range doublings {
		if block >= pow2(dblock) && block < pow2(dblock+1) {
			break
		}
	}

	// now out has our base number. From this point, we need to apply a
	// constant ratio, however many times are required by the difference
	// between the block and the dblock
	var nout int64
	var err error
	for i := uint64(0); i <= (block - pow2(dblock)); i++ {
		nout, err = signed.MulDiv(
			int64(out),
			1000970974193617,
			1000000000000000,
		)
		if err != nil {
			panic(err.Error())
		}
		out = Nanocent(nout)
	}
	return
}

func phase23(block int64) (out Nanocent, err error) {
	// determined by a cubic curvefit for phase 2 and 3
	// y = -41633 - 8.286618*x + 0.00167424*x^2 - 2.654015e-8*x^3
	const (
		a  = 41633
		b  = 8286618
		bD = 1000000
		c  = 167424
		cD = 100000000
		d  = 2654015
		dD = 10000000 // sqrt of the actual divisor, because we apply it twice
	)
	var iout int64

	// zero-order term
	iout = -a * dollar

	// first-order terms
	order1, err := signed.MulDiv(block, b, bD)
	if err != nil {
		return 0, errors.Wrap(err, "order1")
	}
	iout -= order1 * dollar

	// second order term
	order2, err := signed.MulDiv(block*block, c, cD)
	if err != nil {
		return 0, errors.Wrap(err, "order2")
	}
	iout += order2 * dollar

	// third order term
	// compute it over a few rounds to reduce the chance of overflow
	// note that dD is the s
	order3 := block * block
	order3, err = signed.MulDiv(order3, block, dD)
	if err != nil {
		return 0, errors.Wrap(err, "order3 phase 1")
	}
	order3, err = signed.MulDiv(order3, d, dD)
	if err != nil {
		return 0, errors.Wrap(err, "order3 phase 2")
	}

	iout -= order3 * dollar

	out = Nanocent(iout)
	return
}

// PriceAtUnit returns the price of the next ndau given the number already sold
func PriceAtUnit(nunitsSold types.Ndau) (Nanocent, error) {
	ndauSold := nunitsSold / constants.QuantaPerUnit
	block := uint64(ndauSold / SaleBlockQty)

	if block <= phaseBlocks*1 {
		return phase1(block), nil
	}

	if block < phaseBlocks*3 {
		return phase23(int64(block))
	}

	// after the end of phase 3 we don't sell any more ndau so just return the
	// final price
	return Nanocent(50045083 * (dollar / 100)), nil
}
