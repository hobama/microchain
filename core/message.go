package core

import (
	"bytes"
)

var (
	MessageOptionsBufferSize = 4
)

type Message struct {
	Identifier byte
	Options    []byte
	Data       []byte

	Reply chan Message
}

// Test if two messages are equal.
func (m Message) EqualWith(temp Message) bool {
	if m.Identifier != temp.Identifier {
		return false
	}

	if !bytes.Equal(StripBytes(m.Options, 0), StripBytes(temp.Options, 0)) {
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

	buf.Write([]byte{m.Identifier})
	buf.Write(FitBytesIntoSpecificWidth(m.Options, MessageOptionsBufferSize))
	buf.Write(m.Data)

	return buf.Bytes(), nil
}

// Read message from bytes.
func (m *Message) UnmarshalBinary(data []byte) error {
	return nil
}
