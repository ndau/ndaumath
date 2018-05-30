package basics

// OverflowError is returned when the code is invalid and cannot be loaded or run
type OverflowError struct{}

func (e OverflowError) Error() string {
	return "overflow error"
}
