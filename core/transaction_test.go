package core

import (
	"fmt"
	"testing"
	"math/rand"
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
		uint64(rand.Intn(100)),
		uint64(rand.Intn(10))}
}

func GenRandomTransaction() Transaction {
	t := Transaction{}
	t.Header = GenRandomTransactionHeader()
	t.Meta = GenRandomBytes(int(t.Header.MetaLength))

	for i := 0; i < int(t.Header.OutputLength); i ++ {
		t.Outputs.Append(GenRandomTXOutput())
	}

	return t
}

func TestTXOutputEQ(t *testing.T) {
	tx1 := TXOutput{Accepted: 1, Rejected: 1, NextPublicKeyHash: SHA256([]byte{0x00, 0x01, 0x02, 0x03})}
	tx2 := TXOutput{Accepted: 1, Rejected: 1, NextPublicKeyHash: SHA256([]byte{0x00, 0x01, 0x02, 0x03})}

	if !tx1.EqualWith(tx2) {
		panic(fmt.Errorf("(TXOutput) EqualWith() testing failed."))
	}
}

func TestTransactionHeaderEQ(t *testing.T) {
	th1 := GenRandomTransactionHeader()
	th2 := th1

	if !th1.EqualWith(th2) {
		panic(fmt.Errorf("(TransactionHeader) EqualWith() testing failed."))
	}
}

func TestTransactionEQ(t *testing.T) {
	t1 := GenRandomTransaction()
	t2 := t1

	if !t1.EqualWith(t2) {
		panic(fmt.Errorf("(Transaction) EqualWith() testing failed."))
	}
}

func TestTransactionsListEQ(t *testing.T) {
	var trs TransactionsList

	for i := 0; i < 10; i ++ {
		trs.Append(GenRandomTransaction())
	}

	trs2 := trs

	if !trs.EqualWith(trs2) {
		panic(fmt.Errorf("(Transactions) EqualWith() testing failed."))
	}
}

func TestTXOutputMarshal(t *testing.T) {
	for i := 0; i < 5; i ++ {
		tx1 := GenRandomTXOutput()
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
		txos.Append(GenRandomTXOutput())
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
		panic(fmt.Errorf("(*TransactionHeader) UnmarshalBinary() testing failed."))
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
		panic(fmt.Errorf("(*Transaction) UnmarshalBinary() testing failed."))
	}
}

func TestTransactionsListMarshal(t *testing.T) {
	var trs TransactionsList

	for i := 0; i < 10; i ++ {
		trs.Append(GenRandomTransaction())
	}

	trsBytes, err := trs.MarshalBinary()
	if err != nil {
		panic(err)
	}

	trs2 := new(TransactionsList)

	err = trs2.UnmarshalBinary(trsBytes)
	if err != nil {
		panic(err)
	}

	if !trs2.EqualWith(trs) {
		panic(fmt.Errorf("(*TransactionsList) UnmarshalBinary() testing failed."))
	}
}
