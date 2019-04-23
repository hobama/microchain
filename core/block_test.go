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

	var trs TransactionSlice
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

func TestBlockExist(t *testing.T) {
	var bs BlockSlice

	for i := 0; i < 5; i++ {
		b := GenRandomBlock()
		bs = append(bs, b)
	}

	for i := 0; i < 5; i++ {
		b := bs[i]

		ok, index := bs.Contains(b)

		if !ok {
			panic(errors.New("(BlockSlice) Contains() testing failed."))
		}

		if i != index {
			panic(errors.New("(BlockSlice) Contains() testing failed."))
		}
	}
}

func TestBlockSliceEQ(t *testing.T) {
	var bs1 BlockSlice

	for i := 0; i < 5; i++ {
		bs1 = append(bs1, GenRandomBlock())
	}

	bs2 := bs1

	if !bs1.EqualWith(bs2) {
		panic(errors.New("(BlockSlice) EqualWith() testing failed."))
	}
}

func TestBlockSliceMarshal(t *testing.T) {
	var bs1 BlockSlice

	for i := 0; i < 5; i++ {
		bs1 = append(bs1, GenRandomBlock())
	}

	bsBytes, err := bs1.MarshalBinary()
	if err != nil {
		panic(err)
	}

	bs2 := new(BlockSlice)

	err = bs2.UnmarshalBinary(bsBytes)
	if err != nil {
		panic(err)
	}

	if !bs2.EqualWith(bs1) {
		panic(errors.New("(BlocakSlice) MarshalBinary() testing failed."))
	}
}
