package types

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"encoding"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ndau/ndaumath/pkg/constants"
	"github.com/ndau/ndaumath/pkg/signed"
)

//go:generate msgp -tests=0

// A Duration is the difference between two Timestamps.
//
// It can be negative if the timestamps are out of order.
type Duration int64

// ensure Timestamp implements encoding.Text(Un)Marshaler
var _ encoding.TextMarshaler = (*Duration)(nil)
var _ encoding.TextUnmarshaler = (*Duration)(nil)

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
// and `t1m` is one minute. Per RFC3339, leading `p` chars
// are allowed.
//
// There is no `w` symbol for weeks; use multiples of days
// or months instead.
func ParseDuration(s string) (Duration, error) {
	match := constants.DurationRE.FindStringSubmatch(s)
	if match == nil {
		return Duration(0), fmt.Errorf("invalid duration format")
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
				return fmt.Errorf("invalid integer: %s", result[name])
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
	if err := addTime("micros", Microsecond); err != nil {
		return Duration(0), err
	}

	if result["neg"] != "" {
		duration = -duration
	}

	return duration, nil
}

// String represents a Duration as a human-readable string
func (d Duration) String() string {
	value := int64(d)
	out := ""
	if value < 0 {
		out = "-"
		value = -value
	}
	divmod := func(divisor, dividend int64) (int64, int64) {
		return divisor / dividend, divisor % dividend
	}
	extract := func(symbol string, unit int64) {
		var units int64
		units, value = divmod(value, unit)
		if units > 0 {
			out += fmt.Sprintf("%d%s", units, symbol)
		}
	}
	extract("y", Year)
	extract("m", Month)
	extract("d", Day)
	if value > 0 {
		out += "t"
	}
	extract("h", Hour)
	extract("m", Minute)
	extract("s", Second)
	extract("us", Microsecond)

	if out == "" {
		// input duration was 0
		out = "t0s" // seconds are the fundamental unit
	}

	return out
}

// UpdateWeightedAverageAge computes the weighted average age. Note that this
// function may cause order-dependent behavior; it does integer division, and
// for small values, the order in which updates to WAA are applied may be
// significant. We found and fixed a bug related to this in the CreditEAI
// calculation; see duration_test.go for details.
func (d *Duration) UpdateWeightedAverageAge(
	sinceLastUpdate Duration,
	transferQty Ndau,
	previousBalance Ndau,
) error {
	waa := int64(*d) + int64(sinceLastUpdate)
	if int64(transferQty) >= 0 {
		newBalance, err := previousBalance.Add(transferQty)
		if err != nil {
			return err
		}
		// we have to use bigints to prevent the multiplication from overflow
		nb := int64(newBalance)
		pb := int64(previousBalance)
		dur := int64(*d + sinceLastUpdate)
		if nb > 0 {
			waa, err = signed.MulDiv(dur, pb, nb)
			if err != nil {
				return errors.New("Duration overflow in UpdateWeightedAverageAge")
			}
		}
	}
	*d = Duration(waa)
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (d *Duration) UnmarshalText(text []byte) error {
	if d == nil {
		return errors.New("nil Duration")
	}
	dd, err := ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = dd
	return nil
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
