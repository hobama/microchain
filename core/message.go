package core

import (
	"bytes"
	"encoding/json"
)

const (
	MessageOptionsBufferSize int = 4

	// Below are message types definitions
	// TODO:
)

type Message struct {
	Type     byte   `json:"type"` // Message type
	Token    []byte `json:token`  // UUID of message
	SourceID []byte `json:"from"` // Source id
	Data     []byte `json:"data"` // Raw data
}

// Generate new message.
func NewMessage() (Message, []byte) {
	token := GenRandomBytes(32)
	return Message{Token: token}, token
}

// Test if two messages are equal.
func (m Message) EqualWith(temp Message) bool {
	if m.Type != temp.Type {
		return false
	}

	if !bytes.Equal(m.Token, temp.Token) {
		return false
	}

	if !bytes.Equal(m.SourceID, temp.SourceID) {
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
