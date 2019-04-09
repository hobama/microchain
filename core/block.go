package core

import (
	"bytes"
	"fmt"
)

const (
	BlockHeaderLength = 140
	GeneratorIDLength = PublicKeyLength
	PrevBlockIDLength = 32
	MerkelRootLength  = 32
	TransactionsNumLength = 8
)

type BlockHeader struct {
	GeneratorID        []byte   // Block generator ID
	PrevBlockID        []byte   // ID of previoud block
	MerkelRoot         []byte   // We use merkel tree to verify blocks
	Timestamp          uint32   // Timestamp of block generation
	TransactionsLength uint64   // Length of transactions in this block
}

type Block struct {
	BlockHeader              // Block header
	Signature         []byte
	*TransactionsList
}

// Generate new block.
func NewBlock(prevBlockID []byte) Block {
	header := BlockHeader{PrevBlockID: prevBlockID}
	return Block{header, nil, new(TransactionsList)}
}

// Append new transaction to this block.
func (b *Block) AppendNewTransaction(t Transaction) {
	ts := b.TransactionsList.Append(t)
	b.TransactionsList = &ts
}

// Insert new transaction to this block.
func (b *Block) InsertNewTransaction(t Transaction) {
	ts := b.TransactionsList.Insert(t)
	b.TransactionsList = &ts
}

// Sign on block.
func (b *Block) Sign(keypair *KeyPair) []byte {
	s, _ := keypair.Sign(b.Hash())
	return s
}

// Calculate sha256 sum of block header.
func (b *Block) Hash() []byte {
	h, _ := b.BlockHeader.MarshalBinary()
	return SHA256(h)
}

// Serialize block header into bytes.
func (b *BlockHeader) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write(FitBytesIntoSpecificWidth(b.GeneratorID, GeneratorIDLength))
	buf.Write(FitBytesIntoSpecificWidth(b.PrevBlockID, PrevBlockIDLength))
	buf.Write(FitBytesIntoSpecificWidth(b.MerkelRoot, MerkelRootLength))
	buf.Write(UInt32ToBytes(b.Timestamp))
	buf.Write(UInt64ToBytes(b.TransactionsLength))
	return buf.Bytes(), nil
}

// Read block header from bytes.
func (b *BlockHeader) UnMarshalBinary(data []byte) error {
	if len(data) != BlockHeaderLength {
		return fmt.Errorf("Invalid transaction header")
	}

	buf := bytes.NewBuffer(data)

	b.GeneratorID = StripBytes(buf.Next(GeneratorIDLength), 0)
	b.PrevBlockID = StripBytes(buf.Next(PrevBlockIDLength), 0)
	b.MerkelRoot = StripBytes(buf.Next(MerkelRootLength), 0)

	timestamp, err := BytesToUInt32(buf.Next(TimestampLength))
	if err != nil {
		return err
	}

	b.Timestamp = timestamp

	transactionsLength, err := BytesToUInt64(buf.Next(TransactionsNumLength))
	if err != nil {
		return err
	}

	b.TransactionsLength = transactionsLength

	return nil
}
