package hkdf

import (
	"crypto/sha256"

	"golang.org/x/crypto/hkdf"
)

func DeriveAESKeyFromLongKeyAndInfo(shared, info []byte) ([]byte, error) {
	// Create reader for HKDF
	hash := sha256.New
	salt := []byte("EIZBq3CdxfeaGxZ2Zj7QIIhExgbhkdhDW4ePrDheEaEFmzRYdJqrYnddAGk5pqWq")
	// Contextual information
	hkdf := hkdf.New(hash, shared, salt, info)

	// Generate and return the key
	aesKey := make([]byte, 32)
	if _, err := hkdf.Read(aesKey); err != nil {
		return nil, err
	}
	return aesKey, nil
}
