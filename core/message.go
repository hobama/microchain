package core

import (
	"bytes"
	"encoding/json"
)

const (
	// MessageOptionsBufferSize ...
	MessageOptionsBufferSize int = 4

	Ping                 byte = 0x01 // Ping node, test if node is online
	Join                 byte = 0x02 // Join network
	SyncNodes            byte = 0x03 // Sync routing table
	SendTransaction      byte = 0x04 // Send transaction to given node
	PendingTransaction   byte = 0x05 // Send pending transaction
	BroadcastTransaction byte = 0x06 // Broadcast transaction by requestee node
	SyncTransactions     byte = 0x07 // Sync transactions
)

// PingData ... Ping data.
type PingData struct {
	PublicKey []byte `json:"public_key"`
	Address   string `json:"server_addr"`
}

// MarshalJson ... Serialize PingData into Json.
func (p PingData) MarshalJson() ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalJson ... Read PingData from Json.
func (p *PingData) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &p)
}

// NewPingMessage ... Generate new ping message.
func NewPingMessage(pk []byte, addr string) Message {
	data := PingData{pk, addr}
	dataJSON, _ := data.MarshalJson()

	return Message{Type: Ping, Data: dataJSON}
}

// SyncNodesData ... Sync nodes.
type SyncNodesData struct {
	Nodes []RemoteNode `json:"nodes"`
}

// EqualWith ... Test if two SyncNodesData are equal.
func (sn SyncNodesData) EqualWith(temp SyncNodesData) bool {
	if len(sn.Nodes) != len(temp.Nodes) {
		return false
	}

	for i, n := range sn.Nodes {
		if !n.EqualWith(temp.Nodes[i]) {
			return false
		}
	}

	return true
}

// MarshalJson ... Serialize SyncNodesData into Json.
func (sn SyncNodesData) MarshalJson() ([]byte, error) {
	return json.Marshal(sn)
}

// UnmarshalJson ... Read SyncNodesData from Json.
func (sn *SyncNodesData) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &sn)
}

// NewSyncNodesMessage ... Generate new sync nodes message.
func NewSyncNodesMessage(nodes []RemoteNode) Message {
	data := SyncNodesData{nodes}
	dataJSON, _ := data.MarshalJson()

	return Message{Type: SyncNodes, Data: dataJSON}
}

// SendTransactionData ... Send transaction.
type SendTransactionData struct {
	Transaction `json:"transaction"`
}

// EqualWith ... Test if two SendTransactionData are equal.
func (st SendTransactionData) EqualWith(temp SendTransactionData) bool {
	if !st.Transaction.EqualWith(temp.Transaction) {
		return false
	}

	return true
}

// MarshalJson ... Serialize SendTransactionData into Json.
func (st SendTransactionData) MarshalJson() ([]byte, error) {
	return json.Marshal(st)
}

// UnmarshalJson ... Read SendTransactionData from Json.
func (st *SendTransactionData) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &st)
}

// NewSendTransactionMessage ... Generate new send transaction message.
func NewSendTransactionMessage(t Transaction) Message {
	data := SendTransactionData{t}

	dataJSON, _ := data.MarshalJson()

	return Message{Type: SendTransaction, Data: dataJSON}
}

// PendingTransactionData ... Pending transaction.
type PendingTransactionData struct {
	Transaction `json:"transaction"`
}

// EqualWith ... Test if two PendingTransactionData are equal.
func (pt PendingTransactionData) EqualWith(temp PendingTransactionData) bool {
	if !pt.Transaction.EqualWith(temp.Transaction) {
		return false
	}

	return true
}

// MarshalJson ... Serialize PendingTransactionData into Json.
func (pt PendingTransactionData) MarshalJson() ([]byte, error) {
	return json.Marshal(pt)
}

// UnmarshalJson ... Read PendingTransactionData from Json.
func (pt *PendingTransactionData) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &pt)
}

// NewPendingTransactionMessage ... Generate new pending transaction message.
func NewPendingTransactionMessage(t Transaction) Message {
	data := PendingTransactionData{t}

	dataJSON, _ := data.MarshalJson()

	return Message{Type: PendingTransaction, Data: dataJSON}
}

// SyncTransactionsData ... Sync transactions.
type SyncTransactionsData struct {
	Transactions TransactionSlice `json:"transactions"`
}

// EqualWith ... Test if two SyncTransactionsData are equal.
func (st SyncTransactionsData) EqualWith(temp SyncTransactionsData) bool {
	if len(st.Transactions) != len(temp.Transactions) {
		return false
	}

	for i, t := range st.Transactions {
		if !t.EqualWith(temp.Transactions[i]) {
			return false
		}
	}

	return true
}

// MarshalJson ... Serialize SyncTransactionsData into Json.
func (st SyncTransactionsData) MarshalJson() ([]byte, error) {
	return json.Marshal(st)
}

// UnmarshalJson ... Read SyncTransactionsData from Json.
func (st *SyncTransactionsData) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &st)
}

// NewSyncTransactionsMessage ... Generate new sync transactions message.
func NewSyncTransactionsMessage(transactions TransactionSlice) Message {
	data := SyncTransactionsData{transactions}

	dataJSON, _ := data.MarshalJson()

	return Message{Type: SyncTransactions, Data: dataJSON}
}

// Message ... Message carrier.
type Message struct {
	Type byte   `json:"type"` // Message type
	Data []byte `json:"data"` // Raw data
}

// EqualWith ... Test if two messages are equal.
func (m Message) EqualWith(temp Message) bool {
	if m.Type != temp.Type {
		return false
	}

	if !bytes.Equal(m.Data, temp.Data) {
		return false
	}

	return true
}

// MarshalJson ... Serialize message into Json.
func (m Message) MarshalJson() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalJson ... Read message from Json.
func (m *Message) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &m)
}
