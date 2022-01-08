package primitives

import (
	"bytes"
	"crypto/ed25519"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/edwinwalela/jamii-core/jcrypto"
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
		Address:    ed25519.PublicKey{},
		Candidates: []string{},
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
	log.Println("Mining triggered")
	blk := &Block{Nonce: 0, Difficulty: c.Difficulty}
	blk.Timestamp = uint64(time.Now().Unix())
	candidates := []string{}

	voteBase := &Vote{Address: kp.PubKey, Candidates: candidates, Timestamp: blk.Timestamp}

	blk.AddVote(*voteBase) // Add vote

	log.Printf("packaging %d pending vote(s) to new block \n", len(c.PendingVotes))
	for _, v := range c.PendingVotes {
		blk.AddVote(v) // add all pending votes to new block
	}

	blk.PrevHash = (c.LatestBlock().Hash) // Set previous hash

	log.Println("Calculating block hash")
	blk.HashBlk() // Hash block

	log.Printf("Hash: %s......%s", blk.Hash[:12], blk.Hash[len(blk.Hash)-12:len(blk.Hash)-1])
	log.Printf("Hash found after %d tries (nonce) with difficulty %d\n", blk.Nonce, blk.Difficulty)

	c.Chain = append(c.Chain, *blk) // append new block

	log.Println("Writing block to file")
	log.Printf("Current chain length %d\n", len(c.Chain))
	if err := c.writeBlock(blk); err != nil {
		return err
	} // Write to file system

	log.Println("Emptying pending votes")
	c.PendingVotes = []Vote{} // empty pending votes

	// Alert peers to stop mining, broadcast new block

	return nil
}

func (c *Chain) writeBlock(blk *Block) error {

	var blockdump string

	blockdump = fmt.Sprintf("%d,%s,%s,%d,%d|", blk.Difficulty, blk.Hash, blk.PrevHash, blk.Timestamp, blk.Nonce)

	for _, vote := range blk.Votes {
		blockdump += fmt.Sprintf("%s,%s,%s,%d|", string(vote.Address), vote.Hash, vote.Signature, vote.Timestamp)
		for _, candidate := range vote.Candidates {
			blockdump += fmt.Sprintf("%s,", candidate)
		}
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(blockdump); err != nil {
		return err
	}

	blockFileName := fmt.Sprintf("%d", blk.Timestamp) + ".jblock"

	path := filepath.Join(".", c.BlockDir, blockFileName)

	writer, err := os.Create(path)

	if err != nil {
		return err
	}

	if _, writerError := writer.Write(buf.Bytes()); writerError != nil {
		return writerError
	}

	// Decode
	// dec := gob.NewDecoder(buf)
	// var s2 string

	// if err := dec.Decode(&s2); err != nil {
	// 	return err
	// }

	return nil
}
