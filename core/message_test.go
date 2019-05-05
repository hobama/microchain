package core

import (
	"bytes"
	"fmt"
	"testing"
)

// Generate random message
func GenRandomMessage() Message {
	return Message{
		Type: GenRandomBytes(1)[0],
		Data: GenRandomBytes(10),
	}
}

// Test PingData marshal function.
func TestPingDataMarshalJson(t *testing.T) {
	p1 := PingData{PublicKey: GenRandomBytes(32), Address: "127.0.0.1:8000"}

	p1json, err := p1.MarshalJson()
	if err != nil {
		panic(fmt.Errorf("(PingData) MarshalJson() testing failed"))
	}

	var p2 PingData

	err = p2.UnmarshalJson(p1json)
	if err != nil {
		panic(fmt.Errorf("(*PingData) UnmarshalJson() testing failed"))
	}

	if !bytes.Equal(p1.PublicKey, p2.PublicKey) {
		panic(fmt.Errorf("(PingData) MarshalJson()/UnmarshalJson() testing failed"))
	}
}

// Test SyncNodesData marshal function.
func TestSyncNodesDataMarshalJson(t *testing.T) {
	sn1 := SyncNodesData{Nodes: GenRandomRemoteNodes(5)}

	sn1json, err := sn1.MarshalJson()
	if err != nil {
		panic(fmt.Errorf("(SyncNodesData) MarshalJson() testing failed"))
	}

	var sn2 SyncNodesData

	err = sn2.UnmarshalJson(sn1json)
	if err != nil {
		panic(fmt.Errorf("(*SyncNodesData) UnmarshalJson() testing failed"))
	}

	if !sn1.EqualWith(sn2) {
		panic(fmt.Errorf("(SyncNodesData) MarshalJson()/UnmarshalJson() testing failed"))
	}
}

// Test SyncTransactionsData marshal function.
func TestSyncTransactionsDataMarshalJson(t *testing.T) {
	trs := GenRandomTransactionSlice(5)

	st1 := SyncTransactionsData{Transactions: trs}

	st1json, err := st1.MarshalJson()
	if err != nil {
		panic(fmt.Errorf("(SyncTransactionsData) MarshalJson() testing failed"))
	}

	var st2 SyncTransactionsData

	err = st2.UnmarshalJson(st1json)
	if err != nil {
		panic(fmt.Errorf("(*SyncTransactionsData) UnmarshalJson() testing failed"))
	}

	if !st1.EqualWith(st2) {
		panic(fmt.Errorf("(SyncTransactionsData) MarshalJson()/UnmarshalJson() testing failed"))
	}
}

// Test SendTransactionData marshal function.
func TestSendTransactionDataMarshalJson(t *testing.T) {
	tr1 := GenRandomTransaction()

	st1 := SendTransactionData{tr1}

	st1json, err := st1.MarshalJson()
	if err != nil {
		panic(fmt.Errorf("(SendTransactionData) MarshalJson() testing failed"))
	}

	var st2 SendTransactionData

	err = st2.UnmarshalJson(st1json)
	if err != nil {
		panic(fmt.Errorf("(*SendTransactionData) UnmarshalJson() testing failed"))
	}

	if !st1.EqualWith(st2) {
		panic(fmt.Errorf("(SendTransactionData) MarshalJson()/UnmarshalJson() testing failed"))
	}
}

// Test Message marshal function.
func TestMessageMarshalJson(t *testing.T) {
	m1 := GenRandomMessage()

	m1json, err := m1.MarshalJson()
	if err != nil {
		panic(fmt.Errorf("(Message) MarshalJson() testing failed"))
	}

	var m2 Message

	err = m2.UnmarshalJson(m1json)
	if err != nil {
		panic(fmt.Errorf("(*Message) UnmarshalJson() testing failed"))
	}

	if !m1.EqualWith(m2) {
		panic(fmt.Errorf("(Message) MarshalJson()/UnmarshalJson() testing failed"))
	}
}
