package core

import (
	"errors"
	"math/rand"
	"testing"
)

// Generate random TXOutput.
func GenRandomTXOutput() TXOutput {
	return TXOutput{
		rand.Intn(10000),    // Accepted
		rand.Intn(10000),    // Rejected
		GenRandomBytes(100)} // NextPublicKeyHash
}

// Generate random Transaction header.
func GenRandomTransactionHeader() TransactionHeader {
	return TransactionHeader{
		GenRandomBytes(32), // TransactionID
		rand.Intn(10000),   // Timestamp
		GenRandomBytes(32), // PrevTransactionID
		GenRandomBytes(32), // RequesterPublicKey
		GenRandomBytes(32), // RequesterSignature
		GenRandomBytes(32), // RequesteePublicKey
		GenRandomBytes(32)} // RequesteeSignature
}

// Generate random Transaction.
func GenRandomTransaction() Transaction {
	th := GenRandomTransactionHeader()
	txo := GenRandomTXOutput()

	return Transaction{
		th,                 // TransactionHeader
		GenRandomBytes(10), // Meta
		txo}                // TXOutput
}

// Generate random TransactionSlice.
func GenRandomTransactionSlice(n int) TransactionSlice {
	var trs TransactionSlice

	for i := 0; i < n; i++ {
		trs = trs.Append(GenRandomTransaction())
	}

	return trs
}

// Test TXOutput marshal function.
func TestTXOutputMarshalJson(t *testing.T) {
	tx1 := GenRandomTXOutput()

	tx1json, err := tx1.MarshalJson()
	if err != nil {
		panic(errors.New("(TXOutput) MarshalJson() testing failed."))
	}

	var tx2 TXOutput

	err = tx2.UnmarshalJson(tx1json)
	if err != nil {
		panic(errors.New("(*TXOutput) UnmarshalJson() testing failed."))
	}

	if !tx1.EqualWith(tx2) {
		panic(errors.New("(TXOutput) MarshalJson()/UnmarshalJson() testing failed."))
	}
}

// Test TransactionHeader marshal function.
func TestTransactionHeaderMarshalJson(t *testing.T) {
	th1 := GenRandomTransactionHeader()

	th1json, err := th1.MarshalJson()
	if err != nil {
		panic(errors.New("(TransactionHeader) MarshalJson() testing failed."))
	}

	var th2 TransactionHeader

	err = th2.UnmarshalJson(th1json)
	if err != nil {
		panic(errors.New("(*TransactionHeader) UnmarshalJson() testing failed."))
	}

	if !th1.EqualWith(th2) {
		panic(errors.New("(TransactionHeader) MarshalJson()/UnmarshalJson() testing failed."))
	}
}

// Test Transaction marshal function.
func TestTransactionMarshalJson(t *testing.T) {
	tr1 := GenRandomTransaction()

	tr1json, err := tr1.MarshalJson()
	if err != nil {
		panic(errors.New("(Transaction) MarshalJson() testing failed."))
	}

	var tr2 Transaction

	err = tr2.UnmarshalJson(tr1json)
	if err != nil {
		panic(errors.New("(*Transaction) UnmarshalJson() testing failed."))
	}

	if !tr1.EqualWith(tr2) {
		panic(errors.New("(Transaction) MarshalJson()/UnmarshalJson() testing failed."))
	}
}

// Test TransactionSlice marshal function.
func TestTransactionSliceMarshalJson(t *testing.T) {
	trs1 := GenRandomTransactionSlice(10)

	trs1json, err := trs1.MarshalJson()
	if err != nil {
		panic(errors.New("(TransactionSlice) MarshalJson() testing failed."))
	}

	var trs2 TransactionSlice

	err = trs2.UnmarshalJson(trs1json)
	if err != nil {
		panic(errors.New("(*TransactionSlice) UnmarshalJson() testing failed."))
	}

	if !trs1.EqualWith(trs2) {
		panic(errors.New("(TransactionSlice) MarshalJson()/UnmarshalJson() testing failed."))
	}
}
