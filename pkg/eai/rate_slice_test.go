package eai

import (
	"fmt"
	"testing"

	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/stretchr/testify/require"
)

func duFrom(idx int) math.Duration {
	if idx < 0 {
		return DefaultUnlockedEAI[0].From - math.Duration(30*math.Day)
	}
	if idx >= len(DefaultUnlockedEAI) {
		idx = len(DefaultUnlockedEAI) - 1
	}
	return DefaultUnlockedEAI[idx].From + 1
}

func duTo(idx int) math.Duration {
	if idx < 0 {
		return DefaultUnlockedEAI[0].From - 1
	}
	if idx >= len(DefaultUnlockedEAI) {
		idx = len(DefaultUnlockedEAI) - 1
	}
	if idx == len(DefaultUnlockedEAI)-1 {
		return DefaultUnlockedEAI[idx].From + math.Duration(30*math.Day)
	}
	return DefaultUnlockedEAI[idx+1].From - 1
}

func TestRateSliceWithinRatePeriodReturnsSingleElement(t *testing.T) {
	for i := -1; i < len(DefaultUnlockedEAI); i++ {
		fD := duFrom(i)
		tD := duTo(i)
		name := fmt.Sprintf("from %s to %s", fD, tD)
		t.Run(name, func(t *testing.T) {
			rs := DefaultUnlockedEAI.Slice(fD, tD, 0)
			require.Equal(t, 1, len(rs))
			require.Equal(t, DefaultUnlockedEAI.RateAt(fD), DefaultUnlockedEAI.RateAt(tD))
			require.Equal(t, DefaultUnlockedEAI.RateAt(fD), rs[0].Rate)
			require.Equal(t, tD-fD, rs[0].Duration)
		})
	}
}

func TestRateSliceSpanningTwoPeriodsReturnsTwoElements(t *testing.T) {
	for i := -1; i < len(DefaultUnlockedEAI)-1; i++ {
		fD := duFrom(i)
		tD := duTo(i + 1)
		name := fmt.Sprintf("from %s to %s", fD, tD)
		t.Run(name, func(t *testing.T) {
			rs := DefaultUnlockedEAI.Slice(fD, tD, 0)
			require.Equal(t, 2, len(rs))
			require.NotEqual(t, DefaultUnlockedEAI.RateAt(fD), DefaultUnlockedEAI.RateAt(tD))
			require.Equal(t, DefaultUnlockedEAI.RateAt(fD), rs[0].Rate)
			require.Equal(t, DefaultUnlockedEAI[i+1].From-fD, rs[0].Duration)
			require.Equal(t, DefaultUnlockedEAI.RateAt(tD), rs[1].Rate)
			require.Equal(t, tD-DefaultUnlockedEAI[i+1].From, rs[1].Duration)
		})
	}
}

func TestRateSliceSpanningThreePeriodsReturnsThreeElements(t *testing.T) {
	for i := -1; i < len(DefaultUnlockedEAI)-2; i++ {
		fD := duFrom(i)
		tD := duTo(i + 2)
		name := fmt.Sprintf("from %s to %s", fD, tD)
		t.Run(name, func(t *testing.T) {
			rs := DefaultUnlockedEAI.Slice(fD, tD, 0)
			require.Equal(t, 3, len(rs))
			require.NotEqual(t, DefaultUnlockedEAI.RateAt(fD), DefaultUnlockedEAI.RateAt(tD))
			require.Equal(t, DefaultUnlockedEAI.RateAt(fD), rs[0].Rate)
			require.Equal(t, DefaultUnlockedEAI[i+1].From-fD, rs[0].Duration)
			require.Equal(t, DefaultUnlockedEAI[i+1].Rate, rs[1].Rate)
			require.Equal(t, DefaultUnlockedEAI[i+2].From-DefaultUnlockedEAI[i+1].From, rs[1].Duration)
			require.Equal(t, DefaultUnlockedEAI.RateAt(tD), rs[2].Rate)
			require.Equal(t, tD-DefaultUnlockedEAI[i+2].From, rs[2].Duration)
		})
	}
}
