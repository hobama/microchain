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

// TODO: Generate random Transaction.

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
