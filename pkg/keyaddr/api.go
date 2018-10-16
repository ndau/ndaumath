package keyaddr

// This package provides an interface to the ndaumath library for use in React and in particular react-native.
// It is built using the gomobile tool, so the API is constrained to particular types of parameters:
//
// * string
// * signed integer and floating point types
// * []byte
// * functions with specific restrictions
// * structs and interfaces consisting of only these types
//
// Unfortunately, react-native puts additional requirements that makes []byte particularly
// challenging to use. So what we are going to do is use a base-64 encoding of []byte to convert
// it to a string and pass the array of bytes back and forth that way.
//
// This is distinct from using base32 encoding (b32) in a signature; that's something we expect
// to be user-visible, so we're using a specific variant of base 32.

// This package, therefore, consists mainly of wrappers so that we don't have to modify our
// idiomatic Go code to conform to these requirements.

import (
	"encoding/base64"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/oneiro-ndev/ndaumath/pkg/address"
	"github.com/oneiro-ndev/ndaumath/pkg/b32"
	"github.com/oneiro-ndev/ndaumath/pkg/key"
	"github.com/oneiro-ndev/ndaumath/pkg/words"
)

// WordsFromBytes takes an array of bytes and converts it to a space-separated list of
// words that act as a mnemonic. A 16-byte input array will generate a list of 12 words.
func WordsFromBytes(lang string, data string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	sa, err := words.FromBytes(lang, b)
	if err != nil {
		return "", err
	}
	return strings.Join(sa, " "), nil
}

// WordsToBytes takes a space-separated list of words and generates the set of bytes
// from which it was generated (or an error). The bytes are encoded as a base64 string
// using standard base64 encoding, as defined in RFC 4648.
func WordsToBytes(lang string, w string) (string, error) {
	wordlist := strings.Split(w, " ")
	b, err := words.ToBytes(lang, wordlist)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// WordsFromPrefix accepts a language and a prefix string and returns a sorted, space-separated list
// of words that match the given prefix. max can be used to limit the size of the returned list
// (if max is 0 then all matches are returned, which could be up to 2K if the prefix is empty).
func WordsFromPrefix(lang string, prefix string, max int) string {
	return words.FromPrefix(lang, prefix, max)
}

// Key is the object that contains a public or private key
type Key struct {
	Key string
}

// Signature is the result of signing a block of data with a key.
type Signature struct {
	Signature string
}

// Address is an Ndau Address, derived from a public key.
type Address struct {
	Address string
}

// NewKey takes a seed (an array of bytes encoded as a base64 string) and creates a private master
// key from it. The key is returned as a string representation of the key;
// it is converted to and from the internal representation by its member functions.
func NewKey(seedstr string) (*Key, error) {
	seed, err := base64.StdEncoding.DecodeString(seedstr)
	if err != nil {
		return nil, err
	}
	mk, err := key.NewMaster([]byte(seed))
	if err != nil {
		return nil, err
	}
	return &Key{Key: mk.String()}, nil
}

// FromString acts like a constructor so that the wallet can build a Key object
// from a string representation of it.
func FromString(s string) (*Key, error) {
	ekey, err := key.NewKeyFromString(s)
	if err != nil {
		return nil, err
	}
	return &Key{ekey.String()}, nil
}

type pathElement struct {
	id     int32
	harden bool
}

type path []pathElement

func newPath(s string) (path, error) {
	// remove all whitespace
	s = strings.Replace(s, " ", "", -1)
	// treat root specially
	if s == "/" {
		return path{}, nil
	}
	// now validate the path
	// note that other than the pure root marker that we already handled,
	// the numeric part after the slash is not optional
	valpat := regexp.MustCompile("^(/([0-9]+)'?)+$")
	if !valpat.MatchString(s) {
		return nil, errors.New("Not a valid path string")
	}

	parsepat := regexp.MustCompile("/([0-9]+)('?)")
	saa := parsepat.FindAllStringSubmatch(s, -1)
	// saa now has one entry for each path element, and
	// for each entry it has the 0th element as the whole path string,
	// the first as the path ID, and the second as either
	// an apostrophe or an empty string.
	p := make(path, len(saa))
	for i := range saa {
		var err error
		n, err := strconv.ParseInt(saa[i][1], 10, 32)
		if err != nil {
			return nil, err
		}
		p[i].id = int32(n)
		p[i].harden = saa[i][2] == "'"
	}
	return p, nil
}

func (p path) isParentOf(c path) bool {
	// if the parent is not shorter than the purported child, it can't be a parent
	if len(c) <= len(p) {
		return false
	}
	// everything up to the length of the parent has to be the same
	for i := range p {
		if c[i].id != p[i].id || c[i].harden != p[i].harden {
			return false
		}
	}
	return true
}

// DeriveFrom accepts a parent key and its known path, plus a desired child path
// and derives the child key from the parent according to the path info.
// Note that the parent's known path is simply believed -- we have no mechanism to
// check that it's true.
func DeriveFrom(parentKey string, parentPath, childPath string) (*Key, error) {
	ppath, err := newPath(parentPath)
	if err != nil {
		return nil, err
	}
	cpath, err := newPath(childPath)
	if err != nil {
		return nil, err
	}
	if !ppath.isParentOf(cpath) {
		return nil, errors.New("child is not descended from parent")
	}
	// if we get here we know that ppath is a subset of cpath so we can trim cpath
	cpath = cpath[len(ppath):]
	k, err := FromString(parentKey)
	if err != nil {
		return nil, err
	}
	// now iterate
	for _, e := range cpath {
		if e.harden {
			k, err = k.HardenedChild(int32(e.id))
		} else {
			k, err = k.Child(int32(e.id))
		}
		if err != nil {
			return nil, err
		}
	}
	return k, err
}

// ToPublic returns an extended public key from any other extended key.
// If the key is an extended private key, it generates the matching public key.
// If the key is already a public key, it just returns itself.
// It is an error if the key is hardened.
func (k *Key) ToPublic() (*Key, error) {
	ekey, err := key.NewKeyFromString(k.Key)
	if err != nil {
		return nil, err
	}
	nk, err := ekey.Public()
	if err != nil {
		return nil, err
	}
	return &Key{nk.String()}, nil
}

// Child returns the n'th child of the given extended key. The child is of the
// same type (public or private) as the parent. Although n is typed as a signed
// integer, this is due to the limitations of gomobile; n may not be negative.
// It is an error if the given key is a hardened key.
func (k *Key) Child(n int32) (*Key, error) {
	if n < 0 {
		return nil, errors.New("child index cannot be negative")
	}
	ekey, err := key.NewKeyFromString(k.Key)
	if err != nil {
		return nil, err
	}
	ndx := uint32(n)
	nk, err := ekey.Child(ndx)
	if err != nil {
		return nil, err
	}
	return &Key{nk.String()}, nil
}

// HardenedChild returns the n'th hardened child of the given extended key.
// The parent key must be a private key.
// A HardenedChild is guaranteed to have been derived from a private key.
// Although n is typed as a signed integer, this is due to the limitations of gomobile;
// n may not be negative.
// It is an error if the given key is already a hardened key.
func (k *Key) HardenedChild(n int32) (*Key, error) {
	if n < 0 {
		return nil, errors.New("child index cannot be negative")
	}
	ekey, err := key.NewKeyFromString(k.Key)
	if err != nil {
		return nil, err
	}
	ndx := uint32(key.HardenedKeyStart) + uint32(n)
	nk, err := ekey.Child(ndx)
	if err != nil {
		return nil, err
	}
	return &Key{nk.String()}, nil
}

// Sign uses the given key to sign a message; the message must be the
// standard base64 encoding of the bytes of the message.
// It returns a signature object.
// The key must be a private key.
func (k *Key) Sign(msgstr string) (*Signature, error) {
	msg, err := base64.StdEncoding.DecodeString(msgstr)
	if err != nil {
		return nil, err
	}
	ekey, err := key.NewKeyFromString(k.Key)
	if err != nil {
		return nil, err
	}
	pk, err := ekey.ECPrivKey()
	if err != nil {
		return nil, err
	}
	sig, err := pk.Sign(msg)
	if err != nil {
		return nil, err
	}
	return &Signature{b32.Encode(sig.Serialize())}, nil
}

// NdauAddress returns the ndau address associated with the given key.
// Key can be either public or private; if it is private it will be
// converted to a public key first.
func (k *Key) NdauAddress(chainid string) (*Address, error) {
	skind := string(address.KindUser)
	switch chainid {
	case "nd":
		// we're good
	case "tn":
		skind = chainid + string(skind)
	default:
		return nil, errors.New("invalid chain id")
	}

	ekey, err := key.NewKeyFromString(k.Key)
	if err != nil {
		return nil, err
	}

	a, err := address.Generate(address.Kind(skind), ekey.PubKeyBytes())
	if err != nil {
		return nil, err
	}

	return &Address{a.String()}, nil
}

// IsPrivate tests if a given key is a private key; will return non-nil
// error if the key is invalid.
func (k *Key) IsPrivate() (bool, error) {
	ekey, err := key.NewKeyFromString(k.Key)
	if err != nil {
		return false, err
	}
	return ekey.IsPrivate(), nil
}
