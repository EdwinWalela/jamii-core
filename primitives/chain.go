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
	chain        []Block // list of verified blocks
	pendingVotes []Vote  // list of unverified votes recieved from clients via TCP socket
	height       uint64  // current height of the chain
	difficulty   uint64  // node's proof of work difficulty
}
