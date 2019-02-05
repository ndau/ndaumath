package pricecurve

import (
	"math"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/oneiro-ndev/ndaumath/pkg/types"
)

// PriceAtUnit returns the price of the next ndau in dollars given the number already sold
func PriceAtUnit(nunitsSold types.Ndau) float64 {
	const phaseBlocks = 10000
	ndauSold := float64(nunitsSold / constants.QuantaPerUnit)
	saleBlock := ndauSold / (1000)

	if saleBlock <= phaseBlocks*1 {
		// price in phase 1 has 14 doublings, from a starting point of $1 to a finishing price
		// of $16384 at the 10-millionth unit
		var price1 = math.Pow(2.0, saleBlock*14/10000)
		return price1
	}

	// NOTE: this function replaces the elaborate spreadsheet model for phase 2 with a cubic approximation
	// function that was developed // from a curve fit of a few of the key points on the phase 2 and
	// phase 3 data. It is off by a little bit from the initially proposed curve but it's vastly easier to calculate.
	// The difference is a little bit high early in phase 2 (at worst, 13% high) and drifts to about
	// 5% low late in phase 2. It's generally slightly high in phase 3, peaking at 8%, but that's probably
	// a good thing as it makes the curve more s-like.
	// Note that phase 1 is exactly as originally proposed and the slope at entry of phase 2 is deliberately smooth.
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

	// after the end of phase 3 we don't sell any more ndau so just return the final price
	return 500450.83
}

// UnitAtPrice does a binary search for the lowest multiple of 1000 units that exceeds the price
func UnitAtPrice(price float64) int {
	high := 30000
	low := 0
	guess := 15000
	for high-low > 1 {
		p := PriceAtUnit(types.Ndau(guess * 1000 * constants.QuantaPerUnit))
		if p >= price {
			high = guess
		}
		if p < price {
			low = guess
		}
		guess = int((high + low) / 2)
		// console.log('H:', high, 'L:', low, 'G:', guess, 'p:', p, 'wanted: ', price);
	}
	return guess * 1000
}

// TotalPriceFor returns the total price for a group of ndau given the amount to be purchased and the number already sold
// The numbers passed in are integer number of napu NOT ndau
func TotalPriceFor(numNdau, alreadySold types.Ndau) float64 {
	const numPerBlock = 1000 * constants.QuantaPerUnit
	var totalPrice float64
	for {
		var price = PriceAtUnit(alreadySold)
		var availableInThisBlock = alreadySold % numPerBlock
		if availableInThisBlock == 0 {
			availableInThisBlock = numPerBlock
		}

		// if what we're buying fits in the current block, just calculate the total price and we're done
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
