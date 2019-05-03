package core

import (
	"bytes"
	"encoding/json"
)

const (
	// MessageOptionsBufferSize ...
	MessageOptionsBufferSize int = 4

	// Ping ...
	Ping byte = 0x01 // Ping node, test if node is online
	// Join ...
	Join byte = 0x02 // Join network
	// SyncNodes ...
	SyncNodes byte = 0x03 // Sync routing table
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
