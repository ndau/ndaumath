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
	lastEAICalcAge, weightedAverageAge math.Duration,
	lock *math.Lock,
	ageTable, lockTable RateTable,
) math.Ndau {
	factor := CalculateEAIFactor(
		lastEAICalcAge, weightedAverageAge,
		lock, ageTable,
	)

	if lock != nil {
		lockFactor := CalculateEAIFactor(
			lastEAICalcAge, weightedAverageAge,
			// lock is nil here because we don't want to offset the lock table
			nil, lockTable,
		)
		// e^a * e^b == e^(a+b)
		// we can calculate the locked and unlocked rates separately, and
		// multiply them together to get the composite rate.
		factor.Mul(factor, lockFactor)
	}

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

// CalculateEAIFactor calculates the EAI factor for a given table
//
// Factor = e ^ (rate * time)
//
// This calculates unconditionally without worrying about what kind of table
// was used.
func CalculateEAIFactor(
	lastEAICalcAge, weightedAverageAge math.Duration,
	lock *math.Lock,
	table RateTable,
) *decimal.Big {
	bal := decimal.WithContext(decimal.Context128)
	bal.SetUint64(1)

	qty := decimal.WithContext(decimal.Context128)
	offset := AgeOffset(lock)
	for _, row := range table.Slice(lastEAICalcAge+offset, weightedAverageAge+offset) {
		// new balance = balance * e ^ (rate * time)
		// first: what's the time? It's the fraction of a year used
		qty.SetUint64(uint64(row.Duration))
		qty.Quo(qty, decimal.New(math.Year, 0))

		// multiply by rate and exponentiate
		qty.Mul(qty, &row.Rate.Big)
		dmath.Exp(qty, qty)

		// multiply by balance
		bal.Mul(bal, qty)
	}

	return bal
}

// AgeOffset calculates the age offset for an account based on its lock
func AgeOffset(lock *math.Lock) math.Duration {
	if lock == nil {
		return math.Duration(0)
	}
	if lock.EffectiveWeightedAverageAge != nil {
		return *lock.EffectiveWeightedAverageAge
	}
	return lock.NoticePeriod
}
