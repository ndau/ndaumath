package eai

import (
	"fmt"

	"github.com/ericlagergren/decimal"
	dmath "github.com/ericlagergren/decimal/math"
	"github.com/oneiro-ndev/ndaumath/pkg/ndauerr"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
)

// Calculate the EAI due for a given account
//
// EAI is not interest. Ndau never earns interest. However,
// for ease of explanation, we will borrow some terminology
// from interest rate calculations.
//
// EAI is continuously compounded according to the formula
//
//   eai = balance * (e ^ (rate * time) - 1)
//
// Rates are expressed in percent per year.
//
// Thus, 100 ndau at 1 percent rate over 1 year yields 1.00501670 ndau EAI.
//
// The use of continuously compounded interest instead of simple interest
// aids in EAI predictability: using simple interest, an account which
// simply accumulates its EAI, whose node won frequently, would see a higher
// rate of actual return than an identical account whose node won infrequently.
// Continuously compounded interest avoids that issue: both accounts will
// see the same rate of return; the benefit of the one registered to the
// frequent node is that it sees the increase more often.
//
// It is a logic error if `lock != nil && lock.UnlocksOn < blockTime`:
// rates change at the unlock point, so this function must be called once
// at its unlock moment, and once again for the unlocked span from the
// unlock point until the next event.
func Calculate(
	balance math.Ndau,
	blockTime, lastEAICalc math.Timestamp,
	weightedAverageAge math.Duration,
	lock Lock,
	ageTable, lockTable RateTable,
) (math.Ndau, error) {
	factor, err := calculateEAIFactor(
		blockTime,
		lastEAICalc, weightedAverageAge, lock,
		ageTable, lockTable,
	)
	if err != nil {
		return 0, err
	}

	// subtract 1 from the factor: we want just the EAI, not the new balance
	qty := decimal.WithContext(decimal.Context128)
	qty.SetUint64(1)
	factor.Sub(factor, qty)

	// multiply by the ndau balance
	qty.SetUint64(uint64(balance))
	qty.Mul(qty, factor)

	// discard dust
	qty.RoundToInt()

	eai, couldConvert := qty.Uint64()
	if !couldConvert {
		return 0, ndauerr.ErrOverflow
	}
	return math.Ndau(eai), nil
}

// calculateEAIFactor calculates the EAI factor for a given table
//
// Factor = e ^ (rate * time)
//
// Let's diagram the variables in play in here:
//
//  Timestamps
//       │ (unnamed) effective account open
//       │   │         lastEAICalc
//       │   │           │   notify                blockTime    lock.UnlocksOn
// TIME ─┼───┼───────────┼─────┼─────────────────────┼────────────┼──>
//       │   │           │     ├────── freeze ───────┤            │
//       │   │           │     ├──────── lock.NoticePeriod ───────┤
//       │   │           │     └───────────── offset ─────────────┘
//       │   ├── from ───┴───── lastEAICalcAge ──────┤
//       │   └────── weightedAverageAge (to) ────────┘
//   Durations
//
// It is a logic error if lock.UnlocksOn < blockTime;
// in that case, this function will return nil.
func calculateEAIFactor(
	blockTime, lastEAICalc math.Timestamp,
	weightedAverageAge math.Duration,
	lock Lock,
	unlockedTable, lockBonusTable RateTable,
) (*decimal.Big, error) {
	factor := decimal.WithContext(decimal.Context128)
	factor.SetUint64(1)

	lastEAICalcAge := blockTime.Since(lastEAICalc)
	var offset math.Duration
	if lock != nil {
		offset = lock.GetNoticePeriod()
	}
	from := weightedAverageAge - lastEAICalcAge
	qty := decimal.WithContext(decimal.Context128)
	rate := decimal.WithContext(decimal.Context128)
	var rateSlice RateSlice
	if lock != nil && lock.GetUnlocksOn() != nil {
		if *lock.GetUnlocksOn() < blockTime {
			return nil, fmt.Errorf("*lock.UnlocksOn (%s) < blockTime (%s)",
				lock.GetUnlocksOn().String(), blockTime.String(),
			)
		}
		notify := lock.GetUnlocksOn().Sub(lock.GetNoticePeriod())
		freeze := blockTime.Since(notify)
		rateSlice = unlockedTable.SliceF(from, weightedAverageAge, offset, freeze)
	} else {
		rateSlice = unlockedTable.Slice(from, weightedAverageAge, offset)
	}
	for _, row := range rateSlice {
		// fmt.Printf("%s @ %s\n", row.Duration.String(), row.Rate.String())

		// new balance = balance * e ^ (rate * time)
		// first: what's the time? It's the fraction of a year used
		qty.SetUint64(uint64(row.Duration))
		qty.Quo(qty, decimal.New(math.Year, 0))

		// next: what's the actual rate? It's the slice rate plus
		// the lock bonus
		rate.Copy(&row.Rate.Big)
		if lock != nil {
			bonus := lockBonusTable.RateAt(lock.GetNoticePeriod())
			rate.Add(rate, &bonus.Big)
		}

		// multiply by rate and exponentiate
		qty.Mul(qty, rate)
		f, _ := qty.Float64()
		fmt.Println("Qty before exp: ", f)
		dmath.Exp(qty, qty)
		f, _ = qty.Float64()
		fmt.Println("Qty after exp:  ", f)

		// multiply into the current factor
		factor.Mul(factor, qty)
	}

	return factor, nil
}

// CalculateEAIRate accepts a WAA and a lock, plus rate tables,
// and looks up the current EAI rate from that info.
func CalculateEAIRate(
	weightedAverageAge math.Duration,
	lock Lock,
	unlockedTable, lockBonusTable RateTable,
) int64 {
	rate := unlockedTable.RateAt(weightedAverageAge)
	f, _ := rate.Float64()
	i, _ := rate.Int64()
	fmt.Println(rate, f, i)
	if lock != nil {
		bonus := lockBonusTable.RateAt(lock.GetNoticePeriod())
		rate.Add(&rate.Big, &bonus.Big)
	}
	r, _ := rate.Big.Int64()
	return r
}
