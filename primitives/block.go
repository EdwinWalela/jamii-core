package primitives

import (
	"fmt"
	"strconv"

	"github.com/edwinwalela/jamii-core/jcrypto"
)

type Block struct {
	/**
	Genesis block will contain a single Vote(tx) with the candidates field indicating registered
	election candidtes. Address of genesis block will be null

	Genesis block data will be obtained from institution's webserver and sent to nodes prior
	to election

	Extra Genesis block to be included defining registered voters under candidates, voters information
	will be stored in web server and prior to election their addressess sent to nodes and packed into
	a block. The block will be queried prior to casting a vote to ensure the voter was registered
	**/
	Votes      []Vote // list of votes representing candidates/voters (if genesis block) and cast votes if not genesis
	Hash       string // hash of the block (votes hashes + block hash)
	PrevHash   string // hash of previous block in the chain
	Timestamp  uint64 // Unix timestamp of tx in seconds
	Nonce      uint64 // Proof of work tries
	Difficulty uint64 // Proof of work difficulty
}

func (b *Block) HashBlk() {
	hash := ""
	// Retrieve hashes of all votes
	for _, v := range b.Votes {
		hash += v.Hash
	}

	hash += strconv.FormatInt(int64(b.Timestamp), 10)
	hash += strconv.FormatUint(b.Difficulty, 10)
	b.Hash = jcrypto.SHA512(hash)
	for !b.HashValid() {
		b.Nonce++
		b.Hash = jcrypto.SHA512(fmt.Sprintf("%s%d", b.Hash, b.Nonce))

	}
}

func (b *Block) HashValid() bool {
	for i, v := range b.Hash {
		if i == int(b.Difficulty) {
			break
		}
		if string(v) == "0" {
			continue
		} else {
			return false
		}

	}
	return true
}

func (b *Block) AddVote(v Vote) {
	b.Votes = append(b.Votes, v)
}
