package unsigned

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

var bounds []uint64

func init() {
	// Table was constructed in python:
	//
	// >>> from math import ceil, exp
	// >>> i = 0
	// >>> bounds = []
	// >>> while exp(i) <= 2**64:
	// ...     bounds.append(ceil(exp(i))-1)
	// ...     i += 1
	// ...
	// True
	bounds = []uint64{
		0, 2, 7, 20, 54, 148, 403, 1096, 2980, 8103, 22026, 59874, 162754, 442413,
		1202604, 3269017, 8886110, 24154952, 65659969, 178482300, 485165195,
		1318815734, 3584912846, 9744803446, 26489122129, 72004899337, 195729609428,
		532048240601, 1446257064291, 3931334297144, 10686474581524, 29048849665247,
		78962960182680, 214643579785916, 583461742527454, 1586013452313430,
		4311231547115194, 11719142372802611, 31855931757113755, 86593400423993743,
		235385266837019999, 639843493530054911, 1739274941520500991,
		4727839468229346303, 12851600114359308287, ^uint64(0),
	}
}

// LnInt computes the floor of the natural logarithm of its input.
func LnInt(x uint64) (i int) {
	// It uses a linear scan, which for tables of this size on modern
	// architectures will be faster than a binary search.
	for i = -1; ; i++ {
		if x <= bounds[i+1] {
			return
		}
	}
}
