package core

import (
	"bytes"
	"encoding/json"
)

// BlockHeader ...
type BlockHeader struct {
	GeneratorID []byte `json:"generator_id"`  // Block generator ID (Public key of generator)
	PrevBlockID []byte `json:"prev_block_id"` // ID of previoud block
	Timestamp   int    `json:"timestamp"`     // Timestamp of block generation
}

// Block ...
type Block struct {
	Header       BlockHeader      `json:"header"`       // Block header
	Signature    []byte           `json:"signature"`    // Signature by generator
	Transactions TransactionSlice `json:"transactions"` // Transactions
}

// BlockSlice ...
type BlockSlice []Block

// EqualWith ... Test if two block headers are equal.
func (bh BlockHeader) EqualWith(temp BlockHeader) bool {
	if !bytes.Equal(StripBytes(bh.GeneratorID, 0), StripBytes(temp.GeneratorID, 0)) {
		return false
	}

	if !bytes.Equal(StripBytes(bh.PrevBlockID, 0), StripBytes(temp.PrevBlockID, 0)) {
		return false
	}

	if bh.Timestamp != temp.Timestamp {
		return false
	}

	return true
}

// MarshalJson ... Serialize block header into Json.
func (bh BlockHeader) MarshalJson() ([]byte, error) {
	return json.Marshal(bh)
}

// UnmarshalJson ... Read block header from Json.
func (bh *BlockHeader) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &bh)
}

// AppendNewTransaction ... Append new transaction to this block.
func (b *Block) AppendNewTransaction(t Transaction) {
	ts := b.Transactions.Append(t)
	b.Transactions = ts
}

// InsertNewTransaction ... Insert new transaction to this block.
func (b *Block) InsertNewTransaction(t Transaction) {
	ts := b.Transactions.Insert(t)
	b.Transactions = ts
}

// EqualWith ... Test if blocks are equal.
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

// MarshalJson ... Serialize block into bytes.
func (b Block) MarshalJson() ([]byte, error) {
	return json.Marshal(b)
}

// UnmarshalJson ... Read block from Json.
func (b *Block) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &b)
}

// Contains ... Test if one block is existed in blockslice.
func (bs BlockSlice) Contains(b Block) (bool, int) {
	for i, bb := range bs {
		if b.EqualWith(bb) {
			return true, i
		}
	}

	return false, 0
}

// EqualWith ... Test if two block slices are equal.
func (bs BlockSlice) EqualWith(temp BlockSlice) bool {
	for i, b := range bs {
		if !b.EqualWith(temp[i]) {
			return false
		}
	}

	return true
}

// MarshalJson ... Serialize block slice into Json.
func (bs BlockSlice) MarshalJson() ([]byte, error) {
	return json.Marshal(bs)
}

// UnmarshalJson ... Read block slice from Json.
func (bs *BlockSlice) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &bs)
}
