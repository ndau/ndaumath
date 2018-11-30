package signature

import (
	"fmt"

	"github.com/oneiro-ndev/ndaumath/pkg/b32"
	"github.com/pkg/errors"
	"github.com/tinylib/msgp/msgp"
)

// ensure that types here implement msgp marshal types
var _ msgp.Marshaler = (*Signature)(nil)
var _ msgp.Unmarshaler = (*Signature)(nil)
var _ msgp.Sizer = (*Signature)(nil)

// A Signature is a byte slice with known algorithm type
type Signature struct {
	algorithm Algorithm
	data      []byte
}

// Algorithm gets the signature's algorithm
func (signature Signature) Algorithm() Algorithm {
	return signature.algorithm
}

// RawSignature creates a Signature from raw data
//
// This is unsafe and subject to only minimal type-checking; it should
// normally be avoided.
func RawSignature(al Algorithm, data []byte) (*Signature, error) {
	sig := Signature{
		algorithm: al,
		data:      data,
	}
	exsize := sig.Size()
	if exsize >= 0 && len(data) != exsize {
		return nil, fmt.Errorf("wrong signature length: have %d, want %d", len(data), sig.Size())
	}
	return &sig, nil
}

// Size returns the size of this signature
func (signature Signature) Size() int {
	return signature.algorithm.SignatureSize()
}

// Marshal marshals the signature into a serialized binary format
// which includes a type byte for the algorithm.
func (signature Signature) Marshal() (serialized []byte, err error) {
	return marshal(signature.algorithm, signature.data)
}

// Unmarshal unmarshals the serialized binary data into the supplied signature instance
func (signature *Signature) Unmarshal(serialized []byte) error {
	al, b, err := unmarshal(serialized)
	if err == nil && len(b) != al.SignatureSize() {
		err = fmt.Errorf("Wrong size signature: expect len %d, have %d", al.SignatureSize(), len(b))
	}
	if err == nil {
		signature.algorithm = al
		signature.data = b
	}
	return err
}

// Verify is a convenience function to verify from a signature
func (signature Signature) Verify(message []byte, key PublicKey) bool {
	return key.Verify(message, signature)
}

// MarshalMsg implements msgp.Marshaler
func (signature Signature) MarshalMsg(in []byte) (out []byte, err error) {
	out, err = signature.Marshal()
	if err == nil {
		out = append(in, out...)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (signature *Signature) UnmarshalMsg(in []byte) (leftover []byte, err error) {
	var al Algorithm
	var b []byte
	al, b, leftover, err = unmarshalWithLeftovers(in)
	if err == nil {
		signature.algorithm = al
		signature.data = b
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
// Msgsize implements msgp.Sizer
//
// This method was copy-pasted from the IdentifiedData Msgsize implementation,
// as fundamentally a Signature gets serialized as an IdentifiedData, and so should
// have the same size.
func (signature *Signature) Msgsize() (s int) {
	s = 1 + msgp.Uint8Size + msgp.BytesPrefixSize + len(signature.data)
	return
}

// Bytes returns the key's data
func (signature *Signature) Bytes() []byte {
	return signature.data
}

// MarshalText implements encoding.TextMarshaler
//
// This marshaller uses a custom b32 encoding which is case-insensitive and
// lacks certain confusing pairs, for ease of human-friendly handling.
// For the same reason, it embeds a checksum, so it's easy to tell whether
// or not it was received correctly.
func (signature Signature) MarshalText() ([]byte, error) {
	bytes, err := signature.Marshal()
	if err != nil {
		return nil, err
	}
	bytes = AddChecksum(bytes)
	return []byte(b32.Encode(bytes)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (signature *Signature) UnmarshalText(text []byte) error {
	bytes, err := b32.Decode(string(text))
	if err != nil {
		return err
	}
	var checksumOk bool
	bytes, checksumOk = CheckChecksum(bytes)
	if !checksumOk {
		return errors.New("key unmarshal failure: bad checksum")
	}
	return signature.Unmarshal(bytes)
}

// MarshalString is like MarshalText, but to a string
func (signature *Signature) MarshalString() (string, error) {
	// Why doesn't MarshalText produce a string anyway?
	t, err := signature.MarshalText()
	if t == nil {
		t = []byte{}
	}
	return string(t), err
}

// ParseSignature parses a string representation of a signature, if possible
func ParseSignature(s string) (*Signature, error) {
	key := new(Signature)
	err := key.UnmarshalText([]byte(s))
	return key, err
}
