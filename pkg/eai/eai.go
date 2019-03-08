package eai

import (
	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/oneiro-ndev/ndaumath/pkg/unsigned"
	"github.com/pkg/errors"
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
	lock Lock,
	ageTable RateTable,
) (math.Ndau, error) {
	factor, err := calculateEAIFactor(
		blockTime,
		lastEAICalc, weightedAverageAge, lock,
		ageTable,
	)
	if err != nil {
		return 0, err
	}

	// subtract 1 from the factor: we want just the EAI, not the new balance
	// remember that the factor has an implied divisor of RateDivisor
	factor -= constants.RateDenominator
	eai, err := unsigned.MulDiv(uint64(balance), factor, constants.RateDenominator)
	if err != nil {
		return 0, err
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
func calculateEAIFactor(
	blockTime, lastEAICalc math.Timestamp,
	weightedAverageAge math.Duration,
	lock Lock,
	unlockedTable RateTable,
) (uint64, error) {
	if lock != nil && lock.GetUnlocksOn() != nil && *lock.GetUnlocksOn() < blockTime {
		// we need to treat this as two nested calls and return their product
		unlockTs := *lock.GetUnlocksOn()

		atUnlock, err := calculateEAIFactor(
			unlockTs, lastEAICalc,
			weightedAverageAge-blockTime.Since(unlockTs),
			lock,
			unlockedTable,
		)
		if err != nil {
			return 0, errors.Wrap(err, "calculating preUnlock")
		}

		postUnlock, err := calculateEAIFactor(
			blockTime, unlockTs,
			weightedAverageAge,
			nil,
			unlockedTable,
		)
		if err != nil {
			return 0, errors.Wrap(err, "calculating postUnlock")
		}

		factor, err := unsigned.MulDiv(
			atUnlock,
			postUnlock,
			constants.RateDenominator,
		)
		if err != nil {
			return factor, errors.Wrap(err, "calculating composite factor")
		}

		return factor, err
	}

	factor := uint64(constants.RateDenominator) // 1.0, effectively

	lastEAICalcAge := blockTime.Since(lastEAICalc)
	var offset math.Duration
	if lock != nil {
		offset = lock.GetNoticePeriod()
	}
	from := weightedAverageAge - lastEAICalcAge
	if from < 0 {
		// the WAA can be treated as the actual age of the account.
		// if the WAA is more recent than the lastEAICalcAge, from will be negative.
		// this isn't a useful position to take: from any particular account's
		// perspective, the previous state of the blockchain shouldn't matter
		// at all. Therefore, we set it to 0 to get the correct rate period.
		from = 0
	}
	var rateSlice RateSlice
	if lock != nil && lock.GetUnlocksOn() != nil {
		notify := lock.GetUnlocksOn().Sub(lock.GetNoticePeriod())
		freeze := blockTime.Since(notify)
		rateSlice = unlockedTable.SliceF(from, weightedAverageAge, offset, freeze)
	} else {
		rateSlice = unlockedTable.Slice(from, weightedAverageAge, offset)
	}
	for _, row := range rateSlice {
		// fmt.Printf("%s @ %s\n", row.Duration.String(), row.Rate.String())

		// factor = e ^ (rate * time)
		// however, we're operating on rational numbers with implied
		// divisors:
		//              (      rate        duration)       (   rate * duration    )
		// factor = e ^ (--------------- * --------) = e ^ (----------------------)
		//              (RateDenominator     Year  )       (RateDenominator * Year)
		//
		// Computed naively, given that Year = (1_000_000 * 60 * 60 * 24 * 365),
		// any value of RateDenominator > 584_942 will cause the denominator product
		// to overflow uint64.
		//
		// At the same time, we want a big RateDenominator to minimize precision loss,
		// so we calculate in two stages:
		//
		//              ( (rate * duration)                   )
		// factor = e ^ ( ----------------- / RateDenominator )
		//              (        Year                         )

		effectiveRate := row.Rate
		if lock != nil {
			effectiveRate += lock.GetBonusRate()
		}
		divisor, err := unsigned.MulDiv(uint64(effectiveRate), uint64(row.Duration), math.Year)
		if err != nil {
			return 0, err
		}
		rowFactor, err := unsigned.ExpFrac(divisor, constants.RateDenominator)
		if err != nil {
			return 0, err
		}
		factor, err = unsigned.MulDiv(factor, rowFactor, constants.RateDenominator)
		if err != nil {
			return 0, err
		}
	}

	return factor, nil
}

// CalculateEAIRate accepts a WAA, a lock, a rate table, and a calculation
// timestamp, and looks up the current EAI rate from that info. The rate is
// returned as a Rate: a newtype wrapping a uint64, with an implied denominator
// of constants.RateDenominator.
//
// The timestamp is necessary in order to determine whether the lock is still
// notified, or the notice period has expired.
func CalculateEAIRate(
	weightedAverageAge math.Duration,
	lock Lock,
	unlockedTable RateTable,
	at math.Timestamp,
) Rate {
	effectiveWAA := weightedAverageAge
	if lock != nil {
		if lock.GetUnlocksOn() == nil {
			effectiveWAA += lock.GetNoticePeriod()
		} else {
			uo := *lock.GetUnlocksOn()
			if uo < at {
				// notified, which means our effective WAA is frozen at the
				// WAA we'll have at the unlock time, which means we need to
				// add the time until then
				effectiveWAA += uo.Since(at)
			}
			// else we're past the unlock timestamp, so we're back on the normal
			// increase of WAA
		}
	}
	effectiveRate := unlockedTable.RateAt(effectiveWAA)
	if lock != nil {
		effectiveRate += lock.GetBonusRate()
	}
	return effectiveRate
}
