package b32

import (
	"crypto/sha256"

	"github.com/sigurn/crc16"
)

// The CRC16 polynomial used is AUG_CCITT: `0x1021`
var ndauTable = crc16.MakeTable(crc16.CRC16_AUG_CCITT)

// Checksum generates a 2-byte checksum of b.
func Checksum16(b []byte) []byte {
	ck := crc16.Checksum(b, ndauTable)
	return []byte{byte((ck >> 8) & 0xFF), byte(ck & 0xFF)}
}

// Checksum24 generates a 3-byte checksum of buf.
func Checksum24(buf []byte) []byte {
	hasher := sha256.New()
	hasher.Write(buf)
	b := hasher.Sum(nil)
	return b[:3]
}

// Check accepts an array of bytes and a 2-byte checksum and returns true if the checksum
// of b is equal to the value passed in.
func Check(b []byte, ckb []byte) bool {
	ck := crc16.Checksum(b, ndauTable)
	return byte((ck>>8)&0xFF) == ckb[0] && byte(ck&0xFF) == ckb[1]
}
