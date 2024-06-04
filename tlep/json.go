package tlep

import (
	"encoding/json"
	"fmt"
	"git.qowevisa.me/Qowevisa/gotell/tlep/gmyerr"
	"io"
	"os"
)

const (
	TLEPDirName = "tlep"
)

func (t *TLEP) getFileName() string {
	return fmt.Sprintf("%s/%s.tlep", TLEPDirName, t.Name)
}

// don't care for errors for now
func canOpenDir(dirname string) bool {
	dir, err := os.Open(dirname)
	if err != nil {
		// TODO check for errors
		return false
	}
	stat, err := dir.Stat()
	if err != nil {
		// TODO check for errors
		return false
	}
	return stat.IsDir()
}

func (t *TLEP) SaveToFile() error {
	canI := canOpenDir(TLEPDirName)
	if !canI {
		err := os.Mkdir(TLEPDirName, 0755)
		if err != nil {
			return gmyerr.WrapPrefix("os.Mkdir", err)
		}
	}
	s, err := json.Marshal(t)
	if err != nil {
		return gmyerr.WrapPrefix("json.Marshal", err)
	}
	file, err := os.Create(t.getFileName())
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

func LoadFromFileByName(name string) (*TLEP, error) {
	var t TLEP
	t.Name = name
	file, err := os.Open(t.getFileName())
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
	err = json.Unmarshal(result, &t)
	if err != nil {
		return nil, gmyerr.WrapPrefix("json.Unmarshal", err)
	}
	return &t, nil
}
