package core

import (
	"errors"
	"testing"
)

// Generate random message
func GenRandomMessage() Message {
	m, _ := NewMessage()

	m.Type = GenRandomBytes(1)[0]
	m.SourceID = GenRandomBytes(8)
	m.Data = GenRandomBytes(8)
	return m
}

// Test Message marshal functions.
func TestMessageMarshalJson(t *testing.T) {
	m1 := GenRandomMessage()

	m1json, err := m1.MarshalJson()
	if err != nil {
		panic(errors.New("(Message) MarshalJson() testing failed."))
	}

	var m2 Message

	err = m2.UnmarshalJson(m1json)
	if err != nil {
		panic(errors.New("(*Message) UnmarshalJson() testing failed."))
	}

	if !m1.EqualWith(m2) {
		panic(errors.New("(Message) MarshalJson()/UnmarshalJson() testing failed."))
	}
}
