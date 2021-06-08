package main

import (
	"fmt"

	"github.com/edwinwalela/jamii-core/primitives"
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

	// kp := &jcrypto.KeyPair{}

	// decodedPub, err := base64.StdEncoding.DecodeString("oMIwjMspTa8oBTsc/0qJ/GUvc6Fa/MX78F2lwkNdePs=")

	// decodedSig, sErr := base64.StdEncoding.DecodeString("rgBcyHuuyeB9CgwtiG/+PG2JqRZLf3lPpsD5fSgANnQlQBrX78db874Xys5f/JnjZccXVMaUvqpzO3F/utatDw==")

	// if err != nil {
	// 	log.Println(err)
	// }
	// if sErr != nil {
	// 	log.Println(err)
	// }

	// jcrypto.PubKeyFromBytes(decodedPub, kp)

	// valid := jcrypto.VerifySig(decodedSig, []byte("hello world"), kp)

	// fmt.Println(valid)
	v := &primitives.Vote{}
	clientData := "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f|Z/Yg3ETYDRBlXcDy09p/NOyQRGycYRU8kNPsTRkwqRYpiz8ljL87u8fg/x/xzpELh65Af95kIeV2zoV++CbiDw==|rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=|rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=,rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=,rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=,rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=|1623142046"
	v.UnpackClientString(clientData)

	fmt.Println("Address: ", v.Address)
	fmt.Println("Candidates: ", v.Candidates)
	fmt.Println("Signature: ", v.Signature)
	fmt.Println("Hash:", v.Hash)
	fmt.Println("Timestamp: ", v.Timestamp)

}
