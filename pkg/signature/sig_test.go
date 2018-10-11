package signature

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignRoundtripEd25519(t *testing.T) {
	message := make([]byte, 256)
	rand.Read(message)

	public, private, err := Generate(Ed25519, nil)
	require.NoError(t, err)

	signature := private.Sign(message)
	t.Logf("1 data:        %x", signature.data)

	signatureBytes, err := signature.Marshal()
	require.NoError(t, err)
	t.Logf("bytes: %x", signatureBytes)

	signature2 := Signature{}
	err = (&signature2).Unmarshal(signatureBytes)
	require.NoError(t, err)

	t.Logf("2 data:        %x", signature2.data)

	// method 1 to verify a signature
	require.True(t, signature2.Verify(message, public))
	// method 2 to verify a signature
	require.True(t, public.Verify(message, signature2))
}
