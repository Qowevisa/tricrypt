package chaos

import (
	"bytes"
	"encoding/gob"
	"errors"
)

const (
	LATEST_VERSION = 1
)

var (
	invalidVersion = errors.New("Version is incopatible")
)

func (c *ChaosSystem) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(c)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetFromBytes(data []byte) (*ChaosSystem, error) {
	var cs ChaosSystem
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&cs)
	if err != nil {
		return nil, err
	}
	if cs.Version != LATEST_VERSION {
		return nil, invalidVersion
	}

	return &cs, nil
}
