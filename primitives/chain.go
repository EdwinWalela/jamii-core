package primitives

import (
	"bytes"
	"encoding/gob"
	"log"
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
	BlockDir     string
}

func (c *Chain) Genesis() Block {
	blk := Block{Votes: []Vote{}, PrevHash: "", Timestamp: uint64(time.Now().Unix()), Difficulty: c.Difficulty}

	return blk
}

func (c *Chain) Init() error {
	tx := Vote{
		Address:    eddsa.PublicKey{},
		Candidates: []eddsa.PublicKey{},
		Signature:  []byte(""),
		Hash:       "",
		Timestamp:  0,
	}

	kp := jcrypto.KeyPair{}
	if err := jcrypto.GenKeyPair(&kp); err != nil {
		return err
	}

	c.AddTX(tx)

	if mineError := c.Mine(&kp); mineError != nil {
		return mineError
	}
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

func (c *Chain) Mine(kp *jcrypto.KeyPair) error {
	blk := &Block{Nonce: 0, Difficulty: c.Difficulty}
	now := uint64(time.Now().Unix())
	candidates := []eddsa.PublicKey{}

	voteBase := &Vote{Address: kp.PubKey, Candidates: candidates, Timestamp: now}

	blk.AddVote(*voteBase) // Add vote

	for _, v := range c.PendingVotes {
		blk.AddVote(v) // add all pending votes to new block
	}

	blk.PrevHash = (c.LatestBlock().Hash) // Set previous hash

	blk.HashBlk() // Hash block

	c.Chain = append(c.Chain, *blk) // append new block

	if err := c.writeBlock(blk); err != nil {
		return err
	} // Write to file system

	c.PendingVotes = []Vote{} // empty pending votes

	// Alert peers to stop mining, broadcast new block

	return nil
}

func (c *Chain) writeBlock(blk *Block) error {

	// blockFileName := fmt.Sprintf("%d", blk.GetTimestamp()) + ".jblock"

	// path := filepath.Join(".", c.BlockDir, blockFileName)

	// writer, err := os.Create(path)

	// if err != nil {
	// 	return err
	// }

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(blk); err != nil {
		return err
	}

	log.Println(buf.Bytes())

	return nil
}
