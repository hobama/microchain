package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

type KeyPair struct {
	Public  []byte `json: "public"`
	Private []byte `json: "private"`
}

// Generate ECDSA key pair
// Public key format:
// | x ... 32 bytes | y ... 32 bytes |
// Private key format:
// | private key |
func NewECDSAKeyPair() (*KeyPair, error) {
	pubKeyCurve := elliptic.P256()

	privateKey := new(ecdsa.PrivateKey)
	privateKey, err := ecdsa.GenerateKey(pubKeyCurve, rand.Reader)
	if err != nil {
		return nil, err
	}

	xy := privateKey.PublicKey
	xBytes, err := BigIntToBytes(xy.X, 32)
	if err != nil {
		return nil, err
	}

	yBytes, err := BigIntToBytes(xy.Y, 32)
	if err != nil {
		return nil, err
	}

	publicKeyBytes := append(xBytes, yBytes...)

	privateKeyBytes, err := BigIntToBytes(privateKey.D, 32)
	if err != nil {
		return nil, err
	}

	return &KeyPair{Public: publicKeyBytes, Private: privateKeyBytes}, nil
}

// Sign a hash of a file/message
func (kp *KeyPair) Sign(hash []byte) ([]byte, error) {
	// Decode private key
	privateKey := new(big.Int)
	privateKey.SetBytes(kp.Private)

	// Decode public key
	publicKey := kp.Public
	xBytes, yBytes := publicKey[:32], publicKey[32:]
	x := BytesToBigInt(xBytes)
	y := BytesToBigInt(yBytes)

	key := ecdsa.PrivateKey{ecdsa.PublicKey{elliptic.P256(), x, y}, privateKey}

	r, s, err := ecdsa.Sign(rand.Reader, &key, hash)
	if err != nil {
		return nil, err
	}

	rBytes, err := BigIntToBytes(r, 32)
	if err != nil {
		return nil, err
	}

	sBytes, err := BigIntToBytes(s, 32)
	if err != nil {
		return nil, err
	}

	rs := append(rBytes, sBytes...)

	return []byte(rs), nil
}

// Verify signature
func VerifySignature(publicKey, signature, hash []byte) bool {
	xBytes, yBytes := publicKey[:32], publicKey[32:]
	x := BytesToBigInt(xBytes)
	y := BytesToBigInt(yBytes)

	rBytes, sBytes := signature[:32], signature[32:]
	r := BytesToBigInt(rBytes)
	s := BytesToBigInt(sBytes)

	pub := ecdsa.PublicKey{elliptic.P256(), x, y}

	return ecdsa.Verify(&pub, hash, r, s)
}
