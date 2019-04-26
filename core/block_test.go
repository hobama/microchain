package core

import (
	"errors"
	"fmt"
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
	trs := GenRandomTransactionSlice(n)

	return Block{bh, GenRandomBytes(10), trs}
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
	b := GenRandomBlock(10)

	bjson, err := b.MarshalJson()
	if err != nil {
		panic(errors.New("(Block) MarshalJson() testing failed."))
	}

	fmt.Println(string(bjson))
}
