package primitives

import (
	"bytes"
	"crypto/ed25519"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
		duplicate := 0
		for _, b := range c.Chain {
			for _, tx := range b.Votes {
				if tx.Address.Equal(v.Address) {

					duplicate += 1
					// add pending vote to new block
				}
			}
		}
		if duplicate < 2 {
			blk.AddVote(v)
		} else {
			log.Println("Duplicate address found, discarding vote")
		}
	}
	log.Printf("blk size: %d\n", len(blk.Votes))

	if len(blk.Votes) < 2 {
		log.Println("Minimum block size not met, discarding block")
		log.Println("Mining terminated")
		return nil
	}

	blk.PrevHash = (c.LatestBlock().Hash) // Set previous hash

	log.Println("Calculating block hash")
	blk.HashBlk() // Hash block

	log.Printf("Hash: %s......%s", blk.Hash[:12], blk.Hash[len(blk.Hash)-12:len(blk.Hash)-1])
	log.Printf("Hash found after %d tries (nonce) with difficulty %d\n", blk.Nonce, blk.Difficulty)

	c.Chain = append(c.Chain, *blk) // append new block

	log.Println("Writing block to file")
	log.Printf("Current chain length: %d\n", len(c.Chain))
	if err := c.writeBlock(blk); err != nil {
		return err
	} // Write to file system

	log.Println("Emptying pending votes")
	c.PendingVotes = []Vote{} // empty pending votes

	// Alert peers to stop mining, broadcast new block

	return nil
}

func (c *Chain) Result() string {
	res := ""
	var presidential = make(map[string]int)
	var county = make(map[string]int)
	var parliamentary = make(map[string]int)

	for _, b := range c.Chain {
		for _, v := range b.Votes {
			for _, c := range v.Candidates {

				if strings.Contains(c, "Presidential") {

					if count, exists := presidential[c]; !exists {
						presidential[c] = 1
					} else {
						presidential[c] = count + 1
					}
				} else if strings.Contains(c, "Parliamentary") {
					if count, exists := parliamentary[c]; !exists {
						parliamentary[c] = 1
					} else {
						parliamentary[c] = count + 1
					}
				} else if strings.Contains(c, "County") {
					if count, exists := county[c]; !exists {
						county[c] = 1
					} else {
						county[c] = count + 1
					}
				}
			}
		}
	}
	highestPresCount := 0
	topPres := ""
	highestParlCount := 0
	topParl := ""
	highestCountyCount := 0
	topCounty := ""

	for k, v := range presidential {
		if v > highestPresCount {
			highestPresCount = v
			topPres = k
		}

	}
	res += topPres + "-" + strconv.Itoa(highestPresCount) + "|"
	for k, v := range parliamentary {
		if v > highestParlCount {
			highestParlCount = v
			topParl = k
		}
	}
	res += topParl + "-" + strconv.Itoa(highestParlCount) + "|"
	for k, v := range county {
		if v > highestCountyCount {
			highestCountyCount = v
			topCounty = k
		}
	}
	res += topCounty + "-" + strconv.Itoa(highestCountyCount) + "|"

	return res
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
