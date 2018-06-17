package b32

import (
	"encoding/base32"
	"strings"
)

// NdauAlphabet is the encoding alphabet we use for byte32 encoding
// It consists of the lowercase alphabet and digits, without l, 1, 0, and o.
// When decoding, we will accept either upper or lower case.
const NdauAlphabet = "abcdefghijkmnpqrstuvwxyz23456789"

// Index looks up the value of a letter in the ndau encoding alphabet.
func Index(c string) int {
	return strings.Index(NdauAlphabet, c)
}

// Encode converts a byte stream into a base32 string
func Encode(b []byte) string {
	enc := base32.NewEncoding(NdauAlphabet)
	r := enc.EncodeToString(b)
	return r
}

// Decode converts a string back to a byte stream
func Decode(s string) ([]byte, error) {
	enc := base32.NewEncoding(NdauAlphabet)
	return enc.DecodeString(s)
}
