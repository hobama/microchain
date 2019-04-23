package core

type Blockchain struct {
	CurrentBlock     Block            // Current block
	Chain            BlockSlice       // Stored block chain
	TransactionsPool chan Transaction // Transaction queue
	BlocksPool       chan Block       // Block queue
}

// Add block to blockchain.
func (bc *Blockchain) AddBlock(b Block) {
	bc.Chain = append(bc.Chain, b)
}
