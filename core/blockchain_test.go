package core

import (
	"errors"
	"testing"
)

// Test append block to blockchain.
func TestAddBlockToChain(t *testing.T) {
	bc := Blockchain{Block{}, []Block{}, nil, nil}

	b := GenRandomBlock()

	bc.AddBlock(b)

	for _, bb := range bc.Chain {
		if bb.EqualWith(b) {
			return
		}
	}

	panic(errors.New("(*Blockchain) AddBlock(b Block) testing failed."))
}
