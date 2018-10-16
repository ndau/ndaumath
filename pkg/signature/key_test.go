package signature

import (
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
