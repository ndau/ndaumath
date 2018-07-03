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
	//  8%                 ┌────x...
	//  7%         ┌───────┘
	//  6%  ──x────┘
	//      ________________________
	//       174  180     210  213
	//
	// Because the account was locked for 90 days, and 90 days has a bonus
	// rate of 1%, the actual rate used for that period should increase by
	// a constant rate of 1%. We thus get the following calculation to
	// compute the EAI multiplier:
	//
	//  e^(7% * 6 days) * e^(8% * 30 days) * e*(9% * 3 days)

	// calculate the expected value
	expected := decimal.WithContext(decimal.Context128)
	percent := decimal.WithContext(decimal.Context128)
	time := decimal.WithContext(decimal.Context128)

	expected.SetUint64(1)

	time.SetUint64(uint64(6 * math.Day))
	time.Quo(time, decimal.New(1*math.Year, 0))
	t.Logf("Duration period 0: %s", time)
	rfp := RateFromPercent(7.0)
	percent.Copy(&rfp.Big)
	t.Logf("Rate period 0: %s", percent)
	percent.Mul(percent, time)
	dmath.Exp(percent, percent)
	expected.Mul(expected, percent)
	t.Logf("Factor period 0: %s", percent)

	time.SetUint64(uint64(30 * math.Day))
	time.Quo(time, decimal.New(1*math.Year, 0))
	t.Logf("Duration period 1: %s", time)
	rfp = RateFromPercent(8.0)
	percent.Copy(&rfp.Big)
	t.Logf("Rate period 1: %s", percent)
	percent.Mul(percent, time)
	dmath.Exp(percent, percent)
	expected.Mul(expected, percent)
	t.Logf("Factor period 1: %s", percent)

	time.SetUint64(uint64(3 * math.Day))
	time.Quo(time, decimal.New(1*math.Year, 0))
	t.Logf("Duration period 2: %s", time)
	rfp = RateFromPercent(9.0)
	percent.Copy(&rfp.Big)
	t.Logf("Rate period 2: %s", percent)
	percent.Mul(percent, time)
	dmath.Exp(percent, percent)
	expected.Mul(expected, percent)
	t.Logf("Factor period 2: %s", percent)

	// calculate the actual value
	blockTime := math.Timestamp(1 * math.Year)
	lastEAICalc := blockTime.Sub(math.Duration(84 * math.Day))
	weightedAverageAge := math.Duration(123 * math.Day)
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
