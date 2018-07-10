package ndauerr

import "errors"

// ErrOverflow is returned when a math operation would overflow a 64-bit value
var ErrOverflow = errors.New("overflow error")

// ErrDivideByZero is returned when a math operation would cause division by zero
var ErrDivideByZero = errors.New("divide by zero error")

// ErrMath is returned when the result of a decimal math operation could not be converted
// back to a uint64
var ErrMath = errors.New("overflow error")
