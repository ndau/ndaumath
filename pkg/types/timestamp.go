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
	"time"

	"github.com/ndau/ndaumath/pkg/constants"
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

// ensure Timestamp implements encoding.Text(Un)Marshaler
var _ encoding.TextMarshaler = (*Timestamp)(nil)
var _ encoding.TextUnmarshaler = (*Timestamp)(nil)

// ParseTimestamp creates a timestamp from an ISO-3933 string
func ParseTimestamp(s string) (Timestamp, error) {
	err := errors.New("timestamp matched no known format")
	var ts time.Time
	for _, format := range []string{constants.TimestampFormat, time.RFC3339[:len(time.RFC3339)-5]} {
		ts, err = time.Parse(format, s)
		if err == nil {
			break
		}
	}
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

// AsTime converts a Timestamp into a time.Time object
//
// TODO: implement this in a way which ensures its monotonic properties
func (t Timestamp) AsTime() time.Time {
	return constants.Epoch.Add(time.Duration(int64(t)) * time.Microsecond)
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
	return t.AsTime().Format(constants.TimestampFormat)
}

// MarshalText implements encoding.TextMarshaler
func (t Timestamp) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (t *Timestamp) UnmarshalText(text []byte) error {
	if t == nil {
		return errors.New("nil Timestamp")
	}
	tt, err := ParseTimestamp(string(text))
	if err != nil {
		return err
	}
	*t = tt
	return nil
}
