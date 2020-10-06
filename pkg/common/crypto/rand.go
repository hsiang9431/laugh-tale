package crypto

import (
	"encoding/base64"
	"math/rand"
	"time"
)

// The length of returned string is rounded to the
// largest multiple of 4 less than or equal to length
//
// length in this function is for base 64 string,
// which is one third longer than original byte array
// for utilizing the full complexity of AES 256,
// length should be no shorter than 43 charactors in base 64
func RandB64String(length int) string {
	rand.Seed(time.Now().UnixNano())
	length >>= 2          // length /= 4
	length += length << 1 // length *= 3
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)
}

// The length in this function is that of the raw byte array
// base 64 string will be one third longer than it
func RandBytesAndB64(length int) ([]byte, string) {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return bytes, base64.RawURLEncoding.EncodeToString(bytes)
}
