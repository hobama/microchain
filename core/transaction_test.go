package core

import (
	"errors"
	"math/rand"
	"testing"
)

func GenRandomTXOutput() TXOutput {
	return TXOutput{uint64(rand.Intn(100)), uint64(rand.Intn(100)), SHA256(GenRandomBytes(32))}
}

func GenRandomTransactionHeader() TransactionHeader {
	return TransactionHeader{
		SHA256(GenRandomBytes(32)),
		uint32(rand.Intn(10000)),
		GenRandomBytes(32),
		GenRandomBytes(32),
		GenRandomBytes(32),
		GenRandomBytes(32),
		GenRandomBytes(32),
		uint64(rand.Intn(100))}
}

func GenRandomTransaction() Transaction {
	t := Transaction{}
	t.Header = GenRandomTransactionHeader()
	t.Meta = GenRandomBytes(int(t.Header.MetaLength))

	t.Output = GenRandomTXOutput()

	return t
}

func TestTXOutputEQ(t *testing.T) {
	tx1 := TXOutput{Accepted: 1, Rejected: 1, NextPublicKeyHash: SHA256([]byte{0x00, 0x01, 0x02, 0x03})}
	tx2 := TXOutput{Accepted: 1, Rejected: 1, NextPublicKeyHash: SHA256([]byte{0x00, 0x01, 0x02, 0x03})}

	if !tx1.EqualWith(tx2) {
		panic(errors.New("(TXOutput) EqualWith() testing failed."))
	}
}

func TestTransactionHeaderEQ(t *testing.T) {
	th1 := GenRandomTransactionHeader()
	th2 := th1

	if !th1.EqualWith(th2) {
		panic(errors.New("(TransactionHeader) EqualWith() testing failed."))
	}
}

func TestTransactionEQ(t *testing.T) {
	t1 := GenRandomTransaction()
	t2 := t1

	if !t1.EqualWith(t2) {
		panic(errors.New("(Transaction) EqualWith() testing failed."))
	}
}

func TestTransactionSliceEQ(t *testing.T) {
	var trs TransactionSlice

	for i := 0; i < 10; i++ {
		trs.Append(GenRandomTransaction())
	}

	trs2 := trs

	if !trs.EqualWith(trs2) {
		panic(errors.New("(Transactions) EqualWith() testing failed."))
	}
}

func TestTXOutputMarshal(t *testing.T) {
	for i := 0; i < 5; i++ {
		tx1 := GenRandomTXOutput()
		tx1Bytes, err := tx1.MarshalBinary()
		if err != nil {
			panic(errors.New("(TXOutput) MarshalBinary() testing failed."))
		}

		tx2 := new(TXOutput)

		err = tx2.UnmarshalBinary(tx1Bytes)
		if err != nil {
			panic(errors.New("(*TXOutput) UnmarshalBinary() testing failed."))
		}

		if !tx1.EqualWith(*tx2) {
			panic(errors.New("(*TXOutput) UnmarshalBinary() testing failed."))
		}
	}
}

func TestTransactionHeaderMarshal(t *testing.T) {
	th := GenRandomTransactionHeader()

	thBytes, err := th.MarshalBinary()
	if err != nil {
		panic(err)
	}

	th2 := new(TransactionHeader)

	err = th2.UnmarshalBinary(thBytes)
	if err != nil {
		panic(err)
	}

	if !th2.EqualWith(th) {
		panic(errors.New("(*TransactionHeader) UnmarshalBinary() testing failed."))
	}
}

func TestTransactionMarshal(t *testing.T) {
	tr := GenRandomTransaction()

	trBytes, err := tr.MarshalBinary()
	if err != nil {
		panic(err)
	}

	tr2 := new(Transaction)

	err = tr2.UnmarshalBinary(trBytes)
	if err != nil {
		panic(err)
	}

	if !tr2.EqualWith(tr) {
		panic(errors.New("(*Transaction) UnmarshalBinary() testing failed."))
	}
}

func TestTransactionSliceMarshal(t *testing.T) {
	var trs TransactionSlice

	for i := 0; i < 10; i++ {
		trs.Append(GenRandomTransaction())
	}

	trsBytes, err := trs.MarshalBinary()
	if err != nil {
		panic(err)
	}

	trs2 := new(TransactionSlice)

	err = trs2.UnmarshalBinary(trsBytes)
	if err != nil {
		panic(err)
	}

	if !trs2.EqualWith(trs) {
		panic(errors.New("(*TransactionSlice) UnmarshalBinary() testing failed."))
	}
}

// Test diff two transaction slice.
func TestDiffTransactionSlice(t *testing.T) {
	var ts1, ts2 TransactionSlice

	// Insert even indexed item to ts2.
	for i := 0; i < 10; i++ {
		tr := GenRandomTransaction()
		ts1 = append(ts1, tr)
		if i%2 == 0 {
			ts2 = append(ts2, tr)
		}
	}

	// Odd indexed items are passed to ts3.
	ts3 := DiffTransactions(ts1, ts2)

	// Test if odd indexed items are stored in ts3.
	for i, tr := range ts3 {
		if !tr.EqualWith(ts1[2*i+1]) {
			panic(errors.New("DiffTransactions() testing failed."))
		}
	}
}
