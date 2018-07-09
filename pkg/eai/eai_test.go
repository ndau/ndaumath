package eai

import (
	"fmt"
	"testing"

	"github.com/ericlagergren/decimal"
	dmath "github.com/ericlagergren/decimal/math"
	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/stretchr/testify/require"
)

type testLock struct {
	NoticePeriod math.Duration
	UnlocksOn    *math.Timestamp
}

func (l *testLock) GetNoticePeriod() math.Duration {
	if l != nil {
		return l.NoticePeriod
	}
	return 0
}

func (l *testLock) GetUnlocksOn() *math.Timestamp {
	if l != nil {
		return l.UnlocksOn
	}
	return nil
}

var _ Lock = (*testLock)(nil)

func TestEAIFactorUnlocked(t *testing.T) {
	// simple tests that the eai factor for unlocked accounts is e ** (rate * time)
	// there is no period in the rate table shorter than a month, so using a few days
	// should be fine.

	// special case for 0 rate
	blockTime := math.Timestamp(DefaultUnlockedEAI[0].From - 1)
	lastEAICalc := blockTime.Sub(1 * math.Day)
	weightedAverageAge := math.Duration(2 * math.Day)
	zero, err := calculateEAIFactor(
		blockTime, lastEAICalc, weightedAverageAge, nil,
		DefaultUnlockedEAI, DefaultLockBonusEAI,
	)
	t.Run("one", func(t *testing.T) {
		require.NoError(t, err)
		d := decimal.WithContext(decimal.Context128)
		d.SetUint64(1)
		require.Equal(t, d, zero.Reduce())
	})

	// now test each particular rate
	var expect *decimal.Big
	time := decimal.WithContext(decimal.Context128)
	time.SetUint64(1 * math.Day)
	time.Quo(time, decimal.New(1*math.Year, 0))
	for idx, rate := range DefaultUnlockedEAI {
		blockTime := math.Timestamp(rate.From + (2 * math.Day))
		weightedAverageAge = math.Duration(blockTime)
		lastEAICalc = blockTime.Sub(1 * math.Day)

		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Logf("block time: %s", blockTime.String())
			t.Logf("last eai:   %s", lastEAICalc.String())
			t.Logf("WAA:        %s", weightedAverageAge.String())
			t.Logf("rate:       %s", rate.Rate.String())
			expect = decimal.WithContext(decimal.Context128)
			expect.Copy(&rate.Rate.Big)
			expect.Mul(expect, time)
			dmath.Exp(expect, expect)

			// generating expect in this way raises condition flags
			t.Log("expect conditions:", expect.Context.Conditions.Error())
			expect.Context.Conditions = 0

			factor, err := calculateEAIFactor(
				blockTime, lastEAICalc, weightedAverageAge, nil,
				DefaultUnlockedEAI, DefaultLockBonusEAI,
			)
			require.NoError(t, err)
			require.Equal(t, decimal.Condition(0), factor.Context.Conditions)

			expect.Reduce()
			factor.Reduce()

			t.Logf("expect:     %s", expect)
			t.Logf("actual:     %s", factor)

			require.Equal(t, expect, factor)
		})
	}
}

func TestEAIFactorLocked(t *testing.T) {
	// simple tests that the eai factor for locked accounts is e ** (rate * time),
	// where rate is the unlocked rate of the lock period plus the lock duration,
	// plus the lock bonus rate.
	// there is no period in the lock rate table shorter than a month, so using a few days
	// should be fine.

	// for each of these cases we're going to use a WAA old enough to get the
	// max rate anyway, so we don't have to deal with the complexity of including
	// the lock duration in the base rate calculation.
	createdAt := math.Timestamp(0)
	blockTime := math.Timestamp(
		DefaultUnlockedEAI[len(DefaultUnlockedEAI)-1].From + (2 * math.Day),
	)
	lastEAICalc := blockTime.Sub(1 * math.Day)
	weightedAverageAge := blockTime.Since(createdAt)

	baseRate := DefaultUnlockedEAI[len(DefaultUnlockedEAI)-1].Rate
	expect := decimal.WithContext(decimal.Context128)
	oneDay := decimal.WithContext(decimal.Context128)
	oneDay.SetUint64(math.Day)
	oneDay.Quo(oneDay, decimal.New(math.Year, 0))

	// special case for 0 lock rate
	lock := testLock{NoticePeriod: DefaultLockBonusEAI[0].From - math.Day}
	t.Run("no lock bonus", func(t *testing.T) {
		expect.Copy(&baseRate.Big)
		expect.Mul(expect, oneDay)
		dmath.Exp(expect, expect)
		expect.Reduce()
		// generating expect in this way raises condition flags
		t.Log("expect conditions:", expect.Context.Conditions.Error())
		expect.Context.Conditions = 0

		factor, err := calculateEAIFactor(
			blockTime, lastEAICalc, weightedAverageAge, &lock,
			DefaultUnlockedEAI, DefaultLockBonusEAI,
		)
		require.NoError(t, err)
		factor.Reduce()
		require.Equal(t, expect, factor)
	})

	// now test each particular lock bonus rate
	for idx, lockRate := range DefaultLockBonusEAI {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			expect.Copy(&baseRate.Big)
			expect.Add(expect, &lockRate.Rate.Big)
			expect.Mul(expect, oneDay)
			dmath.Exp(expect, expect)
			expect.Reduce()
			// generating expect in this way raises condition flags
			t.Log("expect conditions:", expect.Context.Conditions.Error())
			expect.Context.Conditions = 0

			lock = testLock{NoticePeriod: lockRate.From}
			factor, err := calculateEAIFactor(
				blockTime, lastEAICalc, weightedAverageAge, &lock,
				DefaultUnlockedEAI, DefaultLockBonusEAI,
			)
			require.NoError(t, err)
			factor.Reduce()
			require.Equal(t, expect, factor)
		})
	}
}

func TestEAIFactorSoundness(t *testing.T) {
	days34 := math.Duration(34 * math.Day)
	days90 := math.Duration(90 * math.Day)
	days165 := math.Duration(165 * math.Day)
	days180 := math.Duration(180 * math.Day)

	type ec struct {
		rate float64
		days uint64
	}
	type soundnessCase struct {
		expectCalc       []ec
		lastEAIOffset    math.Duration
		lockPeriod       *math.Duration
		lockNotifyOffset *math.Duration
	}

	cases := []soundnessCase{
		//  Case 1: What happens if an account is:
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
		soundnessCase{
			expectCalc: []ec{
				{6, 21},
				{7, 30},
				{8, 30},
				{9, 3},
			},
			lastEAIOffset: 84 * math.Day,
			lockPeriod:    &days90,
		},
		// Case 2: What happens if an account is:
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
		//  6%              ┌─────|────────────────x-------
		//  5%      ──x─────┘     |
		//         _________________________________________
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
		soundnessCase{
			expectCalc: []ec{
				{6, 21},
				{7, 63},
			},
			lastEAIOffset:    84 * math.Day,
			lockPeriod:       &days90,
			lockNotifyOffset: &days34,
		},
		// Case 3: What happens if an account is:
		//
		// - locked for 180 days
		// - notified to unlock 165 days from now
		// - 84 days since last EAI update
		// - current actual weighted average age is 123 days
		//
		// The difference from case 2 is that there are three steps in the
		// function, and the bonus EAI is 2% instead of 1%.
		//
		// The span of effective average age we care about for the unlocked
		// portion runs from actual day 39 to actual day 123. The notify happens
		// on actual day 108. It expires on actual day 288. At that point, the
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
		// 10%                     ┌────────|────────x-------
		//  9%              ┌──────┘        |
		//  8%      ──x─────┘               |
		//         ___________________________________________
		//  actual    39    60     90      108      123   288
		//  effect.  219   240    270      288......288...288
		//  month    (7)   (8)    (9)                     (9)
		//
		// Because the account was locked for 180 days, and 180 days has a bonus
		// rate of 2%, the actual rate used during the lock and notification
		// periods should increase by a constant rate of 2%.
		// We thus get the following calculation to compute the EAI multiplier:
		//
		//    e^(10% * 21 days)
		//  * e^(11% * 30 days)
		//  * e^(12% * 33 days)
		//
		// The 33 days of the final term are simply the 18 unnotified days
		// of the rate period plus the 15 days notified to date.
		soundnessCase{
			expectCalc: []ec{
				{10, 21},
				{11, 30},
				{12, 33},
			},
			lastEAIOffset:    84 * math.Day,
			lockPeriod:       &days180,
			lockNotifyOffset: &days165,
		},
		// Case 4: What happens if an account is:
		//
		// - locked for 90 days
		// - notified to unlock 34 days from now
		// - 4 days since last EAI update
		// - current actual weighted average age is 123 days
		//
		// The difference from case 2 is that the we've calculated recently,
		// so the rate freeze happens before the last update.
		//
		// Dashed lines in the following graph indicate points in the future,
		// assuming no further transactions are issued.
		//
		//  6%              ┌─────|───────────x────x-------
		//  5%      ────────┘     |
		//         _________________________________________
		//  actual    39    60    67         119  123   157
		//  effect.  129   150   157.........157..157...157
		//  month    (4)   (5)                          (5)
		//
		// Because the account was locked for 90 days, and 90 days has a bonus
		// rate of 1%, the actual rate used during the lock and notification
		// periods should increase by a constant rate of 1%.
		// We thus get the following calculation to compute the EAI multiplier:
		//
		//    e^(7% * 4 days)
		soundnessCase{
			expectCalc: []ec{
				{7, 4},
			},
			lastEAIOffset:    4 * math.Day,
			lockPeriod:       &days90,
			lockNotifyOffset: &days34,
		},
		// Case 5: What happens if an account is:
		//
		// - unlocked
		// - 84 days since last EAI update
		// - current actual weighted average age is 123 days
		//
		// This differs from case 1 in that the account is not locked.
		//
		// The span of effective average age we care about for the unlocked
		// portion runs from day 39 to day 123. Using the example table:
		//
		//  5%                             ┌────x...
		//  4%                     ┌───────┘
		//  3%             ┌───────┘
		//  2%      ──x────┘
		//          _______________________________
		//  actual    39   60      90     120  123
		//  month    (1)  (2)     (3)     (4)
		//
		// Because the account is unlocked, there is no bonus EAI. Our calculation:
		//
		//    e^(2% * 21 days)
		//  * e^(3% * 30 days)
		//  * e^(4% * 30 days)
		//  * e^(5% *  3 days)
		soundnessCase{
			expectCalc: []ec{
				{2, 21},
				{3, 30},
				{4, 30},
				{5, 3},
			},
			lastEAIOffset: 84 * math.Day,
		},
	}

	for i, scase := range cases {
		name := fmt.Sprintf("case %d", i+1)
		t.Run(name, func(t *testing.T) {
			expected := decimal.WithContext(decimal.Context128)
			percent := decimal.WithContext(decimal.Context128)
			time := decimal.WithContext(decimal.Context128)

			expected.SetUint64(1)

			var period int
			calc := func(rate float64, days uint64) {
				t.Logf("Period %d:", period)
				period++
				time.SetUint64(days * math.Day)
				time.Quo(time, decimal.New(1*math.Year, 0))
				t.Logf(" Duration: %s (%d days)", time, days)
				rfp := RateFromPercent(rate)
				percent.Copy(&rfp.Big)
				t.Logf(" Rate: %s", percent)
				percent.Mul(percent, time)
				dmath.Exp(percent, percent)
				expected.Mul(expected, percent)
				t.Logf(" Factor: %s", percent)
			}

			for _, ec := range scase.expectCalc {
				calc(ec.rate, ec.days)
			}
			t.Logf("Total factor: %s", expected)

			// calculate the actual value
			blockTime := math.Timestamp(1 * math.Year)
			lastEAICalc := blockTime.Sub(scase.lastEAIOffset)
			weightedAverageAge := math.Duration(123 * math.Day)
			var lock *testLock
			if scase.lockPeriod != nil {
				lock = &testLock{NoticePeriod: *scase.lockPeriod}
				if scase.lockNotifyOffset != nil {
					uo := blockTime.Add(*scase.lockNotifyOffset)
					lock.UnlocksOn = &uo
				}
			}
			actual, err := calculateEAIFactor(
				blockTime,
				lastEAICalc, weightedAverageAge,
				lock,
				DefaultUnlockedEAI, DefaultLockBonusEAI,
			)
			require.NoError(t, err)

			// simplify
			expected.Reduce()
			actual.Reduce()

			// we require equal contexts here so that if the test fails in the
			// subsequent line, we know that it's not a context mismatch, but a value
			require.Equal(t, expected.Context, actual.Context)
			require.Equal(t, expected, actual)
		})
	}
}

func TestCalculate(t *testing.T) {
	// the meat of the math happens in the calculation of the EAI factor,
	// but it's also worth sanity-checking the actual public function which
	// takes and returns Ndau
	//
	// Let's take the example from case 1:
	// - locked for 90 days
	// - not notified
	// - 84 days since last EAI update
	// - current actual weighted average age is 123 days
	//
	// Let's say additionally that the account contained exactly 1 ndau after
	// its last EAI update. The expected value calculation is easy:
	//
	//      1 ndau
	//    =  100 000 000 napu (from constants)
	//    * 0.01 665 776 679 ... (from scenario 1 log output)
	//    =    1 665 777 napu, as the ndau spec requires rounding dust (towards nearest even on 5)
	expected := math.Ndau(1665777)

	weightedAverageAge := math.Duration(123 * math.Day)
	blockTime := math.Timestamp(weightedAverageAge) // for simplicity
	lastEAICalc := blockTime.Sub(math.Duration(84 * math.Day))
	actual, err := Calculate(
		1*constants.QuantaPerUnit,
		blockTime, lastEAICalc, weightedAverageAge,
		&testLock{
			NoticePeriod: 90 * math.Day,
		},
		DefaultUnlockedEAI, DefaultLockBonusEAI,
	)
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}
