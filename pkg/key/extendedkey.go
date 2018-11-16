// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2018 Oneiro, LLC
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// This source code is basically a subset of some of btcsuite/bctd. As ndau is
// not bitcoin, there are many things in btcd that are unnecessary for ndau and
// add dependencies and complexity; consequently, we have duplicated the source
// and extracted only those portions that we actually need. The license information
// above is from the original licensing terms for btcsuite.

// For package composability reasons, certain functionality originally located
// in this package was moved to the sibling package `bip32`.

package key

// References:
//   [BIP32]: BIP0032 - Hierarchical Deterministic Wallets
//   https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
//   https://github.com/btcsuite/btcd

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"unicode/utf8"

	"github.com/oneiro-ndev/ndaumath/pkg/signature"

	"github.com/btcsuite/btcd/btcec"
	"github.com/oneiro-ndev/ndaumath/pkg/b32"
	"github.com/oneiro-ndev/ndaumath/pkg/bip32"
)

// bip32 constants are described in bip32 documentation
const (
	RecommendedSeedLen = bip32.RecommendedSeedLen
	HardenedKeyStart   = bip32.HardenedKeyStart
	MinSeedBytes       = bip32.MinSeedBytes
	MaxSeedBytes       = bip32.MaxSeedBytes

	// maxUint8 is the max positive integer which can be serialized in a uint8
	maxUint8 = 1<<8 - 1
)

var (
	// ErrDeriveHardFromPublic describes an error in which the caller
	// attempted to derive a hardened extended key from a public key.
	ErrDeriveHardFromPublic = errors.New(
		"cannot derive a hardened key from a public key")

	// ErrDeriveBeyondMaxDepth describes an error in which the caller
	// has attempted to derive more than 255 keys from a root key.
	ErrDeriveBeyondMaxDepth = errors.New(
		"cannot derive a key with more than 255 indices in its path")

	// ErrNotPrivExtKey describes an error in which the caller attempted
	// to extract a private key from a public extended key.
	ErrNotPrivExtKey = errors.New(
		"unable to create private keys from a public extended key")

	// ErrInvalidChild describes an error in which the child at a specific
	// index is invalid due to the derived key falling outside of the valid
	// range for secp256k1 private keys.  This error indicates the caller
	// should simply ignore the invalid child extended key at this index and
	// increment to the next index.
	ErrInvalidChild = errors.New("the extended key at this index is invalid")

	// ErrUnusableSeed describes an error in which the provided seed is not
	// usable due to the derived key falling outside of the valid range for
	// secp256k1 private keys.  This error indicates the caller must choose
	// another seed.
	ErrUnusableSeed = bip32.ErrUnusableSeed

	// ErrInvalidSeedLen describes an error in which the provided seed or
	// seed length is not in the allowed range.
	ErrInvalidSeedLen = bip32.ErrInvalidSeedLen

	// ErrBadChecksum describes an error in which the checksum encoded with
	// a serialized extended key does not match the calculated value.
	ErrBadChecksum = errors.New("bad extended key checksum")

	// ErrInvalidKeyLen describes an error in which the provided serialized
	// key is not the expected length.
	ErrInvalidKeyLen = errors.New(
		"the provided serialized extended key length is invalid")

	// ErrUnknownHDKeyID describes an error where the provided id which
	// is intended to identify the network for a hierarchical deterministic
	// private extended key is not registered.
	ErrUnknownHDKeyID = errors.New("unknown hd private extended key bytes")

	// ErrInvalidKeyEncoding describes an error where the provided key
	// is not properly formatted in base32 encoding
	ErrInvalidKeyEncoding = errors.New("invalid key encoding")

	// NdauPrivateKeyID is the special prefix we use for ndau private keys
	NdauPrivateKeyID = [3]byte{99, 103, 31} // npvt
	// NdauPublicKeyID is another special prefix
	NdauPublicKeyID = [3]byte{99, 100, 16} // npub
	// TestPrivateKeyID is another special prefix
	TestPrivateKeyID = [3]byte{139, 103, 31} // tpvt
	// TestPublicKeyID is another special prefix
	TestPublicKeyID = [3]byte{139, 100, 16} // tpub
)

// fingerprint calculates the checksum of sha256(b).
func fingerprint(buf []byte) []byte {
	hasher := sha256.New()
	hasher.Write(buf)
	return b32.Checksum24(hasher.Sum(nil))
}

// doubleHashB calculates hash(hash(b)) and returns the resulting bytes.
func doubleHashB(b []byte) []byte {
	first := sha256.Sum256(b)
	second := sha256.Sum256(first[:])
	return second[:]
}

// ExtendedKey houses all the information needed to support a hierarchical
// deterministic extended key.  See the package overview documentation for
// more details on how to use extended keys.
type ExtendedKey struct {
	key       []byte // This will be the pubkey for extended pub keys
	pubKey    []byte // This will only be set for extended priv keys
	chainCode []byte
	depth     uint8
	parentFP  []byte
	childNum  uint32
	isPrivate bool
}

// ensure ExtendedKey implements Text(Un)Marshaller
var _ encoding.TextMarshaler = (*ExtendedKey)(nil)
var _ encoding.TextUnmarshaler = (*ExtendedKey)(nil)

// NewExtendedKey returns a new instance of an extended key with the given
// fields.  No error checking is performed here as it's only intended to be a
// convenience method used to create a populated struct. This function should
// only by used by applications that need to create custom ExtendedKeys. All
// other applications should just use NewMaster, Child, or Public.
func NewExtendedKey(key, chainCode, parentFP []byte, depth uint8,
	childNum uint32, isPrivate bool) *ExtendedKey {

	// NOTE: The pubKey field is intentionally left nil so it is only
	// computed and memoized as required.
	return &ExtendedKey{
		key:       key,
		chainCode: chainCode,
		depth:     depth,
		parentFP:  parentFP,
		childNum:  childNum,
		isPrivate: isPrivate,
	}
}

// PubKeyBytes returns bytes for the serialized compressed public key associated
// with this extended key in an efficient manner including memoization as
// necessary.
//
// When the extended key is already a public key, the key is simply returned as
// is since it's already in the correct form.  However, when the extended key is
// a private key, the public key will be calculated and memoized so future
// accesses can simply return the cached result.
func (k *ExtendedKey) PubKeyBytes() []byte {
	// Just return the key if it's already an extended public key.
	if !k.isPrivate {
		return k.key
	}

	// This is a private extended key, so calculate and memoize the public
	// key if needed.
	if len(k.pubKey) == 0 {
		k.pubKey = bip32.PrivateToPublic(k.key)
	}

	return k.pubKey
}

// IsPrivate returns whether or not the extended key is a private extended key.
//
// A private extended key can be used to derive both hardened and non-hardened
// child private and public extended keys.  A public extended key can only be
// used to derive non-hardened child public extended keys.
func (k *ExtendedKey) IsPrivate() bool {
	return k.isPrivate
}

// Depth returns the current derivation level with respect to the root.
//
// The root key has depth zero, and the field has a maximum of 255 due to
// how depth is serialized.
func (k *ExtendedKey) Depth() uint8 {
	return k.depth
}

// ParentFingerprint returns a fingerprint of the parent extended key from which
// this one was derived.
// Since the fingerprint is 3 bytes, we set the high byte to 0 before returning it
// as a uint32
func (k *ExtendedKey) ParentFingerprint() uint32 {
	b := append([]byte{0}, k.parentFP...)
	return binary.BigEndian.Uint32(b)
}

// Child returns a derived child extended key at the given index.  When this
// extended key is a private extended key (as determined by the IsPrivate
// function), a private extended key will be derived.  Otherwise, the derived
// extended key will be also be a public extended key.
//
// When the index is greater to or equal than the HardenedKeyStart constant, the
// derived extended key will be a hardened extended key.  It is only possible to
// derive a hardended extended key from a private extended key.  Consequently,
// this function will return ErrDeriveHardFromPublic if a hardened child
// extended key is requested from a public extended key.
//
// A hardened extended key is useful since, as previously mentioned, it requires
// a parent private extended key to derive.  In other words, normal child
// extended public keys can be derived from a parent public extended key (no
// knowledge of the parent private key) whereas hardened extended keys may not
// be.
//
// NOTE: There is an extremely small chance (< 1 in 2^127) the specific child
// index does not derive to a usable child.  The ErrInvalidChild error will be
// returned if this should occur, and the caller is expected to ignore the
// invalid child and simply increment to the next index.
func (k *ExtendedKey) Child(i uint32) (*ExtendedKey, error) {
	// Prevent derivation of children beyond the max allowed depth.
	if k.depth == maxUint8 {
		return nil, ErrDeriveBeyondMaxDepth
	}

	// There are four scenarios that could happen here:
	// 1) Private extended key -> Hardened child private extended key
	// 2) Private extended key -> Non-hardened child private extended key
	// 3) Public extended key -> Non-hardened child public extended key
	// 4) Public extended key -> Hardened child public extended key (INVALID!)

	// Case #4 is invalid, so error out early.
	// A hardened child extended key may not be created from a public
	// extended key.
	isChildHardened := i >= HardenedKeyStart
	if !k.isPrivate && isChildHardened {
		return nil, ErrDeriveHardFromPublic
	}

	// The data used to derive the child key depends on whether or not the
	// child is hardened per [BIP32].
	//
	// For hardened children:
	//   0x00 || ser256(parentKey) || ser32(i)
	//
	// For normal children:
	//   serP(parentPubKey) || ser32(i)
	keyLen := 33
	data := make([]byte, keyLen+4)
	if isChildHardened {
		// Case #1.
		// When the child is a hardened child, the key is known to be a
		// private key due to the above early return.  Pad it with a
		// leading zero as required by [BIP32] for deriving the child.
		copy(data[1:], k.key)
	} else {
		// Case #2 or #3.
		// This is either a public or private extended key, but in
		// either case, the data which is used to derive the child key
		// starts with the secp256k1 compressed public key bytes.
		copy(data, k.PubKeyBytes())
	}
	binary.BigEndian.PutUint32(data[keyLen:], i)

	// Take the HMAC-SHA512 of the current key's chain code and the derived
	// data:
	//   I = HMAC-SHA512(Key = chainCode, Data = data)
	hmac512 := hmac.New(sha512.New, k.chainCode)
	hmac512.Write(data)
	ilr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = intermediate key used to derive the child
	//   Ir = child chain code
	il := ilr[:len(ilr)/2]
	childChainCode := ilr[len(ilr)/2:]

	// Both derived public or private keys rely on treating the left 32-byte
	// sequence calculated above (Il) as a 256-bit integer that must be
	// within the valid range for a secp256k1 private key.  There is a small
	// chance (< 1 in 2^127) this condition will not hold, and in that case,
	// a child extended key can't be created for this index and the caller
	// should simply increment to the next index.
	ilNum := new(big.Int).SetBytes(il)
	if ilNum.Cmp(btcec.S256().N) >= 0 || ilNum.Sign() == 0 {
		return nil, ErrInvalidChild
	}

	// The algorithm used to derive the child key depends on whether or not
	// a private or public child is being derived.
	//
	// For private children:
	//   childKey = parse256(Il) + parentKey
	//
	// For public children:
	//   childKey = serP(point(parse256(Il)) + parentKey)
	var isPrivate bool
	var childKey []byte
	if k.isPrivate {
		// Case #1 or #2.
		// Add the parent private key to the intermediate private key to
		// derive the final child key.
		//
		// childKey = parse256(Il) + parenKey
		keyNum := new(big.Int).SetBytes(k.key)
		ilNum.Add(ilNum, keyNum)
		ilNum.Mod(ilNum, btcec.S256().N)
		// the bytes function here returns a minimum-length buffer to represent a big-endian
		// value. It actually works to trim leading zeros, but we definitely don't want that
		// as we have an assumption that all keys are the same length. So we might need to
		// put the zeros back on the front.
		childKey = ilNum.Bytes()
		if len(childKey) < 32 {
			buf := make([]byte, 32)
			offset := 32 - len(childKey)
			for i := range childKey {
				buf[i+offset] = childKey[i]
			}
			childKey = buf
		}
		isPrivate = true
	} else {
		// Case #3.
		// Calculate the corresponding intermediate public key for
		// intermediate private key.
		ilx, ily := btcec.S256().ScalarBaseMult(il)
		if ilx.Sign() == 0 || ily.Sign() == 0 {
			return nil, ErrInvalidChild
		}

		// Convert the serialized compressed parent public key into X
		// and Y coordinates so it can be added to the intermediate
		// public key.
		pubKey, err := btcec.ParsePubKey(k.key, btcec.S256())
		if err != nil {
			return nil, err
		}

		// Add the intermediate public key to the parent public key to
		// derive the final child key.
		//
		// childKey = serP(point(parse256(Il)) + parentKey)
		childX, childY := btcec.S256().Add(ilx, ily, pubKey.X, pubKey.Y)
		pk := btcec.PublicKey{Curve: btcec.S256(), X: childX, Y: childY}
		childKey = pk.SerializeCompressed()
	}

	// The fingerprint of the parent for the derived child is the checksum24
	// of the SHA256(parentPubKey).
	parentFP := fingerprint(k.PubKeyBytes())
	return NewExtendedKey(childKey, childChainCode, parentFP,
		k.depth+1, i, isPrivate), nil
}

// Public returns a new extended public key from this extended private key.  The
// same extended key will be returned unaltered if it is already an extended
// public key.
//
// As the name implies, an extended public key does not have access to the
// private key, so it is not capable of signing transactions or deriving
// child extended private keys.  However, it is capable of deriving further
// child extended public keys.
func (k *ExtendedKey) Public() (*ExtendedKey, error) {
	// Already an extended public key.
	if !k.isPrivate {
		return k, nil
	}

	// Convert it to an extended public key.  The key for the new extended
	// key will simply be the pubkey of the current extended private key.
	//
	// This is the function N((k,c)) -> (K, c) from [BIP32].
	return NewExtendedKey(k.PubKeyBytes(), k.chainCode, k.parentFP,
		k.depth, k.childNum, false), nil
}

// HardenedChild returns the n'th hardened child of the given extended key.
//
// The parent key must be a private key.
// A HardenedChild is guaranteed to have been derived from a private key.
// It is an error if the given key is already a hardened key.
func (k *ExtendedKey) HardenedChild(n uint32) (*ExtendedKey, error) {
	ndx := n + HardenedKeyStart
	nk, err := k.Child(ndx)
	return nk, err
}

// ECPubKey converts the extended key to a btcec public key and returns it.
func (k *ExtendedKey) ECPubKey() (*btcec.PublicKey, error) {
	return btcec.ParsePubKey(k.PubKeyBytes(), btcec.S256())
}

// ECPrivKey converts the extended key to a btcec private key and returns it.
//
// As you might imagine this is only possible if the extended key is a private
// extended key (as determined by the IsPrivate function).  The ErrNotPrivExtKey
// error will be returned if this function is called on a public extended key.
func (k *ExtendedKey) ECPrivKey() (*btcec.PrivateKey, error) {
	if !k.isPrivate {
		return nil, ErrNotPrivExtKey
	}

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), k.key)
	return privKey, nil
}

// SPubKey converts the extended key to a signature.PublicKey and returns it.
func (k *ExtendedKey) SPubKey() (*signature.PublicKey, error) {
	pub, err := k.Public()
	if err != nil {
		return nil, err
	}
	sk, err := pub.AsSignatureKey()
	if err != nil {
		return nil, err
	}
	return sk.(*signature.PublicKey), err
}

// SPrivKey converts the extended key to a signature.PrivateKey and returns it.
//
// As you might imagine this is only possible if the extended key is a private
// extended key (as determined by the IsPrivate function).  The ErrNotPrivExtKey
// error will be returned if this function is called on a public extended key.
func (k *ExtendedKey) SPrivKey() (*signature.PrivateKey, error) {
	if !k.isPrivate {
		return nil, ErrNotPrivExtKey
	}

	sk, err := k.AsSignatureKey()
	if err != nil {
		return nil, err
	}
	return sk.(*signature.PrivateKey), err
}

const extraLen = 1 + 3 + 4 + 32

// extra serializes all extra data associated with this key
func (k *ExtendedKey) extra() []byte {
	var childNumBytes [4]byte
	binary.BigEndian.PutUint32(childNumBytes[:], k.childNum)

	// The serialized format is:
	//   field | width (bytes) | notes
	//   ------|-------
	//   depth | 1
	//   parent fingerprint | 3
	//   child num | 4 | serialized as big-endian uint32
	//   chain code | 32
	serializedBytes := make([]byte, 0, extraLen)
	serializedBytes = append(serializedBytes, k.depth)
	serializedBytes = append(serializedBytes, k.parentFP...)
	serializedBytes = append(serializedBytes, childNumBytes[:]...)
	serializedBytes = append(serializedBytes, k.chainCode...)

	return serializedBytes
}

// parseExtra parses extra data associated with this key
func (k *ExtendedKey) parseExtra(data []byte) error {
	// The serialized format is:
	//   field | width (bytes) | notes
	//   ------|-------
	//   depth | 1
	//   parent fingerprint | 3
	//   child num | 4 | serialized as big-endian uint32
	//   chain code | 32
	if len(data) < extraLen {
		return errors.New("cannot parseExtra: too few bytes in data")
	}
	k.depth = data[0]
	k.parentFP = data[1:4]
	k.childNum = binary.BigEndian.Uint32(data[4:8])
	k.chainCode = data[8:40]

	return nil
}

// zero sets all bytes in the passed slice to zero.  This is used to
// explicitly clear private key material from memory.
func zero(b []byte) {
	lenb := len(b)
	for i := 0; i < lenb; i++ {
		b[i] = 0
	}
}

// Zero manually clears all fields and bytes in the extended key.  This can be
// used to explicitly clear key material from memory for enhanced security
// against memory scraping.  This function only clears this particular key and
// not any children that have already been derived.
func (k *ExtendedKey) Zero() {
	zero(k.key)
	zero(k.pubKey)
	zero(k.chainCode)
	zero(k.parentFP)
	k.key = nil
	k.depth = 0
	k.childNum = 0
	k.isPrivate = false
}

// NewMaster creates a new master node for use in creating a hierarchical
// deterministic key chain.  The seed must be between 128 and 512 bits and
// should be generated by a cryptographically secure random generation source.
//
// NOTE: There is an extremely small chance (< 1 in 2^127) the provided seed
// will derive to an unusable secret key.  The ErrUnusable error will be
// returned if this should occur, so the caller must check for it and generate a
// new seed accordingly.
func NewMaster(seed []byte) (*ExtendedKey, error) {
	secretKey, chainCode, err := bip32.NewMaster(seed)
	if err != nil {
		return nil, err
	}

	parentFP := []byte{0x00, 0x00, 0x00}
	return NewExtendedKey(secretKey[:], chainCode[:], parentFP, 0, 0, true), nil
}

// GenerateSeed returns a cryptographically secure random seed that can be used
// as the input for the NewMaster function to generate a new master node.
//
// The length is in bytes and it must be between 16 and 64 (128 to 512 bits).
// The recommended length is 32 (256 bits) as defined by the RecommendedSeedLen
// constant.
func GenerateSeed(length uint8) ([]byte, error) {
	return bip32.GenerateSeed(length, nil)
}

// Bytes returns the bytes of the key
func (k *ExtendedKey) Bytes() []byte {
	return k.key
}

// FromSignatureKey attempts to construct an ExtendedKey from a signature.Key instance
func (k *ExtendedKey) FromSignatureKey(key signature.Key) (err error) {
	k.Zero()

	if signature.NameOf(key.Algorithm()) != signature.NameOf(signature.Secp256k1) {
		err = fmt.Errorf(
			"ExtendedKey must use %s algorithm; provided key uses %s",
			signature.NameOf(signature.Secp256k1),
			signature.NameOf(key.Algorithm()),
		)
		return
	}

	k.isPrivate = signature.IsPrivate(key)
	k.key = key.KeyBytes()
	err = k.parseExtra(key.ExtraBytes())
	return
}

// FromSignatureKey attempts to construct an ExtendedKey from a signature.Key instance
func FromSignatureKey(key signature.Key) (ek *ExtendedKey, err error) {
	ek = new(ExtendedKey)
	err = ek.FromSignatureKey(key)
	return
}

// AsSignatureKey converts this ExtendedKey into a signature.Key instance
func (k ExtendedKey) AsSignatureKey() (signature.Key, error) {
	if k.isPrivate {
		priv, err := signature.RawPrivateKey(signature.Secp256k1, k.key, k.extra())
		return priv, err
	}
	pub, err := signature.RawPublicKey(signature.Secp256k1, k.key, k.extra())
	return pub, err
}

// MarshalText implements encoding.TextMarshaler
func (k ExtendedKey) MarshalText() ([]byte, error) {
	key, err := k.AsSignatureKey()
	if err != nil {
		return nil, err
	}
	return key.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler
func (k *ExtendedKey) UnmarshalText(text []byte) (err error) {
	if !utf8.Valid(text) {
		return errors.New("text not valid utf-8")
	}

	k.Zero()
	key, err := signature.ParseKey(string(text))
	if err != nil {
		return err
	}
	return k.FromSignatureKey(key)
}
