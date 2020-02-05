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
	"fmt"
	"math/rand"
	"testing"

	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/stretchr/testify/require"
)

func TestDuration_UpdateWeightedAverageAge(t *testing.T) {
	// we derive the tests from some canonical data
	// computed in excel and validated by hand
	data := []struct {
		day      int
		transfer int
		balance  int
		waa      int
	}{
		{0, 0, 0, 0},       // dummy entry
		{0, 0, 0, 0},       // create an empty account
		{0, 100, 100, 0},   // give it a balance
		{30, 0, 100, 30},   // eai calculations; no transfer
		{30, 50, 150, 20},  // transfer in
		{40, -50, 100, 30}, // withdraw
		{60, 100, 200, 25}, // transfer in
		{80, -200, 0, 45},  // withdraw everything
		{100, 100, 100, 0}, // start again from 0
	}

	for index := range data {
		if index > 0 {
			sinceLastUpdate := Duration((data[index].day - data[index-1].day) * Day)
			transferQty := Ndau(data[index].transfer * constants.QuantaPerUnit)
			previousBalance := Ndau(data[index-1].balance * constants.QuantaPerUnit)
			waa := Duration(data[index-1].waa * Day)
			expectedWAA := Duration(data[index].waa * Day)

			t.Run(fmt.Sprintf("row %d", index), func(t *testing.T) {
				err := (&waa).UpdateWeightedAverageAge(sinceLastUpdate, transferQty, previousBalance)
				if err != nil {
					t.Errorf("Update weighted average age returned err: %s", err.Error())
				}
				if waa != expectedWAA {
					t.Errorf("WAA: %d; expected %d", waa, expectedWAA)
				}
			})
		}
	}
}

func TestParseDuration(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Duration
		wantErr bool
	}{
		{"<blank>", args{""}, Duration(0), false},
		{"t0s", args{"t0s"}, Duration(0), false},
		{"t1s", args{"t1s"}, Duration(1 * Second), false},
		{"1m", args{"1m"}, Duration(1 * Month), false},
		{"t1m", args{"t1m"}, Duration(1 * Minute), false},
		{"p1y2m3dt4h5m6s", args{"p1y2m3dt4h5m6s"}, Duration(36993906000000), false},
		{"P1Y2M3DT4H5M6S", args{"P1Y2M3DT4H5M6S"}, Duration(36993906000000), false},
		{"1y2m3dt4h5m6s7u", args{"1y2m3dt4h5m6s7u"}, Duration(36993906000007), false},
		{"1h", args{"1h"}, Duration(0), true},               // needs t
		{"100y", args{"100y"}, Duration(100 * Year), false}, // 3 digit year
		{"100m", args{"100m"}, Duration(0), true},           // 3 digit anything else
		{"100d", args{"100d"}, Duration(0), true},           // 3 digit anything else
		{"t100h", args{"t100h"}, Duration(0), true},         // 3 digit anything else
		{"t100m", args{"t100m"}, Duration(0), true},         // 3 digit anything else
		{"t100s", args{"t100s"}, Duration(0), true},         // 3 digit anything else
		{"t1u", args{"t1u"}, Duration(1), false},
		{"t1us", args{"t1us"}, Duration(1), false},
		{"t1μ", args{"t1μ"}, Duration(1), false},
		{"t1μs", args{"t1μs"}, Duration(1), false},
		{"t999999μ", args{"t999999μ"}, Duration(999999), false},
		{"t1000000μ", args{"t1000000μ"}, Duration(0), true},
		{"-t1s", args{"-t1s"}, -Duration(1 * Second), false},
		{"-1m", args{"-1m"}, -Duration(1 * Month), false},
		{"-t1m", args{"-t1m"}, -Duration(1 * Minute), false},
		{"-p1y2m3dt4h5m6s", args{"-p1y2m3dt4h5m6s"}, -Duration(36993906000000), false},
		{"-P1Y2M3DT4H5M6S", args{"-P1Y2M3DT4H5M6S"}, -Duration(36993906000000), false},
		{"-1y2m3dt4h5m6s7u", args{"-1y2m3dt4h5m6s7u"}, -Duration(36993906000007), false},
		{"-100y", args{"-100y"}, -Duration(100 * Year), false}, // 3 digit year
		{"-t1u", args{"-t1u"}, -Duration(1), false},
		{"-t1us", args{"-t1us"}, -Duration(1), false},
		{"-t1μ", args{"-t1μ"}, -Duration(1), false},
		{"-t1μs", args{"-t1μs"}, -Duration(1), false},
		{"-t999999μ", args{"-t999999μ"}, -Duration(999999), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

// MarshalText not tested because it's trivial
func TestDuration_UnmarshalText(t *testing.T) {
	d0 := Duration(0)
	tests := []struct {
		name    string
		t       *Duration
		text    string
		wantErr bool
	}{
		{"nil", nil, "", true},
		{"1234567", new(Duration), "1y2m3dt4h5m6s7us", false},
		{"year", &d0, "1y", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.t.UnmarshalText([]byte(tt.text))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tt.t)
				remarshal := tt.t.String()
				require.Equal(t, tt.text, remarshal)
			}
		})
	}
}

// randomDuration returns a Duration weighted toward short timestamps
func randomDuration() Duration {
	x := 1.0 / (rand.Float64() * 1000)
	return Duration(x*1000000+1) * Millisecond
}

func randomQuantity() Ndau {
	x := 1.0 / (rand.Float64() * 10000)
	n := Ndau(x*100000) * 1000
	n += Ndau(rand.Intn(5) * 100000000)
	return n
}

// TestDuration_UpdateWeightedAverageAge_Fuzz proves that the
// UpdateWeightedAverageAge function is not robust across performing
// calculations in different order for small values. It constructs a sample WAA
// value and updates it twice, with two different quantities, at the same
// timestamp. This is reflective of what happens when a single CreditEAI
// transaction diverts EAI from two different accounts into the same target
// account. Unfortunately, when the amounts are small and the times are small,
// this calculation might be off slightly (so far we've only seen it differ by 1
// microsecond). This causes a hash mismatch if different nodes do the
// calculation in different order -- so we added code to CreditEAI's Apply
// function to sort the list of accounts before iteration.
func TestDuration_UpdateWeightedAverageAge_Fuzz(t *testing.T) {
	ntests := 100
	failureCount := 0
	for i := 0; i < ntests; i++ {
		dur := randomDuration()
		prev := randomQuantity()

		xfer1 := randomQuantity()
		xfer2 := randomQuantity()

		waaA := randomDuration()
		waaB := waaA

		bal := prev
		err := waaA.UpdateWeightedAverageAge(dur, xfer1, bal)
		if err != nil {
			t.Errorf("UpdateWeightedAverageAge(1) returned err: %s", err.Error())
		}
		bal += xfer1
		err = waaA.UpdateWeightedAverageAge(0, xfer2, bal)
		if err != nil {
			t.Errorf("UpdateWeightedAverageAge(2) returned err: %s", err.Error())
		}

		bal = prev
		err = waaB.UpdateWeightedAverageAge(dur, xfer2, bal)
		if err != nil {
			t.Errorf("UpdateWeightedAverageAge(3) returned err: %s", err.Error())
		}
		bal += xfer2
		err = waaB.UpdateWeightedAverageAge(0, xfer1, bal)
		if err != nil {
			t.Errorf("UpdateWeightedAverageAge(4) returned err: %s", err.Error())
		}

		if waaA != waaB {
			failureCount++
			t.Logf("UpdateWeightedAverageAge didn't match (expected):\n %d (%s) != %d (%s)\n dur=%d(%s) prev=%d, xfer1=%d, xfer2=%d\n",
				waaA, waaA, waaB, waaB, dur, dur, prev, xfer1, xfer2)
		}
	}
	if failureCount == 0 || failureCount > ntests/2 {
		t.Errorf("UpdateWeightedAverageAge had a different number of failures than expected -- got %d, should have been 1-%d", failureCount, ntests/2)
	}
}

func TestWAAUpdateCalculation(t *testing.T) {
	priorWAA, err := ParseDuration("t1h52m57s466551us")
	require.NoError(t, err)
	blockTime, err := ParseTimestamp("2019-12-10T20:26:53.194866Z")
	require.NoError(t, err)
	lastWAAUpdate, err := ParseTimestamp("2019-12-10T17:29:15.629235Z")
	require.NoError(t, err)
	balance := Ndau(9305537700000)

	newWAA := priorWAA // copy
	newWAA.UpdateWeightedAverageAge(
		blockTime.Since(lastWAAUpdate),
		0,
		balance,
	)
	// the WAA calculation has to do something
	require.NotEqual(t, priorWAA, newWAA)

	// we can't have a WAA which increases faster than the flow of time, i.e.
	// it must always be true that
	//   (newWAA-priorWAA) <= (blockTime-lastWAAUpdate)
	require.LessOrEqual(t, int64(newWAA-priorWAA), int64(blockTime-lastWAAUpdate))

	// for the transaction which caused this investigation, the new WAA was
	// t4h50m35s32182us. That seems like a lot... but maybe it's accurate?
	expect, err := ParseDuration("t4h50m35s32182us")
	require.Equal(t, expect, newWAA)

	// however, we know that that _can't_ be right, because the account in
	// question was only some 3.5 hours old. Looking at the data, it turns out
	// that there was a bug causing the lastWAAUpdate field to not be updated.
	// If we perform the calculation properly, what comes out?
	realLastWAAUpdate, err := ParseTimestamp("2019-12-10T18:49:51.462752Z")
	require.NoError(t, err)
	acctCreation, err := ParseTimestamp("2019-12-10T16:56:48.877369Z")
	require.NoError(t, err)

	newWAA = priorWAA
	newWAA.UpdateWeightedAverageAge(
		blockTime.Since(realLastWAAUpdate),
		0,
		balance,
	)
	t.Log("real WAA", newWAA)
	require.LessOrEqual(t, int64(acctCreation.Add(newWAA)), int64(blockTime))
}
