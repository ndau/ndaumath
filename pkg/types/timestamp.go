package types

import (
	"time"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
)

// A Timestamp is a single moment in time.
//
// It is monotonically increasing with the passage of time, and represents
// the number of microseconds since the epoch. It has no notion of leap time,
// time zones, or other complicating human factors.
type Timestamp uint64

// A Duration is the difference between two Timestamps
//
// It is an absolute quantity: a negative duration has no meaning.
type Duration uint64

// TimestampFrom creates a Timestamp given a time.Time object
func TimestampFrom(t time.Time) Timestamp {
	// becuase this uses the standard library, it will overflow
	// some 290 years after the epoch
	//
	// TODO: implement this in a way which ensures its monotonic properties
	durationSinceEpoch := t.Sub(constants.Epoch)
	return Timestamp(uint64(durationSinceEpoch / time.Microsecond))
}

// CurrentTimestamp returns the timestamp of the current moment
func CurrentTimestamp() Timestamp {
	// because this uses standard library time, it depends on the
	// system clock's accuracy and timezone setting
	return TimestampFrom(time.Now())
}

// Between measures the Duration between two Timestamps
func (t Timestamp) Between(o Timestamp) Duration {
	var big, small Timestamp
	if uint64(t) < uint64(o) {
		big = o
		small = t
	} else {
		big = t
		small = o
	}
	return Duration(uint64(big) - uint64(small))
}

// Add adds the supplied Duration to the given Timestamp
func (t Timestamp) Add(d Duration) Timestamp {
	return Timestamp(uint64(t) + uint64(d))
}

// Sub subtracts the supplied Duration from the given Timestamp
func (t Timestamp) Sub(d Duration) Timestamp {
	return Timestamp(uint64(t) - uint64(d))
}

const (
	// Millisecond is a thousandth of a second
	Millisecond = 1
	// Second is the duration of 9 192 631 770 periods of the
	// radiation corresponding to the transition between the two
	// hyperfine levels of the ground state of the cesium 133 atom,
	// per the 13th CGPM (1967).
	Second = Millisecond * 1000
	// Day is exactly 86400 Seconds
	Day = Second * 86400
	// Year is exactly 365 days
	Year = Day * 365
)
