package primitives

type Block struct {
	/**
	Genesis block will contain a single Vote with the candidates field indicating registered
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
