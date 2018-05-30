package basics

import (
	"math"
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
