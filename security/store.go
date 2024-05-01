package security

import (
	"crypto/rand"
	"fmt"

	"git.qowevisa.me/Qowevisa/gotell/errors"
	"git.qowevisa.me/Qowevisa/gotell/types"
	"golang.org/x/crypto/curve25519"
)

func generateKeyPair(cfg types.ECDHConfig) ([]byte, []byte, error) {
	var private [32]byte
	if _, err := rand.Read(private[:]); err != nil {
		return nil, nil, errors.WrapErr("rand.Read", err)
	}

	public, err := curve25519.X25519(private[:], curve25519.Basepoint)
	if err != nil {
		return nil, nil, errors.WrapErr("curve25519.X25519", err)
	}

	return private[:], public, nil
}

type KeyPair struct {
	baseName   string
	privateKey []byte
	publicKey  []byte
}

type Store struct {
	Pairs map[string]KeyPair
}

func InitStorage() (*Store, error) {
	err := initFS()
	if err != nil {
		return nil, errors.WrapErr("initFS", err)
	}
	var newStore Store
	newStore.Pairs = make(map[string]KeyPair)

	return &newStore, nil
}

func (s *Store) AddNewPair(name string, ecdhCfg types.ECDHConfig) error {
	_, exists := s.Pairs[name]
	if exists {
		return errors.WrapErr(fmt.Sprintf("Store.Pairs[%s]", name), errors.ALREADY_SET)
	}
	fileExists, err := fs_checkIfPairExists(name)
	if fileExists {
		return errors.WrapErr(fmt.Sprintf("Store.Pairs[%s]. FS found but in store", name), errors.NOT_SET)
	}
	if err != nil {
		return errors.WrapErr("fs_checkIfPairExists", err)
	}
	private, public, err := generateKeyPair(ecdhCfg)
	if err != nil {
		return errors.WrapErr("generateKeyPair", err)
	}
	s.Pairs[name] = KeyPair{
		baseName:   name,
		privateKey: private,
		publicKey:  public,
	}
	return nil
}
