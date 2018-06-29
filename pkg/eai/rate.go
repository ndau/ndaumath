package eai

import (
	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
)

//go:generate msgp

// OnePercent is the Rate of one percent annual interest
const OnePercent = constants.QuantaPerUnit

// A Rate defines a rate of increase over time.
//
// EAI is not interest. Ndau never earns interest. However,
// for ease of explanation, we will borrow some terminology
// from interest rate calculations.
//
// Rates are expressed in terms of effective annual increase.
// EAI for 100 ndau at one percent rate after one year is 1.
// EAI for 400 ndau at one percent rate after 3 months is 1.
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
type Rate uint64

//msgp:tuple RTRow

// RTRow is a single row of a rate table
type RTRow struct {
	From math.Duration
	Rate Rate
}

// A RateTable defines a stepped sequence of EAI rates which apply
// at varying durations.
//
// It is a logic error if the elements of a RateTable are not sorted
// in increasing order by their From field.
type RateTable []RTRow

// RateAt returns the rate in a RateTable for a given duration
func (rt RateTable) RateAt(duration math.Duration) Rate {
	rate := Rate(0)
	for _, row := range rt {
		if duration < row.From {
			return rate
		}
		rate = row.Rate
	}
	return rate
}

var (
	// DefaultUnlockedEAI is the default base rate table for unlocked accounts
	//
	// The UnlockedEAI rate table is a system variable which is adjustable
	// whenever the BPC desires, but for testing purposes, we use this
	// approximation as a default.
	DefaultUnlockedEAI RateTable

	// DefaultLockBonusEAI is the bonus rate for locks of varying length
	//
	// The LockBonusEAI rate table is a system variable which is adjustable
	// whenever the BPC desires, but for testing purposes, we use this
	// approximation as a default.
	DefaultLockBonusEAI RateTable
)

func init() {
	for i := 2; i <= 10; i++ {
		DefaultUnlockedEAI = append(DefaultUnlockedEAI, RTRow{
			Rate: Rate(i * OnePercent),
			From: math.Duration(i * 30 * math.Day),
		})
	}

	maxLBMonths := 3 * 12
	for i := 1; i <= maxLBMonths; i++ {
		DefaultLockBonusEAI = append(DefaultLockBonusEAI, RTRow{
			From: math.Duration(i * 30 * math.Day),
			Rate: Rate(15 * OnePercent * i / maxLBMonths),
		})
	}
}
