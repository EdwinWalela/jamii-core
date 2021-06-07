package main

import (
	"fmt"
	"log"

	"github.com/edwinwalela/jamii-core/crypt"
)

func main() {

	kp := &crypt.KeyPair{}

	crypt.GenKeyPair(kp) // Generate key pair

	hash := crypt.SHA512("hello world") // hash data

	signature, err := kp.Sign(hash) // sign hash

	if err != nil {
		log.Println("signing failed:", err)
	}

	valid := kp.Verify([]byte(hash), signature) // verify hash with signature

	fmt.Println(valid)
}
