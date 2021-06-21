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

	address    eddsa.PublicKey   // public key of tx initiator
	candidates []eddsa.PublicKey // list of candidtes selected by initiator
	signature  []byte            // signature by the initiator's public key
	hash       string            // hash of the tx
	timestamp  uint64            // Unix timestamp of tx in seconds
}

/**
	string data from client ->	pubkey64|sign64|hash|candidate,candidate,candidate|timestamp
	This is broken down by node
**/

func (v *Vote) GetAddress() eddsa.PublicKey {
	return v.address
}

func (v *Vote) GetSignature() []byte {
	return v.signature
}

func (v *Vote) GetHash() string {
	return v.hash
}

func (v *Vote) SetCandidates(candidates []eddsa.PublicKey) {
	v.candidates = candidates
}

func (v *Vote) Hash() string {
	_hash := ""

	_hash += v.address.String()

	for _, candidate := range v.candidates {
		_hash += candidate.String()
	}

	_hash += fmt.Sprintf("%d", v.timestamp)

	return jcrypto.SHA512(_hash)
}

func (v *Vote) isValid() bool {
	return jcrypto.VerifySig(v.signature, []byte(v.hash), v.address)
}
