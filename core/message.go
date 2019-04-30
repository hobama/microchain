package core

import (
	"bytes"
	"encoding/json"
)

const (
	MessageOptionsBufferSize int = 4

	// Below are message types definitions
	Ping byte = 0x01
	Join byte = 0x02
)

// Ping message.
type PingMsg struct {
	PublicKey []byte `json:public_key`
}

// Serialize PingMsg into Json.
func (p PingMsg) MarshalJson() ([]byte, error) {
	return json.Marshal(p)
}

// Read PingMsg from Json.
func (p *PingMsg) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &p)
}

// Message carrier.
type Message struct {
	Type byte   `json:"type"` // Message type
	Data []byte `json:"data"` // Raw data
}

// Test if two messages are equal.
func (m Message) EqualWith(temp Message) bool {
	if m.Type != temp.Type {
		return false
	}

	if !bytes.Equal(m.Data, temp.Data) {
		return false
	}

	return true
}

// Serialize message into Json.
func (m Message) MarshalJson() ([]byte, error) {
	return json.Marshal(m)
}

// Read message from Json.
func (m *Message) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &m)
}
