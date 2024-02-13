package message

import (
	"encoding/json"
	"errors"
)

type Msg struct {
	Type MsgType `json:"type"`
	Msg  Interface
}

type Interface interface {
	GetMsgID() uint64
	GetMsgType() MsgType
	GetPayload() []byte
	// GetReceiverID returns receiverID and if it needs to get id from db
	GetReceiverID() ([]string, bool)
}

type Receiver int
type MsgType string

const (
	Person    MsgType = "person"
	Group     MsgType = "group"
	Broadcast MsgType = "broadcast"
	Ack       MsgType = "ack"
	Auth      MsgType = "auth"
)

type ContentType int

const (
	TEXT = 1
	FILE = 2
)

func IfMsgTypeAllowed(t MsgType) bool {
	if t == Auth || t == Ack || t == Broadcast || t == Group || t == Person {
		return true
	}
	return false
}

func Decode(data []byte) (*Msg, error) {
	var msg Msg
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	if !IfMsgTypeAllowed(msg.Type) {
		return nil, errors.New("invalid msg type")
	}

	return &msg, nil
}
