package jcrypto

import (
	"crypto/sha512"
	"encoding/hex"
)

func SHA512(src string) string {
	hash := sha512.New()
	hash.Write([]byte(src))
	digest := hex.EncodeToString(hash.Sum(nil))
	return digest
}
