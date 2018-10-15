package bip32

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

const (
	// RecommendedSeedLen is the recommended length in bytes for a seed
	// to a master node.
	RecommendedSeedLen = 32 // 256 bits

	// HardenedKeyStart is the index at which a hardened key starts.  Each
	// extended key has 2^31 normal child keys and 2^31 hardned child keys.
	// Thus the range for normal child keys is [0, 2^31 - 1] and the range
	// for hardened child keys is [2^31, 2^32 - 1].
	HardenedKeyStart = 0x80000000 // 2^31

	// MinSeedBytes is the minimum number of bytes allowed for a seed to
	// a master node.
	MinSeedBytes = 16 // 128 bits

	// MaxSeedBytes is the maximum number of bytes allowed for a seed to
	// a master node.
	MaxSeedBytes = 64 // 512 bits
)

var (
	// ErrInvalidSeedLen describes an error in which the provided seed or
	// seed length is not in the allowed range.
	ErrInvalidSeedLen = fmt.Errorf(
		"seed length must be between %d and %d bits",
		MinSeedBytes*8, MaxSeedBytes*8,
	)

	// ErrUnusableSeed describes an error in which the provided seed is not
	// usable due to the derived key falling outside of the valid range for
	// secp256k1 private keys.  This error indicates the caller must choose
	// another seed.
	ErrUnusableSeed = errors.New("unusable seed")
)

// masterKey is the master key used along with a random seed used to generate
// the master node in the hierarchical tree.
var masterKey = []byte("ndau seed")

// NewMaster creates a master secret key and a master chain code per the
// procedure described [in BIP32][bip32-master].
//
// [bip32]: https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
// [bip32-master]: https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#Master_key_generation
func NewMaster(seed []byte) (Il, Ir [32]byte, err error) {
	// Per [bip32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if len(seed) < MinSeedBytes || len(seed) > MaxSeedBytes {
		err = ErrInvalidSeedLen
		return
	}

	// First take the HMAC-SHA512 of the master key and the seed data:
	//   I = HMAC-SHA512(Key = "ndau seed", Data = S)
	hmac512 := hmac.New(sha512.New, masterKey)
	hmac512.Write(seed)
	I := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = master secret key
	//   Ir = master chain code
	size := copy(Il[:], I[:len(I)/2])
	if size != len(Il) {
		panic("programming error in NewMaster: Il")
	}
	size = copy(Ir[:], I[len(I)/2:])
	if size != len(Ir) {
		panic("programming error in NewMaster: Ir")
	}

	// Ensure the key in usable.
	secretKeyNum := new(big.Int).SetBytes(Il[:])
	if secretKeyNum.Cmp(btcec.S256().N) >= 0 || secretKeyNum.Sign() == 0 {
		err = ErrUnusableSeed
	}
	return
}

// GenerateSeed returns a cryptographically secure random seed that can be used
// as the input for the NewMaster function to generate a new master node.
//
// The length is in bytes and it must be between 16 and 64 (128 to 512 bits).
// The recommended length is 32 (256 bits) as defined by the RecommendedSeedLen
// constant.
//
// `rng` should be the best available source of entropy. If nil,
// crypto/rand.Reader is used.
func GenerateSeed(length uint8, rng io.Reader) ([]byte, error) {
	// Per [BIP32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if length < MinSeedBytes || length > MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	// if rand is not set, set it
	if rng == nil {
		rng = rand.Reader
	}

	// generate a seed of the recommended size
	seed := make([]byte, length)
	_, err := io.ReadFull(rng, seed)
	if err != nil {
		return seed, err
	}

	return seed, nil
}

// PrivateToPublic implements the Private -> Public derivation
func PrivateToPublic(private []byte) []byte {
	pkx, pky := btcec.S256().ScalarBaseMult(private)
	pubKey := btcec.PublicKey{Curve: btcec.S256(), X: pkx, Y: pky}
	return pubKey.SerializeCompressed()
}
