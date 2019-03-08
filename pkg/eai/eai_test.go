package eai

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ericlagergren/decimal"
	dmath "github.com/ericlagergren/decimal/math"
	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/stretchr/testify/require"
)

// we accumulate errors faster than a proper bignum implementation does,
// and there's not much we can do about it. All we can do is test that
// our error values are relatively low
const epsilon = 1.0 / 1000000

type testLock struct {
	NoticePeriod math.Duration
	UnlocksOn    *math.Timestamp
	Rate         Rate
}

func newTestLock(period math.Duration, bonusRateTable RateTable) *testLock {
	return &testLock{
		NoticePeriod: period,
		Rate:         bonusRateTable.RateAt(period),
	}
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

func (l *testLock) GetBonusRate() Rate {
	if l != nil {
		return l.Rate
	}
	return Rate(0)
}

var _ Lock = (*testLock)(nil)

func TestEAIFactorUnlocked(t *testing.T) {
	// simple tests that the eai factor for unlocked accounts is e ** (rate * time)
	// there is no period in the rate table shorter than a month, so using a few days
	// should be fine.

	// use decimal library to compute the expected values to double-check
	// ourselves.
	//
	// time is set at a constant 1 day, because that's what's used in this test.
	expect := func(t *testing.T, rate Rate) uint64 {
		time := 1 * math.Day
		// get the applicable rate
		// compute time as fraction of year
		ti := decimal.WithContext(decimal.Context128)
		ti.SetUint64(uint64(time))
		ti.Quo(ti, decimal.New(1*math.Year, 0))

		// big representation of rate denominator
		rd := decimal.New(constants.RateDenominator, 0)

		e := decimal.WithContext(decimal.Context128)
		e.SetUint64(uint64(rate))
		// divide by rd to get the proper fractional rate
		e.Quo(e, rd)
		e.Mul(e, ti)
		dmath.Exp(e, e)
		// multiply by rd to get the integer style again
		e.Mul(e, rd)
		// truncate dust
		e.RoundToInt()

		// generating expect in this way raises condition flags
		t.Log("expect conditions:", e.Context.Conditions.Error())
		e.Context.Conditions = 0

		// compute and return the expected value
		v, ok := e.Uint64()
		require.True(t, ok)
		return v
	}

	// special case for 0 rate
	blockTime := math.Timestamp(DefaultUnlockedEAI[0].From - 1)
	lastEAICalc := blockTime.Sub(1 * math.Day)
	weightedAverageAge := math.Duration(2 * math.Day)
	t.Run("0 rate", func(t *testing.T) {
		zero, err := calculateEAIFactor(
			blockTime, lastEAICalc, weightedAverageAge, nil,
			DefaultUnlockedEAI,
		)
		require.NoError(t, err)
		require.InEpsilon(t, expect(t, 0), zero, epsilon)
	})

	// now test each particular rate
	for idx, rate := range DefaultUnlockedEAI {
		blockTime := math.Timestamp(rate.From + (2 * math.Day))
		weightedAverageAge = math.Duration(blockTime)
		lastEAICalc = blockTime.Sub(1 * math.Day)

		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Logf("block time: %s", blockTime.String())
			t.Logf("last eai:   %s", lastEAICalc.String())
			t.Logf("WAA:        %s", weightedAverageAge.String())
			t.Logf("rate:       %d", rate.Rate)

			factor, err := calculateEAIFactor(
				blockTime, lastEAICalc, weightedAverageAge, nil,
				DefaultUnlockedEAI,
			)
			require.NoError(t, err)
			expectValue := expect(t, rate.Rate)

			t.Logf("expect:     %d", expectValue)
			t.Logf("actual:     %d", factor)

			require.InEpsilon(t, expectValue, factor, epsilon)
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
	rate := DefaultUnlockedEAI[len(DefaultUnlockedEAI)-1].Rate

	// use decimal library to compute the expected values to double-check
	// ourselves.
	//
	// time is set at a constant 1 day, because that's what's used in this test.
	expect := func(lock Lock) uint64 {
		time := 1 * math.Day
		// get the applicable rate
		// compute time as fraction of year
		ti := decimal.WithContext(decimal.Context128)
		ti.SetUint64(uint64(time))
		ti.Quo(ti, decimal.New(1*math.Year, 0))

		// big representation of rate denominator
		rd := decimal.New(constants.RateDenominator, 0)

		e := decimal.WithContext(decimal.Context128)
		e.SetUint64(uint64(rate))
		e.Add(
			e,
			decimal.WithContext(decimal.Context128).SetUint64(uint64(lock.GetBonusRate())),
		)

		// divide by rd to get the proper fractional rate
		e.Quo(e, rd)
		e.Mul(e, ti)
		dmath.Exp(e, e)
		// multiply by rd to get the integer style again
		e.Mul(e, rd)
		// truncate dust
		e.RoundToInt()

		// generating expect in this way raises condition flags
		t.Log("expect conditions:", e.Context.Conditions.Error())
		e.Context.Conditions = 0

		// compute and return the expected value
		v, ok := e.Uint64()
		require.True(t, ok)
		return v
	}

	// special case for 0 lock rate
	lock := newTestLock(DefaultLockBonusEAI[0].From-math.Day, DefaultLockBonusEAI)
	t.Run("no lock bonus", func(t *testing.T) {
		factor, err := calculateEAIFactor(
			blockTime, lastEAICalc, weightedAverageAge, lock,
			DefaultUnlockedEAI,
		)
		require.NoError(t, err)
		require.InEpsilon(t, expect(lock), factor, epsilon)
	})

	// now test each particular lock bonus rate
	for idx, lockRate := range DefaultLockBonusEAI {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			lock = newTestLock(lockRate.From, DefaultLockBonusEAI)
			factor, err := calculateEAIFactor(
				blockTime, lastEAICalc, weightedAverageAge, lock,
				DefaultUnlockedEAI,
			)
			require.NoError(t, err)
			require.InEpsilon(t, expect(lock), factor, epsilon)
		})
	}
}

func TestEAIFactorSoundness(t *testing.T) {
	// note: all cases below have us certain constants:
	// - block time: 1 year
	// - waa: 123 days

	days34 := math.Duration(34 * math.Day)
	days90 := math.Duration(90 * math.Day)
	days165 := math.Duration(165 * math.Day)
	days180 := math.Duration(180 * math.Day)

	type ec struct {
		rate uint64
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
			calc := func(rate uint64, days uint64) {
				t.Logf("Period %d:", period)
				period++
				time.SetUint64(days * math.Day)
				time.Quo(time, decimal.New(1*math.Year, 0))
				t.Logf(" Duration: %s (%d days)", time, days)
				percent.SetUint64(rate)
				percent.Quo(
					percent,
					decimal.WithContext(decimal.Context128).SetUint64(100),
				)
				t.Logf(" Rate: %s", percent)
				percent.Mul(percent, time)
				dmath.Exp(percent, percent)
				expected.Mul(expected, percent)
				t.Logf(" Factor: %s", percent)
			}

			for _, ec := range scase.expectCalc {
				calc(ec.rate, ec.days)
			}
			t.Logf("Total factor:  %s", expected)

			// calculate the actual value
			blockTime := math.Timestamp(1 * math.Year)
			lastEAICalc := blockTime.Sub(scase.lastEAIOffset)
			weightedAverageAge := math.Duration(123 * math.Day)
			var lock *testLock
			if scase.lockPeriod != nil {
				lock = newTestLock(*scase.lockPeriod, DefaultLockBonusEAI)
				if scase.lockNotifyOffset != nil {
					uo := blockTime.Add(*scase.lockNotifyOffset)
					lock.UnlocksOn = &uo
				}
			}
			actual, err := calculateEAIFactor(
				blockTime,
				lastEAICalc, weightedAverageAge,
				lock,
				DefaultUnlockedEAI,
			)
			require.NoError(t, err)

			// log the actual factor
			actualF := decimal.WithContext(decimal.Context128)
			actualF.SetUint64(actual)
			actualF.Quo(actualF, decimal.New(constants.RateDenominator, 0))
			t.Logf("Actual factor: %s", actualF)

			// convert to same format as actual
			expected.Mul(expected, decimal.New(constants.RateDenominator, 0))
			expectedValue, ok := expected.Uint64()
			require.True(t, ok)

			require.InEpsilon(t, expectedValue, actual, epsilon)
		})
	}
}

// we need some custom cases to test situations with different assumptions
func TestSoundnessCustomDates(t *testing.T) {
	daysn35 := math.Duration(-35 * math.Day)
	daysn15 := math.Duration(-15 * math.Day)
	days90 := math.Duration(90 * math.Day)
	days365 := math.Duration(math.Year)

	type ec struct {
		rate uint64
		days uint64
	}
	type customDateCase struct {
		expectCalc         []ec
		lastEAIOffset      math.Duration
		lockPeriod         *math.Duration
		unlocksOnOffset    *math.Duration
		blockTime          math.Timestamp
		weightedAverageAge math.Duration
	}

	// note: all cases below have a fixed block time of 1 year and must be
	// constructed such that things make sense given that constraint.
	cases := []customDateCase{
		//  Case 1: What happens if an account is:
		//
		// - locked for 365 days
		// - notified immediately
		// - 400 days since last EAI update
		// - current actual weighted average age is 400 days
		//
		// In other words, can we correctly handle the case that a genesis
		// account is first processed after it has already unlocked?
		//
		// The span of effective average age we care about for the unlocked
		// portion runs from day 0 to day 400. Using the example table:
		//
		//  13%       x────────────────┐
		//  10%                        x───────x--
		//          _______________________________
		//  actual    0               365     400
		//  effect.  365              365     400
		//
		// Because the account was locked for 365 days, and 365 days has a bonus
		// rate of 3%, the actual rate used for that period should increase by
		// a constant rate of 3%. At the end of the lock period, the bonus expires,
		// returning the account to the basic unlocked rate for its age: 10%.
		// We thus get the following calculation to
		// compute the EAI multiplier:
		//
		//    e^(13% * 365 days)
		//  * e^(10% *  35 days)
		customDateCase{
			expectCalc: []ec{
				{13, 365},
				{10, 35},
			},
			lastEAIOffset:      400 * math.Day,
			lockPeriod:         &days365,
			unlocksOnOffset:    &daysn35,
			blockTime:          400 * math.Day,
			weightedAverageAge: 400 * math.Day,
		},
		//  Case 2: What happens if an account is:
		//
		// - created at day 240
		// - immediately locked for 90 days
		// - notified immediately
		// - block time is 345
		// - 125 days since last EAI update (update was on day 220)
		// - current actual weighted average age is 105 days
		//
		// In other words, can we correctly handle the case that an account
		// is created, locked, notified, and unlocked all in the interval
		// between creditEAI txs?
		//
		// The span of effective average age we care about for the unlocked
		// portion runs from day 0 to day 400. Using the example table:
		//
		//   5%       x─────────────────┐
		//   4%                         └─────x--
		//          _____________________________
		//  actual   240               330   345
		//  effect.   90                90   105
		//  month    (3)               (3)
		//
		// Because the account was locked for 90 days, and 90 days has a bonus
		// rate of 1%, the actual rate used for that period should increase by
		// a constant rate of 1%. At the end of the lock period, the bonus expires,
		// returning the account to the basic unlocked rate for its age: 4%.
		//
		// Because the account was notified immediately, its effective WAA doesn't
		// change through the duration of the notification period.
		//
		// We thus get the following calculation to
		// compute the EAI multiplier:
		//
		//    e^(5% * 90 days)
		//  * e^(4% * 15 days)
		customDateCase{
			expectCalc: []ec{
				{5, 90},
				{4, 15},
			},
			lastEAIOffset:      125 * math.Day,
			lockPeriod:         &days90,
			unlocksOnOffset:    &daysn15,
			blockTime:          345 * math.Day,
			weightedAverageAge: 105 * math.Day,
		},
		//  Case 3: What happens if an account is:
		//
		// - created at day 240
		// - immediately locked for 90 days
		// - notified immediately
		// - block time is 345
		// - 100 days since last EAI update (update was on day 245)
		// - current actual weighted average age is 105 days
		//
		// In other words, can we correctly handle the case that an account
		// is created, locked, notified, and unlocked all in the interval
		// between creditEAI txs?
		//
		// The span of effective average age we care about for the unlocked
		// portion runs from day 0 to day 400. Using the example table:
		//
		//   5%       ├─────x───────────┐
		//   4%                         └─────x--
		//          _____________________________
		//  actual   240   245         330   345
		//  effect.   90    5           90   105
		//  month    (3)               (3)
		//
		// Because the account was locked for 90 days, and 90 days has a bonus
		// rate of 1%, the actual rate used for that period should increase by
		// a constant rate of 1%. At the end of the lock period, the bonus expires,
		// returning the account to the basic unlocked rate for its age: 4%.
		//
		// Because the account was notified immediately, its effective WAA doesn't
		// change through the duration of the notification period.
		//
		// We thus get the following calculation to
		// compute the EAI multiplier:
		//
		//    e^(5% * 85 days)
		//  * e^(4% * 15 days)
		customDateCase{
			expectCalc: []ec{
				{5, 85},
				{4, 15},
			},
			lastEAIOffset:      100 * math.Day,
			lockPeriod:         &days90,
			unlocksOnOffset:    &daysn15,
			blockTime:          345 * math.Day,
			weightedAverageAge: 105 * math.Day,
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
			calc := func(rate uint64, days uint64) {
				t.Logf("Period %d:", period)
				period++
				time.SetUint64(days * math.Day)
				time.Quo(time, decimal.New(1*math.Year, 0))
				t.Logf(" Duration: %s (%d days)", time, days)
				percent.SetUint64(rate)
				percent.Quo(
					percent,
					decimal.WithContext(decimal.Context128).SetUint64(100),
				)
				t.Logf(" Rate: %s", percent)
				percent.Mul(percent, time)
				dmath.Exp(percent, percent)
				expected.Mul(expected, percent)
				t.Logf(" Factor: %s", percent)
			}

			for _, ec := range scase.expectCalc {
				calc(ec.rate, ec.days)
			}
			t.Logf("Total factor:  %s", expected)

			// calculate the actual value
			lastEAICalc := scase.blockTime.Sub(scase.lastEAIOffset)
			var lock *testLock
			if scase.lockPeriod != nil {
				lock = newTestLock(*scase.lockPeriod, DefaultLockBonusEAI)
				if scase.unlocksOnOffset != nil {
					uo := scase.blockTime.Add(*scase.unlocksOnOffset)
					if uo < scase.blockTime.Sub(scase.weightedAverageAge).Add(lock.NoticePeriod) {
						t.Fatal("malformed test case: lock older than waa")
					}
					lock.UnlocksOn = &uo
				}
			}
			actual, err := calculateEAIFactor(
				scase.blockTime,
				lastEAICalc, scase.weightedAverageAge,
				lock,
				DefaultUnlockedEAI,
			)
			require.NoError(t, err)

			// log the actual factor
			actualF := decimal.WithContext(decimal.Context128)
			actualF.SetUint64(actual)
			actualF.Quo(actualF, decimal.New(constants.RateDenominator, 0))
			t.Logf("Actual factor: %s", actualF)

			// convert to same format as actual
			expected.Mul(expected, decimal.New(constants.RateDenominator, 0))
			expectedValue, ok := expected.Uint64()
			require.True(t, ok)

			require.InEpsilon(t, expectedValue, actual, epsilon)
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
		newTestLock(90*math.Day, DefaultLockBonusEAI),
		DefaultUnlockedEAI,
	)
	require.NoError(t, err)

	require.InEpsilon(t, uint64(expected), uint64(actual), epsilon)
}

func TestCalculateEAIRate(t *testing.T) {
	type args struct {
		weightedAverageAge math.Duration
		lock               *testLock
		unlockedTable      RateTable
	}
	tests := []struct {
		name string
		args args
		want Rate
	}{
		{"zero", args{0, nil, DefaultUnlockedEAI}, 0},
		{"65 days unlocked", args{65 * math.Day, nil, DefaultUnlockedEAI}, RateFromPercent(3)},
		{"90 days unlocked", args{90 * math.Day, nil, DefaultUnlockedEAI}, RateFromPercent(4)},
		// lock bonus: 1%. effective WAA: 155d -> 5m -> 6%. Expect 7%.
		{"65 days locked 90", args{65 * math.Day, newTestLock(90*math.Day, DefaultLockBonusEAI), DefaultUnlockedEAI}, RateFromPercent(7)},
		{"90 days locked 90", args{90 * math.Day, newTestLock(90*math.Day, DefaultLockBonusEAI), DefaultUnlockedEAI}, RateFromPercent(8)},
		// lock bonus: 1%. Effective WAA: 90d -> 3m -> 4%. Expect 5%.
		{"0 days locked 90", args{0 * math.Day, newTestLock(90*math.Day, DefaultLockBonusEAI), DefaultUnlockedEAI}, RateFromPercent(5)},
		// lock bonus: 4%. Effective WAA: 1000d -> 2y -> 10%. Expect 14%.
		{"0 days locked 1000", args{0 * math.Day, newTestLock(1000*math.Day, DefaultLockBonusEAI), DefaultUnlockedEAI}, RateFromPercent(14)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// we never test a notified lock in this loop, so we don't need to bother with timestamps
			if got := CalculateEAIRate(tt.args.weightedAverageAge, tt.args.lock, tt.args.unlockedTable, 0); got != tt.want {
				t.Errorf("CalculateEAIRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

const datefmt = "1/2/06"

type realtest struct {
	name         string
	chainDate    math.Timestamp
	lockDuration math.Duration
	quantity     math.Ndau
	expected     math.Ndau
	ken          math.Ndau
}

func getTestRecords(filename string) (string, []realtest) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	recs, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	testdate := ""
	testdata := make([]realtest, 0, len(recs))
	fieldnumbers := map[string]int{}
	for i, r := range recs {
		rec := realtest{}
		switch {
		case i == 0:
			testdate = r[1]
		case i == 4:
			for j, s := range r {
				fieldnumbers[s] = j
			}
		case i > 4:
			// if last column is blank it's not good data
			if r[0] == "" || r[len(r)-1] == "" {
				continue
			}
			for f, j := range fieldnumbers {
				switch f {
				case "chain date":
					startdate, err := time.Parse(datefmt, r[j])
					if err != nil {
						panic(err)
					}
					rec.chainDate, _ = math.TimestampFrom(startdate)
				case "ndau amount in":
					q, _ := strconv.ParseFloat(r[j], 64)
					rec.quantity = math.Ndau(q * constants.QuantaPerUnit)
				case "address ID":
					rec.name = r[j]
				case "days":
					d, _ := strconv.Atoi(r[j])
					rec.lockDuration = math.Duration(d * math.Day)
				case "actual EAI earned":
					q, _ := strconv.ParseFloat(r[j], 64)
					rec.expected = math.Ndau(q * constants.QuantaPerUnit)
				case "Simple EAI":
					q, _ := strconv.ParseFloat(r[j], 64)
					rec.ken = math.Ndau(q * constants.QuantaPerUnit)
				}
			}
			testdata = append(testdata, rec)
		}
	}
	return testdate, testdata
}

func withinEpsilon(x, y, epsilon math.Ndau) bool {
	if x-y > 0 {
		return (x-y < epsilon)
	}
	return (y-x < epsilon)
}

func TestCalculateRealWorld(t *testing.T) {
	// This test reads a CSV file of test data and checks that
	// the calculations the blockchain does match the same
	// calculations in the spreadsheet that generated the CSV.
	// This is all the data in the genesis block.
	testDate, tests := getTestRecords("output_2-11-19.csv")
	enddate, err := time.Parse(datefmt, testDate)
	if err != nil {
		t.Errorf("%s", err)
	}
	endts, _ := math.TimestampFrom(enddate)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			waa := endts.Since(tt.chainDate)
			got, err := Calculate(
				tt.quantity,
				endts, tt.chainDate, waa,
				newTestLock(tt.lockDuration, DefaultLockBonusEAI),
				DefaultUnlockedEAI,
			)
			if err != nil {
				t.Errorf("Calculate had a problem: %s", err)
			}
			if !withinEpsilon(got, tt.expected, math.Ndau(20)) {
				t.Errorf("Calculate() = %v, want Ed:%v Ken:%v", got, tt.expected, tt.ken)
			}
		})
	}
}
