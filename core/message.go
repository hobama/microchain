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
)

// PingData ... Ping data.
type PingData struct {
	PublicKey []byte `json:public_key`
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
func NewPingMessage(n Node) Message {
	data := PingData{n.Keypair.Public}
	dataJSON, _ := data.MarshalJson()

	return Message{Type: Ping, Data: dataJSON}
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
