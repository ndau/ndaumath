package unsigned

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

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
