package eai

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import math "github.com/oneiro-ndev/ndaumath/pkg/types"

// Lock is anything which acts like a lock struct.
//
// If we want to use a Lock struct literal, we have two options:
//   1. Implement it in this package and end up with a dependency on noms.
//   2. Implement it in the ndaunode.backing package and end up with a
//      dependency cycle.
//
// We started with option 1, but a noms dependency from this low in the
// stack is not great. Instead, let's define a Lock interface which we can
// implement elsewhere.
type Lock interface {
	GetNoticePeriod() math.Duration
	GetUnlocksOn() *math.Timestamp
	GetBonusRate() Rate
}
