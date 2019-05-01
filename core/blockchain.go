package core

// Blockchain ...
type Blockchain struct {
	CurrentBlock     Block            // Current block
	Chain            BlockSlice       // Stored block chain
	TransactionsPool chan Transaction // Transaction queue
	BlocksPool       chan Block       // Block queue
}
