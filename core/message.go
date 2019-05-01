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

// PingMsg ... Ping message.
type PingMsg struct {
	PublicKey []byte `json:public_key`
}

// MarshalJson ... Serialize PingMsg into Json.
func (p PingMsg) MarshalJson() ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalJson ... Read PingMsg from Json.
func (p *PingMsg) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &p)
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
