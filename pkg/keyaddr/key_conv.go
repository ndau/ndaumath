package keyaddr

import (
	"github.com/oneiro-ndev/ndaumath/pkg/key"
	"github.com/oneiro-ndev/ndaumath/pkg/signature"
	"github.com/pkg/errors"
)

// KeyFromExtended constructs a `*Key` from a `*key.ExtendedKey`
func KeyFromExtended(k *key.ExtendedKey) (*Key, error) {
	kb, err := k.MarshalText()
	if err != nil {
		return nil, err
	}
	return &Key{Key: string(kb)}, nil
}

// KeyFromPublic constructs a `*Key` from a `signature.PublicKey`
func KeyFromPublic(k signature.PublicKey) (*Key, error) {
	text, err := k.MarshalText()
	if err != nil {
		return nil, errors.Wrap(err, "marshalling into text")
	}
	ekey := new(key.ExtendedKey)
	err = ekey.UnmarshalText(text)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling into ekey")
	}
	return KeyFromExtended(ekey)
}

// KeyFromPrivate constructs a `*Key` from a `signature.PrivateKey`
func KeyFromPrivate(k signature.PrivateKey) (*Key, error) {
	text, err := k.MarshalText()
	if err != nil {
		return nil, errors.Wrap(err, "marshalling into text")
	}
	ekey := new(key.ExtendedKey)
	err = ekey.UnmarshalText(text)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling into ekey")
	}
	return KeyFromExtended(ekey)
}

// ToExtended constructs a `*key.ExtendedKey` from a `Key`
func (k Key) ToExtended() (*key.ExtendedKey, error) {
	ekey := new(key.ExtendedKey)
	err := ekey.UnmarshalText([]byte(k.Key))
	return ekey, err
}

// ToPublicKey constructs a `signature.PublicKey` from a `*Key`
func (k *Key) ToPublicKey() (signature.PublicKey, error) {
	out := signature.PublicKey{}
	ekey, err := k.ToExtended()
	if err != nil {
		return out, errors.Wrap(err, "converting to extendedkey")
	}
	pub, err := ekey.Public()
	if err != nil {
		return out, errors.Wrap(err, "making public")
	}
	text, err := pub.MarshalText()
	if err != nil {
		return out, errors.Wrap(err, "marshalling")
	}
	err = out.UnmarshalText(text)
	return out, err
}

// ToPrivateKey constructs a `signature.PrivateKey` from a `*Key`
func (k *Key) ToPrivateKey() (signature.PrivateKey, error) {
	out := signature.PrivateKey{}
	ekey, err := k.ToExtended()
	if err != nil {
		return out, errors.Wrap(err, "converting to extendedkey")
	}
	if !ekey.IsPrivate() {
		return out, errors.New("cannot convert public key to private key")
	}
	text, err := ekey.MarshalText()
	if err != nil {
		return out, errors.Wrap(err, "marshalling")
	}
	err = out.UnmarshalText(text)
	return out, err
}
