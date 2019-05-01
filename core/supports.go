package core

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
)

// SHA256 ... Calculate SHA-256 sum.
func SHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// BytesToBigInt ... Convert bytes to big.Int.
func BytesToBigInt(b []byte) *big.Int {
	big := new(big.Int)
	big.SetBytes(b)
	return big
}

// BigIntToBytes ... Convert big.Int to bytes with specific width.
func BigIntToBytes(big *big.Int, width int) ([]byte, error) {
	b := big.Bytes()
	if len(b) > width {
		return nil, errors.New("Length of big.Int is larger than given width")
	}

	zeros := make([]byte, width-len(b))
	b = append(zeros, b...)
	return b, nil
}

// UInt32ToBytes ... Convert uint32 to bytes.
func UInt32ToBytes(i uint32) []byte {
	buf := make([]byte, 4) // uint32 => 4 bytes
	binary.LittleEndian.PutUint32(buf, i)
	return buf
}

// BytesToUInt32 ... Convert bytes to uint32.
func BytesToUInt32(bs []byte) (uint32, error) {
	byteslen := len(bs)
	if byteslen > 4 {
		return 0, fmt.Errorf("%d bytes cannot fit into uint32", byteslen)
	}

	u := binary.LittleEndian.Uint32(bs)
	return u, nil
}

// UInt64ToBytes ... Convert uint64 to bytes.
func UInt64ToBytes(i uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i)
	return buf
}

// BytesToUInt64 ... Convert bytes to uint64.
func BytesToUInt64(bs []byte) (uint64, error) {
	byteslen := len(bs)
	if byteslen > 8 {
		return 0, fmt.Errorf("%d bytes cannot fit into uint64", byteslen)
	}

	u := binary.LittleEndian.Uint64(bs)
	return u, nil
}

// JoinBytes ... Concat bytes.
func JoinBytes(bs ...[]byte) []byte {
	var data []byte
	for _, b := range bs {
		data = append(data, b...)
	}
	return data
}

// FitBytesIntoSpecificWidth ... Fit bytes into specific width.
func FitBytesIntoSpecificWidth(data []byte, i int) []byte {
	if len(data) < i {
		zeros := make([]byte, i-len(data))
		return append(zeros, data...)
	}
	return data[:i]
}

// StripBytes ... Strip bytes, usually used to delete leading zeros.
func StripBytes(data []byte, b byte) []byte {

	for i, d := range data {
		if d != b {
			return data[i:]
		}
	}

	return nil
}

// GenRandomBytes ... Generate random bytes.
func GenRandomBytes(l int) []byte {
	p := make([]byte, l)
	_, _ = rand.Read(p)
	return p
}

// Distance ... Distance of two uint64.
func Distance(a uint64, b uint64) uint64 {
	if a >= b {
		return a - b
	}

	return b - a
}
