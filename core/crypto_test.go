package core

import (
	"testing"
)

func TestKeyGen(t *testing.T) {
	for i := 0; i < 5; i++ {
		keyPair, err := NewECDSAKeyPair()
		if err != nil {
			t.Error(err)
		}
		if len(keyPair.Public) != 64 || len(keyPair.Private) != 32 {
			t.Errorf("%d", len(keyPair.Public))
			t.Errorf("%d", len(keyPair.Private))
			t.Error("Invalid key pair")
		}
	}
}

func TestKeySigningAndVerify(t *testing.T) {
	for i := 0; i < 500; i++ {
		keyPair, err := NewECDSAKeyPair()
		if err != nil {
			t.Error(err)
		}

		data := []byte("test" + string(i))
		hash := SHA256(data)

		signature, err := keyPair.Sign(hash)
		if err != nil {
			t.Error(err)
		} else if !VerifySignature(keyPair.Public, signature, hash) {
			t.Error("Invalid signature")
		}
	}
}
