package primitives

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

func (c *Chain) SetDifficulty(diff uint64) {
	c.Difficulty = diff
}

func (c *Chain) AddTX(tx Vote) {
	c.PendingVotes = append(c.PendingVotes, tx)
}
