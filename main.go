package main

import "github.com/edwinwalela/jamii-core/jcrypto"

func main() {
	/** Key pair generation and signing **/

	kp := &jcrypto.KeyPair{}

	jcrypto.GenKeyPair(kp, "") // Generate key pair

	// hash := jcrypto.SHA512("hello world") // hash data

	// signature, err := kp.Sign(hash) // sign hash

	// if err != nil {
	// 	log.Println("signing failed:", err)
	// }

	// valid := kp.Verify([]byte(hash), signature) // verify hash with signature

	// fmt.Println(valid)

	// v := &primitives.Vote{}
	// v.UnpackClientString("pubkeyhex|candidate,candidate,candidate|signature|timestamp")

}
