package hash

import "crypto/sha1"

// FixedKey is the 16-byte magic key used in all hash derivations.
var FixedKey = []byte{
	0x77, 0x21, 0x4d, 0x4b, 0x19, 0x6a, 0x87, 0xcd,
	0x52, 0x00, 0x45, 0xfd, 0x20, 0xa5, 0x1d, 0x67,
}

// Computes the 3-pass SHA-1 derivation for a 4-byte candidate.
// It uses pre-allocated buffers buf1 (20 bytes), buf2 (40 bytes), buf3 (32 bytes)
// to avoid heap allocations in hot loops.
func Audible(candidate uint32, buf1, buf2, buf3 []byte) [20]byte {
	buf1[16] = byte(candidate >> 24)
	buf1[17] = byte(candidate >> 16)
	buf1[18] = byte(candidate >> 8)
	buf1[19] = byte(candidate)

	res1 := sha1.Sum(buf1)

	copy(buf2[16:36], res1[:])
	buf2[36] = buf1[16]
	buf2[37] = buf1[17]
	buf2[38] = buf1[18]
	buf2[39] = buf1[19]
	res2 := sha1.Sum(buf2)

	copy(buf3[0:16], res1[:16])
	copy(buf3[16:32], res2[:16])
	return sha1.Sum(buf3)
}

// NewBuffers creates the 3 pre-allocated buffers used.
// The first 16 bytes of buf1 and buf2 are pre-filled with FixedKey.
func NewBuffers() (buf1, buf2, buf3 []byte) {
	buf1 = make([]byte, 20)
	copy(buf1, FixedKey)
	buf2 = make([]byte, 40)
	copy(buf2, FixedKey)
	buf3 = make([]byte, 32)
	return
}
