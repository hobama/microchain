package core

import (
	"errors"
	"math/rand"
	"testing"
)

func GenRandomBlockHeader() BlockHeader {
	return BlockHeader{
		GenRandomBytes(32),
		GenRandomBytes(32),
		GenRandomBytes(32),
		rand.Uint32(),
		uint64(rand.Intn(10))}
}

func GenRandomBlock() Block {
	bh := GenRandomBlockHeader()

	var trs TransactionsList
	for i := 0; i < int(bh.TransactionsLength); i++ {
		trs = trs.Insert(GenRandomTransaction())
	}

	b := Block{bh, GenRandomBytes(32), trs}

	return b
}

func TestBlockHeaderEQ(t *testing.T) {
	t1 := GenRandomBlockHeader()
	t2 := t1

	if !t1.EqualWith(t2) {
		panic(errors.New("(BlockHeader) EqualWith() testing failed."))
	}
}

func TestBlockHeaderMarshal(t *testing.T) {
	t1 := GenRandomBlockHeader()

	bhBytes, err := t1.MarshalBinary()
	if err != nil {
		panic(err)
	}

	t2 := new(BlockHeader)

	err = t2.UnmarshalBinary(bhBytes)
	if err != nil {
		panic(err)
	}

	if !t2.EqualWith(t1) {
		panic(errors.New("(BlockHeader) MarshalBinary() testing failed."))
	}
}

func TestBlockEQ(t *testing.T) {
	b1 := GenRandomBlock()
	b2 := b1

	if !b1.EqualWith(b2) {
		panic(errors.New("(Block) EqualWith() testing failed."))
	}
}

func TestBlockMarshal(t *testing.T) {
	b1 := GenRandomBlock()

	bBytes, err := b1.MarshalBinary()
	if err != nil {
		panic(err)
	}

	b2 := new(Block)

	err = b2.UnmarshalBinary(bBytes)
	if err != nil {
		panic(err)
	}

	if !b2.EqualWith(b1) {
		panic(errors.New("(Block) MarshalBinary() testing failed."))
	}
}
