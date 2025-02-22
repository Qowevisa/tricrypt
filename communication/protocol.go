package communication

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/gob"
	"errors"

	"git.qowevisa.me/Qowevisa/gotell/gmyerr"
)

// VERSION is 1 byte
const (
	V1 = 1 + iota
)

// ID is 1 byte
const (
	// Client Handles
	ID_SERVER_ASK_CLIENT_NICKNAME = 1
	// Server Handles
	ID_CLIENT_SEND_SERVER_NICKNAME = 2
	// Client Handles
	ID_SERVER_APPROVE_CLIENT_NICKNAME = 3
	ID_SERVER_DECLINE_CLIENT_NICKNAME = 4
	// Server Handles
	ID_CLIENT_SEND_SERVER_LINK = 5
	// Client Handles
	ID_SERVER_APPROVE_CLIENT_LINK = 6
	ID_SERVER_DECLINE_CLIENT_LINK = 7
	// Server Handles
	ID_CLIENT_ASK_SERVER_LINK = 8
	// Client Handles
	ID_SERVER_SEND_CLIENT_ANOTHER_ID = 9
	// Client Handles . Server redirects
	ID_CLIENT_ASK_CLIENT_HANDSHAKE            = 10
	ID_CLIENT_APPROVE_CLIENT_HANDSHAKE        = 11
	ID_CLIENT_DECLINE_CLIENT_HANDSHAKE        = 12
	ID_CLIENT_SEND_CLIENT_ECDH_PUBKEY         = 13
	ID_CLIENT_SEND_CLIENT_CBES_SPECS          = 14
	ID_CLIENT_SEND_CLIENT_MKLG_FINGERPRINT    = 15
	ID_CLIENT_APPROVE_CLIENT_MKLG_FINGERPRINT = 16
	ID_CLIENT_DECLINE_CLIENT_MKLG_FINGERPRINT = 17
	ID_CLIENT_SEND_CLIENT_MESSAGE             = 18
	// SAVED SPACE FOR OTHER DATA
	ID_INTERCOM_SIGNAL_1  = 128
	ID_INTERCOM_SIGNAL_2  = 129
	ID_INTERCOM_SIGNAL_3  = 130
	ID_INTERCOM_SIGNAL_4  = 131
	ID_INTERCOM_SIGNAL_5  = 132
	ID_INTERCOM_SIGNAL_6  = 133
	ID_INTERCOM_SIGNAL_7  = 134
	ID_INTERCOM_SIGNAL_8  = 135
	ID_INTERCOM_SIGNAL_9  = 136
	ID_INTERCOM_SIGNAL_10 = 137
)

// FROM_ID is 2 bytes

// TO_ID is 2 bytes

// DATA_LEN is 2 bytes

type Message struct {
	Version uint8
	ID      uint8
	FromID  uint16
	ToID    uint16
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
		ID:      ID_SERVER_ASK_CLIENT_NICKNAME,
		FromID:  0,
		ToID:    0,
		DataLen: 0,
		Data:    []byte{},
	}
	return c.Bytes()
}

func ClientSendServerNickname(nickname []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		ID:      ID_CLIENT_SEND_SERVER_NICKNAME,
		FromID:  0,
		ToID:    0,
		DataLen: uint16(len(nickname)),
		Data:    nickname,
	}
	return c.Bytes()
}

func ServerSendClientHisID(id []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		ID:      ID_SERVER_APPROVE_CLIENT_NICKNAME,
		FromID:  0,
		ToID:    0,
		DataLen: uint16(len(id)),
		Data:    id,
	}
	return c.Bytes()
}

func ServerSendClientDecline() ([]byte, error) {
	c := Message{
		Version: V1,
		ID:      ID_SERVER_DECLINE_CLIENT_NICKNAME,
		FromID:  0,
		ToID:    0,
		DataLen: 0,
		Data:    []byte{},
	}
	return c.Bytes()
}

func (r *RegisteredUser) GenerateLink(count uint16) (Link, error) {
	if count == 0 {
		return Link{}, ERROR_LINK_ZERO_COUNT
	}
	var l Link
	buf := make([]byte, LINK_LEN_V1)

	_, err := rand.Read(buf)
	if err != nil {
		return Link{}, err
	}
	encoded := base32.StdEncoding.EncodeToString(buf)
	l.Status = LINK_STATUS_CREATED
	l.Data = []byte(encoded)
	l.UseCount = count

	return l, nil
}

func IsThisALinkData(data string) (bool, error) {
	dst := make([]byte, len([]byte(data)))
	n, err := base32.StdEncoding.Decode(dst, []byte(data))
	if err != nil {
		return false, err
	}
	if n != LINK_LEN_V1 {
		return false, errors.New("Link len is not standard")
	}
	return true, nil
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

func DecodeLink(data []byte) (*Link, error) {
	var l Link
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func ClientSendServerLink(from uint16, l Link) ([]byte, error) {
	bb, err := l.Bytes()
	if err != nil {
		return nil, err
	}
	c := Message{
		Version: V1,
		ID:      ID_CLIENT_SEND_SERVER_LINK,
		FromID:  from,
		ToID:    0,
		DataLen: uint16(len(bb)),
		Data:    bb,
	}
	return c.Bytes()
}

func ServerApproveClientLink() ([]byte, error) {
	c := Message{
		Version: V1,
		ID:      ID_SERVER_APPROVE_CLIENT_LINK,
		FromID:  0,
		ToID:    0,
		DataLen: 0,
		Data:    []byte{},
	}
	return c.Bytes()
}

func ServerDeclineClientLink() ([]byte, error) {
	c := Message{
		Version: V1,
		ID:      ID_SERVER_DECLINE_CLIENT_LINK,
		FromID:  0,
		ToID:    0,
		DataLen: 0,
		Data:    []byte{},
	}
	return c.Bytes()
}

func (r *RegisteredUser) GetIDFromLink(l Link) ([]byte, error) {
	bb, err := l.Bytes()
	if err != nil {
		return nil, gmyerr.WrapPrefix("l.Bytes", err)
	}
	c := Message{
		Version: V1,
		ID:      ID_CLIENT_ASK_SERVER_LINK,
		FromID:  r.ID,
		ToID:    0,
		DataLen: uint16(len(bb)),
		Data:    bb,
	}
	return c.Bytes()
}

func ServerSendClientIDFromLink(toID uint16, toName []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		FromID:  0,
		ID:      ID_SERVER_SEND_CLIENT_ANOTHER_ID,
		ToID:    toID,
		DataLen: uint16(len(toName)),
		Data:    toName,
	}
	return c.Bytes()
}

func (r *RegisteredUser) ClientSendThroughServerECDHPubKey(to uint16, pubkey []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		ID:      ID_CLIENT_SEND_CLIENT_ECDH_PUBKEY,
		FromID:  r.ID,
		ToID:    to,
		DataLen: uint16(len(pubkey)),
		Data:    pubkey,
	}
	return c.Bytes()
}

func (r *RegisteredUser) ClientSendThroughServerCBESSpecs(to uint16, data []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		ID:      ID_CLIENT_SEND_CLIENT_CBES_SPECS,
		FromID:  r.ID,
		ToID:    to,
		DataLen: uint16(len(data)),
		Data:    data,
	}
	return c.Bytes()
}

func (r *RegisteredUser) ClientSendThroughServerMKLGPrint(to uint16, mklg []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		ID:      ID_CLIENT_SEND_CLIENT_MKLG_FINGERPRINT,
		FromID:  r.ID,
		ToID:    to,
		DataLen: uint16(len(mklg)),
		Data:    mklg,
	}
	return c.Bytes()
}

func (r *RegisteredUser) SendMessageToID(to uint16, msg []byte) ([]byte, error) {
	c := Message{
		Version: V1,
		ID:      ID_CLIENT_SEND_CLIENT_MESSAGE,
		FromID:  r.ID,
		ToID:    to,
		DataLen: uint16(len(msg)),
		Data:    msg,
	}
	return c.Bytes()
}
