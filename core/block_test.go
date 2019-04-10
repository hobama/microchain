package core

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestGenerateNewBlocks(t *testing.T) {
	for i := 0; i < 5; i ++ {
		_ = NewBlock([]byte{0x01, 0x02, 0x03, 0x04})
	}
}

func genBlockHeaderHelper(d1, d2, d3 []byte, d4 uint32, d5 uint64) BlockHeader {
	return BlockHeader{d1, d2, d3, d4, d5}
}

func TestBlockHeaderEqual(t *testing.T) {
	// Should be equal
	bh1 := genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x03}, 5, 6)
	bh2 := genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x03}, 5, 6)

	if !bh1.EqualWith(bh2) {
		panic(fmt.Errorf("Invalid (BlockHeader) EqualWith()"))
	}

	// Should not be equal
	bh1 = genBlockHeaderHelper([]byte{0x00}, []byte{0x02}, []byte{0x03}, 5, 6)
	bh2 = genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x03}, 5, 6)

	if bh1.EqualWith(bh2) {
		panic(fmt.Errorf("Invalid (BlockHeader) EqualWith()"))
	}

	bh1 = genBlockHeaderHelper([]byte{0x01}, []byte{0x00}, []byte{0x03}, 5, 6)
	bh2 = genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x03}, 5, 6)

	if bh1.EqualWith(bh2) {
		panic(fmt.Errorf("Invalid (BlockHeader) EqualWith()"))
	}

	bh1 = genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x00}, 5, 6)
	bh2 = genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x03}, 5, 6)

	if bh1.EqualWith(bh2) {
		panic(fmt.Errorf("Invalid (BlockHeader) EqualWith()"))
	}

	bh1 = genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x03}, 0, 6)
	bh2 = genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x03}, 5, 6)

	if bh1.EqualWith(bh2) {
		panic(fmt.Errorf("Invalid (BlockHeader) EqualWith()"))
	}

	bh1 = genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x03}, 5, 0)
	bh2 = genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x03}, 5, 6)

	if bh1.EqualWith(bh2) {
		panic(fmt.Errorf("Invalid (BlockHeader) EqualWith()"))
	}
}

func TestBlockHeaderMarshalBinary(t *testing.T) {
	bh1 := genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x03}, 5, 6)
	bh2 := genBlockHeaderHelper([]byte{0x01}, []byte{0x02}, []byte{0x03}, 5, 6)

	bs, err := bh1.MarshalBinary()
	if err != nil {
		panic(err)
	}

	bh3 := new(BlockHeader)
	bh3.UnMarshalBinary(bs)

	if !bh2.EqualWith(*bh3) {
		panic(fmt.Errorf("Invalid (BlockHeader) MarshalBinary/(*BlockHeader) UnMarshalBinary"))
	}
}

func GenRandomBytes(l int) []byte {
	p := make([]byte, l)
	_, _ = rand.Read(p)
	return p
}

func genRandomBlockWithTransactionsNoGreaterThan(n uint64) *Block {
	if n == 0 {
		return nil
	}

	randID, randPrevID, randMerkelID := GenRandomBytes(32), GenRandomBytes(32), GenRandomBytes(32)

	bh := genBlockHeaderHelper(randID, randPrevID, randMerkelID, rand.Uint32(), n)

	b := NewBlock([]byte{0x00, 0x01})

	b.Header = bh

	for i := 0; i < int(n); i ++ {
		t := GenRandomTransaction()
		b.InsertNewTransaction(*t)
	}

	return &b
}

func TestBlockMarshalBinary(t *testing.T) {
	b := genRandomBlockWithTransactionsNoGreaterThan(2)

	bb, err := b.MarshalBinary()
	if err != nil {
		panic(err)
	}

	bbb := new(Block)

	err = bbb.UnMarshalBinary(bb)
	if err != nil {
		panic(err)
	}

	if !b.EqualWith(*bbb) {
		panic(fmt.Errorf("Invalid (Block) MarshalBinary()/(*Block) UnMarshalBinary()"))
	}
}
