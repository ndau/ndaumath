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

// RateAt returns the rate in a RateTable for a given point
func (rt RateTable) RateAt(point math.Duration) Rate {
	rate := Rate(0)
	// the nature of rate tables is that we want the smallest rate
	// for which point >= row.From. The obvious way would be to iterate
	// in reverse, and return the first time that the point >= the row's
	// From value. However, reverse iteration is tedious in go, so we
	// take a different tack instead.
	//
	// We iterate forward, but keep a cache of the row rates, so that the
	// currently active rate always trails behind the table by one row.
	// This means that the first time the point < row.From, we can return
	// the active rate.
	for _, row := range rt {
		if point < row.From {
			return rate
		}
		rate = row.Rate
	}
	return rate
}

// RSRow is a single row of a rate slice.
type RSRow struct {
	Duration math.Duration
	Rate     Rate
}

// A RateSlice is derived from a RateTable, optimized for computation.
//
// Whereas a RateTable is meant to be easy for humans to understand, a
// RateSlice is more efficient for computation. It is a list of rates, and
// the actual duration over which each rate is active.
type RateSlice []RSRow

// Slice transforms a RateTable into a form more amenable for computation.
//
// Rates vary over time, and we want to efficiently compute the sum of interest
// at varying rates. Instead of repeatedly calling RateAt, it's more efficient
// to perform the calculation once to slice the affected data out of the
// RateTable.
func (rt RateTable) Slice(from, to math.Duration) RateSlice {
	if to <= from {
		// when actual duration is 0, it's fine to fake that the actual
		// rate is also 0
		return RateSlice{RSRow{}}
	}

	// the computation can't result in -2, so if after the loop
	// this remains, we know something went wrong
	const uninitialized = -2
	fromI := uninitialized
	toI := uninitialized

	for index, row := range rt {
		if fromI == uninitialized && from < row.From {
			fromI = index - 1
		}
		if toI == uninitialized && to < row.From {
			toI = index - 1
		}
		if fromI != uninitialized && toI != uninitialized {
			break
		}
	}
	// if either variable comes out of the loop wihtout being initialized,
	// the appropriate row index is the highest in the table
	if fromI == uninitialized {
		fromI = len(rt) - 1
	}
	if toI == uninitialized {
		toI = len(rt) - 1
	}

	rateFor := func(idx int) Rate {
		if idx == -1 {
			return Rate(0)
		}
		return rt[idx].Rate
	}

	// take care of the degenerate case, which is nicely simple
	if fromI == toI {
		return RateSlice{RSRow{Rate: rateFor(toI), Duration: to - from}}
	}
	// ok, the rest is relatively straightforward. We need special
	// handling for the first and last rate, because they have partial
	// durations; the rest are just copies from the rate table
	rs := make(RateSlice, toI-fromI+1)
	// it's safe to index rt[fromI+1] because if fromI were the max value,
	// then we would have already returned: fromI must equal toI
	rs[0] = RSRow{Rate: rateFor(fromI), Duration: rt[fromI+1].From - from}
	// it's safe to index rt[toI] because if toI were -1,
	// then we would have already returned: fromI must equal toI
	rs[len(rs)-1] = RSRow{Rate: rateFor(toI), Duration: to - rt[toI].From}

	// indexing rt[fromI+i+1] is safe because fromI+i+1 == toI at max i
	for i := 1; i < toI-fromI; i++ {
		rs[i] = RSRow{
			Rate:     rt[fromI+i].Rate,
			Duration: rt[fromI+i+1].From - rt[fromI+i].From,
		}
	}

	return rs
}

var (
	// DefaultUnlockedEAI is the default base rate table for unlocked accounts
	//
	// The UnlockedEAI rate table is a system variable which is adjustable
	// whenever the BPC desires, but for testing purposes, we use this
	// approximation as a default.
	//
	// Defaults drawn from https://tresor.it/p#0041o9iot7hm4kb5y707es7o/Oneiro%20Company%20Info/Whitepapers%20and%20Presentations/ndau%20Whitepaper%201.3%2020180425%20Final.pdf
	// page 15.
	DefaultUnlockedEAI RateTable

	// DefaultLockBonusEAI is the bonus rate for locks of varying length
	//
	// The LockBonusEAI rate table is a system variable which is adjustable
	// whenever the BPC desires, but for testing purposes, we use this
	// approximation as a default.
	//
	// Defaults drawn from https://tresor.it/p#0041o9iot7hm4kb5y707es7o/Oneiro%20Company%20Info/Whitepapers%20and%20Presentations/ndau%20Whitepaper%201.3%2020180425%20Final.pdf
	// page 15.
	DefaultLockBonusEAI RateTable
)

func init() {
	for i := 2; i < 10; i++ {
		DefaultUnlockedEAI = append(DefaultUnlockedEAI, RTRow{
			Rate: Rate((i + 1) * OnePercent),
			From: math.Duration(i * 30 * math.Day),
		})
	}

	DefaultLockBonusEAI = RateTable{
		RTRow{
			From: math.Duration(3 * 30 * math.Day),
			Rate: Rate(1 * OnePercent),
		},
		RTRow{
			From: math.Duration(6 * 30 * math.Day),
			Rate: Rate(2 * OnePercent),
		},
		RTRow{
			From: math.Duration(1 * math.Year),
			Rate: Rate(3 * OnePercent),
		},
		RTRow{
			From: math.Duration(2 * math.Year),
			Rate: Rate(4 * OnePercent),
		},
		RTRow{
			From: math.Duration(3 * math.Year),
			Rate: Rate(5 * OnePercent),
		},
	}
}
