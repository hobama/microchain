package core

import (
	"bytes"
	"errors"
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
func TestPingMsgMarshalJson(t *testing.T) {
	p1 := PingData{PublicKey: GenRandomBytes(32)}

	p1json, err := p1.MarshalJson()
	if err != nil {
		panic(errors.New("(PingData) MarshalJson() testing failed"))
	}

	var p2 PingData

	err = p2.UnmarshalJson(p1json)
	if err != nil {
		panic(errors.New("(*PingData) UnmarshalJson() testing failed"))
	}

	if !bytes.Equal(p1.PublicKey, p2.PublicKey) {
		panic(errors.New("(PingData) MarshalJson()/UnmarshalJson() testing failed"))
	}
}

// Test Message marshal functions.
func TestMessageMarshalJson(t *testing.T) {
	m1 := GenRandomMessage()

	m1json, err := m1.MarshalJson()
	if err != nil {
		panic(errors.New("(Message) MarshalJson() testing failed"))
	}

	var m2 Message

	err = m2.UnmarshalJson(m1json)
	if err != nil {
		panic(errors.New("(*Message) UnmarshalJson() testing failed"))
	}

	if !m1.EqualWith(m2) {
		panic(errors.New("(Message) MarshalJson()/UnmarshalJson() testing failed"))
	}
}
