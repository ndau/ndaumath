package eai

import (
	"testing"

	"github.com/ericlagergren/decimal"
	dmath "github.com/ericlagergren/decimal/math"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestEAIFactorSoundness1(t *testing.T) {
	//  What happens if an account is:
	//
	// - locked for 90 days
	// - not notified
	// - 84 days since last EAI update
	// - current actual weighted average age is 123 days
	//
	// The span of effective average age we care about for the unlocked
	// portion runs from day 174 to day 213. Using the example table:
	//
	//  8%                             ┌────x...
	//  7%                     ┌───────┘
	//  6%             ┌───────┘
	//  5%      ──x────┘
	//          _______________________________
	//  actual    39   60      90     120  123
	//  effect.  129  150     180     210  213
	//  month    (4)  (5)     (6)     (7)
	//
	// Because the account was locked for 90 days, and 90 days has a bonus
	// rate of 1%, the actual rate used for that period should increase by
	// a constant rate of 1%. We thus get the following calculation to
	// compute the EAI multiplier:
	//
	//    e^(6% * 21 days)
	//  * e^(7% * 30 days)
	//  * e^(8% * 30 days)
	//  * e^(9% *  3 days)

	// calculate the expected value
	expected := decimal.WithContext(decimal.Context128)
	percent := decimal.WithContext(decimal.Context128)
	time := decimal.WithContext(decimal.Context128)

	expected.SetUint64(1)

	calc := func(period int, rate float64, days uint64) {
		time.SetUint64(days * math.Day)
		time.Quo(time, decimal.New(1*math.Year, 0))
		t.Logf("Duration period %d: %s", period, time)
		rfp := RateFromPercent(rate)
		percent.Copy(&rfp.Big)
		t.Logf("Rate period %d: %s", period, percent)
		percent.Mul(percent, time)
		dmath.Exp(percent, percent)
		expected.Mul(expected, percent)
		t.Logf("Factor period %d: %s", period, percent)
	}

	calc(0, 6, 21)
	calc(1, 7, 30)
	calc(2, 8, 30)
	calc(3, 9, 3)

	// calculate the actual value
	weightedAverageAge := math.Duration(123 * math.Day)
	blockTime := math.Timestamp(weightedAverageAge) // for simplicity
	lastEAICalc := blockTime.Sub(math.Duration(84 * math.Day))
	actual := calculateEAIFactor(
		blockTime,
		lastEAICalc, weightedAverageAge,
		&math.Lock{NoticePeriod: 90 * math.Day},
		DefaultUnlockedEAI, DefaultLockBonusEAI,
	)

	// simplify
	expected.Reduce()
	actual.Reduce()

	// we require equal contexts here so that if the test fails in the
	// subsequent line, we know that it's not a context mismatch, but a value
	require.Equal(t, expected.Context, actual.Context)
	require.Equal(t, expected, actual)
}

func TestEAIFactorSoundness2(t *testing.T) {
	//  What happens if an account is:
	//
	// - locked for 90 days
	// - notified to unlock 34 days from now
	// - 84 days since last EAI update
	// - current actual weighted average age is 123 days
	//
	// The difference from case 1 is that the rate freezes the moment the
	// unlock goes through, but time keeps passing and interest keeps
	// accumulating during the notice period.
	//
	// The span of effective average age we care about for the unlocked
	// portion runs from actual day 39 to actual day 123. The notify happens
	// on actual day 67. It expires on actual day 157. At that point, the
	// rate will drop back to the actual weighted average age.
	//
	// The effective period begins on day 129, and runs forward normally
	// until effective day 157. Effective time freezes at that point. On
	// actual day 157, the notice period ends and calculations resume using
	// the actual weighted average age.
	//
	// Dashed lines in the following graph indicate points in the future,
	// assuming no further transactions are issued.
	//
	//  6%              ┌─────|────────────────x-----┐
	//  5%      ──x─────┘     |                      └---
	//  4%                    |
	//         _____________________________________
	//  actual    39    60    67              123   157
	//  effect.  129   150   157..............157...157
	//  month    (4)   (5)                          (5)
	//
	// Because the account was locked for 90 days, and 90 days has a bonus
	// rate of 1%, the actual rate used during the lock and notification
	// periods should increase by a constant rate of 1%.
	// We thus get the following calculation to compute the EAI multiplier:
	//
	//    e^(6% * 21 days)
	//  * e^(7% * 63 days)
	//
	// The 63 days of the final term are simply the seven unnotified days
	// of the rate period plus the 56 days notified to date.

	// calculate the expected value
	expected := decimal.WithContext(decimal.Context128)
	percent := decimal.WithContext(decimal.Context128)
	time := decimal.WithContext(decimal.Context128)

	expected.SetUint64(1)

	calc := func(period int, rate float64, days uint64) {
		t.Logf("Period %d:", period)
		time.SetUint64(days * math.Day)
		time.Quo(time, decimal.New(1*math.Year, 0))
		t.Logf(" Duration: %s", time)
		rfp := RateFromPercent(rate)
		percent.Copy(&rfp.Big)
		t.Logf(" Rate: %s", percent)
		percent.Mul(percent, time)
		dmath.Exp(percent, percent)
		expected.Mul(expected, percent)
		t.Logf(" Factor: %s", percent)
	}

	calc(0, 6, 21)
	calc(1, 7, 63)

	// calculate the actual value
	blockTime := math.Timestamp(1 * math.Year)
	unlocksOn := blockTime.Add(34 * math.Day)
	lastEAICalc := blockTime.Sub(84 * math.Day)
	weightedAverageAge := math.Duration(123 * math.Day)
	actual := calculateEAIFactor(
		blockTime,
		lastEAICalc, weightedAverageAge,
		&math.Lock{
			NoticePeriod: 90 * math.Day,
			UnlocksOn:    &unlocksOn,
		},
		DefaultUnlockedEAI, DefaultLockBonusEAI,
	)

	// simplify
	expected.Reduce()
	actual.Reduce()

	// we require equal contexts here so that if the test fails in the
	// subsequent line, we know that it's not a context mismatch, but a value
	require.Equal(t, expected.Context, actual.Context)
	require.Equal(t, expected, actual)
}
