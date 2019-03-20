package unsigned

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLnInt(t *testing.T) {
	for i := uint(0); i < 64; i++ {
		x := uint64(1) << i
		y := x | (x >> 1)
		for _, z := range []uint64{x, y} {
			t.Run(fmt.Sprint(z), func(t *testing.T) {
				expect := int(math.Floor(math.Log(float64(z))))
				got := LnInt(z)
				require.Equal(t, expect, got)
			})
		}
	}
}
