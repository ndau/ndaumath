package signature

import (
	"encoding"
	"encoding/base64"
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

// ensure that types here implement msgp marshal types
var _ msgp.Marshaler = (*Key)(nil)
var _ msgp.Unmarshaler = (*Key)(nil)
var _ msgp.Sizer = (*Key)(nil)
var _ msgp.Marshaler = (*PublicKey)(nil)
var _ msgp.Unmarshaler = (*PublicKey)(nil)
var _ msgp.Sizer = (*PublicKey)(nil)
var _ msgp.Marshaler = (*PrivateKey)(nil)
var _ msgp.Unmarshaler = (*PrivateKey)(nil)
var _ msgp.Sizer = (*PrivateKey)(nil)

// ensure that types here implement text encoding interfaces
var _ encoding.TextMarshaler = (*Key)(nil)
var _ encoding.TextUnmarshaler = (*Key)(nil)
var _ encoding.TextMarshaler = (*PublicKey)(nil)
var _ encoding.TextUnmarshaler = (*PublicKey)(nil)
var _ encoding.TextMarshaler = (*PrivateKey)(nil)
var _ encoding.TextUnmarshaler = (*PrivateKey)(nil)

// A Key is a byte slice with known algorithm type
type Key struct {
	algorithm Algorithm
	data      []byte
}

// Marshal marshals the key into a serialized binary format
// which includes a type byte for the algorithm.
func (key Key) Marshal() (serialized []byte, err error) {
	return marshal(key.algorithm, key.data)
}

// Unmarshal unmarshals the serialized binary data into the supplied key instance
func (key *Key) Unmarshal(serialized []byte) error {
	al, kb, err := unmarshal(serialized)
	if err == nil {
		key.algorithm = al
		key.data = kb
	}
	return err
}

// MarshalMsg implements msgp.Marshaler
func (key Key) MarshalMsg(in []byte) (out []byte, err error) {
	out, err = key.Marshal()
	if err == nil {
		out = append(in, out...)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (key *Key) UnmarshalMsg(in []byte) (leftover []byte, err error) {
	var al Algorithm
	var kb []byte
	al, kb, leftover, err = unmarshalWithLeftovers(in)
	if err == nil {
		key.algorithm = al
		key.data = kb
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
// Msgsize implements msgp.Sizer
//
// This method was copy-pasted from the IdentifiedData Msgsize implementation,
// as fundamentally a Key gets serialized as an IdentifiedData, and so should
// have the same size.
func (key *Key) Msgsize() (s int) {
	s = 1 + msgp.Uint8Size + msgp.BytesPrefixSize + len(key.data)
	return
}

// MarshalText implements encoding.TextMarshaler
func (key Key) MarshalText() ([]byte, error) {
	bytes, err := key.Marshal()
	if err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(bytes)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (key *Key) UnmarshalText(text []byte) error {
	bytes, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return err
	}
	return key.Unmarshal(bytes)
}

// Bytes returns the key's data
func (key *Key) Bytes() []byte {
	return key.data
}

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

// MarshalText implements encoding.TextMarshaler
func (key PublicKey) MarshalText() ([]byte, error) {
	return Key(key).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler
func (key *PublicKey) UnmarshalText(text []byte) error {
	err := (*Key)(key).UnmarshalText(text)
	if err == nil {
		if len(key.data) != key.Size() {
			err = fmt.Errorf("Wrong size signature: expect len %d, have %d", key.Size(), len(key.data))
		}
	}
	return err
}

// Bytes returns the key's data
func (key *PublicKey) Bytes() []byte {
	return key.data
}

// A PrivateKey is the public half of a keypair
type PrivateKey Key

// RawPrivateKey creates a PrivateKey from raw data
//
// This is unsafe and subject to only minimal type-checking; it should
// normally be avoided.
func RawPrivateKey(al Algorithm, data []byte) (*PrivateKey, error) {
	pk := PrivateKey{
		algorithm: al,
		data:      data,
	}
	if len(data) != pk.Size() {
		return nil, fmt.Errorf("Wrong private key length")
	}
	return &pk, nil
}

// Size returns the size of this key
func (key PrivateKey) Size() int {
	return key.algorithm.PrivateKeySize()
}

// Sign the supplied message
func (key PrivateKey) Sign(message []byte) Signature {
	al := Key(key).algorithm
	return Signature{
		algorithm: al,
		data:      al.Sign(key.data, message),
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
		if len(key.data) != key.Size() {
			err = fmt.Errorf("Wrong size private key: expect len %d, have %d", key.Size(), len(key.data))
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
		if len(key.data) != key.Size() {
			err = fmt.Errorf("Wrong size signature: expect len %d, have %d", key.Size(), len(key.data))
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
	s = 1 + msgp.Uint8Size + msgp.BytesPrefixSize + len(key.data)
	return
}

// MarshalText implements encoding.TextMarshaler
func (key PrivateKey) MarshalText() ([]byte, error) {
	return Key(key).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler
func (key *PrivateKey) UnmarshalText(text []byte) error {
	err := (*Key)(key).UnmarshalText(text)
	if err == nil {
		if len(key.data) != key.Size() {
			err = fmt.Errorf("Wrong size signature: expect len %d, have %d", key.Size(), len(key.data))
		}
	}
	return err
}

// Bytes returns the key's data
func (key *PrivateKey) Bytes() []byte {
	return key.data
}
