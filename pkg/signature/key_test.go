package signature

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Generate high-level keys
func TestGenerateEd25519(t *testing.T) {
	public, private, err := Generate(Ed25519, nil)
	require.NoError(t, err)
	var _ PublicKey = public
	var _ PrivateKey = private
}

func TestGenerateEd25519Consistency(t *testing.T) {
	randbuf := make([]byte, 32)
	_, err := rand.Read(randbuf)
	require.NoError(t, err)
	firstpub, firstprivate, err := Generate(Ed25519, bytes.NewReader(randbuf))
	require.NoError(t, err)
	for i := 1; i <= 10; i++ {
		t.Run(fmt.Sprintf("Pass %d", i), func(t *testing.T) {
			public, private, err := Generate(Ed25519, bytes.NewReader(randbuf))
			require.NoError(t, err)
			require.Equal(t, firstpub, public)
			require.Equal(t, firstprivate, private)
		})
	}
}

func TestGenerateSecp256k1Consistency(t *testing.T) {
	randbuf := make([]byte, 32)
	_, err := rand.Read(randbuf)
	require.NoError(t, err)
	firstpub, firstprivate, err := Generate(Secp256k1, bytes.NewReader(randbuf))
	require.NoError(t, err)
	for i := 1; i <= 10; i++ {
		t.Run(fmt.Sprintf("Pass %d", i), func(t *testing.T) {
			public, private, err := Generate(Secp256k1, bytes.NewReader(randbuf))
			require.NoError(t, err)
			require.Equal(t, firstpub, public)
			require.Equal(t, firstprivate, private)
		})
	}
}

func TestRoundtripEd25519(t *testing.T) {
	public, private, err := Generate(Ed25519, nil)
	require.NoError(t, err)
	publicBytes, err := public.Marshal()
	require.NoError(t, err)
	privateBytes, err := private.Marshal()
	require.NoError(t, err)

	// method 1 to deserialize a key without needing to know in advance
	// which algorithm it uses:
	public2 := new(PublicKey)
	err = public2.Unmarshal(publicBytes)
	require.NoError(t, err)
	require.Equal(t, public.key, public2.key)
	require.Equal(t, NameOf(public.algorithm), NameOf(public2.algorithm))

	// method 2 to deserialize a key without needing to know in advance
	// which algorithm it uses
	private2 := PrivateKey{}
	err = (&private2).Unmarshal(privateBytes)
	require.NoError(t, err)
	require.Equal(t, private.key, private2.key)
	require.Equal(t, NameOf(private.algorithm), NameOf(private2.algorithm))
}

// These examples are for documentation purposes to help show how address generation operates internally
func ExampleMarshalEd25519Public(t *testing.T) {
	public, private, err := Generate(Ed25519, nil)
	require.NoError(t, err)
	var _ PublicKey = public
	var _ PrivateKey = private
	keydata, _ := public.pack()
	b, _ := public.Marshal()
	fmt.Printf(" raw key: `%x` (len %d)\nunpacked: `%x` (len %d)\n", keydata, len(keydata), b, len(b))
	tx, _ := public.MarshalText()
	fmt.Printf("    user: `%s`\n", string(tx))
}

func ExampleMarshalEd25519Private(t *testing.T) {
	public, private, err := Generate(Ed25519, nil)
	require.NoError(t, err)
	var _ PublicKey = public
	var _ PrivateKey = private
	keydata, _ := private.pack()
	b, _ := private.Marshal()
	fmt.Printf(" raw key: `%x` (len %d)\nunpacked: `%x` (len %d)\n", keydata, len(keydata), b, len(b))
	tx, _ := private.MarshalText()
	fmt.Printf("    user: `%s`\n", string(tx))
}

func ExampleMarshalSecp256k1Public(t *testing.T) {
	public, private, err := Generate(Secp256k1, nil)
	require.NoError(t, err)
	var _ PublicKey = public
	var _ PrivateKey = private
	keydata, _ := public.pack()
	b, _ := public.Marshal()
	fmt.Printf(" raw key: `%x` (len %d)\nunpacked: `%x` (len %d)\n", keydata, len(keydata), b, len(b))
	tx, _ := public.MarshalText()
	fmt.Printf("    user: `%s`\n", string(tx))
}

func ExampleMarshalSecp256k1Private(t *testing.T) {
	public, private, err := Generate(Secp256k1, nil)
	require.NoError(t, err)
	var _ PublicKey = public
	var _ PrivateKey = private
	keydata, _ := private.pack()
	b, _ := private.Marshal()
	fmt.Printf(" raw key: `%x` (len %d)\nunpacked: `%x` (len %d)\n", keydata, len(keydata), b, len(b))
	tx, _ := private.MarshalText()
	fmt.Printf("    user: `%s`\n", string(tx))
}
