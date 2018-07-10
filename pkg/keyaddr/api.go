package keyaddr

// This package provides an interface to the ndaumath library for use in React.
// It is built using the gomobile tool, so the API is constrained to particular types of parameters:
//
// * string
// * signed integer and floating point types
// * []byte
// * functions with specific restrictions
// * structs and interfaces consisting of only these types

// This package, therefore, consists mainly of wrappers so that we don't have to modify our
// idiomatic Go code to conform to these requirements.

import (
	"errors"
	"strings"

	"github.com/oneiro-ndev/ndaumath/pkg/address"
	"github.com/oneiro-ndev/ndaumath/pkg/b32"
	"github.com/oneiro-ndev/ndaumath/pkg/key"
	"github.com/oneiro-ndev/ndaumath/pkg/words"
)

// WordsFromBytes takes an array of bytes and converts it to a space-separated list of
// words that act as a mnemonic. A 16-byte input array will generate a list of 12 words.
func WordsFromBytes(lang string, b []byte) (string, error) {
	sa, err := words.FromBytes(lang, b)
	if err != nil {
		return "", err
	}
	return strings.Join(sa, " "), nil
}

// WordsToBytes takes a space-separated list of words and generates the set of bytes
// from which it was generated (or an error).
func WordsToBytes(lang string, w string) ([]byte, error) {
	wordlist := strings.Split(w, " ")
	return words.ToBytes(lang, wordlist)
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

// NewKey takes a seed (an array of bytes) and creates a private master
// key from it. The key is returned as a string representation of the key;
// it is converted to and from the internal representation by its member functions.
func NewKey(seed []byte) (*Key, error) {
	mk, err := key.NewMaster([]byte(seed), key.NdauPrivateKeyID)
	if err != nil {
		return nil, err
	}
	return &Key{Key: mk.String()}, nil
}

// Neuter returns an extended public key from any other extended key.
// If the key is an extended private key, it generates the matching public key.
// If the key is already a public key, it just returns itself.
// It is an error if the key is hardened.
func (k *Key) Neuter() (*Key, error) {
	ekey, err := key.NewKeyFromString(k.Key)
	if err != nil {
		return nil, err
	}
	nk, err := ekey.Neuter()
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

// Sign uses the given key to sign a message; the message will usually be
// the hash of some longer message. It returns a signature object.
// The key must be a private key.
func (k *Key) Sign(msg []byte) (*Signature, error) {
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
func (k *Key) NdauAddress() (*Address, error) {
	ekey, err := key.NewKeyFromString(k.Key)
	if err != nil {
		return nil, err
	}

	a, err := address.Generate(address.KindUser, ekey.PubKeyBytes())
	if err != nil {
		return nil, err
	}

	return &Address{a.String()}, nil
}
