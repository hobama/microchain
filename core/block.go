package core

import (
	"bytes"
)

type BlockHeader struct {
	GeneratorID []byte // Block generator ID                   : 256-bits : 64-bytes
	PrevBlockID []byte // ID of previoud block                 : 256-bits : 32-bytes
	MerkelRoot  []byte // We use merkel tree to verify blocks  : 256-bits : 32-bytes
	Timestamp   int    // Timestamp of block generation        : 32-bits  : 4-bytes
}

type Block struct {
	Header       BlockHeader // Block header
	Signature    []byte      // Signature by generator
	Transactions TransactionSlice
}

type BlockSlice []Block

// Test if one block is existed in blockslice.
func (bs BlockSlice) Contains(b Block) (bool, int) {
	for i, bb := range bs {
		if b.EqualWith(bb) {
			return true, i
		}
	}

	return false, 0
}

// Get previous block that added to the slice.
func (bs BlockSlice) PreviousBlock() *Block {
	l := len(bs)

	if l == 0 {
		return nil
	} else {
		return &bs[l-1]
	}
}

// Test if two block slices are equal.
func (bs BlockSlice) EqualWith(temp BlockSlice) bool {
	for i, b := range bs {
		if !b.EqualWith(temp[i]) {
			return false
		}
	}

	return true
}

// Append new transaction to this block.
func (b *Block) AppendNewTransaction(t Transaction) {
	ts := b.Transactions.Append(t)
	b.Transactions = ts
}

// Insert new transaction to this block.
func (b *Block) InsertNewTransaction(t Transaction) {
	ts := b.Transactions.Insert(t)
	b.Transactions = ts
}

// Test if blocks are equal.
func (b Block) EqualWith(temp Block) bool {
	if !b.Header.EqualWith(temp.Header) {
		return false
	}

	if !bytes.Equal(b.Signature, temp.Signature) {
		return false
	}

	if !b.Transactions.EqualWith(temp.Transactions) {
		return false
	}

	return true
}

// Test if two block headers are equal.
func (bh BlockHeader) EqualWith(temp BlockHeader) bool {
	if !bytes.Equal(StripBytes(bh.GeneratorID, 0), StripBytes(temp.GeneratorID, 0)) {
		return false
	}

	if !bytes.Equal(StripBytes(bh.PrevBlockID, 0), StripBytes(temp.PrevBlockID, 0)) {
		return false
	}

	if !bytes.Equal(StripBytes(bh.MerkelRoot, 0), StripBytes(temp.MerkelRoot, 0)) {
		return false
	}

	if bh.Timestamp != temp.Timestamp {
		return false
	}

	return true
}
