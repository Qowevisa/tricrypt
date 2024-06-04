package monkeylang

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"fmt"
	"git.qowevisa.me/Qowevisa/gotell/tlep/monkeylang/myerr"
	"git.qowevisa.me/Qowevisa/gotell/tlep/shuffle"
	"strings"
)

const (
	VERSION_1 = 1 + iota
)

const LATEST_VERSION = VERSION_1

type WORD_VAL uint16

const DICT_LEN = 0x010000

type Dictionary struct {
	Values [DICT_LEN]WORD_VAL
	Words  [DICT_LEN]string
}

func CreateNewDictionary() (*Dictionary, error) {
	var d Dictionary
	dummyKey := make([]byte, DICT_LEN)

	n, err := rand.Read(dummyKey)
	if err != nil {
		return nil, err
	}
	if n != len(dummyKey) {
		return nil, myerr.INTERNAL_ERROR
	}
	for i := 0; i < DICT_LEN; i++ {
		d.Values[i] = WORD_VAL(i)
	}
	newVals := shuffle.Shuffle(d.Values[:], dummyKey)
	if len(newVals) != len(d.Values) {
		return nil, myerr.INTERNAL_ERROR
	}
	for i, nv := range newVals {
		d.Values[i] = nv
	}
	words := GenerateStrongWords(DICT_LEN)
	for i, word := range words {
		d.Words[i] = word
	}

	return &d, nil
}

func (d *Dictionary) GetFirstWords(n int) []string {
	return d.Words[:n]
}

func (d *Dictionary) GetFirstValues(n int) []WORD_VAL {
	return d.Values[:n]
}

func allUniqueT[T comparable](strings []T) (int, bool) {
	notUnique := 0
	seen := make(map[T]struct{})
	for _, s := range strings {
		if _, ok := seen[s]; ok {
			notUnique++
		}
		seen[s] = struct{}{}
	}
	return notUnique, notUnique == 0
}

func (d *Dictionary) GetStat() string {
	stt := fmt.Sprintf("D Words Len : %d\n", len(d.Values))
	mNU, mIAU := allUniqueT(d.Words[:])
	stt += fmt.Sprintf("D Meanings Not Unique : %d ; %t\n", mNU, mIAU)
	wNU, wIAU := allUniqueT(d.Values[:])
	stt += fmt.Sprintf("D Words Not Unique : %d ; %t\n", wNU, wIAU)
	return stt
}

func (d *Dictionary) GetFingerprint() string {
	mm := strings.Join(d.Words[:], ".")
	var ww []byte
	for _, w := range d.Values {
		ww = append(ww, byte(w>>8), byte(w))
	}
	ww = append(ww, []byte(mm)...)
	hash := sha512.Sum512(ww)
	fingerprint := make([]byte, base32.StdEncoding.EncodedLen(len(hash)))
	base32.StdEncoding.Encode(fingerprint, hash[:])
	return string(fingerprint[:])
}

func (d *Dictionary) GetFingerprintBytes() []byte {
	mm := strings.Join(d.Words[:], ".")
	var ww []byte
	for _, w := range d.Values {
		ww = append(ww, byte(w>>8), byte(w))
	}
	ww = append(ww, []byte(mm)...)
	hash := sha512.Sum512(ww)
	fingerprint := make([]byte, base32.StdEncoding.EncodedLen(len(hash)))
	base32.StdEncoding.Encode(fingerprint, hash[:])
	return fingerprint[:]
}

func (d *Dictionary) GetFingerprintWithInfo(info string) string {
	mm := strings.Join(d.Words[:], ".")
	var ww []byte
	for _, w := range d.Values {
		ww = append(ww, byte(w>>8), byte(w))
	}
	ww = append(ww, []byte(mm)...)
	ww = append(ww, []byte(info)...)
	hash := sha512.Sum512(ww)
	fingerprint := make([]byte, base32.StdEncoding.EncodedLen(len(hash)))
	base32.StdEncoding.Encode(fingerprint, hash[:])
	return string(fingerprint[:])
}
