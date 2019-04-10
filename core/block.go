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
	Header            BlockHeader // Block header
	Signature         []byte
	Transactions      TransactionsList
}

// Generate new block.
func NewBlock(prevBlockID []byte) Block {
	header := BlockHeader{PrevBlockID: prevBlockID}
	return Block{header, nil, []Transaction{}}
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

// Sign on block.
func (b *Block) Sign(keypair *KeyPair) []byte {
	s, _ := keypair.Sign(b.Hash())
	return s
}

// Calculate sha256 sum of block header.
func (b *Block) Hash() []byte {
	h, _ := b.Header.MarshalBinary()
	return SHA256(h)
}

// Serialize block into bytes.
func (b Block) MarshalBinary() ([]byte, error) {
	bh, err := b.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	signature := FitBytesIntoSpecificWidth(b.Signature, SignatureLength)

	transactions, err := b.Transactions.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return append(append(bh, signature...), transactions...), nil
}

// Read block from bytes.
func (b *Block) UnMarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	bh := new(BlockHeader)

	err := bh.UnMarshalBinary(buf.Next(BlockHeaderLength))
	if err != nil {
		return err
	}

	b.Header = *bh
	b.Signature = StripBytes(buf.Next(SignatureLength), 0)

	err = b.Transactions.UnMarshalBinary(buf.Next(buf.Len()))
	if err != nil {
		return err
	}

	if len(b.Transactions) != int(bh.TransactionsLength) {
		return fmt.Errorf("Cannot Unmarshal transactions in this block")
	}

	return nil
}

func (b *Block) GenerateMerkelRoot() []byte {
	return []byte{0x00}
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

// Serialize block header into bytes.
func (bh BlockHeader) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write(FitBytesIntoSpecificWidth(bh.GeneratorID, GeneratorIDLength))
	buf.Write(FitBytesIntoSpecificWidth(bh.PrevBlockID, PrevBlockIDLength))
	buf.Write(FitBytesIntoSpecificWidth(bh.MerkelRoot, MerkelRootLength))
	buf.Write(UInt32ToBytes(bh.Timestamp))
	buf.Write(UInt64ToBytes(bh.TransactionsLength))

	return buf.Bytes(), nil
}

// Read block header from bytes.
func (bh *BlockHeader) UnMarshalBinary(data []byte) error {
	if len(data) != BlockHeaderLength {
		return fmt.Errorf("Invalid transaction header")
	}

	buf := bytes.NewBuffer(data)

	bh.GeneratorID = StripBytes(buf.Next(GeneratorIDLength), 0)
	bh.PrevBlockID = StripBytes(buf.Next(PrevBlockIDLength), 0)
	bh.MerkelRoot = StripBytes(buf.Next(MerkelRootLength), 0)

	timestamp, err := BytesToUInt32(buf.Next(TimestampLength))
	if err != nil {
		return err
	}

	bh.Timestamp = timestamp

	transactionsLength, err := BytesToUInt64(buf.Next(TransactionsNumLength))
	if err != nil {
		return err
	}

	bh.TransactionsLength = transactionsLength

	return nil
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

	if bh.TransactionsLength != temp.TransactionsLength {
		return false
	}

	return true
}

