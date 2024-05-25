package communication

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"

	"git.qowevisa.me/Qowevisa/gotell/gmyerr"
)

// VERSION field
const (
	V1 = 1 + iota
)

// FROM field
const (
	FROM_SERVER = 1 + iota
	FROM_CLIENT
	FROM_MY_ID
)

// FROM_ID is 2 bytes

// TO_ID is 2 bytes

// ACTION field
const (
	ACTION_ASK = 1 + iota
	ACTION_SEND
)

// ABOUT field
const (
	ABOUT_NICKNAME = 1 + iota
	ABOUT_ID
	ABOUT_LINK
	ABOUT_ECDH_PUB_KEY
	ABOUT_CBES_SPECS
	ABOUT_MKLG_FGPRINT
	ABOUT_MESSAGE
)

// DATA_LEN is 2 bytes

type Message struct {
	Version uint8
	From    uint8
	FromID  uint16
	ToID    uint16
	Action  uint8
	About   uint8
	DataLen uint16
	Data    []byte
}

func (c *Message) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(c)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(serverBytes []byte) (*Message, error) {
	var c Message
	buf := bytes.NewBuffer(serverBytes)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func ServerAskClientAboutNickname() ([]byte, error) {
	c := Message{
		Version: V1,
		From:    FROM_SERVER,
		FromID:  0,
		ToID:    0,
		Action:  ACTION_ASK,
		About:   ABOUT_NICKNAME,
		DataLen: 0,
		Data:    []byte{},
	}
	return c.Bytes()
}

func ClientSendServerNickname(nickname []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		From:    FROM_CLIENT,
		FromID:  0,
		ToID:    0,
		Action:  ACTION_SEND,
		About:   ABOUT_NICKNAME,
		DataLen: uint16(len(nickname)),
		Data:    nickname,
	}
	return c.Bytes()
}

func ServerSendClientHisID(id []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		From:    FROM_SERVER,
		FromID:  0,
		ToID:    0,
		Action:  ACTION_SEND,
		About:   ABOUT_ID,
		DataLen: uint16(len(id)),
		Data:    id,
	}
	return c.Bytes()
}

func (r *RegisteredUser) GenerateLink(count uint32) (Link, error) {
	var l Link
	buf := make([]byte, LINK_LEN_V1)

	_, err := rand.Read(buf)
	if err != nil {
		return Link{}, err
	}
	if count == 0 {
		return Link{}, ERROR_LINK_ZERO_COUNT
	}
	l.Status = LINK_STATUS_CREATED
	l.Data = buf
	l.UseCount = count

	return l, nil
}

func (l *Link) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(l)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *RegisteredUser) GetIDFromLink(l Link) ([]byte, error) {
	bb, err := l.Bytes()
	if err != nil {
		return nil, gmyerr.WrapPrefix("l.Bytes", err)
	}
	c := Message{
		Version: V1,
		From:    FROM_CLIENT,
		FromID:  r.ID,
		ToID:    0,
		Action:  ACTION_ASK,
		About:   ABOUT_LINK,
		DataLen: uint16(len(bb)),
		Data:    bb,
	}
	return c.Bytes()
}

func ServerSendClientIDFromLink(toID uint16, toName []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		From:    FROM_SERVER,
		FromID:  0,
		ToID:    toID,
		Action:  ACTION_SEND,
		About:   ABOUT_LINK,
		DataLen: uint16(len(toName)),
		Data:    toName,
	}
	return c.Bytes()
}

func (r *RegisteredUser) ClientSendThroughServerCBESSpecs(to uint16, data []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		From:    FROM_CLIENT,
		FromID:  r.ID,
		ToID:    to,
		Action:  ACTION_SEND,
		About:   ABOUT_CBES_SPECS,
		DataLen: uint16(len(data)),
		Data:    data,
	}
	return c.Bytes()
}

func (r *RegisteredUser) SendMessageToID(to uint16, msg []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		From:    FROM_CLIENT,
		FromID:  r.ID,
		ToID:    to,
		Action:  ACTION_SEND,
		About:   ABOUT_MESSAGE,
		DataLen: uint16(len(msg)),
		Data:    msg,
	}
	return c.Bytes()
}
