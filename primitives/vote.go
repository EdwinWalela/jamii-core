package primitives

import (
	"fmt"

	"github.com/edwinwalela/jamii-core/jcrypto"
	"github.com/katzenpost/core/crypto/eddsa"
)

type Vote struct {
	/**
	Vote represents a tx in the blockchain
	Mobile clients submits vote via TCP socket connection & node packs it into a Vote
	**/

	Address    eddsa.PublicKey   // public key of tx initiator
	Candidates []eddsa.PublicKey // list of candidtes selected by initiator
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

	_hash += v.Address.String()

	for _, candidate := range v.Candidates {
		_hash += candidate.String()
	}

	_hash += fmt.Sprintf("%d", v.Timestamp)

	return jcrypto.SHA512(_hash)
}

func (v *Vote) isValid() bool {
	return jcrypto.VerifySig(v.Signature, []byte(v.Hash), v.Address)
}
