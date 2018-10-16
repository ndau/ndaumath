package signature

import (
	"bytes"
	"encoding"
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

// ensure that PublicKey implements msgp marshal types
var _ msgp.Marshaler = (*PublicKey)(nil)
var _ msgp.Unmarshaler = (*PublicKey)(nil)
var _ msgp.Sizer = (*PublicKey)(nil)

// ensure that PublicKey implements text encoding interfaces
var _ encoding.TextMarshaler = (*PublicKey)(nil)
var _ encoding.TextUnmarshaler = (*PublicKey)(nil)

// ensure that PublicKey implements string shorthand interfaces
var _ fmt.Stringer = (*PublicKey)(nil)

// ensure that PublicKey implements byte export interfaces
var _ byteser = (*PublicKey)(nil)

// A PublicKey is the public half of a keypair
type PublicKey Key

// RawPublicKey creates a PublicKey from raw data
//
// This is unsafe and subject to only minimal type-checking; it should
// normally be avoided.
func RawPublicKey(al Algorithm, data []byte) (*PublicKey, error) {
	pk := PublicKey{
		algorithm: al,
		data:      data,
	}
	if len(data) != pk.Size() {
		return nil, fmt.Errorf("Wrong public key length")
	}
	return &pk, nil
}

// Size returns the size of this key
func (key PublicKey) Size() int {
	return key.algorithm.PublicKeySize()
}

// Verify the supplied message with the given signature
func (key PublicKey) Verify(message []byte, sig Signature) bool {
	if nameOf(key.algorithm) != nameOf(sig.algorithm) {
		return false
	}
	return key.algorithm.Verify(Key(key).data, message, sig.data)
}

// Marshal marshals the PublicKey into a serialized binary format
func (key PublicKey) Marshal() ([]byte, error) {
	return Key(key).Marshal()
}

// Unmarshal unmarshals the serialized bytes into the PublicKey pointer
func (key *PublicKey) Unmarshal(serialized []byte) error {
	err := (*Key)(key).Unmarshal(serialized)
	if err == nil {
		if len(key.data) != key.Size() {
			err = fmt.Errorf("Wrong size public key: expect len %d, have %d", key.Size(), len(key.data))
		}
	}
	return err
}

// MarshalMsg implements msgp.Marshaler
func (key PublicKey) MarshalMsg(in []byte) (out []byte, err error) {
	return Key(key).MarshalMsg(in)
}

// UnmarshalMsg implements msgp.Unmarshaler
func (key *PublicKey) UnmarshalMsg(in []byte) (leftover []byte, err error) {
	leftover, err = (*Key)(key).UnmarshalMsg(in)
	if err == nil {
		if len(key.data) != key.Size() {
			err = fmt.Errorf("Wrong size public key: expect len %d, have %d", key.Size(), len(key.data))
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
	s = 1 + msgp.Uint8Size + msgp.BytesPrefixSize + len(key.data)
	return
}

// MarshalText implements encoding.TextMarshaler.
//
// PublicKeys encode like Keys, with the addition of a human-readable prefix
// for easy identification.
func (key PublicKey) MarshalText() ([]byte, error) {
	bytes, err := Key(key).MarshalText()
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
	err := (*Key)(key).UnmarshalText(text[lep:])
	if err == nil {
		if len(key.data) != key.Size() {
			err = fmt.Errorf("Wrong size public key: expect len %d, have %d", key.Size(), len(key.data))
		}
	}
	return err
}

// Bytes returns the key's data
func (key PublicKey) Bytes() []byte {
	return key.data
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
	return Key(key).String()
}
