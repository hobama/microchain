package core

import (
	//"bytes"
	"fmt"
	"testing"
)

func TestTXOutputEQ(t *testing.T) {
	tx1 := TXOutput{Accepted: 1, Rejected: 1, NextPublicKeyHash: SHA256([]byte{0x00, 0x01, 0x02, 0x03})}
	tx2 := TXOutput{Accepted: 1, Rejected: 1, NextPublicKeyHash: SHA256([]byte{0x00, 0x01, 0x02, 0x03})}

	if !tx1.EqualWith(tx2) {
		panic(fmt.Errorf("(TXOutput) EqualWith() testing failed."))
	}
}

func TestTXOutputMarshal(t *testing.T) {
	for i := 0; i < 5; i ++ {
		tx1 := TXOutput{Accepted: 1, Rejected: 1, NextPublicKeyHash: SHA256(GenRandomBytes(32))}
		tx1Bytes, err := tx1.MarshalBinary()
		if err != nil {
			panic(fmt.Errorf("(TXOutput) MarshalBinary() testing failed."))
		}

		tx2 := new(TXOutput)

		err = tx2.UnmarshalBinary(tx1Bytes)
		if err != nil {
			panic(fmt.Errorf("(*TXOutput) UnmarshalBinary() testing failed."))
		}

		if !tx1.EqualWith(*tx2) {
			panic(fmt.Errorf("(*TXOutput) UnmarshalBinary() testing failed."))
		}
	}
}

func TestTXOutputsMarshal(t *testing.T) {
	var txos TXOutputs

	for i := 0; i < 5; i ++ {
		txos.Append(TXOutput{uint64(i), uint64(i), SHA256(GenRandomBytes(32))})
	}

	b, err := txos.MarshalBinary()
	if err != nil {
		panic(err)
	}

	txoss := new(TXOutputs)

	err = txoss.UnmarshalBinary(b)
	if err != nil {
		panic(err)
	}

	if !txos.EqualWith(*txoss) {
		panic(fmt.Errorf("(*TXOutputs) UnmarshalBinary() testing failed."))
	}
}
