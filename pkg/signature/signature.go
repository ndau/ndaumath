package signature

import (
	"fmt"

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

// RawSignature creates a Signature from raw data
//
// This is unsafe and subject to only minimal type-checking; it should
// normally be avoided.
func RawSignature(al Algorithm, data []byte) (*Signature, error) {
	sig := Signature{
		algorithm: al,
		data:      data,
	}
	if len(data) != sig.Size() {
		return nil, fmt.Errorf("Wrong signature length")
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
