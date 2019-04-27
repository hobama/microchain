package core

import (
	"errors"
	"math/rand"
	"testing"
)

// Generate random block header.
func GenRandomBlockHeader() BlockHeader {
	return BlockHeader{
		GenRandomBytes(64),
		GenRandomBytes(32),
		rand.Intn(10000)}
}

// Generate random block.
func GenRandomBlock(n int) Block {
	bh := GenRandomBlockHeader()
	trs := GenRandomTransactionSlice(rand.Intn(n))

	return Block{bh, GenRandomBytes(10), trs}
}

// Generate random blocks.
func GenRandomBlockSlice(b int, t int) BlockSlice {
	var bs BlockSlice

	for i := 0; i < b; i++ {
		bs = append(bs, GenRandomBlock(t))
	}

	return bs
}

// Test Block marshal function.
func TestBlockHeaderMarshalJson(t *testing.T) {
	bh1 := GenRandomBlockHeader()

	bh1json, err := bh1.MarshalJson()
	if err != nil {
		panic(errors.New("(BlockHeader) MarshalJson() testing failed."))
	}

	var bh2 BlockHeader

	err = bh2.UnmarshalJson(bh1json)
	if err != nil {
		panic(errors.New("(*BlockHeader) UnmarshalJson() testing failed."))
	}

	if !bh1.EqualWith(bh2) {
		panic(errors.New("(BlockHeader) MarshalJson()/UnmarshalJson() testing failed."))
	}
}

// Test Block marshal function.
func TestBlockMarshalJson(t *testing.T) {
	b1 := GenRandomBlock(10)

	b1json, err := b1.MarshalJson()
	if err != nil {
		panic(errors.New("(Block) MarshalJson() testing failed."))
	}

	var b2 Block

	err = b2.UnmarshalJson(b1json)
	if err != nil {
		panic(errors.New("(*Block) UnmarshalJson() testing failed."))
	}

	if !b1.EqualWith(b2) {
		panic(errors.New("(Block) MarshalJson()/UnmarshalJson() testing failed."))
	}
}

// Test BlockSlice marshal function.
func TestBlockSliceMarshalJson(t *testing.T) {
	bs1 := GenRandomBlockSlice(2, 3)

	bs1json, err := bs1.MarshalJson()
	if err != nil {
		panic(errors.New("(BlockSlice) MarshalJson() testing failed."))
	}

	var bs2 BlockSlice

	err = bs2.UnmarshalJson(bs1json)
	if err != nil {
		panic(errors.New("(*BlockSlice) UnmarshalJson() testing failed."))
	}

	if !bs1.EqualWith(bs2) {
		panic(errors.New("(BlockSlice) MarshalJson()/UnmarshalJson() testing failed."))
	}
}
