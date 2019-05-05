package core

// Blockchain ...
type Blockchain struct {
	Blocks BlockSlice // Stored block chain
}

// GetTransactionByID ... Get transaction by id.
func (bc Blockchain) GetTransactionByID(id []byte) (bool, Transaction) {
	if t, tr := bc.Blocks.GetTransactionByID(id); t {
		return true, tr
	}

	return false, Transaction{}
}
