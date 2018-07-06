package types

import "errors"

// ErrorOverflow is returned when a math operation would overflow a 64-bit value
var ErrorOverflow = errors.New("overflow error")

// ErrorDivideByZero is returned when a math operation would cause division by zero
var ErrorDivideByZero = errors.New("divide by zero error")

// ErrorMath is returned when the result of a decimal math operation could not be converted
// back to a uint64
var ErrorMath = errors.New("overflow error")
