package communication

import (
	"bytes"
	"encoding/gob"
)

const (
	SERVER_COMMAND = 1 + iota
	SERVER_MESSAGE
	CLIENT_RESPONSE
	P2P_MESSAGE
)

const (
	V1 = 1 + iota
)

const (
	NICKNAME = 1
)

type communicationMessage struct {
	Type uint8
	Data []byte
}

func AskClientNickname() ([]byte, error) {
	c := communicationMessage{
		Type: SERVER_COMMAND,
		Data: []byte{NICKNAME},
	}
	return c.Bytes()
}

func JustGetMessage(msg []byte) ([]byte, error) {
	c := communicationMessage{
		Type: SERVER_MESSAGE,
		Data: msg,
	}
	return c.Bytes()
}

func (c *communicationMessage) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(c)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(serverBytes []byte) (*communicationMessage, error) {
	var c communicationMessage
	buf := bytes.NewBuffer(serverBytes)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
