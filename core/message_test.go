package core

import (
	"errors"
	"testing"
)

func TestMessageEQ(t *testing.T) {
	m1 := Message{Type: 0x01, Data: []byte{0x01, 0x02, 0x03}}
	m2 := Message{Type: 0x01, Data: []byte{0x01, 0x02, 0x03}}

	if !m2.EqualWith(m1) {
		panic(errors.New("(Message) EqualWith() testing failed."))
	}
}

func TestMessageMarshal(t *testing.T) {
	m1 := Message{Type: 0x01, Data: []byte{0x00, 0x01, 0x02, 0x03}}

	mBytes, err := m1.MarshalBinary()
	if err != nil {
		panic(errors.New("(Message) MarshalBinary() testing failed."))
	}

	m2 := new(Message)

	err = m2.UnmarshalBinary(mBytes)
	if err != nil {
		panic(err)
	}

	if !m2.EqualWith(m1) {
		panic(errors.New("(*Message) UnmarshalBinary() testing failed."))
	}
}
