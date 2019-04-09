package core

import (
	"testing"
)

func TestGenerateNewBlocks(t *testing.T) {
	for i := 0; i < 5; i ++ {
		_ = NewBlock([]byte{0x01, 0x02, 0x03, 0x04})
	}
}
