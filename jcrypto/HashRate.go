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
	// log.Printf("Attempting PoW with difficulty of %d", Difficulty)
	digest := SHA512(SOURCE + fmt.Sprintf("%d", nonce))

	tStart := time.Now().Unix()

	for !HashValid(digest) {
		nonce++
		digest = SHA512(SOURCE + fmt.Sprintf("%d", nonce))
	}

	tEnd := time.Now().Unix()

	elapsed := tEnd - tStart
	Difficulty++
	// log.Printf("Operation completed in %d seconds, diff:%d", elapsed, Difficulty)

	if elapsed >= TARGET_TIME {
		return Difficulty, nonce
	} else {
		return FindDifficulty()
	}
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
