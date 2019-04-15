package core

import (
	"bytes"
	"fmt"
)

var (
	GeneratorIDBufferSize     = PublicKeyBufferSize
	PrevBlockIDBufferSize     = 32
	MerkelRootBufferSize      = 32
	TransactionsNumBufferSize = 8
	BlockHeaderBufferSize     = GeneratorIDBufferSize +
		PrevBlockIDBufferSize +
		MerkelRootBufferSize +
		TimestampBufferSize +
		TransactionsNumBufferSize
)

type BlockHeader struct {
	GeneratorID        []byte // Block generator ID                   : 256-bits : 64-bytes
	PrevBlockID        []byte // ID of previoud block                 : 256-bits : 32-bytes
	MerkelRoot         []byte // We use merkel tree to verify blocks  : 256-bits : 32-bytes
	Timestamp          uint32 // Timestamp of block generation        : 32-bits  : 4-bytes
	TransactionsLength uint64 // Length of transactions in this block : 64-bits  : 8-bytes
}

type Block struct {
	Header       BlockHeader // Block header
	Signature    []byte      // Signature by generator
	Transactions TransactionsList
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

	signature := FitBytesIntoSpecificWidth(b.Signature, SignatureBufferSize)

	transactions, err := b.Transactions.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return append(append(bh, signature...), transactions...), nil
}

// Read block from bytes.
func (b *Block) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	bh := new(BlockHeader)

	err := bh.UnmarshalBinary(buf.Next(BlockHeaderBufferSize))
	if err != nil {
		return err
	}

	b.Header = *bh
	b.Signature = StripBytes(buf.Next(SignatureBufferSize), 0)

	err = b.Transactions.UnmarshalBinary(buf.Next(buf.Len()))
	if err != nil {
		return err
	}

	if len(b.Transactions) != int(bh.TransactionsLength) {
		return fmt.Errorf("Cannot Unmarshal transactions in this block")
	}

	return nil
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

	buf.Write(FitBytesIntoSpecificWidth(bh.GeneratorID, GeneratorIDBufferSize))
	buf.Write(FitBytesIntoSpecificWidth(bh.PrevBlockID, PrevBlockIDBufferSize))
	buf.Write(FitBytesIntoSpecificWidth(bh.MerkelRoot, MerkelRootBufferSize))
	buf.Write(FitBytesIntoSpecificWidth(UInt32ToBytes(bh.Timestamp), TimestampBufferSize))
	buf.Write(FitBytesIntoSpecificWidth(UInt64ToBytes(bh.TransactionsLength), TransactionsNumBufferSize))

	return buf.Bytes(), nil
}

// Read block header from bytes.
func (bh *BlockHeader) UnmarshalBinary(data []byte) error {
	if len(data) != BlockHeaderBufferSize {
		return fmt.Errorf("Invalid transaction header")
	}

	buf := bytes.NewBuffer(data)

	bh.GeneratorID = StripBytes(buf.Next(GeneratorIDBufferSize), 0)
	bh.PrevBlockID = StripBytes(buf.Next(PrevBlockIDBufferSize), 0)
	bh.MerkelRoot = StripBytes(buf.Next(MerkelRootBufferSize), 0)

	timestamp, err := BytesToUInt32(buf.Next(TimestampBufferSize))
	if err != nil {
		return err
	}

	bh.Timestamp = timestamp

	transactionsLength, err := BytesToUInt64(buf.Next(TransactionsNumBufferSize))
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
