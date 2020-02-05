package null_test

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
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/oneiro-ndev/ndaumath/pkg/signature"
	"github.com/stretchr/testify/require"
)

func TestSizes(t *testing.T) {
	public := signature.PublicKey{}
	private := signature.PrivateKey{}

	require.Equal(t, public.Size(), len(public.KeyBytes()))
	require.Equal(t, private.Size(), len(private.KeyBytes()))
}

func TestEncode(t *testing.T) {
	public := signature.PublicKey{}
	private := signature.PrivateKey{}

	pubTextB, err := public.MarshalText()
	require.NoError(t, err)

	pubText := string(pubTextB)
	require.True(t, strings.HasPrefix(pubText, signature.PublicKeyPrefix))
	require.False(t, strings.HasSuffix(pubText, "="))

	privTextB, err := private.MarshalText()
	require.NoError(t, err)
	require.True(t, utf8.Valid(privTextB), "encoding must be valid utf-8")

	privText := string(privTextB)
	require.True(t, strings.HasPrefix(privText, signature.PrivateKeyPrefix))
	require.False(t, strings.HasSuffix(privText, "="))
}

func TestRoundtripText(t *testing.T) {
	// normally we'd call Generate with the byte literal, but deserialization
	// inserts a pointer to the type instead, and require.Equal doesn't think
	// that a value and a pointer to that value are equal, which causes the
	// roundtrip to fail. It's easier to fix by generating with a pointer to
	// the algorithm here than to change the algorithm serialization code.
	public := signature.PublicKey{}
	private := signature.PrivateKey{}

	pubTextB, err := public.MarshalText()
	require.NoError(t, err)
	privTextB, err := private.MarshalText()
	require.NoError(t, err)
	fmt.Printf("%s %s\n", pubTextB, privTextB)

	rtPub := signature.PublicKey{}
	rtPriv := signature.PrivateKey{}

	err = rtPub.UnmarshalText(pubTextB)
	require.NoError(t, err)
	err = rtPriv.UnmarshalText(privTextB)
	require.NoError(t, err)

	// we have made it so that a Null Algorithm behave the same as a nil Algorithm, so
	// the .algorithm value may not be identical; that's OK.
	require.Equal(t, public.KeyBytes(), rtPub.KeyBytes())
	require.Equal(t, public.ExtraBytes(), rtPub.ExtraBytes())
	require.Equal(t, private.KeyBytes(), rtPriv.KeyBytes())
	require.Equal(t, private.ExtraBytes(), rtPriv.ExtraBytes())
}

func TestRoundtripMsg(t *testing.T) {
	// normally we'd call Generate with the byte literal, but deserialization
	// inserts a pointer to the type instead, and require.Equal doesn't think
	// that a value and a pointer to that value are equal, which causes the
	// roundtrip to fail. It's easier to fix by generating with a pointer to
	// the algorithm here than to change the algorithm serialization code.
	public := signature.PublicKey{}
	private := signature.PrivateKey{}

	pubTextB, err := public.MarshalMsg(nil)
	require.NoError(t, err)
	privTextB, err := private.MarshalMsg(nil)
	require.NoError(t, err)

	rtPub := signature.PublicKey{}
	rtPriv := signature.PrivateKey{}

	leftover, err := rtPub.UnmarshalMsg(pubTextB)
	require.NoError(t, err)
	require.Empty(t, leftover)
	leftover, err = rtPriv.UnmarshalMsg(privTextB)
	require.NoError(t, err)
	require.Empty(t, leftover)

	// we have made it so that a Null Algorithm behave the same as a nil Algorithm, so
	// the .algorithm value may not be identical; that's OK.
	require.Equal(t, public.KeyBytes(), rtPub.KeyBytes())
	require.Equal(t, public.ExtraBytes(), rtPub.ExtraBytes())
	require.Equal(t, private.KeyBytes(), rtPriv.KeyBytes())
	require.Equal(t, private.ExtraBytes(), rtPriv.ExtraBytes())
}

func TestRoundtripBare(t *testing.T) {
	// normally we'd call Generate with the byte literal, but deserialization
	// inserts a pointer to the type instead, and require.Equal doesn't think
	// that a value and a pointer to that value are equal, which causes the
	// roundtrip to fail. It's easier to fix by generating with a pointer to
	// the algorithm here than to change the algorithm serialization code.
	public := signature.PublicKey{}
	private := signature.PrivateKey{}

	pubTextB, err := public.Marshal()
	require.NoError(t, err)
	privTextB, err := private.Marshal()
	require.NoError(t, err)

	rtPub := signature.PublicKey{}
	rtPriv := signature.PrivateKey{}

	err = rtPub.Unmarshal(pubTextB)
	require.NoError(t, err)
	err = rtPriv.Unmarshal(privTextB)
	require.NoError(t, err)

	// we have made it so that a Null Algorithm behave the same as a nil Algorithm, so
	// the .algorithm value may not be identical; that's OK.
	require.Equal(t, public.KeyBytes(), rtPub.KeyBytes())
	require.Equal(t, public.ExtraBytes(), rtPub.ExtraBytes())
	require.Equal(t, private.KeyBytes(), rtPriv.KeyBytes())
	require.Equal(t, private.ExtraBytes(), rtPriv.ExtraBytes())
}

func TestChecksum(t *testing.T) {
	public := signature.PublicKey{}
	private := signature.PrivateKey{}

	pubTextB, err := public.MarshalText()
	require.NoError(t, err)
	privTextB, err := private.MarshalText()
	require.NoError(t, err)

	// offset the bytes so that the failures generated are the result of
	// checksums, not fixed prefix rejection.
	for i := len(signature.PublicKeyPrefix); i < len(pubTextB); i++ {
		t.Run(fmt.Sprintf("public@%d", i), func(t *testing.T) {
			text := make([]byte, len(pubTextB))
			copy(pubTextB, text)
			require.Equal(t, pubTextB, text)
			// edit the byte of the public key at i
			// flip the low bit: that should break the checksum without
			// actually forcing anything out of the ascii range
			text[i] = text[i] ^ 1

			rtPub := signature.PublicKey{}
			err = rtPub.UnmarshalText(text)
			require.Error(t, err)
		})
	}

	// offset the bytes so that the failures generated are the result of
	// checksums, not fixed prefix rejection.
	for i := len(signature.PrivateKeyPrefix); i < len(privTextB); i++ {
		t.Run(fmt.Sprintf("private@%d", i), func(t *testing.T) {
			text := make([]byte, len(privTextB))
			copy(privTextB, text)
			require.Equal(t, privTextB, text)
			// edit the byte of the public key at i
			// flip the low bit: that should break the checksum without
			// actually forcing anything out of the ascii range
			text[i] = text[i] ^ 1

			rtPriv := signature.PrivateKey{}
			err = rtPriv.UnmarshalText(text)
			require.Error(t, err)
		})
	}
}

func TestNullGenerationFails(t *testing.T) {
	_, _, err := signature.Generate(signature.Null, nil)
	require.Error(t, err)
}
