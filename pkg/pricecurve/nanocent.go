package pricecurve

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

//go:generate msgp

// A Nanocent is one billionth of one hundredth of one USD.
//
// It is fundamentally an integer and is computed using only integer math, for
// perfect determinism.
type Nanocent int64

const (
	// Dollar is 10^11 nanocents
	Dollar = 100000000000
)

var dollarsRE *regexp.Regexp

func init() {
	dollarsRE = regexp.MustCompile(`^(?P<neg>-?)\$?(?P<dollars>[\d,_]+)(\.(?P<cents>\d{2,11}))?$`)
}

// ParseDollars parses strings expressed in dollars and returns nanocents
func ParseDollars(dollars string) (Nanocent, error) {
	dollars = strings.TrimSpace(dollars)
	// allow for separation by just eliminating spacing chars
	// there isn't a great way to do this within the regex itself
	dollars = strings.Replace(dollars, ",", "", -1)
	dollars = strings.Replace(dollars, "_", "", -1)

	// perform regex matching
	match := dollarsRE.FindStringSubmatch(dollars)
	if len(match) == 0 {
		return 0, fmt.Errorf("'%s' doesn't look like dollars", dollars)
	}

	// get submatches by name
	submatches := make(map[string]string)
	for i, name := range dollarsRE.SubexpNames() {
		if i != 0 && i < len(match) && name != "" {
			submatches[name] = match[i]
		}
	}

	// parse integers
	var err error
	d := int64(0)
	if submatches["dollars"] != "" {
		d, err = strconv.ParseInt(submatches["dollars"], 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, "parsing dollars as integer: "+submatches["dollars"])
		}
	}
	c := int64(0)
	if submatches["cents"] != "" {
		c, err = strconv.ParseInt(submatches["cents"], 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, "parsing cents as integer: "+submatches["cents"])
		}
		for i := 0; i < 11-len(submatches["cents"]); i++ {
			c *= 10
		}
	}

	// add it all up
	nc := Nanocent(d*Dollar + c)

	// handle negatives
	if submatches["neg"] == "-" {
		nc = -nc
	}

	return nc, err
}
