package secp256k1_test

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/oneiro-ndev/ndaumath/pkg/signature"
	"github.com/stretchr/testify/require"
)

func TestSizes(t *testing.T) {
	public, private, err := signature.Generate(signature.Ed25519, nil)
	require.NoError(t, err)

	require.Equal(t, public.Size(), len(public.KeyBytes()))
	require.Equal(t, private.Size(), len(private.KeyBytes()))
}

func TestEncode(t *testing.T) {
	public, private, err := signature.Generate(signature.Secp256k1, nil)
	require.NoError(t, err)

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

func TestRoundtrip(t *testing.T) {
	// normally we'd call Generate with the algorithm literal, but deserialization
	// inserts a pointer to the type instead, and require.Equal doesn't think
	// that a value and a pointer to that value are equal, which causes the
	// roundtrip to fail. It's easier to fix by generating with a pointer to
	// the algorithm here than to change the algorithm serialization code.
	public, private, err := signature.Generate(&signature.Secp256k1, nil)
	require.NoError(t, err)

	pubTextB, err := public.MarshalText()
	require.NoError(t, err)
	privTextB, err := private.MarshalText()
	require.NoError(t, err)

	rtPub := signature.PublicKey{}
	rtPriv := signature.PrivateKey{}

	err = rtPub.UnmarshalText(pubTextB)
	require.NoError(t, err)
	err = rtPriv.UnmarshalText(privTextB)
	require.NoError(t, err)

	require.Equal(t, public, rtPub)
	require.Equal(t, private, rtPriv)
}

func TestRoundtripExtra(t *testing.T) {
	// test that extra data gets de/serialized properly
	// normally we'd call Generate with the algorithm literal, but deserialization
	// inserts a pointer to the type instead, and require.Equal doesn't think
	// that a value and a pointer to that value are equal, which causes the
	// roundtrip to fail. It's easier to fix by generating with a pointer to
	// the algorithm here than to change the algorithm serialization code.
	public0, private0, err := signature.Generate(&signature.Secp256k1, nil)
	require.NoError(t, err)

	extra := make([]byte, 40)
	for i := 1; i < 40; i++ {
		extra[i] = byte(i)
	}

	public, err := signature.RawPublicKey(public0.Algorithm(), public0.KeyBytes(), extra)
	require.NoError(t, err)
	private, err := signature.RawPrivateKey(private0.Algorithm(), private0.KeyBytes(), extra)
	require.NoError(t, err)

	// from here, it's just a matter of running the rest of the normal de/serialization test
	pubTextB, err := public.MarshalText()
	require.NoError(t, err)
	privTextB, err := private.MarshalText()
	require.NoError(t, err)

	rtPub := signature.PublicKey{}
	rtPriv := signature.PrivateKey{}

	err = rtPub.UnmarshalText(pubTextB)
	require.NoError(t, err)
	err = rtPriv.UnmarshalText(privTextB)
	require.NoError(t, err)

	require.Equal(t, *public, rtPub)
	require.Equal(t, *private, rtPriv)

	require.Equal(t, extra, rtPub.ExtraBytes())
	require.Equal(t, extra, rtPriv.ExtraBytes())
}

func TestChecksum(t *testing.T) {
	public, private, err := signature.Generate(signature.Secp256k1, nil)
	require.NoError(t, err)

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
