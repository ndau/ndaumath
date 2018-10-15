package signature

import (
	"bytes"
	"crypto/sha256"
)

const (
	// ChecksumPadWidth is the moduland of which the length of a checksummed
	// byte slice will always equal 0.
	//
	// In other words, we choose a checksum width such that
	// `len(summedMsg) % ChecksumPadWidth == 0`.
	//
	// We choose 5 because we expect to encode this data in a base32 encoding,
	// which needs no padding when the input size is a multiple of 5.
	ChecksumPadWidth = 5
	// ChecksumMinBytes is the minimum number of checksum bytes. The chance
	// that a checksum will accidentally pass is roughly `1 / (256 ^ ChecksumMinBytes)`,
	// so the value of 3 used gives a false positive rate of about 1 / 16 million.
	ChecksumMinBytes = 3
)

func checksumWidth(inputLen int) byte {
	inputLen++ // input will have an additional byte prefixed with this width
	fill := ChecksumPadWidth - byte(inputLen%ChecksumPadWidth)
	if fill < ChecksumMinBytes {
		fill += ChecksumPadWidth
	}
	return fill
}

// return the trailing n bytes of the sha224 checksum of the input bytes
// given that sha-224 is already a simple truncation of sha-256, this means
// that the returned bytes are from the middle of a sha-256 checksum
func cksumN(input []byte, n byte) []byte {
	sum := sha256.Sum224(input)
	return sum[sha256.Size224-int(n):]
}

// AddChecksum adds a checksum to a byte slice.
//
// The number of bytes of the checksum depend on the width of the data:
// an appropriate number will be used, at least `ChecksumMinBytes`,
// such that the total length of the returned slice is a multiple of
// `ChecksumPadWidth`.
//
// The checksum bytes are appended to the end of the message.
//
// A single byte is also added at the head of the byte slice, containing
// the number of padding bytes, for ease of checking the checksum.
func AddChecksum(bytes []byte) []byte {
	n := checksumWidth(len(bytes))
	sum := cksumN(bytes, n)

	out := make([]byte, 1, 1+len(bytes)+int(n))
	out[0] = n
	out = append(out, bytes...)
	out = append(out, sum...)
	return out
}

// CheckChecksum validates the checksum of a summed byte slice.
//
// It returns the wrapped message stripped of checksum data and
func CheckChecksum(checked []byte) (message []byte, checksumOk bool) {
	if len(checked) < 1+ChecksumMinBytes {
		return nil, false
	}
	n := checked[0]
	message = checked[1 : len(checked)-int(n)]

	sumActual := checked[len(checked)-int(n):]
	sumExpect := cksumN(message, n)
	checksumOk = bytes.Equal(sumActual, sumExpect)
	return
}
