package primitives

import (
	"crypto/ecdsa"
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
	string data from client ->	pubkeyhex|candidate,candidate,candidate|signature|timestamp
	This is broken down by node
**/

// func (v *Vote) UnpackClientString(data string) {
// 	vote := strings.Split(data, "|")
// 	for i, val := range vote {
// 		switch i{
// 		case 0:
// 			v.Address
// 		}
// 	}

// }
