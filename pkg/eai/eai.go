package eai

import (
	"github.com/ericlagergren/decimal"
	dmath "github.com/ericlagergren/decimal/math"
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
func Calculate(
	balance math.Ndau,
	blockTime, lastEAICalc math.Timestamp,
	weightedAverageAge math.Duration,
	lock *math.Lock,
	ageTable, lockTable RateTable,
) math.Ndau {
	factor := calculateEAIFactor(
		blockTime,
		lastEAICalc, weightedAverageAge, lock,
		ageTable, lockTable,
	)

	// subtract 1 from the factor: we want just the EAI, not the new balance
	qty := decimal.WithContext(decimal.Context128)
	qty.SetUint64(1)
	factor.Sub(factor, qty)

	// multiply by the ndau balance
	qty.SetUint64(uint64(balance))
	qty.Mul(qty, factor)

	// factor is now no longer the exponentiation factor: it's the rounding
	// increment
	factor.SetUint64(0)
	// discard dust
	dmath.Floor(qty, factor)

	eai, couldConvert := qty.Uint64()
	if !couldConvert {
		panic("Overflow in EAI calculation")
	}
	return math.Ndau(eai)
}

// calculateEAIFactor calculates the EAI factor for a given table
//
// Factor = e ^ (rate * time)
//
// This calculates unconditionally without worrying about what kind of table
// was used.
func calculateEAIFactor(
	blockTime, lastEAICalc math.Timestamp,
	weightedAverageAge math.Duration,
	lock *math.Lock,
	unlockedTable, lockBonusTable RateTable,
) *decimal.Big {
	factor := decimal.WithContext(decimal.Context128)
	factor.SetUint64(1)

	lastEAICalcAge := blockTime.Since(lastEAICalc)
	offset := ageOffset(lock, blockTime)
	qty := decimal.WithContext(decimal.Context128)
	rate := decimal.WithContext(decimal.Context128)
	for _, row := range unlockedTable.Slice(lastEAICalcAge+offset, weightedAverageAge+offset) {
		// new balance = balance * e ^ (rate * time)
		// first: what's the time? It's the fraction of a year used
		qty.SetUint64(uint64(row.Duration))
		qty.Quo(qty, decimal.New(math.Year, 0))

		// next: what's the actual rate? It's the slice rate plus
		// the lock bonus
		rate.Copy(&row.Rate.Big)
		if lock != nil {
			bonus := lockBonusTable.RateAt(lock.NoticePeriod)
			rate.Add(rate, &bonus.Big)
		}

		// multiply by rate and exponentiate
		qty.Mul(qty, rate)
		dmath.Exp(qty, qty)

		// multiply into the current factor
		factor.Mul(factor, qty)
	}

	return factor
}

// ageOffset calculates the age offset for an account based on its lock
func ageOffset(lock *math.Lock, blockTime math.Timestamp) math.Duration {
	if lock == nil {
		return math.Duration(0)
	}
	if lock.UnlocksOn != nil {
		// a notified lock has a offset computed such that the actual
		// weighted average age plus the offset is equal to the value
		// of that calculation at the moment of notification
		return lock.UnlocksOn.Since(blockTime)
	}
	return lock.NoticePeriod
}
