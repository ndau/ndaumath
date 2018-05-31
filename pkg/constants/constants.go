package constants

import (
	"math"
	"time"
)

// CurrencyName is the official name of the currency
const CurrencyName = "ndau"

// CurrencyQuantum is the official name of the smallest possible
// unit of the currency
const CurrencyQuantum = "napu"

// QuantaPerUnit is the number of quanta in a single unit of
// ndau
const QuantaPerUnit = 100000000

// MaxQuantaPerAddress is the number of quanta that
// can be tracked in a single address
const MaxQuantaPerAddress = math.MaxInt64

// Epoch is the basic moment from which Ndau chain time calculations begin
var Epoch time.Time

func init() {
	tz, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	Epoch = time.Date(
		2018, time.January, 18,
		14, 21, 0,
		0, tz,
	)
}
