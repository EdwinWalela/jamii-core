package primitives

import (
	"crypto/ed25519"
	"fmt"

	"github.com/edwinwalela/jamii-core/jcrypto"
)

type Vote struct {
	/**
	Vote represents a tx in the blockchain
	Mobile clients submits vote via TCP socket connection & node packs it into a Vote
	**/

	Address    ed25519.PublicKey // public key of tx initiator
	Candidates []string          // list of candidtes selected by initiator
	Signature  []byte            // signature by the initiator's public key
	Hash       string            // hash of the tx
	Timestamp  uint64            // Unix timestamp of tx in seconds
}

/**
	string data from client ->	pubkey64|sign64|hash|candidate,candidate,candidate|timestamp
	This is broken down by node
**/

func (v *Vote) HashVote() string {
	_hash := ""

	_hash += string(v.Address)

	for _, candidate := range v.Candidates {
		_hash += candidate
	}

	_hash += fmt.Sprintf("%d", v.Timestamp)

	return jcrypto.SHA512(_hash)
}

func (v *Vote) IsValid() bool {

	return jcrypto.VerifySig(v.Signature, []byte(v.Hash), v.Address)
}
