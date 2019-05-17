package core

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

// Generate random TXOutput.
func GenRandomTXOutput() TXOutput {
	return TXOutput{
		rand.Intn(10000), // Accepted
		rand.Intn(10000), // Rejected
	}
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
		panic(fmt.Errorf("(TXOutput) MarshalJson() testing failed"))
	}

	var tx2 TXOutput

	err = tx2.UnmarshalJson(tx1json)
	if err != nil {
		panic(fmt.Errorf("(*TXOutput) UnmarshalJson() testing failed"))
	}

	if !tx1.EqualWith(tx2) {
		panic(fmt.Errorf("(TXOutput) MarshalJson()/UnmarshalJson() testing failed"))
	}
}

// Test TransactionHeader marshal function.
func TestTransactionHeaderMarshalJson(t *testing.T) {
	th1 := GenRandomTransactionHeader()

	th1json, err := th1.MarshalJSON()
	if err != nil {
		panic(fmt.Errorf("(TransactionHeader) MarshalJson() testing failed"))
	}

	var th2 TransactionHeader

	err = th2.UnmarshalJSON(th1json)
	if err != nil {
		panic(fmt.Errorf("(*TransactionHeader) UnmarshalJson() testing failed"))
	}

	if !th1.EqualWith(th2) {
		panic(fmt.Errorf("(TransactionHeader) MarshalJson()/UnmarshalJson() testing failed"))
	}
}

// Test Transaction marshal function.
func TestTransactionMarshalJson(t *testing.T) {
	t1 := GenRandomTransaction()

	t1json, err := t1.MarshalJSON()
	if err != nil {
		panic(fmt.Errorf("(Transaction) MarshalJson() testing failed"))
	}

	var t2 Transaction

	err = t2.UnmarshalJSON(t1json)
	if err != nil {
		panic(fmt.Errorf("(*Transaction) UnmarshalJson() testing failed"))
	}

	if !t1.EqualWith(t2) {
		panic(fmt.Errorf("(Transaction) MarshalJson()/UnmarshalJson() testing failed"))
	}
}

// Test TransactionSlice marshal function.
func TestTransactionSliceMarshalJson(t *testing.T) {
	trs1 := GenRandomTransactionSlice(10)

	trs1json, err := trs1.MarshalJson()
	if err != nil {
		panic(fmt.Errorf("(TransactionSlice) MarshalJson() testing failed"))
	}

	var trs2 TransactionSlice

	err = trs2.UnmarshalJson(trs1json)
	if err != nil {
		panic(fmt.Errorf("(*TransactionSlice) UnmarshalJson() testing failed"))
	}

	if !trs1.EqualWith(trs2) {
		panic(fmt.Errorf("(TransactionSlice) MarshalJson()/UnmarshalJson() testing failed"))
	}
}

// Test sorting transactions.
func TestSortTransactions(t *testing.T) {
	var trs TransactionSlice

	for i := 4; i >= 0; i-- {
		th := TransactionHeader{Timestamp: i}
		t := Transaction{Header: th}

		trs = append(trs, t)
	}

	sort.Sort(trs)

	for i, tr := range trs {
		if tr.Header.Timestamp != i {
			panic(fmt.Errorf("Sort() testing failed"))
		}
	}
}

// Test verify signature.
func TestVerifySignature(t *testing.T) {
	kp, _ := NewECDSAKeyPair()

	tr := GenRandomTransaction()

	tr.Header.RequesterPublicKey = kp.Public
	tr.Header.RequesteePublicKey = kp.Public

	Sig, _ := kp.Sign(tr.Hash())
	tr.Header.RequesterSignature = Sig
	tr.Header.RequesteeSignature = Sig

	if !tr.VerifyRequesterSig() {
		panic(fmt.Errorf("(Transaction) VerifyRequesterSig() testing failed"))
	}

	if !tr.VerifyRequesteeSig() {
		panic(fmt.Errorf("(Transaction) VerifyRequesteeSig() testing failed"))
	}
}
