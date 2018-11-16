package signature

import (
	"bytes"
	"encoding"
	"fmt"
	"strings"

	"github.com/tinylib/msgp/msgp"
)

// PublicKeyPrefix always prefixes Ndau public keys in text serialization
const PublicKeyPrefix = "npub"

// MaybePublic provides a fast way to check whether a string looks like
// it might be an ndau public key.
//
// To get a definitive answer as to whether something is a public key, one
// must attempt to deserialize it using UnmarshalText and check the error
// value. That takes some work; it's faster to use this to get a first impression.
//
// This function will allow some false positives, but no false negatives:
// some values for which it returns `true` may not be actual valid keys,
// but no values for which it returns `false` will return actual valid keys.
func MaybePublic(s string) bool {
	return strings.HasPrefix(s, PublicKeyPrefix)
}

// ensure that PublicKey implements msgp marshal types
var _ msgp.Marshaler = (*PublicKey)(nil)
var _ msgp.Unmarshaler = (*PublicKey)(nil)
var _ msgp.Sizer = (*PublicKey)(nil)

// ensure that PublicKey implements text encoding interfaces
var _ encoding.TextMarshaler = (*PublicKey)(nil)
var _ encoding.TextUnmarshaler = (*PublicKey)(nil)

// ensure that PublicKey implements string shorthand interfaces
var _ fmt.Stringer = (*PublicKey)(nil)

// ensure that PublicKey implements export interfaces
var _ Key = (*PublicKey)(nil)

// A PublicKey is the public half of a keypair
type PublicKey keyBase

// RawPublicKey creates a PublicKey from raw data
//
// This is unsafe and subject to only minimal type-checking; it should
// normally be avoided.
func RawPublicKey(al Algorithm, key, extra []byte) (*PublicKey, error) {
	if key == nil {
		key = []byte{}
	}
	if extra == nil {
		extra = []byte{}
	}
	pk := PublicKey{
		algorithm: al,
		key:       key,
		extra:     extra,
	}
	if len(key) != pk.Size() {
		return nil, fmt.Errorf("wrong public key length: have %d, want %d", len(key), pk.Size())
	}
	return &pk, nil
}

// Size returns the size of this key
func (key PublicKey) Size() int {
	return key.Algorithm().PublicKeySize()
}

// Verify the supplied message with the given signature
func (key PublicKey) Verify(message []byte, sig Signature) bool {
	if NameOf(key.Algorithm()) != NameOf(sig.algorithm) {
		return false
	}
	return key.Algorithm().Verify(keyBase(key).key, message, sig.data)
}

// Marshal marshals the PublicKey into a serialized binary format
func (key PublicKey) Marshal() ([]byte, error) {
	return keyBase(key).Marshal()
}

// Unmarshal unmarshals the serialized bytes into the PublicKey pointer
func (key *PublicKey) Unmarshal(serialized []byte) error {
	err := (*keyBase)(key).Unmarshal(serialized)
	if err == nil {
		if len(key.key) != key.Size() {
			err = fmt.Errorf("Wrong size public key: expect len %d, have %d", key.Size(), len(key.key))
		}
	}
	return err
}

// MarshalMsg implements msgp.Marshaler
func (key PublicKey) MarshalMsg(in []byte) (out []byte, err error) {
	return keyBase(key).MarshalMsg(in)
}

// UnmarshalMsg implements msgp.Unmarshaler
func (key *PublicKey) UnmarshalMsg(in []byte) (leftover []byte, err error) {
	leftover, err = (*keyBase)(key).UnmarshalMsg(in)
	if err == nil {
		if len(key.key) != key.Size() {
			err = fmt.Errorf("Wrong size public key: expect len %d, have %d", key.Size(), len(key.key))
		}
	}
	return leftover, err
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
// Msgsize implements msgp.Sizer
//
// This method was copy-pasted from the IdentifiedData Msgsize implementation,
// as fundamentally a PublicKey gets serialized as an IdentifiedData, and so should
// have the same size.
func (key *PublicKey) Msgsize() (s int) {
	s = 1 + msgp.Uint8Size + msgp.BytesPrefixSize + len(key.key)
	return
}

// MarshalText implements encoding.TextMarshaler.
//
// PublicKeys encode like Keys, with the addition of a human-readable prefix
// for easy identification.
func (key PublicKey) MarshalText() ([]byte, error) {
	bytes, err := keyBase(key).MarshalText()
	bytes = append([]byte(PublicKeyPrefix), bytes...)
	return bytes, err
}

// UnmarshalText implements encoding.TextUnmarshaler
func (key *PublicKey) UnmarshalText(text []byte) error {
	expectPrefix := []byte(PublicKeyPrefix)
	lep := len(expectPrefix)
	if !bytes.Equal(expectPrefix, text[:lep]) {
		return fmt.Errorf("public key must begin with %q; got %q", PublicKeyPrefix, text[:lep])
	}
	err := (*keyBase)(key).UnmarshalText(text[lep:])
	if err == nil {
		if len(key.key) != key.Size() {
			err = fmt.Errorf("Wrong size public key: expect len %d, have %d", key.Size(), len(key.key))
		}
	}
	return err
}

// KeyBytes returns the key's data
func (key PublicKey) KeyBytes() []byte {
	return keyBase(key).KeyBytes()
}

// ExtraBytes returns the key's extra data
func (key PublicKey) ExtraBytes() []byte {
	return keyBase(key).ExtraBytes()
}

// Algorithm returns the key's algorithm
func (key PublicKey) Algorithm() Algorithm {
	return keyBase(key).Algorithm()
}

// String returns a shorthand for the key's data
//
// This returns the first 8 characters of the text serialization,
// an ellipsis, then the final 4 characters of the text serialization.
// Total output size is constant at 15 characters.
//
// This destructively truncates the key, but it is a useful format for
// humans.
func (key PublicKey) String() string {
	return keyBase(key).String(PublicKeyPrefix)
}

// Truncate removes all extra data from this key.
//
// This is a destructive operation which cannot be undone; make copies
// first if you need to.
func (key *PublicKey) Truncate() {
	key.extra = nil
}

// Zeroize removes all data from this key
//
// This is a destructive operation which cannot be undone; make copies
// first if you need to.
func (key *PublicKey) Zeroize() {
	if key == nil {
		return
	}
	kkey := keyBase(*key)
	kkey.Zeroize()
	*key = PublicKey(kkey)
}
