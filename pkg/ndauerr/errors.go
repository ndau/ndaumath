package ndauerr

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import "errors"

// ErrOverflow is returned when a math operation would overflow a 64-bit value
var ErrOverflow = errors.New("overflow error")

// ErrDivideByZero is returned when a math operation would cause division by zero
var ErrDivideByZero = errors.New("divide by zero error")

// ErrMath is returned when the result of a decimal math operation could not be converted
// back to a uint64
var ErrMath = errors.New("overflow error")
