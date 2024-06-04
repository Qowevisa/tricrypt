package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var tooShort = errors.New("Too Short")

type Message struct {
	Data []byte
}

// NOTE: should use only 32 byte key
// key can be derived via hkdf package
func Encrypt(plaintext, key []byte) (*Message, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return &Message{
		Data: append(nonce, ciphertext...),
	}, nil
}

// NOTE: should use only 32 byte key
// key can be derived via hkdf package
func Decrypt(ciphertext, key []byte) (*Message, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, tooShort
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return &Message{
		Data: plaintext,
	}, nil
}
