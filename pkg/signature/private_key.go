package signature

import (
	"bytes"
	"encoding"
	"fmt"
	"strings"

	"github.com/tinylib/msgp/msgp"
)

// PrivateKeyPrefix always prefixes Ndau private keys in text serialization
const PrivateKeyPrefix = "npvt"

// MaybePrivate provides a fast way to check whether a string looks like
// it might be an ndau private key.
//
// To get a definitive answer as to whether something is a private key, one
// must attempt to deserialize it using UnmarshalText and check the error
// value. That takes some work; it's faster to use this to get a first impression.
//
// This function will allow some false positives, but no false negatives:
// some values for which it returns `true` may not be actual valid keys,
// but no values for which it returns `false` will return actual valid keys.
func MaybePrivate(s string) bool {
	return strings.HasPrefix(s, PrivateKeyPrefix)
}

// ensure that PrivateKey implements msgp marshal types
var _ msgp.Marshaler = (*PrivateKey)(nil)
var _ msgp.Unmarshaler = (*PrivateKey)(nil)
var _ msgp.Sizer = (*PrivateKey)(nil)

// ensure that PrivateKey implements text encoding interfaces
var _ encoding.TextMarshaler = (*PrivateKey)(nil)
var _ encoding.TextUnmarshaler = (*PrivateKey)(nil)

// ensure that PrivateKey implements string shorthand interfaces
var _ fmt.Stringer = (*PrivateKey)(nil)

// ensure that PrivateKey implements export interfaces
var _ keyer = (*PrivateKey)(nil)

// A PrivateKey is the private half of a keypair
type PrivateKey Key

// RawPrivateKey creates a PrivateKey from raw data
//
// This is unsafe and subject to only minimal type-checking; it should
// normally be avoided.
func RawPrivateKey(al Algorithm, key, extra []byte) (*PrivateKey, error) {
	pk := PrivateKey{
		algorithm: al,
		key:       key,
		extra:     extra,
	}
	if len(key) != pk.Size() {
		return nil, fmt.Errorf("Wrong private key length")
	}
	return &pk, nil
}

// Size returns the size of this key
func (key PrivateKey) Size() int {
	return key.Algorithm().PrivateKeySize()
}

// Sign the supplied message
func (key PrivateKey) Sign(message []byte) Signature {
	al := key.Algorithm()
	return Signature{
		algorithm: al,
		data:      al.Sign(key.key, message),
	}
}

// Marshal marshals the PrivateKey into a serialized binary format
func (key PrivateKey) Marshal() ([]byte, error) {
	return Key(key).Marshal()
}

// Unmarshal unmarshals the serialized bytes into the PrivateKey pointer
func (key *PrivateKey) Unmarshal(serialized []byte) error {
	err := (*Key)(key).Unmarshal(serialized)
	if err == nil {
		if len(key.key) != key.Size() {
			err = fmt.Errorf("Wrong size private key: expect len %d, have %d", key.Size(), len(key.key))
		}
	}
	return err
}

// MarshalMsg implements msgp.Marshaler
func (key PrivateKey) MarshalMsg(in []byte) (out []byte, err error) {
	return Key(key).MarshalMsg(in)
}

// UnmarshalMsg implements msgp.Unmarshaler
func (key *PrivateKey) UnmarshalMsg(in []byte) (leftover []byte, err error) {
	leftover, err = (*Key)(key).UnmarshalMsg(in)
	if err == nil {
		if len(key.key) != key.Size() {
			err = fmt.Errorf("Wrong size signature: expect len %d, have %d", key.Size(), len(key.key))
		}
	}
	return leftover, err
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
// Msgsize implements msgp.Sizer
//
// This method was copy-pasted from the IdentifiedData Msgsize implementation,
// as fundamentally a PrivateKey gets serialized as an IdentifiedData, and so should
// have the same size.
func (key *PrivateKey) Msgsize() (s int) {
	s = 1 + msgp.Uint8Size + msgp.BytesPrefixSize + len(key.key)
	return
}

// MarshalText implements encoding.TextMarshaler.
//
// PublicKeys encode like Keys, with the addition of a human-readable prefix
// for easy identification.
func (key PrivateKey) MarshalText() ([]byte, error) {
	bytes, err := Key(key).MarshalText()
	bytes = append([]byte(PrivateKeyPrefix), bytes...)
	return bytes, err
}

// UnmarshalText implements encoding.TextUnmarshaler
func (key *PrivateKey) UnmarshalText(text []byte) error {
	expectPrefix := []byte(PrivateKeyPrefix)
	lep := len(expectPrefix)
	if !bytes.Equal(expectPrefix, text[:lep]) {
		return fmt.Errorf("public key must begin with %q; got %q", PublicKeyPrefix, text[:lep])
	}
	err := (*Key)(key).UnmarshalText(text[lep:])
	if err == nil {
		if len(key.key) != key.Size() {
			err = fmt.Errorf("Wrong size signature: expect len %d, have %d", key.Size(), len(key.key))
		}
	}
	return err
}

// KeyBytes returns the key's data
func (key PrivateKey) KeyBytes() []byte {
	return Key(key).KeyBytes()
}

// ExtraBytes returns the key's extra data
func (key PrivateKey) ExtraBytes() []byte {
	return Key(key).ExtraBytes()
}

// Algorithm returns the key's algorithm
func (key PrivateKey) Algorithm() Algorithm {
	return Key(key).Algorithm()
}

// String returns a shorthand for the key's data
//
// This returns the first 8 characters of the text serialization,
// an ellipsis, then the final 4 characters of the text serialization.
// Total output size is constant at 15 characters.
//
// This destructively truncates the key, but it is a useful format for
// humans.
func (key PrivateKey) String() string {
	return Key(key).String()
}

// Truncate removes all extra data from this key.
//
// This is a destructive operation which cannot be undone; make copies
// first if you need to.
func (key *PrivateKey) Truncate() {
	key.extra = nil
}
