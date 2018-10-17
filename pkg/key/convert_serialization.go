package key

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/oneiro-ndev/ndaumath/pkg/b32"
)

// FromOldSerialization attempts to produce an ExtendedKey from the old
// (pre-october-2018) format.
//
// If successful, it produces an ExtendedKey object
// whose MarshalText method will produce the new serialization.
func FromOldSerialization(key string) (*ExtendedKey, error) {
	const serializedKeyLen = 3 + 1 + 3 + 4 + 32 + 33

	// The base32-decoded extended key must consist of a serialized payload
	// plus an additional 4 bytes for the checksum.
	decoded, err := b32.Decode(key)
	if err != nil {
		return nil, ErrInvalidKeyEncoding
	}
	if len(decoded) != serializedKeyLen+4 {
		return nil, ErrInvalidKeyLen
	}

	// The serialized format is:
	//   version (3) || depth (1) || parent fingerprint (3)) ||
	//   child num (4) || chain code (32) || key data (33) || checksum (4)

	// Split the payload and checksum up and ensure the checksum matches.
	payload := decoded[:len(decoded)-4]
	checkSum := decoded[len(decoded)-4:]
	expectedCheckSum := doubleHashB(payload)[:4]
	if !bytes.Equal(checkSum, expectedCheckSum) {
		return nil, ErrBadChecksum
	}

	// Deserialize each of the payload fields.
	depth := payload[3:4][0]
	parentFP := payload[4:7]
	childNum := binary.BigEndian.Uint32(payload[7:11])
	chainCode := payload[11:43]
	keyData := payload[43:76]

	// The key data is a private key if it starts with 0x00.  Serialized
	// compressed pubkeys either start with 0x02 or 0x03.
	isPrivate := keyData[0] == 0x00
	if isPrivate {
		// Ensure the private key is valid.  It must be within the range
		// of the order of the secp256k1 curve and not be 0.
		keyData = keyData[1:]
		keyNum := new(big.Int).SetBytes(keyData)
		if keyNum.Cmp(btcec.S256().N) >= 0 || keyNum.Sign() == 0 {
			return nil, ErrUnusableSeed
		}
	} else {
		// Ensure the public key parses correctly and is actually on the
		// secp256k1 curve.
		_, err := btcec.ParsePubKey(keyData, btcec.S256())
		if err != nil {
			return nil, err
		}
	}

	return NewExtendedKey(keyData, chainCode, parentFP, depth, childNum, isPrivate), nil
}
