package core

type BlockHeader struct {
	Origin     []byte
	PrevBlock  []byte
	MerkelRoot []byte
	Timestamp  uint32
	Nonce      uint32
}

type Block struct {
	*BlockHeader
	Signature []byte
}
