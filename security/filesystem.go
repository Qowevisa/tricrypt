package security

import (
	"os"
	"path/filepath"

	"git.qowevisa.me/Qowevisa/gotell/errors"
)

func findInStore(name string) KeyPair {
	return KeyPair{}
}

const (
	_int_FSPath          = "./key-store"
	_int_FSPairsPath     = _int_FSPath + "/pairs"
	_int_PrivateKeyAddon = "private.Key"
	_int_PublicKeyAddon  = "public.key"
)

func initFS() error {
	var err error
	err = os.Mkdir(_int_FSPath, os.ModePerm)
	if err != nil {
		return errors.WrapErr("Mkdir:initFSPath", err)
	}

	err = os.Mkdir(_int_FSPairsPath, os.ModePerm)
	if err != nil {
		return errors.WrapErr("Mkdir:initFSPairsPath", err)
	}
	return nil
}

func fs_checkIfPairExists(name string) (bool, error) {
	var err error
	privateKeyPath := filepath.Join(_int_FSPairsPath, name+_int_PrivateKeyAddon)
	_, err = os.Stat(privateKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, errors.WrapErr(privateKeyPath, err)
		}
	}
	publicKeyPath := filepath.Join(_int_FSPairsPath, name+_int_PublicKeyAddon)
	_, err = os.Stat(publicKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, errors.WrapErr(publicKeyPath, err)
		}
	}
	return true, nil
}

// TODO: add encryption
func fs_savePair(keyPair KeyPair) error {
	var err error
	err = os.WriteFile(keyPair.baseName+_int_PrivateKeyAddon, keyPair.privateKey, 0644)
	if err != nil {
		return errors.WrapErr("WriteFile:privateKey", err)
	}
	err = os.WriteFile(keyPair.baseName+_int_PublicKeyAddon, keyPair.publicKey, 0644)
	if err != nil {
		return errors.WrapErr("WriteFile:publicKey", err)
	}
	return nil
}

// TODO: add encryption
func fs_readPair(name string, keyPair *KeyPair) error {
	var err error
	privateKeyPath := filepath.Join(_int_FSPairsPath, name+_int_PrivateKeyAddon)
	_, err = os.Stat(privateKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.WrapErr(privateKeyPath, errors.NOT_FOUND)
		} else {
			return errors.WrapErr(privateKeyPath, err)
		}
	}
	privKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return errors.WrapErr("os.ReadFile:private", err)
	}
	publicKeyPath := filepath.Join(_int_FSPairsPath, name+_int_PublicKeyAddon)
	_, err = os.Stat(publicKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.WrapErr(publicKeyPath, errors.NOT_FOUND)
		} else {
			return errors.WrapErr(publicKeyPath, err)
		}
	}
	pubKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return errors.WrapErr("os.ReadFile:publicKey", err)
	}
	keyPair.baseName = name
	keyPair.privateKey = privKey
	keyPair.publicKey = pubKey
	return nil
}
