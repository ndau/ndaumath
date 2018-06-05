package address

import (
	"crypto/sha256"
	"encoding/base32"
	"hash/adler32"
	"strings"
)

// An ndau address is the result of a mathematical process over a public key. It
// is a byte32 encoding, using a custom alphabet, of a portion of the SHA256
// hash of the key, concatenated with some additional marker and checksum
// information. The result is a key that always starts with the letters 'nd' and
// one more character that specifies the type of address.

// NdauAlphabet is the encoding alphabet we use for byte32 encoding
// It consists of the lowercase alphabet and digits, without l, 1, 0, and o.
// When decoding, we will accept either upper or lower case.
const NdauAlphabet = "abcdefghijkmnpqrstuvwxyz23456789"

// ndx looks up the value of a letter in the alphabet.
func ndx(c string) int {
	return strings.Index(NdauAlphabet, c)
}

// Kind indicates the type of address in use; this is an external indication
// designed to help users evaluate their own actions; it may or may not be
// enforced by the blockchain.
type Kind string

// Error is the type of all the errors this package returns
type Error struct {
	msg string
}

func (a *Error) Error() string {
	return "address error: " + a.msg
}

func newError(msg string) error {
	return &Error{msg}
}

// we want the first letters of the address to be "nd?" where ? is the kind of
// address. Valid address types are as follows:
const (
	KindUser      Kind = "a"
	KindNdau      Kind = "n"
	KindEndowment Kind = "e"
	KindExchange  Kind = "x"
)

func IsValidKind(a Kind) bool {
	switch a {
	case KindUser,
		KindNdau,
		KindEndowment,
		KindExchange:
		return true
	}
	return false
}

// We don't want any dead characters, so since we trim the generated
// SHA hash anyway, we trim it to a length that plays well with the above.
// (Pads the result out to a multiple of 5 bytes so that a byte32 encoding has
// no filler).
//
// Note that ETH does this too, and uses a 20-byte subset of a 32-byte hash.
// The possibility of collision is low: As of June 2018, the BTC hashpower is 42
// exahashes per second. If that much hashpower is applied to this problem, the
// likelihood of generating a collision in one year is about 1 in 10^19.
const keyLength = 24

// MinDataLength is the minimum acceptable length for the data to be used
// as input to generate. This will prevent simple errors like trying to
// create an address from an empty key.
const MinDataLength = 12

// Generate creates an address of a given type from an array of bytes (which
// would normally be a public key). It is an error if len(data) < MinDataLength
// or if kind is not a valid kind.
func Generate(kind Kind, data []byte) (string, error) {
	if !IsValidKind(kind) {
		return "", newError("invalid kind")
	}
	if len(data) < MinDataLength {
		return "", newError("insufficient quantity of data")
	}

	prefix := ndx("n")<<11 + ndx("d")<<6 + ndx(string(kind))<<1
	hdr := []byte{byte((prefix >> 8) & 0xFF), byte(prefix & 0xFF)}
	h1 := sha256.Sum256(data)
	h2 := append(hdr, h1[len(h1)-keyLength:]...)
	ck := adler32.New()
	ck.Write(h2)
	h2 = append(h2, ck.Sum(nil)...)

	enc := base32.NewEncoding(NdauAlphabet)
	r := enc.EncodeToString(h2)
	return r, nil
}

// Validate tests if an address is valid on its face.
// It checks the the nd prefix, the address kind, and the checksum.
func Validate(addr string) error {
	addr = strings.ToLower(addr)
	if !strings.HasPrefix(addr, "nd") {
		return newError("not an ndau key")
	}
	if !IsValidKind(Kind(addr[2:3])) {
		return newError("unknown address kind " + addr[2:3])
	}
	enc := base32.NewEncoding(NdauAlphabet)
	h, err := enc.DecodeString(addr)
	if err != nil {
		return err
	}
	k := h[:len(h)-4]
	ck := adler32.New()
	ck.Write(k)
	ckb := ck.Sum(nil)
	if string(ckb) != string(h[len(h)-4:]) {
		return newError("checksum failure")
	}
	return nil
}
