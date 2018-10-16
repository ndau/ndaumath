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

func (k Key) ekey() (*key.ExtendedKey, error) {
	ekey := new(key.ExtendedKey)
	err := ekey.UnmarshalText([]byte(k.Key))
	return ekey, err
}

func asKey(k *key.ExtendedKey) (*Key, error) {
	kb, err := k.MarshalText()
	if err != nil {
		return nil, err
	}
	return &Key{Key: string(kb)}, nil
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
	return asKey(mk)
}

// FromString acts like a constructor so that the wallet can build a Key object
// from a string representation of it.
func FromString(s string) (*Key, error) {
	ekey := new(key.ExtendedKey)
	err := ekey.UnmarshalText([]byte(s))
	if err != nil {
		return nil, err
	}

	// re-marshal for reasons?
	return asKey(ekey)
}

// FromOldString is FromString, but it operates on the old key serialization format.
//
// The returned object will be serialized in the new format, so future calls
// to FromString will succeed.
func FromOldString(s string) (*Key, error) {
	ekey, err := key.FromOldSerialization(s)
	if err != nil {
		return nil, err
	}
	return asKey(ekey)
}

// DeriveFrom accepts a parent key and its known path, plus a desired child path
// and derives the child key from the parent according to the path info.
// Note that the parent's known path is simply believed -- we have no mechanism to
// check that it's true.
func DeriveFrom(parentKey string, parentPath, childPath string) (*Key, error) {
	k, err := FromString(parentKey)
	if err != nil {
		return nil, err
	}
	e, err := k.ekey()
	if err != nil {
		return nil, err
	}
	e, err = e.DeriveFrom(parentPath, childPath)
	if err != nil {
		return nil, err
	}
	return asKey(e)
}

// ToPublic returns an extended public key from any other extended key.
// If the key is an extended private key, it generates the matching public key.
// If the key is already a public key, it just returns itself.
// It is an error if the key is hardened.
func (k *Key) ToPublic() (*Key, error) {
	ekey, err := k.ekey()
	if err != nil {
		return nil, err
	}
	nk, err := ekey.Public()
	if err != nil {
		return nil, err
	}
	return asKey(nk)
}

// Child returns the n'th child of the given extended key. The child is of the
// same type (public or private) as the parent. Although n is typed as a signed
// integer, this is due to the limitations of gomobile; n may not be negative.
// It is an error if the given key is a hardened key.
func (k *Key) Child(n int32) (*Key, error) {
	if n < 0 {
		return nil, errors.New("child index cannot be negative")
	}
	ekey, err := k.ekey()
	if err != nil {
		return nil, err
	}
	ndx := uint32(n)
	nk, err := ekey.Child(ndx)
	if err != nil {
		return nil, err
	}
	return asKey(nk)
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
	ekey, err := k.ekey()
	if err != nil {
		return nil, err
	}
	nk, err := ekey.HardenedChild(uint32(n))
	if err != nil {
		return nil, err
	}
	return asKey(nk)
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
	ekey, err := k.ekey()
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
func (k *Key) NdauAddress(string) (*Address, error) {
	skind := string(address.KindUser)
	ekey, err := k.ekey()
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
	ekey, err := k.ekey()
	if err != nil {
		return false, err
	}
	return ekey.IsPrivate(), nil
}
