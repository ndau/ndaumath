package eai

import (
	"encoding"
	"fmt"
	gomath "math"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/pkg/errors"
)

//go:generate msgp

// A Rate defines a rate of increase over time.
//
// EAI is not interest. Ndau never earns interest. However,
// for ease of explanation, we will borrow some terminology
// from interest rate calculations.
//
// EAI is continuously compounded according to the formula
//
//   eai = balance * (e ^ (rate * time) - 1)
//
// Rates are expressed in percent per year.
//
// Thus, 100 ndau at 1 percent rate over 1 year yields 1.00501670 ndau EAI.
//
// The use of continuously compounded interest instead of simple interest
// aids in EAI predictability: using simple interest, an account which
// simply accumulates its EAI, whose node won frequently, would see a higher
// rate of actual return than an identical account whose node won infrequently.
// Continuously compounded interest avoids that issue: both accounts will
// see the same rate of return; the benefit of the one registered to the
// frequent node is that it sees the increase more often.
//
// We use a signed int so that json2msgp won't need type hints for encoding
// rate tables in system variables.  If we use a rate denominator of 1e12,
// corresponding to a rate of 100%, then 63 bits gives us enough room to
// handle rates in the hundreds of millions of percents.
type Rate int64

var (
	fracdigits int
	ratefmt    string
	ratere     *regexp.Regexp
)

func init() {
	// fracdigits: how many digits go behind the decimal?
	// computed here so that if constants.RateDenominator ever changes,
	// this stays automatically in sync
	fracdigits = int(gomath.Floor(gomath.Log10(constants.RateDenominator))) - 2
	// ratefmt: just there to add the right number of 0s to the fractional part
	// of the text serialization of the rate
	ratefmt = fmt.Sprintf("%%d.%%0%dd", fracdigits)
	// ratere: parse a rate into pct (before the decimal) and frac (after the decimal)
	// strings, which can be used to regenerate the rate
	ratere = regexp.MustCompile(fmt.Sprintf(`^\s*(?P<pct>\d+)(\.(?P<frac>\d{1,%d}))?%%\s*$`, fracdigits))
}

// String writes this Rate as a string
func (r Rate) String() string {
	onePct := RateFromPercent(1)
	rs := fmt.Sprintf(ratefmt, r/onePct, r%onePct)
	for rs[len(rs)-1] == '0' {
		rs = rs[:len(rs)-1]
	}
	if rs[len(rs)-1] == '.' {
		rs = rs[:len(rs)-1]
	}
	return rs + "%"
}

// ParseRate attempts to parse a Rate from the provided string
func ParseRate(s string) (Rate, error) {
	match := ratere.FindStringSubmatch(s)
	result := make(map[string]string)
	for i, name := range ratere.SubexpNames() {
		if i != 0 && name != "" && i < len(match) {
			result[name] = match[i]
		}
	}

	pcts, ok := result["pct"]
	if !ok {
		return Rate(0), errors.New("failed to parse rate")
	}
	pct, err := strconv.ParseUint(pcts, 10, 64)
	if err != nil {
		return Rate(0), errors.Wrap(err, "parsing pct")
	}
	out := RateFromPercent(1) * Rate(pct)

	fracs, ok := result["frac"]
	if ok {
		if len(fracs) > fracdigits {
			fracs = fracs[:fracdigits]
		} else if len(fracs) < fracdigits {
			iters := fracdigits - len(fracs)
			for i := 0; i < iters; i++ {
				fracs += "0"
			}
		}

		frac, err := strconv.ParseUint(fracs, 10, 64)
		if err != nil {
			return Rate(0), errors.Wrap(err, "parsing frac")
		}
		out += Rate(frac)
	}

	return out, nil
}

// RateFromPercent returns a Rate whose value is that of the input, as percent.
//
// i.e. to express 1%, `nPercent` should equal `1`
func RateFromPercent(nPercent uint64) Rate {
	return Rate(nPercent * constants.RateDenominator / 100)
}

//msgp:tuple RTRow

// RTRow is a single row of a rate table
type RTRow struct {
	From math.Duration
	Rate Rate
}

var _ encoding.TextMarshaler = (*RTRow)(nil)
var _ encoding.TextUnmarshaler = (*RTRow)(nil)

// MarshalText implements encoding.TextMarshaler
func (r RTRow) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s:%s", r.From, r.Rate)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (r *RTRow) UnmarshalText(text []byte) error {
	if !utf8.Valid(text) {
		return errors.New("text was not utf-8")
	}
	parts := strings.Split(string(text), ":")
	if len(parts) != 2 {
		return errors.New("invalid fmt: expected single ':'")
	}
	var err error
	r.From, err = math.ParseDuration(parts[0])
	if err != nil {
		return errors.Wrap(err, "parsing from")
	}
	r.Rate, err = ParseRate(parts[1])
	if err != nil {
		return errors.Wrap(err, "parsing rate")
	}

	return nil
}

// A RateTable defines a stepped sequence of EAI rates which apply
// at varying durations.
//
// It is a logic error if the elements of a RateTable are not sorted
// in increasing order by their From field.
type RateTable []RTRow

// RateAt returns the rate in a RateTable for a given point
func (rt RateTable) RateAt(point math.Duration) Rate {
	rate := Rate(0)
	// the nature of rate tables is that we want the smallest rate
	// for which point >= row.From. The obvious way would be to iterate
	// in reverse, and return the first time that the point >= the row's
	// From value. However, reverse iteration is tedious in go, so we
	// take a different tack instead.
	//
	// We iterate forward, but keep a cache of the row rates, so that the
	// currently active rate always trails behind the table by one row.
	// This means that the first time the point < row.From, we can return
	// the active rate.
	for _, row := range rt {
		if point < row.From {
			return rate
		}
		rate = row.Rate
	}
	return rate
}

//msgp:tuple RSRow

// RSRow is a single row of a rate slice.
type RSRow struct {
	Duration math.Duration
	Rate     Rate
}

// A RateSlice is derived from a RateTable, optimized for computation.
//
// Whereas a RateTable is meant to be easy for humans to understand, a
// RateSlice is more efficient for computation. It is a list of rates, and
// the actual duration over which each rate is active.
type RateSlice []RSRow

// Slice transforms a RateTable into a form more amenable for computation.
//
// Rates vary over time, and we want to efficiently compute the sum of interest
// at varying rates. Instead of repeatedly calling RateAt, it's more efficient
// to perform the calculation once to slice the affected data out of the
// RateTable.
func (rt RateTable) Slice(from, to, offset math.Duration) RateSlice {
	return rt.SliceF(from, to, offset, 0)
}

// SliceF transforms a RateTable into a form more amenable for computation.
//
// Rates vary over time, and we want to efficiently compute the sum of interest
// at varying rates. Instead of repeatedly calling RateAt, it's more efficient
// to perform the calculation once to slice the affected data out of the
// RateTable.
//
// Let's diagram the variables in play in here:
// (parentheticized variables are not present)
//
//  Timestamps
//       │ (effective account open)
//       │   │        (lastEAICalc)
//       │   │           │  (notify)              (blockTime)  (lock.UnlocksOn)
// TIME ─┼───┼───────────┼─────┼─────────────────────┼────────────┼──>
//       │   │           │     ├────── freeze ───────┤            │
//       │   │           │     └───────────── offset ─────────────┘
//       │   ├── from ───┴──── (lastEAICalcAge) ─────┤
//       │   └──────────────── to ───────────────────┘
//   Durations
//
// Where freeze == 0, this function returns the rate slice from (from+offset)
// to (to+offset).
//
//   R3                                         ┌────|────────...
//   R2                            ┌────────────┘ / /|
//   R1              ┌────|────────┘ / / / / / / / / |
//   R0  ────────────┘    | / / / / / / / / / / / / /|
//                   (from+offset)                (to+offset)
//
// Where freeze != 0, this function returns the rate slice from (from+offset)
// to (to+offset), but with the actual rate frozen at the freeze point.
//
// (This diagram is not to the same scale as the timeline overview above.)
//
//   R3                                         ┌────|────────...
//   R2                            ┌─────────|───────|──
//   R1              ┌────|────────┘ / / / / | / / / |
//   R0  ────────────┘    | / / / / / / / / /|/ / / /|
//                   (from+offset)           |    (to+offset)
//                                 (to+offset-freeze)
func (rt RateTable) SliceF(from, to, offset, freeze math.Duration) RateSlice {
	if to <= from {
		// when actual duration is 0, it's fine to fake that the actual
		// rate is also 0
		return RateSlice{RSRow{}}
	}

	if freeze < 0 {
		freeze = -freeze
	}

	fromEffective := from + offset
	toEffective := to + offset
	notify := toEffective - freeze

	// the computation can't result in -2, so if after the loop
	// this remains, we know we never touched this var
	const uninitialized = -2
	fromI := uninitialized
	toI := uninitialized
	notifyI := uninitialized

	for index, row := range rt {
		if fromI == uninitialized && fromEffective < row.From {
			fromI = index - 1
		}
		if toI == uninitialized && toEffective < row.From {
			toI = index - 1
		}
		if notifyI == uninitialized && notify < row.From {
			notifyI = index - 1
		}
		if fromI != uninitialized && toI != uninitialized && notifyI != uninitialized {
			break
		}
	}
	// if either variable comes out of the loop wihtout being initialized,
	// the appropriate row index is the highest in the table
	if fromI == uninitialized {
		fromI = len(rt) - 1
	}
	if toI == uninitialized {
		toI = len(rt) - 1
	}
	if notifyI == uninitialized {
		notifyI = len(rt) - 1
	}

	rateFor := func(idx int) Rate {
		if idx == -1 {
			return Rate(0)
		}
		return rt[idx].Rate
	}

	// if we froze before the from point, we have one period at the frozen rate
	if freeze != 0 && notify < fromEffective {
		return RateSlice{RSRow{Rate: rateFor(notifyI), Duration: to - from}}
	}

	// if from and to are in the same rate block, or
	// from and the freeze point are in the same rate block, we have one period
	// at the from rate
	if fromI == toI || fromI == notifyI {
		return RateSlice{RSRow{Rate: rateFor(fromI), Duration: to - from}}
	}
	numRows := 1 - fromI
	if toI <= notifyI {
		numRows += toI
	} else {
		numRows += notifyI
	}

	// ok, the rest is relatively straightforward. We need special
	// handling for the first and last rate, because they have partial
	// durations; the rest are just copies from the rate table
	rs := make(RateSlice, numRows)
	// - it's safe to index rt[fromI+1] because if fromI were the max value,
	//   then we would have already returned: fromI must equal toI
	// - we know that freezePoint > rt[fromI+1].From, because if fromI == notifyI,
	//   we would have already returned. As we're here, we know that the
	//   freeze point isn't in this first block.
	rs[0] = RSRow{Rate: rateFor(fromI), Duration: rt[fromI+1].From - (from + offset)}
	// freezing within the final rate block has no effect on the calculation
	if notifyI == toI {
		// it's safe to index rt[toI] because if toI were -1,
		// then we would have already returned: fromI must equal toI
		rs[numRows-1] = RSRow{Rate: rateFor(toI), Duration: toEffective - rt[toI].From}
	} else {
		rs[numRows-1] = RSRow{Rate: rateFor(notifyI), Duration: freeze + notify - rt[notifyI].From}
	}

	upperBoundI := toI
	if notifyI < toI {
		upperBoundI = notifyI
	}

	// indexing rt[fromI+i+1] is safe because fromI+i+1 == toI at max i
	for i := 1; i < upperBoundI-fromI; i++ {
		rs[i] = RSRow{
			Rate:     rt[fromI+i].Rate,
			Duration: rt[fromI+i+1].From - rt[fromI+i].From,
		}
	}

	return rs
}

var (
	// DefaultUnlockedEAI is the default base rate table for unlocked accounts
	//
	// The UnlockedEAI rate table is a system variable which is adjustable
	// whenever the BPC desires, but for testing purposes, we use this
	// approximation as a default.
	//
	// Defaults drawn from https://tresor.it/p#0041o9iot7hm4kb5y707es7o/Oneiro%20Company%20Info/Whitepapers%20and%20Presentations/ndau%20Whitepaper%201.3%2020180425%20Final.pdf
	// page 15.
	DefaultUnlockedEAI RateTable

	// DefaultLockBonusEAI is the bonus rate for locks of varying length
	//
	// The LockBonusEAI rate table is a system variable which is adjustable
	// whenever the BPC desires, but for testing purposes, we use this
	// approximation as a default.
	//
	// Defaults drawn from https://tresor.it/p#0041o9iot7hm4kb5y707es7o/Oneiro%20Company%20Info/Whitepapers%20and%20Presentations/ndau%20Whitepaper%201.3%2020180425%20Final.pdf
	// page 15.
	DefaultLockBonusEAI RateTable
)

func init() {
	for i := uint64(1); i < 10; i++ {
		DefaultUnlockedEAI = append(DefaultUnlockedEAI, RTRow{
			Rate: RateFromPercent(uint64(i + 1)),
			// the wrapper here for uint64 serves as a notice to gomobile that it should
			// assume this constant is 64 bits. Otherwise this breaks gomobile.
			From: math.Duration(i * 30 * uint64(math.Day)),
		})
	}

	DefaultLockBonusEAI = RateTable{
		RTRow{
			From: math.Duration(3 * 30 * math.Day),
			Rate: RateFromPercent(1),
		},
		RTRow{
			From: math.Duration(6 * 30 * math.Day),
			Rate: RateFromPercent(2),
		},
		RTRow{
			From: math.Duration(1 * math.Year),
			Rate: RateFromPercent(3),
		},
		RTRow{
			From: math.Duration(2 * math.Year),
			Rate: RateFromPercent(4),
		},
		RTRow{
			From: math.Duration(3 * math.Year),
			Rate: RateFromPercent(5),
		},
	}
}
