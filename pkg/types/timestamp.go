package types

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
)

//go:generate msgp -tests=0

// A Timestamp is a single moment in time.
//
// It is monotonically increasing with the passage of time, and represents
// the number of microseconds since the epoch. It has no notion of leap time,
// time zones, or other complicating human factors.
// A timestamp can never be negative, but for mathematical simplicity we represent
// it with an int64. The total range of timestamps is almost 300,000 years.
type Timestamp int64

// A Duration is the difference between two Timestamps.
//
// It can be negative if the timestamps are out of order.
type Duration int64

// ParseTimestamp creates a timestamp from an ISO-3933 string
func ParseTimestamp(s string) (Timestamp, error) {
	ts, err := time.Parse(constants.TimestampFormat, s)
	if err != nil {
		return 0, err
	}
	return TimestampFrom(ts)
}

// TimestampFrom creates a Timestamp given a time.Time object
func TimestampFrom(t time.Time) (Timestamp, error) {
	// because this uses the standard library, it will overflow
	// some 290 years after the epoch
	//
	// TODO: implement this in a way which ensures its monotonic properties
	durationSinceEpoch := t.Sub(constants.Epoch)
	if durationSinceEpoch < 0 {
		return Timestamp(0), errors.New("date is before Epoch start")
	}
	return Timestamp(int64(durationSinceEpoch / time.Microsecond)), nil
}

// Compare implements comparison for Timestamp.
// (useful in sorting)
func (t Timestamp) Compare(o Timestamp) int {
	if t < o {
		return -1
	} else if t > o {
		return 1
	}
	return 0
}

// Since measures the Duration between two Timestamps.
// It will be positive when the argument is older, so present.Since(past) > 0
func (t Timestamp) Since(o Timestamp) Duration {
	return Duration(t - o)
}

// Add adds the supplied Duration to the given Timestamp
// If the result is negative, returns 0
// If the result overflows, returns MaxTimestamp
func (t Timestamp) Add(d Duration) Timestamp {
	ts := Timestamp(int64(t) + int64(d))
	if ts < constants.MinTimestamp {
		if d < 0 {
			return constants.MinTimestamp
		}
		return constants.MaxTimestamp
	}
	return ts
}

// Sub subtracts the supplied Duration from the given Timestamp
func (t Timestamp) Sub(d Duration) Timestamp {
	ts := Timestamp(int64(t) - int64(d))
	if ts < constants.MinTimestamp {
		if d > 0 {
			return constants.MinTimestamp
		}
		return constants.MaxTimestamp
	}
	return ts
}

func (t Timestamp) String() string {
	tt := constants.Epoch.Add(time.Duration(t) * time.Microsecond)
	return tt.Format(constants.TimestampFormat)
}

const (
	// Microsecond is a thousandth of a millisecond
	Microsecond = 1
	// Millisecond is a thousandth of a second
	Millisecond = Microsecond * 1000
	// Second is the duration of 9 192 631 770 periods of the
	// radiation corresponding to the transition between the two
	// hyperfine levels of the ground state of the cesium 133 atom,
	// per the 13th CGPM (1967).
	Second = Millisecond * 1000
	// Minute is exactly 60 Seconds
	Minute = Second * 60
	// Hour is exactly 60 Minutes
	Hour = Minute * 60
	// Day is exactly 24 Hours
	Day = Hour * 24
	// Month is exactly 30 Days
	Month = Day * 30
	// Year is exactly 365 days
	Year = Day * 365
)

// DurationFrom creates a Duration given a time.Duration object
func DurationFrom(d time.Duration) Duration {
	return Duration(d / time.Millisecond * Millisecond)
}

// TimeDuration converts a Duration into a time.Duration
func (d Duration) TimeDuration() time.Duration {
	return time.Duration(int64(d) / Millisecond * int64(time.Millisecond))
}

// ParseDuration creates a duration from a duration string
//
// Allowable durations broadly follow the RFC3339 duration
// specification: `\dy\dm\dd(t\dh\dm\ds)`. Note that `m`
// is used for both months and minutes: `1m` is one month,
// and `t1m` is one minute. Leading `p` chars are allowed.
//
// There is no `w` symbol for weeks; use multiples of days
// or months instead.
//
// Integral seconds are the smallest unit of time which
// can be parsed.
func ParseDuration(s string) (Duration, error) {
	match := constants.DurationRE.FindStringSubmatch(s)
	if match == nil {
		return Duration(0), fmt.Errorf("Invalid duration format")
	}

	// get match groups by name:
	// https://stackoverflow.com/a/20751656/504550
	result := make(map[string]string)
	for i, name := range constants.DurationRE.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	duration := Duration(0)
	addTime := func(name string, unit uint64) error {
		if result[name] != "" {
			value, err := strconv.ParseUint(result[name], 10, 64)
			if err != nil {
				return fmt.Errorf("Invalid integer: %s", result[name])
			}
			duration += Duration(value * unit)
		}
		return nil
	}

	if err := addTime("years", Year); err != nil {
		return Duration(0), err
	}
	if err := addTime("months", Month); err != nil {
		return Duration(0), err
	}
	if err := addTime("days", Day); err != nil {
		return Duration(0), err
	}
	if err := addTime("hours", Hour); err != nil {
		return Duration(0), err
	}
	if err := addTime("minutes", Minute); err != nil {
		return Duration(0), err
	}
	if err := addTime("seconds", Second); err != nil {
		return Duration(0), err
	}

	return duration, nil
}

// String represents a Duration as a human-readable string
func (d Duration) String() string {
	td := d.TimeDuration()
	day := 24 * time.Hour
	if td < day {
		return td.String()
	}
	return fmt.Sprintf("%dd%s", td/day, (td % day).String())
}

// UpdateWeightedAverageAge computes the weighted average age
func (d *Duration) UpdateWeightedAverageAge(
	sinceLastUpdate Duration,
	transferQty Ndau,
	previousBalance Ndau,
) error {
	waa := new(big.Int)
	if int64(transferQty) < 0 {
		waa.Add(big.NewInt(int64(*d)), big.NewInt(int64(sinceLastUpdate)))
	} else {
		newBalance, err := previousBalance.Add(transferQty)
		if err != nil {
			return err
		}
		// we have to use bigints to prevent the multiplication from overflow
		nb := big.NewInt(int64(newBalance))
		pb := big.NewInt(int64(previousBalance))
		dur := big.NewInt(int64(*d + sinceLastUpdate))
		waa.Div(waa.Mul(dur, pb), nb)
	}
	if !waa.IsInt64() {
		return errors.New("Duration overflow in UpdateWeightedAverageAge")
	}
	*d = Duration(waa.Int64())
	return nil
}
