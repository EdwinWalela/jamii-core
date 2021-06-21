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
	votes      []Vote // list of votes representing candidates/voters (if genesis block) and cast votes if not genesis
	hash       string // hash of the block (votes hashes + block hash)
	prevHash   string // hash of previous block in the chain
	timestamp  uint64 // Unix timestamp of tx in seconds
	nonce      uint64 // Proof of work tries
	difficulty uint64 // Proof of work difficulty
}

func (b *Block) GetHash() string {
	return b.hash
}

func (b *Block) GetPreviousHash() string {
	return b.prevHash
}

func (b *Block) SetPreviousHash(hash string) {
	b.prevHash = hash
}

func (b *Block) SetDifficulty(diff uint64) {
	b.difficulty = diff
}

func (b *Block) GetNonce() uint64 {
	return b.nonce
}

func (b *Block) SetHash(hash string) {
	b.hash = hash
}

func (b *Block) Hash() {
	hash := ""
	// Retrieve hashes of all votes
	for _, v := range b.votes {
		hash += v.GetHash()
	}

	hash += strconv.FormatInt(int64(b.timestamp), 10)
	hash += strconv.FormatUint(b.difficulty, 10)
	b.hash = jcrypto.SHA512(hash)
	for !b.HashValid() {
		b.nonce++
		b.hash = jcrypto.SHA512(fmt.Sprintf("%s%d", b.hash, b.nonce))

	}
}

func (b *Block) HashValid() bool {
	for i, v := range b.hash {
		if i == int(b.difficulty) {
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
	b.votes = append(b.votes, v)
}
