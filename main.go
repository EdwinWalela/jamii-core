package main

import (
	"flag"
	"fmt"

	"github.com/edwinwalela/jamii-core/net/peer"
	"github.com/edwinwalela/jamii-core/net/server"
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

	/** Vote Unpacking **/
	// v := &primitives.Vote{}
	// clientData := "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f|Z/Yg3ETYDRBlXcDy09p/NOyQRGycYRU8kNPsTRkwqRYpiz8ljL87u8fg/x/xzpELh65Af95kIeV2zoV++CbiDw==|rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=|rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=,rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=,rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=,rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=|1623142046"
	// v.UnpackClientString(clientData)

	// fmt.Println("Address: ", v.Address)
	// fmt.Println("Candidates: ", v.Candidates)
	// fmt.Println("Signature: ", v.Signature)
	// fmt.Println("Hash:", v.Hash)
	// fmt.Println("Timestamp: ", v.Timestamp)

	exit := make(chan int)
	var connectedPeers []peer.Peer

	localPortPtr := flag.Int("local", 3000, "local Websocket server port")
	flag.Parse()

	server := &server.Server{Host: "localhost", Network: "tcp", Port: *localPortPtr}
	// Initialize local server
	if err := server.Init(); err != nil {
		fmt.Println("Unable to create TCP connection\n", err)
	}
	go func() {
		server.Accept(&connectedPeers)
	}()

	v := &primitives.Vote{}
	clientData := "ea301ceb1ba7df7be32b573fb9c6dc8c7e2121a9d7ecee1d78a94aa201145f39c44b766144bd354b348c7ab7f2d681f3df6095ec577e3cee11a29dea115ee815|vfeMgOgZ6q4QUxTNzvG3/iQKnfCSiRYKiUoexXNEOo7XQJxX1/57L7A1XcLSp6XSIQ70XSYXvYcHCC0cY5qnBA==|StpW0TTiB2G3vaHyRfF35sqhA7misfUw7Uj7lsVK1Hs=|StpW0TTiB2G3vaHyRfF35sqhA7misfUw7Uj7lsVK1Hs=.StpW0TTiB2G3vaHyRfF35sqhA7misfUw7Uj7lsVK1Hs=.StpW0TTiB2G3vaHyRfF35sqhA7misfUw7Uj7lsVK1Hs=.StpW0TTiB2G3vaHyRfF35sqhA7misfUw7Uj7lsVK1Hs=|1623308754005"

	if isValid, err := v.UnpackClientString(clientData); err != nil {
		fmt.Println(err)
	} else if isValid {
		fmt.Println("Vote accepted")
	}

	// cli.MainMenu(&connectedPeers)

	// Generate Key Pair

	// Attempt to connect to peers

	// Read Blocks from memory

	// If no blocks request from peers

	// Check chain

	// Initalize chain

	exit <- 1

}
