package core

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

func checkTransaction(transaction *Transaction) bool {
	timeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(timeBuf, transaction.Header.Timestamp)
	rawid := JoinBytes(timeBuf, transaction.Header.RequesterPublicKey, transaction.Header.RequesteePublicKey)
	return reflect.DeepEqual(transaction.Header.TransactionID, SHA256(rawid))
}

func TestGenerateNewTransaction(t *testing.T) {
	for i := 0; i < 5; i++ {
		from, err := NewECDSAKeyPair()
		if err != nil {
			panic(err)
		}
		to, err := NewECDSAKeyPair()
		if err != nil {
			panic(err)
		}

		transaction := NewTransaction(from.Public, to.Public, []byte{1, 2, 3, 4})
		if err != nil {
			panic(err)
		}

		b := checkTransaction(transaction)
		if b != true {
			panic(fmt.Errorf("Invalid transaction generated"))
		}
	}
}

func GenRandomTransaction() *Transaction {
	from, err := NewECDSAKeyPair()
	if err != nil {
		panic(err)
	}

	to, err := NewECDSAKeyPair()
	if err != nil {
		panic(err)
	}

	meta := GenRandomBytes(20)

	transaction := NewTransaction(from.Public, to.Public, meta)

	return transaction
}

func TestTransactionMarshalBinary(t *testing.T) {
	for i := 0; i < 50; i++ {
		from, err := NewECDSAKeyPair()
		if err != nil {
			panic(err)
		}

		to, err := NewECDSAKeyPair()
		if err != nil {
			panic(err)
		}

		transaction := NewTransaction(from.Public,
		to.Public,
		[]byte{byte(i), byte(i + 1), byte(i + 2), byte(i + 3)})
		transactionBytes, err := transaction.MarshalBinary()
		if err != nil {
			panic(err)
		}

		var newTransaction Transaction
		err = newTransaction.UnMarshalBinary(transactionBytes)
		if err != nil {
			panic(err)
		}

		if !transaction.EqualWith(newTransaction) {
			panic(fmt.Errorf("Cannot marshal/unmarshal transaction"))
		}
	}
}

func TestAppendTransaction(t *testing.T) {
	from, err := NewECDSAKeyPair()
	if err != nil {
		panic(err)
	}
	to, err := NewECDSAKeyPair()
	if err != nil {
		panic(err)
	}

	transaction := NewTransaction(from.Public, to.Public, []byte{1, 2, 3, 4})
	if err != nil {
		panic(err)
	}

	transactions := make(TransactionsList, 10)

	transactions = transactions.Append(*transaction)
	if err != nil {
		panic(err)
	}

	isContained, index := transactions.Contains(*transaction)
	if !isContained || index != 10 || len(transactions) != 11 {
		panic(fmt.Errorf("Append tansaction to list failed"))
	}
}

func TestTransactionsListMarshalBinary(t *testing.T) {
	transactions := new(TransactionsList)
	transactionsbkup := new(TransactionsList)

	for i := 0; i < 50; i++ {
		from, err := NewECDSAKeyPair()
		if err != nil {
			panic(err)
		}

		to, err := NewECDSAKeyPair()
		if err != nil {
			panic(err)
		}

		transaction := NewTransaction(from.Public, to.Public,
		[]byte{byte(i), byte(i + 1), byte(i + 2), byte(i + 3)})
		transactions.Append(*transaction)
		transactionsbkup.Append(*transaction)
	}

	bs, err := transactions.MarshalBinary()
	if err != nil {
		panic(err)
	}

	newTransactions := new(TransactionsList)
	err = newTransactions.UnMarshalBinary(bs)
	if err != nil {
		panic(err)
	}

	if !newTransactions.EqualWith(*transactionsbkup) {
		panic(fmt.Errorf("Cannot marshal/unmarshal transactions"))
	}
}
