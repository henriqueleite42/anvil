package hashing

import (
	"crypto/sha1"
	"encoding/hex"
)

func String(str string) string {
	hashByte := sha1.Sum([]byte(str))
	return hex.EncodeToString(hashByte[:])
}
