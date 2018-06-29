package eai

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
)

var onePercentRate = RateTable{
	RTRow{
		Rate: OnePercent,
		From: 0,
	},
}

func Test100Ndau1Percent1Year(t *testing.T) {
	// EAI for 100 ndau at one percent rate after one year is 1.
	eai := Calculate(
		math.Ndau(100*constants.QuantaPerUnit),
		math.Year,
		nil,
		onePercentRate,
		nil,
	)
	require.Equal(t, math.Ndau(1*constants.QuantaPerUnit), eai)
}
func Test400Ndau1Percent3Months(t *testing.T) {
	// EAI for 400 ndau at one percent rate after 3 months is 1.
	eai := Calculate(
		math.Ndau(400*constants.QuantaPerUnit),
		math.Year/4,
		nil,
		onePercentRate,
		nil,
	)
	require.Equal(t, math.Ndau(1*constants.QuantaPerUnit), eai)
}
