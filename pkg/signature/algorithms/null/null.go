package null

import (
	"io"
)

// This implements algorithm for the Null signature type; this type
// exists to allow serializing and deserializing empty signatures without error

// Never edit this; it would be a const if go were smarter
var Null = null{}

type null struct{}

// PublicKeySize implements Algorithm
func (null) PublicKeySize() int {
	return 0
}

// PrivateKeySize implements Algorithm
func (null) PrivateKeySize() int {
	return 0
}

// SignatureSize implements Algorithm
func (null) SignatureSize() int {
	return 0
}

// Generate implements Algorithm
func (null) Generate(rand io.Reader) (public, private []byte, err error) {
	return []byte{}, []byte{}, nil
}

// Sign implements Algorithm
func (n null) Sign(private, message []byte) []byte {
	return []byte{}
}

// Verify implements Algorithm
func (null) Verify(public, message, sig []byte) bool {
	return false
}
