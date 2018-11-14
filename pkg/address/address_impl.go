package address

import (
	"crypto/sha256"
	"errors"
	"strings"

	"github.com/oneiro-ndev/ndaumath/pkg/b32"
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
type Kind string

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

// we want the first letters of the address to be "XX?" where XX is the code for the particular
// type of network (mainnet==nd and testnet==tn). ? is the kind of
// address. Valid address types are as follows:
const (
	KindUser        Kind = "a"
	KindNdau        Kind = "n"
	KindEndowment   Kind = "e"
	KindExchange    Kind = "x"
	KindBPC         Kind = "b"
	KindMarketMaker Kind = "m"
)

// IsValidKind returns true if the last letter of a is one of the currently-valid Kinds
func IsValidKind(a Kind) bool {
	if len(a) == 0 {
		return false
	}
	switch a[len(a)-1:] {
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

// NewKind constructs an address.Kind, given a string corresponding to the kind value
// The value used is the last letter of the kind string (one of anex)
func NewKind(kind string) (Kind, error) {
	kinds := map[string]Kind{
		string(KindUser):      KindUser,
		string(KindNdau):      KindNdau,
		string(KindExchange):  KindExchange,
		string(KindEndowment): KindEndowment,
	}

	if kind == "" {
		return "", errors.New("kind must not be blank")
	}
	ltr := kind[len(kind)-1:]
	k, ok := kinds[ltr]
	if ok {
		return k, nil
	}
	return k, errors.New("invalid kind character")
}

// splitKind returns the 2-letter prefix and the kind (as a string) for s, which should be a Kind
// containing either 1 or 3 characters.
func splitKind(s Kind) (string, string, error) {
	switch len(s) {
	case 1:
		return "nd", string(s), nil
	case 3:
		return string(s)[0:2], string(s)[2:3], nil
	default:
		return "", "", errors.New("not a valid Kind")
	}
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
func Generate(kind Kind, data []byte) (Address, error) {
	if !IsValidKind(kind) {
		return emptyA(), newError("invalid kind")
	}
	if len(data) < MinDataLength {
		return emptyA(), newError("insufficient quantity of data")
	}
	// the hash contains the last HashTrim bytes of the sha256 of the data
	h := sha256.Sum256(data)
	h1 := h[len(h)-HashTrim:]

	// an ndau address always starts with nd and a "kind" character
	// so we figure out what characters we want and build that into a header
	sprefix, skind, err := splitKind(kind)
	if err != nil {
		return emptyA(), newError("invalid kind")
	}
	prefix := b32.Index(sprefix[0:1])<<11 + b32.Index(sprefix[1:2])<<6 + b32.Index(skind)<<1
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
		return emptyA(), newError("not a valid address length")
	}
	if !IsValidKind(Kind(addr[2:3])) {
		return emptyA(), newError("unknown address kind " + addr[2:3])
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
