package primitives

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"log"
	"math/big"
	"strings"
)

type Vote struct {
	/**
	Vote represents a tx in the blockchain
	Mobile clients submits vote via TCP socket connection & node packs it into a Vote
	**/

	Address    ecdsa.PublicKey   // public key of tx initiator
	Candidates []ecdsa.PublicKey // list of candidtes selected by initiator
	Signature  []byte            // signature by the initiator's public key
	Hash       string            // hash of the tx
	Timestamp  uint64            // Unix timestamp of tx in seconds
}

/**
	string data from client ->	x,y|candidate,candidate,candidate|signature|timestamp
	This is broken down by node
**/

func (v *Vote) UnpackClientString(data string) {
	vote := strings.Split(data, "|")
	ecdsaX := new(big.Int)
	ecdsaY := new(big.Int)

	for i, val := range vote {
		switch i {

		case 0: // Extract curve X value
			ecdsaX, ok := ecdsaX.SetString(val, 10)
			if !ok {
				log.Println("SetString error")
			}
			fmt.Println("X-val:", ecdsaX)

		case 1: // Extract curve Y value
			ecdsaY, ok := ecdsaY.SetString(val, 10)
			if !ok {
				log.Println("SetString error")
			}
			fmt.Println("Y-val:", ecdsaY)
		}
	}

	v.Address = ecdsa.PublicKey{X: ecdsaX, Y: ecdsaY, Curve: elliptic.P256()}

}
