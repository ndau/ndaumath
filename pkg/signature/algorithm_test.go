package signature

import (
	"crypto/rand"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestAlgorithms(t *testing.T) {
	type key interface {
		Key

		Size() int
		Marshal() ([]byte, error)
		Unmarshal([]byte) error
	}

	clone := func(original key) key {
		val := reflect.ValueOf(original)
		for val.Kind() == reflect.Ptr {
			val = reflect.Indirect(val)
		}
		return reflect.New(val.Type()).Interface().(key)
	}

	type testmarshal func(*testing.T, key) []byte
	type testunmarshal func(*testing.T, []byte, key, bool)

	checkerr := func(t *testing.T, err error, expectErr bool) {
		if expectErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}

	rtPairs := []struct {
		name      string
		marshal   testmarshal
		unmarshal testunmarshal
	}{
		{
			name: "bare",
			marshal: func(t *testing.T, k key) []byte {
				bytes, err := k.Marshal()
				require.NoError(t, err)
				return bytes
			},
			unmarshal: func(t *testing.T, b []byte, k key, expectErr bool) {
				err := k.Unmarshal(b)
				checkerr(t, err, expectErr)
			},
		},
		{
			name: "text",
			marshal: func(t *testing.T, k key) []byte {
				tbytes, err := k.MarshalText()
				require.NoError(t, err)
				require.True(
					t, utf8.Valid(tbytes),
					"text marshalling must produce valid utf-8",
				)

				s := string(tbytes)
				require.True(
					t,
					strings.HasPrefix(s, PublicKeyPrefix) || strings.HasPrefix(s, PrivateKeyPrefix),
					"text marshalling must have an appropriate prefix",
				)
				require.False(
					t, strings.HasSuffix(s, "="),
					"text marshalling must be accomplished without padding",
				)

				return tbytes
			},
			unmarshal: func(t *testing.T, tb []byte, k key, expectErr bool) {
				err := k.UnmarshalText(tb)
				checkerr(t, err, expectErr)
			},
		},
		{
			name: "msgp",
			marshal: func(t *testing.T, k key) []byte {
				mbytes, err := k.MarshalMsg(nil)
				require.NoError(t, err)
				return mbytes
			},
			unmarshal: func(t *testing.T, m []byte, k key, expectErr bool) {
				leftover, err := k.UnmarshalMsg(m)
				checkerr(t, err, expectErr)
				require.Empty(
					t, leftover,
					"msgp unmarshalling must consume all bytes produced by msgp marshaller",
				)
			},
		},
	}

	algorithms := []Algorithm{&Ed25519, &Secp256k1}

	for _, algorithm := range algorithms {
		t.Run(NameOf(algorithm), func(t *testing.T) {
			public, private, err := Generate(algorithm, nil)
			require.NoError(t, err)

			pubBuf := make([]byte, 128)
			_, err = io.ReadFull(rand.Reader, pubBuf)
			require.NoError(t, err)
			pubX, err := RawPublicKey(algorithm, public.KeyBytes(), pubBuf)
			require.NoError(t, err)

			privBuf := make([]byte, 128)
			_, err = io.ReadFull(rand.Reader, privBuf)
			require.NoError(t, err)
			privX, err := RawPrivateKey(algorithm, private.KeyBytes(), privBuf)
			require.NoError(t, err)

			for _, kt := range []struct {
				name string
				k    key
			}{
				{"public", &public},
				{"private", &private},
				{"public with extra", pubX},
				{"private with extra", privX},
			} {
				t.Run(kt.name, func(t *testing.T) {
					t.Run("size", func(t *testing.T) {
						require.Equal(t, kt.k.Size(), len(kt.k.KeyBytes()))
					})

					getEmpty := func(t *testing.T, k key) key {
						// set up an empty container of the same publicity
						empty := clone(kt.k)
						empty.Zeroize()
						require.Empty(t, empty.Algorithm())
						require.Empty(t, empty.KeyBytes())
						require.Empty(t, empty.ExtraBytes())
						return empty
					}

					for _, rt := range rtPairs {
						t.Run(rt.name, func(t *testing.T) {
							// ensure marshalling works
							bytes := rt.marshal(t, kt.k)

							t.Log("size of marshalled form", len(bytes))

							t.Run("unmarshal", func(t *testing.T) {
								empty := getEmpty(t, kt.k)
								// ensure that we can unmarshal into the empty container
								rt.unmarshal(t, bytes, empty, false)

								// ensure than unmarshalling preserved our data
								require.Equal(t, kt.k.Algorithm(), empty.Algorithm())
								require.Equal(t, kt.k.KeyBytes(), empty.KeyBytes())
								require.Equal(t, kt.k.ExtraBytes(), empty.ExtraBytes())
							})

							t.Run("checksum", func(t *testing.T) {
								for i := 0; i < len(bytes); i++ {
									t.Run(fmt.Sprintf("err@%d", i), func(t *testing.T) {
										newbytes := make([]byte, len(bytes))
										ncopied := copy(bytes, newbytes)
										require.Equal(t, len(bytes), ncopied)
										require.Equal(t, bytes, newbytes)

										// invert byte i
										newbytes[i] = newbytes[i] ^ 1

										empty := getEmpty(t, kt.k)
										rt.unmarshal(t, newbytes, empty, true)
									})
								}
							})
						})
					}
				})
			}
		})
	}
}

func TestPublic(t *testing.T) {
	algorithms := []Algorithm{&Ed25519, &Secp256k1}

	for _, algorithm := range algorithms {
		t.Run(NameOf(algorithm)+"_Public", func(t *testing.T) {
			// verify that we can create a keypair and then generate the equivalent
			// public key from the private key
			public, private, err := Generate(algorithm, nil)
			require.NoError(t, err)

			pubFromPrivate := algorithm.Public(private.KeyBytes())
			require.Equal(t, public.KeyBytes(), pubFromPrivate)
		})
	}

}
