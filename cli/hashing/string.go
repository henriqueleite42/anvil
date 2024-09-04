package hashing

import "crypto/sha1"

func String(str string) string {
	hashByte := sha1.Sum([]byte(str))
	return string(hashByte[:])
}
