package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/edwinwalela/jamii-core/jcrypto"
)

func main() {
	/** Key pair generation and signing **/

	// kp := &jcrypto.KeyPair{}

	// jcrypto.GenKeyPair(kp, "") // Generate key pair

	// hash := jcrypto.SHA512("hello world") // hash data

	// signature, err := kp.Sign(hash) // sign hash

	// if err != nil {
	// 	log.Println("signing failed:", err)
	// }

	// valid := jcrypto.VerifySig(signature, []byte(hash), kp) // verify hash with signature

	// fmt.Println(valid)

	kp := &jcrypto.KeyPair{}

	decodedPub, err := base64.StdEncoding.DecodeString("oMIwjMspTa8oBTsc/0qJ/GUvc6Fa/MX78F2lwkNdePs=")

	decodedSig, sErr := base64.StdEncoding.DecodeString("rgBcyHuuyeB9CgwtiG/+PG2JqRZLf3lPpsD5fSgANnQlQBrX78db874Xys5f/JnjZccXVMaUvqpzO3F/utatDw==")

	if err != nil {
		log.Println(err)
	}
	if sErr != nil {
		log.Println(err)
	}

	jcrypto.PubKeyFromBytes(decodedPub, kp)

	valid := jcrypto.VerifySig(decodedSig, []byte("hello world"), kp)

	fmt.Println(valid)
	// v := &primitives.Vote{}
	// v.UnpackClientString("32983498|82913|candidate,candidate,candidate|signature|timestamp")

}
