package monkeylang

import (
	"encoding/json"
	"fmt"
	"git.qowevisa.me/Qowevisa/gotell/tlep/gmyerr"
	"io"
	"os"
)

const (
	DictsDirName = "mklgs"
)

const extension = ".dict.mklg"

func getFileName(prefix string) string {
	return fmt.Sprintf("%s/%s%s",
		DictsDirName, prefix, extension,
	)
}

func SaveToFile(d Dictionary, prefix string) error {
	s, err := json.Marshal(d)
	if err != nil {
		return gmyerr.WrapPrefix("json.Marshal", err)
	}
	file, err := os.Create(getFileName(prefix))
	defer file.Close()
	if err != nil {
		return gmyerr.WrapPrefix("os.Create", err)
	}
	_, err = file.Write(s)
	if err != nil {
		return gmyerr.WrapPrefix("file.Write", err)
	}
	return nil
}

func DoesDictExists(prefix string) (bool, error) {
	file, err := os.Open(getFileName(prefix))
	defer file.Close()
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return true, gmyerr.WrapPrefix("os.Open", err)
	}
	return true, nil
}

func LoadFromFile(prefix string) (*Dictionary, error) {
	var d Dictionary
	file, err := os.Open(getFileName(prefix))
	defer file.Close()
	if err != nil {
		return nil, gmyerr.WrapPrefix("os.Create", err)
	}
	var result []byte
	buf := make([]byte, 10240)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, gmyerr.WrapPrefix("file.Write", err)
		}
		result = append(result, buf[:n]...)
	}
	err = json.Unmarshal(result, &d)
	if err != nil {
		return nil, gmyerr.WrapPrefix("json.Unmarshal", err)
	}
	return &d, nil
}
