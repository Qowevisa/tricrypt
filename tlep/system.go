package tlep

import (
	"errors"
	"fmt"
	"log"

	"git.qowevisa.me/Qowevisa/gotell/tlep/chaos"
	"git.qowevisa.me/Qowevisa/gotell/tlep/ecdh"
	"git.qowevisa.me/Qowevisa/gotell/tlep/encrypt"
	"git.qowevisa.me/Qowevisa/gotell/tlep/gmyerr"
	"git.qowevisa.me/Qowevisa/gotell/tlep/hkdf"
	"git.qowevisa.me/Qowevisa/gotell/tlep/monkeylang"
	"git.qowevisa.me/Qowevisa/gotell/tlep/shuffle"
)

type TLEPLevel uint8

var (
	TLEP_LEVEL_NO_CONNECTION  TLEPLevel = 0
	TLEP_LEVEL_ECDH           TLEPLevel = 1
	TLEP_LEVEL_ECDH_CBES      TLEPLevel = 2
	TLEP_LEVEL_ECDH_CBES_MKLG TLEPLevel = 3
)

var (
	ERROR_UNHANDLED_TLEP_LEVEL = errors.New("Unhandled TLEP LEVEL")
)

// Three Layer Encryption Protocol schema
type TLEP struct {
	// Security Layer Level
	SLLevel TLEPLevel
	Name    string
	// Elliptic-Curve Diffie-Hellman
	ECDHConnection *ecdh.ECDHConnection
	// Chaos-Based Encryption System
	CBES *chaos.ChaosSystem
	// MonKeyLanG Dictionary
	MKLGDict *monkeylang.Dictionary
	// Debug for logging
	Debug bool
}

func InitTLEP(name string) (*TLEP, error) {
	var t TLEP
	t.Name = name
	conn, err := ecdh.CreateNewConnection()
	if err != nil {
		return nil, gmyerr.WrapPrefix("ecdh.CreateNewConnection", err)
	}
	t.ECDHConnection = conn
	cbes := chaos.CreateNewChaosSystem()
	t.CBES = cbes
	t.SLLevel = TLEP_LEVEL_NO_CONNECTION
	if t.Debug {
		log.Printf("TLEP initiated")
	}
	return &t, nil
}

func (t *TLEP) ECDHGetPublicKey() ([]byte, error) {
	if t.ECDHConnection == nil {
		return nil, gmyerr.WrapPrefix("t.ECDHConnection", IS_NIL)
	}
	if t.Debug {
		log.Printf("TLEP uses GetMyPublicKeyBytes")
	}
	ar := t.ECDHConnection.GetMyPublicKeyBytes()
	return ar, nil
}

func (t *TLEP) ECDHApplyOtherKeyBytes(otherKey []byte) error {
	if t.ECDHConnection == nil {
		return gmyerr.WrapPrefix("t.ECDHConnection", IS_NIL)
	}
	err := t.ECDHConnection.AcceptOtherPubKeyBytes(otherKey)
	if err != nil {
		return gmyerr.WrapPrefix("t.ECDHConnection.AcceptOtherPubKeyBytes", err)
	}
	t.SLLevel = TLEP_LEVEL_ECDH
	if t.Debug {
		log.Printf("TLEP is upgraded to ECDH level")
	}
	return nil
}

func (t *TLEP) CBESInitRandom() error {
	if t.CBES == nil {
		return gmyerr.WrapPrefix("t.CBES", IS_NIL)
	}
	err := t.CBES.InitRandom()
	if err != nil {
		return gmyerr.WrapPrefix("t.CBES.InitRandom", IS_NIL)
	}
	t.SLLevel = TLEP_LEVEL_ECDH_CBES
	if t.Debug {
		log.Printf("TLEP is upgraded to ECDH_CBES level")
	}
	return nil
}

func (t *TLEP) CBESGetBytes() ([]byte, error) {
	if t.CBES == nil {
		return nil, gmyerr.WrapPrefix("t.CBES", IS_NIL)
	}
	return t.CBES.Bytes()
}

func (t *TLEP) CBESSetFromBytes(bytes []byte) error {
	newCBES, err := chaos.GetFromBytes(bytes)
	if err != nil {
		return gmyerr.WrapPrefix("chaos.GetFromBytes", err)
	}
	t.CBES = newCBES
	t.SLLevel = TLEP_LEVEL_ECDH_CBES
	if t.Debug {
		log.Printf("TLEP is upgraded to ECDH_CBES level")
	}
	return nil
}

func (t *TLEP) CBESGetPassword(passLen uint) ([]byte, error) {
	if t.CBES == nil {
		return nil, gmyerr.WrapPrefix("t.CBES", IS_NIL)
	}
	return t.CBES.GetPassword(passLen), nil
}

func (t *TLEP) EncryptMessageEA(message []byte) ([]byte, error) {
	if t.ECDHConnection == nil {
		return nil, gmyerr.WrapPrefix("t.ECDHConnection", IS_NIL)
	}
	shared, err := t.ECDHConnection.GetShared()
	if errors.Is(err, ecdh.ERROR_SharedNotComputed) {
		return nil, ecdh.ERROR_SharedNotComputed
	}
	aesKey, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(shared, []byte("ECDH_BASE"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	msg, err := encrypt.Encrypt(message, aesKey)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Encrypt", err)
	}
	return msg.Data, nil
}

func (t *TLEP) DecryptMessageEA(message []byte) ([]byte, error) {
	if t.ECDHConnection == nil {
		return nil, gmyerr.WrapPrefix("t.ECDHConnection", IS_NIL)
	}
	shared, err := t.ECDHConnection.GetShared()
	if errors.Is(err, ecdh.ERROR_SharedNotComputed) {
		return nil, ecdh.ERROR_SharedNotComputed
	}
	aesKey, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(shared, []byte("ECDH_BASE"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	msg, err := encrypt.Decrypt(message, aesKey)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Decrypt", err)
	}
	return msg.Data, nil
}

func (t *TLEP) CanIUseEA() (bool, error) {
	if t.ECDHConnection == nil {
		return false, gmyerr.WrapPrefix("t.ECDHConnection", IS_NIL)
	}
	_, err := t.ECDHConnection.GetShared()
	if errors.Is(err, ecdh.ERROR_SharedNotComputed) {
		return false, ecdh.ERROR_SharedNotComputed
	}
	return t.SLLevel > TLEP_LEVEL_ECDH, nil
}

func (t *TLEP) EncryptMessageCAFEA(message []byte) ([]byte, error) {
	// First Layer Encryption
	if t.ECDHConnection == nil {
		return nil, gmyerr.WrapPrefix("t.ECDHConnection", IS_NIL)
	}
	shared, err := t.ECDHConnection.GetShared()
	if errors.Is(err, ecdh.ERROR_SharedNotComputed) {
		return nil, ecdh.ERROR_SharedNotComputed
	}
	aesKey, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(shared, []byte("CAFEA_ENCRYPTION_ECDH"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// First Layer Encryption
	msg, err := encrypt.Encrypt(message, aesKey)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Encrypt", err)
	}
	// Second Layer Encryption
	cbesLongKey, err := t.CBESGetPassword(512)
	if err != nil {
		return nil, gmyerr.WrapPrefix("t.CBESGetPassword", err)
	}
	aesKey2, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(cbesLongKey, []byte("CAFEA_ENCRYPTION_CBES"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// Second Layer Encryption
	msg2, err := encrypt.Encrypt(msg.Data, aesKey2)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Encrypt", err)
	}
	return msg2.Data, nil
}

func (t *TLEP) DecryptMessageCAFEA(message []byte) ([]byte, error) {
	// Second Layer Decryption
	cbesLongKey, err := t.CBESGetPassword(512)
	if err != nil {
		return nil, gmyerr.WrapPrefix("t.CBESGetPassword", err)
	}
	aesKey2, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(cbesLongKey, []byte("CAFEA_ENCRYPTION_CBES"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// Second Layer Decryption
	msg2, err := encrypt.Decrypt(message, aesKey2)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Decrypt", err)
	}
	// First Layer Decryption
	if t.ECDHConnection == nil {
		return nil, gmyerr.WrapPrefix("t.ECDHConnection", IS_NIL)
	}
	shared, err := t.ECDHConnection.GetShared()
	if errors.Is(err, ecdh.ERROR_SharedNotComputed) {
		return nil, ecdh.ERROR_SharedNotComputed
	}
	aesKey, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(shared, []byte("CAFEA_ENCRYPTION_ECDH"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// First Layer Decryption
	msg, err := encrypt.Decrypt(msg2.Data, aesKey)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Decrypt", err)
	}
	return msg.Data, nil
}

// Monkeylang ECDH Shuffle CHaos-based AES-256-GCM
func (t *TLEP) EncryptMessageMESCHA(message []byte) ([]byte, error) {
	// First Layer Key
	if t.ECDHConnection == nil {
		return nil, gmyerr.WrapPrefix("t.ECDHConnection", IS_NIL)
	}
	shared, err := t.ECDHConnection.GetShared()
	if errors.Is(err, ecdh.ERROR_SharedNotComputed) {
		return nil, ecdh.ERROR_SharedNotComputed
	}
	// First Layer Encryption
	aesKey, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(shared, []byte("MESCHA_ENCRYPTION_ECDH"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// First Layer Encryption
	msg, err := encrypt.Encrypt(message, aesKey)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Encrypt", err)
	}
	// First Layer Shuffle
	shfKey, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(shared, []byte("MESCHA_SHUFFLE_ECDH"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// First Layer Shuffle
	shuffledMsg := shuffle.Shuffle(msg.Data, shfKey)
	// Second Layer Key
	cbesLongKey, err := t.CBESGetPassword(512)
	if err != nil {
		return nil, gmyerr.WrapPrefix("t.CBESGetPassword", err)
	}
	// Second Layer Encryption
	aesKey2, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(cbesLongKey, []byte("MESCHA_ENCRYPTION_CBES"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// Second Layer Encryption
	msg2, err := encrypt.Encrypt(shuffledMsg, aesKey2)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Encrypt", err)
	}
	// Second Layer Shuffle
	shfKey2, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(cbesLongKey, []byte("MESCHA_SHUFFLE_CBES"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// Second Layer Shuffle
	shuffledMsg2 := shuffle.Shuffle(msg2.Data, shfKey2)
	// Third Layer Key
	mklgPrint := t.MKLGDict.GetFingerprintBytes()
	// Third Layer Encryption
	aesKey3, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(mklgPrint, []byte("MESCHA_ENCRYPTION_MKLG"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// Third Layer Encryption
	msg3, err := encrypt.Encrypt(shuffledMsg2, aesKey3)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Encrypt", err)
	}
	// Third Layer Shuffle
	shfKey3, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(mklgPrint, []byte("MESCHA_SHUFFLE_MKLG"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// Third Layer Shuffle
	shuffledMsg3 := shuffle.Shuffle(msg3.Data, shfKey3)
	return shuffledMsg3, nil
}

// Monkeylang ECDH Shuffle CHaos-based AES-256-GCM
func (t *TLEP) DecryptMessageMESCHA(message []byte) ([]byte, error) {
	// Third Layer Key
	mklgPrint := t.MKLGDict.GetFingerprintBytes()
	// Third Layer Shuffle
	shfKey3, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(mklgPrint, []byte("MESCHA_SHUFFLE_MKLG"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// Third Layer Shuffle
	unshuffledMsg3 := shuffle.Unshuffle(message, shfKey3)
	// Third Layer Decryption
	aesKey3, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(mklgPrint, []byte("MESCHA_ENCRYPTION_MKLG"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// Third Layer Decryption
	msg3, err := encrypt.Decrypt(unshuffledMsg3, aesKey3)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Decrypt", err)
	}
	// Second Layer Key
	cbesLongKey, err := t.CBESGetPassword(512)
	if err != nil {
		return nil, gmyerr.WrapPrefix("t.CBESGetPassword", err)
	}
	// Second Layer Shuffle
	shfKey2, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(cbesLongKey, []byte("MESCHA_SHUFFLE_CBES"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// Second Layer Shuffle
	unshuffledMsg2 := shuffle.Unshuffle(msg3.Data, shfKey2)
	// Second Layer Decryption
	aesKey2, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(cbesLongKey, []byte("MESCHA_ENCRYPTION_CBES"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// Second Layer Decryption
	msg2, err := encrypt.Decrypt(unshuffledMsg2, aesKey2)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Decrypt", err)
	}
	// First Layer Key
	if t.ECDHConnection == nil {
		return nil, gmyerr.WrapPrefix("t.ECDHConnection", IS_NIL)
	}
	shared, err := t.ECDHConnection.GetShared()
	if errors.Is(err, ecdh.ERROR_SharedNotComputed) {
		return nil, ecdh.ERROR_SharedNotComputed
	}
	// First Layer Shuffle
	shfKey, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(shared, []byte("MESCHA_SHUFFLE_ECDH"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// First Layer Shuffle
	unshuffledMsg := shuffle.Unshuffle(msg2.Data, shfKey)
	// First Layer Decryption
	aesKey, err := hkdf.DeriveAESKeyFromLongKeyAndInfo(shared, []byte("MESCHA_ENCRYPTION_ECDH"))
	if err != nil {
		return nil, gmyerr.WrapPrefix("hkdf.DeriveAESKeyFromLongKeyAndInfo", err)
	}
	// First Layer Decryption
	msg, err := encrypt.Decrypt(unshuffledMsg, aesKey)
	if err != nil {
		return nil, gmyerr.WrapPrefix("encrypt.Decrypt", err)
	}
	return msg.Data, nil
}

func (t *TLEP) EncryptMessageAtMax(msg []byte) ([]byte, error) {
	switch t.SLLevel {
	case TLEP_LEVEL_ECDH:
		if t.Debug {
			log.Printf("Encrypting using EA")
		}
		return t.EncryptMessageEA(msg)
	case TLEP_LEVEL_ECDH_CBES:
		if t.Debug {
			log.Printf("Encrypting using CAFEA")
		}
		return t.EncryptMessageCAFEA(msg)
	case TLEP_LEVEL_ECDH_CBES_MKLG:
		if t.Debug {
			log.Printf("Encrypting using MESCHA")
		}
		return t.EncryptMessageMESCHA(msg)
	}
	return nil, gmyerr.WrapPrefix(fmt.Sprintf("TLEP: %d", t.SLLevel), ERROR_UNHANDLED_TLEP_LEVEL)
}

func (t *TLEP) DecryptMessageAtMax(msg []byte) ([]byte, error) {
	switch t.SLLevel {
	case TLEP_LEVEL_ECDH:
		if t.Debug {
			log.Printf("Decrypting using EA")
		}
		return t.DecryptMessageEA(msg)
	case TLEP_LEVEL_ECDH_CBES:
		if t.Debug {
			log.Printf("Decrypting using CAFEA")
		}
		return t.DecryptMessageCAFEA(msg)
	case TLEP_LEVEL_ECDH_CBES_MKLG:
		if t.Debug {
			log.Printf("Decrypting using MESCHA")
		}
		return t.DecryptMessageMESCHA(msg)
	}
	return nil, gmyerr.WrapPrefix(fmt.Sprintf("TLEP: %d", t.SLLevel), ERROR_UNHANDLED_TLEP_LEVEL)
}
