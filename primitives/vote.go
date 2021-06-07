package primitives

import "crypto/ecdsa"

type Vote struct {
	/**
	Vote represents a tx in the blockchain
	Mobile clients submits vote via TCP socket connection & node packs it into a Vote

	Genesis block will contain a single Vote with the candidates field indicating registered
	election candidtes. Address of genesis block will be null

	Genesis block data will be obtained from institution's webserver and sent to nodes prior
	to election
	**/

	address    ecdsa.PublicKey   // public key of tx initiator
	candidates []ecdsa.PublicKey // list of candidtes selected by initiator
	signature  []byte            // signature by the initiator's public key
	hash       string            // hash of the tx
	timestamp  uint64            // Unix timestamp of tx in seconds
}
