package primitives

import "crypto/ecdsa"

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
