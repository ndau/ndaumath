package eai

import (
	"math/big"

	math "github.com/oneiro-ndev/ndaumath/pkg/types"
)

// Calculate the EAI due for a given account
//
// EAI is not interest. Ndau never earns interest. However,
// for ease of explanation, we will borrow some terminology
// from interest rate calculations.
//
// Rates are expressed in terms of effective annual increase.
// EAI for 100 ndau at one percent rate after one year is 1.
// EAI for 400 ndau at one percent rate after a quarter year is 1.
//
// EAI accrues per the simple interest formula:
//   EAI = Prt
// Where:
//   P is the principal
//   r is the rate
//   t is the time
//
// The use of simple interest math instead of compound interest
// is an intentional incentive for users to choose nodes with a
// high voting power: these nodes will compute EAI more often,
// resulting in effective rates which more closely approximate
// continuous compounding.
func Calculate(
	balance math.Ndau,
	weightedAverageAge math.Duration,
	lock *math.Lock,
	ageTable, lockTable RateTable,
) math.Ndau {
	age := EffectiveAge(weightedAverageAge, lock)
	rate := EffectiveRate(age, lock, ageTable, lockTable)

	eai := new(big.Int)
	eai.SetUint64(uint64(balance)) // P : ndau
	operand := new(big.Int)
	operand.SetUint64(uint64(rate)) // R : annualized rate in percent
	eai.Mul(eai, operand)
	operand.SetUint64(uint64(age)) // T : time
	eai.Mul(eai, operand)

	// because the rate is annualized, we must now divide by 1 year
	operand.SetInt64(math.Year)
	eai.Div(eai, operand)
	// because the rate is expressed in percent, we must now divide
	// by 100%
	operand.SetUint64(100 * OnePercent)
	eai.Div(eai, operand)

	if !eai.IsUint64() {
		panic("eai calculation overflowed uint64")
	}

	return math.Ndau(eai.Uint64())
}

// EffectiveAge calculates the effective age for an account
func EffectiveAge(
	weightedAverageAge math.Duration,
	lock *math.Lock,
) math.Duration {
	if lock == nil {
		return weightedAverageAge
	}
	if lock.EffectiveWeightedAverageAge != nil {
		return *lock.EffectiveWeightedAverageAge
	}
	return weightedAverageAge + lock.NoticePeriod
}

// EffectiveRate calculates the effective EAI rate for an account
func EffectiveRate(
	effectiveAge math.Duration,
	lock *math.Lock,
	ageTable, lockTable RateTable,
) Rate {
	rate := ageTable.RateAt(effectiveAge)
	if lock != nil {
		rate += lockTable.RateAt(lock.NoticePeriod)
	}
	return rate
}
