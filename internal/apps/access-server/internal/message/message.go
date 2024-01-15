package message

import (
	"encoding/json"
	"errors"
)

type Receiver int
type MsgType string

const (
	Person    MsgType = "person"
	Group     MsgType = "group"
	Broadcast MsgType = "broadcast"
	Ack       MsgType = "ack"
	Auth      MsgType = "auth"
)

type ConntentType int

const (
	TEXT = 1
	FILE = 2
)

type Message struct {
	MessageID   int64
	MsgType     MsgType
	ContentType ConntentType
	ReceiverID  []string
	GroupID     []string
	Subscribers []string
	Payload     []byte
	AuthMessage

	// -------- queue --------
	private    int64
	index      int
	retryCount int
}

type AuthMessage struct {
	Token string
}

// TODO add a sequence id such as snowflake to keep the seq id orderly

func MessageAck() {

}

func MessageSync() {

}

func MessageBroadcast() {

}

func MessageAuth(data []byte) (*Message, error) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	if msg.MsgType != Auth {
		return nil, errors.New("the msg is not auth")
	}
	return &msg, err
}
