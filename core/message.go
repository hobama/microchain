package core

import (
	"bytes"
	"errors"
)

const (
	MessageOptionsBufferSize int = 4

	// Below are message types definitions
	MSG_NULL byte = 0x00
	MSG_PING byte = 0xff
	// ...
)

type Message struct {
	Type byte   // Message type
	Data []byte // Raw data

	Reply chan Message // Reply channel
}

// Test if two messages are equal.
func (m Message) EqualWith(temp Message) bool {
	if m.Type != temp.Type {
		return false
	}

	if !bytes.Equal(StripBytes(m.Data, 0), StripBytes(temp.Data, 0)) {
		return false
	}

	return true
}

// Serialize message into bytes.
func (m Message) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write([]byte{m.Type})
	buf.Write(m.Data)

	return buf.Bytes(), nil
}

// Read message from bytes.
func (m *Message) UnmarshalBinary(data []byte) error {

	if len(data) <= 1 {
		return errors.New("Invalid Message.")
	}

	buf := bytes.NewBuffer(data)

	m.Type = buf.Next(1)[0]
	m.Data = buf.Bytes()

	return nil
}
