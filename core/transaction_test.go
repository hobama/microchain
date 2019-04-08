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

func TestTransactionHeaderMarshalBinary(t *testing.T) {
	from, err := NewECDSAKeyPair()
	if err != nil {
		panic(err)
	}

	to, err := NewECDSAKeyPair()
	if err != nil {
		panic(err)
	}

	transaction := NewTransaction(from.Public, to.Public, []byte{1, 2, 3, 4})
	_, err = transaction.MarshalBinary()
	if err != nil {
		panic(err)
	}
}
