package constants

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"math"
	"regexp"
	"time"
)

const (
	// CurrencyName is the official name of the currency.
	CurrencyName = "ndau"

	// CurrencyQuantum is the official name of the smallest possible
	// unit of the currency.
	CurrencyQuantum = "napu"

	// QuantaPerUnit is the number of quanta in a single unit of
	// ndau.
	QuantaPerUnit = 100000000

	// NapuPerNdau is a more human-friendly synonym of QuantaPerUnit
	NapuPerNdau = QuantaPerUnit

	// RateDenominator is the implied denominator for interest rates.
	//
	// EAI rates are expressed as integers, and integer math is performed
	// for determinism. However, the interest-like calculations of EAI
	// expect non-integer rates. We square this circle by treating them
	// as rationals using an implied denominator, specified here.
	//
	// RateDenominator must be tuned to preserve precision. Because of
	// technical details of EAI calculation, it cannot be greater than
	// half of a uint64: approximately 10**18. At the same time, the effective
	// precision of an EAI factor is the precision of the RateDenominator
	// reduced by the fraction of the year over which the EAI is calculated.
	//
	// One hour, as a fraction of a year, is a little more than 1/10000.
	// Therefore, by making the precision of the RateDenominator 10000 times
	// the QuantaPerUnit, we ensure that we don't lose too much precision
	// in EAI calculations. However, EAI calculations much more frequent than
	// an hour can be expected to lose some EAI to dust truncation.
	RateDenominator = QuantaPerUnit * 10000

	// MaxQuantaPerAddress is the number of quanta that
	// can be tracked in a single address.
	MaxQuantaPerAddress = math.MaxInt64

	// TimestampFormat is the format string used to parse timestamps.
	TimestampFormat = "2006-01-02T15:04:05.000000Z"

	// EpochStart is the text representation of the start time of our Epoch.
	EpochStart = "2000-01-01T00:00:00.000000Z"

	// MaxTimestamp is the maximum value a timestamp can take on
	MaxTimestamp = math.MaxInt64

	// MinTimestamp is the minimum value a timestamp can take on
	MinTimestamp = 0

	// DurationFormat is the format regex used to parse durations
	DurationFormat = `(?i)^(?P<neg>-)?p?((?P<years>\d+)y)?((?P<months>\d{1,2})m)?((?P<days>\d{1,2})d)?(t((?P<hours>\d{1,2})h)?((?P<minutes>\d{1,2})m)?((?P<seconds>\d{1,2})s)?((?P<micros>\d{1,6})[μu]s?)?)?$`

	// MaxDuration is the maximum value a duration can contain
	MaxDuration = math.MaxInt64

	// MinDuration is the minimum value a duration can contain
	MinDuration = math.MinInt64
)

var (
	// Epoch is the basic moment from which Ndau chain time calculations begin.
	Epoch time.Time

	// DurationRE is the regular expression used to parse Durations
	DurationRE *regexp.Regexp
)

func init() {
	Epoch, _ = time.Parse(TimestampFormat, EpochStart)
	DurationRE = regexp.MustCompile(DurationFormat)
}
