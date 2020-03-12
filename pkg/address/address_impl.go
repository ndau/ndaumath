package address

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/ndau/ndaumath/pkg/b32"
)

// An ndau address is the result of a mathematical process over a public key. It
// is a byte32 encoding, using a custom alphabet, of a portion of the SHA256
// hash of the key, concatenated with some additional marker and checksum
// information. The result is a key that always starts with a specific 2-letter
// prefix (nd for the main chain and tn for the testnet), plus one more
// character that specifies the type of address.

// Kind indicates the type of address in use; this is an external indication
// designed to help users evaluate their own actions; it may or may not be
// enforced by the blockchain.
// We don't define a type for it so that it serializes naturally.  We give up
// type safety to avoid msgp hassles.  See address.go for details.
//type Kind byte

func emptyA() Address {
	return Address{}
}

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

// All addresses start with this 2-byte prefix, followed by a kind byte.
const addrPrefix string = "nd"
const kindOffset int = len(addrPrefix)

// predefined address kinds
const (
	KindUser        byte = 'a'
	KindNdau        byte = 'n'
	KindEndowment   byte = 'e'
	KindExchange    byte = 'x'
	KindBPC         byte = 'b'
	KindMarketMaker byte = 'm'
)

// IsValidKind returns true if the last letter of a is one of the currently-valid kinds
func IsValidKind(k byte) bool {
	switch k {
	case KindUser,
		KindNdau,
		KindEndowment,
		KindExchange,
		KindBPC,
		KindMarketMaker:
		return true
	}
	return false
}

// ParseKind returns a Kind or an explanation of why the supplied value is not one.
func ParseKind(i interface{}) (byte, error) {
	b := byte(0)
	switch v := i.(type) {
	case string:
		if v == "" {
			return b, fmt.Errorf("empty string is not a valid Kind")
		}
		v = strings.ToLower(v)
		switch v {
		case "u", "user":
			b = KindUser
		case "ndau":
			b = KindNdau
		case "endowment":
			b = KindEndowment
		case "exchange":
			b = KindExchange
		case "bpc":
			b = KindBPC
		case "marketmaker":
			b = KindMarketMaker
		default:
			b = byte(v[0])
		}
	case rune:
		b = byte(strings.ToLower(string(v))[0])
	case byte:
		b = v
	case int8:
		b = byte(v)
	default:
		return b, fmt.Errorf("Kind cannot be parsed from %T", i)
	}

	if !IsValidKind(b) {
		return b, fmt.Errorf("%q is not a valid Kind", string(b))
	}
	return b, nil
}

// HashTrim is the number of bytes that we trim the input hash to.
//
// We don't want any dead characters, so since we trim the generated
// SHA hash anyway, we trim it to a length that plays well with the above,
// meaning that we want it to pad the result out to a multiple of 5 bytes
// so that a byte32 encoding has no filler).
//
// Note that ETH does something similar, and uses a 20-byte subset of a 32-byte hash.
// The possibility of collision is low: As of June 2018, the BTC hashpower is 42
// exahashes per second. If that much hashpower is applied to this problem, the
// likelihood of generating a collision in one year is about 1 in 10^19.
const HashTrim = 26

// AddrLength is the length of the generated address, in characters
const AddrLength = 48

// MinDataLength is the minimum acceptable length for the data to be used
// as input to generate. This will prevent simple errors like trying to
// create an address from an empty key.
const MinDataLength = 12

// Generate creates an address of a given kind from an array of bytes (which
// would normally be a public key). It is an error if len(data) < MinDataLength
// or if kind is not a valid kind.
// Since length changes are explicitly disallowed, we can use a relatively simple
// crc model to have a short (16-bit) checksum and still be quite safe against
// transposition and typos.
func Generate(kind byte, data []byte) (Address, error) {
	if !IsValidKind(kind) {
		return emptyA(), newError(fmt.Sprintf("invalid kind: %x", kind))
	}
	if len(data) < MinDataLength {
		return emptyA(), newError("insufficient quantity of data")
	}
	// the hash contains the last HashTrim bytes of the sha256 of the data
	h := sha256.Sum256(data)
	h1 := h[len(h)-HashTrim:]

	// an ndau address always starts with nd and a "kind" character
	// so we figure out what characters we want and build that into a header
	prefix :=
		b32.Index(addrPrefix[0:1])<<11 +
			b32.Index(addrPrefix[1:2])<<6 +
			b32.Index(string(kind))<<1
	hdr := []byte{byte((prefix >> 8) & 0xFF), byte(prefix & 0xFF)}
	h2 := append(hdr, h1...)
	// then we checksum that result and append the checksum
	h2 = append(h2, b32.Checksum16(h2)...)

	r := b32.Encode(h2)
	return Address{addr: r}, nil
}

// Validate tests if an address is valid on its face.
// It checks the address kind, and the checksum.
// It does NOT test the nd prefix, as that may vary -- clients should test that
// themselves.
func Validate(addr string) (Address, error) {
	addr = strings.ToLower(addr)
	// if !strings.HasPrefix(addr, "nd") {
	// 	return emptyA(), newError("not an ndau key")
	// }
	if len(addr) != AddrLength {
		return emptyA(), newError(fmt.Sprintf("not a valid address length '%s'", addr))
	}
	if kind := addr[kindOffset]; !IsValidKind(kind) {
		return emptyA(), newError(fmt.Sprintf("unknown address kind: %x", kind))
	}
	h, err := b32.Decode(addr)
	if err != nil {
		return emptyA(), err
	}
	// now check the two bytes of the checksum
	if !b32.Check(h[:len(h)-2], h[len(h)-2:]) {
		// uncomment these lines if you want to regenerate a key
		// that matches the main body of the key you gave, but with a proper checksum
		// -------
		// h[len(h)-2] = byte((ck >> 8) & 0xFF)
		// h[len(h)-1] = byte(ck & 0xFF)
		// s := base32.NewEncoding(NdauAlphabet).EncodeToString(h)
		// fmt.Println(s)
		// -------
		return emptyA(), newError("checksum failure")
	}
	return Address{addr: addr}, nil
}

// String gives us a human-readable form of an address, because sometimes we just need that.
func (z Address) String() string {
	return z.addr
}

// Kind returns the kind byte of the address.
func (z Address) Kind() byte {
	return z.addr[kindOffset]
}

// Revalidate this address to ensure it is legitimate
func (z Address) Revalidate() error {
	_, err := Validate(z.addr)
	return err
}
