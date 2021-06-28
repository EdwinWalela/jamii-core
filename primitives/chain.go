package primitives

import (
	"time"

	"github.com/edwinwalela/jamii-core/jcrypto"
	"github.com/katzenpost/core/crypto/eddsa"
)

type Chain struct {
	/**
		Represent the chain of blocks
		node reads chain from memory and initializes chain or requests chain from
		peers through a TCP socket connection

		Proof of work will be determined by a moving average which will be determined by
		running an initial system check to determine an effienct difficulty depenging on
		the node's performance

		Target is to produce a block in 4 minutes
	**/
	Chain        []Block // list of verified blocks
	PendingVotes []Vote  // list of unverified votes recieved from clients via TCP socket
	Height       uint64  // current height of the chain
	Difficulty   uint64  // node's proof of work difficulty
}

func (c *Chain) Genesis() Block {
	blk := Block{votes: []Vote{}, prevHash: "", timestamp: uint64(time.Now().Unix()), difficulty: c.Difficulty}

	return blk
}

func (c *Chain) Init() error {
	tx := Vote{
		address:    eddsa.PublicKey{},
		candidates: []eddsa.PublicKey{},
		signature:  []byte(""),
		hash:       "",
		timestamp:  0,
	}

	kp := jcrypto.KeyPair{}
	if err := jcrypto.GenKeyPair(&kp); err != nil {
		return err
	}

	c.AddTX(tx)
	c.Mine(&kp)
	return nil

}

func (c *Chain) AddTX(tx Vote) {
	c.PendingVotes = append(c.PendingVotes, tx)
}

func (c *Chain) LatestBlock() Block {
	if len(c.Chain) != 0 {
		return c.Chain[len(c.Chain)-1]
	}
	return c.Genesis()
}

func (c *Chain) Mine(kp *jcrypto.KeyPair) {
	blk := &Block{nonce: 0, difficulty: c.Difficulty}
	now := uint64(time.Now().Unix())
	candidates := []eddsa.PublicKey{}

	voteBase := &Vote{address: kp.PubKey, candidates: candidates, timestamp: now}

	blk.AddVote(*voteBase) // Add vote

	for _, v := range c.PendingVotes {
		blk.AddVote(v) // add all pending votes to new block
	}

	blk.SetPreviousHash(c.LatestBlock().hash) // Set previous hash

	blk.Hash() // Hash block

	c.Chain = append(c.Chain, *blk) // append new block

	c.PendingVotes = []Vote{} // empty pending votes

	// Alert peers to stop mining, broadcast new block

}
