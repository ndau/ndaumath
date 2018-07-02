package eai

import "github.com/ericlagergren/decimal"

// shim to assist rate (de)serialization
func parseRateString(s string) Rate {
	r := Rate{Big: decimal.Big{}}
	r.SetString(s)
	return r
}

func (r Rate) toString() string {
	return r.String()
}
