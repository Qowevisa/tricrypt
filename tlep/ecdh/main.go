package ecdh

import (
	"crypto/ecdh"
	"crypto/rand"
	"errors"
)

var Curve = ecdh.P521()

var ERROR_SharedNotComputed = errors.New("Shared secret is not computed!")

type ECDHConnection struct {
	privateKey     *ecdh.PrivateKey
	otherPublicKey *ecdh.PublicKey
	sharedSecret   []byte
}

func CreateNewConnection() (*ECDHConnection, error) {
	privKey, err := generatePrivKey()
	if err != nil {
		return nil, err
	}
	return &ECDHConnection{
		privateKey: privKey,
	}, nil
}

func (c *ECDHConnection) GetMyPublicKeyBytes() []byte {
	return c.privateKey.PublicKey().Bytes()
}

func (c *ECDHConnection) AcceptOtherPubKeyBytes(pubBytes []byte) error {
	pubKey, err := Curve.NewPublicKey(pubBytes)
	if err != nil {
		return err
	}
	c.otherPublicKey = pubKey
	shared, err := c.privateKey.ECDH(pubKey)
	if err != nil {
		return err
	}
	c.sharedSecret = shared
	return nil
}

func (c *ECDHConnection) GetShared() ([]byte, error) {
	if len(c.sharedSecret) == 0 {
		return nil, ERROR_SharedNotComputed
	}
	return c.sharedSecret, nil
}

func generatePrivKey() (*ecdh.PrivateKey, error) {
	privateKey, err := Curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return privateKey, err
}

func AreKeysEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
