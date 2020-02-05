package address

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func getKinds() []byte {
	return []byte{
		KindUser,
		KindNdau,
		KindExchange,
		KindEndowment,
	}
}

func TestArbitraryAddressesAreValid(t *testing.T) {
	kinds := getKinds()
	for i := 0; i < 16; i++ {
		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)

		t.Run(string(i), func(t *testing.T) {
			t.Log("key", fmt.Sprintf("%x", key))
			address, err := Generate(kinds[i&3], key)
			require.NoError(t, err)
			t.Log("address", address)
			address, err = Validate(address.addr)
			t.Log("address", address)
			require.NoError(t, err)
		})
	}
}

func TestArbitraryAddressesDoRoundtrips(t *testing.T) {
	kinds := getKinds()
	for i := 0; i < 16; i++ {
		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)

		t.Run(string(i), func(t *testing.T) {
			t.Log("key", fmt.Sprintf("%x", key))
			address1, err := Generate(kinds[i&3], key)
			require.NoError(t, err)
			t.Log("address1", address1)
			b, err := address1.MarshalMsg(nil)
			require.NoError(t, err)
			address2 := Address{}
			extra, err := address2.UnmarshalMsg(b)
			require.NoError(t, err)
			require.Empty(t, extra)
			t.Log("address2", address2)
			address3, err := Validate(address2.addr)
			t.Log("address3", address3)
			require.NoError(t, err)
		})
	}
}

func TestKnownKeyGeneratesKnownValue(t *testing.T) {
	key := make([]byte, 16)
	for i := byte(0); i < 16; i++ {
		key[i] = i
	}

	address, err := Generate(KindUser, key)
	require.NoError(t, err)
	require.Equal(t, "ndadprx764ciigti8d8whtw2kct733r85qvjukhqhke3dka4", address.String())
}

func TestKnownKeyValidates(t *testing.T) {
	_, err := Validate("ndadprx764ciigti8d8whtw2kct733r85qvjukhqhke3dka4")
	require.NoError(t, err)
	// fail with a minor change
	_, err = Validate("ndxdprx764ciigti8d8whtw2kct733r85qvjukhqhke3dka4")
	require.Error(t, err)
}

func BenchmarkGeneration(b *testing.B) {
	key := make([]byte, 32)
	kinds := getKinds()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		_, err := rand.Read(key)
		if err != nil {
			b.FailNow()
		}
		b.StartTimer()

		_, err = Generate(kinds[n&3], key)
		if err != nil {
			b.FailNow()
		}
	}
}
