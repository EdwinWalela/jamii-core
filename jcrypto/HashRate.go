package jcrypto

import (
	"fmt"
	"time"
)

var Difficulty uint64 = 1
var nonce uint64 = 0

const (
	SOURCE      = "Hello World" // Sample text to hash
	TARGET_TIME = 5             // Target HashRate in seconds
)

func FindDifficulty() (uint64, uint64) {
	digest := SHA512(SOURCE + fmt.Sprintf("%d", nonce))

	tStart := time.Now().Unix()

	for !HashValid(digest) {
		nonce++
		digest = SHA512(SOURCE + fmt.Sprintf("%d", nonce))
	}

	tEnd := time.Now().Unix()

	elapsed := tEnd - tStart
	Difficulty++

	if elapsed > 1 {
		return Difficulty, nonce
	} else {
		FindDifficulty()
	}
	return Difficulty, nonce
}

func HashValid(hash string) bool {
	for i, v := range hash {
		if i == int(Difficulty) {
			break
		}
		if string(v) == "0" {
			continue
		} else {
			return false
		}

	}
	return true
}
