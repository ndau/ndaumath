package constants

import (
	"math"
	"regexp"
	"time"
)

// CurrencyName is the official name of the currency.
const CurrencyName = "ndau"

// CurrencyQuantum is the official name of the smallest possible
// unit of the currency.
const CurrencyQuantum = "napu"

// QuantaPerUnit is the number of quanta in a single unit of
// ndau.
const QuantaPerUnit = 100000000

// MaxQuantaPerAddress is the number of quanta that
// can be tracked in a single address.
const MaxQuantaPerAddress = math.MaxInt64

// TimestampFormat is the format string used to parse timestamps.
const TimestampFormat = "2006-01-02T15:04:05Z"

// EpochStart is the text representation of the start time of our Epoch.
const EpochStart = "2000-01-01T00:00:00Z"

// Epoch is the basic moment from which Ndau chain time calculations begin.
var Epoch time.Time

// MaxTimestamp is the maximum value a timestamp can take on
const MaxTimestamp = math.MaxInt64

// MinTimestamp is the minimum value a timestamp can take on
const MinTimestamp = 0

// DurationFormat is the format regex used to parse durations
const DurationFormat = `(?i)^p?((?P<years>\d+)y)?((?P<months>\d{1,2})m)?((?P<days>\d{1,2})d)?(t((?P<hours>\d{1,2})h)?((?P<minutes>\d{1,2})m)?((?P<seconds>\d{1,2})s)?)?$`

// DurationRE is the regular expression used to parse Durations
var DurationRE *regexp.Regexp

// MaxDuration is the maximum value a duration can contain
const MaxDuration = math.MaxInt64

// MinDuration is the minimum value a duration can contain
const MinDuration = math.MinInt64

func init() {
	Epoch, _ = time.Parse(TimestampFormat, EpochStart)
	DurationRE = regexp.MustCompile(DurationFormat)
}
